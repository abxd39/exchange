package handler

import (
	proto "digicon/proto/rpc"
	"digicon/user_service/log"
	"digicon/user_service/model"

	"golang.org/x/net/context"
)

func (s *RPCServer) ModifyUserLoginPwd(ctx context.Context, req *proto.UserModifyLoginPwdRequest, rsp *proto.UserModifyLoginPwdResponse) (err error) {
	u := model.User{}
	rsp.Err, err = u.ModifyLoginPwd(req)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}

func (s *RPCServer) ModifyPhone1(ctx context.Context, req *proto.UserModifyPhoneRequest, rsp *proto.UserModifyPhoneResponse) (err error) {
	u := model.User{}
	rsp.Err, err = u.ModifyUserPhone1(req)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}

func (s *RPCServer) ModifyPhone2(ctx context.Context, req *proto.UserSetNewPhoneRequest, rsp *proto.UserSetNewPhoneResponse) (err error) {
	u := model.User{}
	rsp.Err, err = u.ModifyUserPhone2(req)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}

func (s *RPCServer) ModifyTradePwd(ctx context.Context, req *proto.UserModifyTradePwdRequest, rsp *proto.UserModifyTradePwdResponse) (err error) {
	u := model.User{}
	rsp.Err, err = u.ModifyTradePwd(req)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}

func (*RPCServer) GetNickName(ctx context.Context, req *proto.UserGetNickNameResquest, rsp *proto.UserGetNickNameResponse) (err error) {
	u := model.UserEx{}
	rsp.Err, err = u.GetNickName(req, rsp)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}

func (*RPCServer) SetNickName(ctx context.Context, req *proto.UserSetNickNameRequest, rsp *proto.UserSetNickNameResponse) (err error) {
	u := model.UserEx{}
	rsp.Err, err = u.SetNickName(req, rsp)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}
