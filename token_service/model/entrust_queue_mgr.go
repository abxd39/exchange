package model

import (
	"digicon/common/genkey"
	log "github.com/sirupsen/logrus"
	"sync"
)

var ins *EntrustQueneMgr
var once sync.Once

//单例获取币币交易管理器
func GetQueneMgr() *EntrustQueneMgr {
	once.Do(func() {
		ins = &EntrustQueneMgr{
			dataMgr:    make(map[string]*EntrustQuene),
			readyQuene: make(chan ConfigQuenes, 1000),
		}
	})
	return ins
}

//币币交易管理器
type EntrustQueneMgr struct {
	dataMgr    map[string]*EntrustQuene
	readyQuene chan ConfigQuenes
}

//获取一个币币交易
func (s *EntrustQueneMgr) GetQueneByUKey(ukey string) (d *EntrustQuene, ok bool) {
	d, ok = s.dataMgr[ukey]
	return
}

//获取一个币币交易
func (s *EntrustQueneMgr) GetQuene(token_id, trade_token_id int) (d *EntrustQuene, ok bool) {
	return s.GetQueneByUKey(genkey.GetUnionKey(token_id, trade_token_id))
}

//添加一个币币交易
func (s *EntrustQueneMgr) AddQuene(e *EntrustQuene) bool {
	_, ok := s.dataMgr[e.TokenQueueId]
	if ok {
		log.Fatalf("insert same quene id is %s", e.TokenQueueId)
	}
	s.dataMgr[e.TokenQueueId] = e
	return ok
}

//遍历每个币币交易
func (s *EntrustQueneMgr) CallBackFunc(f func(*EntrustQuene)) {
	for _, v := range s.dataMgr {
		f(v)
	}
}

//初始化配置
func (s *EntrustQueneMgr) Init() bool {
	InitConfigTokenCny()

	d := new(ConfigQuenes).GetAllQuenes()

	for _, v := range d {
		cny := GetTokenCnyPrice(v.TokenId)
		usd := GetTokenUsdPrice(v.TokenId)
		if cny == 0 {
			panic("err cny config")
		}
		if usd == 0 {
			panic("err usd config")
		}

		p, ok := GetPrice(v.Name)

		if v.SellPoundage==0 {
			log.Fatalf("err SellPoundage config  symbol %s",v.Name)
		}
		if v.BuyPoundage==0 {
			log.Fatalf("err BuyPoundage config  symbol %s",v.Name)
		}
		if ok {
			e := NewEntrustQueue(v.TokenId, v.TokenTradeId, p.Price, v.Name, cny, usd, p.Amount, p.Vol, p.Count, p.UsdVol)
			e.SellPoundage = v.SellPoundage
			e.BuyPoundage = v.BuyPoundage
			s.AddQuene(e)
		} else {
			e := NewEntrustQueue(v.TokenId, v.TokenTradeId, v.Price, v.Name, cny, usd, 0, 0, 0, 0)
			e.SellPoundage = v.SellPoundage
			e.BuyPoundage = v.BuyPoundage
			s.AddQuene(e)
		}
	}

	return true
}
