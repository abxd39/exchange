package model

import (
	. "digicon/token_service/dao"
	. "digicon/token_service/log"
)

type QuenesConfig struct {
	Id           int64
	TokenId      int    `xorm:"unique(union_quene_id) INT(11)"`
	TokenTradeId int    `xorm:"unique(union_quene_id) INT(11)"`
	Switch       int    `xorm:"TINYINT(4)"`
	Price        int64  `xorm:"INT(20)"`
	Name         string `xorm:"varchar(32)"`
	Scope        string `xorm:"varchar(32)"`
}

func (s *QuenesConfig) GetQuenes(uid uint64) []QuenesConfig {
	/*
		t := make([]QuenesConfig, 0)
		err := DB.GetMysqlConn().Where("token_id=? and switch=1", quene_type).Find(&t)
		if err != nil {
			Log.Errorln(err.Error())
			return nil
		}
	*/
	return nil
}

func (s *QuenesConfig) GetAllQuenes() []QuenesConfig {
	t := make([]QuenesConfig, 0)
	err := DB.GetMysqlConn().Where("switch=1").Find(&t)
	if err != nil {
		Log.Errorln(err.Error())
		return nil
	}
	return t
}
