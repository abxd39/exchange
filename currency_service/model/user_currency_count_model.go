package model

// 订单统计表
type UserCurrencyCount struct {
	Uid     uint64 `xorm:"INT(10)" json:"uid"`
	Order   uint32 `xorm:"INT(10)" json:"order"`
	Success uint32 `xorm:"INT(10)" json:"success"`
	Failure uint32 `xorm:"INT(10)" json:"failure"`
	Cancel  uint32 `xorm:"INT(10)" json:"cancel"`
	Good    float64 `xorm:"DECIMAL(10,2)" json:"good"`
}
