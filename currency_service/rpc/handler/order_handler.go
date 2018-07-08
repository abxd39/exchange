package handler

import (
	"context"
	"digicon/common/encryption"
	. "digicon/currency_service/log"
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"encoding/json"
	"fmt"
	"time"
	"digicon/proto/common"

	"digicon/currency_service/rpc/client"
	"strconv"
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
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return nil
	}
	fmt.Println("req order:", req.Order)
	fmt.Println("od num: ", od.Num)

	ads := new(model.Ads)
	var nowAds *model.Ads
	nowAds = ads.Get(od.AdId)

	od.AdType = nowAds.TypeId
	od.Price = int64(nowAds.Price)
	od.TokenId = uint64(nowAds.TokenId)
	od.SellId = nowAds.Uid
	od.BuyId = uint64(nowAds.Uid)
	od.PayId = nowAds.Pays

	//fmt.Println(od.SellId, od.BuyId)

	var uids []uint64
	uids = append(uids, od.SellId)
	uids = append(uids, od.BuyId)

	nickNames, err := client.InnerService.UserSevice.CallGetNickName(uids)    // rpc 获取用户信息

	if err != nil {
		fmt.Println(err)
		Log.Errorln(err.Error())
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return nil
	}else {
		nickUsers := nickNames.User
		for i := 0; i < len(nickUsers); i ++ {
			if nickUsers[i].Uid == od.SellId {
				od.SellName = nickUsers[i].NickName
			}
			if nickUsers[i].Uid == od.BuyId {
				od.BuyName = nickUsers[i].NickName
			}
		}
	}

	od.OrderId = encryption.CreateOrderId(uint64(req.Uid), int32(od.TokenId))
	od.States = 1
	od.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
	od.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")

	//fmt.Println("od:", od)

	id, code := od.Add()
	rsp.Code = code
	rsp.Data = strconv.FormatUint(id, 10)
	return nil
}
