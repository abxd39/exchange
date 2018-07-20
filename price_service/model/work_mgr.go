package model

import (
	"sync"
	. "digicon/price_service/log"
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
	for _, v := range ConfigQuenes {
		p,ok:=GetPrice(v.Name)
		if !ok {
			//Log.Fatalln("please init price")
			g := NewPriceWorkQuene(v.Name, int32(v.TokenId),nil)
			s.AddQuene(g)
		}else {
			g := NewPriceWorkQuene(v.Name, int32(v.TokenId),p.SetProtoData())
			s.AddQuene(g)
		}

	}
}

//添加一个币币交易
func (s *WorkQueneMgr) AddQuene(e *PriceWorkQuene) bool {
	_, ok := s.dataMgr[e.Symbol]
	if ok {
		Log.Fatalf("insert same quene id is %s", e.Symbol)
	}
	s.dataMgr[e.Symbol] = e
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
