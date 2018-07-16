package client

import (
	"context"
	proto "digicon/proto/rpc"
)

///////////   bank pay ////////////

func (s *CurrencyRPCCli) CallBankPay(req *proto.BankPayRequest) (*proto.PaysResponse, error) {
	return s.conn.BankPay(context.TODO(), req)
}

func (s *CurrencyRPCCli) CallGetBankPay(req *proto.PayRequest) (*proto.PaysResponse, error) {
	return s.conn.GetBankPay(context.TODO(), req)
}

func (s *CurrencyRPCCli) CallUpdateBankPay(req *proto.BankPayRequest) (*proto.PaysResponse, error) {
	return s.conn.UpdateBankPay(context.TODO(), req)
}

///////////  ali  pay ////////////

func (s *CurrencyRPCCli) CallAliPay(req *proto.AlipayRequest) (*proto.PaysResponse, error) {
	return s.conn.Alipay(context.TODO(), req)
}

func (s *CurrencyRPCCli) CallGetAliPay(req *proto.PayRequest) (*proto.PaysResponse, error) {
	return s.conn.GetAliPay(context.TODO(), req)
}

func (s *CurrencyRPCCli) CallUpdateAliPay(req *proto.AlipayRequest) (*proto.PaysResponse, error) {
	return s.conn.UpdateAliPay(context.TODO(), req)
}

///////////   wechat pay ////////////

func (s *CurrencyRPCCli) CallWeChatPay(req *proto.WeChatPayRequest) (*proto.PaysResponse, error) {
	return s.conn.WeChatPay(context.TODO(), req)
}

func (s *CurrencyRPCCli) CallGetWeChatPay(req *proto.PayRequest) (*proto.PaysResponse, error) {
	return s.conn.GetWeChatPay(context.TODO(), req)
}

func (s *CurrencyRPCCli) CallUpdateWeChatPay(req *proto.WeChatPayRequest) (*proto.PaysResponse, error) {
	return s.conn.UpdateWeChatPay(context.TODO(), req)
}

///////////  paypal pay ////////////

func (s *CurrencyRPCCli) CallPaypal(req *proto.PaypalRequest) (*proto.PaysResponse, error) {
	return s.conn.Paypal(context.TODO(), req)
}

func (s *CurrencyRPCCli) CallGetPaypal(req *proto.PayRequest) (*proto.PaysResponse, error) {
	return s.conn.GetPaypal(context.TODO(), req)
}
func (s *CurrencyRPCCli) CallUpdatePaypal(req *proto.PaypalRequest) (*proto.PaysResponse, error) {
	return s.conn.UpdatePaypal(context.TODO(), req)
}
