package controller

import (
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	"net/http"
	proto "digicon/proto/rpc"
	"github.com/gin-gonic/gin"
	"github.com/liudng/godump"
)

type MarketGroup struct{}

func (s *MarketGroup) Router(r *gin.Engine) {
	action := r.Group("/market")
	{
		action.GET("/history/kline", s.HistoryKline)
		action.GET("/symbols", s.Symbols)
		action.GET("/entrust_quenes", s.EntrustQuene)
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


func (s *MarketGroup) Symbols(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type SymbolsParam struct {
		TokenId int32  `form:"token_id" binding:"required"`
	}

	var param SymbolsParam

	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallSymbols(&proto.SymbolsRequest{
		Type:param.TokenId,

	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list",rsp.Data)
}

func (s *MarketGroup) EntrustQuene(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	type EntrustQueneParam struct {
		Symbol string  `form:"symbol" binding:"required"`
	}

	var param EntrustQueneParam

	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallEntrustQuene(&proto.EntrustQueneRequest{
		Symbol:param.Symbol,

	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	godump.Dump(len(rsp.Sell))
	ret.SetDataSection("sell",rsp.Sell)
	ret.SetDataSection("buy",rsp.Buy)
}