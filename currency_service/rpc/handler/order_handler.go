package handler

import (
	"context"
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"encoding/json"
	. "digicon/currency_service/log"
	"time"
	"fmt"
	"strconv"
	"bytes"
)


// 产生订单 ID
// 币种 , 年月,  时间秒, 用户id
func createOrderId(userId int32) (orderId string) {
	var buffer bytes.Buffer
	tn := time.Now()
	tnn := tn.UnixNano()
	tnns := strconv.FormatInt(tnn, 10)       // 获取微秒时间
	tnstr  := tn.Format("2006-01-02")
	tnYear := tnstr[:4]
	tnMonth := tnstr[5:7]
	buffer.WriteString(tnYear)
	buffer.WriteString(tnMonth)
	buffer.WriteString(tnns[len(tnns) - 6:])
	buffer.WriteString(strconv.FormatInt(int64(userId), 10))
	orderId = buffer.String()
	return
}




// 获取订单列表
func (s *RPCServer) OrdersList(ctx context.Context, req *proto.OrdersListRequest, rsp *proto.OrdersListResponse) error {
	result := []model.Order{}
	o := new(model.Order)
	rsp.Total,rsp.Page,rsp.PageNum,rsp.Err = o.List(req.Page, req.PageNum,req.AdType, req.States, req.TokenId, req.CreatedTime, &result)

	orders , err := json.Marshal(result)
	if err != nil {
		Log.Errorln(err.Error())
		rsp.Orders = "[]"
		rsp.Message = err.Error()
		return err
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

// 添加订单
func (s *RPCServer) AddOrder(ctx context.Context, req *proto.AddOrderRequest, rsp *proto.OrderResponse) error {
	od := new(model.Order)
	if err := json.Unmarshal([]byte(req.Order), &od); err != nil {
		Log.Println(err.Error())
		fmt.Println(err.Error())
	}

	od.OrderId = createOrderId(req.Uid)
	od.States = 1
	od.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
	od.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")

	id, code := od.Add()
	rsp.Code = code
	rsp.Data = strconv.FormatUint(id,10)
	return nil
}






