package watch

import (
	"digicon/wallet_service/model"
	proto "digicon/proto/rpc"
	"digicon/wallet_service/rpc/client"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Common struct{}

//添加比特币token数量
//**比特币充币
func (p *Common) AddBTCTokenNum(data TranItem) {
	log.Info("AddBTCTokenNum data:",data)
	//查询用户uid
	walletToken := new(models.WalletToken)
	err := walletToken.GetByAddress(data.Address)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("get user token error",err)
		return
	}
	amount := int64(data.Amount * 100000000)
	rsp,errr := client.InnerService.TokenSevice.CallAddTokenNum(&proto.AddTokenNumRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Num:amount,
		Ukey:[]byte(data.Txid),
		OptAddType:0,
	})
	log.Info("btc AddBTCTokenNum result",err,rsp)
	if errr != nil {
		log.Info("AddBTCTokenNum error",err)
	}
}

//确认消耗
//**比特币提币成功调用
func (p *Common) BTCConfirmSubFrozen(data TranItem) {
	//查询用户uid
	walletToken := new(models.WalletToken)
	err := walletToken.GetByAddress(data.Address)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("get user token error",err)
		return
	}
	//根据交易hash查询申请提币数据
	tokenInout := new(models.TokenInout)
	err = tokenInout.GetByHash(data.Txid)
	if err != nil || tokenInout.Uid <= 0 {
		log.Info("get data by hash error",err)
		return
	}
	_,errr := client.InnerService.TokenSevice.CallConfirmSubFrozen(&proto.ConfirmSubFrozenRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Num:tokenInout.Amount,
		Ukey:[]byte(data.Txid),
		Type:1,  //区块入账
	})
	if errr != nil {
		log.Info("BTCConfirmSubFrozen error",err)
	}
}

//确认消耗
//**以太坊提币成功调用
func (p *Common) ETHConfirmSubFrozen(from string,txhash string) {
	//查询用户uid
	walletToken := new(models.WalletToken)
	err := walletToken.GetByAddress(from)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("get user token error",err)
		return
	}
	//根据交易hash查询申请提币数据
	tokenInout := new(models.TokenInout)
	err = tokenInout.GetByHash(txhash)
	if err != nil || tokenInout.Uid <= 0 {
		log.Info("get data by hash error",err)
		return
	}
	_,errr := client.InnerService.TokenSevice.CallConfirmSubFrozen(&proto.ConfirmSubFrozenRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Num:tokenInout.Amount,
		Ukey:[]byte(txhash),
		Type:1,  //区块入账
	})
	if errr != nil {
		log.Info("ETHConfirmSubFrozen error",err)
	}
}

//添加比特币token数量
//**以太坊充币
func (p *Common) AddETHTokenNum(to string,tokenid int,amount string,txhash string) {
	//查询用户uid
	walletToken := new(models.WalletToken)
	err := walletToken.GetByAddress(to)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("get user token error",err)
		return
	}

	bigAmount,err := strconv.ParseInt(amount,10,64)

	if err != nil {
		log.Error("AddETHTokenNum",err)
		return
	}

	rsp,errr := client.InnerService.TokenSevice.CallAddTokenNum(&proto.AddTokenNumRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Opt:1,
		Num:bigAmount,
		Ukey:[]byte(txhash),
		OptAddType:1,
		Type:1,
	})
	log.WithFields(log.Fields{
		"Uid":uint64(walletToken.Uid),
		"TokenId":int32(walletToken.Tokenid),
		"Opt":1,
		"Num":bigAmount,
		"Ukey":[]byte(txhash),
		"OptAddType":1,
		"Type":1,
	}).Info("rpc------")
	log.Info("eth AddETHTokenNum result",errr,rsp)
	if errr != nil {
		log.Info("AddBTCTokenNum error",errr)
	}
}

//测试
//func init() {
//	walletToken := new(models.WalletToken)
//	err := walletToken.GetByAddress("121212112121212")
//	if err != nil {
//		log.Println("get user token error",err.Error())
//		return
//	}
//	fmt.Println("测试数据：",err,walletToken.Uid <= 0)
//}
