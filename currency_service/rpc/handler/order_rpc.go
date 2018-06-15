package handler

import (
	"context"
	. "digicon/currency_service/log"
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"time"
)

type RPCServer struct{}

func (s *RPCServer) AdminCmd(ctx context.Context, req *proto.AdminRequest, rsp *proto.AdminResponse) error {
	Log.Println("Received Say.Hello request ")
	return nil
}

//获取订单列表
func (s *RPCServer) OrdersList(ctx context.Context, req *proto.OrderRequest, rsp *proto.OrdersListResponse) error {
	Log.Println("Received Order request ...")
	result := []model.Order{}
	o := model.Order{}
	rsp.Err = o.List(req.Page, req.PageNum, &result)
	odc := proto.OrdersListResponse_Orders{}
	for i := 0; i < len(result); i++ {
		value := result[i]
		odc.OrderId = value.OrderId
		odc.AdId = value.AdId
		odc.AdType = value.AdType
		odc.Price = value.Price
		odc.Num = value.Num
		odc.TokenId = value.TokenId
		odc.PayId = value.PayId
		odc.States = value.States
		odc.PayStatus = value.PayStatus
		odc.CancelType = value.CancelType
		odc.CreatedTime = value.CreatedTime
		odc.UpdatedTime = value.UpdatedTime
		rsp.Orders = append(rsp.Orders, &odc)
	}
	return nil
}

//产生订单
func (s *RPCServer) AddOrder(ctx context.Context, req *proto.AddOrderRequest, rsp *proto.OrderResponse) error {

	od := new(model.Order)
	od.OrderId = req.OrderId
	od.AdId = req.AdId
	od.AdType = req.AdType
	od.Price = req.Price
	od.Num = req.Num
	od.TokenId = req.TokenId
	od.PayId = req.PayId
	od.SellId = req.SellId
	od.SellName = req.SellName
	od.BuyId = req.BuyId
	od.BuyName = req.BuyName
	od.Fee = req.Fee
	od.States = req.States
	od.PayStatus = req.PayStatus
	od.CancelType = req.CancelType
	od.CreatedTime = time.Now()
	od.UpdatedTime = time.Now()

	rsp.Code = od.Add()
	return nil
}

//更新
func (s *RPCServer) UpdateOrder(ctex context.Context, req *proto.AddOrderRequest, rsp *proto.OrderResponse) error {
	od := new(model.Order)
	od.OrderId = req.OrderId
	od.AdId = req.AdId
	od.AdType = req.AdType
	od.Price = req.Price
	od.Num = req.Num
	od.TokenId = req.TokenId
	od.PayId = req.PayId
	od.SellId = req.SellId
	od.SellName = req.SellName
	od.BuyId = req.BuyId
	od.BuyName = req.BuyName
	od.Fee = req.Fee
	od.States = req.States
	od.PayStatus = req.PayStatus
	od.CancelType = req.CancelType
	od.UpdatedTime = time.Now()

	rsp.Code = od.Update()
	return nil
}
