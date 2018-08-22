package model

import (
	"digicon/common/convert"
	"digicon/common/genkey"
	. "digicon/price_service/dao"
	proto "digicon/proto/rpc"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/jsonpb"
	"github.com/liudng/godump"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
	"encoding/json"
)

const (
	//OneMinPrice

	FivteenMinPrice = iota
	OneHourPrice
	FourHourPrice
	OneDayPrice
	MaxPrice
)

type PriceInfo struct {
	Key     string
	PreData *proto.PriceCache

	UsdPrice int64

	MinPriceKey string
	//MinPriceIdKey string

	MaxPriceKey string
	//MaxPriceIdKey string


	Period    time.Duration
	PeriodSec int64
}

type PriceWorkQuene struct {
	TokenId       int32
	ToekenTradeId int32
	Symbol        string
	PriceChannel  string
	CnyPrice      int64

	entry *proto.PriceCache
	data  []*PriceInfo
}

func NewPriceWorkQuene(name string, token_id, token_trade_id int32, init_price int64, d *proto.PriceCache, ch chan *MsgPricePublish) *PriceWorkQuene {
	var period_key = [MaxPrice]string{
		"15min",
		"1hour",
		"4hour",
		"1day",
	}

	var period_sec = [MaxPrice]time.Duration{
		900 * time.Second,
		3600 * time.Second,
		14400 * time.Second,
		86400 * time.Second,
	}
	m := &PriceWorkQuene{
		Symbol:        name,
		PriceChannel:  genkey.GetPulishKey(name),
		data:          make([]*PriceInfo, 0),
		TokenId:       token_id,
		ToekenTradeId: token_trade_id,
		entry:         d,
	}

	var ok bool
	var err error
	ids := make([]int64, 0)
	one := &Price{}
	ok, err = DB.GetMysqlConn().Where("symbol=? ", name).Desc("id").Limit(1, 0).Get(one)
	if err != nil {
		log.Fatal(err.Error())
	}

	if !ok {
		for i := 0; i < MaxPrice; i++ {
			v := &PriceInfo{}
			v.Key = fmt.Sprintf("%s_%s", name, period_key[i])
			v.PreData = &proto.PriceCache{
				Id:          0,
				Symbol:      name,
				Price:       init_price,
				CreatedTime: 0,
				Amount:      0,
				Vol:         0,
				Count:       0,
				UsdVol:      0,
			}
			v.MinPriceKey = fmt.Sprintf("price:%s:min", v.Key, period_key[i])
		//	v.MinPriceIdKey = fmt.Sprintf("price:%s:min_id", v.Key, period_key[i])
			v.MaxPriceKey = fmt.Sprintf("price:%s:max", v.Key, period_key[i])
		//	v.MaxPriceIdKey = fmt.Sprintf("price:%s:max_id", v.Key, period_key[i])
			v.Period = period_sec[i]
			v.PeriodSec = int64(v.Period / time.Second)
			m.data = append(m.data, v)
		}
		go m.Subscribe(ch)
		return m
	}

	t := time.Unix(one.Id, 0)

	var min15, hour4, hour1 int

	min15 = t.Minute() % 15

	fifteen := t.Add(-time.Duration(60*min15) * time.Second)

	ids = append(ids, fifteen.Unix())

	hour1 = t.Minute() % 60

	onehour := t.Add(-time.Duration(hour1*60) * time.Second)
	ids = append(ids, onehour.Unix())

	hour4 = t.Hour() % 4

	four := onehour.Add(-time.Duration(3600*hour4) * time.Second)
	ids = append(ids, four.Unix())

	left_hour := t.Hour() % 24
	left_min := t.Minute() % 60
	oneday := t.Add(-time.Duration(left_hour*3600)*time.Second - time.Duration(left_min*60)*time.Second)
	ids = append(ids, oneday.Unix())

	for i := 0; i < MaxPrice; i++ {

		v := &PriceInfo{}
		v.Key = fmt.Sprintf("%s_%s", name, period_key[i])

		p := &Price{}
		ok, err := DB.GetMysqlConn().Where("symbol=? and id=?", name, ids[i]).Get(p)
		if err != nil {
			log.Fatal(err.Error())
		}
		if !ok {
			//log.Fatal("please check price id is %d,symbol is %s  base data is null",ids[i],)
			//init

			v.PreData = &proto.PriceCache{
				Id:          ids[i],
				Symbol:      name,
				Price:       init_price,
				CreatedTime: ids[i],
				Amount:      0,
				Vol:         0,
				Count:       0,
				UsdVol:      0,
			}

		} else {
			v.PreData = &proto.PriceCache{
				Id:          p.Id,
				Symbol:      p.Symbol,
				Price:       p.Price,
				CreatedTime: p.CreatedTime,
				Amount:      p.Amount,
				Vol:         p.Vol,
				Count:       p.Count,
				UsdVol:      p.UsdVol,
			}
		}
		v.MinPriceKey = fmt.Sprintf("price:%s:%s:min", name, period_key[i])
	//	v.MinPriceIdKey = fmt.Sprintf("price:%s:min_id", v.Key, period_key[i])
		v.MaxPriceKey = fmt.Sprintf("price:%s:%s:max", name, period_key[i])
	//	v.MaxPriceIdKey = fmt.Sprintf("price:%s:max_id", v.Key, period_key[i])
		v.Period = period_sec[i]
		m.data = append(m.data, v)
	}
	go m.Subscribe(ch)

	ch <- &MsgPricePublish{
		Symbol:       m.Symbol,
		TokenId:      m.TokenId,
		TokenTradeId: m.ToekenTradeId,
		Price:        m.entry.Price,
	}
	return m
}

func (s *PriceWorkQuene) GetEntry() *proto.PriceCache {
	return s.entry
}

func (s *PriceWorkQuene) updatePrice2(k *proto.PriceCache) (err error) {
	g,ok:=GetQueneMgr().PriceMap[s.ToekenTradeId]
	if ok {
		err = InsertPrice(&Price{
			Id:          k.Id,
			Vol:         k.Vol,
			Amount:      k.Amount,
			Price:       k.Price,
			CreatedTime: k.CreatedTime,
			Count:       k.Count,
			Symbol:      k.Symbol,
			UsdVol:      k.UsdVol,
			CnyPrice:    g.CnyPrice,
		})
	}else {
		log.Error("can not find cny price")
		err = InsertPrice(&Price{
			Id:          k.Id,
			Vol:         k.Vol,
			Amount:      k.Amount,
			Price:       k.Price,
			CreatedTime: k.CreatedTime,
			Count:       k.Count,
			Symbol:      k.Symbol,
			UsdVol:      k.UsdVol,
		})
	}

	return
}

func (s *PriceWorkQuene) ReloadKline(stime int64) {

	godump.Dump(stime)
	return
	k := make([]*Price, 0)
	lastt := stime + 43200
	err := DB.GetMysqlConn().Where("created_time>=? and created_time<? and symbol=?", stime, lastt, s.Symbol).Find(&k)
	if err != nil {
		log.Fatalf(err.Error())
	}
	for _, v := range k {

		t := time.Unix(v.Id, 0)
		c := &proto.PriceCache{
			Id:          v.Id,
			Symbol:      v.Symbol,
			Price:       v.Price,
			Amount:      v.Amount,
			Vol:         v.Vol,
			Count:       v.Count,
			UsdVol:      v.UsdVol,
			CreatedTime: v.CreatedTime,
		}

		if t.Second() == 0 {
			s.entry = c
			min := t.Minute()
			if min%15 == 0 {
				s.save(FivteenMinPrice, c)
				if min == 0 {
					s.save(OneHourPrice, c)
					h := t.Hour()
					if h%4 == 0 {
						s.save(FourHourPrice, c)
						if h == 0 {
							s.save(OneDayPrice, c)
						}
					}
				}
			}
		}
	}

	if time.Now().Unix() <= lastt {
		return
	} else {
		s.ReloadKline(lastt)
	}

}

func (s *PriceWorkQuene) Subscribe(mpp chan *MsgPricePublish) {
	pb := DB.GetRedisConn().Subscribe(s.PriceChannel)
	ch := pb.Channel()
	for v := range ch {
		c := &proto.PriceCache{}
		err := jsonpb.UnmarshalString(v.Payload, c)
		if err != nil {
			log.Errorln(err.Error())
			continue
		}

		if c.Price == 0 {
			continue
		}
		log.Infof("rev price %d symbol %s", c.Price, c.Symbol)
		t := time.Unix(c.Id, 0)
		if t.Second() == 0 {
			err = s.updatePrice2(c)
			if err != nil {
				continue
			}
			s.SetPrice(c)
			mpp <- &MsgPricePublish{
				Symbol:       s.Symbol,
				TokenId:      s.TokenId,
				TokenTradeId: s.ToekenTradeId,
				Price:        s.entry.Price,
			}
			min := t.Minute()
			if min%15 == 0 {
				s.save(FivteenMinPrice, c)
				if min == 0 {
					s.save(OneHourPrice, c)
					h := t.Hour()
					if h%4 == 0 {
						s.save(FourHourPrice, c)
						if h == 0 {
							s.save(OneDayPrice, c)
						}
					}
				}
			}
		}
	}

}

func (s *PriceWorkQuene) save(period int, data *proto.PriceCache) {
	log.WithFields(log.Fields{
		"Symbol": s.Symbol,
		"Key":    data.Id,
		"id":     data.Id,
		"amount": data.Amount,
		"vol":    data.Vol,
	}).Info("begin price record ")

	p := s.data[period]
	var h *proto.PeriodPrice
	var close, amount, vol,  open, count int64
	var err error
	/*
		if p.PreData.Count == 0 {
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
			count = data.Count
		} else {
	*/

	high, low := s.GetPeriodMaxPrice(period)
	if high==nil || low==nil {
		log.Error("err when save high or low")
		return
	}
	/*
		high, err = DB.GetRedisConn().Get(p.MaxPriceKey).Int64()
		if err == redis.Nil {
			log.WithFields(log.Fields{
				"symbol":           s.Symbol,
				"data.price ":      data.Price,
				"os_id":            os.Getpid(),
				"data.id":          data.Id,
				"max_key":          p.MaxPriceKey,
				"p.PreData.Dara":   p.PreData.CreatedTime,
				"data.CreatedTime": data.CreatedTime,
			}).Info("record price can not found")
			high = GetHigh(p.PreData.CreatedTime, data.CreatedTime, data.Symbol)
			if high == 0 {
				high = data.Price
			}
		} else if err != nil {
			log.Error(err)
			return
		}

		low, err = DB.GetRedisConn().Get(p.MinPriceKey).Int64()
		if err == redis.Nil {
			log.WithFields(log.Fields{
				"symbol":           s.Symbol,
				"data.price ":      data.Price,
				"os_id":            os.Getpid(),
				"data.id":          data.Id,
				"min_key":          p.MinPriceKey,
				"p.PreData.Dara":   p.PreData.CreatedTime,
				"data.CreatedTime": data.CreatedTime,
			}).Info("record price can not found")
			low = GetLow(p.PreData.CreatedTime, data.CreatedTime, data.Symbol)
			if low == 0 {
				low = data.Price
			}
		} else if err != nil {
			log.Error(err)
			return
		}
	*/
	h = &proto.PeriodPrice{
		Id:     data.Id,
		Open:   convert.Int64ToFloat64By8Bit(p.PreData.Price),
		Close:  convert.Int64ToFloat64By8Bit(data.Price),
		Amount: convert.Int64ToFloat64By8Bit(data.Amount - p.PreData.Amount),
		Vol:    convert.Int64ToFloat64By8Bit(data.Vol - p.PreData.Vol),
		Count:  data.Count - p.PreData.Count,
		Low:    convert.Int64ToFloat64By8Bit(low.Price),
		High:   convert.Int64ToFloat64By8Bit(high.Price),
	}
	open = p.PreData.Price
	close = data.Price
	amount = data.Amount - p.PreData.Amount
	vol = data.Vol - p.PreData.Vol
	count = data.Count - p.PreData.Count

	t := jsonpb.Marshaler{EmitDefaults: true}
	y, err := t.MarshalToString(h)
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	if amount < 0 {
		log.WithFields(log.Fields{
			"Symbol":    s.Symbol,
			"Key":       p.Key,
			"id":        data.Id,
			"amount":    amount,
			"vol":       vol,
			"low":       low,
			"p.PreData": p.PreData.Id,
		}).Errorf("price record error ")
		return
	}

	err = InsertKline(&Kline{
		Symbol: s.Symbol,
		Period: p.Key,
		Id:     h.Id,
		Open:   open,
		Close:  close,
		Amount: amount,
		Vol:    vol,
		Low:    low.Price,
		High:   high.Price,
		Count:  count,
	})

	if err != nil {
		return
	}
	err = DB.GetRedisConn().LPush(p.Key, y).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	p.PreData = data
}


func (s *PriceWorkQuene) SetPrice(data *proto.PriceCache) {
	var err error
	s.entry = data
	b:=make([]byte,0)
	dmin:=new(Price)
	dmax:=new(Price)
	for _, v := range s.data {
		b, err = DB.GetRedisConn().Get(v.MinPriceKey).Bytes()
		if err == redis.Nil {
			dmin = GetLow(v.PreData.CreatedTime, data.CreatedTime, s.Symbol)
			if dmin==nil {
				break
			}
			log.WithFields(log.Fields{
				"min":         dmin.Price ,
				"symbol":       s.Symbol,
				"data.price ":  data.Price,
				"os_id":        os.Getpid(),
				"data.id":      data.Id,
				"PreData.Dara": v.PreData.CreatedTime,
				"CreatedTime":  data.CreatedTime,
				"MinPriceKey":  v.MinPriceKey,
			}).Info("record price case redis nil")

			b,err =json.Marshal(dmin)
			if err!=nil {
				return
			}
			err = DB.GetRedisConn().Set(v.MinPriceKey, b, v.Period).Err()
			if err != nil {
				log.Error(err)
				return
			}
		} else if err != nil {
			log.Error(err)
			return
		}

		err = json.Unmarshal(b,dmin)
		if err!=nil {
			log.Error(err)
			return
		}

		log.WithFields(log.Fields{
			"min":         dmin.Price,
			"symbol":      s.Symbol,
			"data.price ": data.Price,
			"os_id":       os.Getpid(),
			"data.id":     data.Id,
		}).Info("record price")

		if dmin.Price >= data.Price {
			log.WithFields(log.Fields{
				"min":         dmin.Price,
				"symbol":      s.Symbol,
				"data.price ": data.Price,
				"os_id":       os.Getpid(),
				"data.id":     data.Id,
			}).Info("price is low change price")

			p:=&Price{
				Id:data.Id,
				Symbol:data.Symbol,
				Price:data.Price,
				CreatedTime:data.CreatedTime,
				Amount:data.Amount,
				Vol:data.Vol,
			}

			c,ok:=GetQueneMgr().PriceMap[s.ToekenTradeId]
			if ok{
				p.CnyPrice=c.CnyPrice
			}

			b,err =json.Marshal(p)
			if err != nil {
				log.Error(err)
				return
			}

			err = DB.GetRedisConn().Set(v.MinPriceKey, b, v.Period).Err()
			if err != nil {
				log.Error(err)
				return
			}
		}

		b, err = DB.GetRedisConn().Get(v.MaxPriceKey).Bytes()
		if err == redis.Nil {
			h := GetHigh(v.PreData.CreatedTime, data.CreatedTime, s.Symbol)
			log.WithFields(log.Fields{
				"max":              h.Price,
				"symbol":           s.Symbol,
				"data.price ":      data.Price,
				"os_id":            os.Getpid(),
				"data.id":          data.Id,
				"p.PreData.Dara":   v.PreData.CreatedTime,
				"data.CreatedTime": data.CreatedTime,
				"MaxPriceKey":      v.MaxPriceKey,
			}).Info("record price case redis nil")


			b,err =json.Marshal(dmax)
			if err!=nil {
				return
			}
			err = DB.GetRedisConn().Set(v.MaxPriceKey, b, v.Period).Err()
			if err != nil {
				log.Error(err)
				return
			}
		} else if err != nil {
			log.Error(err)
			return
		}


		err = json.Unmarshal(b,dmax)
		if err!=nil {
			log.Error(err)
			return
		}

		if dmax.Price <= data.Price {
			log.WithFields(log.Fields{
				"max":         dmax.Price,
				"symbol":      s.Symbol,
				"data.price ": data.Price,
				"os_id":       os.Getpid(),
				"data.id":     data.Id,
			}).Info("price is high change price")

			p:=&Price{
				Id:data.Id,
				Symbol:data.Symbol,
				Price:data.Price,
				CreatedTime:data.CreatedTime,
				Amount:data.Amount,
				Vol:data.Vol,
			}
			c,ok:=GetQueneMgr().PriceMap[s.ToekenTradeId]
			if ok{
				p.CnyPrice=c.CnyPrice
			}

			b,err =json.Marshal(p)
			if err != nil {
				log.Error(err)
				return
			}

			err = DB.GetRedisConn().Set(v.MaxPriceKey,b, v.Period).Err()
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
}

func (s *PriceWorkQuene) GetPeriodMaxPrice(period int, begin ...int64) (pmax, pmin *Price) {
	t := s.data[period]
	b, err := DB.GetRedisConn().Get(t.MinPriceKey).Bytes()
	var ok bool
	if err == redis.Nil {
		if len(begin) > 0 {
			pmin = GetLow(begin[0]-t.PeriodSec, begin[0], s.Symbol)
			if pmin==nil {
				return
			}
			b,err :=json.Marshal(pmin)
			if err != nil {
				log.Error(err)
				return
			}

			err = DB.GetRedisConn().Set(t.MinPriceKey,b, t.Period).Err()
			if err != nil {
				log.Error(err)
				return
			}

		} else {
				pmin=&Price{}
				ok,err=DB.GetMysqlConn().Where("id=? and symbol=?",t.PreData.Id,s.Symbol).Get(pmin)
				if err!=nil {
					log.Error(err.Error())
					return
				}
			if !ok {
				pmin.Id=s.entry.Id
				pmin.Price=s.entry.Price
				c,ok:=GetQueneMgr().PriceMap[s.ToekenTradeId]
				if ok {
					pmin.CnyPrice=c.CnyPrice
				}

			}
		}
	} else if err != nil {
		log.Error(err)
		return
	}else{
		pmin=new(Price)
		err = json.Unmarshal(b,pmin)
		if err!=nil {
			log.Error(err)
			return
		}
	}

	b, err = DB.GetRedisConn().Get(t.MaxPriceKey).Bytes()
	if err == redis.Nil {
		if len(begin) > 0 {
			pmax = GetHigh(begin[0]-t.PeriodSec, begin[0], s.Symbol)
			if pmax==nil {
				return
			}
			b,err :=json.Marshal(pmax)
			if err != nil {
				log.Error(err)
				return
			}

			err = DB.GetRedisConn().Set(t.MaxPriceKey,b, t.Period).Err()
			if err != nil {
				log.Error(err)
				return
			}

		} else {
			pmax=&Price{}
			ok,err=DB.GetMysqlConn().Where("id=? and symbol=?",t.PreData.Id,s.Symbol).Get(pmax)
			if err!=nil {
				log.Error(err.Error())
				return
			}

			if !ok {
				pmax.Id=s.entry.Id
				pmax.Price=s.entry.Price
				c,ok:=GetQueneMgr().PriceMap[s.ToekenTradeId]
				if ok {
					pmax.CnyPrice=c.CnyPrice
				}
			}
			return
		}
	} else if err != nil {
		log.Error(err)
		return
	}else{
		pmax=new(Price)
		err = json.Unmarshal(b,pmax)
		if err!=nil {
			log.Error(err)
			return
		}
	}
	return
}
