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
			dataMgr: make(map[string]*EntrustQuene),
		}
	})
	return ins
}

//币币交易管理器
type EntrustQueneMgr struct {
	dataMgr map[string]*EntrustQuene
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

	d := new(QuenesConfig).GetAllQuenes()

	for _, v := range d {
		//quene_id := fmt.Sprintf("%d_%d", v.TokenId, v.TokenTradeId)
		e := NewEntrustQueue(v.TokenId, v.TokenTradeId, v.Price, v.Name)
		s.AddQuene(e)
	}

	return true
}
