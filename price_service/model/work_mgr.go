package model

import (
	"digicon/common/convert"
	proto "digicon/proto/rpc"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var ins *WorkQueneMgr
var once sync.Once

//单例获取币币交易管理器
func GetQueneMgr() *WorkQueneMgr {
	once.Do(func() {
		ins = &WorkQueneMgr{
			dataMgr:  make(map[string]*PriceWorkQuene),
			PriceMap: make(map[int32]*PublicPrice),
			msgChan:  make(chan *MsgPricePublish, 1000),
		}
	})
	return ins
}

type MsgPricePublish struct {
	Symbol       string
	TokenId      int32
	TokenTradeId int32
	Price        int64
}

type PublicPrice struct {
	Data       *MsgPricePublish
	CnyPrice   int64
	UsdPrice   int64
	UpdateTime int64
}

type WorkQueneMgr struct {
	dataMgr map[string]*PriceWorkQuene

	PriceMap map[int32]*PublicPrice

	msgChan chan *MsgPricePublish

	UsdRate int64
	CnyRate int64

	IsRun bool
}

func (s *WorkQueneMgr) Init() {
	InitConfig()
	InitConfigTitle()
	InitConfigTokenCny()
	b := GetTokenCnyPrice(1)
	s.UsdRate = b.UsdPrice
	s.CnyRate = b.Price

	s.PriceMap[1] = &PublicPrice{
		CnyPrice:   s.CnyRate,
		UsdPrice:   s.UsdRate,
		UpdateTime: time.Now().Unix(),
		Data: &MsgPricePublish{
			TokenId:      1,
			TokenTradeId: 1,
			Price:        s.CnyRate,
		},
	}

	for _, v := range ConfigQueneData {
		c, ok := ConfigQueneInit[v.Name]
		if !ok {
			log.Fatalf("err init quene price")
		}

		p, ok := GetPrice(v.Name)
		if !ok {
			g := NewPriceWorkQuene(v.Name, v.TokenId, v.TokenTradeId, c.Price, &proto.PriceCache{
				Id:          time.Now().Unix(),
				Symbol:      v.Name,
				Price:       c.Price,
				CreatedTime: time.Now().Unix(),
			}, s.msgChan)
			s.AddQuene(g)
		} else {
			g := NewPriceWorkQuene(v.Name, v.TokenId, v.TokenTradeId, c.Price, p.SetProtoData(), s.msgChan)
			s.AddQuene(g)
		}
	}
	go s.Process()
}

func (s *WorkQueneMgr) Process() {
	for v := range s.msgChan {
		if v.TokenId == 1 { //USDT
			d, ok := s.PriceMap[v.TokenTradeId]
			if !ok {
				s.PriceMap[v.TokenTradeId] = &PublicPrice{
					CnyPrice:   convert.Int64MulInt64By8Bit(v.Price, s.CnyRate),
					UsdPrice:   convert.Int64MulInt64By8Bit(v.Price, s.UsdRate),
					UpdateTime: time.Now().Unix(),
					Data:       v,
				}
			} else {
				d.Data = v
				d.CnyPrice = convert.Int64MulInt64By8Bit(v.Price, s.CnyRate)
				d.UsdPrice = convert.Int64MulInt64By8Bit(v.Price, s.UsdRate)
				d.UpdateTime = time.Now().Unix()
			}
		} else {
			g, ok := s.PriceMap[v.TokenId]
			if !ok {
				log.Errorf("get err price token_id %d", v.TokenId)
				s.msgChan <- v
				continue
			}

			d, ok := s.PriceMap[v.TokenTradeId]
			if !ok {
				s.PriceMap[v.TokenTradeId] = &PublicPrice{
					CnyPrice:   convert.Int64MulInt64MulInt64By16Bit(v.Price, g.Data.Price, s.CnyRate),
					UsdPrice:   convert.Int64MulInt64MulInt64By16Bit(v.Price, g.Data.Price, s.UsdRate),
					UpdateTime: time.Now().Unix(),
					Data:       v,
				}
			} else {
				if d.Data.Symbol != v.Symbol { //遇到同一个币的不同队列只保留一个
					continue
				}
				d.CnyPrice = convert.Int64MulInt64MulInt64By16Bit(v.Price, g.Data.Price, s.CnyRate)
				d.UsdPrice = convert.Int64MulInt64MulInt64By16Bit(v.Price, g.Data.Price, s.UsdRate)
				d.UpdateTime = time.Now().Unix()
				d.Data = v
			}
		}
	}

}

//添加一个币币交易
func (s *WorkQueneMgr) AddQuene(e *PriceWorkQuene) bool {
	_, ok := s.dataMgr[e.Symbol]
	if ok {
		log.Fatalf("insert same quene id is %s", e.Symbol)
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
