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
	"errors"
	"strings"
	"time"
	cf "digicon/wallet_service/conf"
)

type Common struct{}

//添加比特币token数量
//**比特币充币
func (p *Common) AddBTCTokenNum(data TranItem) error {
	log.Info("AddBTCTokenNum data:",data)
	//查询用户uid
	walletToken := new(WalletToken)
	err := walletToken.GetByAddress(data.Address)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("get user token error",err)
		return err
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
		return errr
	}
	return nil
}

//确认消耗
//**比特币提币成功调用
func (p *Common) BTCConfirmSubFrozen(data TranItem) error {
	//根据address获取提币地址所属的用户uid,对于提币来说，address是对方的地址
	tibiAddress := new(TibiAddress)
	boo,err := tibiAddress.GetByAddress(data.Address)
	if boo != true || err != nil {
		log.Error("根据地址获取提币地址失败：",boo,err,data.Address,data)
		return errors.New("根据地址获取提币地址失败")
	}
	//查询用户uid
	walletToken := new(WalletToken)
	err = walletToken.GetByUid(tibiAddress.Uid)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("btc get user token error",err)
		return err
	}

	//根据交易hash查询申请提币数据
	tokenInout := new(TokenInout)
	err = tokenInout.GetByHash(data.Txid)
	if err != nil || tokenInout.Uid <= 0 {
		log.Info("btc get data by hash error",err)
		return err
	}

	amount := decimal.New(tokenInout.Amount,0)
	fee := decimal.New(tokenInout.Fee,0)
	total := amount.Add(fee).IntPart()


	defer func() {
		log.WithFields(log.Fields{
			"uid":walletToken.Uid,
			"token_id":walletToken.Tokenid,
			"num":total,
			"ukey":data.Txid,
			"type":proto.TOKEN_TYPE_OPERATOR_HISTORY_TOKEN_OUT,
		}).Info("比特币冻结数量")
	}()

	rsp,errr := client.InnerService.TokenSevice.CallConfirmSubFrozen(&proto.ConfirmSubFrozenRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Num:total,
		Ukey:[]byte(data.Txid),
		Type:proto.TOKEN_TYPE_OPERATOR_HISTORY_TOKEN_OUT,  //提币成功消耗冻结
	})
	log.Info("比特币确认消耗冻结：",errr,rsp.Err,string(rsp.Message))
	if errr != nil {
		log.Info("BTCConfirmSubFrozen error",err)
		return errr
	}
	return nil
}

//确认消耗
//**以太坊提币成功调用
func (p *Common) ETHConfirmSubFrozen(from string,txhash string,contract string) error {
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
		return err
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
		return errr
	}
	return nil
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
func (p *Common) AddUSDTTokenNum(data USDTTranInfo) {
	log.Info("AddUSDTTokenNum data:",data)
	//查询用户uid
	walletToken := new(WalletToken)
	err := walletToken.GetByAddress(data.Sendingaddress)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("get user token error",err)
		return
	}

	amount,err := convert.StringToInt64By8Bit(data.Amount)
	if err != nil {
		log.Error("StringToInt64By8Bit error",err)
	}

	log.WithFields(log.Fields{
		"uid":walletToken.Uid,
		"token_id":walletToken.Tokenid,
		"num":amount,
		"ukey":data.Txid,
		"optAddType":0,
	}).Info("USDT RPC新增数量")

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
	log.Info("btc AddUSDTTokenNum result",errr,string(rsp.Message))
	if errr != nil {
		log.Info("AddUSDTTokenNum error",errr)
	}
}

//确认消耗
//**比特币提币成功调用
func (p *Common) USDTConfirmSubFrozen(data USDTTranInfo) error {
	//查询用户uid
	walletToken := new(WalletToken)
	err := walletToken.GetByAddress(data.Sendingaddress)
	if err != nil || walletToken.Uid <= 0 {
		log.Info("btc get user token error",err)
		return err
	}

	//根据交易hash查询申请提币数据
	tokenInout := new(TokenInout)
	err = tokenInout.GetByHash(data.Txid)
	if err != nil || tokenInout.Uid <= 0 {
		log.Info("usdt get data by hash error",err)
		return errors.New("usdt get data by hash error")
	}

	amount := decimal.New(tokenInout.Amount,0)
	fee := decimal.New(tokenInout.Fee,0)
	total := amount.Add(fee).IntPart()

	defer func() {
		log.WithFields(log.Fields{
			"uid":walletToken.Uid,
			"token_id":walletToken.Tokenid,
			"num":total,
			"ukey":data.Txid,
			"type":proto.TOKEN_TYPE_OPERATOR_HISTORY_TOKEN_OUT,
		}).Info("USDT冻结数量")
	}()

	rsp,errr := client.InnerService.TokenSevice.CallConfirmSubFrozen(&proto.ConfirmSubFrozenRequest{
		Uid:uint64(walletToken.Uid),
		TokenId:int32(walletToken.Tokenid),
		Num:total,
		Ukey:[]byte(data.Txid),
		Type:proto.TOKEN_TYPE_OPERATOR_HISTORY_TOKEN_OUT,  //提币成功消耗冻结
	})
	log.Info("USDT确认消耗冻结：",errr,rsp.Err,string(rsp.Message))
	if errr != nil {
		log.Info("USDTConfirmSubFrozen error",err)
		return errr
	}
	return nil
}

//提币完成短信通知
func (p *Common) TiBiCompleteSendSms(apply_id int) (err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"apply_id":apply_id,
				"err":err,
			}).Error("TiBiCompleteSendSms error")
		}
		fmt.Println("结果：",err,apply_id)
	}()

	var boo bool

	tokenInout := new(TokenInout)
	err = tokenInout.GetByApplyId(apply_id)
	if err != nil {
		return
	}
	tokens := new(Tokens)
	boo,err = tokens.GetByid(tokenInout.Tokenid)
	if err != nil {
		return
	}
	if boo != true {
		err = errors.New("token not found!")
		return
	}

	user := new(User)
	boo,err = user.GetUser(uint64(tokenInout.Uid))
	if err != nil {
		return
	}
	if boo != true {
		err = errors.New("用户数据为空")
		return
	}

	phone := user.Phone
	mark := tokens.Mark
	num := convert.Int64ToStringBy8Bit(tokenInout.Amount)
	content := strings.Join([]string{"你申请的提币已经完成，币种：",mark,"，到账数量：",num},"")
	_,err = SendInterSms(phone,content)

	fmt.Println(phone,err)

	log.Info("TiBiCompleteSendSms complete")
	return
}

//测试发短信
func TestSms() {
	common := new(Common)
	res := common.TiBiCompleteSendSms(228)
	fmt.Println("发送结果：",res)
}

//汇总提币手续费
func (p *Common) GatherFee(txhash string) (err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"err":err,
			}).Error("汇总数据失败")
		}
	}()
	if txhash == "" {
		return
	}
	tokenInout := new(TokenInout)
	err = tokenInout.GetByHash(txhash)
	if err != nil {
		return
	}
	if tokenInout.Opt == 1 {
		//充币
		return
	}
	//fee := decimal.New(tokenInout.Fee,0)
	//realFee := decimal.New(tokenInout.Real_fee,0)
	//amount := fee.Sub(realFee).IntPart()

	strTokenid := strconv.Itoa(tokenInout.Tokenid)

	//写一条代币手续费
	tokensFreeHistory := new(Token_free_history)
	tokensFreeHistory.Opt = 1
	tokensFreeHistory.Token_id = int64(tokenInout.Tokenid)
	tokensFreeHistory.Type = int64(proto.TOKEN_TYPE_OPERATOR_HISTORY_FEE)
	tokensFreeHistory.Num = tokenInout.Fee
	tokensFreeHistory.Created_time = time.Now().Unix()
	tokensFreeHistory.Ukey = strTokenid + "_" + tokenInout.Txhash
	tokensFreeHistory.Uid = int64(tokenInout.Uid)
	_,err = tokensFreeHistory.InsertThis()
	if err != nil {
		return
	}


	tokenModel := new(Tokens)
	exists, err := tokenModel.GetByName("ETH")
	if err != nil {
		log.Info("tokenModel error",err)
	}
	if !exists {
		log.Info("token not exists eth ...")
	}

	boo,from_uid := p.GetEthAccountUid()
	if boo != true {
		log.Error("获取官方转出uid错误")
		return
	}

	//在写一条以太币出账记录
	tokensEthFreeHistory := new(Token_free_history)
	tokensEthFreeHistory.Opt = 2
	tokensEthFreeHistory.Token_id = int64(tokenModel.Id)
	tokensEthFreeHistory.Type = int64(proto.TOKEN_TYPE_OPERATOR_HISTORY_FEE)
	tokensEthFreeHistory.Num = tokenInout.Real_fee
	tokensEthFreeHistory.Created_time = time.Now().Unix()
	tokensEthFreeHistory.Ukey = strconv.Itoa(tokenModel.Id) + "_" + tokenInout.Txhash
	tokensEthFreeHistory.Uid = from_uid
	_,err = tokensEthFreeHistory.InsertThis()
	if err != nil {
		return
	}

	return
}

//汇总历史手续费，只能调用一次
func GatherHistoryFee() {
	ok := cf.Cfg.MustInt("gather", "history_fee",0)
	if ok != 1 {
		log.Error("报错了")
		return
	}
	tokenInout := new(TokenInout)
	data,err := tokenInout.GetHashs(2)  //表示提币
	if err != nil {
		log.Error("获取hash报错",err)
		return
	}
	p := new(Common)
	for _,v := range data {
		if v.Txhash == "" {
			continue
		}
		if v.Opt != 2 {
			continue
		}
		err := p.GatherFee(v.Txhash)
		if err != nil {
			log.Error("汇总数据出错：",err)
		}
	}
}

//查询以太坊官方账户地址
func (p *Common) GetEthAccountUid() (bool,int64) {
	from_uid := cf.Cfg.MustInt64("accounts","eth_uid",0)
	if from_uid == 0 {
		return false,0
	}
	return true,from_uid
}
