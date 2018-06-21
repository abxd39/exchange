package handler

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"digicon/user_service/model"
	"fmt"

	"golang.org/x/net/context"
)

func (s *RPCServer) ModifyUserLoginPwd(ctx context.Context, req *proto.UserModifyLoginPwdRequest, rsp *proto.UserModifyLoginPwdResponse) error {
	var err error
	var value int32
	u := model.User{}

	value, err = u.ModifyLoginPwd(req)
	if err != nil {
		return err
	}
	rsp.Err = value
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

func (s *RPCServer) ModifyPhone1(ctx context.Context, req *proto.UserModifyPhoneRequest, rsp *proto.UserModifyPhoneResponse) error {
	var err error
	var value int32
	u := model.User{}

	value, err = u.ModifyUserPhone1(req)
	if err != nil {
		return err
	}
	rsp.Err = value
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

func (s *RPCServer) ModifyPhone2(ctx context.Context, req *proto.UserSetNewPhoneRequest, rsp *proto.UserSetNewPhoneResponse) error {
	var err error
	var value int32
	u := model.User{}

	value, err = u.ModifyUserPhone2(req)
	if err != nil {
		return err
	}
	rsp.Err = value
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

func (s *RPCServer) ModifyTradePwd(ctx context.Context, req *proto.UserModifyTradePwdRequest, rsp *proto.UserModifyTradePwdResponse) error {
	var err error
	var value int32
	u := model.User{}
	fmt.Printf("ModifyTradePwd%#v\n", req)
	value, err = u.ModifyTradePwd(req)
	if err != nil {
		return err
	}
	rsp.Err = value
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}
