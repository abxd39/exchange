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
	"github.com/alex023/clock"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

const (
	ADMIN_UID = 1
)

//交易队列类型
type EntrustQuene struct {
	//币币队列ID  格式主要货币_交易货币
	TokenQueueId string

	//队列主要货币ID eg:USDT
	TokenId int
	//队列交易货币ID eg:BTC
	TokenTradeId int
	//限价卖出队列key
	SellQueueId string

	//市价卖出委托队列
	MarketSellQueueId string
	//限价买入队列key
	BuyQueueId string

	//市价买入委托队列
	MarketBuyQueueId string

	//实时交易队列
	TradeQuene string
	//当前队列自增ID
	UUID int64

	//sourceId string
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

	//sellMarketOrderDetail chan *EntrustData
	//上一次成交价格
	price int64

	amout int64
}

type TradeInfo struct {
	 CreateTime int64
	 TradePrice int64
	 Num int64

}

func GenSourceKey(en string) string {
	return fmt.Sprintf("source:%s", en)
}

func NewEntrustQueue(token_id, token_trade_id int, price int64, name string) *EntrustQuene {
	quene_id := name

	m := &EntrustQuene{
		TokenQueueId:      quene_id,
		BuyQueueId:        fmt.Sprintf("%s:1", quene_id),
		SellQueueId:       fmt.Sprintf("%s:2", quene_id),
		MarketBuyQueueId:  fmt.Sprintf("%s:3", quene_id),
		MarketSellQueueId: fmt.Sprintf("%s:4", quene_id),
		TokenId:           token_id,
		TradeQuene:fmt.Sprintf("%s:trade",quene_id),
		TokenTradeId:      token_trade_id,
		UUID:              1,
		//sourceData:   make(map[string]*EntrustData),
		//newOrderDetail:    make(chan *EntrustData, 1000),
		waitOrderDetail: make(chan *EntrustData, 1000),
		//sellMarketOrderDetail: make(chan *EntrustData, 1000),
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

func (s *EntrustQuene) EntrustAl(p *proto.EntrustOrderRequest) (e *EntrustData, ret int32, err error) {
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
		Mount: convert.Int64MulInt64By8Bit( p.Num,p.OnPrice),
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

	e = &EntrustData{
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
func (s *EntrustQuene) SetPrice(price int64) {
	s.price = price
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
		Symbol:     p.Symbol,
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
	if buyer.Opt != proto.ENTRUST_OPT_BUY {
		Log.Fatalln("wrong type")
	}
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

	no := encryption.CreateOrderId(buyer.Uid, int32(s.TokenId))
	trade_time:= time.Now().Unix()
	t := &Trade{
		TradeNo:      no,
		Uid:          buyer.Uid,
		TokenId:      s.TokenId,
		TokenTradeId: s.TokenTradeId,
		Price:        price,
		Num:          num - fee, //记录消耗本来USDT数量
		Fee:          fee,
		DealTime:   trade_time ,
		Opt:          int(proto.ENTRUST_OPT_BUY),
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
		DealTime:     trade_time,
		Opt:          int(proto.ENTRUST_OPT_SELL),
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
	var ret int32
	//USDT left num
	ret, err = buy_token_account.NotifyDelFronzen(session, num, t.TradeNo, FROZEN_LOGIC_TYPE_DEAL)
	if err != nil || ret != ERRCODE_SUCCESS {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	err = buy_trade_token_account.AddMoney(session, t.Num)
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
	ret, err = sell_token_account.NotifyDelFronzen(session, deal_num, o.TradeNo, FROZEN_LOGIC_TYPE_DEAL)
	if err != nil || ret != ERRCODE_SUCCESS {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	err = sell_trade_token_account.AddMoney(session, o.Num)
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

	err = new(Trade).Insert(session, t, o)
	if err != nil {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	err = new(EntrustDetail).UpdateStates(session, buyer.EntrustId, t.States, num)
	if err != nil {
		session.Rollback()
		Log.Errorln(err.Error())
		return
	}

	err = new(EntrustDetail).UpdateStates(session, seller.EntrustId, o.States, deal_num)
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

	b,err:=json.Marshal(&TradeInfo{
		CreateTime:trade_time,
		TradePrice:price,
		Num:deal_num,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	err=DB.GetRedisConn().RPush(s.TradeQuene,b).Err()
	if err != nil {
		Log.Fatalln(err.Error())
		return
	}
	return
}

//匹配交易
func (s *EntrustQuene) match(p *EntrustData) (ret int32, err error) {
	var other *EntrustData
	var others []*EntrustData
	if p.Opt == proto.ENTRUST_OPT_BUY {

		others, err = s.PopFirstEntrust(proto.ENTRUST_OPT_SELL, 1, 1)
		if err == redis.Nil {
			fmt.Printf("get pop nil time=%d", time.Now().Unix())
			//没有对应委托单进入等待区
			if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE {
				s.joinSellQuene(p)
				return
			} else {
				//平台吃单
				var d *EntrustData

				d, ret, err = s.EntrustAl(&proto.EntrustOrderRequest{
					Symbol:  s.TokenQueueId,
					OnPrice: s.price,
					Num:     convert.Int64DivInt64By8Bit(p.SurplusNum, s.price),
					Opt:     proto.ENTRUST_OPT_SELL,
					Type:    proto.ENTRUST_TYPE_LIMIT_PRICE,
					Uid:     ADMIN_UID,
				})

				if err != nil {
					Log.Errorln(err.Error())
					return
				}

				if ret != ERRCODE_SUCCESS {
					s.joinSellQuene(p)
					//s.marketOrderDetail <- p
					return
				}

				other = d

			}

		} else if err != nil {
			Log.Errorln(err.Error())
			return
		} else {
			if len(others) == 1 {
				other = others[0]
			}
		}

		//市价交易撮合
		if p.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
			//num := p.SurplusNum / other.OnPrice //买房愿意用花的比例兑换BTC的数量
			s.delSource(other.Opt, other.Type, other.EntrustId)
			var num, price int64 //BTC数量，成交价格

			if other.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
				num = convert.Int64DivInt64By8Bit(p.SurplusNum, s.price) //计算买家最大买入BTC数量
				price = s.price
			} else {
				num = convert.Int64DivInt64By8Bit(p.SurplusNum, other.OnPrice)
				price = other.OnPrice
			}

			if num > other.SurplusNum { //存在对手单则成交
				s.MakeDeal(p, other, price, other.SurplusNum)
				s.SetPrice(price)
				s.match(p)

			} else if num == other.SurplusNum {
				s.MakeDeal(p, other, price, num)
			} else {
				s.MakeDeal(p, other, price, num)
				s.SetPrice(price)
				s.match(other)
			}
			return

		} else if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE { //限价交易撮合
			var num, price int64 //BTC数量，成交价格

			if other.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
				num = convert.Int64DivInt64By8Bit(p.SurplusNum, p.OnPrice)
				price = p.OnPrice
			} else {
				if p.OnPrice >= other.OnPrice {
					if p.OnPrice <= s.price {
						price = p.OnPrice
					} else if other.OnPrice >= s.price {
						price = other.OnPrice
					} else if s.price > p.OnPrice && s.price < other.OnPrice {
						price = s.price
					}
				} else {
					s.joinSellQuene(p)
					return
				}

				num = convert.Int64DivInt64By8Bit(p.SurplusNum, price) //计算买家最大买入BTC数量
			}
			s.delSource(other.Opt, other.Type, other.EntrustId)

			if num > other.SurplusNum {
				err = s.MakeDeal(p, other, price, other.SurplusNum)
				if err != nil {
					Log.Errorln(err.Error())
					return
				}
				s.SetPrice(price)
				s.match(p)

			} else if num == other.SurplusNum {
				err = s.MakeDeal(p, other, price, num)
				if err != nil {
					Log.Errorln(err.Error())
					return
				}
				s.SetPrice(price)
			} else {
				err = s.MakeDeal(p, other, price, num)
				if err != nil {
					Log.Errorln(err.Error())
					return
				}
				s.SetPrice(price)
				s.joinSellQuene(other)
			}
			return
		}

	} else if p.Opt == proto.ENTRUST_OPT_SELL {
		others, err = s.PopFirstEntrust(proto.ENTRUST_OPT_BUY, 1, 1)

		if err == redis.Nil {
			//没有对应委托单进入等待区
			fmt.Printf("get pop nil time=%d", time.Now().Unix())
			if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE {
				s.joinSellQuene(p)
				return
			} else {
				//平台吃单
				var d *EntrustData

				d, ret, err = s.EntrustAl(&proto.EntrustOrderRequest{
					Symbol:  s.TokenQueueId,
					OnPrice: s.price,
					Num:     convert.Int64MulInt64By8Bit(p.SurplusNum, s.price),
					Opt:     proto.ENTRUST_OPT_BUY,
					Type:    proto.ENTRUST_TYPE_LIMIT_PRICE,
					Uid:     ADMIN_UID,
				})

				if err != nil {
					Log.Errorln(err.Error())
					return
				}

				if ret != ERRCODE_SUCCESS {
					s.joinSellQuene(p)
					//s.marketOrderDetail <- p
					return
				}

				other = d
			}

		} else if err != nil {
			Log.Errorln(err.Error())
			return
		} else {
			if len(others) == 1 {
				other = others[0]
			}
		}

		if p.Type == proto.ENTRUST_TYPE_MARKET_PRICE { //市价交易撮合
			s.delSource(other.Opt, other.Type, other.EntrustId)
			var num, price int64
			if other.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
				num = convert.Int64DivInt64By8Bit(other.SurplusNum, s.price) //计算买家最大买入BTC数量
				price = s.price
			} else {
				num = convert.Int64DivInt64By8Bit(other.SurplusNum, other.OnPrice)
				price = other.OnPrice
			}

			//num := convert.Int64DivInt64By8Bit(other.SurplusNum, other.OnPrice) //买房愿意用花的USDT比例兑换BTC的数量
			if num > p.SurplusNum { //存在限价则成交
				err = s.MakeDeal(other, p, price, p.SurplusNum)
				if err != nil {
					Log.Errorln(err.Error())
					return
				}
				s.SetPrice(price)
				s.match(other)

			} else if num == p.SurplusNum {
				err = s.MakeDeal(other, p, price, num)
				if err != nil {
					Log.Errorln(err.Error())
					return
				}
				s.SetPrice(price)
			} else {
				err = s.MakeDeal(other, p, price, num)
				if err != nil {
					Log.Errorln(err.Error())
					return
				}
				s.SetPrice(price)
				s.match(p)
			}
			return
		} else if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE { //限价交易撮合
			var num, price int64 //BTC数量，成交价格

			if other.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
				num = convert.Int64DivInt64By8Bit(other.SurplusNum, p.OnPrice)
				price = p.OnPrice
			} else {
				if p.OnPrice >= other.OnPrice {
					if p.OnPrice <= s.price {
						s.price = p.OnPrice
					} else if other.OnPrice >= s.price {
						price = other.OnPrice
					} else if s.price > p.OnPrice && s.price < other.OnPrice {
						price = s.price
					}
				} else {
					return
				}

				num = convert.Int64DivInt64By8Bit(other.SurplusNum, price) //买房愿意用花的USDT比例兑换BTC的数量
			}

			s.delSource(other.Opt, other.Type, other.EntrustId)

			if num > p.SurplusNum {
				err = s.MakeDeal(other, p, price, p.SurplusNum)
				if err != nil {
					Log.Errorln(err.Error())
					return
				}
				s.SetPrice(price)

				s.match(p)

			} else if num == p.SurplusNum {
				err = s.MakeDeal(other, p, price, num)
				if err != nil {
					Log.Errorln(err.Error())
					return
				}
				s.SetPrice(price)

			} else {
				err = s.MakeDeal(other, p, price, num)
				if err != nil {
					Log.Errorln(err.Error())
					return
				}
				s.SetPrice(price)
				s.joinSellQuene(other)
			}
			return

		}
	}

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
		case d := <-s.marketOrderDetail:
			go func(data *EntrustData) {
				time.Sleep(10 * time.Second)
				s.waitOrderDetail <- d
			}(d)

		}
	}
}

//定时器
func (s *EntrustQuene) Clock() {
	type MinPrice struct {
		Open   int64
		Close  int64
		Low    int64
		High   int64
		Amount int64
		Vol    int64
		Count  int64
	}

	c := clock.NewClock()
	c.AddJobWithInterval(60*time.Second, func() {

	})

}

//委托入队列
func (s *EntrustQuene) joinSellQuene(p *EntrustData) (ret int, err error) {
	if p.Opt > proto.ENTRUST_OPT_EOMAX {
		ret = ERRCODE_PARAM
		return
	}

	var quene_id string
	var x float64
	if p.Opt == proto.ENTRUST_OPT_BUY {
		if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE {
			quene_id = s.BuyQueueId
			x = convert.Int64ToFloat64By8Bit(p.OnPrice)
		} else {
			quene_id = s.MarketBuyQueueId
			x = -float64(time.Now().Unix())
		}

	} else if p.Opt == proto.ENTRUST_OPT_SELL {
		if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE {
			quene_id = s.SellQueueId
			x = convert.Int64ToFloat64By8Bit(p.OnPrice)
		} else {
			quene_id = s.MarketSellQueueId
			x = -float64(time.Now().Unix())
		}

	}

	//may be not exact

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
func (s *EntrustQuene) delSource(opt proto.ENTRUST_OPT, ty proto.ENTRUST_TYPE, entrust_id string) (err error) {
	var quene_id string

	if opt == proto.ENTRUST_OPT_BUY { //买入类型
		if ty == proto.ENTRUST_TYPE_LIMIT_PRICE {
			quene_id = s.BuyQueueId
		} else {
			quene_id = s.MarketBuyQueueId
		}

	} else if opt == proto.ENTRUST_OPT_SELL {
		if ty == proto.ENTRUST_TYPE_LIMIT_PRICE {
			quene_id = s.SellQueueId
		} else {
			quene_id = s.MarketSellQueueId
		}
	} else {
		return errors.New("opt param err")
	}

	err = DB.GetRedisConn().ZRem(quene_id, entrust_id).Err()
	if err != nil {
		Log.Errorln(err)
		return
	}

	err = DB.GetRedisConn().Del(GenSourceKey(entrust_id)).Err()
	if err != nil {
		Log.Errorln(err)
		return
	}
	return
}

//获取队列首位交易单 sw1表示先取市价单再取限价单，2表示直接获取限价单，count获取数量
func (s *EntrustQuene) PopFirstEntrust(opt proto.ENTRUST_OPT, sw int32, count int64) (en []*EntrustData, err error) {
	var z []redis.Z
	var quene_id string
	//var ok bool
	if opt == proto.ENTRUST_OPT_BUY { //买入类型
		if sw == 1 {
			quene_id = s.MarketBuyQueueId
		} else {
			quene_id = s.BuyQueueId
		}

		z, err = DB.GetRedisConn().ZRevRangeWithScores(quene_id, 0, count).Result()
	} else if opt == proto.ENTRUST_OPT_SELL { //卖出类型
		if sw == 1 {
			quene_id = s.MarketSellQueueId
		} else {
			quene_id = s.SellQueueId
		}

		z, err = DB.GetRedisConn().ZRangeWithScores(quene_id, 0, count).Result()
	}

	if err != nil {
		Log.Errorln(err)
		return
	}

	if len(z) == 0 && sw == 1 {
		return s.PopFirstEntrust(opt, 2, count)
	} else if len(z) == 0 && sw == 2 {
		err = redis.Nil
		return
	}

	for _, v := range z {
		d := v.Member.(string)
		var b []byte
		b, err = DB.GetRedisConn().Get(GenSourceKey(d)).Bytes()
		if err != nil {
			Log.WithFields(logrus.Fields{
				"en_id": d,
				"err":   err.Error(),
			}).Errorln("print data")
			return
		}
		g := &EntrustData{}
		err = json.Unmarshal(b, g)
		if err != nil {
			Log.Errorln(err)
			return
		}
		en = append(en, g)
	}

	/*
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

			} else {
				err = redis.Nil
				return
			}

		err = errors.New("this is sync data err ")
		Log.WithFields(logrus.Fields{
			"quene_id": s.TokenQueueId,
			"opt":      opt,
		}).Errorln(err.Error())
	*/
	return
}

func (s *EntrustQuene) GetTradeList(count int64) []*TradeInfo {
	r,err:=DB.GetRedisConn().LRange(s.TradeQuene,0,count).Result()
	if err==redis.Nil {
		return nil
	}else if err!=nil {
		Log.Errorln(err)
		return nil
	}
	g:=make([]*TradeInfo,0)
	for _,v:=range r {
		data:=&TradeInfo{}

		err = json.Unmarshal([]byte(v),data)
		if err!=nil {
			Log.Errorln(err)
			return nil
		}
		g=append(g,data)
	}
	return g
}