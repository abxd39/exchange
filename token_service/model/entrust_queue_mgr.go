package model

import (
	. "digicon/user_service/dao"
	. "digicon/user_service/log"
	"fmt"
	"sync"
)

var ins *EntrustQueneMgr
var once sync.Once

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
func (s *EntrustQueneMgr) GetQuene(id string) (d *EntrustQuene, ok bool) {
	d, ok = s.dataMgr[id]
	if !ok {
		return
	}
	return
}

func (s *EntrustQueneMgr) AddQuene(e *EntrustQuene) bool {
	_, ok := s.dataMgr[e.TokenQueneId]
	if !ok {
		Log.Fatalf("insert same quene id is %s", e.TokenQueneId)
	}
	s.dataMgr[e.TokenQueneId] = e
	return ok
}

func (s *EntrustQueneMgr) CallBackFunc(f func(*EntrustQuene)) {
	for _, v := range s.dataMgr {
		f(v)
	}
}

//初始化配置
func (s *EntrustQueneMgr) Init() bool {
	d := make([]QuenesConfig, 0)
	err := DB.GetMysqlConn().Find(&d)
	if err != nil {
		Log.Fatalln(err.Error())
	}

	for _, v := range d {
		quene_id := fmt.Sprintf("%s_%s", v.TokenId, v.TokenTradeId)
		e := NewEntrustQuene(quene_id)
		s.AddQuene(e)
	}

	return true
}
