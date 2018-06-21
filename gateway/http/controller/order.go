package controller

import (
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderRequest struct {
	Id uint64 `form:"id" json:"id"  binding:"required"` //order 表Id
}

type CancelOrderRequest struct {
	//Id         uint64  `form:"id" json:"id"  binding:"required"`          //order 表Id
	OrderRequest
	CancelType uint32 `form:"id" json:"cancel_type" binding:"required"` //取消类型: 1卖方 2 买方
}

type OneOrder struct {
	AdId     uint64  `form:"ad_id"   json:"ad_id"   binding:"required"`       // 广告id
	AdType   uint32  `form:"ad_type" json:"ad_type" binding:"required"`       // 广告类型：1出售 2购买
	Price    float64 `form:"price"   json:"price"   binding:"required"`       // 单价
	Num      float64 `form:"num"        json:"num"     binding:"required"`    // 交易数量
	TokenId  uint64  `form:"token_id"   json:"token_id"   binding:"required"` // 货币类型
	PayId    uint64  `form:"pay_id"     json:"pay_id"     binding:"required"` // 支付类型
	SellId   uint64  `form:"sell_id"    json:"sell_id"    binding:"required"` // 卖家id
	SellName string  `form:"sell_name"  json:"sell_name"  binding:"required"` // 卖家昵称
	BuyId    uint64  `form:"buy_id"     json:"buy_id"     binding:"required"` // 买家id
	BuyName  string  `form:"buy_name"   json:"buy_name"   binding:"required"` // 买家昵称
}

type AddOrder struct {
	Uid int32 `form:"uid"   json:"uid"  binding:"required"` // 用户 id
	OneOrder
}

type Order struct {
	OrderRequest
	OneOrder
	Fee         float64 `form:"fee"          json:"fee"  `        // 手续费用
	OrderId     string  `form:"order_id"     json:"order_id" `    // 订单id
	States      uint32  `form:"states"       json:"states"`       // 订单状态 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消
	PayStatus   uint32  `form:"pay_status"   json:"pay_status"`   // 支付状态 1待支付 2待放行(已支付) 3确认支付(已完成)
	CancelType  uint32  `form:"cancel_type"  json:"cancel_type"`  // 取消类型 1卖方 2 买方
	CreatedTime string  `form:"created_time" json:"created_time"` //
	UpdatedTime string  `form:"updated_time" json:"updated_time"`
}

func (this *CurrencyGroup) OrdersList(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	type OrderListParam struct {
		Page      int32   `form:"page"       json:"page"`
		PageNum   int32   `form:"page_num"   json:"page_num"`
		TokenId   float64 `form:"token_id"   json:"token_id"`
		AdType    uint32  `form:"ad_type"    json:"ad_type"`
		States    uint32  `form:"states"     json:"states"`
		StartTime string  `form:"start_time" json:"start_time"`
		EndTime   string  `form:"end_time"   json:"end_time"`
		Id        uint64  `form:"id"           json:"id"`
	}
	var param OrderListParam
	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	var tmpStates uint32
	if param.States != 0 {
		tmpStates = param.States
	} else {
		tmpStates = 100
	}

	rsp, err := rpc.InnerService.CurrencyService.CallOrdersList(&proto.OrdersListRequest{
		Page:      param.Page,
		PageNum:   param.PageNum,
		TokenId:   param.TokenId,
		AdType:    param.AdType,
		States:    tmpStates,
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		Id:        param.Id,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	var orders []Order
	if err = json.Unmarshal([]byte(rsp.Orders), &orders); err != nil {
		Log.Errorln(err.Error())
	}

	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list", orders)
	ret.SetDataSection("total", rsp.Total)
	ret.SetDataSection("page", rsp.Page)
	ret.SetDataSection("page_num", rsp.PageNum)
}

// 取消订单
func (this CurrencyGroup) CancelOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var param CancelOrderRequest
	err := c.ShouldBind(&param)
	if err != nil {
		Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallCancelOrder(&proto.CancelOrderRequest{
		Id:         param.Id,
		CancelType: param.CancelType,
	})
	ret.SetErrCode(rsp.Code, rsp.Message)
}

// 删除订单
func (this CurrencyGroup) DeleteOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var param OrderRequest
	err := c.ShouldBind(&param)
	if err != nil {
		Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallDeleteOrder(&proto.OrderRequest{
		Id: param.Id,
	})
	ret.SetErrCode(rsp.Code, rsp.Message)
}

// 待放行
func (this CurrencyGroup) ReadyOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var param OrderRequest
	err := c.ShouldBind(&param)
	if err != nil {
		Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallReadyOrder(&proto.OrderRequest{
		Id: param.Id,
	})
	ret.SetErrCode(rsp.Code, rsp.Message)
}

// 确认放行
func (this CurrencyGroup) ConfirmOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var param OrderRequest
	err := c.ShouldBind(&param)
	if err != nil {
		Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallConfirmOrder(&proto.OrderRequest{
		Id: param.Id,
	})
	ret.SetErrCode(rsp.Code, rsp.Message)
}

// 添加订单
func (this CurrencyGroup) AddOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var param AddOrder
	err := c.ShouldBind(&param)
	if err != nil {
		Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	orderStr, _ := json.Marshal(param)
	//fmt.Println("params uid:", param.Uid)
	rsp, err := rpc.InnerService.CurrencyService.CallAddOrder(&proto.AddOrderRequest{
		Order: string(orderStr),
		Uid:   param.Uid,
	})
	ret.SetErrCode(rsp.Code, rsp.Message)
	ret.SetDataSection("id", rsp.Data)

}
