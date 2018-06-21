package handler

import (
	"digicon/common/genkey"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"digicon/token_service/model"
	"github.com/liudng/godump"
	"golang.org/x/net/context"
	"log"
)

type RPCServer struct{}

func (s *RPCServer) AdminCmd(ctx context.Context, req *proto.AdminRequest, rsp *proto.AdminResponse) error {
	log.Print("Received Say.Hello request")
	rsp.Data = "Hello " + req.Cmd
	return nil
}

func (s *RPCServer) EntrustOrder(ctx context.Context, req *proto.EntrustOrderRequest, rsp *proto.EntrustOrderResponse) error {
	q, ok := model.GetQueneMgr().GetQuene(int(req.TokenId), int(req.TokenTradeId))
	if !ok {
		rsp.Err = ERR_TOKEN_QUENE_CONF
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	q.JoinSellQuene(&model.EntrustDetail{
		EntrustId:  genkey.GetTimeUnionKey(q.GetUUID()),
		Uid:        int(req.Uid),
		AllNum:     req.Num,
		SurplusNum: req.Num,
		Opt:        int(req.Opt),
		OnPrice:    req.OnPrice,
		States:     0,
	})

	return nil
}

func (s *RPCServer) AddTokenNum(ctx context.Context, req *proto.AddTokenNumRequest, rsp *proto.CommonErrResponse) error {
	u := &model.UserToken{}
	var err error
	var ret int32
	err = u.GetUserToken(int(req.Uid), int(req.TokenId))
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}

	if req.Opt { //减少类型
		ret, err = u.SubMoney(req.Num, string(req.Hash))
	} else { //增加类型
		ret, err = u.AddMoney(req.Num, string(req.Hash))
	}

	rsp.Err = ret
	rsp.Message = err.Error()
	return nil
}

func Test() {
	req := &proto.AddTokenNumRequest{
		Uid:     8,
		TokenId: 4,
		Num:     10000000,
		Hash:    []byte("dasfdsaonzz11opqqq11+="),
		Opt:     false,
	}

	u := &model.UserToken{}
	var err error
	var ret int32
	err = u.GetUserToken(int(req.Uid), int(req.TokenId))
	if err != nil {
		godump.Dump(err.Error())
		return
	}

	if req.Opt { //减少类型
		ret, err = u.SubMoney(req.Num, string(req.Hash))

	} else { //增加类型
		ret, err = u.AddMoney(req.Num, string(req.Hash))

	}
	godump.Dump(ret)
}
