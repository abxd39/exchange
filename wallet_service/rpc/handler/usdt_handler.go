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
	"digicon/wallet_service/watch"
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
		fmt.Println("USDT 提币结果：",rsp.Code,rsp.Message)
		if rsp.Code != ERRCODE_SUCCESS {
			log.WithFields(log.Fields{
				"code":rsp.Code,
				"msg":rsp.Message,
			}).Error("UsdtTiBi error")
			log.Error("USDT广播失败，改回状态:",rsp.Message)
			//把状态改回去
			//更新申请单记录
			_,err := new(TokenInout).UpdateApplyTiBi2(int(req.Applyid),4,rsp.Message)  //正在提币中
			if err != nil {
				log.Error("USDT UpdateApplyTiBi error:",err)
			}
			//取消冻结
			s.tiBiErrorUnfreeze(req.Applyid)
		}
	}()
	
	fromAddress := req.FromAddress
	toAddress := req.ToAddress
	protertyid := req.Protertyid
	mount := req.Amount
	tokenId := req.Tokenid

	if req.Tokenid != 1 {
		rsp.Message = "Tokenid错误"
		rsp.Data = ""
		rsp.Code = ERRCODE_UNKNOWN
		return errors.New(rsp.Message)
	}

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

	fmt.Println("-----------------------USDT:",url)

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
	fmt.Println("USDT交易结果：",txHash,err)
	if err != nil {
		fmt.Println("USDT发送交易失败",err.Error())
		log.Errorf("USDT发送交易失败",err.Error())
		return "", err
	}

	//添加txhash到监控队列
	new(watch.USDTTiBiWatch).PushRedisList(txHash)

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