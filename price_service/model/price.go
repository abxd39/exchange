package model

import (
	"digicon/common/convert"
	. "digicon/price_service/dao"
	. "digicon/price_service/log"
	proto "digicon/proto/rpc"
	"time"
)

type Price struct {
	Id          int64  `xorm:"index(keep) BIGINT(20)"`
	Symbol      string `xorm:"index(keep) VARCHAR(32)"`
	Price       int64  `xorm:"BIGINT(20)"`
	CreatedTime int64  `xorm:"BIGINT(20)"`
	Amount      int64  `xorm:"BIGINT(20)"`
	Vol         int64  `xorm:"BIGINT(20)"`
	Count       int64  `xorm:"BIGINT(20)"`
}

func InsertPrice(p *Price) {
	_, err := DB.GetMysqlConn().InsertOne(p)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func GetHigh(begin, end int64) (high int64) {
	hp := &Price{}
	ok, err := DB.GetMysqlConn().Where("created_time>? and end_time<?", begin, end).Decr("price").Limit(1, 0).Get(hp)
	if err != nil {
		Log.Errorln(err.Error())
		return 0
	}
	if ok {
		return hp.Price
	}
	return 0
}

func GetLow(begin, end int64) (low int64) {
	hp := &Price{}
	ok, err := DB.GetMysqlConn().Where("created_time>? and end_time<?", begin, end).Asc("price").Limit(1, 0).Get(hp)
	if err != nil {
		Log.Errorln(err.Error())
		return 0
	}
	if ok {
		return hp.Price
	}
	return 0
}

func Calculate(price, amount ,cny_price int64, symbol string) *proto.PriceBaseData {
	t := time.Now()
	l := t.Add(-86400 * time.Second)
	h := GetHigh(l.Unix(), t.Unix())

	j := GetLow(l.Unix(), t.Unix())

	p := &Price{}
	_, err := DB.GetMysqlConn().Where("symbol=? amd created_time>=? created_time<? ", symbol, t.Unix()).Asc("created_time").Limit(1, 0).Get(p)
	if err != nil {
		Log.Errorln(err.Error())
		return nil
	}

	return &proto.PriceBaseData{
		Symbol: symbol,
		High:   convert.Int64ToStringBy8Bit(h),
		Low:    convert.Int64ToStringBy8Bit(j),
		Scope:  convert.Int64DivInt64By8BitString(price-p.Price, p.Price),
		Amount: convert.Int64ToStringBy8Bit(amount - p.Amount),
		Price:  convert.Int64ToStringBy8Bit(price),
		CnyPrice:convert.Int64MulInt64By8BitString(cny_price,price),
	}

}
