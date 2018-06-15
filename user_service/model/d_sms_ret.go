package model

import (
	"digicon/common/random"
	//"github.com/sirupsen/logrus"
	. "digicon/user_service/dao"
	"digicon/user_service/tools"
)

func SendSms(phone string, ty int32) (ret int32, err_msg string) {
	code := random.Random6dec()

	err := DB.SetSmsCode(phone, code, ty)
	if err != nil {
		err_msg = err.Error()
		return
	}

	ret, msg := tools.Send253YunSms(phone, code)
	err_msg = msg
	return
}
