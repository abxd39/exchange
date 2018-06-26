package model

import (
	"digicon/common/random"
	//"github.com/sirupsen/logrus"
	//. "digicon/user_service/dao"
	. "digicon/proto/common"
	"digicon/user_service/tools"
	"fmt"
	"github.com/go-redis/redis"
)

const (
	SMS_REGISTER   = 1 //注册业务
	SMS_FORGET     = 2
	SMS_CHANGE_PWD = 3

	SMS_MAX = 4
)

//发送短信
func SendSms(phone, region string, ty int32) (ret int32, err error) {
	code := random.Random6dec()
	r := &RedisOp{}
	err = r.SetSmsCode(phone, code, ty)
	if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}

	ret, err = tools.Send253YunSms(fmt.Sprintf("%s%s", region, phone), code)

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

//短信通用处理
func ProcessSmsLogic(ty int32, phone, region string) (ret int32, err error) {
	switch ty {
	case SMS_REGISTER:
		//TODO判断
		u := User{}
		ret, err = u.CheckUserExist(phone, "phone")
		if err != nil {
			return
		}

		if ret != ERRCODE_SUCCESS {
			return
		}

		ret, err = SendSms(phone, region, ty)
	case SMS_FORGET:
		ret, err = SendSms(phone, region, ty)
	case SMS_CHANGE_PWD:
		ret, err = SendSms(phone, region, ty)
	default:
		return

	}
	return

}
