package model

import (
	"digicon/common/check"
	"digicon/common/encryption"
	"digicon/common/google"
	"digicon/common/random"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	. "digicon/user_service/log"
	"fmt"
	"strconv"
	"time"

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
}


const (
	AUTH_EMAIL  = 2//00000010
	AUTH_PHONE  = 1//00000001
	AUTH_GOOGLE = 8//00001000
)

func (s *User) GetUser(uid uint64) (ret int32, err error) {
	ok, err := DB.GetMysqlConn().Where("uid=?", uid).Get(s)
	if err != nil {
		Log.Errorln(err.Error())
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
		Log.Errorln(err.Error())
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
		Log.Errorln(err.Error())
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
		Log.Errorln(err.Error())
		return
	}
	if ret != ERRCODE_SUCCESS {
		Log.Errorln("db err when find user_ex,uid=?", s.Uid)
		return
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
		},

		Real: &proto.UserRealData{
			RealName:     ex.RealName,
			IdentifyCard: ex.IdentifyCard,
		},

		Invite: &proto.UserInviteData{
			InviteCode: ex.InviteCode,
			Invites:    int32(ex.Invites),
		},
	}

	m := jsonpb.Marshaler{EmitDefaults: true}

	data, err = m.MarshalToString(r)
	if err != nil {
		Log.Errorln(err.Error())
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
	/*
	   	out = &proto.UserAllData{}

	   	err = json.Unmarshal([]byte(d), out)
	   	if err != nil {
	   		ret = ERRCODE_UNKNOWN
	   		return
	   	}
	   	godump.Dump("ll")
	   godump.Dump(out)
	*/
	out = &proto.UserAllData{}
	err = jsonpb.UnmarshalString(d, out)
	if err != nil {
		return
	}
	ret = ERRCODE_SUCCESS
	return
}

//通用注册
func (s *User) Register(req *proto.RegisterRequest, filed string) int32 {
	ret, err := s.CheckUserExist(req.Ukey, filed)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	if ret != ERRCODE_SUCCESS {
		return ret
	}

	/*
		e := &User{
			Pwd:     req.Pwd,
			Phone:   req.Ukey,
			Account: req.Ukey,
			Country: req.Country,
		}
		_, err = DB.GetMysqlConn().Cols("pwd", filed, "account", "country").Insert(e)
		if err != nil {
			Log.Errorln(err.Error())
			return ERRCODE_UNKNOWN
		}
	*/
	var chmod int
	if filed == "phone" {
		chmod = AUTH_PHONE
	} else if filed == "email" {
		chmod = AUTH_EMAIL
	} else {
		Log.Fatalf("register error filed %s", filed)
	}

	sql := fmt.Sprintf("INSERT INTO `user` (`account`,`pwd`,`country`,`%s`,`security`) VALUES ('%s','%s','%s','%s',%d)", filed, req.Ukey, req.Pwd, req.Country, req.Ukey, chmod)
	_, err = DB.GetMysqlConn().Exec(sql)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	e := &User{}
	sql = fmt.Sprintf("%s=?", filed)
	_, err = DB.GetMysqlConn().Where(sql, req.Ukey).Get(e)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	code := random.Krand(6, random.KC_RAND_KIND_UPPER)
	str_code := string(code)
	d := &UserEx{} //主邀请人
	if req.InviteCode != "" {
		ok, err := DB.GetMysqlConn().Where("invite_code=?", req.InviteCode).Get(d)
		if err != nil {
			Log.Errorln(err.Error())
			return ERRCODE_UNKNOWN
		}
		if !ok {
			return ERRCODE_UNKNOWN
		}
		m := &UserEx{
			Uid:          e.Uid,
			RegisterTime: time.Now().Unix(),
			InviteCode:   str_code,
			InviteId:     d.Uid,
		}

		_, err = DB.GetMysqlConn().Insert(m)
		if err != nil {
			Log.Errorln(err.Error())
			return ERRCODE_UNKNOWN
		}

		_, err = DB.GetMysqlConn().Where("uid=?", d.Uid).Cols("invites").Incr("invites", 1).Update(d)
		if err != nil {
			Log.Errorln(err.Error())
			return ERRCODE_UNKNOWN
		}
	} else {
		m := &UserEx{
			Uid:          e.Uid,
			RegisterTime: time.Now().Unix(),
			InviteCode:   str_code,
		}

		_, err = DB.GetMysqlConn().Insert(m)
		if err != nil {
			Log.Errorln(err.Error())
			return ERRCODE_UNKNOWN
		}
	}

	return ERRCODE_SUCCESS

}

/*
//通过手机注册
func (s *User) RegisterByPhone(req *proto.RegisterRequest) int32 {
	if ret := s.CheckUserExist(req.Ukey, "phone"); ret != ERRCODE_SUCCESS {
		return ret
	}

	e := &User{
		Pwd:     req.Pwd,
		Phone:   req.Ukey,
		Account: req.Ukey,
	}
	_, err := DB.GetMysqlConn().Cols("pwd", "phone", "account").Insert(e)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	_, err = DB.GetMysqlConn().Where("phone=?", req.Ukey).Get(e)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	m := &UserEx{
		Uid:          e.Uid,
		RegisterTime: time.Now().Unix(),
		InviteCode:   req.InviteCode,
	}

	_, err = DB.GetMysqlConn().Insert(m)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	return ERRCODE_SUCCESS
}

//通过邮箱注册


func (s *User) RegisterByEmail(req *proto.RegisterEmailRequest) int32 {
	if ret := s.CheckUserExist(req.Email, "email"); ret != ERRCODE_SUCCESS {
		return ret
	}

	e := &User{
		Pwd:     req.Pwd,
		Email:   req.Email,
		Account: req.Email,
	}
	_, err := DB.GetMysqlConn().Cols("pwd", "email", "account").Insert(e)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	_, err = DB.GetMysqlConn().Where("email=?", req.Email).Get(e)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	m := &UserEx{
		Uid:          e.Uid,
		RegisterTime: time.Now().Unix(),
		InviteCode:   req.InviteCode,
	}

	_, err = DB.GetMysqlConn().Insert(m)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	return ERRCODE_SUCCESS
}
*/

//检查用户注册过没
func (s *User) CheckUserExist(param string, col string) (ret int32, err error) {
	sql := fmt.Sprintf("%s=?", col)
	ok, err := DB.GetMysqlConn().Where(sql, param).Exist(&User{})
	if err != nil {
		Log.Errorln(err.Error())
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
	if err != nil {
		Log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}
	if ok {
		if s.Pwd == pwd {
			token, err = s.refreshToken()
			if err != nil {
				Log.Errorln(err.Error())
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
		Log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}
	if ok {
		if s.Pwd == pwd {
			token, err = s.refreshToken()
			if err != nil {
				Log.Errorln(err.Error())
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
	//_,err:=DB.GetMysqlConn().Where("uid=?",s.Uid).Cols("token").Update(s)
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
	s.Pwd = newpwd
	_, err = DB.GetMysqlConn().Where("uid=?", s.Uid).Cols("pwd").Update(s)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

//修改谷歌验证密钥
func (s *User) SetGoogleSecertKey(uid uint64, secert_key string) (ret int32) {
	s.GoogleVerifyId = secert_key

	_, err := DB.GetMysqlConn().Where("uid=?", uid).Cols("google_verify_id").Update(s)
	if err != nil {
		Log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}
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
		err = s.SecurityChmod(AUTH_GOOGLE)
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


func (s *User) authSecurityCode(code int) bool {
	g := s.SecurityAuth & code
	if g > 0 {
		return true
	}
	return false
}

//自动判断验证方式
func (s *User) AuthCodeByAl(ukey, code string,ty int32)(ret int32, err error) {
	m := s.GetAuthMethod()
	switch m {
	case AUTH_EMAIL:
		return AuthEmail(ukey,ty,code)
	case AUTH_PHONE:
		return AuthSms(ukey,ty,code)
	case AUTH_GOOGLE:
		var code_ int
		code_,err = strconv.Atoi(code)
		if err!=nil {
			Log.Errorln(err.Error())
			return
		}
		return s.AuthGoogleCode(s.GoogleVerifyId,uint32(code_))
	default:
		break
	}

	return ERRCODE_UNKNOWN,errors.New("err auth methon")
}

//验证通过修改权限
func (s *User) SecurityChmod(code int) (err error) {
	s.SecurityAuth=s.SecurityAuth^code
	_,err = DB.GetMysqlConn().Where("uid=?",s.Uid).Cols("security_auth").Update(s)
	if err!=nil {
		Log.Errorln(err.Error())
		return
	}
	return nil
}



/*
	func:
*/
//func (s *User) BindUserEmail()
