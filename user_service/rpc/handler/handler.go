package handler

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	"golang.org/x/net/context"
	"log"
)

type RPCServer struct{}

func (s *RPCServer) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	log.Print("Received Say.Hello request")
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func (s *RPCServer) Register(ctx context.Context, req *proto.RegisterRequest, rsp *proto.RegisterResponse) error {
	rsp.Err = DB.Register(req)
	return nil
}

func (s *RPCServer) Login(ctx context.Context, req *proto.LoginRequest, rsp *proto.LoginResponse) error {
	rsp.Err = DB.Login(req.Phone, req.Pwd)
	return nil
}

func (s *RPCServer) ForgetPwd(ctx context.Context, req *proto.ForgetRequest, rsp *proto.ForgetResponse) error {
	u, ret := DB.GetUserByPhone(req.Phone)
	if ret != ERRCODE_SUCCESS {
		rsp.Err = ret
		rsp.Message = GetErrorMessage(ret)
		return nil
	}
	rsp.Err=ret
	rsp.Message=GetErrorMessage(ret)
	rsp.Email=u.Email
	rsp.Phone=u.Phone
	return nil
}

func (s *RPCServer) AuthSecurity(ctx context.Context, req *proto.SecurityRequest, rsp *proto.SecurityResponse) error {
	security_key,err :=DB.GenSecurityKey(req.Phone)
	if  err!=nil{
		return nil
	}
	rsp.Err=ERRCODE_SUCCESS
	rsp.Message=GetErrorMessage(rsp.Err)
	rsp.SecurityKey=security_key
	return nil
}

