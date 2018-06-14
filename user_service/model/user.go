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
	Pwd              string `xorm:"VARCHAR(255)"`
	Phone            string `xorm:"unique VARCHAR(64)"`
	PhoneVerifyTime  int    `xorm:"comment('手机验证时间') INT(11)"`
	Email            string `xorm:"unique VARCHAR(128)"`
	EmailVerifyTime  int    `xorm:"INT(11)"`
	GoogleVerifyId   string `xorm:"VARCHAR(128)"`
	GoogleVerifyTime int    `xorm:"INT(255)"`
}

/*
type ArticlesStruct struct {
	ID             int32
	Description    string //重要 、一般
	Title          string
	CreateDateTime string
}

type ArticlesDetailStruct struct {
	ID            int32
	Title         string
	Description   string
	Content       string
	Covers        []byte
	ContentImages []byte
	Type          int32
	TypeName      string
	Author        string
	Weight        int32
	Shares        int32
	Hits          int32
	Comments      int32
	DisplayMark   bool
	CreateTime    string
	UpdateTime    string
	AdminID       int32
	AdminNickname string
}
*/

func (s *User) RegisterByPhone(req *proto.RegisterPhoneRequest) int32 {
	if ret := s.CheckUserExist(req.Phone, "phone"); ret != ERRCODE_SUCCESS {
		return ret
	}

	e := &User{
		Pwd:   req.Pwd,
		Phone: req.Phone,
	}
	_, err := DB.GetMysqlConn().Cols("pwd", "phone").Insert(e)
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

func (s *User) RegisterByEmail(req *proto.RegisterEmailRequest) int32 {
	if ret := s.CheckUserExist(req.Email, "email"); ret != ERRCODE_SUCCESS {
		return ret
	}

	e := &User{
		Pwd:   req.Pwd,
		Email: req.Email,
	}
	_, err := DB.GetMysqlConn().Cols("pwd", "email").Insert(e)
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

func (s *User) GetUserExByPhone(phone string) (u *UserEx, ret int32) {
	u = &UserEx{}
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

func (s *User) ModifyPwd(phone string, pwd string) (ret int32) {
	return 0
}
