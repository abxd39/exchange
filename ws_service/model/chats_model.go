package model

import (
	"digicon/ws_service/dao"
	. "digicon/ws_service/log"
	. "digicon/proto/common"
	"fmt"
)

// 订单聊天
type Chats struct {
	Id          uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	OrderId     string `xorm:"VARBINARY(64)" json:"order_id"`
	IsOrderUser string `xorm:"TINYINT(1)" json:"is_order_user"`
	Uid         uint64 `xorm:"INT(10)" json:"uid"`
	Uname       string `xorm:"VARBINARY(10)" json:"uname"`
	Content     string `xorm:"VARBINARY(10)" json:"content"`
	States      uint32 `xorm:"TINYINT(1)" json:"states"`
	CreatedTime string `xorm:"DATETIME" json:"created_time"`
}

func (this *Chats) Add() int {

	//
	ord := new(Order)
	isOrd, err := dao.DB.GetMysqlConn().Where("order_id=?", this.OrderId).Get(ord)
	if err != nil {
		Log.Errorln(err.Error())
		fmt.Println(err.Error())
		return ERRCODE_UNKNOWN
	}
	if !isOrd {
		Log.Errorln(isOrd)
		//return ERRCODE_ORDER_NOTEXIST
	}

	// 是否为订单主用户: 0否 1是
	if this.Uid == ord.SellId {
		this.IsOrderUser = "1"
	} else {
		this.IsOrderUser = "0"
	}

	_, err = dao.DB.GetMysqlConn().Insert(this)
	if err != nil {
		fmt.Println("insert chat error:", err.Error())
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	return ERRCODE_SUCCESS


}

