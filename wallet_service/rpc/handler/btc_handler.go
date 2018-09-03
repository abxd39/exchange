package handler

import (
	"context"
	"digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/wallet_service/model"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	. "digicon/wallet_service/utils"
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
	txHash, err := BtcSendToAddress(req.Address, req.Amount, req.Tokenid, int(req.Uid), int(req.Applyid))
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
	toAddress := req.Address
	mount := req.Amount
	tokenId := req.Tokenid
	txHash, err := BtcSendToAddress(toAddress, mount, tokenId,int(req.Uid), int(req.Applyid))
	if err != nil {
		log.Error("BtcTibi error:",err)
		rsp.Message = "提币失败"
		rsp.Data = ""
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return nil
	} else {
		rsp.Message = "提币成功"
		rsp.Data = string(txHash)
		rsp.Code = errdefine.ERRCODE_SUCCESS
	}

	return nil
}

/////////////////////////////////////////
//比特币相关
//比特币交易
//
// btc send to address
//
func BtcSendToAddress(toAddress string, mount string, tokenId int32, uid int, applyid int) (string, error) {
	log.Info("比特币交易request:",toAddress,mount,tokenId,uid,applyid)
	wToken := new(WalletToken)
	wToken.GetByUid(uid)
	password := wToken.Password

	token := Tokens{}
	token.GetByid(int(tokenId))
	url := token.Node

	//fmt.Println("----------------------------")
	err := BtcWalletPhrase(url, password, 1*60*60)

	if err != nil {
		msg := "钱包解锁失败!"
		log.Errorln("比特币钱包解锁失败",msg)
		return "", nil
	}

	enough, err := BtcCheckBalance(url, mount)
	if !enough {
		msg := "balance not enough!"
		err = errors.New(msg)
		log.Errorln("比特币余额不足",msg)
		return "", err
	}
	//fmt.Println("btc send before ...")
	txHash, err := BtcSendToAddressFunc(url, toAddress, mount)
	if err != nil {
		fmt.Println(err.Error())
		log.Errorf("比特币发送交易失败",err.Error())
		return "", err
	}
	//amount, err := convert.StringToInt64By8Bit(mount)
	//if err != nil {
	//	fmt.Println(err)
	//}
	tio := new(TokenInout) //
	//更新
	row, err := tio.UpdateApplyTiBi(applyid, txHash)
	//row, err := tio.BtcInsert(txHash, wToken.Address, toAddress, "BTC", amount,
	//	wToken.Chainid, int(tokenId), 0, int(uid),
	//)

	if err != nil || row <= 0 {
		log.Errorln(err.Error())
		fmt.Println(err.Error())
	}

	return txHash, err
}

