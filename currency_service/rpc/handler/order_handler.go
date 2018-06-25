package handler

import (
	"bytes"
	"context"
	. "digicon/currency_service/log"
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// 获取订单列表
func (s *RPCServer) OrdersList(ctx context.Context, req *proto.OrdersListRequest, rsp *proto.OrdersListResponse) error {
	result := []model.Order{}
	o := new(model.Order)

	rsp.Total, rsp.Page, rsp.PageNum, rsp.Err = o.List(req.Page, req.PageNum, req.AdType, req.States, req.Id, req.TokenId, req.StartTime, req.EndTime, &result)

	orders, err := json.Marshal(result)
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
func (s *RPCServer) CancelOrder(ctx context.Context, req *proto.CancelOrderRequest, rsp *proto.OrderResponse) error {
	updateTimeStr := time.Now().Format("2006-01-02 15:04:05")
	code, msg := new(model.Order).Cancel(req.Id, req.CancelType, updateTimeStr)
	rsp.Code = code
	rsp.Message = msg
	return nil
}

// 删除订单
func (s *RPCServer) DeleteOrder(ctx context.Context, req *proto.OrderRequest, rsp *proto.OrderResponse) error {
	fmt.Println(req.Id)
	updateTimeStr := time.Now().Format("2006-01-02 15:04:05")
	code, msg := new(model.Order).Delete(req.Id, updateTimeStr)
	rsp.Code = code
	rsp.Message = msg
	return nil
}

// 确认放行
func (s *RPCServer) ConfirmOrder(ctx context.Context, req *proto.OrderRequest, rsp *proto.OrderResponse) error {
	updateTimeStr := time.Now().Format("2006-01-02 15:04:05")
	code, msg := new(model.Order).Confirm(req.Id, updateTimeStr)
	rsp.Code = code
	rsp.Message = msg
	return nil
}

// 待放行
func (s *RPCServer) ReadyOrder(ctx context.Context, req *proto.OrderRequest, rsp *proto.OrderResponse) error {
	updateTimeStr := time.Now().Format("2006-01-02 15:04:05")
	code, msg := new(model.Order).Ready(req.Id, updateTimeStr)
	rsp.Code = code
	rsp.Message = msg
	return nil
}

// 添加订单
func (s *RPCServer) AddOrder(ctx context.Context, req *proto.AddOrderRequest, rsp *proto.OrderResponse) error {
	od := new(model.Order)
	if err := json.Unmarshal([]byte(req.Order), &od); err != nil {
		Log.Println(err.Error())
		fmt.Println(err.Error())
	}

	od.OrderId = createOrderId(req.Uid, od.TokenId)
	od.States = 1
	od.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
	od.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")

	id, code := od.Add()
	rsp.Code = code
	rsp.Data = strconv.FormatUint(id, 10)
	return nil
}

// 产生订单 ID
//  uid, 时间秒,
func createOrderId(userId int32, tokenId uint64) (orderId string) {
	tn := time.Now()
	tnn := tn.UnixNano()
	tnnStr := strconv.FormatInt(tnn, 10) // 获取微秒时间
	tnnStrLen := len(tnnStr)
	tStr := tn.Format("2006-01-02 15:04:05.34234")

	var buffer bytes.Buffer
	buffer.WriteString(strconv.FormatInt(int64(userId), 10))
	buffer.WriteString(tStr[2:4])    // 年
	buffer.WriteString(tStr[5:7])    // 月
	buffer.WriteString(tStr[8:10])   // 日
	buffer.WriteString(tStr[17:19])  // 秒
	//buffer.WriteString(tStr[22:24])
	buffer.WriteString(tnnStr[tnnStrLen-5 : tnnStrLen-2]) // 微秒
	orderId = buffer.String()
	return
}




