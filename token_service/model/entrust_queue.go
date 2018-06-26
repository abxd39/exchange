package model

import (
	"digicon/common/convert"
	"digicon/common/genkey"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	. "digicon/token_service/log"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

//交易队列类型
type EntrustQuene struct {
	//币币队列ID  格式主要货币_交易货币
	TokenQueneId string

	TokenId      int
	TokenTradeId int
	//卖出队列key
	SellQueneId string
	//买入队列key
	BuyQueneId string

	//当前队列自增ID
	UUID int64

	//key 是委托ID，委托数据源
	sourceData map[string]*EntrustData

	//缓存将要保存的DB的委托请求
	//newOrderDetail chan *EntrustData

	//缓存将要更新的DB的委托请求
	updateOrderDetail chan *EntrustData

	//等待处理的委托请求
	waitOrderDetail chan *EntrustData

	//市价等待队列
	marketOrderDetail chan *EntrustData
	//上一次成交价格
	price int64
}

/*
const (
	ORDER_OPT_BUY  = 0 //买类型
	ORDER_OPT_SELL = 1 //卖类型

)

const (
	ENTRUST_MARKET_PRICE = 1 //市价委托
	ENTRUST_LIMIT_PRICE  = 2 //限价委托

)

*/
func NewEntrustQuene(quene_id string) *EntrustQuene {
	m := &EntrustQuene{
		TokenQueneId: quene_id,
		BuyQueneId:   fmt.Sprintf("%s:1", quene_id),
		SellQueneId:  fmt.Sprintf("%s:2", quene_id),
		UUID:         1,
		sourceData:   make(map[string]*EntrustData),
		//newOrderDetail:    make(chan *EntrustData, 1000),
		waitOrderDetail:   make(chan *EntrustData, 1000),
		updateOrderDetail: make(chan *EntrustData, 1000),
		marketOrderDetail: make(chan *EntrustData, 1000),
	}
	go m.process()
	return m
}

//获取自增ID
func (s *EntrustQuene) GetUUID() int64 {
	return atomic.AddInt64(&s.UUID, 1)
}

//委托请求检查
func (s *EntrustQuene) EntrustReq(p *proto.EntrustOrderRequest) (ret int32, err error) {
	g := &EntrustDetail{
		EntrustId:  genkey.GetTimeUnionKey(s.GetUUID()),
		TokenId:    s.TokenId,
		Uid:        int(p.Uid),
		AllNum:     p.Num,
		SurplusNum: p.Num,
		Opt:        int(p.Opt),
		OnPrice:    p.OnPrice,
		States:     0,
	}

	m := &UserToken{}
	err = m.GetUserToken(int(p.Uid), int(p.TokenId))
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	if m.Balance < p.Num { //检查余额
		ret = ERR_TOKEN_LESS
		return
	}

	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()

	//冻结资金
	ret, err = m.SubMoneyWithFronzen(session, p.Num, g.EntrustId)
	if err != nil || ret != ERRCODE_SUCCESS {
		session.Rollback()
		return
	}

	//记录委托
	err = g.Insert(session)
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		return
	}

	//交易流水
	err = new(MoneyRecord).InsertRecord(session, &MoneyRecord{
		Uid:     int(p.Uid),
		TokenId: int(p.TokenId),
		Ukey:    g.EntrustId,
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
		Type:    MONEY_UKEY_TYPE_ENTRUST,
	})
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		return
	}

	err = session.Commit()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	//委托请求进入等待处理队列
	s.waitOrderDetail <- &EntrustData{
		EntrustId:  g.EntrustId,
		Uid:        g.Uid,
		SurplusNum: g.SurplusNum,
		Opt:        p.Opt,
		OnPrice:    g.Price,
		States:     0,
	}

	return
}

//开始交易加入买入USDT-》BTC  ，卖出USDT-》BTC
func (s *EntrustQuene) MakeDeal(buyer *EntrustData, seller *EntrustData, price int64, deal_num int64) (err error) {
	//var ret int32
	buy_token_account := &UserToken{} //买方主账户余额
	err = buy_token_account.GetUserToken(buyer.Uid, s.TokenId)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	buy_trade_token_account := &UserToken{} //买方交易账户余额
	err = buy_trade_token_account.GetUserToken(buyer.Uid, s.TokenTradeId)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	sell_token_account := &UserToken{} //卖方主账户余额
	err = sell_token_account.GetUserToken(seller.Uid, s.TokenId)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	sell_trade_token_account := &UserToken{} //卖方交易账户余额
	err = sell_trade_token_account.GetUserToken(buyer.Uid, s.TokenTradeId)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	t := &Trade{
		Uid:          buyer.Uid,
		TokenId:      s.TokenId,
		TokenTradeId: s.TokenTradeId,
		Price:        price,
		Num:          deal_num,
		Fee:          price / 1000,
		DealTime:     time.Now().Unix(),
	}

	o := &Trade{
		Uid:          seller.Uid,
		TokenId:      s.TokenId,
		TokenTradeId: s.TokenTradeId,
		Price:        price,
		Num:          deal_num,
		Fee:          price / 1000,
		DealTime:     time.Now().Unix(),
	}

	num := deal_num * seller.OnPrice //计算此次交易USDT

	if seller.SurplusNum < deal_num { //卖方部分成交
		t.States = TRADE_STATES_PART
		o.States = TRADE_STATES_ALL
		buyer.SurplusNum -= num
	} else if seller.SurplusNum == deal_num {
		t.States = TRADE_STATES_ALL
		o.States = TRADE_STATES_ALL
	} else {
		t.States = TRADE_STATES_ALL
		o.States = TRADE_STATES_PART
		seller.SurplusNum -= deal_num
	}

	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()

	//USDT left num

	_, err = buy_token_account.NotifyDelFronzen(session, num, buyer.EntrustId)
	if err != nil {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	err = buy_trade_token_account.AddMoney(session, deal_num)
	if err != nil {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}
	_, err = sell_trade_token_account.NotifyDelFronzen(session, deal_num, seller.EntrustId)
	if err != nil {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	err = sell_token_account.AddMoney(session, deal_num)
	if err != nil {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	err = new(Trade).Insert(session, t, o)
	if err != nil {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	err = session.Commit()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	return
}

func (s *EntrustQuene) match(p *EntrustData) (ret int, err error) {
	var other *EntrustData
	if p.Opt == proto.ENTRUST_OPT_BUY {
		other, err = s.popFirstEntrust(proto.ENTRUST_OPT_SELL)
		if err == redis.Nil {
			//没有对应委托单进入等待区
			s.marketOrderDetail <- p
			return
		} else if err != nil {
			Log.Errorln(err.Error())
			return
		}

		//市价交易撮合
		if p.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
			num := p.SurplusNum / other.OnPrice

			if num > other.SurplusNum {//存在限价则成交
				s.MakeDeal(p, other, other.OnPrice, other.SurplusNum)
				s.match(p)
			} else if num == other.SurplusNum {
				s.MakeDeal(p, other, other.OnPrice, other.SurplusNum)
			} else {
				s.MakeDeal(p, other, other.OnPrice, p.SurplusNum)
				s.joinSellQuene(other)
			}

		} else if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE { //限价交易撮合
			if p.OnPrice >= other.OnPrice {
				if p.OnPrice <= s.price {
					s.price = p.OnPrice
				} else if other.OnPrice >= s.price {
					s.price = other.OnPrice
				} else if s.price > p.OnPrice && s.price < other.OnPrice {
					s.price = s.price
				}

				if p.SurplusNum > other.SurplusNum {
					s.MakeDeal(p, other, s.price, other.SurplusNum)
					s.match(p)
				} else if p.SurplusNum == other.SurplusNum {
					s.MakeDeal(p, other, s.price, other.SurplusNum)
				} else {
					s.MakeDeal(p, other, s.price, p.SurplusNum)
					s.joinSellQuene(other)
				}
			}
		}

	}else 	if p.Opt == proto.ENTRUST_OPT_SELL {
		other, err = s.popFirstEntrust(proto.ENTRUST_OPT_BUY)
		if err == redis.Nil {
			//没有对应委托单进入等待区
			s.marketOrderDetail <- p
			return
		} else if err != nil {
			Log.Errorln(err.Error())
			return
		}

		if p.Type == proto.ENTRUST_TYPE_MARKET_PRICE {//市价交易撮合

			if p.SurplusNum > other.SurplusNum {//存在限价则成交
				s.MakeDeal(p, other, other.OnPrice, other.SurplusNum)
				s.match(p)
			} else if p.SurplusNum  == other.SurplusNum {
				s.MakeDeal(p, other, other.OnPrice, other.SurplusNum)
			} else {
				s.MakeDeal(p, other, other.OnPrice, p.SurplusNum)
				s.joinSellQuene(other)
			}

		}else if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE { //限价交易撮合
			if p.OnPrice >= other.OnPrice {
				if p.OnPrice <= s.price {
					s.price = p.OnPrice
				} else if other.OnPrice >= s.price {
					s.price = other.OnPrice
				} else if s.price > p.OnPrice && s.price < other.OnPrice {
					s.price = s.price
				}

				if p.SurplusNum > other.SurplusNum {
					s.MakeDeal(p, other, s.price, other.SurplusNum)
					s.match(p)
				} else if p.SurplusNum == other.SurplusNum {
					s.MakeDeal(p, other, s.price, other.SurplusNum)
				} else {
					s.MakeDeal(p, other, s.price, p.SurplusNum)
					s.joinSellQuene(other)
				}
			}
		}

	}

	return
}

//限价委托入队列 opt 0 buy ,1 sell
func (s *EntrustQuene) joinSellQuene(p *EntrustData) (ret int, err error) {
	if p.Opt > proto.ENTRUST_OPT_EOMAX {
		ret = ERRCODE_PARAM
		return
	}

	var quene_id string
	if p.Opt == proto.ENTRUST_OPT_BUY {
		quene_id = s.BuyQueneId
	} else if p.Opt == proto.ENTRUST_OPT_SELL {
		quene_id = s.SellQueneId
	}

	//may be not exact
	x := convert.Int64ToFloat64By8Bit(p.OnPrice)

	err = DB.GetRedisConn().ZAdd(quene_id, redis.Z{
		Member: p.Uid,
		Score:  x,
	}).Err()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	if ok := s.insertOrderDetail(p); ok {
	}

	return
}

//获取订单详情
func (s *EntrustQuene) GetOrderData(order_id string) (d *EntrustData, ok bool) {
	d, ok = s.sourceData[order_id]
	if !ok {
		return
	}
	return
}

//保存订单详情
func (s *EntrustQuene) insertOrderDetail(d *EntrustData) bool {
	if _, ok := s.GetOrderData(d.EntrustId); ok {
		return false
	}
	s.sourceData[d.EntrustId] = d
	return true
}

//删除队列数据源模拟弹出操作
func (s *EntrustQuene) delOrderDetail(order_id string) bool {
	delete(s.sourceData, order_id)
	return true
}

//处理请求数据
func (s *EntrustQuene) process() {
	for {
		select {
		case w := <-s.waitOrderDetail:
			s.match(w)
		}
	}
}

//获取弹出队列首位交易单
func (s *EntrustQuene) popFirstEntrust(opt proto.ENTRUST_OPT) (en *EntrustData, err error) {
	var z []redis.Z
	var ok bool
	if opt == proto.ENTRUST_OPT_BUY { //买入类型
		z, err = DB.GetRedisConn().ZRangeWithScores(s.BuyQueneId, 0, 1).Result()
	} else if opt == proto.ENTRUST_OPT_SELL { //卖出类型
		z, err = DB.GetRedisConn().ZRevRangeWithScores(s.BuyQueneId, 0, 1).Result()
	}

	if err != nil {
		Log.Errorln(err)
		return
	}

	if len(z) > 0 {
		d := z[0].Member.(string)
		en, ok = s.GetOrderData(d)
		if ok {
			err = DB.GetRedisConn().ZRem(s.BuyQueneId, d).Err()
			if err != nil {
				return
			}
			s.delOrderDetail(d)
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
