package handler

import (
	"digicon/common/convert"
	//"digicon/price_service/exchange"
	. "digicon/price_service/dao"
	"digicon/price_service/model"
	proto "digicon/proto/rpc"
	"fmt"
	"github.com/alex023/clock"
	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"strings"
	"time"
	"github.com/shopspring/decimal"
)

type RPCServer struct {
	Publisher micro.Publisher
	topic     string
}

func NewRPCServer(pb micro.Publisher) *RPCServer {
	r := &RPCServer{
		Publisher: pb,
		topic:     "topic.go.micro.srv.price",
	}
	go r.Process()
	return r
}

func (s *RPCServer) Process() {
	for {
		c := clock.NewClock()
		t := time.Now()
		diff := 60 - t.Second()
		job := func() {
			d := make([]*proto.CnyBaseData, 0)
			for _, v := range model.GetQueneMgr().PriceMap {

				d = append(d, &proto.CnyBaseData{
					TokenId:     v.Data.TokenTradeId,
					CnyPrice:    convert.Int64ToStringBy8Bit(v.CnyPrice),
					UsdPrice:    convert.Int64ToStringBy8Bit(v.UsdPrice),
					CnyPriceInt: v.CnyPrice,
					UsdPriceInt: v.UsdPrice,
				})
			}

			if len(d) > 1 {
				s.publishEvent(&proto.CnyPriceResponse{
					Data: d,
				})
			}
		}

		d := time.Duration(diff+10) * time.Second
		c.AddJobWithInterval(d, job)

		time.Sleep(d)
		//log.Info("circle process send price")
	}
}

func (s *RPCServer) publishEvent(data *proto.CnyPriceResponse) error {
	// Marshal to JSON string
	m := jsonpb.Marshaler{EmitDefaults: true}

	g, err := m.MarshalToString(data)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	err = DB.GetRedisConn().Set("history.price.go.micro", g, 0).Err()

	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	if err := s.Publisher.Publish(context.TODO(), data); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *RPCServer) AdminCmd(ctx context.Context, req *proto.AdminRequest, rsp *proto.AdminResponse) error {
	log.Print("Received Say.Hello request")
	rsp.Data = "Hello " + req.Cmd
	return nil
}

func (s *RPCServer) CurrentPrice(ctx context.Context, req *proto.CurrentPriceRequest, rsp *proto.CurrentPriceResponse) error {
	q, ok := model.GetQueneMgr().GetQueneByUKey(req.Symbol)
	if !ok {
		return nil
	}
	e := q.GetEntry()
	h, l, err := q.GetPeriodLHPrice(model.OneDayPrice)
	if err != nil {
		return nil
	}
	rsp.Data = model.Calculate(q.ToekenTradeId, e.Price, e.Amount, q.Symbol, h.Price, l.Price)
	return nil
}

func (s *RPCServer) SelfSymbols(ctx context.Context, req *proto.SelfSymbolsRequest, rsp *proto.SelfSymbolsResponse) error {
	//t := new(model.QuenesConfig).GetQuenes(req.Uid)
	return nil
}

func (s *RPCServer) LastPrice(ctx context.Context, req *proto.LastPriceRequest, rsp *proto.LastPriceResponse) error {
	p, ok := model.GetPrice(req.Symbol)
	rsp.Ok = ok
	if ok {
		rsp.Data = &proto.PriceCache{
			Id:          p.Id,
			Symbol:      p.Symbol,
			Amount:      p.Amount,
			Vol:         p.Vol,
			CreatedTime: p.CreatedTime,
			Count:       p.Count,
		}

		return nil
	}
	return nil
}

func (s *RPCServer) SymbolTitle(ctx context.Context, req *proto.NullRequest, rsp *proto.SymbolTitleResponse) error {
	g := make([]*proto.TitleBaseData, 0)

	for _, v := range model.ConfigTitles {
		g = append(g, &proto.TitleBaseData{
			Mark:    v.Mark,
			TokenId: int32(v.TokenId),
		})
	}
	rsp.Data = g
	return nil
}

func (s *RPCServer) SymbolsById(ctx context.Context, req *proto.SymbolsByIdRequest, rsp *proto.SymbolsByIdResponse) error {
	g := model.GetConfigQuenesByType(req.TokenId)
	rsp.Data = make([]*proto.SymbolBaseData, 0)

	for _, v := range g {
		q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
		if !ok {
			return nil
		}

		p, ok := model.Get24HourPrice(v.Name)
		if !ok {
			g := model.ConfigQueneInit[v.Name]
			p.Price = g.Price
			p.Amount = 0
			log.Errorf("SymbolsById not found name %s", v.Name)
		}

		price := q.GetEntry().Price

		c, ok := model.GetQueneMgr().PriceMap[v.TokenTradeId]
		if !ok {
			log.Errorf("SymbolsById not found TokenTradeId %s", v.TokenTradeId)
			continue
		}

		rsp.Data = append(rsp.Data, &proto.SymbolBaseData{
			Symbol: v.Name,
			Price:  convert.Int64ToStringBy8Bit(price),
			//CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(q.CnyPrice, price)),
			CnyPrice:     convert.Int64ToStringBy8Bit(c.CnyPrice),
			Scope:        convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
			TradeTokenId: int32(v.TokenTradeId),
		})
	}

	return nil
}

func (s *RPCServer) Symbols(ctx context.Context, req *proto.NullRequest, rsp *proto.SymbolsResponse) error {

	rsp.Usdt = new(proto.SymbolsBaseData)
	rsp.Usdt.Data = make([]*proto.SymbolBaseData, 0)
	rsp.Btc = new(proto.SymbolsBaseData)
	rsp.Btc.Data = make([]*proto.SymbolBaseData, 0)
	rsp.Eth = new(proto.SymbolsBaseData)
	rsp.Eth.Data = make([]*proto.SymbolBaseData, 0)
	rsp.Sdc = new(proto.SymbolsBaseData)
	rsp.Sdc.Data = make([]*proto.SymbolBaseData, 0)
	for _, v := range model.ConfigQueneData {
		if v.TokenId == 1 {
			rsp.Usdt.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}

			price := q.GetEntry().Price

			rsp.Usdt.Data = append(rsp.Usdt.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(q.CnyPrice, price)),
				Scope:        convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		} else if v.TokenId == 2 {
			rsp.Btc.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}
			price := q.GetEntry().Price
			rsp.Btc.Data = append(rsp.Btc.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(q.CnyPrice, price)),
				Scope:        convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		} else if v.TokenId == 3 {
			rsp.Eth.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}
			price := q.GetEntry().Price
			rsp.Eth.Data = append(rsp.Eth.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(q.CnyPrice, price)),
				Scope:        convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		} else if v.TokenId == 4 {
			rsp.Sdc.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}
			price := q.GetEntry().Price
			rsp.Sdc.Data = append(rsp.Sdc.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(q.CnyPrice, price)),
				Scope:        convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		}
	}
	return nil
}

func (s *RPCServer) GetCnyPrices(ctx context.Context, req *proto.CnyPriceRequest, rsp *proto.CnyPriceResponse) error {

	d := make([]*proto.CnyBaseData, 0)
	for _, v := range req.TokenTradeId {
		g, ok := model.GetQueneMgr().PriceMap[v]
		if !ok {
			continue
		}

		d = append(d, &proto.CnyBaseData{
			TokenId:     v,
			CnyPrice:    convert.Int64ToStringBy8Bit(g.CnyPrice),
			UsdPrice:    convert.Int64ToStringBy8Bit(g.UsdPrice),
			CnyPriceInt: g.CnyPrice,
			UsdPriceInt: g.UsdPrice,
		})
	}
	rsp.Data = d
	return nil
}

func (s *RPCServer) Quotation(ctx context.Context, req *proto.QuotationRequest, rsp *proto.QuotationResponse) error {
	g := model.GetConfigQuenesByType(req.TokenId)

	for _, v := range g {
		q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
		if !ok {
			return nil
		}

		price := q.GetEntry().Price
		h, l, err := q.GetPeriodLHPrice(model.OneDayPrice)
		if err != nil {
			return nil
		}
		r := model.Calculate(q.ToekenTradeId, price, q.GetEntry().Amount, q.Symbol, h.Price, l.Price)

		cny_price, ok := model.GetQueneMgr().PriceMap[q.ToekenTradeId]
		if ok {
			rsp.Data = append(rsp.Data, &proto.QutationBaseData{
				Symbol:   v.Name,
				Price:    r.Price,
				Scope:    r.Scope,
				Low:      r.Low,
				High:     r.High,
				Amount:   r.Amount,
				CnyPrice: convert.Int64ToStringBy8Bit(cny_price.CnyPrice),
				CnyLow:   convert.Int64ToStringBy8Bit(l.CnyPrice),
				CnyHigh:  convert.Int64ToStringBy8Bit(h.CnyPrice),
			})
		} else {
			rsp.Data = append(rsp.Data, &proto.QutationBaseData{
				Symbol: v.Name,
				Price:  r.Price,
				Scope:  r.Scope,
				Low:    r.Low,
				High:   r.High,
				Amount: r.Amount,
			})
		}
	}
	return nil
}

/*
	获取一个币对的价格比
*/
func (s *RPCServer) GetSymbolsRate(ctx context.Context, req *proto.GetSymbolsRateRequest, rsp *proto.GetSymbolsRateResponse) error {
	type BaseData struct {
		Symbol   string `json:"symbol"`    //  btc/usdt
		Price    string `json:"price"`     //  1btc = xxx usdt
		CnyPrice string `json:"cny_price"` //  cny
	}
	data := map[string]*proto.RateBaseData{}
	for _, symbol := range req.Symbols {
		var ok bool
		data[symbol], ok = getSymbolRate(symbol)
		if !ok {
			tmpdata, ok := getOtherSymbolRage(symbol)
			if !ok {
				continue
			}
			data[symbol] = tmpdata
		}

	}
	rsp.Data = data
	return nil
}

func getSymbolRate(symbol string) (data *proto.RateBaseData, ok bool) {
	q, ok := model.GetQueneMgr().GetQueneByUKey(symbol)
	if !ok {
		fmt.Println(symbol, ok)
		//return getOtherSymbolRage(symbol)
		return
	} else {
		e := q.GetEntry()
		data = &proto.RateBaseData{
			Symbol:   q.Symbol,
			Price:    convert.Int64ToStringBy8Bit(e.Price),
			CnyPrice: convert.Int64ToStringBy8Bit(q.CnyPrice),
		}
	}
	return
}

/*
   如果没有找到币对
*/
func getOtherSymbolRage(symbol string) (data *proto.RateBaseData, ok bool) {
	//fmt.Println("other symbol rage ...", symbol)
	tmpSym := strings.Split(symbol, "/")
	if len(tmpSym) < 2 {
		return
	}
	tmpSymToUSDT := fmt.Sprintf("%s/USDT", tmpSym[0])
	toUSDTQ, ok := model.GetQueneMgr().GetQueneByUKey(tmpSymToUSDT)
	if !ok {
		fmt.Println(tmpSymToUSDT, " not ok!!!")
		return
	}
	usdtToTmpSym := fmt.Sprintf("%s/USDT", tmpSym[1])
	usdtToQ, ok := model.GetQueneMgr().GetQueneByUKey(usdtToTmpSym)
	if !ok {
		fmt.Println(usdtToTmpSym, " not ok!!!!")
		return
	}
	BTCPrice := toUSDTQ.GetEntry()
	OtherPrice := usdtToQ.GetEntry()
	price := convert.Int64DivInt64By8Bit(BTCPrice.Price, OtherPrice.Price)
	cnyPrice := convert.Int64MulInt64By8Bit(OtherPrice.Price, toUSDTQ.CnyPrice)
	data = &proto.RateBaseData{
		Symbol:   symbol,
		Price:    convert.Int64ToStringBy8Bit(price),
		CnyPrice: convert.Int64ToStringBy8Bit(cnyPrice),
	}
	return
}

func (s *RPCServer) Volume(ctx context.Context, req *proto.VolumeRequest, rsp *proto.VolumeResponse) error {
	nowSum,daySum,weekSum,monthSum := model.GetVolumeTotal()

	w,_ := decimal.NewFromString("100000000")
	now_sum,_ := decimal.NewFromString(nowSum)

	day_sum,_ := decimal.NewFromString(daySum)
	day_sum = now_sum.Sub(day_sum)
	day_sum = day_sum.Div(w)

	week_sum,_ := decimal.NewFromString(weekSum)
	week_sum = now_sum.Sub(week_sum)
	week_sum = week_sum.Div(w)

	month_sum,_ := decimal.NewFromString(monthSum)
	month_sum = now_sum.Sub(month_sum)
	month_sum = month_sum.Div(w)

	rsp.DayVolume = day_sum.IntPart()
	rsp.WeekVolume = week_sum.IntPart()
	rsp.MonthVolume = month_sum.IntPart()
	return nil
}
