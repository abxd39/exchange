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
	TokenQueueId string

	TokenId      int
	TokenTradeId int
	//卖出队列key
	SellQueueId string
	//买入队列key
	BuyQueueId string

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


func NewEntrustQueue(token_id, token_trade_id int,price int64) *EntrustQuene {
	quene_id:=fmt.Sprintf("%s%s",token_id,token_trade_id)
	m := &EntrustQuene{
		TokenQueueId: quene_id,
		BuyQueueId:   fmt.Sprintf("%s:1", quene_id),
		SellQueueId:  fmt.Sprintf("%s:2", quene_id),
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

//开始交易加入举例买入USDT-》BTC  ，卖出USDT-》BTC  ,deal_num 买方实际获得BTC数量
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
	err = sell_token_account.GetUserToken(seller.Uid, s.TokenTradeId)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	sell_trade_token_account := &UserToken{} //卖方交易账户余额
	err = sell_trade_token_account.GetUserToken(seller.Uid, s.TokenId)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	num := deal_num * price //计算此次买家交易USDT 1

	fee := num * 5 / 1000 //买家消耗手续费0.005个USDT
	t := &Trade{
		Uid:          buyer.Uid,
		TokenId:      s.TokenId,
		TokenTradeId: s.TokenTradeId,
		Price:        price,
		Num:          num - fee, //记录消耗本来USDT数量
		Fee:          fee,
		DealTime:     time.Now().Unix(),
		Type:         int(proto.ENTRUST_OPT_BUY),
	}

	sell_fee := deal_num * 5 / 1000
	o := &Trade{
		Uid:          seller.Uid,
		TokenId:      s.TokenId,
		TokenTradeId: s.TokenTradeId,
		Price:        price,
		Num:          deal_num - sell_fee,
		Fee:          sell_fee,
		DealTime:     time.Now().Unix(),
		Type:         int(proto.ENTRUST_OPT_SELL),
	}

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

	err = sell_token_account.AddMoney(session, num)
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
			if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE {
				s.joinSellQuene(p)
			} else {
				//平台吃单
				s.marketOrderDetail <- p
			}
			return
		} else if err != nil {
			Log.Errorln(err.Error())
			return
		}

		//市价交易撮合
		if p.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
			num := p.SurplusNum / other.OnPrice //买房愿意用花的比例兑换BTC的数量

			if num > other.SurplusNum { //存在限价则成交
				s.MakeDeal(p, other, other.OnPrice, other.SurplusNum)
				s.match(p)
			} else if num == other.SurplusNum {
				s.MakeDeal(p, other, other.OnPrice, other.SurplusNum)
			} else {
				s.MakeDeal(p, other, other.OnPrice, p.SurplusNum)
				s.joinSellQuene(other)
			}
			return
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
				return
			}

		}

	} else if p.Opt == proto.ENTRUST_OPT_SELL {
		other, err = s.popFirstEntrust(proto.ENTRUST_OPT_BUY)
		if err == redis.Nil {
			//没有对应委托单进入等待区
			if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE {
				s.joinSellQuene(p)
			} else {
				//平台吃单
				s.marketOrderDetail <- p
			}
			return
		} else if err != nil {
			Log.Errorln(err.Error())
			return
		}

		if p.Type == proto.ENTRUST_TYPE_MARKET_PRICE { //市价交易撮合

			if p.SurplusNum > other.SurplusNum { //存在限价则成交
				s.MakeDeal(p, other, other.OnPrice, other.SurplusNum)
				s.match(p)

			} else if p.SurplusNum == other.SurplusNum {
				s.MakeDeal(p, other, other.OnPrice, other.SurplusNum)
			} else {
				s.MakeDeal(p, other, other.OnPrice, p.SurplusNum)
				s.joinSellQuene(other)
			}
			return
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
					return
				} else if p.SurplusNum == other.SurplusNum {
					s.MakeDeal(p, other, s.price, other.SurplusNum)
					return
				} else {
					s.MakeDeal(p, other, s.price, p.SurplusNum)
					s.joinSellQuene(other)
					return
				}
			}
		}
	}

	s.joinSellQuene(other)
	s.joinSellQuene(p)
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
		quene_id = s.BuyQueueId
	} else if p.Opt == proto.ENTRUST_OPT_SELL {
		quene_id = s.SellQueueId
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
		z, err = DB.GetRedisConn().ZRangeWithScores(s.BuyQueueId, 0, 1).Result()
	} else if opt == proto.ENTRUST_OPT_SELL { //卖出类型
		z, err = DB.GetRedisConn().ZRevRangeWithScores(s.BuyQueueId, 0, 1).Result()
	}

	if err != nil {
		Log.Errorln(err)
		return
	}

	if len(z) > 0 {
		d := z[0].Member.(string)
		en, ok = s.GetOrderData(d)
		if ok {
			err = DB.GetRedisConn().ZRem(s.BuyQueueId, d).Err()
			if err != nil {
				return
			}
			s.delOrderDetail(d)
			return
		}
		err = errors.New("this is unrealize err when get order detail  ")
		Log.WithFields(logrus.Fields{
			"quene_id": s.TokenQueueId,
			"opt":      opt,
			"member":   d,
		}).Errorln(err.Error())

		return
	}

	err = errors.New("this is sync data err ")
	Log.WithFields(logrus.Fields{
		"quene_id": s.TokenQueueId,
		"opt":      opt,
	}).Errorln(err.Error())
	return
}
