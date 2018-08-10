package model

import (
	"bytes"
	"database/sql"
	"digicon/currency_service/conf"
	"digicon/currency_service/dao"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
	. "digicon/proto/common"
	"digicon/common/convert"
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

	NumTotalPrice  int64       `xorm:"default 0 comment('后台需要的数量总价格') BIGINT(64)"             json:"num_total_price"`
	FeePrice       int64       `xorm:"default 0 comment('后台需要的计算出费用价格') BIGINT(64)"          json:"fee_price"`

}


type OrderModel struct {
	Order         `xorm:"extends"`
	TokenName     string        `xorm:"VARBINARY(36)"  form:"token_name" json:"token_name"`
}



type BuyTotal struct {
	BuyTotalAll       int64    `json:"buy_total_all"`
	Price             int64    `json:"price"`
}

type SellTotal struct {
	SellTotalAll      int64    `json:"sell_total_all"`
	Price             int64    `json:"price"`
}

type SumBuyTotal struct {
	BuyTotalAll       int64    `json:"buy_total_all"`
	BuyTotalAllCny    int64    `json:"buy_total_all_cny"`
}

type SumSellTotal struct {
	SellTotalAll      int64    `json:"sell_total_all"`
	SellTotalAllCny   int64    `json:"sell_total_all_cny"`
}


/*
 ====================  sum buy total all =================
*/
func (this *Order)  GetCurDaySum(tokenId uint32, startTime, endTime string) (sumBuy SumBuyTotal, sumSell SumSellTotal, err error) {
	sql := "select sum(num) as  buy_total_all,price  from `order`  where  token_id =? and ad_type=? and states=3 and created_time >= ? and created_time <= ?"
	engine := dao.DB.GetMysqlConn()
	var btotal BuyTotal
	_, err = engine.Table("order").SQL(sql, tokenId, SellType, startTime, endTime).Get(&btotal)

	sumBuy.BuyTotalAll = btotal.BuyTotalAll
	sumBuy.BuyTotalAllCny = convert.Int64MulInt64By8Bit(btotal.BuyTotalAll, btotal.Price)


	buysql := "select sum(num) as  sell_total_all, price  from `order` where  token_id =? and ad_type=? and states=3 and created_time >= ? and created_time <= ?"
	var stotal SellTotal
	_, err = engine.Table("order").SQL(buysql, tokenId, BuyType, startTime, endTime).Get(&stotal)

	log.Printf("tokenId: %v, stotal: %v \n", tokenId,  stotal)

	sumSell.SellTotalAll = stotal.SellTotalAll
	sumSell.SellTotalAllCny = convert.Int64MulInt64By8Bit(stotal.SellTotalAll, stotal.Price)

	log.Printf("tokenId: %v, sum sell: %v \n", tokenId, sumSell)
	return

}



/*
====================== sum end =======================
*/



//列出订单
func (this *Order) List(Page, PageNum int32,
	AdType, States uint32, Id, Uid uint64,
	TokenId float64, StartTime, EndTime string, o *[]OrderModel) (total int64, rPage int32, rPageNum int32, code int32) {

	engine := dao.DB.GetMysqlConn()
	if Page <= 1 {
		Page = 1
	}
	if PageNum <= 0 {
		PageNum = 10
	}

	orderModel := new(Order)

	query := engine.Table("order").Join("LEFT", "ads", "ads.id=order.ad_id").Desc("order.id")
	query = query.Where("order.sell_id=? or order.buy_id=?", Uid, Uid)

	if States == 0 { // 状态为0，表示已经删除
		return 0, 0, 0, ERRCODE_SUCCESS
	} else if States == 100 { // 默认传递的States
		query = query.Where("order.states != 0")
	} else {
		query = query.Where("order.states = ?", States)
	}

	if Id != 0 {
		query = query.Where("order.id = ?", Id)
	}
	if AdType != 0 {
		query = query.Where("order.ad_type = ?", AdType)
	}
	if TokenId != 0 {
		query = query.Where("order.token_id = ?", TokenId)
	}
	//fmt.Println(StartTime, EndTime)
	if StartTime != "" && EndTime != "" {
		query = query.Where("order.created_time >= ? AND order.created_time <= ?", StartTime, EndTime)
	}

	tmpQuery := *query
	countQuery := &tmpQuery
	err := query.Limit(int(PageNum), (int(Page)-1)*int(PageNum)).Find(o)
	total, _ = countQuery.Count(orderModel)

	if err != nil {
		log.Errorln(err.Error())
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
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN, "delete Error!"
	}
	return ERRCODE_SUCCESS, ""
}

// 取消订单
// set state == 4
// params: id userid, CancelType: 取消类型: 1卖方 2 买方
func (this *Order) Cancel(Id uint64, CancelType uint32, updateTimeStr string, uid int32) (code int32, msg string) {
	engine := dao.DB.GetMysqlConn()
	var err error

	tmpOrder := new(Order)
	code, err = tmpOrder.GetOrder(Id)
	if err != nil {
		msg := "order not exists!"
		return ERRCODE_ORDER_NOTEXIST, msg
	}
	if tmpOrder.States == 3 {
		msg := "订单已完成!"
		return ERRCODE_TRADE_HAS_COMPLETED, msg
	}


	uCurrencyCount := new(UserCurrencyCount)
	has, err := uCurrencyCount.CheckUserCurrencyCountExists(uint64(uid))
	if err != nil {
		log.Error(err.Error())
	}
	if has {
		cancelSql := "UPDATE  `user_currency_count` set `cancel` = `cancel` + 1, `orders`=`orders`+1  where uid=? "
		_, err = engine.Exec(cancelSql, uid)
	} else {
		insertSql := "INSERT INTO `user_currency_count` (uid,orders, cancel, good) values(?,?,?, ?)"
		_, err = dao.DB.GetMysqlConn().Exec(insertSql, uid, 1, 1, 100)
	}

	err  = CancelAction(Id, CancelType)


	if err != nil {
		code = ERRCODE_UNKNOWN
		msg = "cancel error!"
	}else{
		code = ERRCODE_SUCCESS
	}
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
		log.Errorln(err.Error())
		code, msg = ERRCODE_UNKNOWN, "Approved Error!"
		return
	}
	code, msg = ERRCODE_SUCCESS, ""
	return
}

// 添加订单
func (this *Order) Add(curUId int32) (id uint64, code int32) {
	var err error
	/////////////

	engine := dao.DB.GetMysqlConn()

	uCurrency := new(UserCurrency)
	fmt.Println(this.SellId, this.TokenId)


	_, err = engine.Table("user_currency").Where("uid =? and token_id =?", this.SellId, this.TokenId).Get(uCurrency)
	if err != nil {
		log.Errorln("查询用户余额失败!", err.Error())
		code = ERRCODE_USER_BALANCE
		return
	}

	rate := conf.Cfg.MustValue("rate", "fee_rate")
	rateFloat, _ := strconv.ParseInt(rate, 10, 64)
	freeze := this.Num * (1 + rateFloat)

	//fmt.Println(uCurrency)
	fmt.Println(this.BuyName, freeze, uCurrency.Balance)
	if freeze > uCurrency.Balance {
		log.Errorln("卖家余额不足!")
		code = ERRCODE_SELLER_LESS
		return
	}

	nTime := time.Now()
	nowTime := nTime.Format("2006-01-02 15:04:05")

	/// 0
	adsM := new(Ads).Get(this.AdId)
	if adsM.States == 0 {
		msg := "订单已下架"
		log.Println(msg)
		code = ERRCODE_ADS_NOTEXIST
	}

	if adsM.Num < uint64(this.Num) {
		msg := "下单失败,购买的数量大于订单的数量!"
		//fmt.Println(msg)
		log.Println(msg)
		code = ERRCODE_TRADE_ERROR_ADS_NUM
		return
	}

	curCnyPrice := convert.Int64MulInt64By8Bit(this.Num, this.Price)

	int64MinLimit := convert.Float64ToInt64By8Bit(float64(adsM.MinLimit))
	fmt.Println(int64MinLimit, curCnyPrice)
	if  curCnyPrice  < int64MinLimit{
		msg := "下单失败,买价小于允许的最小价格!"
		log.Println(msg)
		code = ERRCODE_TRADE_LOWER_PRICE
		return
	}

	int64MaxLimit := convert.Float64ToInt64By8Bit(float64(adsM.MaxLimit))
	fmt.Println(int64MaxLimit, curCnyPrice)
	if  curCnyPrice > int64MaxLimit {
		msg := "下单失败,买价大于允许的最大价格!"
		log.Println(msg)
		code = ERRCODE_TRADE_LARGE_PRICE
		return
	}



	session := engine.NewSession()
	/// 1. 卖家冻结
	balanceCny := convert.Int64MulInt64By8Bit(uCurrency.Balance, this.Price)
	freezeCny := convert.Int64MulInt64By8Bit(freeze, this.Price)
	var newBalanceCny int64
	if balanceCny >= freezeCny{
		newBalanceCny = balanceCny - freezeCny
	}else{
		newBalanceCny = 0
	}
	sellSql := "update user_currency set  `freeze`=`freeze`+ ?, `balance`=`balance`-?,`balance_cny`=?, `freeze_cny`=`freeze_cny`+? , `version`=`version`+1  WHERE  `uid` = ? and `token_id` = ? and `version`=?"
	sqlRest, err := session.Exec(sellSql, freeze, freeze, freezeCny,  newBalanceCny,this.SellId, this.TokenId, uCurrency.Version) // 卖家 扣除平台费用
	//sellSql := "update user_currency set  `freeze`=`freeze`+ ?, `balance`=`balance`-?,`version`=`version`+1  WHERE  `uid` = ? and `token_id` = ? and `version`=?"
	//sqlRest, err := session.Exec(sellSql, freeze, freeze, this.SellId, this.TokenId, uCurrency.Version) // 卖家 扣除平台费用
	if err != nil {
		log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return
	}
	if rst, _ := sqlRest.RowsAffected(); rst == 0 {
		log.Errorln("冻结失败!")
		session.Rollback()
		code = ERRCODE_ORDER_ERROR
		return
	}
	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return
	}

	/// 2. 记录卖家冻结
	var buffer bytes.Buffer
	buffer.WriteString("insert into user_currency_history ")
	buffer.WriteString("(uid, trade_uid, order_id, token_id, num, fee, surplus, operator, address, states, created_time ,updated_time )")
	buffer.WriteString("values (?, ? ,?, ?, ?, ?,   ?, ?, ?, ?, ? , ?)")
	insertSql := buffer.String()

	_, err = session.Table(`user_currency_history`).Exec(insertSql,
		this.SellId, this.BuyId, this.OrderId, this.TokenId, this.Num, 0, uCurrency.Balance, 5, "", this.States, nowTime, nowTime, // 卖家记录 , 5 为冻结
	)

	if err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
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
		log.Errorln(err.Error())
		code = ERRCODE_UNKNOWN
		session.Rollback()
		return
	}

	//////////////////////////////////////////////////////
	// 4. 广告个数减去相应的数量
	//////////////////////////////////////////////////////

	updateAdsSql := "update `ads` set `num`=`num`-?, `updated_time`=?  WHERE `id`=? "
	_,err = session.Exec(updateAdsSql, freeze, nowTime, this.AdId)
	if err != nil {
		fmt.Println("update ads num states error:", err.Error())
		log.Println(err.Error())
		session.Rollback()
		code = ERRCODE_TRADE_ERROR_ADS_NUM
		return
	}



	/// 4. 自动回复的消息
	fmt.Println("ads reply: ", adsM.Reply)
	if adsM.Reply != ""{
		replaySql := "insert into chats (order_id, is_order_user, uid, uname, content, states, created_time) values (?,?,?,?, ?,?,?)"
		var isOrderUser  int
		if adsM.Uid ==  uint64(curUId) {
			isOrderUser = 1
		}else{
			isOrderUser = 0
		}
		var uname string
		if adsM.Uid == this.SellId  {
			uname = this.SellName
		}else{
			uname = this.BuyName
		}
		_, err = session.Exec(replaySql, this.OrderId, isOrderUser, adsM.Uid, uname ,adsM.Reply, 1, nowTime)
		if err != nil {
			log.Error(err.Error())
		}
	}


	err = session.Commit()
	if err != nil {
		log.Errorln(err.Error())
		code = ERRCODE_UNKNOWN
		session.Rollback()
		return
	}

	id = this.Id
	code = ERRCODE_SUCCESS

	/////    检查是否超时
	go CheckOrderExiryTime(id, this.ExpiryTime)
	/////    检查广告是否需要下架
	go AdsAutoDownline(adsM.Id)

	return
}

// 确认放行(支付完成)
// set state = 3
func (this *Order) Confirm(Id uint64, updateTimeStr string, uid int32) (int32, string) {
	code, err := this.ConfirmSession(Id, updateTimeStr, uid)
	if err != nil {
		log.Errorln(err.Error())
		return code, "confirm Error!"
	} else {
		return ERRCODE_SUCCESS, ""
	}
}


// 确认放行事务
func (this *Order) ConfirmSession(Id uint64, updateTimeStr string, uid int32) (code int32, err error) {
	engine := dao.DB.GetMysqlConn()

	_, err = engine.Table(`order`).Where("id=?", Id).Get(this)
	if err != nil {
		log.Errorln(err.Error())
		code = ERRCODE_UNKNOWN
		return code, err
	}

	rate := conf.Cfg.MustValue("rate", "fee_rate")
	rateFloat, _ := strconv.ParseInt(rate, 10, 64)

	allNum := this.Num
	rateFee := allNum * rateFloat

	//tokens := new(Tokens)

	//_, err = engine.Where("id=?", this.TokenId).Get(tokens)
	tokenModel := new(CommonTokens)
	tokens := tokenModel.Get(uint32(this.TokenId), "")

	if err != nil {
		log.Errorln("获取币种id, token_id 失败")
		code = ERRCODE_UNKNOWN
		return code, err
	}

	//tokenName := tokens.Name
	tokenName := tokens.Mark //

	uCurrency := new(UserCurrency)
	_, err = engine.Where("uid =? and token_id =?", this.SellId, this.TokenId).Get(uCurrency)
	if err != nil {
		log.Errorln("查询用户余额失败!", err.Error())
		code = ERRCODE_USER_BALANCE
		return code, err
	}

	buyCurrency := new(UserCurrency)
	_, err = engine.Table("user_currency").Where("uid =? and token_id =?", this.BuyId, this.TokenId).Get(buyCurrency)
	if err != nil {
		log.Errorln("查询用户余额失败!", err.Error())
		code = ERRCODE_USER_BALANCE
		return
	}

	sellNum := allNum + rateFee
	if uCurrency.Freeze < sellNum || uCurrency.Freeze < 0 {
		log.Println("余额不足!")
		err = errors.New("余额不足")
		code = ERRCODE_SELLER_LESS
		return code, err
	}

	// 事务开始
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()

	if err != nil {
		log.Println(err.Error())
		session.Rollback()
		code = ERRCODE_TRADE_ERROR
		return
	}

	/////////////////////////////////////////////////////////
	// 1. user_currency扣, 买家减交易金额，卖家加交易金额和减去费用
	////////////////////////////////////////////////////////
	freezeCny := convert.Int64MulInt64By8Bit(sellNum, this.Price)
	sellSql := "update user_currency set  `freeze`=`freeze` - ?,`freeze_cny`=`freeze_cny`-? ,`version`=`version`+1  WHERE  uid = ? and token_id = ? and version = ?"
	sqlRest, err := session.Exec(sellSql, sellNum,  freezeCny, this.SellId, this.TokenId, uCurrency.Version) // 卖家 扣除平台费用
	if err != nil {
		log.Println(err.Error())
		session.Rollback()
		return
	}
	if rst, _ := sqlRest.RowsAffected(); rst == 0 {
		fmt.Println("卖家扣除失败!", err.Error())
		log.Errorln("卖家扣除失败!", err.Error())
		session.Rollback()
		code = ERRCODE_TRADE_ERROR
		err = errors.New("卖家扣除失败!")
		return code, err
	}

	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
	}

	has, err := session.Table(`user_currency`).Where("uid=? AND token_id=? ", this.BuyId, this.TokenId).Exist()
	if err != nil {
		log.Errorln("find user_currency error!")
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return
	}
	if has {
		fmt.Println("has ....")

		buySql := "update user_currency set `balance`=`balance`+?, `balance_cny`=`balance_cny`+?, `version`=`version`+1   WHERE uid=? and token_id=? and version=?"
		buyRest, err := session.Exec(buySql, allNum, freezeCny, this.BuyId, this.TokenId, buyCurrency.Version)
		if err != nil {
			fmt.Println("买家添加余额失败, ", err.Error())
			log.Println(err.Error())
			session.Rollback()
			code = ERRCODE_TRADE_ERROR
			return code, err
		}
		if rst, _ := buyRest.RowsAffected(); rst == 0 {
			fmt.Println("买家添加余额失败失败!", err.Error())
			log.Errorln("买家添加余额失败失败!", err.Error())
			session.Rollback()
			code = ERRCODE_TRADE_ERROR
			err = errors.New("买家添加余额失败失败!")
			return code, err
		}
	} else {
		fmt.Println("no has....")
		insertSql := "insert into user_currency (uid, token_id, token_name, balance,balance_cny , version) values (?, ?, ?, ?, ?, 0)"
		buyRest, err := session.Exec(insertSql, this.BuyId, this.TokenId, tokenName, allNum, freezeCny)
		if err != nil {
			fmt.Println("买家添加余额失败, ", err.Error())
			log.Println(err.Error())
			session.Rollback()
			code = ERRCODE_TRADE_ERROR
			return code, err
		}
		if rst, _ := buyRest.RowsAffected(); rst == 0 {
			fmt.Println("买家添加余额失败失败!", err.Error())
			log.Errorln("买家添加余额失败失败!", err.Error())
			session.Rollback()
			code = ERRCODE_TRADE_ERROR
			err = errors.New("买家添加余额失败失败!")
			return code, err
		}
		if err != nil {
			fmt.Println("insert error!, ", err.Error())
			log.Println(err.Error())
			session.Rollback()
			code = ERRCODE_TRADE_ERROR
			return code, err
		}
	}

	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
	}

	//////////////////////////////////////////////////
	// 2. 插入记录 user_currency_history, 转入，转出，卖家扣费用
	/////////////////////////////////////////////////
	var buffer bytes.Buffer
	buffer.WriteString("insert into user_currency_history ")
	buffer.WriteString("(uid, trade_uid,order_id, token_id, num, fee, surplus,  operator, address, states,created_time, updated_time)")
	buffer.WriteString("values (?, ?, ?, ?, ?, ?,  ?, ?, ?, ? ,?, ?), (?,?, ?, ?, ?, ?,   ?, ?, ?, ? , ?, ?)")
	insertSql := buffer.String()

	nowCreate := time.Now().Format("2006-01-02 15:04:05")
	_, err = session.Table(`user_currency_history`).Exec(insertSql,
		this.SellId, this.BuyId, this.OrderId, this.TokenId, sellNum, rateFee, uCurrency.Balance, 2, "", this.States, nowCreate, updateTimeStr, // 卖家记录 , 2订单转出
		this.BuyId, this.SellId, this.OrderId, this.TokenId, this.Num, 0, buyCurrency.Balance, 1, "", this.States, nowCreate, updateTimeStr, // 买家记录 , 1订单转入
	)
	if err != nil {
		fmt.Println("insert into history error: ", err.Error())
		log.Println(err.Error())
		session.Rollback()
		code = ERRCODE_TRADE_ERROR
		return code, err
	}
	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
	}

	////////////////////////////////////////////////////////
	// 3. 统计表 加 1
	////////////////////////////////////////////////////////

	has, err = session.Exist(&UserCurrencyCount{Uid: this.SellId})
	if err != nil {
		log.Println(err.Error())
	}
	if has {
		countAddOneSql := "UPDATE  `user_currency_count` set `success` = `success` + 1, `orders`=`orders`+1  where uid=? "
		_, err = session.Exec(countAddOneSql, this.SellId)
		if err != nil {
			log.Println(err.Error())
			session.Rollback()
			code = ERRCODE_TRADE_ERROR
			return code, err
		}
	} else {
		insertSql := "INSERT INTO `user_currency_count` (uid,orders, success, good) values(?,?,?, ?)"
		_, err = session.Exec(insertSql, this.SellId, 1, 1, 100)
		if err != nil {
			log.Println(err.Error())
			session.Rollback()
			code = ERRCODE_TRADE_ERROR
			return code, err
		}
	}

	buyHas, err := session.Exist(&UserCurrencyCount{Uid: this.BuyId})
	if err != nil {
		log.Println(err.Error())
	}
	if buyHas {
		countAddOneSql := "UPDATE  `user_currency_count` set `success` = `success` + 1, `orders`=`orders`+1  where uid=? "
		_, err = session.Exec(countAddOneSql, this.BuyId)
		if err != nil {
			log.Println(err.Error())
			session.Rollback()
			code = ERRCODE_TRADE_ERROR
			return code, err
		}
	} else {
		insertSql := "INSERT INTO `user_currency_count` (uid,orders, success, good) values(?,?,?, ?)"
		_, err = session.Exec(insertSql, this.BuyId, 1, 1, 100)
		if err != nil {
			log.Println(err.Error())
			session.Rollback()
			code = ERRCODE_TRADE_ERROR
			return code, err
		}
	}

	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	adsM := new(Ads).Get(this.AdId)


	//////////////////////////////////////////////////////
	// 5. 更新状态
	//////////////////////////////////////////////////////

	updateStatesSql := "UPDATE   `order`   SET   `states`=?, `updated_time`=?,`release_time`=?,`fee`=?  WHERE  `id`=?"
	_, err = session.Exec(updateStatesSql, 3, updateTimeStr, now, rateFee, Id)
	if err != nil {
		fmt.Println("update order states error:", err.Error())
		log.Println(err.Error())
		session.Rollback()
		code = ERRCODE_TRADE_ERROR
		return code, err
	}
	err = engine.ClearCache(new(UserCurrency))
	if err != nil {
		log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
	}
	err = session.Commit()
	if err != nil {
		log.Errorln(err.Error())
		session.Rollback()
		code = ERRCODE_UNKNOWN
		return code, err
	}

	code = ERRCODE_SUCCESS

	/////    检查广告是否需要下架
	go AdsAutoDownline(adsM.Id)
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
		log.Errorln(err.Error())
		code = ERRCODE_ORDER_NOTEXIST
	}
	return
}



func (this *Order) GetOrdersByStatus()(ods []Order, err error){
	engine := dao.DB.GetMysqlConn()
	err = engine.Where("states = 1 OR states = 2").Find(&ods)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}





/*

 */
func (this *Order) GetOrderByTime(uid uint64, startTime, endTime string) (ods []Order, err error) {
	engine := dao.DB.GetMysqlConn()
	err = engine.Where("(sell_id = ? OR buy_id=? ) AND created_time >= ? AND created_time <= ? AND states=3", uid, uid, startTime, endTime).Find(&ods)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}



/*
	根据tokenid来获取订单
*/
func (this *Order) GetOrderByTokenIdByTime(tokenid uint32, startTime, endTime string) (ods []Order, err error) {
	engine := dao.DB.GetMysqlConn()
	err = engine.Where("token_id = ?  AND created_time >= ? AND created_time <= ? AND states = 3", tokenid, startTime, endTime).Find(&ods)
	return
}







/*
    获取历史记录
*/
func (this *Order) GetOrderHistory(startTime, endTime string, limit int32) (uhistory []Order, err error) {
	now := time.Now()
	if startTime != "" {
		startTime = now.Format("2006-01-02")
	}
	if endTime != "" {
		endTime = now.Format("2006-01-02")
	}
	token := new(CommonTokens).Get(0, "BTC")
	var tokenId uint32
	if token != nil {
		tokenId = token.Id
	}else{
		tokenId = 2
	}
	engine := dao.DB.GetMysqlConn()
	if limit != 0  && startTime == "" && endTime == ""{
		err = engine.Where("token_id=?", tokenId).Limit(int(limit)).Desc("created_time").Find(&uhistory)
	} else if limit == 0 && startTime != "" && endTime != ""{
		err = engine.Where("token_id = ? AND created_time >= ? && created_time <= ?",tokenId, startTime, endTime).Desc("created_time").Find(&uhistory)
	}else{
		err = engine.Where("token_id =? AND created_time >= ? && created_time <= ?", tokenId, startTime, endTime).Desc("created_time").Limit(int(limit)).Find(&uhistory)
	}
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}



////////////////////////////////    func ////////////////////////////////////////

/*

 */
func CheckOrderExiryTime(id uint64, exiryTime string) {
	od := new(Order)
	for {
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println("id: ", id, " order Exitry time(min): ", GetHourDiffer(now, exiryTime))
		if GetHourDiffer(now, exiryTime) <= 0 {
			engine := dao.DB.GetMysqlConn()
			_, err := engine.Where("id=?", id).Get(od)
			if err != nil {
				log.Errorln("get order states error!")
			} else {
				if od.States == 0 || od.States == 2 || od.States == 3 || od.States == 4 { // 0删除 2待放行(已支付) 3确认支付(已完成) 4取消
					fmt.Println("break: od stats", od.States)
					break
				}
			}
			err = CancelAction(od.Id, 3)
			if err != nil {
				err = CancelAction(od.Id, 3)
			}
			break
		}
		time.Sleep(2 * 60 * time.Second)
	}
	log.Println(id, " break .................")
	return
}


func CancelAction(id uint64, CancelType uint32) (err error){

	engine := dao.DB.GetMysqlConn()
	od := new(Order)
	_, err = od.GetOrder(id)

	uCurrency := new(UserCurrency)
	_, err = engine.Where("uid =? and token_id =?", od.SellId, od.TokenId).Get(uCurrency)
	if err != nil {
		log.Errorln("查询用户余额失败!", err.Error())
		return
	}

	rate := conf.Cfg.MustValue("rate", "fee_rate")
	rateFloat, _ := strconv.ParseInt(rate, 10, 64)

	allNum := od.Num
	rateFee := allNum * rateFloat

	// 事务开始
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()

	if err != nil {
		fmt.Println("session start error!!!!!!!!!!!")
		log.Println(err.Error())
		session.Rollback()
		return
	}

	/////////////////////////////////////////////////////////
	// 1. user_currency 还原 freeze
	////////////////////////////////////////////////////////

	sellNum := rateFee + allNum
	freezeCny := convert.Int64MulInt64By8Bit(sellNum, od.Price)

	if uCurrency.Freeze - sellNum < 0 {
		sellNum = 0
	}
	if uCurrency.FreezeCny - freezeCny <0 {
		freezeCny = 0
	}

	sellSql := "update user_currency set   `balance`=`balance`+?, `freeze`=`freeze` - ?,`balance_cny`=`balance_cny`+?, `freeze_cny`=`freeze_cny`-? ,  `version`=`version`+1  WHERE  uid = ? and token_id = ? and version = ?"
	sqlRest, err := session.Exec(sellSql, sellNum, sellNum,freezeCny, freezeCny, od.SellId, od.TokenId, uCurrency.Version) // 卖家 扣除平台费用

	if err != nil {
		fmt.Println("user_currency 还原 freeze error")
		log.Println(err.Error())
		log.Errorln("od.Id:", od.Id, od.ExpiryTime)
		session.Rollback()
		return
	}
	if rst, _ := sqlRest.RowsAffected(); rst == 0 {
		fmt.Println("卖家扣除失败!", err.Error())
		log.Errorln("卖家扣除失败!", err.Error())
		session.Rollback()
		err = errors.New("卖家扣除失败!")
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err = session.Where("id = ?", id).Update(&Order{States: 4, UpdatedTime: now})
	if err != nil {
		log.Errorln("order id exiry time out update status = 4 error! id:", id, err.Error())
		session.Rollback()
		return
	}

	sql := "UPDATE   `order`   SET   `states`=? , `cancel_type`=?, `updated_time`=?  WHERE  `id`=?"
	_, err = session.Exec(sql, 4, CancelType, now, od.Id)
	if err != nil {
		log.Println("CancelType:", CancelType)
		log.Errorln(err.Error())
		session.Rollback()
		return
	}
	fmt.Println("session commit ....")
	err = session.Commit()
	if err != nil {
		log.Println(err.Error())
		session.Rollback()
		return
	}

	return
}





//获取相差时间
func GetHourDiffer(start_time, end_time string) int64 {
	var minute int64
	t1, err := time.Parse("2006-01-02 15:04:05", start_time)
	t2, err := time.Parse("2006-01-02 15:04:05", end_time)
	//fmt.Println(t1,t2)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		minute = diff / 60            //3600
		return minute
	} else {
		return minute
	}
}
