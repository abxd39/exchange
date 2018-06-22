package handler

import (
	"digicon/currency_service/log"
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"fmt"

	"golang.org/x/net/context"
)

func (*RPCServer) Alipay(ctx context.Context, req *proto.AlipayRequest, rsp *proto.PaysResponse) (err error) {
	p := model.UserCurrencyAlipayPay{}
	rsp.Code, err = p.SetAlipay(req)
	fmt.Println(rsp.Code)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}

func (*RPCServer) BankPay(ctx context.Context, req *proto.BankPayRequest, rsp *proto.PaysResponse) (err error) {
	p := model.UserCurrencyBankPay{}
	rsp.Code, err = p.SetBankPay(req)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}

func (*RPCServer) Paypal(ctx context.Context, req *proto.PaypalRequest, rsp *proto.PaysResponse) (err error) {
	p := model.UserCurrencyPaypalPay{}
	rsp.Code, err = p.SetPaypal(req)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}

func (*RPCServer) WeChatPay(ctx context.Context, req *proto.WeChatPayRequest, rsp *proto.PaysResponse) (err error) {
	p := model.UserCurrencyWechatPay{}
	rsp.Code, err = p.SetWechatPay(req)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}
