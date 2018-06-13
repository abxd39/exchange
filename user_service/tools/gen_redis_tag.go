package tools

import (
	"github.com/sirupsen/logrus"
	"fmt"
	. "digicon/user_service/log"
)

const (
	SMS_REGISTER 	=1//注册业务
)
const (
	LOGIC_SMS      = 1 //短信业务
	LOGIC_SECURITY = 2 //重置密码业务
)

//获取用户redis中逻辑标签信息
func GetUserTagByLogic(phone string, logic int, others ...interface{}) string {
	switch logic {
	case LOGIC_SMS:
		if len(others) > 0 {
			ty := others[0]
			return getUserTagSms(phone, ty.(int32))
		}
	case LOGIC_SECURITY:
		return getUserTagSecurity(phone)
	default:
		break
	}

	Log.WithFields(logrus.Fields{
		"phone": phone,
		"logic": logic,
	}).Error("获取用户redis中逻辑标签信息")
	return ""
}

func getUserTagSecurity(phone string) string {
	return fmt.Sprintf("%s:SecurityKey", phone)
}

func getUserTagSms(phone string, ty int32) string {
	return fmt.Sprintf("%s:Sms:%d", phone, ty)
}
