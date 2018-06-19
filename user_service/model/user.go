package model

import (
	"digicon/common/check"
	"digicon/common/google"
	"digicon/common/random"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	. "digicon/user_service/log"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/jsonpb"
	"strconv"
	"time"
)

type User struct {
	Uid              int    `xorm:"not null pk autoincr INT(11)"`
	Account          string `xorm:"unique VARCHAR(128)"`
	Pwd              string `xorm:"VARCHAR(255)"`
	Phone            string `xorm:"unique VARCHAR(64)"`
	PhoneVerifyTime  int    `xorm:"comment('手机验证时间') INT(11)"`
	Email            string `xorm:"unique VARCHAR(128)"`
	EmailVerifyTime  int    `xorm:"INT(11)"`
	GoogleVerifyId   string `xorm:"VARCHAR(128)"`
	GoogleVerifyTime int    `xorm:"INT(255)"`
	SmsTip           bool   `xorm:"INT(4)"`
	PayPwd           string `xorm:"comment('支付密码') VARCHAR(255)"`
	NeedPwd          bool   `xorm:"INT(4)"`
	NeedPwdTime      int    `xorm:"INT(11)"`
}

func (s *User) GetUser(uid int32) (ret int32, err error) {
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
			Uid:            int32(s.Uid),
			Account:        s.Account,
			Phone:          s.Phone,
			Email:          s.Email,
			SmsTip:         s.SmsTip,
			GoogleVerifyId: google_switch,
			PaySwitch:      pay_switch,
			NeedPwd:        s.NeedPwd,
			NeedPwdTime:    int32(s.NeedPwdTime),
			LoginPwdLevel:  pwd_level,
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

	/*
		b,err:=json.Marshal(r)
		if err != nil {
			Log.Errorln(err.Error())
		}
		data=string(b)
	*/
	return
}

//刷新用户缓存
func (s *User) RefreshCache(uid int32) (out *proto.UserAllData, ret int32, err error) {
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
	if ret := s.CheckUserExist(req.Ukey, filed); ret != ERRCODE_SUCCESS {
		return ret
	}

	e := &User{
		Pwd:     req.Pwd,
		Phone:   req.Ukey,
		Account: req.Ukey,
	}
	_, err := DB.GetMysqlConn().Cols("pwd", filed, "account").Insert(e)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	sql := fmt.Sprintf("%s=?", filed)
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

		d.Invites += 1
		_, err = DB.GetMysqlConn().Where("uid=?", d.Uid).Cols("invites").Update(d)
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
	ok, err := DB.GetMysqlConn().Where(sql, param).Get(&User{})
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
func (s *User) LoginByPhone(phone, pwd string) int32 {
	if ok := check.CheckPhone(phone); !ok {
		return ERRCODE_SMS_PHONE_FORMAT
	}
	m := &User{}
	ok, err := DB.GetMysqlConn().Where("phone=?", phone).Get(m)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	if ok {
		if m.Pwd == pwd {
			return ERRCODE_SUCCESS
		}
		return ERRCODE_PWD
	}
	return ERRCODE_ACCOUNT_NOTEXIST
}

//通过邮箱登陆
func (s *User) LoginByEmail(eamil, pwd string) int32 {
	if ok := check.CheckEmail(eamil); !ok {
		return ERRCODE_SMS_EMAIL_FORMAT
	}
	m := &User{}
	ok, err := DB.GetMysqlConn().Where("email=?", eamil).Get(m)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	if ok {
		if m.Pwd == pwd {
			return ERRCODE_SUCCESS
		}
		return ERRCODE_PWD
	}
	return ERRCODE_ACCOUNT_NOTEXIST
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
func (s *User) SetGoogleSecertKey(uid int32, secert_key string) (ret int32) {
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
