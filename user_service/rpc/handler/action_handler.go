package handler

import (
	"digicon/common/random"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	"digicon/user_service/model"
	"fmt"
	"github.com/go-redis/redis"
	"golang.org/x/net/context"
)

//提交谷歌验证码
func (s *RPCServer) GetGoogleSecretKey(ctx context.Context, req *proto.GoogleAuthRequest, rsp *proto.GoogleAuthResponse) error {
	u := model.GetUser(req.Uid)
	if u.CheckGoogleExist() { //检查是否已经有谷歌私钥，有的话不能再次申请
		rsp.Err = ERRCODE_GOOGLE_CODE_EXIST
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	code := random.Krand(16, random.KC_RAND_KIND_UPPER)
	str_code := string(code)

	DB.SetTmpGoogleSecertKey(req.Uid, str_code)
	str := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=sdun.com", u.Account, str_code)
	rsp.SecretKey = str_code
	rsp.Url = str
	return nil
}

//提交谷歌验证码
func (s *RPCServer) AuthGoogleSecretKey(ctx context.Context, req *proto.AuthGoogleSecretKeyRequest, rsp *proto.CommonErrResponse) error {
	key, err := DB.GetTmpGoogleSecertKey(req.Uid)
	if err == redis.Nil {
		rsp.Err = ERRCODE_SMS_CODE_NIL
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	} else if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}

	u := model.GetUser(req.Uid)
	ret, err := u.AuthGoogleCode(key, req.Code)
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
	u := model.GetUser(req.Uid)

	if !u.CheckGoogleExist() {
		rsp.Err = ERRCODE_GOOGLE_CODE_NOT_EXIST
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	ret, err := u.DelGoogleCode(req.Code)
	if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return err
	}
	rsp.Err = ret
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}
