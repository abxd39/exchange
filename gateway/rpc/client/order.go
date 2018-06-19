package client

import (
	proto "digicon/proto/rpc"
	"context"
	. "digicon/gateway/log"
)

func (s *CurrencyRPCCli) CallOrdersList(req *proto.OrdersListRequest)(rsp *proto.OrdersListResponse, err error) {
	rsp, err = s.conn.OrdersList(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}




func (s *CurrencyRPCCli) CallDeleteOrder(req *proto.OrderRequest) (rsp *proto.OrderResponse, err error){
	rsp, err = s.conn.DeleteOrder(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *CurrencyRPCCli) CallCancelOrder(req *proto.CancelOrderRequest)(rsp *proto.OrderResponse, err error){
	rsp, err = s.conn.CancelOrder(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *CurrencyRPCCli) CallConfirmOrder(req *proto.OrderRequest)(rsp *proto.OrderResponse, err error) {
	rsp, err = s.conn.ConfirmOrder(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}








