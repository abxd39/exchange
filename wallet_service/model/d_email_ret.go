package models

import (
	. "digicon/proto/common"
	"github.com/go-redis/redis"
)


//验证邮箱
func AuthEmail(email string, ty int32, code string) (ret int32, err error) {
	r := RedisOp{}
	auth_code, err := r.GetEmailCode(email, ty)
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
