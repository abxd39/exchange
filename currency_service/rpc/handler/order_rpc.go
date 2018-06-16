package handler

import (
	"context"
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"fmt"
)

//获取订单列表
func (s *RPCServer) OrdersList(ctx context.Context, req *proto.OrdersListRequest, rsp *proto.OrdersListResponse) error {
	result := []model.Order{}

<<<<<<< HEAD
	Log.Println("Received Order request ...")
	result := []model.Order{}
	o := model.Order{}
	rsp.Err = o.List(req.Page, req.PageNum, &result)
	odc := proto.OrdersListResponse_Orders{}
	for i := 0; i < len(result); i++ {
=======
	o := new(model.Order)
	rsp.Total,rsp.Page,rsp.PageNum,rsp.Err = o.List(req.Page, req.PageNum,req.AdType, req.Status, req.TokenId, req.CreatedTime, &result)

	for i := 0; i< len(result);i++{
>>>>>>> 5251c8c21a1e93a567e3eea899dcc3a9d453e922
		value := result[i]
		odc := proto.OrdersListResponse_Orders{}
		odc.Id = value.Id
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
		//fmt.Println(value)
		rsp.Orders = append(rsp.Orders, &odc)
		fmt.Println(rsp.Orders)
	}
	return nil
}

// 取消订单
func (s *RPCServer) CancelOrder( ctx context.Context, req *proto.OrderRequest, rsp *proto.OrderResponse) error {
	code := new(model.Order).Cancel(req.Id)
	rsp.Code = code
	return nil
}

// 删除订单
func (s *RPCServer) DeleteOrder(ctx context.Context, req *proto.OrderRequest, rsp *proto.OrderResponse) error {
	code := new(model.Order).Delete(req.Id)
	rsp.Code = code
	return nil
}
<<<<<<< HEAD
=======







>>>>>>> 5251c8c21a1e93a567e3eea899dcc3a9d453e922
