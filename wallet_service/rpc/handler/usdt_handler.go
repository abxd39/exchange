package handler

import (
	"context"
	proto "digicon/proto/rpc"

	"fmt"
	log "github.com/sirupsen/logrus"
	. "digicon/wallet_service/model"
	. "digicon/wallet_service/utils"
	"digicon/common/errors"
	."digicon/proto/common"
)

//usdt提币

func (s *WalletHandler) CreateUSDTWallet(ctx context.Context, req *proto.CreateWalletRequest, rsp *proto.CreateWalletResponse) error {
	fmt.Println("create usdt wallet ...")
	var err error
	fmt.Println(err)

	return nil
}


func (s *WalletHandler) UsdtTiBi(ctx context.Context, req *proto.UsdtTiBiRequest, rsp *proto.UsdtTiBiResponse) error {
	
	defer func() {
		if rsp.Code != ERRCODE_SUCCESS {
			log.WithFields(log.Fields{
				"code":rsp.Code,
				"msg":rsp.Message,
			}).Error("UsdtTiBi error")
		}
		fmt.Println("提币结果：",rsp.Code,rsp.Message)
	}()
	
	fromAddress := req.FromAddress
	toAddress := req.ToAddress
	protertyid := req.Protertyid
	mount := req.Amount
	tokenId := req.Tokenid
	txHash, err := s.UsdtSendToAddress(fromAddress,toAddress, int(protertyid),mount, tokenId, int(req.Uid), int(req.Applyid))
	if err != nil {
		rsp.Message = err.Error()
		rsp.Data = ""
		rsp.Code = ERRCODE_UNKNOWN
		return errors.New(rsp.Message)
	}
	if txHash == "" {
		rsp.Message = "提币失败，hash为空"
		rsp.Data = ""
		rsp.Code = ERRCODE_UNKNOWN
		return errors.New(rsp.Message)
	}
	rsp.Message = "提币成功"
	rsp.Data = string(txHash)
	rsp.Code = ERRCODE_SUCCESS

	return nil
}

func (s *WalletHandler) UsdtSendToAddress(fromAddress string,toAddress string,propertyid int, mount string, tokenId int32, uid int, applyid int) (string, error) {
	log.Info("USDT交易request:",fromAddress,toAddress,propertyid,mount,tokenId,uid,applyid)
	wToken := new(WalletToken)
	wToken.GetByUid(uid)
	password := wToken.Password

	token := Tokens{}
	token.GetByid(int(tokenId))
	url := token.Node

	err := UsdtWalletPhrase(url, password, 1*60*60)

	if err != nil {
		msg := "USDT钱包解锁失败!"
		log.Errorln("USDT钱包解锁失败",msg)
		return "", err
	}

	enough, err := UsdtCheckBalance(url,fromAddress,propertyid, mount)
	if err != nil {
		log.Error("余额检查失败：",err)
		return "",err
	}
	if !enough {
		msg := "balance not enough!"
		err = errors.New(msg)
		log.Errorln("USDT余额不足",msg)
		return "", err
	}
	//fmt.Println("btc send before ...")
	txHash, err := UsdtSendToAddressFunc(url, fromAddress,toAddress,propertyid,mount)
	if err != nil {
		fmt.Println("USDT发送交易失败",err.Error())
		log.Errorf("USDT发送交易失败",err.Error())
		return "", err
	}

	tio := new(TokenInout) //
	//更新
	row, err := tio.UpdateApplyTiBi(applyid, txHash)

	if err != nil {
		log.Errorln(err.Error())
		fmt.Println(err.Error())
		return "",err
	}
	if row <= 0 {
		log.Println("更新失败：",applyid)
		fmt.Println("更新失败：",applyid)
	}

	return txHash, err
}