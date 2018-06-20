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
	OrderId     string       `xorm:"not null pk comment('订单ID') INT(10)"   json:"order_id"`  // hash( type_id, 6( user_id, + 时间秒）
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
	AdType, States uint32, Id uint64,
	TokenId float64, CreatedTime string, o *[]Order) (int64,int32, int32, int32) {

	engine := dao.DB.GetMysqlConn()
	if Page <= 1 {
		Page = 1
	}
	if PageNum <= 0{
		PageNum = 10
	}

	query := engine.Desc("id")

	orderModel := new(Order)
	fmt.Println("States:", States)

	if States == 0 {                            // 状态为0，表示已经删除
		return 0, 0, 0, ERRCODE_SUCCESS
	}else if States == 100{                     // 默认传递的States
		query = query.Where("states != 0")
	}else{
		query = query.Where("states = ?", States)
	}

	fmt.Println("id:", Id)
	if Id != 0 {
		query = query.Where("id = ?", Id)
	}
	if AdType != 0 {
		query = query.Where("ad_type = ?", AdType)
	}
	if TokenId != 0 {
		query = query.Where("token_id = ?", TokenId)
	}

	if CreatedTime !=  ``  {
		query = query.Where("created_time = ?", CreatedTime)
	}
	tmpQuery := *query
	countQuery := &tmpQuery
	err := query.Limit(int(PageNum), (int(Page) - 1) * int(PageNum)).Find(o)
	total, _:= countQuery.Count(orderModel)

	if err != nil {
		Log.Errorln(err.Error())
		return 0,0, 0, ERRCODE_UNKNOWN
	}
	return total, Page, PageNum, ERRCODE_SUCCESS
}


// 删除订单(将states设置成0)
// id     uint64
// set state = 0
func (this *Order) Delete(Id uint64,  updateTimeStr string) (int32, string){
	var err error
	sql := "UPDATE   `order`   SET   `states`=?, `updated_time`=?  WHERE  `id`=? and `states` != 2 "
	_, err = dao.DB.GetMysqlConn().Exec(sql,0, updateTimeStr, Id)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN, err.Error()
	}
	return ERRCODE_SUCCESS, ""
}



// 取消订单
// set state == 4
// params: id userid, CancelType: 取消类型: 1卖方 2 买方
func (this *Order) Cancel(Id uint64, CancelType uint32,  updateTimeStr string) (int32,string ){
	var err error
	sql := "UPDATE   `order`   SET   `states`=? , `cancel_type`=?, `updated_time`=?  WHERE  `id`=?"
	_, err = dao.DB.GetMysqlConn().Exec(sql, 4,CancelType, updateTimeStr ,Id)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN, err.Error()
	}
	return ERRCODE_SUCCESS, ""
}

// 确认放行(支付完成)
// set state = 3
func (this *Order) Confirm(Id uint64, updateTimeStr string) (int32, string){
	//  去调用

	var err error
	sql := "UPDATE   `order`   SET   `states`=?, `updated_time`=?  WHERE  `id`=?"
	_, err = dao.DB.GetMysqlConn().Exec(sql, 3, updateTimeStr,Id)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN, err.Error()
	}
	return ERRCODE_SUCCESS, ""
}


// 待放行
// set states=2
func (this *Order)Ready(Id uint64,  updateTimeStr string) (int32, string) {
	var err error
	sql := "UPDATE   `order`   SET   `states`=?, `updated_time`=?  WHERE  `id`=?"
	_, err = dao.DB.GetMysqlConn().Exec(sql, 2, updateTimeStr,Id)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN, err.Error()
	}
	return ERRCODE_SUCCESS, ""
}


//添加订单
func (this *Order) Add() (id uint64, code int32) {
	_, err := dao.DB.GetMysqlConn().Insert(this)
	if err != nil {
		Log.Errorln(err.Error())
		code = ERRCODE_UNKNOWN
	}else{
		id = this.Id
	}
	code =  ERRCODE_SUCCESS
	return
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


