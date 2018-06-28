package model

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"digicon/user_service/dao"
	"fmt"
	"strings"
)

func (*User) ModifyLoginPwd(req *proto.UserModifyLoginPwdRequest) (result int32, err error) {

	fmt.Println("lllllllllllllllllllllllllllllllllll")
	fmt.Println(req)
	if va := strings.Compare(req.ConfirmPwd, req.NewPwd); va != 0 {
		return ERRCODE_PWD_COMFIRM, nil
	}
	//验证短信
	result, err = AuthSms(req.Phone, SMS_CHANGE_PWD, req.Verify)
	if err != nil {
		return

	}
	if result != ERRCODE_SUCCESS {
		return
	}
	//token

	//modify DB
	engine := dao.DB.GetMysqlConn()
	_, err = engine.ID(req.Uid).Update(&User{Pwd: req.NewPwd})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	return ERRCODE_SUCCESS, nil

}

func (*User) ModifyUserPhone1(req *proto.UserModifyPhoneRequest) (result int32, err error) {
	//验证短信
	fmt.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")
	fmt.Println("电话号码为：", req.Phone, "验证码为：", req.Verify)
	result, err = AuthSms(req.Phone, SMS_CHANGE_PWD, req.Verify)
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

func (*User) ModifyUserPhone2(req *proto.UserSetNewPhoneRequest) (result int32, err error) {
	result, err = AuthSms(req.Phone, SMS_CHANGE_PWD, req.Verify)
	if err != nil {
		return
	}
	if result != ERRCODE_SUCCESS {
		return
	}
	//token
	//修改数据库字段
	engine := dao.DB.GetMysqlConn()
	_, err = engine.ID(req.Uid).Update(&User{Phone: req.Phone})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	return ERRCODE_SUCCESS, nil
}

func (s *User) ModifyTradePwd(req *proto.UserModifyTradePwdRequest) (result int32, err error) {
	result, err = AuthSms(req.Phone, SMS_CHANGE_PWD, req.Verify)
	if err != nil {
		return
	}
	if result != ERRCODE_SUCCESS {
		return
	}
	//验证token
	//修改数据库字段
	engine := dao.DB.GetMysqlConn()
	_, err = engine.ID(req.Uid).Update(&User{PayPwd: req.NewPwd})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	return ERRCODE_SUCCESS, nil
}
