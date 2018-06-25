package model

import (
	"digicon/currency_service/conf"
	"digicon/currency_service/dao"
	. "digicon/currency_service/log"
	. "digicon/proto/common"
	"strconv"
	"bytes"
	"time"
)


// 订单表
type Order struct {
	Id          uint64       `xorm:"not null pk autoincr comment('ID')  INT(10)"  json:"id"`
	OrderId     string       `xorm:"not null pk comment('订单ID') INT(10)"   json:"order_id"`  // hash( type_id, 6( user_id, + 时间秒）
	AdId        uint64       `xorm:"not null default 0 comment('广告ID') index INT(10)"  json:"ad_id"`
	AdType      uint32       `xorm:"not null default 0 comment('广告类型:1出售 2购买') TINYINT(1)"  json:"ad_type"`
	Price       int64        `xorm:"not null default 0 comment('价格') BIGINT(64)"   json:"price"`
	Num         int64        `xorm:"not null default 0 comment('数量') BIGINT(64)"   json:"num"`
	TokenId     uint64       `xorm:"not null default 0 comment('货币类型') INT(10)"       json:"token_id"`
	PayId       uint64       `xorm:"not null default 0 comment('支付类型') INT(10)"       json:"pay_id"`
	SellId      uint64       `xorm:"not null default 0 comment('卖家id') INT(10)"         json:"sell_id"`
	SellName    string       `xorm:"not null default '' comment('卖家昵称') VARCHAR(64)"  json:"sell_name"`
	BuyId       uint64       `xorm:"not null default 0 comment('买家id') INT(10)"    json:"buy_id"`
	BuyName     string       `xorm:"not null default '' comment('买家昵称') VARCHAR(64)"   json:"buy_name"`
	Fee         int64        `xorm:"not null default 0 comment('手续费用') BIGINT(64)"  json:"fee"`
	States      uint32       `xorm:"not null default 0 comment('订单状态: 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消') TINYINT(1)"   json:"states"`
	PayStatus   uint32       `xorm:"not null default 0 comment('支付状态: 1待支付 2待放行(已支付) 3确认支付(已完成)') TINYINT(1)"  json:"pay_status"`
	CancelType  uint32       `xorm:"not null default 0 comment('取消类型: 1卖方 2 买方') TINYINT(1)"   json:"cancel_type"`
	CreatedTime string       `xorm:"not null comment('创建时间') DATETIME"  json:"created_time"`
	UpdatedTime string       `xorm:"comment('修改时间') DATETIME"    json:"updated_time"`
}




//列出订单
func (this *Order)  List(Page, PageNum int32,
	AdType, States uint32, Id uint64,
	TokenId float64, StartTime, EndTime string, o *[]Order) (int64,int32, int32, int32) {

	engine := dao.DB.GetMysqlConn()
	if Page <= 1 {
		Page = 1
	}
	if PageNum <= 0{
		PageNum = 10
	}

	query := engine.Desc("id")
	orderModel := new(Order)

	if States == 0 {                            // 状态为0，表示已经删除
		return 0, 0, 0, ERRCODE_SUCCESS
	}else if States == 100{                     // 默认传递的States
		query = query.Where("states != 0")
	}else{
		query = query.Where("states = ?", States)
	}

	if Id != 0 {
		query = query.Where("id = ?", Id)
	}
	if AdType != 0 {
		query = query.Where("ad_type = ?", AdType)
	}
	if TokenId != 0 {
		query = query.Where("token_id = ?", TokenId)
	}
	//fmt.Println(StartTime, EndTime)
	if StartTime  !=  ``   {
		query = query.Where("created_time >= ?", StartTime)
	}
	if EndTime != ``{
		query = query.Where("created_time <= ?", EndTime)
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
		return ERRCODE_UNKNOWN, "delete Error!"
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
		//fmt.Println(err.Error())
		return ERRCODE_UNKNOWN, "cancel Error!"
	}
	return ERRCODE_SUCCESS, ""
}



// 待放行
// set states=2
func (this *Order)Ready(Id uint64,  updateTimeStr string) (int32, string) {
	engine := dao.DB.GetMysqlConn()
	var err error
	sql := "UPDATE   `order`   SET   `states`=?, `updated_time`=?  WHERE  `id`=?"
	_, err = engine.Exec(sql, 2, updateTimeStr,Id)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN, "Approved Error!"
	}
	return ERRCODE_SUCCESS, ""
}


// 添加订单
func (this *Order) Add() (id uint64, code int32) {
	var err error
	engine := dao.DB.GetMysqlConn()

	uCurrency := new(UserCurrency)
	engine.Table("user_currency").Where("uid =? and token_id =?", this.SellId, this.TokenId).Get(uCurrency)

	rate := conf.Cfg.MustValue("rate", "fee_rate")
	rateFloat, _ := strconv.ParseInt(rate, 10, 64)
	freeze := this.Num * (1 + rateFloat)

	if freeze > uCurrency.Balance {
		Log.Errorln("卖家余额不足!")
		code = ERRCODE_SELLER_LESS
		return
	}

	nowTime := time.Now().Format("2006-01-02 15:04:05")
	session := engine.NewSession()

	/// 1. 卖家冻结
	sellSql := "update user_currency set  `freeze`=`freeze`+ ?, `balance`=`balance`-?,`version`=`version`+1  WHERE  uid = ? and token_id = ? and version = ?"
	_, err = session.Exec(sellSql, freeze, freeze, this.SellId, this.TokenId,  uCurrency.Version )         // 卖家 扣除平台费用
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return
	}

	/// 2. 记录卖家冻结
	var buffer bytes.Buffer
	buffer.WriteString("insert into user_currency_history ")
	buffer.WriteString("(uid, order_id, token_id, num, fee, operator, address, states, created_time ,updated_time)")
	buffer.WriteString("values (?, ?, ?, ?, ?,  ?, ?, ?, ? , ?)")
	insertSql := buffer.String()
	_, err = session.Table(`user_currency_history`).Exec(insertSql,
		this.SellId,this.OrderId, this.TokenId,  this.Num , 0 , 5, "", this.States, nowTime, nowTime,      // 卖家记录 , 5 为冻结
	)

	if err != nil{
		Log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return
	}

	/// 3. 下单，添加订单
	_, err = session.Table("order").InsertOne(this)

	if err != nil {
		Log.Errorln(err.Error())
		code = ERRCODE_UNKNOWN
		session.Rollback()
		return
	}

	err = session.Commit()
	if err != nil {
		Log.Errorln(err.Error())
		code = ERRCODE_UNKNOWN
		session.Rollback()
		return
	}
	id = this.Id
	code =  ERRCODE_SUCCESS
	return
}



// 确认放行(支付完成)
// set state = 3
func (this *Order) Confirm(Id uint64, updateTimeStr string) (int32, string){
	code, err := this.ConfirmSession(Id, updateTimeStr)
	if err != nil {
		Log.Errorln(err.Error())
		return code, "confirm Error!"
	}else{
		return ERRCODE_SUCCESS, ""
	}
}





// 确认放行事务
func (this *Order) ConfirmSession (Id uint64, updateTimeStr string) (code int32, err error) {
	engine := dao.DB.GetMysqlConn()

	engine.Table(`order`).Where("id=?", Id).Get(this)

	rate := conf.Cfg.MustValue("rate", "fee_rate")
	rateFloat, _ := strconv.ParseInt(rate, 10, 64)

	allNum := this.Num
	rateFee := allNum * rateFloat

	tokens := new(Tokens)
	engine.Where("id=?", this.TokenId).Get(tokens)
	tokenName := tokens.Name

	uCurrency := new(UserCurrency)
	engine.Where("uid =? and token_id =?", this.SellId, this.TokenId).Get(uCurrency)

	sellNum := allNum + rateFee
	if  uCurrency.Freeze < sellNum || uCurrency.Freeze < 0 {
		Log.Println("余额不足!")
		code = ERRCODE_SELLER_LESS
		return
	}

	// 事务开始
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()

	if err != nil {
		Log.Println(err.Error())
		session.Rollback()
		return
	}

	/////////////////////////////////////////////////////////
	// 1. user_currency扣, 买家减交易金额，卖家加交易金额和减去费用
	////////////////////////////////////////////////////////

	sellSql := "update user_currency set  `freeze`=`freeze` - ?,`version`=`version`+1  WHERE  uid = ? and token_id = ? and version = ?"
	_, err = session.Exec(sellSql, sellNum, this.SellId, this.TokenId,  uCurrency.Version )         // 卖家 扣除平台费用
	if err != nil {
		Log.Println(err.Error())
		session.Rollback()
		return
	}

	buySql := "INSERT  INTO user_currency(`uid`, `token_id`, `token_name`,  `balance`)  values(?, ?, ?, ? ) ON DUPLICATE  KEY  UPDATE  `balance`=`balance`+?  WHERE `uid`=? and token_id=? and version=?"
	_, err = session.Exec(buySql, this.BuyId, this.TokenId, tokenName ,  allNum, allNum, this.BuyId, this.TokenId, uCurrency.Version)

	if err != nil {
		Log.Println(err.Error())
		session.Rollback()
		return
	}

	//////////////////////////////////////////////////
	// 2. 插入记录 user_currency_history, 转入，转出，卖家扣费用
	/////////////////////////////////////////////////
	var buffer bytes.Buffer
	buffer.WriteString("insert into user_currency_history ")
	buffer.WriteString("(uid, order_id, token_id, num, fee, operator, address, states, updated_time)")
	buffer.WriteString("values (?, ?, ?, ?, ?,   ?, ?, ? ,  ?), (?, ?, ?, ?, ?,   ?, ?, ? ,  ?)")
	insertSql := buffer.String()

	_, err = session.Table(`user_currency_history`).Exec(insertSql,
		this.SellId,this.OrderId, this.TokenId, sellNum , rateFee , 2, "", this.States,  updateTimeStr,   // 卖家记录 , 2订单转出
		this.BuyId, this.OrderId, this.TokenId, this.Num, 0       , 1, "", this.States,  updateTimeStr,   // 买家记录 , 1订单转入
		)
	if err != nil{
		Log.Println(err.Error())
		session.Rollback()
		return
	}

	////////////////////////////////////////////////////////
	// 3. 统计表 加 1
	////////////////////////////////////////////////////////
	//countAddOneSql := "INSERT INTO `user_currency_count` (uid,orders, success, good) values(?,?,?,?) ON DUPLICATE KEY UPDATE `success` = `success`+1"
	//_, err = session.Exec(countAddOneSql, this.SellId, 1, 1, 100.0)
	countAddOneSql := "UPDATE  `user_currency_count` set `success` = `success` + 1  where uid=? "
	_, err = session.Exec(countAddOneSql, this.SellId)
	if err != nil {
		Log.Println(err.Error())
		session.Rollback()
		return
	}


	//////////////////////////////////////////////////////
	// 4. 更新状态
	//////////////////////////////////////////////////////
	updateStatesSql := "UPDATE   `order`   SET   `states`=?, `updated_time`=?, `fee`=?  WHERE  `id`=?"
	_, err = session.Exec(updateStatesSql,3, updateTimeStr, rateFee, Id )
	if err != nil {
		Log.Println(err.Error())
		session.Rollback()
		return
	}
	err = session.Commit()
	code = ERRCODE_SUCCESS
	return
}
