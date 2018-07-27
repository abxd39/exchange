package model

import (
	"digicon/common/convert"
	"digicon/common/encryption"
	"digicon/common/genkey"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alex023/clock"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/jsonpb"
	"github.com/liudng/godump"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
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
	//价格推送到键值
	PriceChannel string
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

	lock sync.Mutex
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
	//成交量
	amount int64
	//成交额
	vol int64
	//成交笔数
	count int64
	//主币兑人民币价格
	cny int64
	//人民币成交额
	cny_vol int64
}

type TradeInfo struct {
	CreateTime int64
	TradePrice int64
	Num        int64
}

func GenSourceKey(en string) string {
	return fmt.Sprintf("source:%s", en)
}

func NewEntrustQueue(token_id, token_trade_id int, price int64, name string, cny int64, amount, vol, count int64) *EntrustQuene {
	quene_id := name

	m := &EntrustQuene{
		TokenQueueId:      quene_id,
		PriceChannel:      genkey.GetPulishKey(quene_id),
		BuyQueueId:        fmt.Sprintf("%s:1", quene_id),
		SellQueueId:       fmt.Sprintf("%s:2", quene_id),
		MarketBuyQueueId:  fmt.Sprintf("%s:3", quene_id),
		MarketSellQueueId: fmt.Sprintf("%s:4", quene_id),
		TokenId:           token_id,
		TradeQuene:        fmt.Sprintf("%s:trade", quene_id),
		TokenTradeId:      token_trade_id,
		UUID:              1,
		//sourceData:   make(map[string]*EntrustData),
		//newOrderDetail:    make(chan *EntrustData, 1000),
		waitOrderDetail: make(chan *EntrustData, 1000),
		//sellMarketOrderDetail: make(chan *EntrustData, 1000),
		//updateOrderDetail: make(chan *EntrustData, 1000),
		marketOrderDetail: make(chan *EntrustData, 1000),
		price:             price,
		cny:               cny,
	}
	go m.process()
	go m.Clock()
	return m
}

//获取自增ID
func (s *EntrustQuene) GetUUID() int64 {
	return atomic.AddInt64(&s.UUID, 1)
}

//平台自动委托
func (s *EntrustQuene) EntrustAl(p *proto.EntrustOrderRequest) (e *EntrustData, ret int32, err error) {
	g := &EntrustDetail{
		EntrustId:  genkey.GetTimeUnionKey(s.GetUUID()),
		TokenId:    s.TokenTradeId,
		Uid:        p.Uid,
		AllNum:     p.Num,
		SurplusNum: p.Num,
		Opt:        int(p.Opt),
		OnPrice:    p.OnPrice,
		States:     int(proto.TRADE_STATES_TRADE_NONE),
		Type:       int(p.Type),
		Mount:      convert.Int64MulInt64By8Bit(p.Num, p.OnPrice),
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
		log.Errorln(err.Error())
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
		log.Errorln(err.Error())
		session.Rollback()
		return
	}

	//交易流水
	err = InsertRecord(session, &MoneyRecord{
		Uid:     p.Uid,
		TokenId: token_id,
		Ukey:    g.EntrustId,
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
		Type:    MONEY_UKEY_TYPE_ENTRUST,
		Num:     p.Num,
		Balance: m.Balance,
	})
	if err != nil {
		log.Errorln(err.Error())
		session.Rollback()
		return
	}

	err = session.Commit()
	if err != nil {
		log.Errorln(err.Error())
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
func (s *EntrustQuene) SetTradeInfo(price int64, deal_num int64) {
	s.price = price
	s.count += 1
	s.vol += convert.Int64MulInt64By8Bit(price, deal_num)
	s.amount += deal_num

	s.cny_vol = convert.Int64MulInt64By8Bit(s.vol, s.cny)
}

//委托请求检查
func (s *EntrustQuene) EntrustReq(p *proto.EntrustOrderRequest) (ret int32, err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"type":     p.Type,
				"uid":      p.Uid,
				"symbol":   p.Symbol,
				"opt":      p.Opt,
				"on_price": p.OnPrice,
				"num":      p.Num,
			}).Errorf("EntrustReq error %s", err.Error())
		}
	}()
	g := &EntrustDetail{
		EntrustId:  genkey.GetTimeUnionKey(s.GetUUID()),
		TokenId:    s.TokenTradeId,
		Uid:        p.Uid,
		AllNum:     p.Num,
		SurplusNum: p.Num,
		Opt:        int(p.Opt),
		OnPrice:    p.OnPrice,
		States:     int(proto.TRADE_STATES_TRADE_NONE),
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
		log.Errorln(err.Error())
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
		session.Rollback()
		return
	}

	//交易流水
	err = InsertRecord(session, &MoneyRecord{
		Uid:     p.Uid,
		TokenId: token_id,
		Ukey:    g.EntrustId,
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
		Type:    MONEY_UKEY_TYPE_ENTRUST,
		Num:     p.Num,
		Balance: m.Balance,
	})
	if err != nil {
		session.Rollback()
		return
	}

	err = session.Commit()
	if err != nil {
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

//开始交易加入举例买入USDT-》BTC  ，卖出USDT-》BTC  ,deal_num 卖方实际消耗BTC数量
func (s *EntrustQuene) MakeDeal(buyer *EntrustData, seller *EntrustData, price int64, buy_num, deal_num int64) (err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"buy_id":       buyer.Uid,
				"buy_entrust":  buyer.EntrustId,
				"sell_id":      seller.Uid,
				"sell_entrust": seller.EntrustId,
				"price":        price,
				"buy_num":      buy_num,
				"deal_num":     deal_num,
			}).Errorf("MakeDeal error %s", err.Error())
		}
	}()
	//var ret int32
	if buyer.Opt != proto.ENTRUST_OPT_BUY {
		return errors.New("wrong type")
	}

	buy_token_account := &UserToken{} //买方主账户余额 USDT
	err = buy_token_account.GetUserToken(buyer.Uid, s.TokenId)
	if err != nil {
		return
	}

	buy_trade_token_account := &UserToken{} //买方交易账户余额 BTC
	err = buy_trade_token_account.GetUserToken(buyer.Uid, s.TokenTradeId)
	if err != nil {
		return
	}

	sell_token_account := &UserToken{} //卖方主账户余额  BTC
	err = sell_token_account.GetUserToken(seller.Uid, s.TokenTradeId)
	if err != nil {
		return
	}

	sell_trade_token_account := &UserToken{} //卖方交易账户余额 USDT
	err = sell_trade_token_account.GetUserToken(seller.Uid, s.TokenId)
	if err != nil {
		return
	}

	//num := convert.Int64MulInt64By8Bit(deal_num, price) //买家消耗USDT数量
	//fmt.Printf("price =%d,deal_num=%d ,num =%d \n", price, deal_num, num)

	fee := buy_num * 5 / 1000 //买家消耗手续费0.005个USDT

	no := encryption.CreateOrderId(buyer.Uid, int32(s.TokenId))
	trade_time := time.Now().Unix()
	t := &Trade{
		TradeNo:      no,
		Uid:          buyer.Uid,
		TokenId:      s.TokenId,
		TokenTradeId: s.TokenTradeId,
		Price:        price,
		Num:          buy_num - fee, //记录消耗本来USDT数量
		Fee:          fee,
		DealTime:     trade_time,
		Opt:          int(proto.ENTRUST_OPT_BUY),
		TokenName:    s.TokenQueueId,
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
		TokenName:    s.TokenQueueId,
	}

	var buy_surplus, sell_surplus int64
	if seller.SurplusNum < deal_num { //卖方部分成交
		t.States = int(proto.TRADE_STATES_TRADE_PART)
		o.States = int(proto.TRADE_STATES_TRADE_ALL)

	} else if seller.SurplusNum == deal_num && buyer.SurplusNum == buy_num {
		t.States = int(proto.TRADE_STATES_TRADE_ALL)
		o.States = int(proto.TRADE_STATES_TRADE_ALL)

	} else {
		t.States = int(proto.TRADE_STATES_TRADE_ALL)
		o.States = int(proto.TRADE_STATES_TRADE_PART)

	}
	buy_surplus = buyer.SurplusNum - buy_num
	sell_surplus = seller.SurplusNum - deal_num
	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()
	var ret int32
	//USDT left num

	ret, err = buy_token_account.NotifyDelFronzen(session, buy_num, t.TradeNo, FROZEN_LOGIC_TYPE_DEAL)
	if err != nil {
		session.Rollback()
		return
	}
	if ret != ERRCODE_SUCCESS {
		session.Rollback()
		return
	}

	err = buy_trade_token_account.AddMoney(session, t.Num)
	if err != nil {
		session.Rollback()
		return
	}

	err = InsertRecord(session, &MoneyRecord{
		Uid:     buyer.Uid,
		TokenId: buy_trade_token_account.TokenId,
		Ukey:    t.TradeNo,
		Opt:     int(proto.ENTRUST_OPT_BUY),
		Type:    MONEY_UKEY_TYPE_TRADE,
		Num:     deal_num,
		Balance: buy_trade_token_account.Balance,
	})

	if err != nil {
		session.Rollback()
		return
	}

	if buyer.Uid == seller.Uid { //还没处理
		sell_trade_token_account = buy_token_account
		sell_token_account = buy_trade_token_account

	}
	ret, err = sell_token_account.NotifyDelFronzen(session, deal_num, o.TradeNo, FROZEN_LOGIC_TYPE_DEAL)
	if err != nil || ret != ERRCODE_SUCCESS {
		session.Rollback()
		return
	}

	err = sell_trade_token_account.AddMoney(session, o.Num)
	if err != nil {
		session.Rollback()
		return
	}

	err = InsertRecord(session, &MoneyRecord{
		Uid:     seller.Uid,
		TokenId: sell_trade_token_account.TokenId,
		Ukey:    o.TradeNo,
		Opt:     int(proto.ENTRUST_OPT_SELL),
		Type:    MONEY_UKEY_TYPE_TRADE,
		Num:     buy_num,
		Balance: sell_trade_token_account.Balance,
	})
	if err != nil {
		session.Rollback()
		return
	}

	err = new(Trade).Insert(session, t, o)
	if err != nil {
		session.Rollback()
		return
	}

	err = new(EntrustDetail).UpdateStates(session, buyer.EntrustId, t.States, buy_surplus)
	if err != nil {
		session.Rollback()
		return
	}

	err = new(EntrustDetail).UpdateStates(session, seller.EntrustId, o.States, sell_surplus)
	if err != nil {
		session.Rollback()
		return
	}
	err = session.Commit()
	if err != nil {
		return
	}

	b, err := json.Marshal(&TradeInfo{
		CreateTime: trade_time,
		TradePrice: price,
		Num:        deal_num,
	})
	if err != nil {
		return
	}

	err = DB.GetRedisConn().LPush(s.TradeQuene, b).Err()
	if err != nil {
		return
	}

	return
}

func (s *EntrustQuene) match2(p *EntrustData) (err error) {
	var buyer *EntrustData
	var seller *EntrustData
	var others []*EntrustData

	defer func() (err2 error) {
		if err == redis.Nil {
			if buyer != nil && buyer.SurplusNum > 0 {
				log.WithFields(logrus.Fields{
					"buyer_id":         buyer.Uid,
					"buyer_entrust_id": buyer.EntrustId,
					"sulplus":          buyer.SurplusNum,
					"os_id":            os.Getpid(),
				}).Info("NULL quene match")
				s.joinSellQuene(buyer)
			}

			if seller != nil && seller.SurplusNum > 0 {
				log.WithFields(logrus.Fields{
					"seller_id":         seller.Uid,
					"seller_entrust_id": seller.EntrustId,
					"sulplus":           seller.SurplusNum,
					"os_id":             os.Getpid(),
				}).Info("NULL quene match")
				s.joinSellQuene(seller)
			}
		} else if err != nil {

			if buyer != nil && buyer.SurplusNum > 0 {
				log.WithFields(logrus.Fields{
					"buy_uid":        buyer.Uid,
					"buy_entrust_id": buyer.EntrustId,
					"sulplus":        buyer.SurplusNum,
					"os_id":          os.Getpid(),
				}).Errorln(err.Error())
				s.joinSellQuene(buyer)
			}

			if seller != nil && seller.SurplusNum > 0 {
				log.WithFields(logrus.Fields{
					"sell_uid":        seller.Uid,
					"sell_entrust_id": seller.EntrustId,
					"sulplus":         seller.SurplusNum,
					"os_id":           os.Getpid(),
				}).Errorln(err.Error())
				s.joinSellQuene(seller)
			}

		} else {
			if buyer != nil && buyer.SurplusNum > 0 {
				return s.match2(buyer)
			}

			if seller != nil && seller.SurplusNum > 0 {
				return s.match2(seller)
			}
		}
		return
	}()

	if p.Opt == proto.ENTRUST_OPT_BUY {
		buyer = p

		others, err = s.PopFirstEntrust(proto.ENTRUST_OPT_SELL, 1, 1)
		if err != nil {
			return
		}
		if len(others) > 0 {
			seller = others[0]
		} else {
			return
		}

	} else {
		seller = p
		others, err = s.PopFirstEntrust(proto.ENTRUST_OPT_BUY, 1, 1)
		if err != nil {
			return
		}
		if len(others) > 0 {
			buyer = others[0]
		} else {
			return
		}
	}

	var buy_num, sell_num, g_num, price int64 //BTC数量，成交价格

	if buyer.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
		if seller.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
			price = s.price
			g_num = convert.Int64DivInt64By8Bit(buyer.SurplusNum, price)
		} else {
			price = seller.OnPrice
			g_num = convert.Int64DivInt64By8Bit(buyer.SurplusNum, seller.OnPrice)
		}
	} else {
		if seller.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
			price = buyer.OnPrice
			g_num = convert.Int64DivInt64By8Bit(buyer.SurplusNum, buyer.OnPrice)
		} else {
			//检查价格匹配规则
			if buyer.OnPrice >= seller.OnPrice {
				if buyer.OnPrice <= s.price {
					price = buyer.OnPrice
				} else if seller.OnPrice >= s.price {
					price = seller.OnPrice
				} else if s.price > buyer.OnPrice && s.price < seller.OnPrice {
					price = s.price
				}
				g_num = convert.Int64DivInt64By8Bit(buyer.SurplusNum, price)
			} else {
				return
			}

		}
	}

	//计算交易数量
	if g_num > seller.SurplusNum {
		buy_num = convert.Int64MulInt64By8Bit(seller.SurplusNum, price)
		sell_num = seller.SurplusNum

	} else if g_num == seller.SurplusNum {
		buy_num = convert.Int64MulInt64By8Bit(seller.SurplusNum, price)
		sell_num = seller.SurplusNum
	} else {
		buy_num = convert.Int64MulInt64By8Bit(g_num, price)
		sell_num = g_num
	}

	log.WithFields(logrus.Fields{
		"buyer_id":         buyer.Uid,
		"seller_id":        seller.Uid,
		"buyer_entrust_id": buyer.EntrustId,
		"sell_entrust_id":  seller.EntrustId,
		"sell_num":         sell_num,
		"buy_num":          buy_num,
		"price":            price,
		"g_num":            g_num,
		"buy_surplus_num":  buyer.SurplusNum,
		"sell_surplus_num": seller.SurplusNum,
		"os_id":            os.Getpid(),
	}).Info("record match trade")

	if buy_num == 0 || sell_num == 0 || price == 0 {

		log.WithFields(logrus.Fields{
			"symbol":           s.TokenQueueId,
			"buyer_id":         buyer.Uid,
			"seller_id":        seller.Uid,
			"buyer_entrust_id": buyer.EntrustId,
			"sell_entrust_id":  seller.EntrustId,
			"sell_num":         sell_num,
			"buy_num":          buy_num,
			"price":            price,
			"g_num":            g_num,
			"buyer_type":       buyer.Type,
			"seller_type":      seller.Type,
		}).Info("please check logic")
		err = errors.New("please check logic")
		return
	}
	err = s.delSource(others[0].Opt, others[0].Type, others[0].EntrustId)
	if err != nil {
		return
	}

	err = s.MakeDeal(buyer, seller, price, buy_num, sell_num)
	if err != nil {
		return
	}

	s.SetTradeInfo(price, sell_num)

	buyer.SurplusNum -= buy_num
	seller.SurplusNum -= sell_num

	log.WithFields(logrus.Fields{
		"buyer_id":         buyer.Uid,
		"seller_id":        seller.Uid,
		"buyer_entrust_id": buyer.EntrustId,
		"sell_entrust_id":  seller.EntrustId,
		"buy_surplus_num":  buyer.SurplusNum,
		"sell_surplus_num": seller.SurplusNum,
		"os_id":            os.Getpid(),
	}).Info("finish match trade")
	return
}

func (s *EntrustQuene) SurplusBack(e *EntrustData) (err error) {
	u := &UserToken{}
	entry := GetEntrust(e.EntrustId)
	if entry != nil {
		err = u.GetUserToken(e.Uid, entry.TokenId)
		if err != nil {
			return nil
		}

		session := DB.GetMysqlConn().NewSession()
		defer session.Close()
		err = session.Begin()
		if err != nil {
			return nil
		}
		err = u.ReturnFronzen(session, e.SurplusNum, e.EntrustId, proto.TOKEN_TYPE_OPERATOR_HISTORY_FRONZE_SYS_SURPLUS)
		if err != nil {
			session.Rollback()
			return nil
		}

	}
	return nil
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
					return
				}

				if ret != ERRCODE_SUCCESS {
					s.joinSellQuene(p)
					return
				}
				other = d
			}

		} else if err != nil {
			log.Errorln(err.Error())
			return
		} else {
			if len(others) == 0 {
				log.WithFields(logrus.Fields{
					"entrust_id": p.EntrustId,
				}).Info("match get other data")
				s.joinSellQuene(p)
				return
			} else {
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
				buy_num := convert.Int64MulInt64By8Bit(other.SurplusNum, price)
				err = s.MakeDeal(p, other, price, buy_num, other.SurplusNum)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}

				s.SetTradeInfo(price, other.SurplusNum)
				other.SurplusNum -= other.SurplusNum
				p.SurplusNum -= buy_num
				s.match(p)

			} else if num == other.SurplusNum {
				buy_num := convert.Int64MulInt64By8Bit(num, price)
				err = s.MakeDeal(p, other, price, buy_num, num)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}
			} else {
				buy_num := convert.Int64MulInt64By8Bit(num, price)
				err = s.MakeDeal(p, other, price, buy_num, num)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}
				s.SetTradeInfo(price, num)
				other.SurplusNum -= other.SurplusNum
				p.SurplusNum -= buy_num
				s.match(other)
			}
			return

		} else if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE { //限价交易撮合
			var num, price int64 //BTC数量，成交价格

			if other.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
				s.delSource(other.Opt, other.Type, other.EntrustId)
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
					s.delSource(other.Opt, other.Type, other.EntrustId)
				} else {
					s.joinSellQuene(p)
					return
				}

				num = convert.Int64DivInt64By8Bit(p.SurplusNum, price) //计算买家最大买入BTC数量
			}

			if num > other.SurplusNum {

				buy_num := convert.Int64MulInt64By8Bit(other.SurplusNum, price)
				err = s.MakeDeal(p, other, price, buy_num, other.SurplusNum)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}
				s.SetTradeInfo(price, other.SurplusNum)
				other.SurplusNum -= other.SurplusNum
				p.SurplusNum -= buy_num
				s.match(p)

			} else if num == other.SurplusNum {
				buy_num := convert.Int64MulInt64By8Bit(other.SurplusNum, price)
				err = s.MakeDeal(p, other, price, buy_num, num)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}
				s.SetTradeInfo(price, num)
			} else {
				buy_num := convert.Int64MulInt64By8Bit(other.SurplusNum, price)
				err = s.MakeDeal(p, other, price, buy_num, num)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}
				s.SetTradeInfo(price, num)
				other.SurplusNum -= other.SurplusNum
				p.SurplusNum -= buy_num
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
					log.Errorln(err.Error())
					return
				}

				if ret != ERRCODE_SUCCESS {
					s.joinSellQuene(p)
					return
				}

				other = d
			}

		} else if err != nil {
			log.Errorln(err.Error())
			return
		} else {
			if len(others) == 0 {
				s.joinSellQuene(p)
				log.WithFields(logrus.Fields{
					"entrust_id": p.EntrustId,
				}).Errorln(err.Error())
				return
			} else {
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

				buy_num := convert.Int64MulInt64By8Bit(p.SurplusNum, price)
				err = s.MakeDeal(other, p, price, buy_num, p.SurplusNum)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}
				s.SetTradeInfo(price, p.SurplusNum)
				other.SurplusNum -= buy_num
				p.SurplusNum -= p.SurplusNum
				s.match(other)

			} else if num == p.SurplusNum {
				buy_num := convert.Int64MulInt64By8Bit(p.SurplusNum, price)
				err = s.MakeDeal(other, p, price, buy_num, num)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}
				s.SetTradeInfo(price, num)
				other.SurplusNum -= buy_num
				p.SurplusNum -= p.SurplusNum
			} else {
				buy_num := convert.Int64MulInt64By8Bit(p.SurplusNum, price)
				err = s.MakeDeal(other, p, price, buy_num, num)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}
				s.SetTradeInfo(price, num)
				other.SurplusNum -= buy_num
				p.SurplusNum -= p.SurplusNum
				s.match(p)
			}
			return
		} else if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE { //限价交易撮合
			var num, price int64 //BTC数量，成交价格

			if other.Type == proto.ENTRUST_TYPE_MARKET_PRICE {
				s.delSource(other.Opt, other.Type, other.EntrustId)
				num = convert.Int64DivInt64By8Bit(other.SurplusNum, p.OnPrice)
				price = p.OnPrice
			} else {
				if other.OnPrice >= p.OnPrice {
					if p.OnPrice <= s.price {
						price = p.OnPrice
					} else if p.OnPrice >= s.price {
						price = p.OnPrice
					} else if s.price < p.OnPrice && s.price > other.OnPrice {
						price = s.price
					}
					s.delSource(other.Opt, other.Type, other.EntrustId)
				} else {
					s.joinSellQuene(p)
					return
				}

				log.WithFields(logrus.Fields{
					"price": price,
					"uid":   p.Uid,
					"num":   p.SurplusNum,
					"oid":   other.Uid,
					"onum":  other.SurplusNum,
				}).Info("print data")
				if price == 0 {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					return
				}

				num = convert.Int64DivInt64By8Bit(other.SurplusNum, price) //买房愿意用花的USDT比例兑换BTC的数量
			}

			if num > p.SurplusNum {

				buy_num := convert.Int64MulInt64By8Bit(p.SurplusNum, price)
				err = s.MakeDeal(other, p, price, buy_num, p.SurplusNum)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}
				s.SetTradeInfo(price, p.SurplusNum)
				other.SurplusNum -= buy_num
				p.SurplusNum -= p.SurplusNum
				s.match(other)

			} else if num == p.SurplusNum {
				buy_num := convert.Int64MulInt64By8Bit(num, price)
				err = s.MakeDeal(other, p, price, buy_num, num)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}
				s.SetTradeInfo(price, num)

			} else {
				buy_num := convert.Int64MulInt64By8Bit(num, price)
				err = s.MakeDeal(other, p, price, buy_num, num)
				if err != nil {
					s.joinSellQuene(p)
					s.joinSellQuene(other)
					log.WithFields(logrus.Fields{
						"uid":          p.Uid,
						"entrust_id":   p.EntrustId,
						"oid":          other.Uid,
						"o_entrust_id": other.EntrustId,
					}).Errorln(err.Error())
					return
				}
				s.SetTradeInfo(price, num)
				other.SurplusNum -= buy_num
				p.SurplusNum -= p.SurplusNum
				s.match(p)
			}
			return

		}
	}

	return
}

//处理请求数据
func (s *EntrustQuene) process() {
	for {
		select {
		case w := <-s.waitOrderDetail:
			s.match2(w)
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
	godump.Dump(s.PriceChannel)
	c := clock.NewClock()
	job := func() {
		m := &proto.PriceCache{
			Id:          time.Now().Unix(),
			Symbol:      s.TokenQueueId,
			Price:       s.price,
			CreatedTime: time.Now().Unix(),
			Amount:      s.amount,
			Count:       s.count,
			Vol:         s.vol,
			CnyVol:      s.cny_vol,
		}

		t := jsonpb.Marshaler{EmitDefaults: true}
		data, err := t.MarshalToString(m)
		if err != nil {
			log.Errorln(err.Error())
			return
		}

		err = DB.GetRedisConn().Publish(s.PriceChannel, data).Err()

		if err != nil {
			log.Errorln(err.Error())
			return
		}

	}

	c.AddJobRepeat(1*time.Second, 0, job)
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
		log.Errorln(err.Error())
		return
	}

	err = DB.GetRedisConn().ZAdd(quene_id, redis.Z{
		Member: p.EntrustId,
		Score:  x,
	}).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	rsp := DB.GetRedisConn().Set(GenSourceKey(p.EntrustId), b, 0)
	err = rsp.Err()
	if err != nil {
		log.Errorln(err.Error())
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
		log.Errorln(err)
		return
	}

	err = DB.GetRedisConn().Del(GenSourceKey(entrust_id)).Err()
	if err != nil {
		log.Errorln(err)
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
		log.Errorln(err)
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
			log.WithFields(logrus.Fields{
				"en_id": d,
				"err":   err.Error(),
			}).Errorln("print data")
			return
		}
		godump.Dump(string(b))
		g := &EntrustData{}
		err = json.Unmarshal(b, g)
		if err != nil {
			log.Errorln(err)
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
					log.WithFields(logrus.Fields{
						"en_id": d,
						"err":   err.Error(),
					}).Errorln("print data")
					return
				}

				en = &EntrustData{}
				err = json.Unmarshal(b, en)
				if err != nil {
					log.Errorln(err)
					return
				}
				return

			} else {
				err = redis.Nil
				return
			}

		err = errors.New("this is sync data err ")
		log.WithFields(logrus.Fields{
			"quene_id": s.TokenQueueId,
			"opt":      opt,
		}).Errorln(err.Error())
	*/
	return
}

func (s *EntrustQuene) GetTradeList(count int64) []*TradeInfo {
	r, err := DB.GetRedisConn().LRange(s.TradeQuene, 0, count).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		log.Errorln(err)
		return nil
	}
	g := make([]*TradeInfo, 0)
	for _, v := range r {
		data := &TradeInfo{}

		err = json.Unmarshal([]byte(v), data)
		if err != nil {
			log.Errorln(err)
			return nil
		}
		g = append(g, data)
	}
	return g
}

func (s *EntrustQuene) GetCnyPrice(price int64) string {
	return convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(s.cny, price))
}

//撤销委托
func (s *EntrustQuene) DelEntrust(e *EntrustDetail) (err error) {
	u := &UserToken{}
	err = u.GetUserToken(e.Uid, e.TokenId)
	if err != nil {
		return err
	}

	sess := DB.GetMysqlConn().NewSession()

	e.States = int(proto.TRADE_STATES_TRADE_DEL)
	_, err = sess.Where("entrust_id=?", e.EntrustId).Update(e)
	if err != nil {
		sess.Rollback()
		return err
	}
	var ret int32
	err = u.ReturnFronzen(sess, e.SurplusNum, e.EntrustId, proto.TOKEN_TYPE_OPERATOR_HISTORY_ENTRUST)
	if err != nil {
		sess.Rollback()
		return err
	}

	if ret != ERRCODE_SUCCESS {
		sess.Rollback()
		return nil
	}
	err = sess.Commit()
	if err != nil {
		return err
	}

	err = s.delSource(proto.ENTRUST_OPT(e.Opt), proto.ENTRUST_TYPE(e.Type), e.EntrustId)
	if err != nil {
		return err
	}
	return nil
}
