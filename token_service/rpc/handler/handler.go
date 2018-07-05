package handler

import (
	"digicon/common/convert"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"digicon/token_service/model"
	"github.com/go-redis/redis"
	"golang.org/x/net/context"
	"log"
	"github.com/liudng/godump"
)

type RPCServer struct {
}

func (s *RPCServer) AdminCmd(ctx context.Context, req *proto.AdminRequest, rsp *proto.AdminResponse) error {
	log.Print("Received Say.Hello request")
	rsp.Data = "Hello " + req.Cmd
	return nil
}

func (s *RPCServer) EntrustOrder(ctx context.Context, req *proto.EntrustOrderRequest, rsp *proto.CommonErrResponse) error {
	q, ok := model.GetQueneMgr().GetQueneByUKey(req.Symbol)
	if !ok {
		rsp.Err = ERR_TOKEN_QUENE_CONF
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	ret, err := q.EntrustReq(req)
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}
	rsp.Err = ret
	rsp.Message = GetErrorMessage(rsp.Err)
	/*
		q.JoinSellQuene(&model.EntrustDetail{
			EntrustId:  genkey.GetTimeUnionKey(q.GetUUID()),
			Uid:        int(req.Uid),
			AllNum:     req.Num,
			SurplusNum: req.Num,
			Opt:        int(req.Opt),
			OnPrice:    req.OnPrice,
			States:     0,
		})
	*/
	return nil
}

func (s *RPCServer) Symbols(ctx context.Context, req *proto.SymbolsRequest, rsp *proto.SymbolsResponse) error {
	t := new(model.QuenesConfig).GetQuenes(req.Type)

	for _, v := range t {
		//g := convert.Int64ToStringBy8Bit(v.Price)
		rsp.Data = append(rsp.Data, &proto.SymbolBaseData{
			Symbol:   v.Name,
			Price:    convert.Int64ToStringBy8Bit(v.Price),
			CnyPrice: convert.Int64ToStringBy8Bit(7 * v.Price),
			Scope:    v.Scope,
		})
	}

	return nil
}

func (s *RPCServer) AddTokenNum(ctx context.Context, req *proto.AddTokenNumRequest, rsp *proto.CommonErrResponse) error {
	ret, err := model.AddTokenSess(req)
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
	}
	rsp.Err = ret
	rsp.Message = GetErrorMessage(ret)
	return nil
}

func (s *RPCServer) HistoryKline(ctx context.Context, req *proto.HistoryKlineRequest, rsp *proto.HistoryKlineResponse) error {
	return nil
}

func (s *RPCServer) EntrustQuene(ctx context.Context, req *proto.EntrustQueneRequest, rsp *proto.EntrustQueneResponse) error {
	q, ok := model.GetQueneMgr().GetQueneByUKey(req.Symbol)
	if !ok {
		return nil
	}
	others, err := q.PopFirstEntrust(proto.ENTRUST_OPT_BUY, 2, 5)
	if err == redis.Nil {

	} else if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	} else {
		for _, v := range others {
			g := &proto.EntrustBaseData{
				OnPrice:    convert.Int64ToStringBy8Bit(v.OnPrice),
				SurplusNum: convert.Int64ToStringBy8Bit(v.SurplusNum),
			}
			g.Price = convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(v.OnPrice, v.SurplusNum))
			rsp.Buy = append(rsp.Buy, g)
		}
	}

	others, err = q.PopFirstEntrust(proto.ENTRUST_OPT_SELL, 2, 5)
	if err == redis.Nil {

	} else if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	} else {
		for _, v := range others {
			g := &proto.EntrustBaseData{
				OnPrice:    convert.Int64ToStringBy8Bit(v.OnPrice),
				SurplusNum: convert.Int64ToStringBy8Bit(v.SurplusNum),
			}
			g.Price = convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(v.OnPrice, v.SurplusNum))
			rsp.Sell = append(rsp.Sell, g)
		}
	}

	rsp.Err = ERRCODE_SUCCESS
	return nil
}

/*
func (s *RPCServer) EntrustList(ctx context.Context, req *proto.EntrustQueneRequest, rsp *proto.EntrustQueneResponse) error {
	a := make([][]interface{}, 0)
	godump.Dump(a)
	return nil
}
*/
func (s *RPCServer) EntrustHistory(ctx context.Context, req *proto.EntrustHistoryRequest, rsp *proto.EntrustHistoryResponse) error {

	r := new(model.EntrustDetail).GetHistory(req.Uid, int(req.Limit), int(req.Page))
	for _, v := range r {

		rsp.Data=append(rsp.Data,&proto.EntrustHistoryBaseData{
			EntrustId:v.EntrustId,
			Symbol:v.Symbol,
			Opt:proto.ENTRUST_OPT(v.Opt),
			Type:proto.ENTRUST_TYPE(v.Type),
			AllNum:v.AllNum,
			OnPrice:v.OnPrice,
			SurplusNum:v.SurplusNum,
			CreateTime:v.CreatedTime,
			States:int32(v.States),
		})
	}
	return nil
}

func (s *RPCServer) EntrustList(ctx context.Context, req *proto.EntrustHistoryRequest, rsp *proto.EntrustHistoryResponse) error {
	r := new(model.EntrustDetail).GetList(req.Uid, int(req.Limit), int(req.Page))
	godump.Dump(len(r))
	for _, v := range r {

		rsp.Data=append(rsp.Data,&proto.EntrustHistoryBaseData{
			EntrustId:v.EntrustId,
			Symbol:v.Symbol,
			Opt:proto.ENTRUST_OPT(v.Opt),
			Type:proto.ENTRUST_TYPE(v.Type),
			AllNum:v.AllNum,
			OnPrice:v.OnPrice,
			SurplusNum:v.SurplusNum,
			CreateTime:v.CreatedTime,
			States:int32(v.States),
		})
	}
	return nil
}