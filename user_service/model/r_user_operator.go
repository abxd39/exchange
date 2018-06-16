package model

import (
	. "digicon/user_service/dao"
	. "digicon/user_service/log"
	"digicon/user_service/tools"
	"time"
)

type RedisOp struct {
}

func (s *RedisOp) SetSmsCode(phone string, code string, ty int32) (err error) {
	err = DB.GetRedisConn().Set(tools.GetPhoneTagByLogic(phone, ty), code, 600*time.Second).Err()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) GetSmsCode(phone string, ty int32) (code string, err error) {
	code, err = DB.GetRedisConn().Get(tools.GetPhoneTagByLogic(phone, ty)).Result()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) SetTmpGoogleSecertKey(uid int32, code string) (err error) {
	err = DB.GetRedisConn().Set(tools.GetUserTagByLogic(uid, tools.UID_TAG_GOOGLE_SECERT_KEY), code, 600*time.Second).Err()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) GetTmpGoogleSecertKey(uid int32) (key string, err error) {
	key, err = DB.GetRedisConn().Get(tools.GetUserTagByLogic(uid, tools.UID_TAG_GOOGLE_SECERT_KEY)).Result()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) SetUserBaseInfo(uid int32, data string) (err error) {
	err = DB.GetRedisConn().Set(tools.GetUserTagByLogic(uid, tools.UID_TAG_BASE_INFO), data, 60*time.Second).Err()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *RedisOp) GetUserBaseInfo(uid int32) (rsp string, err error) {
	rsp, err = DB.GetRedisConn().Get(tools.GetUserTagByLogic(uid, tools.UID_TAG_BASE_INFO)).Result()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}
