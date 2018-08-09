package model

import "digicon/currency_service/dao"

type CurrencyDailySheet struct {
	Id              int32 `xorm:"not null pk comment('自增id') TINYINT(4)"                     json:"id"`
	TokenId         int32 `xorm:"not null comment('币种ID') INT(11)"                           json:"token_id"`
	SellTotal       int64 `xorm:"not null default 0 comment('法币卖出总数') BIGINT(20)"         json:"sell_total"`
	SellCny         int64 `xorm:"not null comment('法币卖出总额折合cny') BIGINT(20)"             json:"sell_cny"`
	BuyTotal        int64 `xorm:"not null default 0 comment('法币买入总数') BIGINT(20)"          json:"buy_total"`
	BuyCny          int64 `xorm:"not null default 0 comment('法币买入总额折合cny') BIGINT(20)"    json:"buy_cny"`
	FeeSellTotal    int64 `xorm:"not null default 0 comment('法币卖出手续费总数') BIGINT(20)"     json:"fee_sell_total"`
	FeeSellCny      int64 `xorm:"not null comment('法币卖出手续费折合cny') BIGINT(20)"            json:"fee_sell_cny"`
	FeeBuyTotal     int64 `xorm:"not null default 0 comment('法币买入手续费总数') BIGINT(20)"     json:"fee_buy_total"`
	FeeBuyCny       int64 `xorm:"not null default 0 comment('法币买入手续费折合cny') BIGINT(20)"  json:"fee_buy_cny"`
	BuyTotalAll     int64 `xorm:"not null comment('累计买入总额') BIGINT(20)"                    json:"buy_total_all"`
	BuyTotalAllCny  int64 `xorm:"not null comment('累计买入总额折合cny') BIGINT(20)"              json:"buy_total_all_cny"`
	SellTotalAll    int64 `xorm:"not null comment('累计卖出总额') BIGINT(20)"                    json:"sell_total_all"`
	SellTotalAllCny int64 `xorm:"not null comment('累计卖出总额折合') BIGINT(20)"                 json:"sell_total_all_cny"`
	Total           int64 `xorm:"not null comment('总数') BIGINT(20)"                           json:"total"`
	TotalCny        int64 `xorm:"not null comment('总数折合') BIGINT(20)"                       json:"total_cny"`
	Date            int64 `xorm:"not null comment('时间戳，精确到天') BIGINT(10)"                json:"date"`
}


type FindDailySheet struct {
	Id                 int32 `xorm:"not null pk comment('自增id') TINYINT(4)"                     json:"id"`
	DateStr            int64 `xorm:"not null comment('时间戳，精确到天') VARCHAT(20)"                json:"date_str"`
}


func (this *CurrencyDailySheet) Insert () (err error){
	engine := dao.DB.GetMysqlConn()
	_, err = engine.InsertOne(this)
	return
}

func (this *CurrencyDailySheet) InsertOneDay() (err error){
	sql := "insert into  `currency_daily_sheet` " +
	"(`token_id`, `sell_total`, `sell_cny`, `buy_total`, `buy_cny`, `fee_sell_total`, `fee_sell_cny`, `fee_buy_total`, " +
		"`fee_buy_cny`,  `buy_total_all`, `buy_total_all_cny`, `sell_total_all`, `sell_total_all_cny`, `total`, `total_cny`, `date` " +
		") value (?, ?, ?, ?,    ?, ?, ?, ?,    ?, ?, ?, ?,   ?, ?, ?, ?)"
	engine := dao.DB.GetMysqlConn()
	_, err = engine.Exec(sql, this.TokenId, this.SellTotal, this.SellCny, this.BuyTotal, this.BuyCny, this.FeeSellTotal, this.FeeSellCny, this.FeeBuyTotal,
		this.FeeBuyCny, this.BuyTotalAll, this.BuyTotalAllCny, this.SellTotalAll, this.SellTotalAllCny, this.Total, this.TotalCny, this.Date)
	return
}



func (this *CurrencyDailySheet) GetOneDay(tokenid uint32, today int64) (result FindDailySheet,err error) {
	sql := "SELECT  id, FROM_UNIXTIME(?, \"%Y-%m-%d\") as date_str  FROM  `currency_daily_sheet` WHERE token_id=? "
	engine := dao.DB.GetMysqlConn()
	_, err = engine.SQL(sql, today, tokenid).Get(&result)
	return
}

