package model

import (
	. "digicon/token_service/dao"
	. "digicon/token_service/log"
)

type QuenesConfig struct {
	Id           int64  `xorm:"pk autoincr BIGINT(20)"`
	TokenId      int    `xorm:"comment('交易币') unique(union_quene_id) INT(11)"`
	TokenTradeId int    `xorm:"comment('实际交易币') unique(union_quene_id) INT(11)"`
	Switch       int    `xorm:"comment('开关0关1开') TINYINT(4)"`
	Price        int64  `xorm:"comment('初始价格') BIGINT(20)"`
	Name         string `xorm:"comment('USDT/BTC') VARCHAR(32)"`
	Scope        string `xorm:"comment('振幅') DECIMAL(6,2)"`
	Low          int64  `xorm:"comment('最低价') BIGINT(20)"`
	High         int64  `xorm:"comment('最高价') BIGINT(20)"`
	Amount       int64  `xorm:"comment('成交量') BIGINT(20)"`
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

func (s *QuenesConfig) GetQuenesByType(token_id int32) []QuenesConfig {
	t := make([]QuenesConfig, 0)
	err := DB.GetMysqlConn().Where("switch=1 and token_id=?", token_id).Find(&t)
	if err != nil {
		Log.Errorln(err.Error())
		return nil
	}
	return t
}
