package model

import (
	"digicon/common/genkey"
	. "digicon/price_service/dao"
	"github.com/golang/protobuf/jsonpb"
	proto "digicon/proto/rpc"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
	"errors"
	"digicon/common/convert"
)


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

type PriceInfo struct {
	Key      string
	PreData  *proto.PriceCache
	MinPrice int64
	UsdPrice int64
}

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

	c:=&ConfigQuenes{}
	ok,err := DB.GetMysqlConn2().Where("switch=1").Get(c)
	if err!=nil {
		log.Fatal(err.Error())
	}
	if !ok {
		log.Fatal(errors.New("err config"))
	}
	ids:=make([]int64,0)
	one:=&Price{}
	ok,err =DB.GetMysqlConn().Where("symbol=? ",name).Desc("id").Get(one)
	if err!=nil {
		log.Fatal(err.Error())
	}



	if !ok {
		for i := 0; i < MaxPrice; i++ {
			v := &PriceInfo{}
			v.Key = fmt.Sprintf("%s_%s", name, period_key[i])
			v.PreData=&proto.PriceCache{
				Id:0,
				Symbol:name,
				Price:c.Price,
				CreatedTime:0,
				Amount:0,
				Vol:0,
				Count:0,
				UsdVol:0,
			}
			m.data = append(m.data, v)
		}
		go m.Subscribe()
		return m
	}
	t := time.Unix(one.Id, 0)

	var min5 ,min15 ,hour4 int
	if t.Minute()%5 !=0 {
		min5=t.Minute()%5
	}

	five:=t.Add(-time.Duration(60*min5)*time.Second)
	log.Debug("five add %v",&five)

	ids=append(ids,five.Unix())

	if t.Minute()%15 !=0 {
		min15= t.Minute()%15
	}

	fifteen:=t.Add(-time.Duration(60*min15)*time.Second)

	log.Debug("fifteen add %v",&fifteen)
	ids=append(ids,fifteen.Unix())

	if t.Hour()%4 !=0 {
		hour4 = t.Hour()%4
	}

	four:=t.Add(-time.Duration(t.Minute()*60)*time.Second -time.Duration(3600*hour4)*time.Second)
	ids=append(ids,four.Unix())


	left_hour:=t.Hour()%24
	left_min:=t.Minute()%60
	oneday:=t.Add(-time.Duration(left_hour*3600)*time.Second-time.Duration(left_min*60)*time.Second)
	ids=append(ids,oneday.Unix())


	week:=t.Weekday()
	one_week:=oneday.Add( -(time.Duration(week))*86400*time.Second )
	ids=append(ids,one_week.Unix())

	month:=t.Day()
	one_month:=oneday.Add(-time.Duration(month-1)*86400*time.Second )
	ids=append(ids,one_month.Unix())


	for i := 0; i < MaxPrice; i++ {

		v := &PriceInfo{}
		v.Key = fmt.Sprintf("%s_%s", name, period_key[i])

		p:=&Price{}
		ok,err := DB.GetMysqlConn().Where("symbol=? and id=?",name,ids[i]).Get(p)
		if err!=nil {
			log.Fatal(err.Error())
		}
		if !ok {
			//init
			v.PreData=&proto.PriceCache{
				Id:ids[i],
				Symbol:name,
				Price:c.Price,
				CreatedTime:ids[i],
				Amount:0,
				Vol:0,
				Count:0,
				UsdVol:0,
			}
		}else{
			v.PreData=&proto.PriceCache{
				Id:p.Id,
				Symbol:p.Symbol,
				Price:p.Price,
				CreatedTime:p.CreatedTime,
				Amount:p.Amount,
				Vol:p.Vol,
				Count:p.Count,
				UsdVol:p.UsdVol,
			}
		}

		m.data = append(m.data, v)
	}
	go m.Subscribe()

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


func (s *PriceWorkQuene)ReloadKline()  {
	var period_key = [MaxPrice]string{
		"5min",
		"15min",
		"4hour",
		"1day",
		"1week",
		"1month",
	}
	k:=make([]*Kline,0)

	DB.GetMysqlConn().Find(k)

	for _,v:=range k  {

		t := time.Unix(v.Id, 0)
		if t.Second() == 0 {
			min := t.Minute()
			if min%15 == 0 {
				DB.GetMysqlConn().Insert(v)
				h := t.Hour()
				if h%4 == 0 {
					s.save(FourHourPrice, k)
					if h == 0 {
						s.save(OneDayPrice, k)
						w := t.Weekday()
						if w == 0 {
							s.save(OneWeekPrice, k)
						}
						if t.Day() == 1 {
							s.save(OneMonthPrice, k)
						}
					}
				}
			}
		}
	}
}

func (s *PriceWorkQuene) Subscribe() {
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

		t := time.Unix(k.Id, 0)
		if t.Second() == 0 {
			s.entry = k
			s.updatePrice2(k)
			min := t.Minute()
			if min%15 == 0 {
				s.save(FivteenMinPrice, k)
				h := t.Hour()
				if h%4 == 0 {
					s.save(FourHourPrice, k)
					if h == 0 {
						s.save(OneDayPrice, k)
						w := t.Weekday()
						if w == 0 {
							s.save(OneWeekPrice, k)
						}
						if t.Day() == 1 {
							s.save(OneMonthPrice, k)
						}
					}
				}
			}
		}
	}

}

func (s *PriceWorkQuene) save(period int, data *proto.PriceCache) {
	log.WithFields(log.Fields{
		"Symbol":       s.Symbol,
		"Key":  data.Id,
		"id":data.Id,
		"amount": data.Amount,
		"vol":        data.Vol,
	}).Info("begin price record ")

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

	if amount<0 {
		log.WithFields(log.Fields{
			"Symbol":       s.Symbol,
			"Key":  p.Key,
			"id":data.Id,
			"amount": amount,
			"vol":        vol,
			"low":      low,
			"p.PreData":p.PreData.Id,
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
		Low:    low,
		High:   high,
	})

	if err!=nil {
		return
	}
	err = DB.GetRedisConn().LPush(p.Key, y).Err()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	
	p.PreData = data
}
