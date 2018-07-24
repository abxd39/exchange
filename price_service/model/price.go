package model

import (
	"digicon/common/convert"
	. "digicon/price_service/dao"
	proto "digicon/proto/rpc"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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

/*
type Price struct {
	Id     int64 `xorm:"BIGINT(20)"`
	Open   int64 `xorm:"comment('开盘价') BIGINT(20)"`
	Close  int64 `xorm:"comment('收盘价') BIGINT(20)"`
	Low    int64 `xorm:"comment('最低价') BIGINT(20)"`
	High   int64 `xorm:"comment('最高价') BIGINT(20)"`
	Amount int64 `xorm:"comment('成交量') BIGINT(20)"`
	Vol    int64 `xorm:"comment('成交额') BIGINT(20)"`
	Count  int64 `xorm:"BIGINT(20)"`
}
*/
func InsertPrice(p *Price) {
	_, err := DB.GetMysqlConn().InsertOne(p)
	if err != nil {
		log.WithFields(logrus.Fields{
			"id": p.Id,
		}).Errorln(err.Error())
		return
	}
	return
}

func GetHigh(begin, end int64, symbol string) (high int64) {
	hp := &Price{}
	ok, err := DB.GetMysqlConn().Where("created_time>? and created_time<=? and symbol=?", begin, end, symbol).Desc("price").Limit(1, 0).Get(hp)
	if err != nil {
		log.Errorln(err.Error())
		return 0
	}
	if ok {
		return hp.Price
	}
	return 0
}

func GetLow(begin, end int64, symbol string) (low int64) {
	hp := &Price{}
	ok, err := DB.GetMysqlConn().Where("created_time>? and created_time<=? and symbol=?", begin, end, symbol).Asc("price").Limit(1, 0).Get(hp)
	if err != nil {
		log.Errorln(err.Error())
		return 0
	}
	if ok {
		return hp.Price
	}
	return 0
}

//计算当前价格数据
func Calculate(price, amount, cny_price int64, symbol string) *proto.PriceBaseData {
	t := time.Now()
	l := t.Add(-86400 * time.Second)

	begin := l.Unix()
	end := t.Unix()
	h := GetHigh(begin, end, symbol)

	j := GetLow(begin, end, symbol)

	p := &Price{}
	_, err := DB.GetMysqlConn().Where("symbol=? and created_time>=? and created_time<? ", symbol, begin, end).Asc("created_time").Limit(1, 0).Get(p)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}

	return &proto.PriceBaseData{
		Symbol:   symbol,
		High:     convert.Int64ToStringBy8Bit(h),
		Low:      convert.Int64ToStringBy8Bit(j),
		Scope:    convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
		Amount:   convert.Int64ToStringBy8Bit(amount - p.Amount),
		Price:    convert.Int64ToStringBy8Bit(price),
		CnyPrice: convert.Int64MulInt64By8BitString(cny_price, price),
	}

}

//从数据库获取最新价格
func GetPrice(symbol string) (*Price, bool) {
	m := &Price{}
	ok, err := DB.GetMysqlConn().Where("symbol=?", symbol).Desc("created_time").Limit(1, 0).Get(m)
	if err != nil {
		log.Errorf(err.Error())
	}
	return m, ok
}

func Get24HourPrice(symbol string) (*Price, bool) {
	t := time.Now()
	l := t.Add(-86400 * time.Second)

	begin := l.Unix()
	end := t.Unix()
	p := &Price{}

	ok, err := DB.GetMysqlConn().Where("symbol=? and created_time>=? and created_time<? ", symbol, begin, end).Asc("created_time").Limit(1, 0).Get(p)
	if err != nil {
		log.Errorln(err.Error())
		return nil, ok
	}

	return p, ok
}

func (s *Price) SetProtoData() *proto.PriceCache {
	return &proto.PriceCache{
		Id:          s.Id,
		Symbol:      s.Symbol,
		Price:       s.Price,
		CreatedTime: s.CreatedTime,
		Amount:      s.Amount,
		Vol:         s.Vol,
		Count:       s.Count,
	}
}
