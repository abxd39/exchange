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
		//action.GET("/market/history/kline", s.HistoryKline)

		action.GET("/self_symbols", s.SelfSymbols)

		action.GET("/entrust_list", s.EntrustList)

		action.GET("/entrust_history", s.EntrustHistory)

		action.GET("/balance", s.TokenBalance)

		action.GET("/balance_list", s.TokenBalanceList)
	}
}

func (s *TokenGroup) EntrustOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type EntrustOrderParam struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		Symbol  string `form:"symbol" binding:"required"`
		Opt     int32  `form:"opt"  binding:"required"`
		OnPrice string `form:"on_price"  `
		Type    int32  `form:"type" binding:"required"`
		Num     string `form:"num" binding:"required"`
	}
	var param EntrustOrderParam

	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	var o int64
	var err error
	if param.Type == int32(proto.ENTRUST_TYPE_LIMIT_PRICE) {
		o, err = convert.StringToInt64By8Bit(param.OnPrice)
		if err != nil {
			ret.SetErrCode(ERRCODE_PARAM, err.Error())
			return
		}
		if o == 0 {
			ret.SetErrCode(ERRCODE_PARAM)
			return
		}
	} else {
		if param.OnPrice == "" {
			o = 0
		}
	}

	n, err := convert.StringToInt64By8Bit(param.Num)
	if err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallEntrustOrder(&proto.EntrustOrderRequest{
		Symbol:  param.Symbol,
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

func (s *TokenGroup) SelfSymbols(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type SelfSymbolsParam struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		TokenId int32  `form:"token_id" binding:"required"`
	}

	var param SelfSymbolsParam

	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallSelfSymbols(&proto.SelfSymbolsRequest{
		Uid: param.Uid,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection(RET_DATA, rsp.Data)
}

func (s *TokenGroup) EntrustList(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type EntrustListParam struct {
		Uid   uint64 `form:"uid" binding:"required"`
		Token string `form:"token" binding:"required"`
		Limit int32  `form:"limit" `
		Page  int32  `form:"page" `
	}

	var param EntrustListParam

	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Limit == 0 {
		param.Limit = 5
	}
	if param.Page == 0 {
		param.Page = 1
	}
	rsp, err := rpc.InnerService.TokenService.CallEntrustList(&proto.EntrustHistoryRequest{
		Uid:   param.Uid,
		Limit: param.Limit,
		Page:  param.Page,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list", rsp.Data)
}

func (s *TokenGroup) EntrustHistory(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type EntrustListParam struct {
		Uid   uint64 `form:"uid" binding:"required"`
		Token string `form:"token" binding:"required"`
		Limit int32  `form:"limit" `
		Page  int32  `form:"page" `
	}

	var param EntrustListParam

	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Limit == 0 {
		param.Limit = 5
	}
	if param.Page == 0 {
		param.Page = 1
	}

	rsp, err := rpc.InnerService.TokenService.CallEntrustHistory(&proto.EntrustHistoryRequest{
		Uid:   param.Uid,
		Limit: param.Limit,
		Page:  param.Page,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list", rsp.Data)
}

func (s *TokenGroup) TokenBalance(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type TokenBalanceParam struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		TokenId int32  `form:"token_id" binding:"required"`
	}

	var param TokenBalanceParam

	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallTokenBalance(&proto.TokenBalanceRequest{
		Uid:     param.Uid,
		TokenId: param.TokenId,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("balance", rsp.Balance)
}

func (s *TokenGroup) TokenBalanceList(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type TokenBalanceListParam struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		TokenId int32  `form:"token_id"`
		NoZero  bool   `form:"no_zero"`
	}

	var param TokenBalanceListParam

	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallTokenBalanceList(&proto.TokenBalanceListRequest{
		Uid:     param.Uid,
		NoZero:  param.NoZero,
		TokenId: param.TokenId,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list", rsp.ListData)
}
