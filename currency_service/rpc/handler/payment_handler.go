package handler

import (
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"fmt"
	log "github.com/sirupsen/logrus"

	"digicon/proto/common"
	"github.com/gin-gonic/gin/json"
	"golang.org/x/net/context"
)



/////////////   ali pay  //////////////////

func (*RPCServer) Alipay(ctx context.Context, req *proto.AlipayRequest, rsp *proto.PaysResponse) (err error) {
	p := model.UserCurrencyAlipayPay{
		Uid:         req.Uid,
		Name:        req.Name,
		Alipay:      req.Alipay,
		ReceiptCode: req.ReceiptCode,
	}
	rsp.Code, err = p.SetAlipay(req)
	rsp.Data = p.ReceiptCode
	fmt.Println(rsp.Code)
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func (*RPCServer) GetAliPay(ctx context.Context, req *proto.PayRequest, rsp *proto.PaysResponse) (err error) {
	amd := new(model.UserCurrencyAlipayPay)
	err = amd.GetByUid(req.Uid)
	data, err := json.Marshal(amd)
	if err != nil {
		log.Errorln(err.Error())
		rsp.Code = errdefine.ERRCODE_UNKNOWN
	} else {
		rsp.Code = errdefine.ERRCODE_SUCCESS
		rsp.Data = string(data)
	}
	return
}

func (*RPCServer) UpdateAliPay(ctx context.Context, req *proto.AlipayRequest, rsp *proto.PaysResponse) (err error) {
	p := model.UserCurrencyAlipayPay{
		Uid:         req.Uid,
		Name:        req.Name,
		Alipay:      req.Alipay,
		ReceiptCode: req.ReceiptCode,
	}
	rsp.Code, err = p.SetAlipay(req)
	rsp.Data = p.ReceiptCode
	fmt.Println(rsp.Code)
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

////////////////  bank  pay //////////////////

func (*RPCServer) BankPay(ctx context.Context, req *proto.BankPayRequest, rsp *proto.PaysResponse) (err error) {
	p := model.UserCurrencyBankPay{
		BankName: req.BankName,
		CardNum:  req.CardNum,
		BankInfo: req.BankInfo,
		Name:     req.Name,
	}
	rsp.Code, err = p.SetBankPay(req)
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func (*RPCServer) GetBankPay(ctx context.Context, req *proto.PayRequest, rsp *proto.PaysResponse) (err error) {
	p := new(model.UserCurrencyBankPay)
	err = p.GetByUid(req.Uid)

	data, err := json.Marshal(p)
	if err != nil {
		rsp.Code = errdefine.ERRCODE_UNKNOWN
	} else {
		rsp.Code = errdefine.ERRCODE_SUCCESS
		rsp.Data = string(data)
	}
	return nil
}

func (*RPCServer) UpdateBankPay(ctx context.Context, req *proto.BankPayRequest, rsp *proto.PaysResponse) (err error) {
	p := new(model.UserCurrencyBankPay)
	p.BankName = req.BankName
	p.CardNum = req.CardNum
	p.BankInfo = req.BankInfo
	p.Name = req.Name
	rsp.Code, err = p.SetBankPay(req)
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

////////////////////  paypal /////////////////////////
func (*RPCServer) Paypal(ctx context.Context, req *proto.PaypalRequest, rsp *proto.PaysResponse) (err error) {
	p := model.UserCurrencyPaypalPay{}
	p.Paypal = req.Paypal
	rsp.Code, err = p.SetPaypal(req)
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func (*RPCServer) GetPaypal(ctx context.Context, req *proto.PayRequest, rsp *proto.PaysResponse) (err error) {
	paypal := new(model.UserCurrencyPaypalPay)
	err = paypal.GetByUid(req.Uid)
	data, err := json.Marshal(paypal)
	if err != nil {
		rsp.Code = errdefine.ERRCODE_UNKNOWN
	} else {
		rsp.Code = errdefine.ERRCODE_SUCCESS
		rsp.Data = string(data)
	}
	return nil
}

func (*RPCServer) UpdatePaypal(ctx context.Context, req *proto.PaypalRequest, rsp *proto.PaysResponse) (err error) {
	p := model.UserCurrencyPaypalPay{}
	p.Paypal = req.Paypal
	rsp.Code, err = p.SetPaypal(req)
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

///////////////  wechat pay ////////////////

func (*RPCServer) WeChatPay(ctx context.Context, req *proto.WeChatPayRequest, rsp *proto.PaysResponse) (err error) {
	p := model.UserCurrencyWechatPay{
		Name:        req.Name,
		Wechat:      req.Wechat,
		ReceiptCode: req.ReceiptCode,
	}
	fmt.Println(req)
	rsp.Code, err = p.SetWechatPay(req)
	if err != nil {
		fmt.Println(err.Error())
		log.Errorf(err.Error())
	}

	return nil
}

func (*RPCServer) GetWeChatPay(ctx context.Context, req *proto.PayRequest, rsp *proto.PaysResponse) (err error) {
	wcp := new(model.UserCurrencyWechatPay)
	err = wcp.GetByUid(req.Uid)
	data, err := json.Marshal(wcp)
	if err != nil {
		rsp.Code = errdefine.ERRCODE_UNKNOWN
	} else {
		rsp.Code = errdefine.ERRCODE_SUCCESS
		rsp.Data = string(data)
	}

	return nil
}

func (*RPCServer) UpdateWeChatPay(ctx context.Context, req *proto.WeChatPayRequest, rsp *proto.PaysResponse) (err error) {
	p := model.UserCurrencyWechatPay{
		Name:        req.Name,
		Wechat:      req.Wechat,
		ReceiptCode: req.ReceiptCode,
	}
	rsp.Code, err = p.SetWechatPay(req)
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}



///////////////////////////////////////////////////////////


func (*RPCServer) GetPaySet (ctx context.Context, req *proto.PayRequest, rsp *proto.PaysResponse) ( err error){
	var bankpaySet uint32
	var alipaySet  uint32
	var wechatpaySet uint32
	var paypalSet    uint32

	bankpay := new(model.UserCurrencyBankPay)
	err = bankpay.GetByUid(req.Uid)

	if bankpay.CardNum != ""{
		bankpaySet  = 1
	}
	alipay := new(model.UserCurrencyAlipayPay)
	alipay.GetByUid(req.Uid)
	if alipay.Alipay != ""{
		alipaySet = 1
	}

	wechatpay := new(model.UserCurrencyWechatPay)
	wechatpay.GetByUid(req.Uid)
	if wechatpay.Wechat != ""{
		wechatpaySet = 1
	}

	paypal := new(model.UserCurrencyPaypalPay)
	paypal.GetByUid(req.Uid)

	if paypal.Paypal != ""{
		paypalSet = 1
	}

	type PaySet struct {
		BankPay     uint32  `form:"bank_pay"    json:"bank_pay"`
		AliPay      uint32  `form:"ali_pay"     json:"ali_pay"`
		WeChatPay   uint32  `form:"wechat_pay"  json:"wechat_pay"`
		PaypalPay   uint32  `form:"paypal_pay"  json:"paypal_pay"`
	}

	var payset PaySet
	payset.PaypalPay = paypalSet
	payset.WeChatPay = wechatpaySet
	payset.AliPay = alipaySet
	payset.BankPay = bankpaySet

	data, err := json.Marshal(payset)
	if err != nil {
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return
	}
	rsp.Code = errdefine.ERRCODE_SUCCESS
	rsp.Data = string(data)
	return
}