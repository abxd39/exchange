package handler

import (
	"digicon/common/random"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	//. "digicon/user_service/dao"
	"digicon/user_service/model"
	"fmt"

	"digicon/common/constant"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/jsonpb"
	"golang.org/x/net/context"
	log "github.com/sirupsen/logrus"
	"github.com/GeeTeam/GtGoSdk"
	cf "digicon/user_service/conf"

)

//获取谷歌验密钥
func (s *RPCServer) GetGoogleSecretKey(ctx context.Context, req *proto.GoogleAuthRequest, rsp *proto.GoogleAuthResponse) error {
	u := model.User{}
	ret, err := u.GetUser(req.Uid)
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}
	if ret != ERRCODE_SUCCESS {
		rsp.Err = ret
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	code := random.Krand(16, random.KC_RAND_KIND_UPPER)
	str_code := string(code)
	r := model.RedisOp{}
	r.SetTmpGoogleSecertKey(req.Uid, str_code)
	str := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=sdun.com", u.Account, str_code)
	rsp.SecretKey = str_code
	rsp.Url = str
	return nil
}

//提交谷歌验证码
func (s *RPCServer) AuthGoogleSecretKey(ctx context.Context, req *proto.AuthGoogleSecretKeyRequest, rsp *proto.CommonErrResponse) error {
	r := model.RedisOp{}
	key, err := r.GetTmpGoogleSecertKey(req.Uid)
	if err == redis.Nil {
		rsp.Err = ERRCODE_SMS_CODE_NIL
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	} else if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}

	u := model.User{}
	ret, err := u.GetUser(req.Uid)
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}
	if ret != ERRCODE_SUCCESS {
		rsp.Err = ret
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	ret, err = u.AuthGoogleCode(key, req.Code)
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}
	if ret == ERRCODE_SUCCESS {
		ret := u.SetGoogleSecertKey(req.Uid, key)
		rsp.Err = ret
		return nil
	}
	rsp.Err = ret
	rsp.Message = GetErrorMessage(rsp.Err)
	/*
		code, _ := google.GenGoogleCode(key)
		//code是16进制数据需要转成10进制
		g := strconv.Itoa(int(code))
		r, err := strconv.Atoi(g)
		if err != nil {
			rsp.Err = ERRCODE_UNKNOWN
			rsp.Message = err.Error()
			return err
		}
		if req.Code == uint32(r) {
			u := &model.User{}
			ret := u.SetGoogleSecertKey(req.Uid, key)
			rsp.Err = ret
		} else {
			rsp.Err = ERRCODE_GOOGLE_CODE
		}
	*/
	return nil
}

//解绑谷歌接口
func (s *RPCServer) DelGoogleSecretKey(ctx context.Context, req *proto.DelGoogleSecretKeyRequest, rsp *proto.CommonErrResponse) error {
	u := model.User{}
	ret, err := u.GetUser(req.Uid)
	if ret != ERRCODE_SUCCESS {
		rsp.Err = ret
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	if !u.CheckGoogleExist() {
		rsp.Err = ERRCODE_GOOGLE_CODE_NOT_EXIST
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	ret, err = u.DelGoogleCode(req.Code)
	if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return err
	}
	rsp.Err = ret
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

func (s *RPCServer) ResetGoogleSecretKey(ctx context.Context, req *proto.ResetGoogleSecretKeyRequest, rsp *proto.CommonErrResponse) error {
	u := model.User{}
	ret, err := u.GetUser(req.Uid)
	if ret != ERRCODE_SUCCESS {
		rsp.Err = ret
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	r := model.RedisOp{}
	key, err := r.GetTmpGoogleSecertKey(req.Uid)
	if err == redis.Nil {
		rsp.Err = ERRCODE_SMS_CODE_NIL
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	} else if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}

	m := u.GetAuthMethodExpectGoogle()
	if m == constant.AUTH_PHONE {
		ret, err = u.AuthCodeByAl(u.Phone, req.SmsCode, constant.SMS_SET_GOOGLE_CODE, true)
	} else {
		ret, err = u.AuthCodeByAl(u.Email, req.SmsCode, constant.SMS_SET_GOOGLE_CODE, true)
	}
	//ret, err = u.AuthCodeByAl(req.Ukey, req.SmsCode, model.SMS_FORGET,true)
	if err == redis.Nil {
		rsp.Err = ERRCODE_SMS_CODE_NIL
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	} else if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return err
	}

	if ret != ERRCODE_SUCCESS {
		rsp.Err = ret
		return nil
	}

	ret, err = u.AuthGoogleCode(key, req.AuthCode)
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}
	if ret == ERRCODE_SUCCESS {
		ret = u.SetGoogleSecertKey(req.Uid, key)
		rsp.Err = ret
	}
	return nil
}

//获取用户基础信息
func (s *RPCServer) GetUserInfo(ctx context.Context, req *proto.UserInfoRequest, rsp *proto.UserInfoResponse) error {
	u := &model.User{}
	d, ret, err := u.RefreshCache(req.Uid)
	if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return nil
	}

	m := jsonpb.Marshaler{EmitDefaults: true}

	data, err := m.MarshalToString(d.Base)
	if err != nil {
		return nil
	}

	rsp.Src = data
	rsp.Data = d.Base
	rsp.Err = ERRCODE_SUCCESS
	rsp.Message = GetErrorMessage(ERRCODE_SUCCESS)
	return nil
}

//获取实名信息
func (s *RPCServer) GetUserRealName(ctx context.Context, req *proto.UserInfoRequest, rsp *proto.UserRealNameResponse) error {
	u := &model.User{}
	d, ret, err := u.RefreshCache(req.Uid)
	if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return nil
	}

	m := jsonpb.Marshaler{EmitDefaults: true}

	data, err := m.MarshalToString(d.Real)
	if err != nil {
		return nil
	}
	fmt.Println("----->123456",data)
	rsp.Src = data
	rsp.Data = d.Real
	rsp.Err = ERRCODE_SUCCESS
	rsp.Message = GetErrorMessage(ERRCODE_SUCCESS)
	return nil
}

//邀请信息
func (s *RPCServer) GetUserInvite(ctx context.Context, req *proto.UserInfoRequest, rsp *proto.UserInviteResponse) error {
	u := &model.User{}
	d, ret, err := u.RefreshCache(req.Uid)
	if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return nil
	}

	m := jsonpb.Marshaler{EmitDefaults: true}

	data, err := m.MarshalToString(d.Invite)
	if err != nil {
		return nil
	}

	rsp.Src = data
	rsp.Data = d.Invite
	rsp.Err = ERRCODE_SUCCESS
	rsp.Message = GetErrorMessage(ERRCODE_SUCCESS)
	return nil
}

func (s *RPCServer) Api1(ctx context.Context, req *proto.Api1Request, rsp *proto.Api1Response) error {
	Gt := GtGoSdk.GeetestLib(cf.GtPrivateKey, cf.GtCaptchaID)
	Gt.PreProcess(req.Phone)
	responseMap := Gt.GetResponseMap()

	rsp.Data = &proto.Api1BaseData{}
	r, ok := responseMap["gt"]
	if !ok {
		rsp.Err=ERRCODE_UNKNOWN
		return nil
	}
	rsp.Data.Gt = r.(string)
	r, ok = responseMap["challenge"]
	if !ok {
		rsp.Err=ERRCODE_UNKNOWN
		return nil
	}
	rsp.Data.Challenge = r.(string)

	r, ok = responseMap["success"]
	if !ok {
		rsp.Err=ERRCODE_UNKNOWN
		return nil
	}
	rsp.Data.Success = int32(r.(int))

	log.WithFields(log.Fields{
		"Challenge":   rsp.Data.Challenge,
		"Gt":    rsp.Data.Gt,
		"Phone":req.Phone,
	}).Info("Api1")
	return nil
}

func (s *RPCServer) Api2(ctx context.Context, req *proto.Api2Request, rsp *proto.Api2Response) error {

	Gt := GtGoSdk.GeetestLib(cf.GtPrivateKey, cf.GtCaptchaID)
	var result bool
	if req.Status==0 {
		result = Gt.FailbackValidate(req.Challenge, req.Validate, req.Seccode)
	}else{
		result = Gt.SuccessValidate(req.Challenge, req.Validate, req.Seccode,req.Phone)
	}

	log.WithFields(log.Fields{
		"Challenge":    req.Challenge,
		"Validate":    req.Validate,
		"Seccode":  req.Seccode,
		"Status":    req.Status,
		"Phone": req.Phone,
		"result":result,
	}).Info("Api2")
	if result {
		model.SetGreeSuccess(req.Phone)
		rsp.Err = ERRCODE_SUCCESS
	} else {
		rsp.Err = ERRCODE_UNKNOWN
	}
	return nil
}


func (s *RPCServer) Refresh(ctx context.Context, req *proto.RefreshRequest, rsp *proto.CommonErrResponse) error {
	new(model.User).ForceRefreshCache(req.Uid)
	return nil
}