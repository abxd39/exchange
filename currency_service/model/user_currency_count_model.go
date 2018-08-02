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

/*
	检查用户记录是否存在
*/
func (this *UserCurrencyCount) CheckUserCurrencyCountExists(uid uint64) (has bool, err error) {
	engine := dao.DB.GetMysqlConn()
	has, err = engine.Exist(&UserCurrencyCount{Uid: uid})
	return
}

/*
	添加用户好评率
*/
func (this *UserCurrencyCount) AddUserCurrencyCount(uid uint64, order uint32, success uint32) (err error) {
	insertSql := "INSERT INTO `user_currency_count` (uid,orders, success) values(?,?,?)"
	engine := dao.DB.GetMysqlConn()
	_, err = engine.Exec(insertSql, this.Uid, this.Orders, this.Success)
	return
}
