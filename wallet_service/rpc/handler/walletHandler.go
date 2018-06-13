package handler

import (
	proto "digicon/proto/rpc"
	"context"
	"log"
	. "digicon/wallet_service/model"
	"errors"
	"fmt"
)

type Wallet struct{}

func (s *Wallet) Hello(ctx context.Context, req *proto.HelloRequest2, rsp *proto.HelloResponse2) error {
	log.Print("Received Say.Hello request")
	rsp.Greeting = "Hello darnice mani" + req.Name
	return nil
}

func (s *Wallet) CreateWallet(ctx context.Context, req *proto.CreateWalletRequest, rsp *proto.CreateWalletResponse) error {
	log.Print("Received Say.CreateWallet request")
	var err error
	var addr string

	tokenModel := &Tokens{Id:int(req.Tokenid)}
	_,err =tokenModel.GetByid()


	switch tokenModel.Type {
	case "eth":
		addr,err = Neweth(int(req.Userid),int(req.Tokenid),"123456",)
	default:
		err = errors.New("unknow type")
	}

	if err != nil {
		rsp.Code = string("1")
		rsp.Type = "eth"
		rsp.Msg	 = err.Error()
		fmt.Println(rsp)
		return nil
	}
	rsp.Code = "0"
	rsp.Type = "eth"
	rsp.Addr = addr
	rsp.Msg= "生成成功"
	return nil
}