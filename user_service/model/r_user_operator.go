package model

import (
	. "digicon/user_service/dao"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"time"
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
return fmt.Sprintf("account:%s:gree:exist",account)
}


func (s *RedisOp) SetSmsCode(phone string, code string, ty int32) (err error) {
	err = DB.GetRedisConn().Set(GetPhoneTagByLogic(phone, ty), code, 600*time.Second).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) GetEmailCode(email string, ty int32) (code string, err error) {
	code, err = DB.GetRedisConn().Get(GetEmailTagByLogic(email, ty)).Result()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) SetEmailCode(email string, code string, ty int32) (err error) {
	err = DB.GetRedisConn().Set(GetEmailTagByLogic(email, ty), code, 600*time.Second).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) GetSmsCode(phone string, ty int32) (code string, err error) {
	code, err = DB.GetRedisConn().Get(GetPhoneTagByLogic(phone, ty)).Result()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) SetTmpGoogleSecertKey(uid uint64, code string) (err error) {
	err = DB.GetRedisConn().Set(GetUserTagByLogic(uid, UID_TAG_GOOGLE_SECERT_KEY), code, 600*time.Second).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) GetTmpGoogleSecertKey(uid uint64) (key string, err error) {
	key, err = DB.GetRedisConn().Get(GetUserTagByLogic(uid, UID_TAG_GOOGLE_SECERT_KEY)).Result()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) SetUserBaseInfo(uid uint64, data string) (err error) {
	err = DB.GetRedisConn().Set(GetUserTagByLogic(uid, UID_TAG_BASE_INFO), data, 1800*time.Second).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) GetUserBaseInfo(uid uint64) (rsp string, err error) {
	rsp, err = DB.GetRedisConn().Get(GetUserTagByLogic(uid, UID_TAG_BASE_INFO)).Result()
	if err == redis.Nil {

	} else if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) SetUserToken(token string, uid uint64) (err error) {
	err = DB.GetRedisConn().Set(getUserToken(uid), token, 604800*time.Second).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

//验证
func (s *RedisOp) GetUserToken(uid uint64) (ret string, err error) {
	ret, err = DB.GetRedisConn().Get(getUserToken(uid)).Result()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	/*
		if err == redis.Nil {
			ret = ERRCODE_TokenVerify
			log.Errorln(err.Error())
			return err, ret

		} else if err != nil {
			ret = ERRCODE_UNKNOWN
			log.Errorln(err.Error())
			return err, ret
		}

		if rsp != token {
			ret = ERRCODE_TokenVerify
			return err, ret
		}

	*/
	return
}

/*
func (s *RedisOp) SetUserToken(uid int32, token []byte) (err error) {
	err = DB.GetRedisConn().Set(GetUserTagByLogic(uid, UID_TAG_TOKEN), token, 604800*time.Second).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) GetUserToken(uid int32) (token []byte, err error) {
	token, err = DB.GetRedisConn().Get(GetUserTagByLogic(uid, UID_TAG_TOKEN)).Bytes()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
*/


func SetGreeSuccess(account string) (err error) {
	err = DB.GetRedisConn().Set(getGree(account),1, 300*time.Second).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}


func  GetGreeSuccess(account string) (err error) {
	err = DB.GetRedisConn().Set(getGree(account),1, 300*time.Second).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func  DelGreeSuccess(account string) (err error) {
	err = DB.GetRedisConn().Del(getGree(account)).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
