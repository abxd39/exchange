package client

import (
	"context"
	cf "digicon/gateway/conf"
	proto "digicon/proto/rpc"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"strings"
)

type UserRPCCli struct {
	conn proto.UserRPCService
}

func (s *UserRPCCli) CallFirstVerify(req *proto.FirstVerifyRequest) (rsp *proto.FirstVerifyResponse, err error) {
	return s.conn.FirstRealNameVerify(context.TODO(), req)
}

func (s *UserRPCCli) CallSecondVerify(req *proto.SecondRequest) (rsp *proto.SecondResponse, err error) {
	return s.conn.SecondVerify(context.TODO(), req)
}

func (s *UserRPCCli) CallGetVerifyCount(uid uint64) (rsp *proto.VerifyCountResponse, err error) {
	return s.conn.GetVerifyCount(context.TODO(), &proto.VerifyCountRequest{
		Uid: uid,
	})
}

func (s *UserRPCCli) CallApi1(phone string) (rsp *proto.Api1Response, err error) {
	rsp, err = s.conn.Api1(context.TODO(), &proto.Api1Request{Phone: phone})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallApi2(req *proto.Api2Request) (rsp *proto.Api2Response, err error) {
	rsp, err = s.conn.Api2(context.TODO(), req)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallGreet(name string) (rsp *proto.HelloResponse, err error) {
	rsp, err = s.conn.Hello(context.TODO(), &proto.HelloRequest{Name: name})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallRegister(ukey string, pwd, invite_code string, country string, code string, ty int32) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.Register(context.TODO(), &proto.RegisterRequest{
		Ukey:       ukey,
		Pwd:        pwd,
		InviteCode: invite_code,
		Code:       code,
		Type:       ty,
		Country:    country,
	})
	if err != nil {
		log.Errorln(err.Error())
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
		log.Errorln(err.Error())
		return
	}
	return
}

type LoginUserBaseData struct {
	Uid   uint64 `json:"uid"`
	Token string `json:"token"`
}

func (s *UserRPCCli) CallLogin(ukey, pwd string, ty int32, ip string) (rsp *proto.LoginResponse, err error) {
	rsp, err = s.conn.Login(context.TODO(), &proto.LoginRequest{
		Ukey: ukey,
		Pwd:  pwd,
		Type: ty,
		Ip:   ip,
	})
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	return
}

func (s *UserRPCCli) CallTokenVerify(uid uint64, token []byte) (rsp *proto.TokenVerifyResponse, err error) {
	rsp, err = s.conn.TokenVerify(context.TODO(), &proto.TokenVerifyRequest{
		Uid:   uid,
		Token: token,
	})
	if err != nil {
		log.Errorln(err.Error())
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
		log.Errorln(err.Error())
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
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallSendSms(phone, region string, ty int32) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.SendSms(context.TODO(), &proto.SmsRequest{
		Phone:  phone,
		Type:   ty,
		Region: region,
	})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallSendEmail(email string, ty int32) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.SendEmail(context.TODO(), &proto.EmailRequest{
		Email: email,
		Type:  ty,
	})
	if err != nil {
		log.Errorln(err.Error())
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
// 		log.Errorln(err.Error())
// 		return
// 	}
// 	return
// }

func (s *UserRPCCli) CallGoogleSecretKey(uid uint64) (rsp *proto.GoogleAuthResponse, err error) {
	rsp, err = s.conn.GetGoogleSecretKey(context.TODO(), &proto.GoogleAuthRequest{
		Uid: uid,
	})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallAuthGoogleSecretKey(uid uint64, code uint32) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.AuthGoogleSecretKey(context.TODO(), &proto.AuthGoogleSecretKeyRequest{
		Uid:  uid,
		Code: code,
	})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallDelGoogleSecretKey(uid uint64, code uint32) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.DelGoogleSecretKey(context.TODO(), &proto.DelGoogleSecretKeyRequest{
		Uid:  uid,
		Code: code,
	})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
func (s *UserRPCCli) CallResetGoogleSecretKey(p *proto.ResetGoogleSecretKeyRequest) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.ResetGoogleSecretKey(context.TODO(), p)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

type UserBaseData struct {
	Uid            uint64 `json:"uid"`
	Account        string `json:"account"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	GoogleVerifyId bool   `json:"google_verify_id"`
	LoginPwdLevel  int32  `json:"login_pwd_level"`
	SmsTip         bool   `json:"sms_tip"`
	PaySwitch      bool   `json:"pay_switch"`
	NeedPwd        bool   `json:"need_pwd"`
	NeedPwdTime    int32  `json:"need_pwd_time"`
	Country        string `json:"country"`
	GoogleExist    bool   `json:"google_exist"`
	NickName       string `json:"nick_name"`
	HeadSculpture  string `json:"head_scul"`
}

func replaceNickName(nickname string) (rpname string) {
	if nickname != "" {
		nickLen := len(nickname)
		if nickLen <= 4 {
			rpname = strings.Replace(nickname, nickname[1:nickLen-1], "****", -1)
		} else if nickLen < 7 {
			rpname = strings.Replace(nickname, nickname[2:nickLen-2], "****", -1)
		} else {
			if strings.Contains(nickname, "@") {
				rpname = strings.Replace(nickname, nickname[3:nickLen-4], "****", -1)
			} else {
				rpname = strings.Replace(nickname, nickname[3:nickLen-4], "***", -1)
			}
		}
	} else {
		rpname = ""
	}
	return
}

func (s *UserRPCCli) CallGetUserBaseInfo(uid uint64) (rsp *proto.UserInfoResponse, u *UserBaseData, err error) {
	rsp, err = s.conn.GetUserInfo(context.TODO(), &proto.UserInfoRequest{
		Uid: uid,
	})

	if err != nil {
		log.Errorln(err.Error())
		return
	}

	out := &proto.UserBaseData{}
	err = jsonpb.UnmarshalString(rsp.Src, out)
	if err != nil {
		return
	}
	var nickname string
	if out.NickName == "" {
		nickname = replaceNickName(out.Account)
	} else {
		nickname = out.NickName
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
		Country:        out.Country,
		GoogleExist:    out.GoogleExist,
		NickName:       nickname,
		HeadSculpture:  out.HeadSculpture,
	}

	return
}

type UserRealData struct {
	RealName     string `json:"real_name"`
	IdentifyCard string `json:"identify_card"`
	SecondMark   int32  `json:"second_mark"`
	CheckMarkFirst int32 `json:"check_mark_first"`
	CheckMarkSecond int32 `json:"check_mark_second"`
}

func (s *UserRPCCli) CallGetUserRealName(uid uint64) (rsp *proto.UserRealNameResponse, u *UserRealData, err error) {
	rsp, err = s.conn.GetUserRealName(context.TODO(), &proto.UserInfoRequest{
		Uid: uid,
	})

	if err != nil {
		log.Errorln(err.Error())
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
		CheckMarkFirst: out.CheckMarkFirst,
		CheckMarkSecond:out.CheckMarkSecond,
	}

	return
}

type UserInviteData struct {
	InviteCode string `json:"invite_code"`
	Invites    int32  `json:"invites"`
}

func (s *UserRPCCli) CallGetUserInvite(uid uint64) (rsp *proto.UserInviteResponse, u *UserInviteData, err error) {
	rsp, err = s.conn.GetUserInvite(context.TODO(), &proto.UserInfoRequest{
		Uid: uid,
	})

	if err != nil {
		log.Errorln(err.Error())
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

func (s *UserRPCCli) CallGetIpRecord(uid uint64, limit, page int32) (rsp *proto.IpRecordResponse, err error) {
	rsp, err = s.conn.GetIpRecord(context.TODO(), &proto.CommonPageRequest{
		Uid:   uid,
		Limit: limit,
		Page:  page,
	})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallTokensList() (rsp *proto.TokenListResponse, err error) {
	rsp, err = s.conn.TokenList(context.TODO(), &proto.NullRequest{})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallCheckAuthSecurity(p *proto.CheckSecurityRequest) (rsp *proto.CheckSecurityResponse, err error) {
	rsp, err = s.conn.CheckSecurity(context.TODO(), p)
	if err != nil {
		log.Errorln(err.Error())
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

func (s *UserRPCCli) CallModifyUserLoginPwd(req *proto.UserModifyLoginPwdRequest) (*proto.UserModifyLoginPwdResponse, error) {
	return s.conn.ModifyUserLoginPwd(context.TODO(), req)
}

func (s *UserRPCCli) CallModifyPhone1(req *proto.UserModifyPhoneRequest) (*proto.UserModifyPhoneResponse, error) {
	fmt.Println("1111111111111111111111")
	fmt.Println(req)
	return s.conn.ModifyPhone1(context.TODO(), req)
}

func (s *UserRPCCli) CallModifyPhone2(req *proto.UserSetNewPhoneRequest) (*proto.UserSetNewPhoneResponse, error) {
	return s.conn.ModifyPhone2(context.TODO(), req)
}

func (s *UserRPCCli) CallModifyTradePwd(req *proto.UserModifyTradePwdRequest) (*proto.UserModifyTradePwdResponse, error) {
	return s.conn.ModifyTradePwd(context.TODO(), req)
}

func (s *UserRPCCli) CallSetNickName(req *proto.UserSetNickNameRequest) (*proto.UserSetNickNameResponse, error) {
	return s.conn.SetNickName(context.TODO(), req)
}
func (s *UserRPCCli) CallGetNickName(req *proto.UserGetNickNameRequest) (*proto.UserGetNickNameResponse, error) {
	return s.conn.GetNickName(context.TODO(), req)
}

func (s *UserRPCCli) CallBindEmail(req *proto.BindEmailRequest) (*proto.BindPhoneEmailResponse, error) {
	return s.conn.BindEmail(context.TODO(), req)
}

func (s *UserRPCCli) CallBindPhone(req *proto.BindPhoneRequest) (*proto.BindPhoneEmailResponse, error) {
	return s.conn.BindPhone(context.TODO(), req)
}
