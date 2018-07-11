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
		/*
			r := model.RedisOp{}
			code, err := r.GetSmsCode(req.Ukey, tools.SMS_FORGET)
			if err == redis.Nil {
				rsp.Err = ERRCODE_SMS_CODE_NIL
				rsp.Message = GetErrorMessage(rsp.Err)
				return nil
			} else if err != nil {
				rsp.Err = ERRCODE_UNKNOWN
				rsp.Message = err.Error()
				return nil
			}
			if code == req.Code {
				u := model.User{}
				ret, err := u.GetUserByPhone(req.Ukey)
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

			} else {
				rsp.Err = ERRCODE_SMS_CODE_DIFF
				rsp.Message = GetErrorMessage(rsp.Err)
				return nil
			}
		*/

	} else if req.Type == 2 { //邮箱找回

	}

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
	g:=new(model.Tokens).GetTokens()
	for _,v:=range g{
		rsp.Data=append(rsp.Data,&proto.TokenMarkBaseData{
			TokenId:int32(v.Id),
			Mark:v.Mark,
		})
	}
	return nil
}