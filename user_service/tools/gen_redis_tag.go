package tools

import (
	. "digicon/user_service/log"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	SMS_REGISTER   = 1 //注册业务
	SMS_FORGET     = 2
	SMS_CHANGE_PWD = 3
	SMS_MAX        = 4
)
const (
	LOGIC_SMS      = 1 //短信业务
	LOGIC_SECURITY = 2 //重置密码业务
)

const (
	UID_TAG_GOOGLE_SECERT_KEY = "google_secert"
)

//获取用户redis中逻辑标签信息
func GetPhoneTagByLogic(phone string, ty int32) string {
	switch ty {
	case LOGIC_SMS:
		return getUserTagSms(phone, ty)
	case LOGIC_SECURITY:
		return getUserTagSecurity(phone)
	default:
		break
	}

	Log.WithFields(logrus.Fields{
		"phone": phone,
	}).Error("获取用户redis中逻辑标签信息")
	return ""
}

func getUserTagSecurity(phone string) string {
	return fmt.Sprintf("%s:SecurityKey", phone)
}

func getUserTagSms(phone string, ty int32) string {
	return fmt.Sprintf("%s:Sms:%d", phone, ty)
}

func GetUserTagByLogic(uid int32, tag string) string {
	return fmt.Sprintf("%s:%s", uid, tag)
}
