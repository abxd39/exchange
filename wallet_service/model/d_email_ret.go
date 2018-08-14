package models

import (
	. "digicon/proto/common"
	"fmt"
	"github.com/go-redis/redis"
)


//验证邮箱
func AuthEmail(email string, ty int32, code string) (ret int32, err error) {
	fmt.Println(email, ty, code)
	r := RedisOp{}
	auth_code, err := r.GetEmailCode(email, ty)
	fmt.Println(auth_code)
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
