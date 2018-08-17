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
	"math/big"
	"digicon/wallet_service/watch"
	"digicon/wallet_service/rpc/client"
	log "github.com/sirupsen/logrus"
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
		rsp.Msg = "Token暂不可用"
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

	switch tokenModel.Signature {
	case "eip155", "eth":
		addr, err = Neweth(int(req.Userid), int(req.Tokenid), "123456", tokenModel.Chainid)
	case "btc":
		addr, err = NewBTC(int(req.Userid), int(req.Tokenid), "123456", tokenModel.Chainid)
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
		rsp.Msg = "创建失败"
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
	log.Print("Received Say.Signtx request")
	rsp.Data = new(proto.SigntxPos)

	keystore := &WalletToken{Uid: int(req.Userid), Tokenid: int(req.Tokenid)}

	ok, err := utils.Engine_wallet.Where("uid=? and tokenid=?", req.Userid, int(req.Tokenid)).Get(keystore)

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
		rsp.Msg = "获取随机数失败"
		return nil
	}

	//获取gasprice
	gasPrice,gasErr := utils.RpcGetGasPrice(tokenData.Node)
	log.Info("获取gasprice:",gasErr,gasPrice)
	if gasErr != nil  {
		rsp.Code = "1"
		rsp.Msg = "获取gasprice失败"
		return nil
	}
	
	signtxstr, err := keystore.Signtx(req.To, mount, gasPrice,nonce)

	if err != nil {
		rsp.Code = "1"
		rsp.Msg = err.Error()
		return nil
	}

	signtxStr := hex.EncodeToString(signtxstr)

	log.Info("生成结果：",signtxStr)
	if signtxStr == "" {
		rsp.Code = "1"
		rsp.Msg = "签名字符串为空，请稍后重试"
		return nil
	}

	rsp.Code = "0"
	rsp.Msg = "生成成功"
	//rsp.Data = new(proto.SigntxPos)
	rsp.Data.Signtx = hex.EncodeToString(signtxstr)
	return nil
}

func (this *WalletHandler) SendRawTx(ctx context.Context, req *proto.SendRawTxRequest, rsp *proto.SendRawTxResponse) error {
	log.Print("Received Say.SendRawTx request")
	TokenModel := new(Tokens)
	ok, err := TokenModel.GetByid(int(req.TokenId))
	if err != nil || !ok {
		rsp.Code = "1"
		rsp.Msg = "token不存在"
		return nil
	}

	rets, err := utils.RpcSendRawTx(TokenModel.Node, req.Signtx)
	if err != nil {
		rsp.Code = "1"
		rsp.Msg = err.Error()
		return nil
	}
	txhash, ok := rets["result"]
	if ok {
		//更新申请单记录
		new(TokenInout).UpdateApplyTiBi(int(req.Applyid),txhash.(string))
		//添加txhash到监控队列
		new(watch.EthTiBiWatch).PushRedisList(txhash.(string))
		rsp.Code = "0"
		rsp.Msg = "发送成功"
		rsp.Data = new(proto.SendRawTxPos)
		rsp.Data.Result = txhash.(string)
		return nil
	}
	if !ok {
		error := rets["error"].(map[string]interface{})
		rsp.Code = strconv.Itoa(int(error["code"].(float64)))
		rsp.Msg = error["message"].(string)
		return nil
	}
	return nil
}

func (this *WalletHandler) Tibi(ctx context.Context, req *proto.TibiRequest, rsp *proto.TibiResponse) error {
	rsp.Code = "0"
	rsp.Msg = "生成成功"
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
	rsp.Msg = "添加1成功"
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
		return nil
	}
	////验证邮箱验证码
	ret,err = AuthEmail(req.Email,13,req.EmailCode)
	if ret != ERRCODE_SUCCESS {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = GetErrorMessage(ret)  //"邮箱验证码错误"
		return nil
	}
	//验证资金密码
	ret,err = tokenInoutMD.AuthPayPwd(req.Uid,req.Password)
	if ret != ERRCODE_SUCCESS {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = "支付密码错误"
		return nil
	}

	//先冻结资金
	tmp1,_ := new(big.Int).SetString(req.Amount,10)
	fee1 := decimal.NewFromBigInt(tmp1, int32(8)).IntPart()
	c,rErr := client.InnerService.TokenSevice.CallSubTokenWithFronze(&proto.SubTokenWithFronzeRequest{
		Uid:uint64(req.Uid),
		TokenId:req.Tokenid,
		Num:fee1,
		Opt:1, //卖
		Ukey:[]byte{byte(req.Uid)},
		Type:12,  //提币
	})
	fmt.Println("资金冻结结果：",rErr,req.Uid,fee1,c)
	if rErr != nil {
		log.Error(rErr.Error())
		rsp.Code = 1
		rsp.Msg = "冻结资金失败"
		return nil
	}

	var amountCny int64
	var feeCny int64
	//查询币种和人民币汇率
	err,cnyPriceInt := utils.GetCnyPrice(int(req.Tokenid))
	if err == nil || cnyPriceInt > 0 {
		//计算折合人民币
		a,err := strconv.ParseFloat(req.RealAmount,10)
		if err != nil {
			log.Error(err)
			rsp.Code = ERRCODE_UNKNOWN
			rsp.Msg = "解析失败"+err.Error()
			return nil
		}
		t1 := decimal.NewFromFloat(a)
		t1_c := decimal.NewFromFloat(float64(cnyPriceInt))
		amountCny = t1.Mul(t1_c).IntPart()

		b,err := strconv.ParseFloat(req.Gasprice,10)
		if err != nil {
			log.Error(err)
			rsp.Code = ERRCODE_UNKNOWN
			rsp.Msg = "解析失败"+err.Error()
			return nil
		}
		t2 := decimal.NewFromFloat(b)
		t2_c := decimal.NewFromFloat(float64(cnyPriceInt))
		feeCny = t2.Mul(t2_c).IntPart()
	} else {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = "获取人民币价格出错"+err.Error()
		return nil
	}


	//保存数据
	_,err = tokenInoutMD.TiBiApply(int(req.Uid),int(req.Tokenid),req.To,req.RealAmount,req.Gasprice,amountCny,feeCny)
	if err != nil {
		log.Error(err.Error())
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = err.Error()
		return nil
	}

	log.Info("提币完成")

	rsp.Code = 0
	rsp.Msg = "保存成功"
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
	tokenInoutMD := new(TokenInout)
	//保存数据
	_,err := tokenInoutMD.CancelTiBi(int(req.Uid),int(req.Id))
	if err != nil {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = err.Error()
		return nil
	}
	//解除冻结账户金额，查询
	tokenInout := new(TokenInout)
	boo,err := tokenInout.GetApplyInOut(int(req.Uid),int(req.Id))
	if boo != true || err != nil {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = "查询失败"
		return nil
	}
	//调用rpc解冻
	_,errr := client.InnerService.TokenSevice.CallCancelSubTokenWithFronze(&proto.CancelFronzeTokenRequest{
		Uid:uint64(req.Uid),
		TokenId:int32(tokenInout.Tokenid),
		Num:tokenInout.Amount + tokenInout.Fee,
	})
	if errr != nil {
		log.WithFields(log.Fields{
			"uid":req.Uid,
			"tokenid":tokenInout.Tokenid,
			"amount":tokenInout.Amount,
			"fee":tokenInout.Fee,
		}).Error("解冻失败")
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = "解冻失败"
		return nil
	}
	rsp.Code = ERRCODE_SUCCESS
	rsp.Msg = "解冻成功"
	return nil
}

//同步以太坊区块交易
func (this *WalletHandler) SyncEthBlockTx(ctx context.Context, req *proto.SyncEthBlockTxRequest, rsp *proto.SyncEthBlockTxResponse) error {
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
