package handler

import (
	"digicon/currency_service/model"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"

	"golang.org/x/net/context"
)

func (*RPCServer) Alipay(ctx context.Context, req *proto.AlipayRequest, rsp *proto.PaysResponse) error {
	var result int32
	var err error
	p := model.UserCurrencyAlipayPay{}
	result, err = p.SetAlipay(req)
	rsp.Code = result
	rsp.Message = GetErrorMessage(rsp.Code)
	if err != nil {
		return err
	}
	return nil
}

func (*RPCServer) BankPay(ctx context.Context, req *proto.BankPayRequest, rsp *proto.PaysResponse) error {
	var result int32
	var err error
	p := model.UserCurrencyBankPay{}
	result, err = p.SetBankPay(req)
	rsp.Code = result
	rsp.Message = GetErrorMessage(rsp.Code)
	if err != nil {
		return err
	}
	return nil
}

func (*RPCServer) Paypal(ctx context.Context, req *proto.PaypalRequest, rsp *proto.PaysResponse) error {
	var result int32
	var err error
	p := model.UserCurrencyPaypalPay{}
	result, err = p.SetPaypal(req)
	rsp.Code = result
	rsp.Message = GetErrorMessage(rsp.Code)
	if err != nil {
		return err
	}
	return nil
}

func (*RPCServer) WeChatPay(ctx context.Context, req *proto.WeChatPayRequest, rsp *proto.PaysResponse) error {
	var result int32
	var err error
	p := model.UserCurrencyWechatPay{}
	result, err = p.SetWechatPay(req)
	rsp.Code = result
	rsp.Message = GetErrorMessage(rsp.Code)
	if err != nil {
		return err
	}
	return nil
}
