package handler

import (
	"digicon/common/convert"
	//"digicon/price_service/exchange"
	"digicon/price_service/model"
	proto "digicon/proto/rpc"
	"golang.org/x/net/context"
	"log"
)

type RPCServer struct{}

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

	cny := model.GetTokenCnyPrice(q.TokenId)
	e := q.GetEntry()
	rsp.Data = model.Calculate(e.Price, e.Amount, cny, q.Symbol)
	return nil
}

func (s *RPCServer) SelfSymbols(ctx context.Context, req *proto.SelfSymbolsRequest, rsp *proto.SelfSymbolsResponse) error {
	//t := new(model.QuenesConfig).GetQuenes(req.Uid)
	return nil
}

func (s *RPCServer) LastPrice(ctx context.Context, req *proto.LastPriceRequest, rsp *proto.LastPriceResponse) error {
	p, ok := model.GetPrice(req.Symbol)
	rsp.Ok=ok
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

func (s *RPCServer) Symbols(ctx context.Context, req *proto.NullRequest, rsp *proto.SymbolsResponse) error {

	rsp.Usdt = new(proto.SymbolsBaseData)
	rsp.Usdt.Data = make([]*proto.SymbolBaseData, 0)
	rsp.Btc = new(proto.SymbolsBaseData)
	rsp.Btc.Data = make([]*proto.SymbolBaseData, 0)
	rsp.Eth = new(proto.SymbolsBaseData)
	rsp.Eth.Data = make([]*proto.SymbolBaseData, 0)
	rsp.Sdc = new(proto.SymbolsBaseData)
	rsp.Sdc.Data = make([]*proto.SymbolBaseData, 0)
	for _, v := range model.ConfigQuenes {
		if v.TokenId == 1 {
			rsp.Usdt.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			cny := model.GetTokenCnyPrice(q.TokenId)

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}
			price := q.GetEntry().Price
			rsp.Usdt.Data = append(rsp.Usdt.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(cny, price)),
				Scope:        convert.Int64DivInt64By8BitString(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		} else if v.TokenId == 2 {
			rsp.Btc.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			cny := model.GetTokenCnyPrice(q.TokenId)

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}
			price := q.GetEntry().Price
			rsp.Btc.Data = append(rsp.Btc.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(cny, price)),
				Scope:        convert.Int64DivInt64By8BitString(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		} else if v.TokenId == 3 {
			rsp.Eth.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			cny := model.GetTokenCnyPrice(q.TokenId)

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}
			price := q.GetEntry().Price
			rsp.Eth.Data = append(rsp.Eth.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(cny, price)),
				Scope:        convert.Int64DivInt64By8BitString(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		} else if v.TokenId == 4 {
			rsp.Sdc.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			cny := model.GetTokenCnyPrice(q.TokenId)

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}
			price := q.GetEntry().Price
			rsp.Sdc.Data = append(rsp.Sdc.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(cny, price)),
				Scope:        convert.Int64DivInt64By8BitString(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		}
	}
	return nil
}
