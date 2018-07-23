package model

import (
	. "digicon/user_service/dao"
	. "digicon/user_service/log"
)

type LoginRecord struct {
	Id          int64  `xorm:"pk autoincr BIGINT(20)"`
	Ip          string `xorm:"VARCHAR(64)"`
	Uid         uint64 `xorm:"BIGINT(20)"`
	CreatedTime int64  `xorm:"BIGINT(20) created"`
}

func (s *LoginRecord) AddLoginRecord(uid uint64, ip string) {

	_, err := DB.GetMysqlConn().InsertOne(&LoginRecord{
		Uid: uid,
		Ip:  ip,
	})

	if err != nil {
		Log.Errorln(err.Error())
		return
	}
}

func (s *LoginRecord) GetLoginRecord(uid uint64, page, limit int) []LoginRecord {
	g := make([]LoginRecord, 0)
	err := DB.GetMysqlConn().Where("uid=?", uid).Desc("created_time").Limit(limit, page-1).Find(&g)
	if err != nil {
		Log.Errorln(err.Error())
		return nil
	}
	return g
}
