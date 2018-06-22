package client

import (
	"context"
	proto "digicon/proto/rpc"
)

func (s *CurrencyRPCCli) CallBankPay(req *proto.BankPayRequest) (*proto.PaysResponse, error) {
	return s.conn.BankPay(context.TODO(), req)
}

func (s *CurrencyRPCCli) CallAliPay(req *proto.AlipayRequest) (*proto.PaysResponse, error) {
	return s.conn.Alipay(context.TODO(), req)
}

func (s *CurrencyRPCCli) CallWeChatPay(req *proto.WeChatPayRequest) (*proto.PaysResponse, error) {
	return s.conn.WeChatPay(context.TODO(), req)
}

func (s *CurrencyRPCCli) CallPaypal(req *proto.PaypalRequest) (*proto.PaysResponse, error) {
	return s.conn.Paypal(context.TODO(), req)
}
