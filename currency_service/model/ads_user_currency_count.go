package model

// 法币交易列表
type AdsUserCurrencyCount struct {
	Ads `xorm:"extends"`
	Success uint32
}
func (AdsUserCurrencyCount) TableName() string {
	return "ads"
}
