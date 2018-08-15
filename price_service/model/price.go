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
	UsdVol      int64  `xorm:"BIGINT(20)"`
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

type Volume struct {
	Sum    int64 `xorm:"BIGINT(20)"`
	Amount int64 `xorm:"BIGINT(20)"`
}

func InsertPrice(p *Price) error {
	_, err := DB.GetMysqlConn().InsertOne(p)
	if err != nil {
		log.WithFields(logrus.Fields{
			"id": p.Id,
		}).Errorln(err.Error())
		return err
	}
	return nil
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
func Calculate(token_id int32 ,price, amount, cny_price int64, symbol string, high, low int64) *proto.PriceBaseData {
	t := time.Now()

	s := t.Second()
	min := t.Add(-time.Duration(s) * time.Second)
	same := min.Unix()
	log.Info(same)
	l := min.Add(-86400 * time.Second)
	yestday := l.Unix()
	p := &Price{}
	ok, err := DB.GetMysqlConn().Where("id=? and symbol=?", yestday, symbol).Get(p)
	//_, err := DB.GetMysqlConn().Where("symbol=? and created_time>=? and created_time<? ", symbol, begin, end).Asc("created_time").Limit(1, 0).Get(p)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	if !ok {
		g := ConfigQueneInit[symbol]
		p.Price = g.Price
		p.Amount = 0
	}
	log.WithFields(log.Fields{
		"high":     high,
		"low":      low,
		"price ":   price,
		"p.Price":  p.Price,
		"amount":   amount,
		"p.Amount": p.Amount,
		"price":    price,
		"id":       p.Id,
		"symbol":   symbol,
	}).Info("current price print")
	var cny string
	c,ok:=GetQueneMgr().PriceMap[token_id]
	if !ok {
		cny="0"
	}else{
		cny=convert.Int64ToStringBy8Bit(c.CnyPrice)
	}

	return &proto.PriceBaseData{
		Symbol:   symbol,
		High:     convert.Int64ToStringBy8Bit(high),
		Low:      convert.Int64ToStringBy8Bit(low),
		Scope:    convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
		Amount:   convert.Int64ToStringBy8Bit(amount - p.Amount),
		Price:    convert.Int64ToStringBy8Bit(price),
		//CnyPrice: convert.Int64MulInt64By8BitString(cny_price, price),
		CnyPrice :cny,
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

//查询交易量
func GetVolumeTotal() *proto.VolumeResponse {
	t := time.Now().Local()
	//nowUnix := time.Now().Unix()
	dayUnix := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	weekUnix := time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()), 0, 0, 0, 0, t.Location()).Unix()
	mondayUnix := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Unix()

	var nowVolume Volume
	var dayVolume Volume
	var weekVolume Volume
	var monthVolume Volume

	var err error
	var res bool
	res, err = DB.GetMysqlConn().SQL("select sum(usd_vol) Sum,sum(amount) Amount from (select * from price where id >= (select max(id) from price) group by symbol order by id desc) as a").Get(&nowVolume)
	if err != nil || res != true {
		log.Warningln(err.Error())
		return nil
	}
	res, err = DB.GetMysqlConn().SQL("select sum(usd_vol) Sum,sum(amount) Amount from (select * from price where id > ? group by symbol order by id desc) as a", dayUnix).Get(&dayVolume)
	if err != nil || res != true {
		log.Warningln(err.Error())
		return nil
	}
	res, err = DB.GetMysqlConn().SQL("select sum(usd_vol) Sum,sum(amount) Amount from (select * from price where id > ? group by symbol order by id desc) as a", weekUnix).Get(&weekVolume)
	if err != nil || res != true {
		log.Warningln(err.Error())
		return nil
	}
	res, err = DB.GetMysqlConn().SQL("select sum(usd_vol) Sum,sum(amount) Amount from (select * from price where id > ? group by symbol order by id desc) as a", mondayUnix).Get(&monthVolume)
	if err != nil || res != true {
		log.Warningln(err.Error())
		return nil
	}

	data := &proto.VolumeResponse{
		DayVolume:   nowVolume.Sum - dayVolume.Sum,
		WeekVolume:  nowVolume.Sum - weekVolume.Sum,
		MonthVolume: nowVolume.Sum - monthVolume.Sum}
	return data

}
