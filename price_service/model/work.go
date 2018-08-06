package model

import (
	"digicon/common/convert"
	"digicon/common/genkey"
	. "digicon/price_service/dao"

	proto "digicon/proto/rpc"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"
	"time"
)

type PriceInfo struct {
	Key      string
	PreData  *proto.PriceCache
	MinPrice int64
	UsdPrice int64
}

const (
	//OneMinPrice
	FiveMinPrice = iota
	FivteenMinPrice
	FourHourPrice
	OneDayPrice
	OneWeekPrice
	OneMonthPrice
	MaxPrice
)

type PriceWorkQuene struct {
	TokenId      int32
	Symbol       string
	PriceChannel string
	CnyPrice     int64

	entry *proto.PriceCache
	data  []*PriceInfo
}

func NewPriceWorkQuene(name string, token_id int32, cny int64, d *proto.PriceCache) *PriceWorkQuene {
	var period_key = [MaxPrice]string{
		"5min",
		"15min",
		"4hour",
		"1day",
		"1week",
		"1month",
	}
	m := &PriceWorkQuene{
		Symbol:       name,
		PriceChannel: genkey.GetPulishKey(name),
		data:         make([]*PriceInfo, 0),
		TokenId:      token_id,
		entry:        d,
		CnyPrice:     cny,
	}

	for i := 0; i < MaxPrice; i++ {
		v := &PriceInfo{}
		v.Key = fmt.Sprintf("%s_%s", name, period_key[i])
		m.data = append(m.data, v)
	}
	go m.Publish()

	return m
}

func (s *PriceWorkQuene) GetEntry() *proto.PriceCache {
	return s.entry
}

func (s *PriceWorkQuene) updatePrice2(k *proto.PriceCache) {
	InsertPrice(&Price{
		Id:          k.Id,
		Vol:         k.Vol,
		Amount:      k.Amount,
		Price:       k.Price,
		CreatedTime: k.CreatedTime,
		Count:       k.Count,
		Symbol:      k.Symbol,
		UsdVol:      k.UsdVol,
	})

	s.entry = k
}

func (s *PriceWorkQuene) Publish() {
	pb := DB.GetRedisConn().Subscribe(s.PriceChannel)
	ch := pb.Channel()
	for v := range ch {
		k := &proto.PriceCache{}
		err := jsonpb.UnmarshalString(v.Payload, k)
		if err != nil {
			log.Errorln(err.Error())
			continue
		}

		if k.Price == 0 {
			continue
		}
		s.entry = k

		t := time.Unix(k.Id, 0)
		if t.Second() == 0 {
			s.updatePrice2(k)
			min := t.Minute()
			if min%15 == 0 {
				s.save(FivteenMinPrice, k)
			}

			h := t.Hour()
			if h%4 == 0 {
				s.save(FourHourPrice, k)
			}

			if h == 0 {
				s.save(OneDayPrice, k)
				w := t.Weekday()
				if w == 1 {
					s.save(OneWeekPrice, k)
				}

				if t.Day() == 1 {
					s.save(OneMonthPrice, k)
				}
			}

		}

	}

}

func (s *PriceWorkQuene) save(period int, data *proto.PriceCache) {
	p := s.data[period]
	var h *proto.PeriodPrice
	var close, amount, vol, low, high, open int64
	if p.PreData == nil {
		h = &proto.PeriodPrice{
			Id:     data.Id,
			Open:   0,
			Close:  convert.Int64ToFloat64By8Bit(data.Price),
			Amount: convert.Int64ToFloat64By8Bit(data.Amount),
			Vol:    convert.Int64ToFloat64By8Bit(data.Vol),
			Count:  data.Count,
			Low:    convert.Int64ToFloat64By8Bit(data.Price),
			High:   convert.Int64ToFloat64By8Bit(data.Price),
		}
		open = data.Price
		close = data.Price
		amount = data.Amount
		vol = data.Vol
		low = data.Price
		high = data.Price
	} else {
		high = GetHigh(p.PreData.CreatedTime, data.CreatedTime, data.Symbol)

		low = GetLow(p.PreData.CreatedTime, data.CreatedTime, data.Symbol)
		h = &proto.PeriodPrice{
			Id:     data.Id,
			Open:   convert.Int64ToFloat64By8Bit(p.PreData.Price),
			Close:  convert.Int64ToFloat64By8Bit(data.Price),
			Amount: convert.Int64ToFloat64By8Bit(data.Amount - p.PreData.Amount),
			Vol:    convert.Int64ToFloat64By8Bit(data.Vol - p.PreData.Vol),
			Count:  data.Count - p.PreData.Count,
			Low:    convert.Int64ToFloat64By8Bit(low),
			High:   convert.Int64ToFloat64By8Bit(high),
		}
		open = p.PreData.Price
		close = data.Price
		amount = data.Amount - p.PreData.Amount
		vol = data.Vol - p.PreData.Vol

	}

	t := jsonpb.Marshaler{EmitDefaults: true}
	y, err := t.MarshalToString(h)
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	err = DB.GetRedisConn().LPush(p.Key, y).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	InsertKline(&Kline{
		Symbol: s.Symbol,
		Period: p.Key,
		Id:     h.Id,
		Open:   open,
		Close:  close,
		Amount: amount,
		Vol:    vol,
		Low:    low,
		High:   high,
	})

	p.PreData = data
}
