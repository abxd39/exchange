package model

import (
	"digicon/currency_service/dao"
	. "digicon/currency_service/log"
	. "digicon/proto/common"
)

// 订单表
type Order struct {
	OrderId     uint64       `xorm:"not null pk autoincr comment('订单ID') INT(10)"`
	AdId        uint64       `xorm:"not null default 0 comment('广告ID') index INT(10)"`
	AdType      uint32       `xorm:"not null default 0 comment('广告类型:1出售 2购买') TINYINT(1)"`
	Price       float64      `xorm:"not null default 0.000000 comment('价格') DECIMAL(20,6)"`
	Num         float64      `xorm:"not null default 0.000000 comment('数量') DECIMAL(20,6)"`
	TokenId     uint64       `xorm:"not null default 0 comment('货币类型') INT(10)"`
	PayId       uint64       `xorm:"not null default 0 comment('支付类型') INT(10)"`
	SellId      uint64       `xorm:"not null default 0 comment('卖家id') INT(10)"`
	SellName    string       `xorm:"not null default '' comment('卖家昵称') VARCHAR(64)"`
	BuyId       uint64       `xorm:"not null default 0 comment('买家id') INT(10)"`
	BuyName     string       `xorm:"not null default '' comment('买家昵称') VARCHAR(64)"`
	Fee         float64      `xorm:"not null default 0.000000 comment('手续费用') DECIMAL(20,6)"`
	States      uint32       `xorm:"not null default 0 comment('订单状态: 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消') TINYINT(1)"`
	PayStatus   uint32       `xorm:"not null default 0 comment('支付状态: 1待支付 2待放行(已支付) 3确认支付(已完成)') TINYINT(1)"`
	CancelType  uint32       `xorm:"not null default 0 comment('取消类型: 1卖方 2 买方') TINYINT(1)"`
	CreatedTime string       `xorm:"not null comment('创建时间') DATETIME"`
	UpdatedTime string       `xorm:"comment('修改时间') DATETIME"`
}



//列出订单
func (this *Order)  List(startRow, endRow int32, o *[]Order) int32 {
	 engine := dao.DB.GetMysqlConn()
	 err := engine.Limit(int(endRow), int(startRow)).Find(o)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	return ERRCODE_SUCCESS
}


//添加订单
func (this *Order) Add() int32 {
	_, err := dao.DB.GetMysqlConn().Insert(this)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	return ERRCODE_SUCCESS
}


//更新订单
func (this *Order) Update() int32 {
	_, err := dao.DB.GetMysqlConn().Id(this.OrderId).Update(this)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	return ERRCODE_SUCCESS
}

func (this *Order) Delete() int32 {
	//_, err := dao.DB.GetMysqlConn().Id(this.Id).Delete(this)
	//if err != nil {
	//	Log.Errorln(err.Error())
	//	return ERRCODE_UNKNOWN
	//}

	return ERRCODE_SUCCESS
}




