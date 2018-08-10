package handler

import (
	"context"
	"digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/wallet_service/model"
	"errors"
	"fmt"
	//"strconv"
)

func (s *WalletHandler) CreateBTCWallet(ctx context.Context, req *proto.CreateWalletRequest, rsp *proto.CreateWalletResponse) error {
	fmt.Println(" create btc wallet ...")
	fmt.Println(req.String())

	var err error
	var addr string
	rsp.Data = new(proto.CreateWalletPos)
	tokenId := int(req.Tokenid)
	tokenModel := &Tokens{Id: tokenId}
	_, err = tokenModel.GetByid(tokenId)
	switch tokenModel.Signature {
	case "btc":
		addr, err = NewBTC(int(req.Userid), tokenId, "123456", tokenModel.Chainid)
	default:
		err = errors.New("unknow type ..")
	}

	if err != nil {
		rsp.Code = "1"
		rsp.Msg = err.Error()
		rsp.Data.Type = tokenModel.Signature
		rsp.Data.Addr = ""
		return nil
	}
	rsp.Code = "0"
	rsp.Msg = addr
	rsp.Data.Type = tokenModel.Signature
	rsp.Data.Addr = addr
	return nil
}

func (s *WalletHandler) BtcSigntx(ctx context.Context, req *proto.BtcSigntxRequest, rsp *proto.BtcSigntxResponse) error {
	//fmt.Println("btc signtx request ...")
	//fmt.Println(req.Uid)
	txHash, err := BtcSendToAddress(req.Address, req.Amount, req.Tokenid, int(req.Uid),int(req.Applyid))
	if err != nil {
		rsp.Data = ""
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return nil
	} else {
		rsp.Data = string(txHash)
		rsp.Code = errdefine.ERRCODE_SUCCESS
	}

	return nil
}

/*
	Ti bi提币
*/
func (s *WalletHandler) BtcTibi(ctx context.Context, req *proto.BtcTibiRequest, rsp *proto.BtcResponse) error {
	//fmt.Println(" btc tibi ...")
	//fmt.Println(req.Amount, req.To, req.Tokenid, req.Uid)
	toAddress := req.Address
	mount := req.Amount
	tokenId := req.Tokenid
	txHash, err := BtcTiBiToAddress(toAddress, mount, tokenId, req.Uid,int(req.Applyid))
	if err != nil {
		rsp.Data = ""
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return nil
	} else {
		rsp.Data = string(txHash)
		rsp.Code = errdefine.ERRCODE_SUCCESS
	}

	return nil
}









/////////////////////////////////////////
//比特币相关
//比特币交易











