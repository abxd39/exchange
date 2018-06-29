package model

import (
	"digicon/common/convert"
	"digicon/common/encryption"
	"digicon/common/genkey"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	. "digicon/token_service/log"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/liudng/godump"
	"github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

//交易队列类型
type EntrustQuene struct {
	//币币队列ID  格式主要货币_交易货币
	TokenQueueId string

	//队列主要货币ID eg:USDT
	TokenId int
	//队列交易货币ID eg:BTC
	TokenTradeId int
	//卖出队列key
	SellQueueId string
	//买入队列key
	BuyQueueId string

	//当前队列自增ID
	UUID int64

	sourceId string
	//key 是委托ID，委托数据源
	//sourceData map[string]*EntrustData

	//缓存将要保存的DB的委托请求
	//newOrderDetail chan *EntrustData

	//缓存将要更新的DB的委托请求
	//updateOrderDetail chan *EntrustData

	//等待处理的委托请求
	waitOrderDetail chan *EntrustData

	//市价等待队列
	marketOrderDetail chan *EntrustData
	//上一次成交价格
	price int64
}

func GenSourceKey(en string) string {
	return fmt.Sprintf("source:%s", en)
}

func NewEntrustQueue(token_id, token_trade_id int, price int64) *EntrustQuene {
	quene_id := fmt.Sprintf("token:%d/%d", token_id, token_trade_id)
	m := &EntrustQuene{
		TokenQueueId: quene_id,
		BuyQueueId:   fmt.Sprintf("%s:1", quene_id),
		SellQueueId:  fmt.Sprintf("%s:2", quene_id),
		TokenId:      token_id,
		TokenTradeId: token_trade_id,
		UUID:         1,
		//sourceData:   make(map[string]*EntrustData),
		//newOrderDetail:    make(chan *EntrustData, 1000),
		waitOrderDetail: make(chan *EntrustData, 1000),
		//updateOrderDetail: make(chan *EntrustData, 1000),
		marketOrderDetail: make(chan *EntrustData, 1000),
		price:             price,
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
		TokenId:    s.TokenTradeId,
		Uid:        p.Uid,
		AllNum:     p.Num,
		SurplusNum: p.Num,
		Opt:        int(p.Opt),
		OnPrice:    p.OnPrice,
		States:     0,
		Type:       int(p.Type),
	}

	m := &UserToken{}
	var token_id int
	if p.Opt == proto.ENTRUST_OPT_BUY {
		token_id = s.TokenId
	} else {
		token_id = s.TokenTradeId
	}

	err = m.GetUserToken(p.Uid, token_id)
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
	ret, err = m.SubMoneyWithFronzen(session, p.Num, g.EntrustId, FROZEN_LOGIC_TYPE_ENTRUST)
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
		Uid:     p.Uid,
		TokenId: token_id,
		Ukey:    g.EntrustId,
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
		Type:    MONEY_UKEY_TYPE_ENTRUST,
		Num:     p.Num,
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
		OnPrice:    g.OnPrice,
		States:     0,
		Type:       p.Type,
	}

	return
}

//开始交易加入举例买入USDT-》BTC  ，卖出USDT-》BTC  ,deal_num 买方实际获得BTC数量
func (s *EntrustQuene) MakeDeal(buyer *EntrustData, seller *EntrustData, price int64, deal_num int64) (err error) {
	//var ret int32
	buy_token_account := &UserToken{} //买方主账户余额 USDT
	err = buy_token_account.GetUserToken(buyer.Uid, s.TokenId)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	buy_trade_token_account := &UserToken{} //买方交易账户余额 BTC
	err = buy_trade_token_account.GetUserToken(buyer.Uid, s.TokenTradeId)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	sell_token_account := &UserToken{} //卖方主账户余额  BTC
	err = sell_token_account.GetUserToken(seller.Uid, s.TokenTradeId)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	sell_trade_token_account := &UserToken{} //卖方交易账户余额 USDT
	err = sell_trade_token_account.GetUserToken(seller.Uid, s.TokenId)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	num := convert.Int64MulInt64By8Bit(deal_num, price)
	fmt.Printf("price =%d,deal_num=%d ,num =%d \n", price, deal_num, num)

	fee := num * 5 / 1000 //买家消耗手续费0.005个USDT
	godump.Dump(fee)

	no := encryption.CreateOrderId(buyer.Uid, int32(s.TokenId))
	t := &Trade{
		TradeNo:      no,
		Uid:          buyer.Uid,
		TokenId:      s.TokenId,
		TokenTradeId: s.TokenTradeId,
		Price:        price,
		Num:          num - fee, //记录消耗本来USDT数量
		Fee:          fee,
		DealTime:     time.Now().Unix(),
		Opt:         int(proto.ENTRUST_OPT_BUY),
	}

	sell_fee := deal_num * 5 / 1000
	o := &Trade{
		TradeNo:      no,
		Uid:          seller.Uid,
		TokenId:      s.TokenId,
		TokenTradeId: s.TokenTradeId,
		Price:        price,
		Num:          deal_num - sell_fee,
		Fee:          sell_fee,
		DealTime:     time.Now().Unix(),
		Opt:         int(proto.ENTRUST_OPT_SELL),
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
	godump.Dump(time.Now().Unix())
	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()

	//USDT left num
	_, err = buy_token_account.NotifyDelFronzen(session, num, t.TradeNo, FROZEN_LOGIC_TYPE_DEAL)
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

	err = new(MoneyRecord).InsertRecord(session, &MoneyRecord{
		Uid:     buyer.Uid,
		TokenId: buy_trade_token_account.TokenId,
		Ukey:    t.TradeNo,
		Opt:     int(proto.ENTRUST_OPT_BUY),
		Type:    MONEY_UKEY_TYPE_TRADE,
		Num:     deal_num,
	})

	if err != nil {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	if buyer.Uid == seller.Uid { //还没处理
		sell_trade_token_account = buy_token_account
		sell_token_account = buy_trade_token_account

	}
	_, err = sell_token_account.NotifyDelFronzen(session, deal_num, o.TradeNo, FROZEN_LOGIC_TYPE_DEAL)
	if err != nil {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	err = sell_trade_token_account.AddMoney(session, num)
	if err != nil {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	err = new(MoneyRecord).InsertRecord(session, &MoneyRecord{
		Uid:     seller.Uid,
		TokenId: sell_trade_token_account.TokenId,
		Ukey:    o.TradeNo,
		Opt:     int(proto.ENTRUST_OPT_SELL),
		Type:    MONEY_UKEY_TYPE_TRADE,
		Num:     num,
	})
	if err != nil {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	Log.Debugf("uid=? and other=?",t.Uid,o.Uid)
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
	godump.Dump(time.Now().Unix())
	return
}

func (s *EntrustQuene) match(p *EntrustData) (ret int, err error) {
	godump.Dump(p)
	var other *EntrustData
	if p.Opt == proto.ENTRUST_OPT_BUY {
		other, err = s.popFirstEntrust(proto.ENTRUST_OPT_SELL)
		if err == redis.Nil {
			fmt.Printf("get pop nil time=%d", time.Now().Unix())
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
			//num := p.SurplusNum / other.OnPrice //买房愿意用花的比例兑换BTC的数量

			s.delSource(proto.ENTRUST_OPT_SELL, other.EntrustId)

			num := convert.Int64DivInt64By8Bit(p.SurplusNum, other.OnPrice)
			if num > other.SurplusNum { //存在限价则成交
				s.MakeDeal(p, other, other.OnPrice, other.SurplusNum)
				s.match(p)
			} else if num == other.SurplusNum {
				s.MakeDeal(p, other, other.OnPrice, other.SurplusNum)
			} else {
				s.MakeDeal(p, other, other.OnPrice, num)
				s.joinSellQuene(other)
			}
			return

		} else if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE { //限价交易撮合
			if p.OnPrice >= other.OnPrice {

				s.delSource(proto.ENTRUST_OPT_SELL, other.EntrustId)

				godump.Dump("match1")
				if p.OnPrice <= s.price {
					s.price = p.OnPrice
				} else if other.OnPrice >= s.price {
					s.price = other.OnPrice
				} else if s.price > p.OnPrice && s.price < other.OnPrice {
					s.price = s.price
				}

				//num := p.SurplusNum / s.price//买房愿意用花的比例兑换BTC的数量
				num := convert.Int64DivInt64By8Bit(p.SurplusNum, s.price)
				godump.Dump(num)
				if num > other.SurplusNum {
					s.MakeDeal(p, other, s.price, other.SurplusNum)
					s.match(p)

				} else if num == other.SurplusNum {
					s.MakeDeal(p, other, s.price, num)

				} else {
					s.MakeDeal(p, other, s.price, num)
					s.joinSellQuene(other)
				}
				return
			}
		}

	} else if p.Opt == proto.ENTRUST_OPT_SELL {
		other, err = s.popFirstEntrust(proto.ENTRUST_OPT_BUY)
		if err == redis.Nil {
			//没有对应委托单进入等待区
			fmt.Printf("get pop nil time=%d", time.Now().Unix())
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
			s.delSource(proto.ENTRUST_OPT_SELL, other.EntrustId)

			num := convert.Int64DivInt64By8Bit(other.SurplusNum, other.OnPrice) //买房愿意用花的USDT比例兑换BTC的数量
			if num > p.SurplusNum {                                             //存在限价则成交
				s.MakeDeal(p, other, other.OnPrice, num)
				s.match(p)

			} else if num == p.SurplusNum {
				s.MakeDeal(p, other, other.OnPrice, num)
			} else {
				s.MakeDeal(p, other, other.OnPrice, num)
				s.joinSellQuene(other)
			}
			return
		} else if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE { //限价交易撮合
			if p.OnPrice >= other.OnPrice {
				godump.Dump("match2")

				s.delSource(proto.ENTRUST_OPT_SELL, other.EntrustId)

				if p.OnPrice <= s.price {
					s.price = p.OnPrice
				} else if other.OnPrice >= s.price {
					s.price = other.OnPrice
				} else if s.price > p.OnPrice && s.price < other.OnPrice {
					s.price = s.price
				}

				num := convert.Int64DivInt64By8Bit(other.SurplusNum, s.price) //买房愿意用花的USDT比例兑换BTC的数量

				if num > p.SurplusNum {
					s.MakeDeal(other, p, s.price, p.SurplusNum)
					s.match(p)

				} else if num == p.SurplusNum {
					s.MakeDeal(other, p, s.price, num)
				} else {
					s.MakeDeal(other, p, s.price, num)
					s.joinSellQuene(other)
				}
				return
			}
		}
	}

	s.joinSellQuene(other)
	s.joinSellQuene(p)
	return
}

/*
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
*/
//处理请求数据
func (s *EntrustQuene) process() {
	for {
		select {
		case w := <-s.waitOrderDetail:
			s.match(w)
		}
	}
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

	b, err := json.Marshal(p)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	err = DB.GetRedisConn().ZAdd(quene_id, redis.Z{
		Member: p.EntrustId,
		Score:  x,
	}).Err()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	rsp := DB.GetRedisConn().Set(GenSourceKey(p.EntrustId), b, 0)
	err = rsp.Err()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	/*
		if ok := s.insertOrderDetail(p); ok {
		}
	*/
	return
}

//弹出数据
func (s *EntrustQuene) delSource(opt proto.ENTRUST_OPT, entrust_id string) (err error) {
	if opt == proto.ENTRUST_OPT_BUY { //买入类型
		err = DB.GetRedisConn().ZRem(s.BuyQueueId, entrust_id).Err()
		if err != nil {
			Log.Errorln(err)
			return
		}
	} else if opt == proto.ENTRUST_OPT_SELL {
		err = DB.GetRedisConn().ZRem(s.SellQueueId, entrust_id).Err()
		if err != nil {
			Log.Errorln(err)
			return
		}
	} else {
		return errors.New("opt param err")
	}

	err = DB.GetRedisConn().Del(GenSourceKey(entrust_id)).Err()
	if err != nil {
		Log.Errorln(err)
		return
	}
	return
}

//获取队列首位交易单
func (s *EntrustQuene) popFirstEntrust(opt proto.ENTRUST_OPT) (en *EntrustData, err error) {
	var z []redis.Z
	//var ok bool
	if opt == proto.ENTRUST_OPT_BUY { //买入类型
		z, err = DB.GetRedisConn().ZRevRangeWithScores(s.BuyQueueId, 0, 0).Result()
	} else if opt == proto.ENTRUST_OPT_SELL { //卖出类型
		z, err = DB.GetRedisConn().ZRangeWithScores(s.BuyQueueId, 0, 0).Result()
	}

	if err != nil {
		Log.Errorln(err)
		return
	}

	if len(z) > 0 {
		d := z[0].Member.(string)
		var b []byte
		b, err = DB.GetRedisConn().Get(GenSourceKey(d)).Bytes()
		if err != nil {
			Log.WithFields(logrus.Fields{
				"en_id": d,
				"err":   err.Error(),
			}).Errorln("print data")
			return
		}

		en = &EntrustData{}
		err = json.Unmarshal(b, en)
		if err != nil {
			Log.Errorln(err)
			return
		}

		return

		/*
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
		*/

	} else {
		err = redis.Nil
		return
	}

	err = errors.New("this is sync data err ")
	Log.WithFields(logrus.Fields{
		"quene_id": s.TokenQueueId,
		"opt":      opt,
	}).Errorln(err.Error())
	return
}
