package model

import (
	. "digicon/user_service/dao"
	log "github.com/sirupsen/logrus"
)
/*
type UserLoginLog struct {
	Id          int64  `xorm:"pk autoincr BIGINT(20)"`
	LoginIp     string `xorm:"VARCHAR(64)"`
	Uid         uint64 `xorm:"BIGINT(20)"`
	LoginTime int64  `xorm:"BIGINT(20) created"`
}
*/
type UserLoginLog struct {
	Id           int    `xorm:"not null pk autoincr INT(10)"`
	Uid          uint64  `xorm:"not null comment('用户uid') BIGINT(20)"`
	TerminalType int    `xorm:"not null comment('终端类型') TINYINT(4)"`
	TerminalName string `xorm:"not null default 'web' comment('登录的终端名称') VARCHAR(100)"`
	LoginIp      string `xorm:"not null comment('登录IP') VARCHAR(15)"`
	LoginTime    int64  `xorm:"comment('登录时间戳') BIGINT(11)"`
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
