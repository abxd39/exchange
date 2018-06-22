package model

import (
	"digicon/common/convert"
	. "digicon/proto/common"
	. "digicon/token_service/dao"
	. "digicon/token_service/log"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"sync/atomic"

)

//交易队列类型
type EntrustQuene struct {
	//币币队列ID  格式主要货币_交易货币
	TokenQueneId string

	TokenId int
	TokenTradeId int
	//卖出队列key
	SellQueneId string
	//买入队列key
	BuyQueneId string

	//当前队列自增ID
	UUID int64

	//key 是委托ID
	sourceData map[string]*EntrustDetail

	//缓存将要保存的DB的委托请求
	newOrderDetail chan *EntrustDetail

	//缓存将要更新的DB的委托请求
	updateOrderDetail chan *EntrustDetail

	//等待处理的委托请求
	waitOrderDetail chan *EntrustDetail
	//上一次成交价格
	price int64
}

const (
	ORDER_OPT_BUY  = 0 //买类型
	ORDER_OPT_SELL = 1 //卖类型

)

const (
	ENTRUST_MARKET_PRICE = 1 //市价委托
	ENTRUST_LIMIT_PRICE  = 2 //限价委托

)

func NewEntrustQuene(quene_id string) *EntrustQuene {
	m := &EntrustQuene{
		TokenQueneId:      quene_id,
		BuyQueneId:        fmt.Sprintf("%s:0", quene_id),
		SellQueneId:       fmt.Sprintf("%s:1", quene_id),
		UUID:              1,
		sourceData:        make(map[string]*EntrustDetail),
		newOrderDetail:    make(chan *EntrustDetail),
		waitOrderDetail:   make(chan *EntrustDetail),
		updateOrderDetail: make(chan *EntrustDetail),
	}
	go m.process()
	return m
}

//获取自增ID
func (s *EntrustQuene) GetUUID() int64 {
	return atomic.AddInt64(&s.UUID, 1)
}

func (s *EntrustQuene) JoinWaitChann(p *EntrustDetail) {
	s.newOrderDetail <- p
	s.waitOrderDetail <- p

}

func (s *EntrustQuene) MakeDeal(buyer *EntrustDetail,seller *EntrustDetail,price int64 )  {
	/*
	if buyer.SurplusNum < seller.SurplusNum {

		deal_num:=buyer.SurplusNum

		t:=&Trade{
			Uid:buyer.Uid,
			TokenId:s.TokenId,
			TokenTradeId:s.TokenTradeId,
			Price:price,
			Num:deal_num,
			Fee:price/1000,
			DealTime:time.Now().Unix(),
			States:TRADE_STATES_ALL,
		}


		o:=&Trade{
			Uid:seller.Uid,
			TokenId:s.TokenId,
			TokenTradeId:s.TokenTradeId,
			Price:price,
			Num:deal_num,
			Fee:price/1000,
			DealTime:time.Now().Unix(),
			States:TRADE_STATES_PART,
		}


		s.waitOrderDetail<-
	}else if buyer.SurplusNum == seller.SurplusNum{
		deal_num:=buyer.SurplusNum

		t:=&Trade{
			Uid:buyer.Uid,
			TokenId:s.TokenId,
			TokenTradeId:s.TokenTradeId,
			Price:price,
			Num:deal_num,
			Fee:price/1000,
			DealTime:time.Now().Unix(),
			States:TRADE_STATES_ALL,
		}


		o:=&Trade{
			Uid:seller.Uid,
			TokenId:s.TokenId,
			TokenTradeId:s.TokenTradeId,
			Price:price,
			Num:deal_num,
			Fee:price/1000,
			DealTime:time.Now().Unix(),
			States:TRADE_STATES_ALL,
		}
	}else{
		deal_num:=seller.SurplusNum

		t:=&Trade{
			Uid:buyer.Uid,
			TokenId:s.TokenId,
			TokenTradeId:s.TokenTradeId,
			Price:price,
			Num:deal_num,
			Fee:price/1000,
			DealTime:time.Now().Unix(),
			States:TRADE_STATES_PART,
		}


		o:=&Trade{
			Uid:seller.Uid,
			TokenId:s.TokenId,
			TokenTradeId:s.TokenTradeId,
			Price:price,
			Num:deal_num,
			Fee:price/1000,
			DealTime:time.Now().Unix(),
			States:TRADE_STATES_ALL,
		}
	}

*/
}

func (s *EntrustQuene) Match(p *EntrustDetail) (ret int, err error) {
	if p.Opt > 2 {
		ret = ERRCODE_PARAM
		return
	}
	var other *EntrustDetail
	if p.Opt == ORDER_OPT_BUY {
		other, err = s.GetFirstEntrust(ORDER_OPT_SELL)
		if err==redis.Nil {
			s.JoinSellQuene(p)
		} else if err != nil {
			Log.Errorln(err.Error())
			return
		}

		if p.Type == ENTRUST_MARKET_PRICE {
			s.MakeDeal(p,other,other.OnPrice)
		}

	}
	ret, err = s.JoinSellQuene(p)
	return
}

//限价委托入队列 opt 0 buy ,1 sell
func (s *EntrustQuene) JoinSellQuene(p *EntrustDetail) (ret int, err error) {
	if p.Opt > 2 {
		ret = ERRCODE_PARAM
		return
	}
	var quene_id string
	if p.Opt == ORDER_OPT_BUY {
		quene_id = s.BuyQueneId
	} else if p.Opt == ORDER_OPT_SELL {
		quene_id = s.SellQueneId
	}

	//may be not exact
	x := convert.Int64ToFloat64By8Bit(p.Price)

	err = DB.GetRedisConn().ZAdd(quene_id, redis.Z{
		Member: p.Uid,
		Score:  x,
	}).Err()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}


	if ok := s.insertOrderDetail(p); ok {
		//s.pushOrderDetail <- p
	}


	return
}

//获取订单详情
func (s *EntrustQuene) GetOrderDetail(order_id string) (d *EntrustDetail, ok bool) {
	d, ok = s.sourceData[order_id]
	if !ok {
		return
	}
	return
}

//保存订单详情
func (s *EntrustQuene) insertOrderDetail(d *EntrustDetail) bool {
	if _, ok := s.GetOrderDetail(d.EntrustId); ok {
		return false
	}
	s.sourceData[d.EntrustId] = d
	return true
}

func (s *EntrustQuene) delOrderDetail(d *EntrustDetail) bool {
	return true
}
//延时保存委托数据
func (s *EntrustQuene) process() {

	for {
		select {
		//case d = <-s.pushOrderDetail:
			//d.Save()
		case w := <-s.waitOrderDetail:
			s.Match(w)
		case t:= <-s.newOrderDetail:
			t.Save()

		}
	}

}

//获取队列首位交易单
func (s *EntrustQuene) GetFirstEntrust(opt int) (en *EntrustDetail, err error) {
	var z []redis.Z
	var ok bool
	if opt == ORDER_OPT_BUY { //买入类型
		z, err = DB.GetRedisConn().ZRangeWithScores(s.BuyQueneId, 0, 1).Result()
	} else if opt == ORDER_OPT_SELL { //卖出类型
		z, err = DB.GetRedisConn().ZRevRangeWithScores(s.BuyQueneId, 0, 1).Result()
	}

	if err != nil {
		Log.Errorln(err)
		return
	}

	if len(z) > 0 {
		d := z[0].Member.(string)
		en, ok = s.GetOrderDetail(d)
		if ok {
			err = DB.GetRedisConn().ZRem(s.BuyQueneId,d).Err()

			return
		}
		err = errors.New("this is unrealize err when get order detail  ")
		Log.WithFields(logrus.Fields{
			"quene_id": s.TokenQueneId,
			"opt":      opt,
			"member":   d,
		}).Errorln(err.Error())

		return
	}

	err = errors.New("this is sync data err ")
	Log.WithFields(logrus.Fields{
		"quene_id": s.TokenQueneId,
		"opt":      opt,
	}).Errorln(err.Error())
	return
}
