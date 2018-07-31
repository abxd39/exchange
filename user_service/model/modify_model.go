package model

import (
	. "digicon/common/constant"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"digicon/user_service/dao"
	"fmt"
	"strings"
)

func (s *User) ModifyLoginPwd(req *proto.UserModifyLoginPwdRequest) (int32, error) {
	fmt.Println("0..0.0.0.0.0.0.0.0.0.0.00.0.0.0.0000.0.0.0.")
	fmt.Println(req)
	if va := strings.Compare(req.ConfirmPwd, req.NewPwd); va != 0 {
		return ERRCODE_PWD_COMFIRM, nil
	}
	//modify DB
	engine := dao.DB.GetMysqlConn()
	ph := new(User)
	var ok bool
	ok, err := engine.Where("uid=?", req.Uid).Get(ph)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if !ok {
		return ERRCODE_ACCOUNT_NOTEXIST, nil
	}
	//旧密码的判断
	if b := strings.Compare(req.OldPwd, ph.Pwd); b != 0 {
		return ERRCODE_OLDPWD, nil
	}
	fmt.Println("lllllllllllllllllllllllllllllllllll")
	fmt.Println(ph.Phone)
	//验证短信
	//model.AuthSms()
	result, err := AuthSms(ph.Phone, SMS_MODIFY_LOGIN_PWD, req.Verify)
	if err != nil {
		return ERRCODE_UNKNOWN, err

	}
	if result != ERRCODE_SUCCESS {
		return result, nil
	}
	//token
	_, err = engine.Where("uid=?", req.Uid).Update(&User{Pwd: req.NewPwd})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	s.ForceRefreshCache(req.Uid)
	return ERRCODE_SUCCESS, nil

}

func (s *User) ModifyUserPhone1(req *proto.UserModifyPhoneRequest) (int32, error) {
	//验证短信
	engine := dao.DB.GetMysqlConn()
	ph := new(User)
	var ok bool
	ok, err := engine.Where("uid=?", req.Uid).Get(ph)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if !ok {
		return ERRCODE_ACCOUNT_NOTEXIST, nil
	}
	fmt.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")
	fmt.Println("电话号码为：", ph.Phone, "验证码为：", req.Verify)
	result, err := AuthSms(ph.Phone, SMS_MODIFY_PHONE, req.Verify)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if result != ERRCODE_SUCCESS {
		return result, nil
	}
	//token
	//旧的电话号码验证通过
	return ERRCODE_SUCCESS, nil
}

func (s *User) ModifyUserPhone2(req *proto.UserSetNewPhoneRequest) (int32, error) {

	result, err := AuthSms(req.Phone, SMS_MODIFY_PHONE, req.Verify)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if result != ERRCODE_SUCCESS {
		return result, err
	}
	//token
	//修改数据库字段
	engine := dao.DB.GetMysqlConn()
	u := new(User)
	has, err := engine.Where("uid=?", req.Uid).Get(u)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if !has {
		return ERRCODE_ACCOUNT_NOTEXIST, nil
	}
	//全表检索该号码是否已经有账号绑定
	has, err = engine.Where("phone=?", req.Phone).Exist(&User{})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if has {
		return ERRCODE_PHONE_EXIST, nil
	}
	_, err = engine.Where("uid=?", req.Uid).Update(&User{Phone: req.Phone})
	if err != nil {
		fmt.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv", err)
		return ERRCODE_UNKNOWN, err
	}
	s.ForceRefreshCache(req.Uid)
	return ERRCODE_SUCCESS, nil
}

func (s *User) ModifyTradePwd(req *proto.UserModifyTradePwdRequest) (int32, error) {
	engine := dao.DB.GetMysqlConn()
	ph := new(User)
	var ok bool
	ok, err := engine.Where("uid=?", req.Uid).Get(ph)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if !ok {
		return ERRCODE_ACCOUNT_NOTEXIST, nil
	}
	//
	if eq := strings.Compare(req.ConfirmPwd, req.NewPwd); eq != 0 {
		return ERRCODE_PWD_COMFIRM, nil
	}

	result, err := AuthSms(ph.Phone, SMS_RESET_TRADE_PWD, req.Verify)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if result != ERRCODE_SUCCESS {
		return result, nil
	}
	//验证token
	//修改数据库字段
	ph.SetTardeMark = ph.SetTardeMark ^ AUTH_TRADEMARK
	_, err = engine.Where("uid=?", req.Uid).Update(&User{
		PayPwd:       req.NewPwd,
		SetTardeMark: ph.SetTardeMark,
	})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	s.ForceRefreshCache(req.Uid)
	return ERRCODE_SUCCESS, nil
}
