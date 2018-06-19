package model

import (
	"digicon/currency_service/dao"
	. "digicon/currency_service/log"
	. "digicon/proto/common"
	"fmt"
)

// 订单表
type Order struct {
	Id          uint64       `xorm:"not null pk autoincr comment('ID')  INT(10)"  json:"id"`
	OrderId     uint64       `xorm:"not null pk comment('订单ID') INT(10)"   json:"order_id"`  // hash( type_id, 6( user_id, + 时间秒）
	AdId        uint64       `xorm:"not null default 0 comment('广告ID') index INT(10)"  json:"ad_id"`
	AdType      uint32       `xorm:"not null default 0 comment('广告类型:1出售 2购买') TINYINT(1)"  json:"ad_type"`
	Price       float64      `xorm:"not null default 0.000000 comment('价格') DECIMAL(20,6)"   json:"price"`
	Num         float64      `xorm:"not null default 0.000000 comment('数量') DECIMAL(20,6)"   json:"num"`
	TokenId     uint64       `xorm:"not null default 0 comment('货币类型') INT(10)"       json:"token_id"`
	PayId       uint64       `xorm:"not null default 0 comment('支付类型') INT(10)"       json:"pay_id"`
	SellId      uint64       `xorm:"not null default 0 comment('卖家id') INT(10)"         json:"sell_id"`
	SellName    string       `xorm:"not null default '' comment('卖家昵称') VARCHAR(64)"  json:"sell_name"`
	BuyId       uint64       `xorm:"not null default 0 comment('买家id') INT(10)"    json:"buy_id"`
	BuyName     string       `xorm:"not null default '' comment('买家昵称') VARCHAR(64)"   json:"buy_name"`
	Fee         float64      `xorm:"not null default 0.000000 comment('手续费用') DECIMAL(20,6)"  json:"fee"`
	States      uint32       `xorm:"not null default 0 comment('订单状态: 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消') TINYINT(1)"   json:"states"`
	PayStatus   uint32       `xorm:"not null default 0 comment('支付状态: 1待支付 2待放行(已支付) 3确认支付(已完成)') TINYINT(1)"  json:"pay_status"`
	CancelType  uint32       `xorm:"not null default 0 comment('取消类型: 1卖方 2 买方') TINYINT(1)"   json:"cancel_type"`
	CreatedTime string       `xorm:"not null comment('创建时间') DATETIME"  json:"created_time"`
	UpdatedTime string       `xorm:"comment('修改时间') DATETIME"    json:"updated_time"`

}



//列出订单
func (this *Order)  List(Page, PageNum int32,
	AdType, States uint32,
	TokenId float64, CreatedTime string, o *[]Order) (int64,int32, int32, int32) {

	engine := dao.DB.GetMysqlConn()
	if Page <= 0 {
		Page = 0
	}
	if PageNum <= 0{
		PageNum = 10
	}

	displaySql := "states != 0"        // 0 为已删除，不显示

	query := engine.Where(displaySql)
	orderModel := new(Order)
	total, _:= query.Count(orderModel)

	if AdType != 0 {
		query = query.Where("ad_type = ?", AdType)
	}
	if TokenId != 0 {
		query = query.Where("token_id = ?", TokenId)
	}
	if States != 0 {
		fmt.Println("States:", States)
		query = query.Where("states = ?", States)
	}
	if CreatedTime !=  ``  {
		query = query.Where("created_time = ?", CreatedTime)
	}
	err := query.Where(displaySql).Limit(int(PageNum), int(Page)).Find(o)
	if err != nil {
		Log.Errorln(err.Error())
		return 0,0, 0, ERRCODE_UNKNOWN
	}
	return total, Page, PageNum, ERRCODE_SUCCESS
}


// 删除订单(将states设置成0)
// id     uint64
// set state = 0
func (this *Order) Delete(Id uint64) int32 {
	var err error
	sql := "UPDATE   `order`   SET   `states`=? WHERE  `id`=?"
	_, err = dao.DB.GetMysqlConn().Exec(sql,0, Id)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	return ERRCODE_SUCCESS
}



// 取消订单
// set state == 4
func (this *Order) Cancel(Id uint64, CancelType uint32) int32 {
	var err error
	sql := "UPDATE   `order`   SET   `states`=? , `cancel_type`=?  WHERE  `id`=?"
	_, err = dao.DB.GetMysqlConn().Exec(sql, 4,CancelType, Id)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	return ERRCODE_SUCCESS
}

// 确认放行(支付完成)
// set state = 3
func (this *Order) Confirm(Id uint64) int32{
	var err error
	sql := "UPDATE   `order`   SET   `states`=?  WHERE  `id`=?"
	_, err = dao.DB.GetMysqlConn().Exec(sql, 3, Id)
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
func (this *Order) Update(Id uint64) int32 {
	_, err := dao.DB.GetMysqlConn().Id(Id).Update(this)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	return ERRCODE_SUCCESS
}


