package handler

import (
	"context"
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
)

//获取订单列表
func (s *RPCServer) OrdersList(ctx context.Context, req *proto.OrdersListRequest, rsp *proto.OrdersListResponse) error {
	result := []model.Order{}
	o := new(model.Order)
	rsp.Total,rsp.Page,rsp.PageNum,rsp.Err = o.List(req.Page, req.PageNum,req.AdType, req.Status, req.TokenId, req.CreatedTime, &result)

	for i := 0; i < len(result);i++{
		value := result[i]
		od := proto.OrdersListResponse_Orders{}
		od.Id          = value.Id
		od.OrderId     = value.OrderId
		od.AdId        = value.AdId
		od.AdType      = value.AdType
		od.Price       = value.Price
		od.Num         = value.Num
		od.TokenId     = value.TokenId
		od.PayId       = value.PayId
		od.States      = value.States
		od.PayStatus   = value.PayStatus
		od.CancelType  = value.CancelType
		od.CreatedTime = value.CreatedTime
		od.UpdatedTime = value.UpdatedTime
		rsp.Orders = append(rsp.Orders, &od)
	}
	return nil
}

// 取消订单
func (s *RPCServer) CancelOrder( ctx context.Context, req *proto.CancelOrderRequest, rsp *proto.OrderResponse) error {
	code := new(model.Order).Cancel(req.Id, req.CancelType)
	rsp.Code = code
	return nil
}

// 删除订单
func (s *RPCServer) DeleteOrder(ctx context.Context, req *proto.OrderRequest, rsp *proto.OrderResponse) error {
	code := new(model.Order).Delete(req.Id)
	rsp.Code = code
	return nil
}

// 确认放行
func (s *RPCServer) ConfirmOrder(ctx context.Context, req *proto.OrderRequest, rsp *proto.OrderResponse) error{
	code := new(model.Order).Confirm(req.Id)
	rsp.Code = code
	return nil
}





