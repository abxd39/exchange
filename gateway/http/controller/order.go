package controller

import (
	"github.com/gin-gonic/gin"
	. "digicon/proto/common"
	"digicon/gateway/rpc"
	. "digicon/gateway/log"
	"net/http"
	proto "digicon/proto/rpc"
	"encoding/json"
)


type OrderRequest  struct {
	Id    uint64  `form:"id" json:"id"  binding:"required"`     //order 表Id
}


type CancelOrderRequest  struct {
	Id         uint64  `form:"id" json:"id"  binding:"required"`          //order 表Id
	CancelType uint32  `form:"id" json:"cancel_type" binding:"required"`  //取消类型: 1卖方 2 买方
}


type Order struct {
	Id          uint64       `json:"id"`                      // id
	OrderId     uint64       `json:"order_id" `               // 订单id
	AdId        uint64       `json:"ad_id"  `                 // 广告id
	AdType      uint32       `json:"ad_type"`                 // 广告类型
	Price       float64      `json:"price"  `                 // 单价
	Num         float64      `json:"num"  `                   // 交易数量
	TokenId     uint64       `json:"token_id"  `              // 类型：1出售 2购买
	PayId       uint64       `json:"pay_id" `                 // 支付类型
	SellId      uint64       `json:"sell_id"  `               // 卖家id
	SellName    string       `json:"sell_name"  `             // 卖家昵称
	BuyId       uint64       `json:"buy_id" `                 // 买家id
	BuyName     string       `json:"buy_name" `               // 买家昵称
	Fee         float64      `json:"fee"  `                   // 手续费用
	States      uint32       `json:"states"`                  // 订单状态 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消
	PayStatus   uint32       `json:"pay_status"`			  // 支付状态 1待支付 2待放行(已支付) 3确认支付(已完成)
	CancelType  uint32       `json:"cancel_type"`             // 取消类型 1卖方 2 买方
	CreatedTime string       `json:"created_time"`            //
	UpdatedTime string       `json:"updated_time"`
}



func (this *CurrencyGroup) OrdersList(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	type OrderListParam  struct {
		Page        int32       `form:"page" `
		PageNum     int32       `form:"page_num" `
		TokenId     float64   	`form:"token_id"`
		AdType      uint32      `form:"ad_type"`
		States      uint32      `form:"states"`
		CreatedTime string      `form:"created_time"`
	}
	var param OrderListParam
	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	var tmpStates uint32
	if param.States != 0{
		tmpStates = param.States
	}else{
		tmpStates = 100
	}

	rsp, err := rpc.InnerService.CurrencyService.CallOrdersList(&proto.OrdersListRequest{
		Page :       param.Page,
		PageNum:     param.PageNum,
		TokenId:     param.TokenId,
		AdType:      param.AdType,
		Status:      tmpStates,
		CreatedTime: param.CreatedTime,
	})
	if err != nil{
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
func (this CurrencyGroup) CancelOrder(c *gin.Context){
	ret := NewPublciError()
	defer func(){
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
		Id:param.Id,
		CancelType:param.CancelType,
	})
	ret.SetErrCode(rsp.Code, rsp.Message)
}


// 删除订单
func (this CurrencyGroup) DeleteOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func(){
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
		Id:param.Id,
	})
	ret.SetErrCode(rsp.Code, rsp.Message)
}



// 确认放行
func (this CurrencyGroup) ConfirmOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func(){
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
		Id:param.Id,
	})
	ret.SetErrCode(rsp.Code, rsp.Message)
}

