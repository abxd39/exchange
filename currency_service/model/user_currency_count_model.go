package model

import (
	"digicon/currency_service/dao"
	"fmt"
)

// 订单统计表
type UserCurrencyCount struct {
	Uid     uint64  `xorm:"INT(10)" json:"uid"`
	Orders  uint32  `xorm:"INT(10)" json:"orders"`
	Success uint32  `xorm:"INT(10)" json:"success"`
	Failure uint32  `xorm:"INT(10)" json:"failure"`
	Cancel  uint32  `xorm:"INT(10)" json:"cancel"`
	Good    float64 `xorm:"DECIMAL(10,2)" json:"good"`
}

//
//  获取用户好评率
//
func (this *UserCurrencyCount) GetUserCount(uid uint64) (data UserCurrencyCount, err error) {
	//data := new(UserCurrencyCount)
	engine := dao.DB.GetMysqlConn()
	_, err = engine.Where("uid = ?", uid).Get(&data)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
