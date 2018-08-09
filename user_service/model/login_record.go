package model

import (
	. "digicon/user_service/dao"
	log "github.com/sirupsen/logrus"
)

type UserLoginLog struct {
	Id          int64  `xorm:"pk autoincr BIGINT(20)"`
	LoginIp     string `xorm:"VARCHAR(64)"`
	Uid         uint64 `xorm:"BIGINT(20)"`
	LoginTime int64  `xorm:"BIGINT(20) created"`
}

func (s *UserLoginLog) AddLoginRecord(uid uint64, ip string) {
	_, err := DB.GetMysqlConn().InsertOne(&UserLoginLog{
		Uid: uid,
		LoginIp:  ip,
	})

	if err != nil {
		log.Errorln(err.Error())
		return
	}
}

func (s *UserLoginLog) GetLoginRecord(uid uint64, page, limit int) []UserLoginLog {
	g := make([]UserLoginLog, 0)
	err := DB.GetMysqlConn().Where("uid=?", uid).Desc("login_time").Limit(limit, page-1).Find(&g)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	return g
}
