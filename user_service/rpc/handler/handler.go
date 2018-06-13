package handler

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	"golang.org/x/net/context"

	"digicon/common/random"
	cf "digicon/user_service/conf"
	. "digicon/user_service/log"
	"digicon/user_service/tools"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type RPCServer struct{}

func (s *RPCServer) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	Log.Print("Received Say.Hello request")
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func (s *RPCServer) RegisterByPhone(ctx context.Context, req *proto.RegisterPhoneRequest, rsp *proto.CommonErrResponse) error {
	code, err := DB.GetSmsCode(req.Phone, tools.SMS_REGISTER)
	if err == redis.Nil {
		rsp.Err = ERRCODE_SMS_CODE_NIL
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	} else if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}

	if req.Code == code {
		rsp.Err = DB.RegisterByPhone(req)
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	rsp.Err = ERRCODE_SMS_CODE_DIFF
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

func (s *RPCServer) RegisterByEmail(ctx context.Context, req *proto.RegisterEmailRequest, rsp *proto.CommonErrResponse) error {
	rsp.Err = DB.RegisterByEmail(req)
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

func (s *RPCServer) Login(ctx context.Context, req *proto.LoginRequest, rsp *proto.LoginResponse) error {
	rsp.Err = DB.Login(req.Phone, req.Pwd)
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

func (s *RPCServer) ForgetPwd(ctx context.Context, req *proto.ForgetRequest, rsp *proto.ForgetResponse) error {
	u, ret := DB.GetUserByPhone(req.Phone)
	if ret != ERRCODE_SUCCESS {
		rsp.Err = ret
		rsp.Message = GetErrorMessage(ret)
		return nil
	}
	rsp.Err = ret
	rsp.Message = GetErrorMessage(ret)
	rsp.Email = u.Email
	rsp.Phone = u.Phone
	return nil
}

func (s *RPCServer) AuthSecurity(ctx context.Context, req *proto.SecurityRequest, rsp *proto.SecurityResponse) error {
	security_key, err := DB.GenSecurityKey(req.Phone)
	if err != nil {
		return nil
	}
	rsp.Err = ERRCODE_SUCCESS
	rsp.Message = GetErrorMessage(rsp.Err)
	rsp.SecurityKey = security_key
	return nil
}

func (s *RPCServer) SendSms(ctx context.Context, req *proto.SmsRequest, rsp *proto.CommonErrResponse) error {
	code := random.Random6dec()

	err := DB.SetSmsCode(req.Phone, code, req.Type)
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}
	Log.WithFields(logrus.Fields{
		"account": cf.SmsAccount,
		"pwd":     cf.SmsPwd,
		"code":    code,
	}).Info("print data")

	ret, msg := tools.Send253YunSms(req.Phone, code)
	rsp.Err = ret
	rsp.Message = msg
	return nil
}

func (s *RPCServer) SendEmail(ctx context.Context, req *proto.EmailRequest, rsp *proto.CommonErrResponse) error {
	return nil
}

func (s *RPCServer) ChangePwd(ctx context.Context, req *proto.ChangePwdRequest, rsp *proto.CommonErrResponse) error {
	security_key, err := DB.GetSecurityKeyByPhone(req.Phone)
	if err != nil {
		return nil
	}
	if string(security_key) == string(req.SecurityKey) {
		rsp.Err = ERRCODE_SUCCESS
		rsp.Message = GetErrorMessage(rsp.Err)
	} else {
		rsp.Err = ERRCODE_SECURITY_KEY
		rsp.Message = GetErrorMessage(rsp.Err)
	}
	return nil
}
