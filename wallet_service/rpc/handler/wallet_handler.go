package handler

import (
	"context"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/wallet_service/model"
	"digicon/wallet_service/utils"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"digicon/wallet_service/watch"
	"digicon/wallet_service/rpc/client"
	log "github.com/sirupsen/logrus"
	"digicon/common/random"
	cf "digicon/wallet_service/conf"
	"digicon/common/convert"
)

type WalletHandler struct{}

func (s *WalletHandler) Hello(ctx context.Context, req *proto.HelloRequest2, rsp *proto.HelloResponse2) error {
	log.Print("Received Say.Hello request")
	rsp.Greeting = "Hello darnice mani" + req.Name
	return nil
}

func (s *WalletHandler) CreateWallet(ctx context.Context, req *proto.CreateWalletRequest, rsp *proto.CreateWalletResponse) error {
	log.Print("Received Say.CreateWallet request")

	var err error
	var addr string
	rsp.Data = new(proto.CreateWalletPos)

	//查询token状态
	tokenP := new(Tokens)
	b,e := tokenP.GetByid(int(req.Tokenid))
	if b != true || e != nil {
		rsp.Code = "1"
		rsp.Msg = err.Error()
		rsp.Data.Type = ""
		rsp.Data.Addr = ""
		return nil
	}
	if tokenP.Status == 2 {
		rsp.Code = "1"
		rsp.Msg = GetErrorMessage(ERRCODE_TOKEN_INVALID)
		rsp.Data.Type = ""
		rsp.Data.Addr = ""
		return nil
	}

	//判断钱包是否存在
	existsWalletToken := new(WalletToken)
	boo,address,signature := existsWalletToken.WalletTokenExist(int(req.Userid), int(req.Tokenid))
	if boo == true {
		rsp.Code = "0"
		rsp.Msg = address
		rsp.Data.Type = signature
		rsp.Data.Addr = address
		return nil
	}

	tokenModel := &Tokens{Id: int(req.Tokenid)}
	_, err = tokenModel.GetByid(int(req.Tokenid))

	//如果是以太坊类型的钱包，则直接查询，达到以太币和ERC20代币共用同一个地址
	if tokenModel.Signature == "eip155" || tokenModel.Signature == "eth" {
		walletToken := new(WalletToken)
		boo,err := walletToken.GetByTypeUid("eth",int(req.Userid))
		if boo == true && err == nil {
			//查询到了
			addr, err := walletToken.CopyEth(int(req.Userid), int(req.Tokenid), "123456", tokenModel.Chainid)
			if err != nil {
				rsp.Code = "1"
				rsp.Msg = err.Error()
				rsp.Data.Type = tokenModel.Signature
				rsp.Data.Addr = ""
				return nil
			}
			if addr == "" {
				rsp.Code = "1"
				rsp.Msg = GetErrorMessage(ERRCODE_CREATE_ERROR)
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
	}

	switch tokenModel.Signature {
	case "eip155", "eth":
		addr, err = Neweth(int(req.Userid), int(req.Tokenid), "123456", tokenModel.Chainid)
	case "btc":
		addr, err = NewBTC(int(req.Userid), int(req.Tokenid), "123456", tokenModel.Chainid)
	case "omni":
		addr, err = NewUSDT(int(req.Userid), int(req.Tokenid), "123456", tokenModel.Chainid)
	default:
		err = errors.New("unknow type")
	}

	if err != nil {
		rsp.Code = "1"
		rsp.Msg = err.Error()
		rsp.Data.Type = tokenModel.Signature
		rsp.Data.Addr = ""
		return nil
	}
	if addr == "" {
		rsp.Code = "1"
		rsp.Msg = GetErrorMessage(ERRCODE_CREATE_ERROR)
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

func (this *WalletHandler) Signtx(ctx context.Context, req *proto.SigntxRequest, rsp *proto.SigntxResponse) error {
	rsp.Data = new(proto.SigntxPos)

	//查询转出账号uid
	from_uid := cf.Cfg.MustInt("accounts","eth_uid",0)
	if from_uid == 0 {
		rsp.Code = "1"
		rsp.Msg = "转出账号未定义"
		return errors.New("转出账号未定义")
	}


	keystore := &WalletToken{Uid: from_uid, Tokenid: int(req.Tokenid)}

	ok, err := utils.Engine_wallet.Where("uid=? and tokenid=?", from_uid, int(req.Tokenid)).Get(keystore)

	if err != nil {
		rsp.Code = "1"
		rsp.Msg = err.Error()
		return nil
	}
	if !ok {
		rsp.Code = "1"
		rsp.Msg = "connot find keystore"
		return nil
	}
	deciml, err := new(Tokens).GetDecimal(int(req.Tokenid))
	if deciml == 0 || err != nil {
		rsp.Code = "1"
		rsp.Msg = "connot find tokens"
		return nil
	}
	deci_temp, err := decimal.NewFromString(req.Mount)
	if err != nil {
		rsp.Code = "1"
		rsp.Msg = "parse error"
		return nil
	}
	mount := deci_temp.Round(int32(deciml)).Coefficient()
	
	//获取随机数
	tokenData := new(Tokens)
	tokenData.GetByid(int(req.Tokenid))
	nonce,nonce_err := utils.RpcGetNonce(tokenData.Node,keystore.Address)

	if nonce_err != nil  {
		rsp.Code = "1"
		rsp.Msg = "get nonce error"
		return nil
	}

	//获取gasprice
	gasPrice,gasErr := utils.RpcGetGasPrice(tokenData.Node)
	if gasErr != nil  {
		rsp.Code = "1"
		rsp.Msg = "get gasprice error"
		return nil
	}
	
	signtxstr, err := keystore.Signtx(req.To, mount, gasPrice,nonce)

	if err != nil {
		rsp.Code = "1"
		rsp.Msg = err.Error()
		return nil
	}

	signtxStr := hex.EncodeToString(signtxstr)
	if signtxStr == "" {
		rsp.Code = "1"
		rsp.Msg = "signtxStr is null"
		return nil
	}

	rsp.Code = "0"
	rsp.Msg = GetErrorMessage(ERRCODE_SUCCESS)
	rsp.Data.Signtx = hex.EncodeToString(signtxstr)
	return nil
}

//提币失败，解冻用户冻结数据
func (this *WalletHandler) tiBiErrorUnfreeze(apply_id int32) {
	//根据申请id，查询数据
	tokenInout := new(TokenInout)
	err := tokenInout.GetByApplyId(int(apply_id))
	if err != nil {
		log.Error("根据申请id查询数据失败")
		return
	}
	//调用rpc解冻
	res,errr := client.InnerService.TokenSevice.CallCancelSubTokenWithFronze(&proto.CancelFronzeTokenRequest{
		Uid:uint64(tokenInout.Uid),
		TokenId:int32(tokenInout.Tokenid),
		Num:tokenInout.Amount + tokenInout.Fee,
		Ukey:[]byte(random.Random6dec()),
		Type:13,//取消提币
	})
	if errr != nil {
		log.Error("根据申请id解冻用户数据失败")
		return
	}
	log.Info("解冻成功：",res.Err,res.Message)
}

func (this *WalletHandler) SendRawTx(ctx context.Context, req *proto.SendRawTxRequest, rsp *proto.SendRawTxResponse) error {
	//log.Info("广播交易：",req.TokenId,req.Signtx)
	defer func() {
		if rsp.Code != "0" {
			log.WithFields(log.Fields{
				"signtx":req.Signtx,
				"token_id":req.TokenId,
				"code":rsp.Code,
				"msg":rsp.Msg,
			}).Error("SendRawTx error")
		}
		if rsp.Code != "0" {
			log.Error("广播失败，改回状态:",rsp.Msg)
			//把状态改回去
			//更新申请单记录
			_,err := new(TokenInout).UpdateApplyTiBi2(int(req.Applyid),4,rsp.Msg)  //正在提币中
			if err != nil {
				log.Error("UpdateApplyTiBi error:",err)
			}
			//取消冻结
			this.tiBiErrorUnfreeze(req.Applyid)
		}
	}()
	TokenModel := new(Tokens)
	ok, err := TokenModel.GetByid(int(req.TokenId))
	if err != nil || !ok {
		log.Error("token not find",req.TokenId)
		rsp.Code = "1"
		rsp.Msg = "token not find"
		return nil
	}

	rets, err := utils.RpcSendRawTx(TokenModel.Node, req.Signtx)
	if err != nil {
		log.Error("HTTP ERROR：",err,rets)
		rsp.Code = "1"
		rsp.Msg = err.Error()
		return nil
	}
	txhash, ok := rets["result"]
	if ok {
		//更新申请单记录
		_,err := new(TokenInout).UpdateApplyTiBi(int(req.Applyid),txhash.(string))
		if err != nil {
			log.Error("UpdateApplyTiBi error:",err)
		}
		//添加txhash到监控队列
		new(watch.EthTiBiWatch).PushRedisList(txhash.(string))
		rsp.Code = "0"
		rsp.Msg = GetErrorMessage(ERRCODE_SUCCESS)
		rsp.Data = new(proto.SendRawTxPos)
		rsp.Data.Result = []byte(txhash.(string))
		log.Info("广播交易成功：",rsp.Code,rsp.Msg,rsp.Data.Result)
		return nil
	}
	if !ok {
		log.Error("SendRawTx success：",rets,err)
		error := rets["error"].(map[string]interface{})
		rsp.Code = strconv.Itoa(int(error["code"].(float64)))
		rsp.Msg = error["message"].(string)
		return nil
	}

	return nil
}

func (this *WalletHandler) Tibi(ctx context.Context, req *proto.TibiRequest, rsp *proto.TibiResponse) error {
	rsp.Code = "0"
	rsp.Msg = GetErrorMessage(ERRCODE_SUCCESS)
	rsp.Data = new(proto.NilPos)
	return nil
}

func (this *WalletHandler) GetValue(ctx context.Context, req *proto.GetValueRequest, rsp *proto.GetValueResponse) error {
	rsp.Code = "0"
	rsp.Msg = "操作成功"
	rsp.Data = new(proto.GetValuePos)
	rsp.Data.Value = ""

	WalletTokenModel := &WalletToken{Uid: int(req.Uid), Tokenid: int(req.Tokenid)}
	ok, err := utils.Engine_wallet.Where("uid=? and tokenid=?", req.Uid, int(req.Tokenid)).Get(WalletTokenModel)
	if !ok || err != nil {
		rsp.Code = "1"
		rsp.Msg = "数据不存在"
		return nil
	}

	TokenModel := &Tokens{Id: int(req.Tokenid)}
	ok, err = utils.Engine_common.Where("id=?", int(req.Tokenid)).Get(TokenModel)
	if !ok || err != nil {
		rsp.Code = "1"
		rsp.Msg = "数据不存在"
		return nil
	}
	var value string
	var err2 error
	switch TokenModel.Signature {
	case "eip155", "eth":
		value, err2 = utils.RpcGetValue(TokenModel.Node, WalletTokenModel.Address, TokenModel.Contract, TokenModel.Decimal)
	default:
		value = ""
		err2 = errors.New("unavailable token type")
	}
	if err2 != nil {
		rsp.Code = "1"
		rsp.Msg = err2.Error()
		return nil
	}
	rsp.Data.Value = value
	return nil
}

//添加提币地址
func (this *WalletHandler) AddressSave(ctx context.Context, req *proto.AddressSaveRequest, rsp *proto.AddressSaveResponse) error {
	fmt.Println(req.Address)
	rsp.Code = "0"
	rsp.Msg = "添加成功"
	rsp.Data = new(proto.NilPos)
	TibiAddressModel := new(TibiAddress)
	_, err := TibiAddressModel.Save(int(req.Uid), int(req.Tokenid), req.Address, req.Mark)
	if err != nil {
		rsp.Code = "1"
		log.Info("add address save error!", err.Error())
		rsp.Msg = err.Error()
	}
	return nil
}

//提币地址列表
func (this *WalletHandler) AddressList(ctx context.Context, req *proto.AddressListRequest, rsp *proto.AddressListResponse) error {
	rsp.Code = "0"
	rsp.Msg = "列表成功"
	//rsp.Data = []AddrlistPos
	TibiAddressModel := new(TibiAddress)
	//rets, err := TibiAddressModel.List(int(req.Uid), int(req.Tokenid))
	rets, err := TibiAddressModel.List(int(req.Uid))
	if err != nil {
		return err
	}
	rsp.Data = rets
	//fmt.Println(rets)
	return nil
}

//提币地址删除
func (this *WalletHandler) AddressDelete(ctx context.Context, req *proto.AddressDeleteRequest, rsp *proto.AddressDeleteResponse) error {
	TibiAddressModel := new(TibiAddress)
	_, err := TibiAddressModel.DeleteByid(int(req.Id), int(req.Uid))
	if err != nil {
		rsp.Code = "1"
		rsp.Msg = err.Error()
		rsp.Data = new(proto.NilPos)
		return nil
	}
	rsp.Code = "0"
	rsp.Msg = "操作成功"
	rsp.Data = new(proto.NilPos)
	return nil
}

func (this *WalletHandler) InList(ctx context.Context, req *proto.InListRequest, rsp *proto.InListResponse) error {
	filter := map[string]interface{}{
		"uid": req.Uid,
		"opt": 1,
	}

	tokenInoutMD := new(TokenInout)
	modelList, list, err := tokenInoutMD.GetInOutList(int(req.Page), int(req.PageNum), filter)
	if err != nil {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = err.Error()
		return nil
	}

	// 拼接返回数据
	rsp.Data = new(proto.InListResponse_Data)
	rsp.Data.Items = make([]*proto.InListResponse_Data_Item, 0)

	rsp.Data.PageIndex = int32(modelList.PageIndex)
	rsp.Data.PageSize = int32(modelList.PageSize)
	rsp.Data.TotalPage = int32(modelList.PageCount)
	rsp.Data.Total = int32(modelList.Total)

	for _, v := range list {
		rsp.Data.Items = append(rsp.Data.Items, &proto.InListResponse_Data_Item{
			Id:          int32(v.Id),
			TokenId:     int32(v.Tokenid),
			TokenName:   v.TokenName,
			Amount:      v.Amount,
			Address:     v.To,
			States:      int32(v.States),
			CreatedTime: v.CreatedTime.Unix(),
		})
	}

	return nil
}

func (this *WalletHandler) OutList(ctx context.Context, req *proto.OutListRequest, rsp *proto.OutListResponse) error {
	filter := map[string]interface{}{
		"uid": req.Uid,
		"opt": 2,
	}

	tokenInoutMD := new(TokenInout)
	modelList, list, err := tokenInoutMD.GetInOutList(int(req.Page), int(req.PageNum), filter)
	if err != nil {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = err.Error()
		return nil
	}

	// 拼接返回数据
	rsp.Data = new(proto.OutListResponse_Data)
	rsp.Data.Items = make([]*proto.OutListResponse_Data_Item, 0)

	rsp.Data.PageIndex = int32(modelList.PageIndex)
	rsp.Data.PageSize = int32(modelList.PageSize)
	rsp.Data.TotalPage = int32(modelList.PageCount)
	rsp.Data.Total = int32(modelList.Total)

	for _, v := range list {
		rsp.Data.Items = append(rsp.Data.Items, &proto.OutListResponse_Data_Item{
			Id:          int32(v.Id),
			TokenId:     int32(v.Tokenid),
			TokenName:   v.TokenName,
			Amount:      v.Amount,
			Fee:         v.Fee,
			Remarks:     v.Remarks,
			Address:     v.To,
			States:      int32(v.States),
			CreatedTime: v.CreatedTime.Unix(),
		})
	}

	return nil
}

func (this *WalletHandler) TibiApply(ctx context.Context, req *proto.TibiApplyRequest, rsp *proto.TibiApplyResponse) error {
	tokenInoutMD := new(TokenInout)
	//验证短信验证码
	ret,err := AuthSms(req.Phone,13,req.SmsCode)
	if ret != ERRCODE_SUCCESS {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = GetErrorMessage(ret)  //"短信验证码错误"
		return errors.New(GetErrorMessage(ret))
	}

	////验证邮箱验证码
	ret,err = AuthEmail(req.Email,13,req.EmailCode)
	if ret != ERRCODE_SUCCESS {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = GetErrorMessage(ret)  //"邮箱验证码错误"
		return errors.New(GetErrorMessage(ret))
	}

	//验证资金密码
	ret,err = tokenInoutMD.AuthPayPwd(req.Uid,req.Password)
	if ret != ERRCODE_SUCCESS {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = GetErrorMessage(ERRCODE_PAY_PWD)
		return errors.New(GetErrorMessage(ERRCODE_PAY_PWD))
	}

	//检查资金是否足够
	userToken := new(UserToken)
	boo,err := userToken.GetByUidTokenid(int(req.Uid),int(req.Tokenid))
	if boo == false || err != nil {
		log.Error(GetErrorMessage(ERRCODE_TOKEN_NOT_ENOUGH),boo,err,req.Uid,req.Tokenid)
		rsp.Code = ERRCODE_TOKEN_NOT_ENOUGH
		rsp.Msg = GetErrorMessage(ERRCODE_TOKEN_NOT_ENOUGH)
		return errors.New(GetErrorMessage(ERRCODE_TOKEN_NOT_ENOUGH))
	}

	amount,err1 := convert.StringToInt64By8Bit(req.Amount)
	balance,err2 := strconv.ParseInt(userToken.Balance,10,64)
	if err1 != nil || err2 != nil {
		log.Error("格式化错误：",req.Amount,userToken.Balance,err,err1)
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = "格式化错误"
		return errors.New("格式化错误")
	}
	if balance < amount {
		log.Error("余额不足：",balance,amount)
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = GetErrorMessage(ERRCODE_TOKEN_NOT_ENOUGH)
		return errors.New(GetErrorMessage(ERRCODE_TOKEN_NOT_ENOUGH))
	}

	var amountCny int64
	var feeCny int64
	//查询币种和人民币汇率
	err,cnyPriceInt := utils.GetCnyPrice(int(req.Tokenid))
	if err == nil || cnyPriceInt > 0 {
		//计算折合人民币
		a,err := strconv.ParseFloat(req.RealAmount,10)
		if err != nil {
			log.Error(GetErrorMessage(ERRCODE_PARSE),err)
			rsp.Code = ERRCODE_UNKNOWN
			rsp.Msg = GetErrorMessage(ERRCODE_PARSE)
			return errors.New(GetErrorMessage(ERRCODE_PARSE))
		}
		t1 := decimal.NewFromFloat(a)
		t1_c := decimal.NewFromFloat(float64(cnyPriceInt))
		amountCny = t1.Mul(t1_c).IntPart()

		b,err := strconv.ParseFloat(req.Gasprice,10)
		if err != nil {
			log.Error(GetErrorMessage(ERRCODE_PARSE),err)
			rsp.Code = ERRCODE_UNKNOWN
			rsp.Msg = GetErrorMessage(ERRCODE_PARSE)
			return errors.New(GetErrorMessage(ERRCODE_PARSE))
		}
		t2 := decimal.NewFromFloat(b)
		t2_c := decimal.NewFromFloat(float64(cnyPriceInt))
		feeCny = t2.Mul(t2_c).IntPart()
	} else {
		rsp.Code = ERRCODE_UNKNOWN
		log.Error(GetErrorMessage(ERRCODE_CNY_PRICE)+err.Error())
		rsp.Msg = GetErrorMessage(ERRCODE_CNY_PRICE)
		return errors.New(GetErrorMessage(ERRCODE_CNY_PRICE))
	}


	//先冻结资金


	a,err := strconv.ParseFloat(req.Amount,10)
	if err != nil {
		log.Error(GetErrorMessage(ERRCODE_FORMAT),req.Amount,err)
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = GetErrorMessage(ERRCODE_FORMAT)
		return errors.New(GetErrorMessage(ERRCODE_FORMAT))
	}
	t1 := decimal.NewFromFloat(a)
	t1_c := decimal.NewFromFloat(float64(100000000))
	fee := t1.Mul(t1_c).IntPart()

	c,rErr := client.InnerService.TokenSevice.CallSubTokenWithFronze(&proto.SubTokenWithFronzeRequest{
		Uid:uint64(req.Uid),
		TokenId:req.Tokenid,
		Num:fee,
		Opt:1, //卖
		Ukey:[]byte(random.Random6dec()),
		Type:12,  //提币
	})
	log.Info("资金冻结结果：",rErr,req.Uid,fee,c)
	if rErr != nil {
		log.Error(GetErrorMessage(ERRCODE_FREEZE),rErr)
		rsp.Code = 1
		rsp.Msg = GetErrorMessage(ERRCODE_FREEZE)
		return errors.New(GetErrorMessage(ERRCODE_FREEZE))
	}

	//查询代币类型
	tokens := new(Tokens)
	boo,err = tokens.GetByid(int(req.Tokenid))
	if boo != true || err != nil {
		log.Error("查询token错误：",boo,err,req.Tokenid)
		rsp.Code = 1
		rsp.Msg = "查询token错误"
		return errors.New("查询token错误")
	}

	fromAddress := ""
	if tokens.Signature == "eip" || tokens.Signature == "eip155" {
		//查询配置的提币地址
		fromAddress = cf.Cfg.MustValue("accounts","eth_address","")
		if fromAddress == "" {
			rsp.Code = 1
			rsp.Msg = GetErrorMessage(ERRCODE_TIBI_ADDRESS)
			return errors.New(GetErrorMessage(ERRCODE_TIBI_ADDRESS))
		}
	}

	//usdt
	if tokens.Signature == "omni" {
		//查询配置的提币地址
		fromAddress = cf.Cfg.MustValue("accounts","usdt_address","")
		if fromAddress == "" {
			rsp.Code = 1
			rsp.Msg = GetErrorMessage(ERRCODE_TIBI_ADDRESS)
			return errors.New(GetErrorMessage(ERRCODE_TIBI_ADDRESS))
		}
	}

	//保存数据
	_,err = tokenInoutMD.TiBiApply(int(req.Uid),int(req.Tokenid),req.To,req.RealAmount,req.Gasprice,amountCny,feeCny,fromAddress)
	if err != nil {
		log.Error(err.Error())
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = err.Error()
		return err
	}

	log.Info("提币完成")

	rsp.Code = 0
	rsp.Msg = GetErrorMessage(ERRCODE_SUCCESS)
	return nil
}

func (s *WalletHandler) GetAddress(ctx context.Context, req *proto.GetAddressRequest, rsp *proto.GetAddressResponse) error {
	log.Print("Received Say.CreateWallet request")
	fmt.Println(req.String())
	walletToken := new(WalletToken)
	err := walletToken.GetByUidTokenid(int(req.Userid), int(req.Tokenid))
	if err != nil {
		rsp.Code = "1"
		rsp.Msg = err.Error()
		rsp.Addr = ""
		return nil
	}
	rsp.Code = "0"
	rsp.Msg = "成功"
	rsp.Addr = walletToken.Address
	rsp.Type = walletToken.Type
	return nil
}

func (this *WalletHandler) CancelTiBi(ctx context.Context, req *proto.CancelTiBiRequest, rsp *proto.CancelTiBiResponse) error {

	defer func() {
		if rsp.Code != ERRCODE_SUCCESS {
			log.WithFields(log.Fields{
				"uid":req.Uid,
				"id":req.Id,
				"code":rsp.Code,
				"msg":rsp.Msg,
			}).Error("CancelTiBi error")
		}
	}()

	//解除冻结账户金额，查询
	tokenInout := new(TokenInout)
	boo,err := tokenInout.GetApplyInOut(int(req.Uid),int(req.Id))
	if boo != true || err != nil {
		rsp.Code = ERRCODE_SELECT_ERROR
		rsp.Msg = GetErrorMessage(ERRCODE_SELECT_ERROR)
		return nil
	}

	//调用rpc解冻
	ukey := strconv.Itoa(int(req.Uid)) + random.Random6dec()
	res,errr := client.InnerService.TokenSevice.CallCancelSubTokenWithFronze(&proto.CancelFronzeTokenRequest{
		Uid:uint64(req.Uid),
		TokenId:int32(tokenInout.Tokenid),
		Num:tokenInout.Amount + tokenInout.Fee,
		Ukey:[]byte(ukey),
		Type:proto.TOKEN_TYPE_OPERATOR_TRANSFER_FROM_CANCELTIBI,//取消提币
	})
	if errr != nil {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = "RPC ERROR"
		return nil
	}
	if res.Err != 0 {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = res.Message
		return nil
	}

	//修改状态
	tokenInoutMD := new(TokenInout)
	//保存数据
	_,err = tokenInoutMD.CancelTiBi(int(req.Uid),int(req.Id))
	if err != nil {
		rsp.Code = ERRCODE_SAVE_ERROR
		rsp.Msg = GetErrorMessage(ERRCODE_SAVE_ERROR)
		return nil
	}

	rsp.Code = ERRCODE_SUCCESS
	rsp.Msg = GetErrorMessage(ERRCODE_SUCCESS)

	return nil
}

//同步以太坊区块交易
func (this *WalletHandler) SyncEthBlockTx(ctx context.Context, req *proto.SyncEthBlockTxRequest, rsp *proto.SyncEthBlockTxResponse) error {

	//判断key是否正确
	key := cf.Cfg.MustValue("keys","sync_block","")
	if key == "" {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = "KEY未配置"
		return nil
	}
	if req.Key != key {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = "key error"
		return nil
	}

	err,msg := new(watch.EthOperate).WorkerHander(int(req.Block))
	if err != nil {
		rsp.Code = 1
		rsp.Msg = msg
		return err
	}
	rsp.Code = 0
	rsp.Msg = msg
	return nil
}

//获取提币手续费
func (this *WalletHandler) GetOutTokenFee(ctx context.Context, req *proto.GetOutTokenFeeRequest, rsp *proto.GetOutTokenFeeResponse) error {
	token := new(Tokens)
	data,err := token.GetAllTokenFee()
	if err != nil {
		rsp.Code = 1
		rsp.Msg = err.Error()
		return err
	}

	for _,v := range data {
		rsp.Data = append(rsp.Data,&proto.GetOutTokenFeeResponseList{
			Tokenid:int32(v.Id),
			Fee:strconv.FormatFloat(v.Out_token_fee,'f',-1,64),
		})
	}

	rsp.Code = 0
	rsp.Msg = GetErrorMessage(ERRCODE_SUCCESS)
	return nil
}

//解冻用户数据
func (this *WalletHandler) CancelSubTokenWithFronze(ctx context.Context, req *proto.CancelSubTokenWithFronzeRequest, rsp *proto.CancelSubTokenWithFronzeResponse) error {
	//判断key是否正确
	key := cf.Cfg.MustValue("keys","cancel_fronze","")
	if key == "" {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = "KEY未配置"
		return nil
	}
	if req.Key != key {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = "key error"
		return nil
	}
	//调用rpc解冻
	res,errr := client.InnerService.TokenSevice.CallCancelSubTokenWithFronze(&proto.CancelFronzeTokenRequest{
		Uid:uint64(req.Uid),
		TokenId:int32(req.Tokenid),
		Num:req.Num,
		Ukey:[]byte(req.Ukey),
		Type:proto.TOKEN_TYPE_OPERATOR_TRANSFER_FROM_CANCELTIBI,//取消提币
	})
	if errr != nil {
		log.Error("RPC ERROR")
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = "RPC ERROR"
		return nil
	}
	if res.Err != 0 {
		log.Error(res.Message)
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = res.Message
		return nil
	}
	rsp.Code = ERRCODE_SUCCESS
	rsp.Msg = GetErrorMessage(ERRCODE_SUCCESS)
	return nil
}
