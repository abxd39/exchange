package model

import (
	. "digicon/token_service/dao"
	log "github.com/sirupsen/logrus"
)

type ConfigQuenes struct {
	Id           int64  `xorm:"pk autoincr BIGINT(20)"`
	TokenId      int    `xorm:"comment('交易币') unique(union_quene_id) INT(11)"`
	TokenTradeId int    `xorm:"comment('实际交易币') unique(union_quene_id) INT(11)"`
	Switch       int    `xorm:"comment('开关0关1开') TINYINT(4)"`
	Name         string `xorm:"comment('USDT/BTC') VARCHAR(32)"`
	Price        int64  `xorm:"BIGINT(20)"`
	SellPoundage         int64  `xorm:"comment('卖出手续费') BIGINT(20)"`
	BuyPoundage          int64  `xorm:"comment('买入手续费') BIGINT(20)"`
}

func (s *ConfigQuenes) GetQuenes(uid uint64) []ConfigQuenes {
	return nil
}

func (s *ConfigQuenes) GetAllQuenes() []ConfigQuenes {
	t := make([]ConfigQuenes, 0)
	err := DB.GetMysqlConn().Where("switch=1").Find(&t)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	return t
}

func (s *ConfigQuenes) GetQuenesByType(token_id int32) []ConfigQuenes {
	t := make([]ConfigQuenes, 0)
	err := DB.GetMysqlConn().Where("switch=1 and token_id=?", token_id).Find(&t)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	return t
}
