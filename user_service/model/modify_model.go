package model

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"digicon/user_service/dao"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func (s *User) ModifyLoginPwd(req *proto.UserModifyLoginPwdRequest) (result int32, err error) {
	fmt.Println("0..0.0.0.0.0.0.0.0.0.0.00.0.0.0.0000.0.0.0.")
	fmt.Println(req)
	if va := strings.Compare(req.ConfirmPwd, req.NewPwd); va != 0 {
		return ERRCODE_PWD_COMFIRM, nil
	}
	//modify DB
	engine := dao.DB.GetMysqlConn()
	ph := new(User)
	var ok bool
	ok, err = engine.ID(req.Uid).Get(ph)
	if err != nil {
		result = ERRCODE_UNKNOWN
		return
	}
	if !ok {
		result = ERRCODE_UNKNOWN
		err = errors.New("get phone number failed")
	}
	//旧密码的判断
	if b := strings.Compare(req.OldPwd, ph.Pwd); b != 0 {
		return ERRCODE_OLDPWD, nil
	}
	fmt.Println("lllllllllllllllllllllllllllllllllll")
	fmt.Println(ph.Phone)
	//验证短信
	//model.AuthSms()
	result, err = AuthSms(ph.Phone, SMS_MODIFY_LOGIN_PWD, req.Verify)
	if err != nil {
		return

	}
	if result != ERRCODE_SUCCESS {
		return
	}
	//token

	_, err = engine.ID(req.Uid).Update(&User{Pwd: req.NewPwd})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	s.RefreshCache(req.Uid)
	return ERRCODE_SUCCESS, nil

}

func (s *User) ModifyUserPhone1(req *proto.UserModifyPhoneRequest) (result int32, err error) {
	//验证短信
	engine := dao.DB.GetMysqlConn()
	ph := new(User)
	var ok bool
	ok, err = engine.ID(req.Uid).Get(ph)
	if err != nil {
		 return ERRCODE_UNKNOWN,err
	}
	if !ok {
		return ERRCODE_PHONE_NOT_EXIST,nil
	}
	fmt.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")
	fmt.Println("电话号码为：", ph.Phone, "验证码为：", req.Verify)
	result, err = AuthSms(ph.Phone, SMS_MODIFY_PHONE, req.Verify)
	if err != nil {
		return
	}
	if result != ERRCODE_SUCCESS {
		return
	}
	//token
	//旧的电话号码验证通过
	return 0, nil
}

func (s *User) ModifyUserPhone2(req *proto.UserSetNewPhoneRequest) (result int32, err error) {

	result, err = AuthSms(req.Phone, SMS_MODIFY_PHONE, req.Verify)
	if err != nil {
		return
	}
	if result != ERRCODE_SUCCESS {
		return ERRCODE_SMS_CODE_DIFF, err
	}
	//token
	//修改数据库字段
	engine := dao.DB.GetMysqlConn()
	u := new(User)
	has, err := engine.ID(req.Uid).Get(u)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if !has {
		return ERRCODE_ACCOUNT_NOTEXIST, nil
	}
	if strings.Compare(req.Phone, u.Phone) == 0 {
		return ERRCODE_PHONE_EXIST, nil
	}
	_, err = engine.ID(req.Uid).Update(&User{Phone: req.Phone})
	if err!=nil{
		return ERRCODE_UNKNOWN,err
	}
	s.RefreshCache(req.Uid)
	return ERRCODE_SUCCESS, nil
}

func (s *User) ModifyTradePwd(req *proto.UserModifyTradePwdRequest) (result int32, err error) {
	engine := dao.DB.GetMysqlConn()
	ph := new(User)
	var ok bool
	ok, err = engine.ID(req.Uid).Get(ph)
	if err != nil {
		result = ERRCODE_UNKNOWN
		return
	}
	if !ok {
		result = ERRCODE_UNKNOWN
		err = errors.New("get phone number failed")
	}
	//
	if eq := strings.Compare(req.ConfirmPwd, req.NewPwd); eq != 0 {
		return ERRCODE_PWD_COMFIRM, nil
	}
	if eq := strings.Compare(req.NewPwd, ph.PayPwd); eq != 0 {
		return ERRCODE_OLDPWD, nil
	}
	result, err = AuthSms(ph.Phone, SMS_RESET_TRADE_PWD, req.Verify)
	if err != nil {
		return
	}
	if result != ERRCODE_SUCCESS {
		return
	}
	//验证token
	//修改数据库字段

	_, err = engine.ID(req.Uid).Update(&User{PayPwd: req.NewPwd})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	s.RefreshCache(req.Uid)
	return ERRCODE_SUCCESS, nil
}
