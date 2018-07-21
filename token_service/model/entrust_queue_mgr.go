package model

import (
	"digicon/common/genkey"
	. "digicon/token_service/log"
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
	if !ok {
		return
	}
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
		Log.Fatalf("insert same quene id is %s", e.TokenQueueId)
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

			p,ok := GetPrice(v.Name)
			if ok {
				e := NewEntrustQueue(v.TokenId, v.TokenTradeId, p.Price, v.Name, cny,p.Amount,p.Vol,p.Count)
				s.AddQuene(e)
			}else{
				e := NewEntrustQueue(v.TokenId, v.TokenTradeId, 100000000, v.Name, cny,0,0,0)
				s.AddQuene(e)
			}

		}

	return true
}
