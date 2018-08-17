package models

import (
	. "digicon/common/constant"
	"fmt"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"digicon/wallet_service/utils"
)

type RedisOp struct {
}

const (
	UID_TAG_BASE_INFO         = "base"
	UID_TAG_GOOGLE_SECERT_KEY = "google_key"
	UID_TAG_TOKEN             = "token"
)

//获取手机redis中逻辑标签信息
func GetPhoneTagByLogic(phone string, ty int32) string {
	if ty >= SMS_MAX {
		log.WithFields(logrus.Fields{
			"phone": phone,
		}).Error("获取手机redis中逻辑标签信息")
		return ""
	}
	return getUserTagSms(phone, ty)
}

//获取email中逻辑标签信息
func GetEmailTagByLogic(email string, ty int32) string {
	return fmt.Sprintf("email:%s:code:%d", email, ty)
}

//获取用户标签
func GetUserTagByLogic(uid uint64, tag string) string {
	return fmt.Sprintf("user:%d:info:%s", uid, tag)
}

//获取用户短信标签
func getUserTagSms(phone string, ty int32) string {
	return fmt.Sprintf("phone:%s:Sms:%d", phone, ty)
}

func getUserToken(uid uint64) string {
	return fmt.Sprintf("uid:%d:token", uid)
}

func getGree(account string) string {
	return fmt.Sprintf("account:%s:gree:exist", account)
}

func (s *RedisOp) GetEmailCode(email string, ty int32) (code string, err error) {
	code, err = utils.Redis.Get(GetEmailTagByLogic(email, ty)).Result()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) GetSmsCode(phone string, ty int32) (code string, err error) {
	code,err = utils.Redis.Get(GetPhoneTagByLogic(phone, ty)).Result()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
