package model

import (
	. "digicon/price_service/log"

	"sync"
)

var ins *WorkQueneMgr
var once sync.Once

//单例获取币币交易管理器
func GetQueneMgr() *WorkQueneMgr {
	once.Do(func() {
		ins = &WorkQueneMgr{
			dataMgr: make(map[string]*PriceWorkQuene),
		}
	})
	return ins
}

type WorkQueneMgr struct {
	dataMgr map[string]*PriceWorkQuene
}

func (s *WorkQueneMgr) Init() {
	d := new(ConfigQuenes).GetAllQuenes()
	for _, v := range d {
		g := NewPriceWorkQuene(v.Name,int32(v.TokenId))
		s.AddQuene(g)
	}

}

//添加一个币币交易
func (s *WorkQueneMgr) AddQuene(e *PriceWorkQuene) bool {
	_, ok := s.dataMgr[e.PriceChannel]
	if ok {
		Log.Fatalf("insert same quene id is %s", e.PriceChannel)
	}
	s.dataMgr[e.PriceChannel] = e
	return ok
}

//获取一个币币交易
func (s *WorkQueneMgr) GetQueneByUKey(ukey string) (d *PriceWorkQuene, ok bool) {
	d, ok = s.dataMgr[ukey]
	if !ok {
		return
	}
	return
}
