package model

type Price struct {
	Id     int64 `xorm:"BIGINT(20)"`
	Open   int64 `xorm:"comment('开盘价') BIGINT(20)"`
	Close  int64 `xorm:"comment('收盘价') BIGINT(20)"`
	Low    int64 `xorm:"comment('最低价') BIGINT(20)"`
	High   int64 `xorm:"comment('最高价') BIGINT(20)"`
	Amount int64 `xorm:"comment('成交量') BIGINT(20)"`
	Vol    int64 `xorm:"comment('成交额') BIGINT(20)"`
	Count  int64 `xorm:"BIGINT(20)"`
}


