package model

import (
	"digicon/common/random"
	//"github.com/sirupsen/logrus"
	//. "digicon/user_service/dao"
	. "digicon/proto/common"
	"digicon/user_service/tools"
	"github.com/go-redis/redis"
)

const (
	SMS_REGISTER   = 1 //注册业务
	SMS_FORGET     = 2
	SMS_CHANGE_PWD = 3

	SMS_MAX        = 4
)

//发送短信
func SendSms(phone string, ty int32) (ret int32, err_msg string) {
	code := random.Random6dec()
	r := &RedisOp{}
	err := r.SetSmsCode(phone, code, ty)
	if err != nil {
		err_msg = err.Error()
		return
	}

	ret, msg := tools.Send253YunSms(phone, code)
	err_msg = msg
	return
}

//验证短信
func AuthSms(phone string, ty int32, code string) (ret int32, err error) {
	r := RedisOp{}
	auth_code, err := r.GetSmsCode(phone, ty)
	if err == redis.Nil {
		ret = ERRCODE_SMS_CODE_NIL
		return
	} else if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}

	if code != auth_code {
		ret = ERRCODE_SMS_CODE_DIFF
		return
	}

	ret = ERRCODE_SUCCESS
	return
}
