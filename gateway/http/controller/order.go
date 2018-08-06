package controller

import (
	"digicon/common/convert"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"fmt"
	"digicon/gateway/utils"
)

type OrderRequest struct {
	Id uint64 `form:"id" json:"id"  binding:"required"` //order 表Id
}

type CancelOrderRequest struct {
	OrderRequest
	CancelType uint32 `form:"id" json:"cancel_type" binding:"required"` //取消类型: 1卖方 2 买方
	Uid        int32  `form:"id" json:"uid"  binding:"required"`        //
}

type OtherType struct {
	AdId     uint64 `form:"ad_id"   json:"ad_id"   binding:"required"`       // 广告id
	AdType   uint32 `form:"ad_type" json:"ad_type" binding:"required"`       // 广告类型：1出售 2购买
	TokenId  uint64 `form:"token_id"   json:"token_id"   binding:"required"` // 货币类型
	PayId    string `form:"pay_id"     json:"pay_id"     binding:"required"` // 支付类型
	SellId   uint64 `form:"sell_id"    json:"sell_id"    binding:"required"` // 卖家id
	SellName string `form:"sell_name"  json:"sell_name"  binding:"required"` // 卖家昵称
	BuyId    uint64 `form:"buy_id"     json:"buy_id"     binding:"required"` // 买家id
	BuyName  string `form:"buy_name"   json:"buy_name"   binding:"required"` // 买家昵称
}

type OtherOrderType struct {
	TokenName   string `form:"token_name"   json:"token_name"`
	OrderId     string `form:"order_id"     json:"order_id" `    // 订单id
	States      uint32 `form:"states"       json:"states"`       // 订单状态 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消
	PayStatus   uint32 `form:"pay_status"   json:"pay_status"`   // 支付状态 1待支付 2待放行(已支付) 3确认支付(已完成)
	CancelType  uint32 `form:"cancel_type"  json:"cancel_type"`  // 取消类型 1卖方 2 买方
	CreatedTime string `form:"created_time" json:"created_time"` //
	UpdatedTime string `form:"updated_time" json:"updated_time"`
}

type OneOrder struct {
	OtherType
	Num   float64 `form:"num"        json:"num"     binding:"required"` // 交易数量
	Price float64 `form:"price"      json:"price"   binding:"required"` // 货币类型
}

type Order struct {
	OrderRequest
	OneOrder
	TotalPrice float64 `form:"total_price"  json:"total_price"   ` //
	Fee        float64 `form:"fee"          json:"fee"  `          // 手续费用
	OtherOrderType
}

/////  转后台数据类型

type BackOrder struct {
	OrderRequest
	OtherType
	Num   int64 `form:"num"        json:"num"     binding:"required"` // 交易数量
	Price int64 `form:"price"      json:"price"   binding:"required"` // 货币类型
	Fee   int64 `form:"fee"        json:"fee"     binding:"required"` // 手续费用
	OtherOrderType
}

//type AddOrder
type AddOrder struct {
	Uid  int32   `form:"uid"       json:"uid"        binding:"required"` // 用户 id
	AdId uint64  `form:"ad_id"     json:"ad_id"      binding:"required"` // 广告id
	Num  float64 `form:"num"       json:"num"        binding:"required"` // 交易数量
	//TypeId  int32   `form:"type_id"   json:"type_id"    binding:"required"` // 当前用户交易类型
}

type BackAddOrder struct {
	Uid  int32  `form:"uid"       json:"uid"        binding:"required"` // 用户 id
	AdId uint64 `form:"ad_id"     json:"ad_id"      binding:"required"` // 广告id
	Num  int64  `form:"num"       json:"num"        binding:"required"` // 交易数量
	//TypeId  int32   `form:"type_id"   json:"type_id"    binding:"required"` // 当前用户交易类型
	//PayId    uint64 `form:"pay_id"     json:"pay_id"     binding:"required"`         // 支付类型
	//Price int64 `form:"price"      json:"price"   binding:"required"` // 货币类型
	//OtherType
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
		Uid       uint64  `form:"uid"        json:"uid"`
	}
	var param OrderListParam
	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
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
		Uid:       param.Uid,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}

	var backOrders []BackOrder
	if err = json.Unmarshal([]byte(rsp.Orders), &backOrders); err != nil {
		log.Errorln(err.Error())
	}
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	var orders []Order
	for i := 0; i < len(backOrders); i++ {
		var o Order
		bod := backOrders[i]
		o.Id = bod.Id
		o.AdId = bod.AdId
		o.AdType = bod.AdType
		o.Price = convert.Int64ToFloat64By8Bit(bod.Price)
		o.Num = convert.Int64ToFloat64By8Bit(bod.Num)
		o.Fee = convert.Int64ToFloat64By8Bit(bod.Fee)
		o.TokenId = bod.TokenId
		o.TokenName = bod.TokenName
		o.PayId = bod.PayId
		o.SellId = bod.SellId
		o.SellName = bod.SellName
		o.BuyId = bod.BuyId
		o.BuyName = bod.BuyName
		o.States = bod.States
		o.OrderId = bod.OrderId
		o.PayStatus = bod.PayStatus
		o.CancelType = bod.CancelType
		o.CreatedTime = bod.CreatedTime
		o.UpdatedTime = bod.UpdatedTime
		o.TotalPrice = convert.Int64ToFloat64By8Bit(convert.Int64MulInt64By8Bit(bod.Num, bod.Price))
		orders = append(orders, o)
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
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallCancelOrder(&proto.CancelOrderRequest{
		Id:         param.Id,
		CancelType: param.CancelType,
		Uid:        param.Uid,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
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
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallDeleteOrder(&proto.OrderRequest{
		Id: param.Id,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
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
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallReadyOrder(&proto.OrderRequest{
		Id: param.Id,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	return
}

// 确认放行
func (this CurrencyGroup) ConfirmOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	//var param OrderRequest
	param := struct {
		Id  uint64 `form:"id" json:"id"   binding:"required"` //order 表Id
		Uid int32  `form:"id" json:"uid"  binding:"required"` //
	}{}
	err := c.ShouldBind(&param)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallConfirmOrder(&proto.OrderRequest{
		Id:  param.Id,
		Uid: param.Uid,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}

	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
}

// 添加订单
func (this CurrencyGroup) AddOrder(c *gin.Context) {
	fmt.Println("add order:,....")
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var param AddOrder
	var backParam BackAddOrder
	err := c.ShouldBind(&param)

	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	fmt.Println("param: ", param)
	backParam.Uid = param.Uid
	backParam.AdId = param.AdId
	backParam.Num = convert.Float64ToInt64By8Bit(param.Num)
	orderStr, _ := json.Marshal(backParam)
	//fmt.Println("params uid:", param.Uid)
	rsp, err := rpc.InnerService.CurrencyService.CallAddOrder(&proto.AddOrderRequest{
		Order: string(orderStr),
		Uid:   param.Uid,
		//TypeId:param.TypeId,
	})
	fmt.Println("rsp:", rsp )
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	ret.SetDataSection("id", rsp.Data)
	return
}

// TradeDetail

func (this *CurrencyGroup) TradeDetail(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var param OrderRequest
	err := c.ShouldBind(&param)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetTradeDetail(&proto.TradeDetailRequest{
		Id: param.Id,
	})

	type Data struct {
		SellId     uint64 `form:"sell_id"                json:"sell_id"`
		SellName   string `form:"sell_name"              json:"sell_name"`
		BuyId      uint64 `form:"buy_id"                 json:"buy_id"`
		BuyName    string `form:"buy_name"                json:"buy_name"`
		States     uint32 `form:"states"                 json:"states"`
		ExpiryTime string `xorm:"expiry_time"            json:"expiry_time" `
		TokenId    uint64 `form:"token_id"              json:"token_id"`
		TokenName  string `form:"token_name"            json:"token_name"`

		OrderId           string `form:"order_id"               json:"order_id"`
		PayPrice          int64  `form:"pay_price"              json:"pay_price"`
		Num               int64  `form:"num"                    json:"num"`
		Price             int64  `form:"price"                  json:"price"`
		AliPayName        string `form:"alipay_name"            json:"alipay_name"`
		Alipay            string `form:"alipay"                 json:"alipay"`
		AliReceiptCode    string `form:"ali_receipt_code"       json:"ali_receipt_code"`
		BankpayName       string `form:"bankpay_name"           json:"bankpay_name"`
		CardNum           string `form:"card_num"               json:"card_num"`
		BankName          string `form:"bank_name"              json:"bank_name"`
		BankInfo          string `form:"bank_info"              json:"bank_info"`
		WechatName        string `form:"wechat_name"            json:"wechat_name"`
		Wechat            string `form:"wechat"                 json:"wechat"`
		WechatReceiptCode string `form:"wechat_receipt_code"    json:"wechat_receipt_code"`
		PaypalNum         string `form:"paypal_num"             json:"paypal_num"`
	}
	var dt Data
	if err = json.Unmarshal([]byte(rsp.Data), &dt); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
	} else {
		pay_price := utils.Round2(convert.Int64ToFloat64By8Bit(dt.PayPrice), 2)

		ret.SetDataSection("sell_id", dt.SellId)
		ret.SetDataSection("sell_name", dt.SellName)
		ret.SetDataSection("buy_id", dt.BuyId)
		ret.SetDataSection("buy_name", dt.BuyName)
		ret.SetDataSection("status", dt.States)
		ret.SetDataSection("expiry_time", dt.ExpiryTime)
		ret.SetDataSection("token_id", dt.TokenId)
		ret.SetDataSection("token_name", dt.TokenName)

		ret.SetDataSection("order_id", dt.OrderId)
		ret.SetDataSection("pay_price", pay_price)
		ret.SetDataSection("num", convert.Int64ToFloat64By8Bit(dt.Num))
		ret.SetDataSection("price", convert.Int64ToFloat64By8Bit(dt.Price))

		ret.SetDataSection("alipay_name", dt.AliPayName)
		ret.SetDataSection("alipay", dt.Alipay)
		ret.SetDataSection("ali_receipt_code", dt.AliReceiptCode)

		ret.SetDataSection("bankpay_name", dt.BankpayName)
		ret.SetDataSection("card_num", dt.CardNum)
		ret.SetDataSection("bank_info", dt.BankInfo)
		ret.SetDataSection("bank_name", dt.BankName)

		ret.SetDataSection("wechat_name", dt.WechatName)
		ret.SetDataSection("wechat", dt.Wechat)
		ret.SetDataSection("wechat_receipt_code", dt.WechatReceiptCode)

		ret.SetDataSection("paypal_num", dt.PaypalNum)
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}
}

/*
  func: GetTradeHistory
*/
func (this *CurrencyGroup) GetTradeHistory(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		StartTime string `form:"start_time"    json:"start_time"`
		EndTime   string `form:"end_time"      json:"end_time"`
		Limit     int32  `form:"limit"        json:"limit"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetTradeHistory(&proto.GetTradeHistoryRequest{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Limit:     req.Limit,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}

	type UserCurrencyHistory struct {
		Price         int64  `json:"price"           `
		Fee         int64  `json:"fee"           `
		CreatedTime string `json:"created_time"  `
	}
	type RespUserCurrencyHistory struct {
		Price         float64 `json:"price"           `
		Fee         float64 `json:"fee"           `
		CreatedTime string  `json:"created_time"  `
	}
	var uCurrencyHistoryList []UserCurrencyHistory
	err = json.Unmarshal([]byte(rsp.Data), &uCurrencyHistoryList)

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	var rspCuHistory []RespUserCurrencyHistory
	for _, v := range uCurrencyHistoryList {
		var tmp RespUserCurrencyHistory
		tmp.CreatedTime = v.CreatedTime
		tmp.Price = convert.Int64ToFloat64By8Bit(v.Price)
		tmp.Fee = convert.Int64ToFloat64By8Bit(v.Fee)
		rspCuHistory = append(rspCuHistory, tmp)
	}
	ret.SetDataSection("list", rspCuHistory)
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))

	return
}

/*
	获取近期交易价格
*/
func (this *CurrencyGroup) GetRecentTransactionPrice(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		PriceType uint32 `form:"price_type"   json:"price_type"   binding:"required"` // 货币类型
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetRecentTransactionPrice(&proto.GetRecentTransactionPriceRequest{
		PriceType: req.PriceType,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}

	type TransactionPrice struct {
		MarketPrice float64 `json:"market_price"`
		LatestPrice float64 `json:"latest_price"`
	}
	var transprice TransactionPrice
	err = json.Unmarshal([]byte(rsp.Data), &transprice)
	if err != nil {
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	ret.SetDataSection("market_price", transprice.MarketPrice)
	ret.SetDataSection("latest_price", transprice.LatestPrice)
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	return
}
