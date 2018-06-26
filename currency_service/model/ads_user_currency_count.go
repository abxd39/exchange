package model

// 法币交易列表 - 用户虚拟币-订单统计 - 用户虚拟货币资产
type AdsUserCurrencyCount struct {
	Ads     `xorm:"extends"`
	Balance int64
	Freeze  int64
	Success uint32
}

func (AdsUserCurrencyCount) TableName() string {
	return "ads"
}
