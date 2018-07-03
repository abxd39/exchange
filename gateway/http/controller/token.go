package controller

import (
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	"github.com/gin-gonic/gin"
	"net/http"
	proto "digicon/proto/rpc"
	"digicon/common/convert"
	. "digicon/gateway/log"
)
type TokenGroup struct{}

func (s *TokenGroup) Router(r *gin.Engine) {
	action := r.Group("/token")
	{
		action.POST("/entrust_order", s.EntrustOrder)
		//action.GET("/market/history/kline", s.HistoryKline)
	}
}

func (s *TokenGroup) EntrustOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type EntrustOrderParam struct {
		Uid uint64 `form:"uid" binding:"required"`
		Symbol string `form:"symbol" binding:"required"`
		Opt int32 `form:"opt" `
		OnPrice string `form:"on_price" `
		Type int32 `form:"type" binding:"required"`
		Num string`form:"num" binding:"required"`
	}
	var param EntrustOrderParam

	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	o,err:=convert.StringToInt64By8Bit(param.OnPrice)
	if err!=nil {
		ret.SetErrCode(ERRCODE_PARAM,err.Error())
		return
	}

	n,err:=convert.StringToInt64By8Bit(param.Num)
	if err!=nil {
		ret.SetErrCode(ERRCODE_PARAM,err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallEntrustOrder(&proto.EntrustOrderRequest{
		Symbol:param.Symbol,
		Opt:proto.ENTRUST_OPT(param.Opt),
		OnPrice:o,
		Num:n,
		Uid:param.Uid,
		Type: proto.ENTRUST_TYPE(param.Type),
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}
