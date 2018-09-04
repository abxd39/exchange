package models

import (
	"digicon/wallet_service/utils"
	"digicon/bak/templete2/log"
	"fmt"
)

type TokenInoutDailySheet struct {
	Id             int    `xorm:"not null pk autoincr comment('自增id') BIGINT(20)"   json:"id"`
	TokenId        int    `xorm:"not null comment('货币id') TINYINT(4)"               json:"token_id"`
	TokenName      string `xorm:"not null comment('货币名称') VARCHAR(20)"             json:"token_name"`
	TotalDayNum    int64  `xorm:"not null comment('日提币总量') BIGINT(20)"             json:"total_day_num"`
	TotalDayCny    int64  `xorm:"not null comment('日提币总数折合') BIGINT(20)"         json:"total_day_cny"`
	TotalDayNumFee int64  `xorm:"not null comment('日提币手续费数量') BIGINT(20)"       json:"total_day_num_fee"`
	TotalDayFeeCny int64  `xorm:"not null comment('日提币手续费总数折合') BIGINT(20)"    json:"total_day_fee_cny"`
	TotalDayPut    int64  `xorm:"not null comment('日充币总额') BIGINT(20)"             json:"total_day_put"`
	TotalDayPutCny int64  `xorm:"not null default 0 comment('日充币折合') BIGINT(20)"  json:"total_day_put_cny"`
	Total          int64  `xorm:"not null comment('提币累计总金额') BIGINT(20)"         json:"total"`
	TotalFee       int64  `xorm:"not null comment('提币手续费累计总金额') BIGINT(20)"    json:"total_fee"`
	TotalPut       int64  `xorm:"not null comment('充币累计总额') BIGINT(20)"           json:"total_put"`
	Date           string  `xorm:"not null comment('时间戳') VARCHAR(20)"                json:"date"`
}

type FindDailySheet struct {
	Id      int     `xorm:"not null pk comment('自增id') BIGINT(20)"                     json:"id"`
	//DateStr string  `xorm:"not null comment('时间戳，精确到天') VARCHAT(20)"                json:"date_str"`
}


func (this *TokenInoutDailySheet) InsertOneDayTotal() (err error){
	sql := "insert into `token_inout_daily_sheet` " +
		"(`token_id`, `token_name`, `total_day_num`, `total_day_cny`, `total_day_num_fee`, `total_day_fee_cny`," +
		"`total_day_put`, `total_day_put_cny`, `total`, `total_fee`, `total_put`, `date`) "+
		" values(?, ?, ?, ?, ?, ?,   ?, ?, ?, ?, ?, ?)"
	engine := utils.Engine_wallet
	engine.ShowSQL(true)
	_, err = engine.Exec(sql, this.TokenId, this.TokenName, this.TotalDayNum, this.TotalDayCny, this.TotalDayNumFee, this.TotalDayFeeCny,
		this.TotalDayPut, this.TotalDayPutCny, this.Total, this.TotalFee,this.TotalPut, this.Date)
	engine.ShowSQL(false)
	if err!=nil{
		log.Log.Infof("InsertOneDayTotal",err.Error())
	}
	return
}

func (this *TokenInoutDailySheet) UpdateOneDayTotal(id int64) (err error){
	sql := "UPDATE  `token_inout_daily_sheet` " +
		"set `token_name`=?, `total_day_num`=?, `total_day_cny`=?, `total_day_num_fee`=?, `total_day_fee_cny`=? , " +
		"`total_day_put`=?, `total_day_put_cny`=?, `total`=?, `total_fee`=?, `total_put`=?  WHERE id=?"
	engine := utils.Engine_wallet
	engine.ShowSQL(true)
	_, err = engine.Exec(sql, this.TokenName, this.TotalDayNum, this.TotalDayCny, this.TotalDayNumFee, this.TotalDayFeeCny,
		this.TotalDayPut, this.TotalDayPutCny, this.Total, this.TotalFee,this.TotalPut,  id)
	engine.ShowSQL(false)
	if err!=nil{
		log.Log.Infof("UpdateOneDayTotal",err.Error())
	}
	return
}





func (this *TokenInoutDailySheet) CheckOneDay(tkId uint32, today string)(result FindDailySheet, err error) {
	engine := utils.Engine_wallet
	sql := "SELECT id FROM g_wallet.`token_inout_daily_sheet` WHERE token_id=? AND `date` = ?;"
	_, err = engine.SQL(sql, tkId, today ).Get(&result)
	if err!=nil{
		log.Log.Infof("CheckOneDay",err.Error())
		fmt.Println("CheckOneDay",err.Error())
	}
	return
}
