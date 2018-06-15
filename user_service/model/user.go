package model

import (
	"digicon/common/check"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	. "digicon/user_service/log"
	"fmt"
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
}


func GetUser(uid int32) *User {
	u := &User{}
	DB.GetMysqlConn().Where("uid=?", uid).Get(u)
	return u
}
//通过手机注册
func (s *User) RegisterByPhone(req *proto.RegisterPhoneRequest) int32 {
	if ret := s.CheckUserExist(req.Phone, "phone"); ret != ERRCODE_SUCCESS {
		return ret
	}

	e := &User{
		Pwd:     req.Pwd,
		Phone:   req.Phone,
		Account: req.Phone,
	}
	_, err := DB.GetMysqlConn().Cols("pwd", "phone", "account").Insert(e)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	_, err = DB.GetMysqlConn().Where("phone=?", req.Phone).Get(e)
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

//检查用户注册过没
func (s *User) CheckUserExist(param string, col string) int32 {
	sql := fmt.Sprintf("%s=?", col)
	ok, err := DB.GetMysqlConn().Where(sql, param).Get(&User{})
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	if ok {
		return ERRCODE_ACCOUNT_EXIST
	}
	return ERRCODE_SUCCESS
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

//根据手机查询用户
func (s *User) GetUserByPhone(phone string) (u *User, ret int32) {
	u = &User{}
	ok, err := DB.GetMysqlConn().Where("phone=?", phone).Get(u)
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


//修改密码
func (s *User) ModifyPwd(phone string, pwd string) (ret int32) {
	return 0
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
