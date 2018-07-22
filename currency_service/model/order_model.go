package model

import (
	"bytes"
	"database/sql"
	"digicon/currency_service/conf"
	"digicon/currency_service/dao"
	. "digicon/currency_service/log"
	. "digicon/proto/common"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// 订单表
type Order struct {
	Id          uint64         `xorm:"not null pk autoincr comment('ID')  INT(10)"                 json:"id"`
	OrderId     string         `xorm:"not null pk comment('订单ID') INT(10)"                        json:"order_id"`
	AdId        uint64         `xorm:"not null default 0 comment('广告ID') index INT(10)"           json:"ad_id"`
	AdType      uint32         `xorm:"not null default 0 comment('广告类型:1出售 2购买') TINYINT(1)"  json:"ad_type"`
	Price       int64          `xorm:"not null default 0 comment('价格') BIGINT(64)"                json:"price"`
	Num         int64          `xorm:"not null default 0 comment('数量') BIGINT(64)"                json:"num"`
	TokenId     uint64         `xorm:"not null default 0 comment('货币类型') INT(10)"               json:"token_id"`
	PayId       string         `xorm:"not null default 0 comment('支付类型') VARCHAR(64)"           json:"pay_id"`
	SellId      uint64         `xorm:"not null default 0 comment('卖家id') INT(10)"                json:"sell_id"`
	SellName    string         `xorm:"not null default '' comment('卖家昵称') VARCHAR(64)"          json:"sell_name"`
	BuyId       uint64         `xorm:"not null default 0 comment('买家id') INT(10)"                json:"buy_id"`
	BuyName     string         `xorm:"not null default '' comment('买家昵称') VARCHAR(64)"         json:"buy_name"`
	Fee         int64          `xorm:"not null default 0 comment('手续费用') BIGINT(64)"           json:"fee"`
	States      uint32         `xorm:"not null default 0 comment('订单状态: 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消') TINYINT(1)"   json:"states"`
	PayStatus   uint32         `xorm:"not null default 0 comment('支付状态: 1待支付 2待放行(已支付) 3确认支付(已完成)') TINYINT(1)"  json:"pay_status"`
	CancelType  uint32         `xorm:"not null default 0 comment('取消类型: 1卖方 2 买方') TINYINT(1)"   json:"cancel_type"`
	CreatedTime string         `xorm:"not null comment('创建时间') DATETIME"                       json:"created_time"`
	UpdatedTime string         `xorm:"comment('修改时间')     DATETIME"                           json:"updated_time"`
	ConfirmTime sql.NullString `xorm:"default null comment('确认支付时间')  DATETIME"             json:"confirm_time"`
	ReleaseTime sql.NullString `xorm:"default null comment('放行时间')     DATETIME"              json:"release_time"`
	ExpiryTime  string         `xorm:"comment('过期时间')     DATETIME"                           json:"expiry_time"`
}

//列出订单
func (this *Order) List(Page, PageNum int32,
	AdType, States uint32, Id, Uid uint64,
	TokenId float64, StartTime, EndTime string, o *[]Order) (total int64, rPage int32, rPageNum int32, code int32) {

	engine := dao.DB.GetMysqlConn()
	if Page <= 1 {
		Page = 1
	}
	if PageNum <= 0 {
		PageNum = 10
	}

	query := engine.Desc("id")
	orderModel := new(Order)

	query = query.Where("sell_id=? or buy_id=?", Uid, Uid)

	if States == 0 { // 状态为0，表示已经删除
		return 0, 0, 0, ERRCODE_SUCCESS
	} else if States == 100 { // 默认传递的States
		query = query.Where("states != 0")
	} else {
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
	if StartTime != `` {
		query = query.Where("created_time >= ?", StartTime)
	}
	if EndTime != `` {
		query = query.Where("created_time <= ?", EndTime)
	}

	tmpQuery := *query
	countQuery := &tmpQuery
	err := query.Limit(int(PageNum), (int(Page)-1)*int(PageNum)).Find(o)
	total, _ = countQuery.Count(orderModel)

	if err != nil {
		Log.Errorln(err.Error())
		total = 0
		rPage = 0
		rPageNum = 0
		code = ERRCODE_UNKNOWN
	} else {
		total = total
		rPage = Page
		rPageNum = PageNum
		code = ERRCODE_SUCCESS
	}
	return
}

// 删除订单(将states设置成0)
// id     uint64
// set state = 0
func (this *Order) Delete(Id uint64, updateTimeStr string) (int32, string) {
	var err error
	sql := "UPDATE   `order`   SET   `states`=?, `updated_time`=?  WHERE  `id`=? and `states` != 2 "
	_, err = dao.DB.GetMysqlConn().Exec(sql, 0, updateTimeStr, Id)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN, "delete Error!"
	}
	return ERRCODE_SUCCESS, ""
}

// 取消订单
// set state == 4
// params: id userid, CancelType: 取消类型: 1卖方 2 买方
func (this *Order) Cancel(Id uint64, CancelType uint32, updateTimeStr string) (code int32, msg string) {
	var err error
	sql := "UPDATE   `order`   SET   `states`=? , `cancel_type`=?, `updated_time`=?  WHERE  `id`=?"
	_, err = dao.DB.GetMysqlConn().Exec(sql, 4, CancelType, updateTimeStr, Id)
	if err != nil {
		Log.Errorln(err.Error())
		//fmt.Println(err.Error())
		code = ERRCODE_UNKNOWN
		msg = "cancel Error!"
		//return ERRCODE_UNKNOWN,
	}
	code = ERRCODE_SUCCESS
	msg = ""
	return
}

// 待放行
// set states=2
func (this *Order) Ready(Id uint64, updateTimeStr string) (code int32, msg string) {
	engine := dao.DB.GetMysqlConn()
	var err error
	now := time.Now().Format("2006-01-02 15:04:05")
	sql := "UPDATE   `order`   SET   `states`=?, `updated_time`=?,`confirm_time`=?  WHERE  `id`=?"
	_, err = engine.Exec(sql, 2, updateTimeStr, now, Id)
	if err != nil {
		fmt.Println(err.Error())
		Log.Errorln(err.Error())
		code, msg = ERRCODE_UNKNOWN, "Approved Error!"
		return
	}
	code, msg = ERRCODE_SUCCESS, ""
	return
}

// 添加订单
func (this *Order) Add() (id uint64, code int32) {
	var err error
	/////////////

	engine := dao.DB.GetMysqlConn()

	uCurrency := new(UserCurrency)

	fmt.Println(this.SellId, this.TokenId)
	_, err = engine.Table("user_currency").Where("uid =? and token_id =?", this.SellId, this.TokenId).Get(uCurrency)
	if err != nil {
		Log.Errorln("查询用户余额失败!", err.Error())
		code = ERRCODE_USER_BALANCE
		return
	}

	rate := conf.Cfg.MustValue("rate", "fee_rate")
	rateFloat, _ := strconv.ParseInt(rate, 10, 64)
	freeze := this.Num * (1 + rateFloat)

	//fmt.Println(uCurrency)
	fmt.Println(this.BuyName, freeze, uCurrency.Balance)
	if freeze > uCurrency.Balance {
		Log.Errorln("卖家余额不足!")
		code = ERRCODE_SELLER_LESS
		return
	}

	nTime := time.Now()
	nowTime := nTime.Format("2006-01-02 15:04:05")

	session := engine.NewSession()

	/// 1. 卖家冻结
	sellSql := "update user_currency set  `freeze`=`freeze`+ ?, `balance`=`balance`-?,`version`=`version`+1  WHERE  `uid` = ? and `token_id` = ? and `version`=?"
	sqlRest, err := session.Exec(sellSql, freeze, freeze, this.SellId, this.TokenId, uCurrency.Version) // 卖家 扣除平台费用
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return
	}
	if rst, _ := sqlRest.RowsAffected(); rst == 0 {
		Log.Errorln("冻结失败!")
		session.Rollback()
		code = ERRCODE_ORDER_ERROR
		return
	}
	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return
	}

	/// 2. 记录卖家冻结
	var buffer bytes.Buffer
	buffer.WriteString("insert into user_currency_history ")
	buffer.WriteString("(uid, order_id, token_id, num, fee, operator, address, states, created_time ,updated_time )")
	buffer.WriteString("values (?, ?, ?, ?, ?,  ?, ?, ?, ? , ?)")
	insertSql := buffer.String()
	_, err = session.Table(`user_currency_history`).Exec(insertSql,
		this.SellId, this.OrderId, this.TokenId, this.Num, 0, 5, "", this.States, nowTime, nowTime, // 卖家记录 , 5 为冻结
	)

	if err != nil {
		fmt.Println(err.Error())
		Log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return
	}

	/// 3. 下单，添加订单
	this.ConfirmTime = sql.NullString{String: "", Valid: false}
	this.ReleaseTime = sql.NullString{String: "", Valid: false}
	_, err = session.Table("order").InsertOne(this)

	if err != nil {
		fmt.Println("order error....", err.Error())
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
	code = ERRCODE_SUCCESS
	go CheckOrderExiryTime(id, this.ExpiryTime)
	return
}

// 确认放行(支付完成)
// set state = 3
func (this *Order) Confirm(Id uint64, updateTimeStr string) (int32, string) {
	code, err := this.ConfirmSession(Id, updateTimeStr)
	if err != nil {
		Log.Errorln(err.Error())
		return code, "confirm Error!"
	} else {
		return ERRCODE_SUCCESS, ""
	}
}

// 确认放行事务
func (this *Order) ConfirmSession(Id uint64, updateTimeStr string) (code int32, err error) {
	engine := dao.DB.GetMysqlConn()

	_, err = engine.Table(`order`).Where("id=?", Id).Get(this)
	if err != nil {
		Log.Errorln(err.Error())
		code = ERRCODE_UNKNOWN
		return code, err
	}

	rate := conf.Cfg.MustValue("rate", "fee_rate")
	rateFloat, _ := strconv.ParseInt(rate, 10, 64)

	allNum := this.Num
	rateFee := allNum * rateFloat

	tokens := new(Tokens)
	_, err = engine.Where("id=?", this.TokenId).Get(tokens)
	if err != nil {
		Log.Errorln("获取币种id, token_id 失败")
		code = ERRCODE_UNKNOWN
		return code, err
	}
	tokenName := tokens.Name

	uCurrency := new(UserCurrency)
	_, err = engine.Where("uid =? and token_id =?", this.SellId, this.TokenId).Get(uCurrency)
	if err != nil {
		Log.Errorln("查询用户余额失败!", err.Error())
		code = ERRCODE_USER_BALANCE
		return code, err
	}

	sellNum := allNum + rateFee
	if uCurrency.Freeze < sellNum || uCurrency.Freeze < 0 {
		Log.Println("余额不足!")
		err = errors.New("余额不足")
		code = ERRCODE_SELLER_LESS
		return code, err
	}

	// 事务开始
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()

	if err != nil {
		Log.Println(err.Error())
		session.Rollback()
		code = ERRCODE_TRADE_ERROR
		return
	}

	/////////////////////////////////////////////////////////
	// 1. user_currency扣, 买家减交易金额，卖家加交易金额和减去费用
	////////////////////////////////////////////////////////

	sellSql := "update user_currency set  `freeze`=`freeze` - ?,`version`=`version`+1  WHERE  uid = ? and token_id = ? and version = ?"
	sqlRest, err := session.Exec(sellSql, sellNum, this.SellId, this.TokenId, uCurrency.Version) // 卖家 扣除平台费用
	if err != nil {
		Log.Println(err.Error())
		session.Rollback()
		return
	}
	if rst, _ := sqlRest.RowsAffected(); rst == 0 {
		fmt.Println("卖家扣除失败!", err.Error())
		Log.Errorln("卖家扣除失败!", err.Error())
		session.Rollback()
		code = ERRCODE_TRADE_ERROR
		err = errors.New("卖家扣除失败!")
		return code, err
	}

	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
	}

	has, err := session.Table(`user_currency`).Where("uid=? AND token_id=? ", this.BuyId, this.TokenId).Exist()
	if err != nil {
		Log.Errorln("find user_currency error!")
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return
	}
	if has {
		fmt.Println("has ....")
		_ , err := session.Where("uid=? AND  token_id=? AND version=?",this.BuyId, this.TokenId, uCurrency.Version).Update(&UserCurrency{Balance:allNum})
		if err != nil {
			fmt.Println("insert error!, ", err.Error())
			Log.Println(err.Error())
			session.Rollback()
			code = ERRCODE_TRADE_ERROR
			return code, err
		}
	} else {
		_, err := session.InsertOne(&UserCurrency{Uid:this.BuyId, TokenId: uint32(this.TokenId), TokenName:tokenName, Balance:allNum})
		if err != nil {
			fmt.Println("insert error!, ", err.Error())
			Log.Println(err.Error())
			session.Rollback()
			code = ERRCODE_TRADE_ERROR
			return code, err
		}
	}


	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
	}

	//////////////////////////////////////////////////
	// 2. 插入记录 user_currency_history, 转入，转出，卖家扣费用
	/////////////////////////////////////////////////
	var buffer bytes.Buffer
	buffer.WriteString("insert into user_currency_history ")
	buffer.WriteString("(uid, order_id, token_id, num, fee, operator, address, states,created_time, updated_time)")
	buffer.WriteString("values (?, ?, ?, ?, ?,   ?, ?, ? ,?, ?), (?, ?, ?, ?, ?,   ?, ?, ? , ?, ?)")
	insertSql := buffer.String()

	nowCreate := time.Now().Format("2006-01-02 15:04:05")
	_, err = session.Table(`user_currency_history`).Exec(insertSql,
		this.SellId, this.OrderId, this.TokenId, sellNum, rateFee, 2, "", this.States,nowCreate, updateTimeStr, // 卖家记录 , 2订单转出
		this.BuyId, this.OrderId, this.TokenId, this.Num, 0, 1, "", this.States,nowCreate,  updateTimeStr, // 买家记录 , 1订单转入
	)
	if err != nil {
		fmt.Println("insert into history error: ",  err.Error())
		Log.Println(err.Error())
		session.Rollback()
		code = ERRCODE_TRADE_ERROR
		return code, err
	}
	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
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
		code = ERRCODE_TRADE_ERROR
		return code, err
	}
	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
	}

	//////////////////////////////////////////////////////
	// 4. 更新状态
	//////////////////////////////////////////////////////
	now := time.Now().Format("2006-01-02 15:04:05")
	updateStatesSql := "UPDATE   `order`   SET   `states`=?, `updated_time`=?,`release_time`=?,`fee`=?  WHERE  `id`=?"
	_, err = session.Exec(updateStatesSql, 3, updateTimeStr, now, rateFee, Id)
	if err != nil {
		fmt.Println("update order states error:", err.Error())
		Log.Println(err.Error())
		session.Rollback()
		code = ERRCODE_TRADE_ERROR
		return code, err
	}
	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
	}
	err = session.Commit()
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
	}

	code = ERRCODE_SUCCESS
	return
}

/*
	获取订单付款信息
*/

func (this *Order) GetOrder(Id uint64) (code int32, err error) {
	//var order Order
	engine := dao.DB.GetMysqlConn()
	//order := new(Order)
	_, err = engine.Where("id = ?", Id).Get(this)
	if err != nil {
		Log.Errorln(err.Error())
		code = ERRCODE_ORDER_NOTEXIST
	}
	return
}

///////////////////// func ////////////////
/*

 */
func CheckOrderExiryTime(id uint64, exiryTime string) {
	fmt.Println("go run check order Exiry time ...............................")
	od := new(Order)
	for {
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println(now, exiryTime)
		fmt.Println("id: ", id, " order Exitry time(min): ", getHourDiffer(now, exiryTime))
		if getHourDiffer(now, exiryTime) <= 0 {
			engine := dao.DB.GetMysqlConn()
			_, err := engine.Where("id=?", id).Get(od)
			if err != nil {
				Log.Errorln("get order states error!")
			} else {
				if od.States == 0 || od.States == 2 || od.States == 3 || od.States == 4 { // 0删除 2待放行(已支付) 3确认支付(已完成) 4取消
					fmt.Println("break: od stats", od.States)
					break
				}
			}
			_, err = engine.Where("id = ?", id).Update(&Order{States: 4, UpdatedTime: now})
			if err != nil {
				Log.Errorln("order id exiry time out update status = 4 error! id:", id, err.Error())
			}
			break
		}
		time.Sleep(2 * 60 * time.Second)
	}
	fmt.Println(id, " break .................")
	return
}

//获取相差时间
func getHourDiffer(start_time, end_time string) int64 {
	var minute int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		minute = diff / 60            //3600
		return minute
	} else {
		return minute
	}
}
