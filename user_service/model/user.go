package model

import (
	"digicon/common/check"
	"digicon/common/encryption"
	"digicon/common/google"
	"digicon/common/random"
	proto "digicon/proto/rpc"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	. "digicon/common/constant"
	. "digicon/proto/common"
	. "digicon/user_service/dao"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/jsonpb"
	"github.com/pkg/errors"
)

type User struct {
	Uid              uint64 `xorm:"not null pk autoincr comment('用户ID') BIGINT(11)"`
	Account          string `xorm:"comment('账号') unique VARCHAR(64)"`
	Pwd              string `xorm:"comment('密码') VARCHAR(255)"`
	Country          string `xorm:"comment('地区号') VARCHAR(32)"`
	Phone            string `xorm:"comment('手机') unique VARCHAR(64)"`
	PhoneVerifyTime  int    `xorm:"comment('手机验证时间') INT(11)"`
	Email            string `xorm:"comment('邮箱') unique VARCHAR(128)"`
	EmailVerifyTime  int    `xorm:"comment('邮箱验证时间') INT(11)"`
	GoogleVerifyId   string `xorm:"comment('谷歌私钥') VARCHAR(128)"`
	GoogleVerifyTime int    `xorm:"comment('谷歌验证时间') INT(255)"`
	SmsTip           bool   `xorm:"default 0 comment('短信提醒') TINYINT(1)"`
	PayPwd           string `xorm:"comment('支付密码') VARCHAR(255)"`
	NeedPwd          bool   `xorm:"comment('免密设置1开启0关闭') TINYINT(1)"`
	NeedPwdTime      int    `xorm:"comment('免密周期') INT(11)"`
	Status           int    `xorm:"default 0 comment('用户状态，1正常，2冻结') INT(11)"`
	SecurityAuth     int    `xorm:"comment('认证状态1110') TINYINT(8)"`
	SetTardeMark     int    `xorm:"comment('资金密码设置状态标识') INT(8)"`
	WhiteList        int    `xorm:"not null default 2 comment('用户白名单 1为白名单 免除交易手续费，2 需要缴纳交易手续费') TINYINT(4)"`
}

type UidAndInviteID struct {
	Uid      uint64
	InviteId uint64
}

func (*User) TableName() string {
	return "user"
}


func (s *User) GetUser(uid uint64) (ret int32, err error) {
	ok, err := DB.GetMysqlConn().Where("uid=?", uid).Get(s)
	if err != nil {
		log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}

	if ok {
		ret = ERRCODE_SUCCESS
		return
	}
	ret = ERRCODE_ACCOUNT_NOTEXIST
	return
}

//根据手机查询用户
func (s *User) GetUserByPhone(phone string) (ret int32, err error) {
	ok, err := DB.GetMysqlConn().Where("phone=?", phone).Get(s)
	if err != nil {
		log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}

	if ok {
		ret = ERRCODE_SUCCESS
		return
	}
	ret = ERRCODE_ACCOUNT_NOTEXIST
	return
}

//根据邮箱查询用户
func (s *User) GetUserByEmail(email string) (ret int32, err error) {
	ok, err := DB.GetMysqlConn().Where("email=?", email).Get(s)
	if err != nil {
		log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}

	if ok {
		ret = ERRCODE_SUCCESS
		return
	}
	ret = ERRCODE_ACCOUNT_NOTEXIST
	return
}

//根据邀请码属性获取用户
func (s *User) GetUserByInviteCode(inviteCode string) (ret int32, err error) {
	ok, err := DB.GetMysqlConn().
		Alias("u").
		Select("u.*").
		Join("INNER", []string{new(UserEx).TableName(), "ue"}, "ue.uid=u.uid").
		Where("ue.invite_code=?", inviteCode).
		Get(s)

	if err != nil {
		log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}

	if ok {
		ret = ERRCODE_SUCCESS
		return
	}
	ret = ERRCODE_ACCOUNT_NOTEXIST
	return
}

// 获取注册未获得奖励
func (s *User) GetRegisterNoRewardUser() ([]*UidAndInviteID, error) {
	var uidAndInviteID []*UidAndInviteID
	err := DB.GetMysqlConn().SQL(fmt.Sprintf("SELECT uid,invite_id FROM %s"+
		" WHERE uid NOT IN (SELECT uid FROM %s WHERE type=%d)"+
		" ORDER BY register_time ASC", new(UserEx).TableName(), new(TokenFrozen).TableName(), proto.TOKEN_TYPE_OPERATOR_HISTORY_REGISTER)).
		Limit(1000).
		Find(&uidAndInviteID)
	if err != nil {
		return nil, err
	}

	return uidAndInviteID, nil
}

//序列化用户基础数据
func (s *User) SerialJsonData() (data string, err error) {
	var (
		google_switch bool
		pay_switch    bool
		pwd_level     int32
	)
	if s.GoogleVerifyId != "" {
		google_switch = true
	}
	if s.PayPwd != "" {
		pay_switch = true
	}
	pwd_level = 1

	ex := &UserEx{}
	ret, err := ex.GetUserEx(s.Uid)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	if ret != ERRCODE_SUCCESS {
		log.Errorln("db err when find user_ex,uid=?", s.Uid)
		return
	}
	var first int32
	var second int32
	//if s.SecurityAuth^AUTH_TWO == AUTH_TWO {
	//	mark = 1
	//} else {
	//	mark = 0
	//}

	for {
		//有没有通过实名认证
		if s.SecurityAuth & AUTH_FIRST ==AUTH_FIRST{
			first = FIRST_ALREADY //已经通过认证
			break
		}
		//是否有申请认证
		if s.SetTardeMark & APPLY_FOR_FIRST == APPLY_FOR_FIRST{
			first = FIRST_VERIFYING
			break
		}
		//如果没有提交申请 检查是否有申请过认证但是没有通过
		if s.SetTardeMark & APPLY_FOR_FIRST_NOT_ALREADY == APPLY_FOR_FIRST_NOT_ALREADY {
			first = FIRST_NOT_ALREADY
			break
		}
		first =FIRST_NOT_V
		break
	}

	for{
		//有无通过二级认证
		if s.SecurityAuth & AUTH_TWO ==AUTH_TWO{
			second = SECOND_ALREADY
			break
		}
		//是否有申请二级认证
		if s.SetTardeMark & APPLY_FOR_SECOND ==APPLY_FOR_SECOND{
			second = SECOND_VERIFYING
			break
		}
		//没有通过二级认证
		if s.SetTardeMark & APPLY_FOR_SECOND_NOT_ALREADY ==APPLY_FOR_SECOND_NOT_ALREADY{
			second = SECOND_NOT_ALREADY
			break
		}
		second = SECOND_NOT_V
		break
	}

	r := &proto.UserAllData{
		Base: &proto.UserBaseData{
			Uid:            s.Uid,
			Account:        s.Account,
			Phone:          s.Phone,
			Email:          s.Email,
			SmsTip:         s.SmsTip,
			GoogleVerifyId: google_switch,
			PaySwitch:      pay_switch,
			NeedPwd:        s.NeedPwd,
			NeedPwdTime:    int32(s.NeedPwdTime),
			LoginPwdLevel:  pwd_level,
			Country:        s.Country,
			GoogleExist:    s.authSecurityCode(AUTH_GOOGLE),
			NickName:       ex.NickName,
			HeadSculpture:  ex.HeadSculpture,
			TradePwd:		s.authTradeCode(AUTH_TRADEMARK),
		},

		Real: &proto.UserRealData{
			RealName:        ex.RealName,
			IdentifyCard:    ex.IdentifyCard,
			CheckMarkFirst:  first,
			CheckMarkSecond: second,
		},

		Invite: &proto.UserInviteData{
			InviteCode: ex.InviteCode,
			Invites:    int32(ex.Invites),
		},
	}

	m := jsonpb.Marshaler{EmitDefaults: true}

	data, err = m.MarshalToString(r)
	if err != nil {
		log.Errorln(err.Error())
	}

	return
}

//强制刷新用户缓存
func (s *User) ForceRefreshCache(uid uint64) (out *proto.UserAllData, ret int32, err error) {
	var d string
	u := &User{}
	r := RedisOp{}
	ret, err = u.GetUser(uid)
	log.Info(uid)
	if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}
	if ret != ERRCODE_SUCCESS {
		return
	}

	d, err = u.SerialJsonData()
	if err != nil {

		ret = ERRCODE_UNKNOWN
		return
	}

	err = r.SetUserBaseInfo(uid, d)
	if err != nil {

		ret = ERRCODE_UNKNOWN
		return
	}

	return
}

//刷新用户缓存
func (s *User) RefreshCache(uid uint64) (out *proto.UserAllData, ret int32, err error) {
	r := RedisOp{}
	d, err := r.GetUserBaseInfo(uid)

	if err == redis.Nil { //找不到缓存记录则取mysql数据
		u := &User{}
		ret, err = u.GetUser(uid)
		if err != nil {
			ret = ERRCODE_UNKNOWN
			return
		}
		if ret != ERRCODE_SUCCESS {
			return
		}

		d, err = u.SerialJsonData()
		if err != nil {
			ret = ERRCODE_UNKNOWN
			return
		}

		err = r.SetUserBaseInfo(uid, d)
		if err != nil {
			ret = ERRCODE_UNKNOWN
			return
		}

	} else if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}

	out = &proto.UserAllData{}
	err = jsonpb.UnmarshalString(d, out)
	if err != nil {
		return
	}
	ret = ERRCODE_SUCCESS
	return
}

//通用注册
func (s *User) Register(req *proto.RegisterRequest, filed string) (errCode int32, uid uint64, referUid uint64) {
	var ret int32
	var err error
	var ok bool

	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"ukey":    req.Ukey,
				"type":    req.Type,
				"code":    req.Code,
				"invite":  req.InviteCode,
				"country": req.Country,
			}).Errorf("Register error %s", err.Error())
		}
	}()
	ret, err = s.CheckUserExist(req.Ukey, filed)
	if err != nil {
		return ERRCODE_UNKNOWN, 0, 0
	}

	if ret != ERRCODE_SUCCESS {
		return ret, 0, 0
	}

	pwd := encryption.GenMd5AndReverse(req.Pwd)

	d := &UserEx{} //主邀请人
	if req.InviteCode != "" {
		ok, err = DB.GetMysqlConn().Where("invite_code=?", req.InviteCode).Get(d)
		if err != nil {
			return ERRCODE_UNKNOWN, 0, 0
		}
		if !ok {
			return ERRCODE_INVITE, 0, 0
		}
	}
	var chmod int
	var sql string
	if filed == "phone" {
		chmod = AUTH_PHONE
		sql = fmt.Sprintf("INSERT INTO `user` (`account`,`pwd`,`country`,`%s`,`security_auth`) VALUES ('%s','%s','%s','%s','%d')", filed, req.Ukey, pwd, req.Country, req.Ukey, chmod)

	} else if filed == "email" {
		chmod = AUTH_EMAIL
		sql = fmt.Sprintf("INSERT INTO `user` (`account`,`pwd`,`%s`,`security_auth`) VALUES ('%s','%s','%s','%d')", filed, req.Ukey, pwd, req.Ukey, chmod)
	} else {
		log.Errorln("register error filed %s", filed)
	}

	//sql := fmt.Sprintf("INSERT INTO `user` (`account`,`pwd`,`country`,`%s`,`security_auth`) VALUES ('%s','%s','%s','%s',%d)", filed, req.Ukey, req.Pwd, req.Country, req.Ukey, chmod)
	_, err = DB.GetMysqlConn().Exec(sql)
	if err != nil {
		return ERRCODE_UNKNOWN, 0, 0
	}

	e := &User{}
	sql = fmt.Sprintf("%s=?", filed)
	_, err = DB.GetMysqlConn().Where(sql, req.Ukey).Get(e)
	if err != nil {
		return ERRCODE_UNKNOWN, 0, 0
	}

	code := random.Krand(6, random.KC_RAND_KIND_UPPER)
	str_code := string(code)

	if req.InviteCode != "" {
		m := &UserEx{
			Uid:          int64(e.Uid),
			RegisterTime: time.Now().Unix(),
			InviteCode:   str_code,
			InviteId:     d.Uid,
		}

		_, err = DB.GetMysqlConn().Insert(m)
		if err != nil {
			return ERRCODE_UNKNOWN, 0, 0
		}

		_, err = DB.GetMysqlConn().Where("uid=?", d.Uid).Cols("invites").Incr("invites", 1).Update(d)
		if err != nil {
			return ERRCODE_UNKNOWN, 0, 0
		}
	} else {
		m := &UserEx{
			Uid:          int64(e.Uid),
			RegisterTime: time.Now().Unix(),
			InviteCode:   str_code,
			NickName:     e.Account,                         //  注册的时候，昵称直接等与账户名
			HeadSculpture: random.SetRegisterRandHeader(),   //  注册时候，默认头像
		}

		_, err = DB.GetMysqlConn().Insert(m)
		if err != nil {
			return ERRCODE_UNKNOWN, 0, 0
		}
	}

	return ERRCODE_SUCCESS, e.Uid, uint64(d.Uid)

}

//检查用户注册过没
func (s *User) CheckUserExist(param string, col string) (ret int32, err error) {
	sql := fmt.Sprintf("%s=?", col)
	ok, err := DB.GetMysqlConn().Where(sql, param).Exist(&User{})
	if err != nil {
		log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}
	if ok {
		ret = ERRCODE_ACCOUNT_EXIST
		return
	}
	ret = ERRCODE_SUCCESS
	return
}

//通过手机登陆
func (s *User) LoginByPhone(phone, pwd string) (token string, ret int32) {
	ok, err := DB.GetMysqlConn().Where("phone=?", phone).Get(s)
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"phone": phone,
				"pwd":   pwd,
			}).Errorf("LoginByPhone error %s", err.Error())
		}
	}()
	if err != nil {
		log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}

	if ok {
		if s.Pwd == pwd {
			token, err = s.refreshToken()
			if err != nil {
				log.Errorln(err.Error())
				ret = ERRCODE_UNKNOWN
				return
			}
			ret = ERRCODE_SUCCESS
			return
		}
		ret = ERRCODE_PWD
		return
	}
	ret = ERRCODE_ACCOUNT_NOTEXIST
	return
}

//通过邮箱登陆
func (s *User) LoginByEmail(eamil, pwd string) (token string, ret int32) {
	if ok := check.CheckEmail(eamil); !ok {
		ret = ERRCODE_SMS_EMAIL_FORMAT
		return
	}

	ok, err := DB.GetMysqlConn().Where("email=?", eamil).Get(s)
	if err != nil {
		log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}
	if ok {
		if s.Pwd == pwd {
			token, err = s.refreshToken()
			if err != nil {
				log.Errorln(err.Error())
				ret = ERRCODE_UNKNOWN
				return
			}
			ret = ERRCODE_SUCCESS
			return
		}
		ret = ERRCODE_PWD
		return
	}
	ret = ERRCODE_ACCOUNT_NOTEXIST
	return
}

//更新token
func (s *User) refreshToken() (token string, err error) {
	uid_ := fmt.Sprintf("%d", s.Uid)
	salt := random.Krand(6, random.KC_RAND_KIND_NUM)
	b := encryption.Gensha256(uid_, time.Now().Unix(), string(salt))
	err = new(RedisOp).SetUserToken(string(b), s.Uid)
	if err != nil {
		return
	}

	token = string(b)
	return
}

//通用获取用户登陆信息
func (s *User) GetLoginUser(p *proto.LoginUserBaseData) (err error) {
	p.Uid = s.Uid

	return
}

//修改密码
func (s *User) ModifyPwd(newpwd string) (err error) {

	s.Pwd = encryption.GenMd5AndReverse(newpwd)
	_, err = DB.GetMysqlConn().Where("uid=?", s.Uid).Cols("pwd").Update(s)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

//修改谷歌验证密钥
func (s *User) SetGoogleSecertKey(uid uint64, secert_key string) (ret int32) {
	s.GoogleVerifyId = secert_key

	_, err := DB.GetMysqlConn().Where("uid=?", uid).Cols("google_verify_id").Update(s)
	if err != nil {
		log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}
	s.SecurityChmod(AUTH_GOOGLE)
	s.ForceRefreshCache(s.Uid)
	return ERRCODE_SUCCESS
}

//检查谷歌私钥存在
func (s *User) CheckGoogleExist() bool {
	if s.GoogleVerifyId == "" {
		return false
	}
	return true
}

//验证谷歌验证码
func (s *User) AuthGoogleCode(key string, input uint32) (ret int32, err error) {
	code, _ := google.GenGoogleCode(key)
	//code是16进制数据需要转成10进制
	g := strconv.Itoa(int(code))
	r, err := strconv.Atoi(g)
	if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}

	if input == uint32(r) {
		ret = ERRCODE_SUCCESS
		return
	}
	ret = ERRCODE_GOOGLE_CODE
	return
}

//解绑谷歌私钥
func (s *User) DelGoogleCode(input uint32) (ret int32, err error) {
	ret, err = s.AuthGoogleCode(s.GoogleVerifyId, input)
	if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}
	s.GoogleVerifyId = ""
	if ret == ERRCODE_SUCCESS {
		_, err = DB.GetMysqlConn().Where("uid=?", s.Uid).Cols("google_verify_id").Update(s)
		if err != nil {
			ret = ERRCODE_UNKNOWN
			log.Errorln(err.Error())
			return
		}
	}
	return
}

//获取验证类型
func (s *User) GetAuthMethod() int32 {
	if s.authSecurityCode(AUTH_GOOGLE) {
		return AUTH_GOOGLE
	} else if s.authSecurityCode(AUTH_PHONE) {
		return AUTH_PHONE
	} else if s.authSecurityCode(AUTH_EMAIL) {
		return AUTH_EMAIL
	}
	return 0
}

//获取排除谷歌验证类型
func (s *User) GetAuthMethodExpectGoogle() int32 {
	if s.authSecurityCode(AUTH_PHONE) {
		return AUTH_PHONE
	} else if s.authSecurityCode(AUTH_EMAIL) {
		return AUTH_EMAIL
	}
	return 0
}

//验证类型是否设置
func (s *User) authSecurityCode(code int) bool {
	g := s.SecurityAuth & code
	if g > 0 {
		return true
	}
	return false
}


//验证类型是否设置
func (s *User) authTradeCode(code int) bool {
	g := s.SetTardeMark & code
	if g > 0 {
		return true
	}
	return false
}

//自动判断验证方式
func (s *User) AuthCodeByAl(ukey, code string, ty int32, need bool) (ret int32, err error) {
	var m int32
	if !need {
		m = s.GetAuthMethod()
	} else {
		m = s.GetAuthMethodExpectGoogle()
	}
	switch m {
	case AUTH_EMAIL:
		return AuthEmail(ukey, ty, code)
	case AUTH_PHONE:
		return AuthSms(ukey, ty, code)
	case AUTH_GOOGLE:
		var code_ int
		code_, err = strconv.Atoi(code)
		if err != nil {
			log.Errorln(err.Error())
			return
		}
		return s.AuthGoogleCode(s.GoogleVerifyId, uint32(code_))
	default:
		break
	}

	return ERRCODE_UNKNOWN, errors.New("err auth method")
}

//验证通过修改权限
func (s *User) SecurityChmod(code int) (err error) {
	s.SecurityAuth = s.SecurityAuth ^ code
	_, err = DB.GetMysqlConn().Where("uid=?", s.Uid).Cols("security_auth").Update(s)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return nil
}

//删除认证权限
func (s *User) DelSecurityChmod(code int) (err error) {
	s.SecurityAuth = s.SecurityAuth &^ code
	_, err = DB.GetMysqlConn().Where("uid=?", s.Uid).Cols("security_auth").Update(s)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return nil
}

/*
	func: bind user email
*/
func (s *User) BindUserEmail(email string, uid uint64) (has bool, err error) {
	engine := DB.GetMysqlConn()
	s.Email = email
	eu := User{Email: email}
	has, err = engine.Where("account=? or email =?", email, email).Exist(&eu)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	if has {
		msg := "邮箱已经存在"
		log.Println(msg)
		return
	}
	_, err = engine.Where("uid=? ", uid).Update(s)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

/*
	func: bind user phone
*/
func (s *User) BindUserPhone(phone, country string, uid uint64) (has bool, err error) {
	engine := DB.GetMysqlConn()
	s.Phone = phone
	nu := User{Phone: phone, Country: country}
	has, err = engine.Where(" account=? or phone=?", phone, phone).Exist(&nu)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	if has {
		msg := "电话已经存在"
		log.Println(msg)
		return
	}
	_, err = engine.Where("uid=? ", uid).Update(s)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}


/*
	verify user pay_pwd
*/
func (s *User) VerifyPayPwd(uid uint64, paypwd string) (ret int32, err error) {
	ret, err = s.GetUser(uid)
	if err != nil {
		return ret, err
	}
	newpaypwd := encryption.GenMd5AndReverse(paypwd)
	//fmt.Println("newpaypwd:", newpaypwd, " spay: ", s.PayPwd)
	//log.Println("newpaypwd:", newpaypwd, " spay: ", s.PayPwd)
	if s.PayPwd == newpaypwd {
		return ERRCODE_SUCCESS, nil
	}else{
		return ERRCODE_PWD, nil
	}
}