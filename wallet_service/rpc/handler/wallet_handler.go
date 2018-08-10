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
	"log"
	"strconv"
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

	fmt.Println(req.String())
	tokenModel := &Tokens{Id: int(req.Tokenid)}
	_, err = tokenModel.GetByid(int(req.Tokenid))

	switch tokenModel.Signature {
	case "eip155", "eth":
		addr, err = Neweth(int(req.Userid), int(req.Tokenid), "123456", tokenModel.Chainid)
	case "btc":
		//fmt.Println(tokenModel.Chainid)
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
	fmt.Println("获取gasprice:",gasErr,gasPrice)
	if gasErr != nil  {
		rsp.Code = "1"
		rsp.Msg = "获取gasprice失败"
		return nil
	}
	
	signtxstr, err := keystore.Signtx(req.To, mount, gasPrice,nonce)

	fmt.Println("生成的签名----------------：", gasPrice,nonce,err,signtxstr)

	if err != nil {
		rsp.Code = "1"
		rsp.Msg = err.Error()
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

	fmt.Println("发送交易：",TokenModel.Node,req.Signtx)

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
		fmt.Println("add address save error!", err.Error())
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
	//ret,err := AuthSms(req.Phone,SMS_CARRY_COIN,req.SmsCode)
	//if ret != ERRCODE_SUCCESS {
	//	rsp.Code = ERRCODE_UNKNOWN
	//	rsp.Msg = err.Error()
	//	return nil
	//}
	////验证邮箱验证码
	//ret,err = AuthEmail(req.Email,SMS_CARRY_COIN,req.EmailCode)
	//if ret != ERRCODE_SUCCESS {
	//	rsp.Code = ERRCODE_UNKNOWN
	//	rsp.Msg = err.Error()
	//	return nil
	//}
	//验证资金密码
	//ret,err := tokenInoutMD.AuthPayPwd(req.Uid,req.Password)
	//if ret != ERRCODE_SUCCESS {
	//	rsp.Code = ERRCODE_UNKNOWN
	//	rsp.Msg = err.Error()
	//	return nil
	//}
	//保存数据
	_,err := tokenInoutMD.TiBiApply(int(req.Uid),int(req.Tokenid),req.To,req.RealAmount,req.Gasprice)
	if err != nil {
		rsp.Code = ERRCODE_UNKNOWN
		rsp.Msg = err.Error()
		return nil
	}
	//冻结账户金额

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
	//解除冻结账户金额

	rsp.Code = ERRCODE_SUCCESS
	rsp.Msg = "修改成功"
	return nil
}
