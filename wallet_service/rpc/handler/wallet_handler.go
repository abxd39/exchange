package handler

import (
	"context"
	proto "digicon/proto/rpc"
	. "digicon/wallet_service/model"
	"digicon/wallet_service/utils"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
)

type WalletHandler struct{}

func (s *WalletHandler) Hello(ctx context.Context, req *proto.HelloRequest2, rsp *proto.HelloResponse2) error {
	log.Print("Received Say.Hello request")
	rsp.Greeting = "Hello darnice mani" + req.Name
	return nil
}

func (s *WalletHandler) CreateWallet(ctx context.Context, req *proto.CreateWalletRequest, rsp *proto.CreateWalletResponse) error {
	log.Print("Received Say.CreateWallet request")
	fmt.Println(req.String())
	var err error
	var addr string
	rsp.Data = new(proto.CreateWalletPos)
	tokenModel := &Tokens{Id: int(req.Tokenid)}
	_, err = tokenModel.GetByid(int(req.Tokenid))

	switch tokenModel.Signature {
	case "eip155", "eth":
		addr, err = Neweth(int(req.Userid), int(req.Tokenid), "123456", tokenModel.Chainid)
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
	signtxstr, err := keystore.Signtx(req.To, mount, req.Gasprice)

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
	rets, err := utils.RpcSendRawTx(TokenModel.Node, req.Signtx)
	if err != nil {
		rsp.Code = "1"
		rsp.Msg = err.Error()
		return nil
	}
	txhash, ok := rets["result"]
	if ok {
		rsp.Code = "0"
		rsp.Msg = "发送成功"
		rsp.Data = new(proto.SendRawTxPos)
		rsp.Data.Result = txhash.(string)
		return nil
	}
	if !ok {
		error := rets["error"].(map[string]string)
		rsp.Code = error["code"]
		rsp.Msg = error["message"]
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
	rsp.Code = "0"
	rsp.Msg = "添加1成功"
	rsp.Data = new(proto.NilPos)
	TibiAddressModel := new(TibiAddress)
	_, err := TibiAddressModel.Save(int(req.Uid), int(req.Tokenid), req.Address, req.Mark)
	if err != nil {
		rsp.Code = "1"
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
	rets, err := TibiAddressModel.List(int(req.Uid), int(req.Tokenid))
	if err != nil {
		return err
	}
	rsp.Data = rets
	fmt.Println(rets)
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
