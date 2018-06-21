package client

import (
	"context"
	cf "digicon/gateway/conf"
	. "digicon/gateway/log"
	proto "digicon/proto/rpc"

	"github.com/golang/protobuf/jsonpb"
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

func (s *UserRPCCli) CallRegister(ukey string, pwd, invite_code string, country int32, code string, ty int32) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.Register(context.TODO(), &proto.RegisterRequest{
		Ukey:       ukey,
		Pwd:        pwd,
		InviteCode: invite_code,
		Code:       code,
		Type:       ty,
		Country:    country,
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

func (s *UserRPCCli) CallLogin(ukey, pwd string, ty int32) (rsp *proto.LoginResponse, err error) {
	rsp, err = s.conn.Login(context.TODO(), &proto.LoginRequest{
		Ukey: ukey,
		Pwd:  pwd,
		Type: ty,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallForgetPwd(ukey, pwd, code string, ty int32) (rsp *proto.ForgetResponse, err error) {
	rsp, err = s.conn.ForgetPwd(context.TODO(), &proto.ForgetRequest{
		Ukey: ukey,
		Type: ty,
		Pwd:  pwd,
		Code: code,
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

// func (s *UserRPCCli) CallChangePwd(phone, security_key string) (rsp *proto.CommonErrResponse, err error) {
// 	rsp, err = s.conn.ChangePwd(context.TODO(), &proto.ChangePwdRequest{
// 		Phone:       phone,
// 		SecurityKey: []byte(security_key),
// 	})
// 	if err != nil {
// 		Log.Errorln(err.Error())
// 		return
// 	}
// 	return
// }

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

func (s *UserRPCCli) CallDelGoogleSecretKey(uid int32, code uint32) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.DelGoogleSecretKey(context.TODO(), &proto.DelGoogleSecretKeyRequest{
		Uid:  uid,
		Code: code,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

type UserBaseData struct {
	Uid            int32  `json:"uid"`
	Account        string `json:"account"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	GoogleVerifyId bool   `json:"google_verify_id"`
	LoginPwdLevel  int32  `json:"login_pwd_level"`
	SmsTip         bool   `json:"sms_tip"`
	PaySwitch      bool   `json:"pay_switch"`
	NeedPwd        bool   `json:"need_pwd"`
	NeedPwdTime    int32  `json:"need_pwd_time"`
}

func (s *UserRPCCli) CallGetUserBaseInfo(uid int32) (rsp *proto.UserInfoResponse, u *UserBaseData, err error) {
	rsp, err = s.conn.GetUserInfo(context.TODO(), &proto.UserInfoRequest{
		Uid: uid,
	})

	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	out := &proto.UserBaseData{}
	err = jsonpb.UnmarshalString(rsp.Src, out)
	if err != nil {
		return
	}

	u = &UserBaseData{
		Uid:            out.Uid,
		Account:        out.Account,
		Phone:          out.Phone,
		Email:          out.Email,
		GoogleVerifyId: out.GoogleVerifyId,
		LoginPwdLevel:  out.LoginPwdLevel,
		SmsTip:         out.SmsTip,
		PaySwitch:      out.PaySwitch,
		NeedPwd:        out.NeedPwd,
		NeedPwdTime:    out.NeedPwdTime,
	}

	return
}

type UserRealData struct {
	RealName     string `json:"real_name"`
	IdentifyCard string `json:"identify_card"`
}

func (s *UserRPCCli) CallGetUserRealName(uid int32) (rsp *proto.UserRealNameResponse, u *UserRealData, err error) {
	rsp, err = s.conn.GetUserRealName(context.TODO(), &proto.UserInfoRequest{
		Uid: uid,
	})

	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	out := &proto.UserRealData{}
	err = jsonpb.UnmarshalString(rsp.Src, out)
	if err != nil {
		return
	}

	u = &UserRealData{
		RealName:     out.RealName,
		IdentifyCard: out.IdentifyCard,
	}

	return
}

type UserInviteData struct {
	InviteCode string `json:"invite_code"`
	Invites    int32  `json:"invites"`
}

func (s *UserRPCCli) CallGetUserInvite(uid int32) (rsp *proto.UserInviteResponse, u *UserInviteData, err error) {
	rsp, err = s.conn.GetUserInvite(context.TODO(), &proto.UserInfoRequest{
		Uid: uid,
	})

	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	out := &proto.UserInviteData{}
	err = jsonpb.UnmarshalString(rsp.Src, out)
	if err != nil {
		return
	}

	u = &UserInviteData{
		InviteCode: out.InviteCode,
		Invites:    out.Invites,
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

func (s *UserRPCCli) CallModifyUserLoginPwd(req *proto.UserModifyLoginPwdRequest) (*proto.UserModifyLoginPwdResponse, error) {
	rsp, err := s.conn.ModifyUserLoginPwd(context.TODO(), req)
	return rsp, err
}

func (s *UserRPCCli) CallModifyPhone1(req *proto.UserModifyPhoneRequest) (*proto.UserModifyPhoneResponse, error) {
	rsp, err := s.conn.ModifyPhone1(context.TODO(), req)
	return rsp, err
}

func (s *UserRPCCli) CallModifyPhone2(req *proto.UserSetNewPhoneRequest) (*proto.UserSetNewPhoneResponse, error) {
	rsp, err := s.conn.ModifyPhone2(context.TODO(), req)
	return rsp, err
}
func (s *UserRPCCli) CallModifyTradePwd(req *proto.UserModifyTradePwdRequest) (*proto.UserModifyTradePwdResponse, error) {
	rsp, err := s.conn.ModifyTradePwd(context.TODO(), req)
	return rsp, err
}
