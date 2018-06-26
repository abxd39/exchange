package controller

import (
	"digicon/common/convert"
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TokenGroup struct{}

func (s *TokenGroup) Router(r *gin.Engine) {
	action := r.Group("/token")
	{
		action.POST("/entrust_order", s.EntrustOrder)
	}
}

func (s *TokenGroup) EntrustOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type EntrustOrderParam struct {
		Uid          int32  `form:"uid" binding:"required"`
		TokenId      int32  `form:"token_id" binding:"required"`
		TokenTradeId int32  `form:"token_trade_id" binding:"required"`
		Opt          int32  `form:"opt" `
		OnPrice      string `form:"on_price" binding:"required"`
		Type         int32  `form:"type" binding:"required"`
		Num          string `form:"num" binding:"required"`
	}
	var param EntrustOrderParam

	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	o, err := convert.StringToInt64By8Bit(param.OnPrice)
	if err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	n, err := convert.StringToInt64By8Bit(param.Num)
	if err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallEntrustOrder(&proto.EntrustOrderRequest{
		TokenId: param.TokenId,
		Opt:     proto.ENTRUST_OPT(param.Opt),
		OnPrice: o,
		Num:     n,
		Uid:     param.Uid,
		Type:    proto.ENTRUST_TYPE(param.Type),
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}
