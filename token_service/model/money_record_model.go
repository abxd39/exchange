package model

import (
	. "digicon/token_service/dao"
	. "digicon/token_service/log"
)

type MoneyRecord struct {
	Id          int64  `xorm:"pk autoincr BIGINT(20)"`
	Uid         int    `xorm:"INT(11)"`
	TokenId     int    `xorm:"INT(11)"`
	Hash        string `xorm:"unique VARCHAR(128)"`
	Opt         int    `xorm:"TINYINT(1)"`
	Num         int64  `xorm:"BIGINT(20)"`
	CreatedTime int64  `xorm:"BIGINT(20)"`
}

func (s *MoneyRecord) CheckExist(hash string) (ok bool, err error) {
	ok, err = DB.GetMysqlConn().Where("hash=?", hash).Get(&MoneyRecord{})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *MoneyRecord) InsertRecord(p *MoneyRecord) {
	_, err := DB.GetMysqlConn().InsertOne(p)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}
