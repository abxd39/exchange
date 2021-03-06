package model

import (
	"digicon/common/convert"
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

	"digicon/common/random"
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

	//堆积买入统计
	HeapBuySort string
	//堆积卖出统计
	HeapSellSort string

	//
	HeapBuyHash  string
	HeapSellHash string


	//当前队列自增ID
	UUID int64

	lock sync.Mutex

	//等待处理的委托请求
	waitOrderDetail chan *EntrustDetail

	//市价等待队列
	marketOrderDetail chan *EntrustDetail

	//上一次成交价格
	price_c int64
	//成交量
	amount int64
	//成交额
	vol int64
	//成交笔数
	count int64
	//主币兑人民币价格
	//cny int64

	//主币兑美元价格
	//usd int64
	//人民币成交额
	usd_vol string

	//买入手续费
	BuyPoundage float64
	//卖出手续费
	SellPoundage float64
}

type TradeInfo struct {
	CreateTime int64
	TradePrice int64
	Num        int64
}

//func NewEntrustQueue(token_id, token_trade_id int, price int64, name string, cny, usd int64, amount, vol, count, usd_vol int64) *EntrustQuene {
func NewEntrustQueue(token_id, token_trade_id int, price int64, name string, amount, vol, count int64,usd_vol string) *EntrustQuene {
	quene_id := name
	log.Infof("load config symbol %s", name)
	m := &EntrustQuene{
		TokenQueueId:      quene_id,
		PriceChannel:      genkey.GetPulishKey(quene_id),
		BuyQueueId:        fmt.Sprintf("%s:1", quene_id),
		SellQueueId:       fmt.Sprintf("%s:2", quene_id),
		MarketBuyQueueId:  fmt.Sprintf("%s:3", quene_id),
		MarketSellQueueId: fmt.Sprintf("%s:4", quene_id),
		HeapBuyHash:       fmt.Sprintf("%s:5", quene_id),
		HeapSellHash:      fmt.Sprintf("%s:6", quene_id),
		HeapBuySort:       fmt.Sprintf("%s:7", quene_id),
		HeapSellSort:      fmt.Sprintf("%s:8", quene_id),
		//InnerBuyQuene:     fmt.Sprintf("%s:9", quene_id),
		//InnerSellQuene:    fmt.Sprintf("%s:10", quene_id),
		TokenId:           token_id,
		TradeQuene:        fmt.Sprintf("%s:trade", quene_id),
		TokenTradeId:      token_trade_id,
		UUID:              1,
		waitOrderDetail:   make(chan *EntrustDetail, 1000),
		marketOrderDetail: make(chan *EntrustDetail, 1000),
		price_c:           price,
		//cny:               cny,
		//usd:               usd,
		usd_vol: usd_vol,
		amount:  amount,
		vol:     vol,
		count:   count,
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
/*
func (s *EntrustQuene) EntrustAl(p *proto.EntrustOrderRequest) (e *EntrustData, ret int32, err error) {
	return
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
	ret, err = m.SubMoneyWithFronzen(session, p.Num, g.EntrustId, proto.TOKEN_TYPE_OPERATOR_HISTORY_ENTRUST)
	if err != nil || ret != ERRCODE_SUCCESS {
		session.Rollback()
		return
	}

	//记录委托
	err = Insert(session, g)
	if err != nil {
		log.Errorln(err.Error())
		session.Rollback()
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
*/
//实时更新交易数据
func (s *EntrustQuene) SetTradeInfo(price int64, deal_num int64) {
	s.price_c = price
	s.count += 1
	s.vol += convert.Int64MulInt64By8Bit(price, deal_num)
	s.amount += deal_num

	s.usd_vol =convert.Int64MulInt64By8BitString(s.vol, GetUsdPrice(int32(s.TokenTradeId)))
	//s.usd_vol += convert.Int64MulInt64By8Bit(s.vol, GetUsdPrice(int32(s.TokenTradeId)))
}

//委托请求检查
func (s *EntrustQuene) EntrustReq(p *proto.EntrustOrderRequest, isFree bool) (ret int32, err error) {
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
	var token_id int
	var sum int64
	var fee_precent float64

	var d proto.TOKEN_TYPE_OPERATOR
	if p.Opt == proto.ENTRUST_OPT_BUY {
		token_id = s.TokenId
		fee_precent = s.BuyPoundage
		if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE {
			sum = convert.Int64DivInt64By8Bit(p.Num, p.OnPrice)
		}
		d = proto.TOKEN_TYPE_OPERATOR_HISTORY_ENTRUST_BUY
	} else {
		token_id = s.TokenTradeId
		fee_precent = s.SellPoundage
		if p.Type == proto.ENTRUST_TYPE_LIMIT_PRICE {
			sum = convert.Int64MulInt64By8Bit(p.Num, p.OnPrice)
		}
		d = proto.TOKEN_TYPE_OPERATOR_HISTORY_ENTRUST_SELL
	}

	g := &EntrustDetail{
		EntrustId:  genkey.GetTimeUnionKey(s.GetUUID()),
		TokenId:    token_id,
		Uid:        p.Uid,
		AllNum:     p.Num,
		SurplusNum: p.Num,
		Opt:        int(p.Opt),
		OnPrice:    p.OnPrice,
		States:     int(proto.TRADE_STATES_TRADE_UN),
		Type:       int(p.Type),
		Symbol:     p.Symbol,
		Sum:        sum,
		FeePercent: fee_precent,
		TradeNum:   0,
		IsFree:     isFree,
	}

	m := &UserToken{}

	err = m.GetUserToken(p.Uid, token_id)
	if err != nil {
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
	ret, err = m.SubMoneyWithFronzen(session, p.Num, g.EntrustId, d)
	if err != nil || ret != ERRCODE_SUCCESS {
		session.Rollback()
		return
	}

	//记录委托
	err = Insert(session, g)
	if err != nil {
		session.Rollback()
		return
	}

	err = session.Commit()
	if err != nil {
		return
	}

	s.waitOrderDetail <- g
	return
}

//开始交易加入举例买入USDT-》BTC  ，卖出USDT-》BTC  ,deal_num 卖方实际消耗BTC数量
func (s *EntrustQuene) MakeDeal(buyer *EntrustDetail, seller *EntrustDetail, price int64, buy_num, deal_num int64) (err error) {
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
	if buyer.Opt != int(proto.ENTRUST_OPT_BUY) {
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

	var main_fee, trade_fee int64

	buyCnyRate := GetCnyPrice(int32(s.TokenTradeId))
	if buyCnyRate == 0 {
		err = errors.New("not found cny price")
		return
	}

	sellCnyRate := GetCnyPrice(int32(s.TokenId))
	if sellCnyRate == 0 {
		err = errors.New("not found cny price")
		return
	}

	if !seller.IsFree {
		main_fee = convert.Int64MulFloat64(buy_num, s.SellPoundage) //卖家消耗手续费0.005个USDT（得到的货币）
	}

	if !buyer.IsFree {
		trade_fee = convert.Int64MulFloat64(deal_num, s.BuyPoundage)
	}

	rand := random.Krand(6, random.KC_RAND_KIND_LOWER)
	no := fmt.Sprintf("%d_%s", time.Now().Unix(), rand)

	log.Infof("trade_fee %d,main_fee %d,deal_num %d,buy_num %d", trade_fee, main_fee, deal_num, buy_num)
	trade_time := time.Now().Unix()
	t := &Trade{
		TradeNo:          no,
		Uid:              buyer.Uid,
		TokenId:          s.TokenId,
		TokenTradeId:     s.TokenTradeId,
		TokenAdmissionId: s.TokenTradeId,
		Price:            price,
		Num:              deal_num - trade_fee,
		Fee:              trade_fee,
		DealTime:         trade_time,
		Opt:              int(proto.ENTRUST_OPT_BUY),
		Symbol:           s.TokenQueueId,
		EntrustId:        buyer.EntrustId,
		FeeCny:           convert.Int64MulInt64By8Bit(trade_fee, buyCnyRate),
		TotalCny:         convert.Int64MulInt64By8Bit(deal_num, buyCnyRate),
	}

	o := &Trade{
		TradeNo:          no,
		Uid:              seller.Uid,
		TokenId:          s.TokenId,
		TokenTradeId:     s.TokenTradeId,
		TokenAdmissionId: s.TokenId,
		Price:            price,
		Num:              buy_num - main_fee,
		Fee:              main_fee,
		DealTime:         trade_time,
		Opt:              int(proto.ENTRUST_OPT_SELL),
		Symbol:           s.TokenQueueId,
		EntrustId:        seller.EntrustId,
		FeeCny:           convert.Int64MulInt64By8Bit(main_fee, sellCnyRate),
		TotalCny:         convert.Int64MulInt64By8Bit(deal_num, sellCnyRate),
	}

	if buyer.SurplusNum < buy_num {
		err = errors.New("error when check surplus num")
		return
	}

	buy_trade, err := GetUserTradeByEntrustId(buyer.EntrustId)
	if err != nil {
		return err
	}

	buy_trade = append(buy_trade, t)
	buyer.Price = CaluateAvgPrice(buy_trade)

	sell_trade, err := GetUserTradeByEntrustId(seller.EntrustId)
	if err != nil {
		return err
	}

	sell_trade = append(sell_trade, o)
	seller.Price = CaluateAvgPrice(sell_trade)

	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}
	//var ret int32 =ERRCODE_SUCCESS
	//USDT left num

	err = buyer.SubSurplus(session, buy_num)
	if err != nil {
		session.Rollback()
		return
	}
	//t.States = buyer.States

	err = seller.SubSurplus(session, deal_num)
	if err != nil {
		session.Rollback()
		return
	}
	//o.States = seller.States
	err = buy_token_account.NotifyDelFronzen(session, buy_num, t.TradeNo, proto.TOKEN_TYPE_OPERATOR_FROZEN_COMFIRM_DEL)
	if err!=nil {
		session.Rollback()
		return
	}


	err = buy_trade_token_account.AddMoney(session, t.Num, t.TradeNo, proto.TOKEN_TYPE_OPERATOR_HISTORY_TRADE)
	if err != nil {
		session.Rollback()
		return
	}

	if buyer.Uid == seller.Uid { //还没处理
		sell_trade_token_account = buy_token_account
		sell_token_account = buy_trade_token_account

	}
	err = sell_token_account.NotifyDelFronzen(session, deal_num, o.TradeNo, proto.TOKEN_TYPE_OPERATOR_FROZEN_COMFIRM_DEL)
	if err != nil {
		session.Rollback()
		return
	}

	err = sell_trade_token_account.AddMoney(session, o.Num, o.TradeNo, proto.TOKEN_TYPE_OPERATOR_HISTORY_TRADE)
	if err != nil {
		session.Rollback()
		return
	}

	err = new(Trade).Insert(session, t, o)
	if err != nil {
		session.Rollback()
		return
	}

	tfree := &TokenFreeHistory{
		TokenId: s.TokenTradeId,
		Opt:     int(proto.TOKEN_OPT_TYPE_ADD),
		Type:    int(proto.TOKEN_TYPE_OPERATOR_HISTORY_TRADE),
		Num:     trade_fee,

		Ukey: t.TradeNo,
		//TradeId:t.TradeId,
	}

	ofree := &TokenFreeHistory{
		TokenId: s.TokenId,
		Opt:     int(proto.TOKEN_OPT_TYPE_ADD),
		Type:    int(proto.TOKEN_TYPE_OPERATOR_HISTORY_TRADE),
		Num:     main_fee,

		Ukey: o.TradeNo,
		//TradeId:o.TradeId,
	}

	err = InsertIntoTokenFreeHistory(session, tfree, ofree)
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

func (s *EntrustQuene) match2(entrust_id string) (err error) {
	//func (s *EntrustQuene) match2(p *EntrustDetail) (err error) {
	p := GetEntrust(entrust_id)
	if p == nil {
		return errors.New(fmt.Sprintf("error entrust_id %d", p.EntrustId))
	}
	var buyer *EntrustDetail
	var seller *EntrustDetail
	var others []*EntrustDetail

	defer func() (err2 error) {
		if err == redis.Nil {
			if buyer != nil && buyer.SurplusNum > 0 {
				log.WithFields(logrus.Fields{
					"buyer_id":         buyer.Uid,
					"buyer_entrust_id": buyer.EntrustId,
					"sulplus":          buyer.SurplusNum,
					"os_id":            os.Getpid(),
				}).Info("NULL quene match")
				s.joinSellQuene(buyer.EntrustId)
			}

			if seller != nil && seller.SurplusNum > 0 {
				log.WithFields(logrus.Fields{
					"seller_id":         seller.Uid,
					"seller_entrust_id": seller.EntrustId,
					"sulplus":           seller.SurplusNum,
					"os_id":             os.Getpid(),
				}).Info("NULL quene match")
				s.joinSellQuene(seller.EntrustId)
			}
		} else if err != nil {
			buyer = GetEntrust(buyer.EntrustId)
			seller = GetEntrust(seller.EntrustId)
			if buyer != nil && buyer.SurplusNum > 0 {
				log.WithFields(logrus.Fields{
					"buy_uid":        buyer.Uid,
					"buy_entrust_id": buyer.EntrustId,
					"sulplus":        buyer.SurplusNum,
					"os_id":          os.Getpid(),
				}).Errorln(err.Error())
				s.joinSellQuene(buyer.EntrustId)
			}

			if seller != nil && seller.SurplusNum > 0 {
				log.WithFields(logrus.Fields{
					"sell_uid":        seller.Uid,
					"sell_entrust_id": seller.EntrustId,
					"sulplus":         seller.SurplusNum,
					"os_id":           os.Getpid(),
				}).Errorln(err.Error())
				s.joinSellQuene(seller.EntrustId)
			}

		} else {
			if buyer != nil && buyer.SurplusNum > 0 {
				err = s.match2(buyer.EntrustId)
				if err == redis.Nil {

				} else if err != nil {
					log.WithFields(logrus.Fields{
						"buy_uid":        buyer.Uid,
						"buy_entrust_id": buyer.EntrustId,
						"sulplus":        buyer.SurplusNum,
						"os_id":          os.Getpid(),
					}).Errorln(err.Error())
				}
			}

			if seller != nil && seller.SurplusNum > 0 {
				err = s.match2(seller.EntrustId)
				if err == redis.Nil {

				} else if err != nil {
					log.WithFields(logrus.Fields{
						"sell_uid":        seller.Uid,
						"sell_entrust_id": seller.EntrustId,
						"sulplus":         seller.SurplusNum,
						"os_id":           os.Getpid(),
					}).Errorln(err.Error())
				}
			}
		}
		return
	}()

	if p.Opt == int(proto.ENTRUST_OPT_BUY) {
		buyer = p
		others, err = s.PopFirstEntrust(proto.ENTRUST_OPT_SELL, 1, 1)
		if err != nil {
			return
		}
		if len(others) > 0 {
			seller = others[0]
		} else {
			err = redis.Nil
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
			err = redis.Nil
			return
		}
	}

	err = s.delSource(proto.ENTRUST_OPT(others[0].Opt), proto.ENTRUST_TYPE(others[0].Type), others[0].EntrustId, others[0])
	if err != nil {
		return
	}
	var buy_num, sell_num, g_num, price int64 //BTC数量，成交价格

	if buyer.Type == int(proto.ENTRUST_TYPE_MARKET_PRICE) {
		if seller.Type == int(proto.ENTRUST_TYPE_MARKET_PRICE) {
			price = s.price_c
			g_num = convert.Int64DivInt64By8Bit(buyer.SurplusNum, price)
		} else {
			price = seller.OnPrice
			g_num = convert.Int64DivInt64By8Bit(buyer.SurplusNum, price)
		}
	} else {
		if seller.Type == int(proto.ENTRUST_TYPE_MARKET_PRICE) {
			price = buyer.OnPrice
			g_num = convert.Int64DivInt64By8Bit(buyer.SurplusNum, price)
		} else {
			//检查价格匹配规则
			if buyer.OnPrice >= seller.OnPrice {

				if buyer.OnPrice >= seller.OnPrice && seller.OnPrice >= s.price_c {
					price = seller.OnPrice
				} else if buyer.OnPrice >= s.price_c && s.price_c >= seller.OnPrice {
					price = s.price_c
				} else if s.price_c >= buyer.OnPrice && buyer.OnPrice >= seller.OnPrice {
					price = buyer.OnPrice
				} else {
					log.WithFields(logrus.Fields{
						"buyer_price":     buyer.OnPrice,
						"seller_price":    seller.OnPrice,
						"price ":          s.price_c,
						"os_id":           os.Getpid(),
						"buy_entrust_id":  buyer.EntrustId,
						"sell_entrust_id": seller.EntrustId,
					}).Errorln("err price please check")
					return nil
				}

				/*
					if buyer.OnPrice <= s.price_c {
						price = buyer.OnPrice
					} else if seller.OnPrice >= s.price_c {
						price = seller.OnPrice
					} else if s.price_c > buyer.OnPrice && s.price_c < seller.OnPrice {
						price = s.price_c
					}else{
						log.WithFields(logrus.Fields{
							"buyer_price":        buyer.OnPrice,
							"seller_price": 	seller.OnPrice,
							"price ":         s.price_c,
							"os_id":           os.Getpid(),
							"buy_entrust_id": buyer.EntrustId,
							"sell_entrust_id": seller.EntrustId,
						}).Errorln("err price please check")
						return
					}

				*/
				g_num = convert.Int64DivInt64By8Bit(buyer.SurplusNum, price)
			} else {
				log.WithFields(logrus.Fields{
					"buyer_price":     buyer.OnPrice,
					"seller_price":    seller.OnPrice,
					"price ":          s.price_c,
					"os_id":           os.Getpid(),
					"buy_entrust_id":  buyer.EntrustId,
					"sell_entrust_id": seller.EntrustId,
				}).Errorln("test check  price please check")
				err = redis.Nil
				return
			}
		}
	}

	//计算交易数量
	if g_num != 0 {
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
	} else {
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
		}).Warn("please check logic")
		err = s.SurplusBack(buyer)
		if err != nil {
			return
		}
		return
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

	if sell_num == 0 {
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
		}).Warn("please check logic 1")
		err = s.SurplusBack(seller)
		if err != nil {
			return
		}
		return
	}

	if buy_num == 0 {
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
		}).Warn("please check logic 2")
		err = s.SurplusBack(buyer)
		if err != nil {
			return
		}
		return
	}



	err = s.MakeDeal(buyer, seller, price, buy_num, sell_num)

	log.WithFields(logrus.Fields{
		"buyer_id":         buyer.Uid,
		"seller_id":        seller.Uid,
		"buyer_entrust_id": buyer.EntrustId,
		"sell_entrust_id":  seller.EntrustId,
		"buy_surplus_num":  buyer.SurplusNum,
		"sell_surplus_num": seller.SurplusNum,
		"os_id":            os.Getpid(),
	}).Info("finish match trade")

	if err != nil {
		return
	}

	s.SetTradeInfo(price, sell_num)

	return
}

//剩余小额退回
func (s *EntrustQuene) SurplusBack(e *EntrustDetail) (err error) {
	u := &UserToken{}
	entry := GetEntrust(e.EntrustId)
	if entry == nil {
		return errors.New("entrust id is not exist")
	}
	err = u.GetUserToken(e.Uid, entry.TokenId)
	if err != nil {
		return
	}
	//num:=e.SurplusNum
	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}
	err = u.ReturnFronzen(session, e.SurplusNum, e.EntrustId, proto.TOKEN_TYPE_OPERATOR_HISTORY_FRONZE_SYS_SURPLUS)
	if err != nil {
		session.Rollback()
		return
	}

	entry.States = int(proto.TRADE_STATES_TRADE_ALL)
	//e.SurplusNum -= e.SurplusNum
	_, err = session.Where("entrust_id=?", e.EntrustId).Decr("surplus_num", e.SurplusNum).Cols("states", "surplus_num").Update(entry)
	if err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		return
	}
/*
	err = s.delSource(proto.ENTRUST_OPT(e.Opt), proto.ENTRUST_TYPE(e.Type), e.EntrustId, e)
	if err != nil {
		log.Error(err.Error())
		return
	}
*/
	e.SurplusNum -= e.SurplusNum
	return
}

//匹配交易
/*
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
*/
//处理请求数据
func (s *EntrustQuene) process() {
	for {
		select {
		case w := <-s.waitOrderDetail:
			s.match2(w.EntrustId)
		case d := <-s.marketOrderDetail:
			go func(data *EntrustDetail) {
				time.Sleep(10 * time.Second)
				s.waitOrderDetail <- d
			}(d)

		}
	}
}

func (s *EntrustQuene) Clock() {
	for {
		t := time.Now()
		if t.Second()%60 == 0 {
			m := &proto.PriceCache{
				Id:          time.Now().Unix(),
				Symbol:      s.TokenQueueId,
				Price:       s.price_c,
				CreatedTime: time.Now().Unix(),
				Amount:      s.amount,
				Count:       s.count,
				Vol:         s.vol,
				UsdVol:      s.usd_vol,
			}

			t := jsonpb.Marshaler{EmitDefaults: true}
			data, err := t.MarshalToString(m)
			if err != nil {
				log.Errorln(err.Error())
				return
			}
			//log.Infof("begin publish symbol %s", s.TokenQueueId)
			err = DB.GetRedisConn().Publish(s.PriceChannel, data).Err()

			if err != nil {
				log.Errorln(err.Error())
				return
			}
		}
		time.Sleep(1 * time.Second)
	}

}

//定时器
func (s *EntrustQuene) Clock2() {
	for {
		c := clock.NewClock()
		t := time.Now()
		diff := 60 - t.Second()
		job := func() {
			m := &proto.PriceCache{
				Id:          time.Now().Unix(),
				Symbol:      s.TokenQueueId,
				Price:       s.price_c,
				CreatedTime: time.Now().Unix(),
				Amount:      s.amount,
				Count:       s.count,
				Vol:         s.vol,
				UsdVol:      s.usd_vol,
			}

			t := jsonpb.Marshaler{EmitDefaults: true}
			data, err := t.MarshalToString(m)
			if err != nil {
				log.Errorln(err.Error())
				return
			}
			log.Infof("begin to publish  trade symbol %s", s.TokenQueueId)
			err = DB.GetRedisConn().Publish(s.PriceChannel, data).Err()

			if err != nil {
				log.Errorln(err.Error())
				return
			}

		}

		d := time.Duration(diff) * time.Second

		inter := time.Duration(1 * time.Second)
		c.AddJobWithInterval(d, job)
		time.Sleep(inter)
		//c.AddJobWithInterval(time.Duration(diff+30)*time.Second,  job)
		log.Infof("circle process send trade symbol %s", s.TokenQueueId)
	}

	/*
		c := clock.NewClock()
		job := func() {
			m := &proto.PriceCache{
				Id:          time.Now().Unix(),
				Symbol:      s.TokenQueueId,
				Price:       s.price_c,
				CreatedTime: time.Now().Unix(),
				Amount:      s.amount,
				Count:       s.count,
				Vol:         s.vol,
				UsdVol:      s.usd_vol,
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

		t:=time.Now()
		diff:=60-t.Second()
		c.AddJobRepeat(time.Duration(diff)*time.Second, 0, job)
	*/
}

//委托入队列
//func (s *EntrustQuene) joinSellQuene(p *EntrustDetail) (ret int, err error) {
func (s *EntrustQuene) joinSellQuene(entrust_id string) (ret int, err error) {

	p := GetEntrust(entrust_id)
	if p == nil {
		err = errors.New(fmt.Sprintf("null entrust_id %s", entrust_id))
		return
	}
	if p.Opt > int(proto.ENTRUST_OPT_EOMAX) {
		ret = ERRCODE_PARAM
		return
	}
	log.WithFields(logrus.Fields{
		"entrust": p.EntrustId,
		"uid":     p.Uid,
		"all_num": p.AllNum,
		"sulplus": p.SurplusNum,
		"symbol":  s.TokenQueueId,
		"os_id":   os.Getpid(),
	}).Info("joinSellQuene record")
	if p.SurplusNum == 0 {
		log.WithFields(logrus.Fields{
			"entrust": p.EntrustId,
			"uid":     p.Uid,
			"all_num": p.AllNum,
			"sulplus": p.SurplusNum,
			"os_id":   os.Getpid(),
		}).Errorf("surplus null join quene")
	}

	var quene_id, hash_id, sort_id string
	var x float64
	if p.Opt == int(proto.ENTRUST_OPT_BUY) {
		hash_id = s.HeapBuyHash
		sort_id = s.HeapBuySort

		if p.Type == int(proto.ENTRUST_TYPE_LIMIT_PRICE) {
			quene_id = s.BuyQueueId
			x = convert.Int64ToFloat64By8Bit(p.OnPrice)
		} else {
			quene_id = s.MarketBuyQueueId
			x = float64(time.Now().Unix())
		}

	} else if p.Opt == int(proto.ENTRUST_OPT_SELL) {
		hash_id = s.HeapSellHash
		sort_id = s.HeapSellSort

		if p.Type == int(proto.ENTRUST_TYPE_LIMIT_PRICE) {
			quene_id = s.SellQueueId
			x = convert.Int64ToFloat64By8Bit(p.OnPrice)
		} else {
			quene_id = s.MarketSellQueueId
			x = float64(time.Now().Unix())
		}
	}

	err = DB.GetRedisConn().ZAdd(quene_id, redis.Z{
		Member: p.EntrustId,
		Score:  x,
	}).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	//同步记录堆积数量
	if p.Type == int(proto.ENTRUST_TYPE_MARKET_PRICE) {
		return
	}

	key := fmt.Sprintf("%d", p.OnPrice)

	log.Infof("add begin record on price  %d ,add hash_id %s,key %s,entrust_id %s,num %d", p.OnPrice,hash_id,key,entrust_id,p.SurplusNum)
	var exist bool
	exist, err = DB.GetRedisConn().HExists(hash_id, key).Result()
	if err != nil {
		log.Infof("add begin err  record on price  %d ,add hash_id %s,key %s,entrust_id %s,num %d,err %s", p.OnPrice,hash_id,key,entrust_id,p.SurplusNum,err.Error())
		log.Errorln(err.Error())
		return
	}

	if !exist {
		log.Infof("not exist add begin record on price  %d ,hash_id %s,key %s,entrust_id %s,num %d", p.OnPrice,hash_id,key,entrust_id,p.SurplusNum)
		err = DB.GetRedisConn().ZAdd(sort_id, redis.Z{
			Member: p.OnPrice,
			Score:  x,
		}).Err()
		if err != nil {
			log.Infof("not exist err add begin record on price  %d ,hash_id %s,key %s,entrust_id %s,num %d", p.OnPrice,hash_id,key,entrust_id,p.SurplusNum)
			log.Errorln(err.Error())
			return
		}
	}
	var val int64
	val ,err = DB.GetRedisConn().HIncrBy(hash_id, key, p.SurplusNum).Result()
	if err != nil {
		log.Errorf("add err hash_id %s begin record on price  %d ,key %s,entrust_id %s,num %d ,add val %d,err  %s", hash_id,p.OnPrice,key,entrust_id,p.SurplusNum,val,err.Error())
		return
	}
	log.Infof("add hash_id %s begin record on price  %d ,key %s,entrust_id %s,num %d ,add val %d", hash_id,p.OnPrice,key,entrust_id,p.SurplusNum,val)
	return
}

func (s *EntrustQuene) Test(opt proto.ENTRUST_OPT, OnPrice int64, num int64) {
	var err error
	var x float64
	var hash_id, sort_id string
	if opt == proto.ENTRUST_OPT_BUY {
		hash_id = s.HeapBuyHash
		sort_id = s.HeapBuySort
		x = convert.Int64ToFloat64By8Bit(OnPrice)
	} else {
		hash_id = s.HeapSellHash
		sort_id = s.HeapSellSort
		x = convert.Int64ToFloat64By8Bit(OnPrice)
	}

	key := fmt.Sprintf("%d", OnPrice)
	_, err = DB.GetRedisConn().HGet(hash_id, key).Int64()
	if err == redis.Nil {
		err = DB.GetRedisConn().ZAdd(sort_id, redis.Z{
			Member: OnPrice,
			Score:  x,
		}).Err()
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	} else if err != nil {
		log.Errorln(err.Error())
		return
	}

	err = DB.GetRedisConn().HIncrBy(hash_id, key, num).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}

}

//弹出数据
func (s *EntrustQuene) delSource(opt proto.ENTRUST_OPT, ty proto.ENTRUST_TYPE, entrust_id string, p *EntrustDetail) (err error) {
	var quene_id, hash_id, sort_id string

	if opt == proto.ENTRUST_OPT_BUY { //买入类型
		if ty == proto.ENTRUST_TYPE_LIMIT_PRICE {
			hash_id = s.HeapBuyHash
			sort_id = s.HeapBuySort
			quene_id = s.BuyQueueId
		} else {
			quene_id = s.MarketBuyQueueId
		}

	} else if opt == proto.ENTRUST_OPT_SELL {
		if ty == proto.ENTRUST_TYPE_LIMIT_PRICE {
			hash_id = s.HeapSellHash
			sort_id = s.HeapSellSort
			quene_id = s.SellQueueId
		} else {
			quene_id = s.MarketSellQueueId
		}
	} else {
		return errors.New("opt param err")
	}

	log.Infof("zrem quene %s,entrust_id %s", quene_id, entrust_id)
	err = DB.GetRedisConn().ZRem(quene_id, entrust_id).Err()
	if err != nil {
		log.Errorln(err)
		return
	}

	if hash_id == "" {
		return
	}
	var val int64
	key := fmt.Sprintf("%d", p.OnPrice)

	val, err = DB.GetRedisConn().HIncrBy(hash_id, key, -p.SurplusNum).Result()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	log.Infof("del hash_id %s ,price %d,quene %s,entrust_id %s,val %d,,num %d", hash_id,p.OnPrice,quene_id, entrust_id, val, p.SurplusNum)
	if val == 0 {
		err = DB.GetRedisConn().HDel(hash_id, key).Err()
		if err != nil {
			log.Errorf(		"del hash_id finish 1 %s ,quene %s,entrust_id %s,val %d,,num %d,err %s", hash_id,quene_id, entrust_id, val, p.SurplusNum,err.Error())
			return
		}

		err = DB.GetRedisConn().ZRem(sort_id, key).Err()
		if err != nil {
			log.Errorf(		"del hash_id finish 2 %s ,quene %s,entrust_id %s,val %d,,num %d,err %s", hash_id,quene_id, entrust_id, val, p.SurplusNum,err.Error())
			return
		}

	}

	return
}

func Testu() {
	val, err := DB.GetRedisConn().HIncrBy("ONT/USDT:6", "238000000", -415966388).Result()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
//WRONGTYPE Operation against a key holding the wrong kind of value
	if val <= 0 {
		err = DB.GetRedisConn().HDel("ONT/USDT:6", "238000000").Err()
		if err != nil {
			log.Errorln(err.Error())
			return
		}
		err = DB.GetRedisConn().ZRem("ONT/USDT:8", "238000000").Err()
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}

}

type TempPrice struct {
	OnPrice    string
	SurplusNum int64
}

func (s *EntrustQuene) PopFirstEntrust2(opt proto.ENTRUST_OPT, sw proto.ENTRUST_TYPE, count int64) (en []*TempPrice, err error) {
	var hash_id, sort_id string

	g := make([]string, 0)
	en = make([]*TempPrice, 0)
	if opt == proto.ENTRUST_OPT_BUY { //买入类型
		hash_id = s.HeapBuyHash
		sort_id = s.HeapBuySort
		g, err = DB.GetRedisConn().ZRevRange(s.HeapBuySort, 0, count).Result()

	} else if opt == proto.ENTRUST_OPT_SELL { //卖出类型
		hash_id = s.HeapSellHash
		sort_id = s.HeapSellSort
		g, err = DB.GetRedisConn().ZRange(s.HeapSellSort, 0, count).Result()
	} else {
		return
	}


		if err!=nil {
			log.WithFields(logrus.Fields{
				"hash_id": hash_id,
				"sort_id": sort_id,
				"symbol":  s.TokenQueueId,
			}).Errorln(err)
			return
		}



	for _, v := range g {
		var i int64
		i, err = DB.GetRedisConn().HGet(hash_id, v).Int64()
		if err != nil {
			if err!=nil {
				log.WithFields(logrus.Fields{
					"hash_id": hash_id,
					"sort_id": sort_id,
					"symbol":  s.TokenQueueId,
				}).Errorln(err)
			}
			return
		}
		if i == 0 {
			err = DB.GetRedisConn().ZRem(sort_id, v).Err()
			if err != nil {
				if err!=nil {
					log.WithFields(logrus.Fields{
						"hash_id": hash_id,
						"sort_id": sort_id,
						"symbol":  s.TokenQueueId,
					}).Errorln(err)
				}
				return
			}
			err = DB.GetRedisConn().HDel(hash_id, v).Err()
			if err != nil {
				if err!=nil {
					log.WithFields(logrus.Fields{
						"hash_id": hash_id,
						"sort_id": sort_id,
						"symbol":  s.TokenQueueId,
					}).Errorln(err)
				}
				return
			}
			continue
		}
		en = append(en, &TempPrice{
			OnPrice:    v,
			SurplusNum: i,
		})
	}

	return
}

//获取队列首位交易单 sw1表示先取市价单再取限价单，2表示直接获取限价单，count获取数量
func (s *EntrustQuene) PopFirstEntrust(opt proto.ENTRUST_OPT, sw proto.ENTRUST_TYPE, count int64) (en []*EntrustDetail, err error) {
	g := make([]string, 0)
	//var z []redis.Z
	var quene_id string
	//var ok bool
	if opt == proto.ENTRUST_OPT_BUY { //买入类型
		if sw == proto.ENTRUST_TYPE_MARKET_PRICE {
			quene_id = s.MarketBuyQueueId
		} else {
			quene_id = s.BuyQueueId
		}

		//z, err = DB.GetRedisConn().ZRevRangeWithScores(quene_id, 0, count).Result()
		g, err = DB.GetRedisConn().ZRevRange(quene_id, 0, count).Result()

	} else if opt == proto.ENTRUST_OPT_SELL { //卖出类型
		if sw == proto.ENTRUST_TYPE_MARKET_PRICE {
			quene_id = s.MarketSellQueueId
		} else {
			quene_id = s.SellQueueId
		}

		//z, err = DB.GetRedisConn().ZRangeWithScores(quene_id, 0, count).Result()
		g, err = DB.GetRedisConn().ZRange(quene_id, 0, count).Result()
	}

	if err != nil {
		log.Errorln(err)
		return
	}

	if len(g) == 0 && sw == proto.ENTRUST_TYPE_MARKET_PRICE {
		return s.PopFirstEntrust(opt, proto.ENTRUST_TYPE_LIMIT_PRICE, count)
	} else if len(g) == 0 && sw == proto.ENTRUST_TYPE_LIMIT_PRICE {
		err = redis.Nil
		return
	}

	en2 := make(map[string]*EntrustDetail, 0)
	err = DB.GetMysqlConn().In("entrust_id", g).Find(&en2)
	if err != nil {
		log.Errorln(err)
		return
	}

	//调整顺序
	en = make([]*EntrustDetail, 0)
	for _, k := range g {
		v, ok := en2[k]
		if ok {
			en = append(en, v)
		}
	}

	if len(g) > 0 && len(en) == 0 {
		log.WithFields(log.Fields{
			"opt":         opt,
			"symbol":      s.TokenQueueId,
			"sw":          sw,
			"entrusdt_id": g[0],
		}).Errorf("data is not consist please check")
		//s.delSource(opt, sw, g[0],)
	}
	/*
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
			g := &EntrustDetail{}
			err = json.Unmarshal(b, g)
			if err != nil {
				log.Errorln(err)
				return
			}
			en = append(en, g)
		}
	*/
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

//撤销委托
func (s *EntrustQuene) DelEntrust(e *EntrustDetail) (err error) {

	defer func() {
		if err != nil {
			log.WithFields(logrus.Fields{
				"buyer_id":         e.EntrustId,
				"seller_id":        e.Uid,
				"buyer_entrust_id": e.SurplusNum,
				"os_id":            os.Getpid(),
			}).Errorf("DelEntrust err %s", err.Error())
		}
	}()
	log.WithFields(logrus.Fields{
		"buyer_id":         e.EntrustId,
		"seller_id":        e.Uid,
		"buyer_entrust_id": e.SurplusNum,
		"os_id":            os.Getpid(),
	}).Info("DelEntrust")

	u := &UserToken{}
	err = u.GetUserToken(e.Uid, e.TokenId)
	if err != nil {
		return err
	}

	sess := DB.GetMysqlConn().NewSession()
	defer sess.Close()
	err = sess.Begin()

	err = u.ReturnFronzen(sess, e.SurplusNum, e.EntrustId, proto.TOKEN_TYPE_OPERATOR_HISTORY_ENTRUST_DEL)
	if err != nil {
		sess.Rollback()
		return err
	}
	e.States = int(proto.TRADE_STATES_TRADE_DEL)
	var aff int64
	aff, err = sess.Where("entrust_id=?", e.EntrustId).Decr("surplus_num", e.SurplusNum).Cols("states", "surplus_num").Update(e)
	if err != nil {
		sess.Rollback()
		return err
	}
	if aff == 0 {
		err = errors.New("version is err")
		return
	}
	err = sess.Commit()
	if err != nil {
		return err
	}

	err = s.delSource(proto.ENTRUST_OPT(e.Opt), proto.ENTRUST_TYPE(e.Type), e.EntrustId, e)
	if err != nil {
		return err
	}

	e.SurplusNum = 0

	return nil
}

func Test9(begin, end int64) {
	//DB.GetMysqlConn().Where("states =1 or states =4").GroupBy("symbol").GroupBy("opt").GroupBy("on_price").Select("sum(surplus_num),symbol").Find()
	log.Infof("load time $d", begin)
	g := 1
	for index := 0; g != 0; index++ {
		log.Infof("load index $d", index)
		res, err := DB.GetMysqlConn().Query("SELECT SUM(`surplus_num`) as num ,`symbol`,`opt` ,`on_price`   FROM `entrust_detail` WHERE  `type`=2  AND `states`  IN(1,4) and created_time>=? and created_time<?  GROUP BY `symbol` ,`opt` ,`on_price` limit ?,?", begin, end, index*1000, 1000*(index+1))
		if err != nil {
			log.Error(err.Error())
			return
		}

		g = len(res)
		if g > 0 {
			for _, v := range res {
				g := convert.BytesToInt64Ascii(v["num"])

				h := convert.BytesToStringAscii(v["symbol"])

				j := convert.BytesToInt64Ascii(v["opt"])

				k := convert.BytesToInt64Ascii(v["on_price"])

				q, _ := GetQueneMgr().GetQueneByUKey(h)
				q.Test(proto.ENTRUST_OPT(j), k, g)
			}
		}
	}

	if end > time.Now().Unix() {
		return
	}

	Test9(begin+86400, end+86400)
}

func Test10() {
	c := new(ConfigQuenes)
	g := c.GetAllQuenes()
	index := 0
	var err error
	res := make([]string, 0)
	flag := true
	for _, v := range g {
		q, _ := GetQueneMgr().GetQueneByUKey(v.Name)

		log.Infof("GetQueneByUKey  proc name %s", v.Name)
		DB.GetRedisConn().Del(q.HeapSellHash)
		DB.GetRedisConn().Del(q.HeapBuyHash)
		DB.GetRedisConn().Del(q.HeapSellSort)
		DB.GetRedisConn().Del(q.HeapBuySort)
		for {
			if !flag {
				index = 0
				flag = true
				break
			}

			res, err = DB.GetRedisConn().ZRange(q.BuyQueueId, int64(index*1000), int64((index+1)*1000)).Result()
			if err == redis.Nil {
				flag = false
			} else if err != nil {
				log.Errorf("mysql err %s", err.Error())
				return
			} else {
				if len(res) == 0 {
					flag = false
				} else {
					ret := make([]EntrustDetail, 0)
					err = DB.GetMysqlConn().In("entrust_id", res).Find(&ret)
					if err != nil {
						log.Errorf("redis err %s", err.Error())
						return
					}
					log.Infof("BuyQueueId begin proc index %d", index)
					for _, n := range ret {
						q.Test(proto.ENTRUST_OPT_BUY, n.OnPrice, n.SurplusNum)
					}
					log.Infof("BuyQueueId end proc index %d", index)
					index++
				}
			}

		}

		for {
			if !flag {
				index = 0
				flag = true
				break
			}

			res, err = DB.GetRedisConn().ZRange(q.SellQueueId, int64(index*1000), int64((index+1)*1000)).Result()
			if err == redis.Nil {
				flag = false

			} else if err != nil {
				log.Errorf("mysql err %s", err.Error())
				return
			} else {
				if len(res) == 0 {
					flag = false
				} else {
					ret := make([]EntrustDetail, 0)
					err = DB.GetMysqlConn().In("entrust_id", res).Find(&ret)
					if err != nil {
						log.Errorf("redis err %s", err.Error())
						return
					}
					log.Infof("SellQueueId begin proc index %d", index)
					for _, n := range ret {
						q.Test(proto.ENTRUST_OPT_SELL, n.OnPrice, n.SurplusNum)
					}
					log.Infof(" SellQueueId end proc index %d", index)
					index++
				}
			}

		}

	}
}
