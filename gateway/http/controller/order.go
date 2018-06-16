package controller

import (
	"github.com/gin-gonic/gin"
	. "digicon/proto/common"
	"digicon/gateway/rpc"
	. "digicon/gateway/log"
	"net/http"
	proto "digicon/proto/rpc"
)


type OrderRequest  struct {
	Id    uint64  `form:"id" json:"id"  binding:"required"`     //order 表Id
}


type CancelOrderRequest  struct {
	Id         uint64  `form:"id" json:"id"  binding:"required"`          //order 表Id
	CancelType uint32  `form:"id" json:"cancel_type" binding:"required"`  //取消类型: 1卖方 2 买方
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
		Status      uint32      `form:"status"`
		CreatedTime string      `form:"created_time"`
	}
	var param OrderListParam
	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallOrdersList(&proto.OrdersListRequest{
		Page :       param.Page,
		PageNum:     param.PageNum,
		TokenId:     param.TokenId,
		AdType:      param.AdType,
		Status:      param.Status,
		CreatedTime: param.CreatedTime,
	})
	if err != nil{
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list", rsp.Orders)
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

