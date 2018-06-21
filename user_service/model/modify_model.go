package model

import (
	proto "digicon/proto/rpc"
	"digicon/user_service/dao"
	"strings"
)

func (s *User) ModifyLoginPwd(req *proto.UserModifyLoginPwdRequest) (int32, error) {
	var value int32
	var err error
	//
	if va := strings.Compare(req.ConfirmPwd, req.NewPwd); va != 0 {
		return 205, nil
	}
	//验证短信
	value, err = AuthSms(req.Phone, SMS_CHANGE_PWD, req.Verify)
	if err != nil {
		return value, err
	}
	//token

	//modify DB
	engine := dao.DB.GetMysqlConn()
	var op int64
	op, err = engine.ID(req.Uid).Update(&User{Pwd: req.NewPwd})
	if err != nil {
		return int32(op), err
	}
	return 0, nil

}

func (s *User) ModifyUserPhone1(req *proto.UserModifyPhoneRequest) (int32, error) {
	var value int32
	var err error
	//验证短信
	value, err = AuthSms(req.Phone, SMS_CHANGE_PWD, req.Verify)
	if err != nil {
		return value, err
	}
	//token
	//旧的电话号码验证通过
	return 0, nil
}

func (s *User) ModifyUserPhone2(req *proto.UserSetNewPhoneRequest) (int32, error) {
	var value int32
	var err error
	value, err = AuthSms(req.Phone, SMS_CHANGE_PWD, req.Verify)
	if err != nil {
		return value, err
	}
	//token
	//修改数据库字段
	engine := dao.DB.GetMysqlConn()
	var op int64
	op, err = engine.ID(req.Uid).Update(&User{Phone: req.Phone})
	if err != nil {
		return int32(op), err
	}
	return 0, nil
}

func (s *User) ModifyTradePwd(req *proto.UserModifyTradePwdRequest) (int32, error) {
	var value int32
	var err error
	value, err = AuthSms(req.Phone, SMS_CHANGE_PWD, req.Verify)
	if err != nil {
		return value, err
	}
	//验证token
	//修改数据库字段
	engine := dao.DB.GetMysqlConn()
	var op int64
	op, err = engine.ID(req.Uid).Update(&User{PayPwd: req.NewPwd})
	if err != nil {
		return int32(op), err
	}
	return 0, nil
}
