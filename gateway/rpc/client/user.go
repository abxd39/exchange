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
	conn proto.UserRPCService
}

func (s *UserRPCCli) CallGreet(name string) (rsp *proto.HelloResponse, err error) {
	rsp, err = s.conn.Hello(context.TODO(), &proto.HelloRequest{Name: name})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallRegisterByPhone(phone, pwd, invite_code string, country int, code string) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.RegisterByPhone(context.TODO(), &proto.RegisterPhoneRequest{
		Phone:      phone,
		Pwd:        pwd,
		InviteCode: invite_code,
		Code:       code,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallRegisterByEmail(email, pwd, invite_code string, country int) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.RegisterByEmail(context.TODO(), &proto.RegisterEmailRequest{
		Email:      email,
		Pwd:        pwd,
		InviteCode: invite_code,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallLogin(phone, email, pwd string, ty int32) (rsp *proto.LoginResponse, err error) {
	rsp, err = s.conn.Login(context.TODO(), &proto.LoginRequest{
		Phone: phone,
		Email: email,
		Pwd:   pwd,
		Type:  ty,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallForgetPwd(phone string) (rsp *proto.ForgetResponse, err error) {
	rsp, err = s.conn.ForgetPwd(context.TODO(), &proto.ForgetRequest{
		Phone: phone,
	})

	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallAuthSecurity(phone, phone_code, email_code string) (rsp *proto.SecurityResponse, err error) {
	rsp, err = s.conn.AuthSecurity(context.TODO(), &proto.SecurityRequest{
		Phone:         phone,
		PhoneAuthCode: phone_code,
		EmailAuthCode: email_code,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallSendSms(phone string, ty int32) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.SendSms(context.TODO(), &proto.SmsRequest{
		Phone: phone,
		Type:  ty,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallSendEmail(email string) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.SendEmail(context.TODO(), &proto.EmailRequest{
		Email: email,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallChangePwd(phone, security_key string) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.ChangePwd(context.TODO(), &proto.ChangePwdRequest{
		Phone:       phone,
		SecurityKey: []byte(security_key),
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallGoogleSecretKey(uid int32) (rsp *proto.GoogleAuthResponse, err error) {
	rsp, err = s.conn.GetGoogleSecretKey(context.TODO(), &proto.GoogleAuthRequest{
		Uid: uid,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallAuthGoogleSecretKey(uid int32, code uint32) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.AuthGoogleSecretKey(context.TODO(), &proto.AuthGoogleSecretKeyRequest{
		Uid:  uid,
		Code: code,
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
	greeter := proto.NewUserRPCService(service_name, service.Client())
	u = &UserRPCCli{
		conn: greeter,
	}
	return
}
