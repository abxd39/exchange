package model

import (
	"digicon/common/random"
	. "digicon/proto/common"
	cf "digicon/user_service/conf"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dm"
	"github.com/go-redis/redis"
	"errors"
	. "digicon/user_service/log"
)

func SendEmail(email string, ty int32) (err error) {
	code := random.Random6dec()
	r := &RedisOp{}
	err = r.SetEmailCode(email, code, ty)
	if err != nil {
		return
	}
	return sendAliEmail(email, code)
}

func sendAliEmail(email, code string) (err error) {
	d, err := dm.NewClientWithAccessKey("cn-hangzhou", cf.EmailAppKey, cf.EmailSecretKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	r := dm.CreateSingleSendMailRequest()
	r.AccountName = "emailsdun@email.sdun.io"
	r.AddressType = requests.NewInteger(1)
	r.ToAddress = email
	r.FromAlias = "shendun"
	r.Subject = "欢迎注册神盾"
	r.TextBody = fmt.Sprintf("您好，您正在注册神盾账号。【神盾】安全验证: %s 出于安全原因，该验证码将于10分钟后失效。请勿将验证码透露给他人。", code)
	r.ReplyToAddress = requests.NewBoolean(false)

	h, err := d.SingleSendMail(r)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	ok := h.IsSuccess()
	if ok {
		return
	}

	return errors.New("error send email")
}

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

func ProcessEmailLogic(ty int32, email string) (ret int32, err error) {
	switch ty {
	case SMS_REGISTER:
		u := User{}
		ret, err = u.CheckUserExist(email, "email")
		if err != nil {
			return
		}

		if ret != ERRCODE_SUCCESS {
			return
		}

		err = SendEmail(email, ty)
		if err!=nil {
			return
		}

		return
	default:
		err = SendEmail(email, ty)
		if err!=nil {
			return
		}
		return
	}
}