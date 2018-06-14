package handler

import (
	"context"
	proto "digicon/proto/rpc"
	. "digicon/wallet_service/model"
	"errors"
	"fmt"
	"log"
)

type Wallet struct{}

func (s *Wallet) Hello(ctx context.Context, req *proto.HelloRequest2, rsp *proto.HelloResponse2) error {
	log.Print("Received Say.Hello request")
	rsp.Greeting = "Hello darnice mani" + req.Name
	return nil
}

func (s *Wallet) CreateWallet(ctx context.Context, req *proto.CreateWalletRequest, rsp *proto.CreateWalletResponse) error {
	log.Print("Received Say.CreateWallet request")
	fmt.Println(req.String())
	var err error
	var addr string
	rsp.Data = new(proto.CreateWalletPos)
	tokenModel := &Tokens{Id: int(req.Tokenid)}
	_, err = tokenModel.GetByid()

	switch tokenModel.Type {
	case "eth":
		addr, err = Neweth(int(req.Userid), int(req.Tokenid), "123456")
	default:
		err = errors.New("unknow type")
	}

	if err != nil {
		rsp.Code = string("1")
		rsp.Msg = err.Error()
		rsp.Data.Type = "eth"
		rsp.Data.Addr = ""

		return nil
	}
	rsp.Code = "0"
	rsp.Msg = addr
	rsp.Data.Type = "eth"
	rsp.Data.Addr = addr

	return nil
}
