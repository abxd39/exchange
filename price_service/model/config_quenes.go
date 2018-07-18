package model

import (
	. "digicon/price_service/dao"
	. "digicon/price_service/log"
)

type ConfigQuenes struct {
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

func (s *ConfigQuenes) GetAllQuenes() []ConfigQuenes {
	t := make([]ConfigQuenes, 0)
	err := DB.GetMysqlConn().Where("switch=1").Find(&t)
	if err != nil {
		Log.Errorln(err.Error())
		return nil
	}
	return t
}
