package client

import (
	"context"
	cf "digicon/gateway/conf"
	. "digicon/gateway/log"
	proto "digicon/proto/rpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

type UserRPCCli struct {
	conn proto.Gateway2UserService
}

func (s *UserRPCCli) CallGreet(name string) (rsp *proto.HelloResponse, err error) {
	rsp, err = s.conn.Hello(context.TODO(), &proto.HelloRequest{Name: name})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallRegister(phone, pwd, invite_code string, country int) (rsp *proto.RegisterResponse, err error) {
	rsp, err = s.conn.Register(context.TODO(), &proto.RegisterRequest{
		Phone:      phone,
		Pwd:        pwd,
		InviteCode: invite_code,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallLogin(phone, pwd string) (rsp *proto.LoginResponse, err error) {
	rsp, err = s.conn.Login(context.TODO(), &proto.LoginRequest{
		Phone: phone,
		Pwd:   pwd,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallForgetPwd(phone string) (rsp *proto.ForgetResponse, err error)  {
	rsp, err =s.conn.ForgetPwd(context.TODO(),&proto.ForgetRequest{
		Phone:phone,
	})

	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallAuthSecurity(phone,phone_code,email_code string) (rsp *proto.SecurityResponse, err error)  {
	rsp, err =s.conn.AuthSecurity(context.TODO(),&proto.SecurityRequest{
		Phone:phone,
		PhoneAuthCode:phone_code,
		EmailAuthCode:email_code,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}




func NewUserRPCCli() (u *UserRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("user.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_user")
	greeter := proto.NewGateway2UserService(service_name, service.Client())
	u = &UserRPCCli{
		conn: greeter,
	}
	return
}
