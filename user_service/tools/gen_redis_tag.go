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

/*
const (
	LOGIC_SMS      = 1 //短信业务
	LOGIC_SECURITY = 2 //重置密码业务
)
*/
/*
const (
	UID_TAG_GOOGLE_SECERT_KEY = "google_secert"
	UID_TAG_BASE_INFO="user:info"
)
*/

const (
	UID_TAG_BASE_INFO         = "base"
	UID_TAG_GOOGLE_SECERT_KEY = "google_key"
)

//获取手机redis中逻辑标签信息
func GetPhoneTagByLogic(phone string, ty int32) string {
	if ty >= SMS_MAX {
		Log.WithFields(logrus.Fields{
			"phone": phone,
		}).Error("获取手机redis中逻辑标签信息")
		return ""
	}
	return getUserTagSms(phone, ty)
}

//获取用户标签
func GetUserTagByLogic(uid int32, tag string) string {
	return fmt.Sprintf("user:%s:info:%s", uid, tag)
}

/*
func getUserTagSecurity(phone string) string {
	return fmt.Sprintf("%s:SecurityKey", phone)
}
*/
//获取用户短信标签
func getUserTagSms(phone string, ty int32) string {
	return fmt.Sprintf("phone:%s:Sms:%d", phone, ty)
}
