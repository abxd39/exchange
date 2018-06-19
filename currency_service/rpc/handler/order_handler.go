package handler

import (
	"context"
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"encoding/json"
	. "digicon/currency_service/log"
)

//获取订单列表
func (s *RPCServer) OrdersList(ctx context.Context, req *proto.OrdersListRequest, rsp *proto.OrdersListResponse) error {
	result := []model.Order{}
	o := new(model.Order)
	rsp.Total,rsp.Page,rsp.PageNum,rsp.Err = o.List(req.Page, req.PageNum,req.AdType, req.Status, req.TokenId, req.CreatedTime, &result)

	orders , err := json.Marshal(result)
	if err != nil {
		Log.Errorln(err.Error())
	}
	rsp.Orders = string(orders)
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





