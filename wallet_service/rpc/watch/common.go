package watch

import (
	"log"
	"digicon/wallet_service/model"
	proto "digicon/proto/rpc"
	"digicon/wallet_service/rpc/client"
	"fmt"
)

type Common struct{}

//添加比特币token数量
//**比特币充币
func (p *Common) AddBTCTokenNum(data TranItem) {
	//查询用户uid
	walletToken := new(models.WalletToken)
	err := walletToken.GetByAddress(data.Address)
	if err != nil || walletToken.Uid <= 0 {
		log.Println("get user token error",err.Error())
		return
	}
	amount := int64(data.Amount * 100000000)
	_,errr := client.InnerService.TokenSevice.CallAddTokenNum(&proto.AddTokenNumRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Num:amount,
		Ukey:[]byte(data.Txid),
		OptAddType:0,
	})
	if errr != nil {
		log.Println("AddBTCTokenNum error",err.Error())
	}
}

//确认消耗
//**比特币提币成功调用
func (p *Common) BTCConfirmSubFrozen(data TranItem) {
	//查询用户uid
	walletToken := new(models.WalletToken)
	err := walletToken.GetByAddress(data.Address)
	fmt.Println("错误数据：",err,walletToken.Uid)
	if err != nil || walletToken.Uid <= 0 {
		log.Println("get user token error",err.Error())
		return
	}
	//根据交易hash查询申请提币数据
	tokenInout := new(models.TokenInout)
	err = tokenInout.GetByHash(data.Txid)
	if err != nil || tokenInout.Uid <= 0 {
		log.Println("get data by hash error",err.Error())
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
		log.Println("BTCConfirmSubFrozen error",err.Error())
	}
}

//确认消耗
//**以太坊提币成功调用
func (p *Common) ETHConfirmSubFrozen(from string,txhash string) {
	//查询用户uid
	walletToken := new(models.WalletToken)
	err := walletToken.GetByAddress(from)
	if err != nil || walletToken.Uid <= 0 {
		log.Println("get user token error",err.Error())
		return
	}
	//根据交易hash查询申请提币数据
	tokenInout := new(models.TokenInout)
	err = tokenInout.GetByHash(txhash)
	if err != nil || tokenInout.Uid <= 0 {
		log.Println("get data by hash error",err.Error())
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
		log.Println("ETHConfirmSubFrozen error",err.Error())
	}
}

//添加比特币token数量
//**以太坊充币
func (p *Common) AddETHTokenNum(to string,tokenid int,amount int64,txhash string) {
	//查询用户uid
	walletToken := new(models.WalletToken)
	err := walletToken.GetByAddress(to)
	if err != nil || walletToken.Uid <= 0 {
		log.Println("get user token error",err.Error())
		return
	}
	amount = amount * 100000000
	_,errr := client.InnerService.TokenSevice.CallAddTokenNum(&proto.AddTokenNumRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Num:amount,
		Ukey:[]byte(txhash),
		OptAddType:0,
	})
	if errr != nil {
		log.Println("AddBTCTokenNum error",err.Error())
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
