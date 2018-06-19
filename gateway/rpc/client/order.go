package client

import (
	"context"
	. "digicon/gateway/log"
	proto "digicon/proto/rpc"
	"fmt"
)

func (s *CurrencyRPCCli) CallOrdersList(req *proto.OrdersListRequest) (rsp *proto.OrdersListResponse, err error) {
	rsp, err = s.conn.OrdersList(context.TODO(), req)
	fmt.Println(len(rsp.Orders))
	fmt.Println(rsp.Orders)
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

<<<<<<< HEAD
func (s *CurrencyRPCCli) CallCancelOrder(req *proto.OrderRequest) (rsp *proto.OrderResponse, err error) {
=======
func (s *CurrencyRPCCli) CallCancelOrder(req *proto.CancelOrderRequest)(rsp *proto.OrderResponse, err error){
>>>>>>> efbae774ef18301514fd2b43f7d2dbf165bb5624
	rsp, err = s.conn.CancelOrder(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}
<<<<<<< HEAD
=======

func (s *CurrencyRPCCli) CallConfirmOrder(req *proto.OrderRequest)(rsp *proto.OrderResponse, err error) {
	rsp, err = s.conn.ConfirmOrder(context.TODO(), req)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}








>>>>>>> efbae774ef18301514fd2b43f7d2dbf165bb5624
