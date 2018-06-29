package controller

import (
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MarketGroup struct{}

func (s *MarketGroup) Router(r *gin.Engine) {
	action := r.Group("/market")
	{
		action.GET("/history/kline", s.HistoryKline)
	}
}

func (s *MarketGroup) HistoryKline(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type KlineParam struct {
		Symbol string `form:"symbol" binding:"required"`
		Period string `form:"period" binding:"required"`
		Size   int32  `form:"size" binding:"required"`
	}

	var param KlineParam

	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallHistoryKline(param.Symbol, param.Period, param.Size)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(ERRCODE_SUCCESS)
	ret.SetDataSection("list", rsp)
}
