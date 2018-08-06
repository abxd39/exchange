package model

import (
	"digicon/common/random"
	. "digicon/proto/common"
	cf "digicon/user_service/conf"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dm"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	. "digicon/common/constant"
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
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"email": email,
				"code":  code,
			}).Errorf("sendAliEmail error %s", err.Error())
		}
	}()

	d, err := dm.NewClientWithAccessKey("cn-hangzhou", cf.EmailAppKey, cf.EmailSecretKey)
	if err != nil {
		return
	}

	r := dm.CreateSingleSendMailRequest()
	r.AccountName = "emailsdun@email.sdun.io"
	r.AddressType = requests.NewInteger(1)
	r.ToAddress = email
	r.FromAlias = "shendun"
	r.Subject = "欢迎注册UNT"
	r.TextBody = fmt.Sprintf("您好，您正在注册UNT账号。【UNT】安全验证: %s 出于安全原因，该验证码将于10分钟后失效。请勿将验证码透露给他人。", code)
	r.ReplyToAddress = requests.NewBoolean(false)

	h, err := d.SingleSendMail(r)
	if err != nil {
		return
	}

	ok := h.IsSuccess()
	if ok {
		return
	}

	return errors.New("error send email:[" + h.String() + "]")
}

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
		if err != nil {
			return
		}

		return
	default:
		err = SendEmail(email, ty)
		if err != nil {
			return
		}
		return
	}
}
