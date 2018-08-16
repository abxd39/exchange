package handler

import (
	proto "digicon/proto/rpc"
	//"digicon/user_service/log"
	"digicon/user_service/model"

	"github.com/liudng/godump"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (s *RPCServer) ModifyUserLoginPwd(ctx context.Context, req *proto.UserModifyLoginPwdRequest, rsp *proto.UserModifyLoginPwdResponse) (err error) {
	u := model.User{}
	rsp.Err, err = u.ModifyLoginPwd(req)
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func (s *RPCServer) ModifyPhone1(ctx context.Context, req *proto.UserModifyPhoneRequest, rsp *proto.UserModifyPhoneResponse) (err error) {
	u := model.User{}
	rsp.Err, err = u.ModifyUserPhone1(req)
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func (s *RPCServer) ModifyPhone2(ctx context.Context, req *proto.UserSetNewPhoneRequest, rsp *proto.UserSetNewPhoneResponse) (err error) {
	u := model.User{}
	godump.Dump(req)
	rsp.Err, err = u.ModifyUserPhone2(req)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}
	return nil
}

func (s *RPCServer) ModifyTradePwd(ctx context.Context, req *proto.UserModifyTradePwdRequest, rsp *proto.UserModifyTradePwdResponse) (err error) {
	u := model.User{}
	rsp.Err, err = u.ModifyTradePwd(req)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}
	return nil
}

func (*RPCServer) GetNickName(ctx context.Context, req *proto.UserGetNickNameRequest, rsp *proto.UserGetNickNameResponse) (err error) {
	u := model.UserEx{}
	rsp.Err, err = u.GetNickName(req, rsp)
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func (s *RPCServer) SetNickName(ctx context.Context, req *proto.UserSetNickNameRequest, rsp *proto.UserSetNickNameResponse) (err error) {
	u := model.UserEx{}
	rsp.Err, err = u.SetNickName(req, rsp)
	if err != nil {
		log.Errorf(err.Error())
	}
	user := new(model.User)
	user.ForceRefreshCache(req.Uid) // 设置用户昵称时候，刷新缓存
	return nil
}

/*
	短信验证rpc
*/
func (*RPCServer) AuthVerify(ctx context.Context, req *proto.AuthVerifyRequest, rsp *proto.AuthVerifyResponse) (err error) {
	u := model.User{}
	_, err = u.GetUser(req.Uid)
	rsp.Code, err = model.AuthSms(u.Phone, req.AuthType, req.Code)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

func (*RPCServer) FirstRealNameVerify(c context.Context, req *proto.FirstVerifyRequest, rsp *proto.FirstVerifyResponse) (err error) {
	u := model.UserEx{}
	rsp.Code, err = u.SetFirstVerify(req, rsp)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

func (*RPCServer) SecondVerify(c context.Context, req *proto.SecondRequest, rsp *proto.SecondResponse) (err error) {
	u := model.UserSecondaryCertification{}
	rsp.Code, err = u.SetSecondVerify(req, rsp)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	return nil
}

func (*RPCServer) GetVerifyCount(c context.Context, req *proto.VerifyCountRequest, rsp *proto.VerifyCountResponse) (err error) {
	u := model.UserEx{}
	rsp.Code, err = u.GetVerifyCount(req, rsp)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	return nil
}


/*
	VerifyPaypwd
*/
func (*RPCServer) GetVerifyPayPwd(c context.Context, req *proto.VerifyPayPwdRequest, rsp *proto.VerifyPayPwdRespose) ( err error) {
	u := model.User{}
	rsp.Code, err = u.VerifyPayPwd(req.Uid, req.PayPwd)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}