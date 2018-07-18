package handler

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	//. "digicon/user_service/dao"
	"golang.org/x/net/context"

	. "digicon/user_service/log"
	"digicon/user_service/model"

	"time"
)

type RPCServer struct{}

func (s *RPCServer) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	Log.Print("Received Say.Hello request")
	rsp.Greeting = "Hello " + req.Name
	return nil
}

//注册
func (s *RPCServer) Register(ctx context.Context, req *proto.RegisterRequest, rsp *proto.CommonErrResponse) error {
	if req.Type == 1 { //手机注册
		/*
			r := model.RedisOp{}
			code, err := r.GetSmsCode(req.Ukey, model.SMS_REGISTER)
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
				u := &model.User{}
				rsp.Err = u.Register(req, "phone")
				rsp.Message = GetErrorMessage(rsp.Err)
				return nil
			}

			rsp.Err = ERRCODE_SMS_CODE_DIFF
			rsp.Message = GetErrorMessage(rsp.Err)

		*/
		ret, err := model.AuthSms(req.Ukey, model.SMS_REGISTER, req.Code)
		if err != nil {
			rsp.Err = ERRCODE_UNKNOWN
			rsp.Message = err.Error()
			return nil
		}
		if ret != ERRCODE_SUCCESS {
			rsp.Err = ret
			return nil
		}
		u := &model.User{}
		rsp.Err = u.Register(req, "phone")
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	} else if req.Type == 2 {
		ret, err := model.AuthEmail(req.Ukey, model.SMS_REGISTER, req.Code)
		if err != nil {
			rsp.Err = ERRCODE_UNKNOWN
			rsp.Message = err.Error()
			return nil
		}
		if ret != ERRCODE_SUCCESS {
			rsp.Err = ret
			return nil
		}
		u := &model.User{}
		rsp.Err = u.Register(req, "email")
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	rsp.Err = ERRCODE_SMS_CODE_DIFF
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

//注册by email
func (s *RPCServer) RegisterByEmail(ctx context.Context, req *proto.RegisterEmailRequest, rsp *proto.CommonErrResponse) error {
	return nil
}

//登陆
func (s *RPCServer) Login(ctx context.Context, req *proto.LoginRequest, rsp *proto.LoginResponse) error {
	u := &model.User{}
	var token string
	var ret int32
	if req.Type == 1 { //手机登陆
		token, ret = u.LoginByPhone(req.Ukey, req.Pwd)
	} else if req.Type == 2 { //邮箱登陆
		token, ret = u.LoginByEmail(req.Ukey, req.Pwd)
	}

	if ret == ERRCODE_SUCCESS {
		new(model.LoginRecord).AddLoginRecord(u.Uid, req.Ip)

		var p proto.LoginUserBaseData
		u.GetLoginUser(&p)
		p.Token = []byte(token)
		rsp.Data = &p
	}
	rsp.Err = ret
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

//忘记密码
func (s *RPCServer) ForgetPwd(ctx context.Context, req *proto.ForgetRequest, rsp *proto.ForgetResponse) error {
	var ret int32
	var err error

	u := &model.User{}
	if req.Type == 1 { //电话找回
		ret, err = u.GetUserByPhone(req.Ukey)
	} else if req.Type == 2 {
		ret, err = u.GetUserByEmail(req.Ukey)
	} else {
		ret = ERRCODE_PARAM
		return nil
	}

	if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return err
	}

	ret, err = u.AuthCodeByAl(req.Ukey, req.Code, req.Type)
	if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return err
	}

	err = u.ModifyPwd(req.Pwd)
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}
	rsp.Err = ERRCODE_SUCCESS
	return nil

	/*
		if req.Type == 1 { //电话找回
			ret, err := model.AuthSms(req.Ukey, req.Type, req.Code)
			if err != nil {
				rsp.Err = ERRCODE_UNKNOWN
				rsp.Message = err.Error()
				return nil
			}

			if ret != ERRCODE_SUCCESS {
				rsp.Err = ret
				rsp.Message = GetErrorMessage(ret)
			}

			u := model.User{}
			ret, err = u.GetUserByPhone(req.Ukey)
			if err != nil {
				rsp.Err = ret
				rsp.Message = err.Error()
				return err
			}

			if ret != ERRCODE_SUCCESS {
				rsp.Err = ret
				rsp.Message = GetErrorMessage(rsp.Err)
				return nil
			}
			err = u.ModifyPwd(req.Pwd)
			if err != nil {
				rsp.Err = ERRCODE_UNKNOWN
				rsp.Message = err.Error()
				return nil
			}
			rsp.Err = ERRCODE_SUCCESS
			rsp.Message = GetErrorMessage(rsp.Err)
			return nil

		} else if req.Type == 2 { //邮箱找回

		}
	*/
	rsp.Err = ERRCODE_PARAM
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

//安全认证
func (s *RPCServer) AuthSecurity(ctx context.Context, req *proto.SecurityRequest, rsp *proto.SecurityResponse) error {
	/*
		security_key, err := DB.GenSecurityKey(req.Phone)
		if err != nil {
			return nil
		}
		rsp.Err = ERRCODE_SUCCESS
		rsp.Message = GetErrorMessage(rsp.Err)
		rsp.SecurityKey = security_key
	*/
	return nil
}

//发生短信验证码
func (s *RPCServer) SendSms(ctx context.Context, req *proto.SmsRequest, rsp *proto.CommonErrResponse) error {
	ret, err := model.ProcessSmsLogic(req.Type, req.Phone, req.Region)
	if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return nil
	}
	rsp.Err = ret
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

//发送邮箱验证码
func (s *RPCServer) SendEmail(ctx context.Context, req *proto.EmailRequest, rsp *proto.CommonErrResponse) error {
	ret, err := model.ProcessEmailLogic(req.Type, req.Email)
	if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return nil
	}
	rsp.Err = ret
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

//改变密码
func (s *RPCServer) ChangePwd(ctx context.Context, req *proto.EmailRequest, rsp *proto.CommonErrResponse) error {
	/*
			security_key, err := DB.GetSecurityKeyByPhone(req.Phone)
			if err != nil {
				return nil
			}
			if string(security_key) == string(req.SecurityKey) {
				u := model.User{}
				ret := u.GetUserByPhone(req.Phone)
				if ret != ERRCODE_SUCCESS {
					rsp.Err = ret
					rsp.Message = GetErrorMessage(rsp.Err)
					return nil
				}

				err = u.ModifyPwd(req.Pwd)
				if err != nil {
					rsp.Err = ERRCODE_UNKNOWN
					rsp.Message = err.Error()
					return nil
				}
				rsp.Err = ERRCODE_SUCCESS
				rsp.Message = GetErrorMessage(rsp.Err)
			} else {
				rsp.Err = ERRCODE_SECURITY_KEY
				rsp.Message = GetErrorMessage(rsp.Err)
			}
	// 	*/
	return nil
}

//获取登陆记录
func (s *RPCServer) GetIpRecord(ctx context.Context, req *proto.CommonPageRequest, rsp *proto.IpRecordResponse) error {
	g := new(model.LoginRecord).GetLoginRecord(req.Uid, int(req.Page), int(req.Limit))
	for _, v := range g {
		rsp.Data = append(rsp.Data, &proto.IpRecordBaseData{
			Ip:          v.Ip,
			CreatedTime: time.Unix(v.CreatedTime, 0).Format("2006-01-02 15:04:05"),
		})
	}

	return nil
}

func (this *RPCServer) TokenList(ctx context.Context, req *proto.NullRequest, rsp *proto.TokenListResponse) error {
	g := new(model.Tokens).GetTokens()
	for _, v := range g {
		rsp.Data = append(rsp.Data, &proto.TokenMarkBaseData{
			TokenId: int32(v.Id),
			Mark:    v.Mark,
		})
	}
	return nil
}

func (this *RPCServer) CheckSecurity(ctx context.Context, req *proto.CheckSecurityRequest, rsp *proto.CheckSecurityResponse) error {
	u := &model.User{}
	var ret int32
	var err error
	if req.Type == 3 {
		ret, err = u.GetUser(req.Uid)

	} else if req.Type == 2 {
		ret, err = u.GetUserByEmail(req.Ukey)

	} else if req.Type == 1 {
		ret, err = u.GetUserByPhone(req.Ukey)
	} else {
		rsp.Err = ERRCODE_PARAM
		return nil
	}
	if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return nil
	}

	if ret != ERRCODE_SUCCESS {
		rsp.Err = ret
		return nil
	}
	rsp.Auth = u.GetAuthMethod()
	return nil
}



/*
	// bind user email
*/
func (this *RPCServer) BindEmail(ctx context.Context, req *proto.BindEmailRequest, rsp *proto.BindPhoneEmailResponse) error{
	u := new(model.User)
	u.GetUser(req.Uid)
	phone := u.Phone
	var err error
	rsp.Code, err = u.AuthCodeByAl(req.Email, req.EmailCode,  model.SMS_BIND_EMAIL )
	if err != nil {
		Log.Errorln("auth code by email error!")
		return err
	}
	if req.VerifyType == 1 {       // 3: 短信校验
		rsp.Code, err = u.AuthCodeByAl(phone, req.VerifyCode, model.SMS_BIND_EMAIL)
		if err != nil {
			return err
		}
	}else if req.VerifyType == 2 {  // 4 谷歌验证
		rsp.Code, err = u.AuthCodeByAl(u.GoogleVerifyId, req.VerifyCode, model.SMS_BIND_EMAIL)
		if err != nil {
			return err
		}
	} else{
		Log.Errorln(" not found verifyType!")
		rsp.Code = ERRCODE_UNKNOWN
		return nil
	}
	err = u.BindUserEmail(req.Email, req.Uid)
	if err != nil {
		Log.Errorln("bind user email err!", err.Error())
		rsp.Code = ERRCODE_UNKNOWN
		return nil
	}
	err = u.SecurityChmod(model.AUTH_EMAIL)
	if err != nil {
		msg := "after bind user email, security chmod error!"
		Log.Errorln(msg)
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Message = msg
	}
	return nil
}

func(this *RPCServer) BindPhone(ctx context.Context, req *proto.BindPhoneRequest, rsp *proto.BindPhoneEmailResponse) error{
	u := new(model.User)
	u.GetUser(req.Uid)
	phone := u.Phone
	var err error
	rsp.Code, err = u.AuthCodeByAl(req.Phone, req.PhoneCode, model.SMS_BIND_PHONE)
	if err != nil {
		return err
	}
	if req.VerifyType == 1 {       //  1. email verify
		rsp.Code, err = u.AuthCodeByAl(phone, req.VerifyCode, model.SMS_BIND_PHONE)
		if err != nil {
			return err
		}
	}else if req.VerifyType == 2 {  // 2. google verify
		rsp.Code, err = u.AuthCodeByAl(u.GoogleVerifyId, req.VerifyCode, model.SMS_BIND_PHONE)
		if err != nil {
			return err
		}
	} else{
		Log.Errorln(" not found verifyType!")
		rsp.Code = ERRCODE_UNKNOWN
		return nil
	}
	err = u.BindUserPhone(req.Phone, req.Uid)
	if err != nil {
		Log.Errorln("bind user phone err!", err.Error())
		rsp.Code = ERRCODE_UNKNOWN
		return nil
	}
	err = u.SecurityChmod(model.AUTH_PHONE)
	if err != nil {
		msg := "after bind user phone, security chmod error!"
		Log.Errorln(msg)
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Message = msg
		return nil
	}
	return nil
}