package watch

import (
	."digicon/wallet_service/model"
	proto "digicon/proto/rpc"
	"digicon/wallet_service/rpc/client"
	log "github.com/sirupsen/logrus"
	"strconv"
	"fmt"
	"digicon/common/convert"
	"github.com/shopspring/decimal"
)

type Common struct{}

//添加比特币token数量
//**比特币充币
func (p *Common) AddBTCTokenNum(data TranItem) {
	log.Info("AddBTCTokenNum data:",data)
	//查询用户uid
	walletToken := new(WalletToken)
	err := walletToken.GetByAddress(data.Address)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("get user token error",err)
		return
	}

	amount := convert.Float64ToInt64By8Bit(data.Amount)

	log.WithFields(log.Fields{
		"uid":walletToken.Uid,
		"token_id":walletToken.Tokenid,
		"num":amount,
		"ukey":data.Txid,
		"optAddType":0,
	}).Info("比特币RPC新增数量")

	//amount := int64(data.Amount * 100000000)
	rsp,errr := client.InnerService.TokenSevice.CallAddTokenNum(&proto.AddTokenNumRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Num:amount,
		Ukey:[]byte(data.Txid),
		OptAddType:1,
		Type:proto.TOKEN_TYPE_OPERATOR_HISTORY_HASH,
		Opt:proto.TOKEN_OPT_TYPE_ADD,
	})
	log.Info("btc AddBTCTokenNum result",errr,string(rsp.Message))
	if errr != nil {
		log.Info("AddBTCTokenNum error",errr)
	}
}

//确认消耗
//**比特币提币成功调用
func (p *Common) BTCConfirmSubFrozen(data TranItem) {
	//查询用户uid
	walletToken := new(WalletToken)
	err := walletToken.GetByAddress(data.Address)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("btc get user token error",err)
		return
	}

	//根据交易hash查询申请提币数据
	tokenInout := new(TokenInout)
	err = tokenInout.GetByHash(data.Txid)
	if err != nil || tokenInout.Uid <= 0 {
		log.Info("btc get data by hash error",err)
		return
	}

	defer func() {
		log.WithFields(log.Fields{
			"uid":walletToken.Uid,
			"token_id":walletToken.Tokenid,
			"num":tokenInout.Amount,
			"ukey":data.Txid,
			"type":proto.TOKEN_TYPE_OPERATOR_HISTORY_TOKEN_OUT,
		}).Info("比特币冻结数量")
	}()

	rsp,errr := client.InnerService.TokenSevice.CallConfirmSubFrozen(&proto.ConfirmSubFrozenRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Num:tokenInout.Amount,
		Ukey:[]byte(data.Txid),
		Type:proto.TOKEN_TYPE_OPERATOR_HISTORY_TOKEN_OUT,  //提币成功消耗冻结
	})
	log.Info("比特币确认消耗冻结：",errr,rsp.Err,string(rsp.Message))
	if errr != nil {
		log.Info("BTCConfirmSubFrozen error",err)
	}
}

//确认消耗
//**以太坊提币成功调用
func (p *Common) ETHConfirmSubFrozen(from string,txhash string,contract string) {
	//查询用户uid
	//walletToken := new(models.WalletToken)
	//boo,err := walletToken.GetByAddressContract(from,contract)
	//if err != nil || boo != true {
	//	log.Error("get user token error",err,from)
	//	return
	//}
	//根据交易hash查询申请提币数据
	tokenInout := new(TokenInout)
	err := tokenInout.GetByHash(txhash)
	if err != nil || tokenInout.Uid <= 0 {
		log.Error("get data by hash error",err,txhash)
		return
	}

	amount := decimal.New(tokenInout.Amount,0)
	fee := decimal.New(tokenInout.Fee,0)
	total := amount.Add(fee).IntPart()

	rsp,errr := client.InnerService.TokenSevice.CallConfirmSubFrozen(&proto.ConfirmSubFrozenRequest{
		Uid:uint64(tokenInout.Uid),
		TokenId:int32(tokenInout.Tokenid),
		Num:total,
		Ukey:[]byte(txhash),
		Type:proto.TOKEN_TYPE_OPERATOR_HISTORY_TOKEN_OUT,  //提币成功消耗冻结
	})
	log.WithFields(log.Fields{
		"uid":uint64(tokenInout.Uid),
		"token_uid":int32(tokenInout.Tokenid),
		"num":total,
		"ukey":txhash,
		"type":proto.TOKEN_TYPE_OPERATOR_HISTORY_TOKEN_OUT,
	}).Info("ETHConfirmSubFrozen result:",rsp,errr)
	if errr != nil {
		log.Error("ETHConfirmSubFrozen error",err)
	}
}

//添加比特币token数量
//**以太坊充币
func (p *Common) AddETHTokenNum(uid int,to string,tokenid int,amount string,txhash string) (bool,error) {
	fmt.Println("开始写到链上：",tokenid,",",uid,",",txhash)
	//查询用户uid
	//walletToken := new(models.WalletToken)
	//err := walletToken.GetByAddress(to)
	//if err != nil || walletToken.Uid <= 0 {
	//	log.Info("get user token error",err)
	//	return false,err
	//}

	bigAmount,err := strconv.ParseInt(amount,10,64)

	if err != nil {
		log.Error("AddETHTokenNum",err)
		return false,err
	}

	rsp,errr := client.InnerService.TokenSevice.CallAddTokenNum(&proto.AddTokenNumRequest{
		Uid:uint64(uid),
		TokenId:int32(tokenid),
		Opt:1,
		Num:bigAmount,
		Ukey:[]byte(txhash),
		OptAddType:1,
		Type:1,
	})
	log.WithFields(log.Fields{
		"Uid":uint64(uid),
		"TokenId":int32(tokenid),
		"Opt":1,
		"Num":bigAmount,
		"Ukey":[]byte(txhash),
		"OptAddType":1,
		"Type":1,
	}).Info("rpc------")
	log.Info("eth AddETHTokenNum result",errr,rsp)
	if errr != nil {
		log.Info("AddBTCTokenNum error",errr)
		return false,err
	}
	return true,nil
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

//添加比特币token数量
//**比特币充币
func (p *Common) AddUSDTTokenNum(data TranItem) {
	log.Info("AddBTCTokenNum data:",data)
	//查询用户uid
	walletToken := new(WalletToken)
	err := walletToken.GetByAddress(data.Address)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("get user token error",err)
		return
	}

	amount := convert.Float64ToInt64By8Bit(data.Amount)

	log.WithFields(log.Fields{
		"uid":walletToken.Uid,
		"token_id":walletToken.Tokenid,
		"num":amount,
		"ukey":data.Txid,
		"optAddType":0,
	}).Info("比特币RPC新增数量")

	//amount := int64(data.Amount * 100000000)
	rsp,errr := client.InnerService.TokenSevice.CallAddTokenNum(&proto.AddTokenNumRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Num:amount,
		Ukey:[]byte(data.Txid),
		OptAddType:1,
		Type:proto.TOKEN_TYPE_OPERATOR_HISTORY_HASH,
		Opt:proto.TOKEN_OPT_TYPE_ADD,
	})
	log.Info("btc AddBTCTokenNum result",errr,string(rsp.Message))
	if errr != nil {
		log.Info("AddBTCTokenNum error",errr)
	}
}

//确认消耗
//**比特币提币成功调用
func (p *Common) USDTConfirmSubFrozen(data TranItem) {
	//查询用户uid
	walletToken := new(WalletToken)
	err := walletToken.GetByAddress(data.Address)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("btc get user token error",err)
		return
	}

	//根据交易hash查询申请提币数据
	tokenInout := new(TokenInout)
	err = tokenInout.GetByHash(data.Txid)
	if err != nil || tokenInout.Uid <= 0 {
		log.Info("btc get data by hash error",err)
		return
	}

	defer func() {
		log.WithFields(log.Fields{
			"uid":walletToken.Uid,
			"token_id":walletToken.Tokenid,
			"num":tokenInout.Amount,
			"ukey":data.Txid,
			"type":proto.TOKEN_TYPE_OPERATOR_HISTORY_TOKEN_OUT,
		}).Info("比特币冻结数量")
	}()

	rsp,errr := client.InnerService.TokenSevice.CallConfirmSubFrozen(&proto.ConfirmSubFrozenRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Num:tokenInout.Amount,
		Ukey:[]byte(data.Txid),
		Type:proto.TOKEN_TYPE_OPERATOR_HISTORY_TOKEN_OUT,  //提币成功消耗冻结
	})
	log.Info("比特币确认消耗冻结：",errr,rsp.Err,string(rsp.Message))
	if errr != nil {
		log.Info("BTCConfirmSubFrozen error",err)
	}
}

//提币完成短信通知
func (p *Common) TiBiCompleteSendSms(apply_id int) () {
	tokenInout := new(TokenInout)
	err := tokenInout.GetByApplyId(apply_id)
	if err != nil {

	}
}
