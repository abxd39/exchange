package client

import (
	"context"
	. "digicon/gateway/log"
	proto "digicon/proto/rpc"
)

func (s *CurrencyRPCCli) CallOrdersList(req *proto.OrdersListRequest) (rsp *proto.OrdersListResponse, err error) {
	rsp, err = s.conn.OrdersList(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *CurrencyRPCCli) CallDeleteOrder(req *proto.OrderRequest) (rsp *proto.OrderResponse, err error) {
	rsp, err = s.conn.DeleteOrder(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *CurrencyRPCCli) CallCancelOrder(req *proto.CancelOrderRequest) (rsp *proto.OrderResponse, err error) {
	rsp, err = s.conn.CancelOrder(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *CurrencyRPCCli) CallConfirmOrder(req *proto.OrderRequest) (rsp *proto.OrderResponse, err error) {

	rsp, err = s.conn.ConfirmOrder(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *CurrencyRPCCli) CallReadyOrder(req *proto.OrderRequest) (rsp *proto.OrderResponse, err error) {
	rsp, err = s.conn.ReadyOrder(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *CurrencyRPCCli) CallAddOrder(req *proto.AddOrderRequest) (rsp *proto.OrderResponse, err error) {
	rsp, err = s.conn.AddOrder(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}
