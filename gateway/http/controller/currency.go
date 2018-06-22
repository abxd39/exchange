package controller

import (
	"digicon/common/convert"
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CurrencyGroup struct{}

func (this *CurrencyGroup) Router(r *gin.Engine) {
	Currency := r.Group("/currency")
	{
		Currency.GET("/ads", this.GetAds)                           // 获取广告(买卖)
		Currency.POST("/add_ads", this.AddAds)                      // 新增广告(买卖)
		Currency.POST("/updated_ads", this.UpdatedAds)              // 修改广告(买卖)
		Currency.POST("/updated_ads_status", this.UpdatedAdsStatus) // 修改广告(买卖)状态
		Currency.POST("/ads_list", this.AdsList)                    // 法币交易列表 - (广告(买卖))
		Currency.POST("/ads_user_list", this.AdsUserList)           // 个人法币交易列表 - (广告(买卖))
		Currency.GET("/tokens", this.GetTokens)                     // 获取货币类型
		Currency.GET("/tokens_list", this.GetTokensList)            // 获取货币类型列表
		Currency.GET("/pays", this.GetPays)                         // 获取支付方式
		Currency.GET("/pays_list", this.GetPaysList)                // 获取支付方式列表
		Currency.POST("/add_chats", this.AddChats)                  // 新增订单聊天
		Currency.GET("/chats_list", this.GetChatsList)              // 获取订单聊天列表

		//// order ////
		Currency.GET("/orders", this.OrdersList)           // 获取订单列表
		Currency.POST("/add_order", this.AddOrder)         // 添加订单
		Currency.POST("/ready_order", this.ReadyOrder)     // 待放行
		Currency.POST("/confirm_order", this.ConfirmOrder) // 确认放行
		Currency.POST("/cancel_order", this.CancelOrder)   // 取消订单
		Currency.POST("/delete_order", this.CancelOrder)   // 删除订单

		////payment///
		Currency.POST("/bank_pay", this.BankPay)
		Currency.POST("/alipay", this.Alipay)
		Currency.POST("/wechatpay", this.WeChatPay)
		Currency.POST("/paypal", this.Paypal)

	}
}

// 买卖(广告)
type CurrencyAds struct {
	Id          uint64  `json:"id"`
	Uid         uint64  `json:"uid"`          // 用户ID
	TypeId      uint32  `json:"type_id"`      // 类型:1出售 2购买
	TokenId     uint32  `json:"token_id"`     // 货币类型
	TokenName   string  `json:"token_name"`   // 货币名称
	Price       float64 `json:"price"`        // 单价
	Num         float64 `json:"num"`          // 数量
	Premium     int32   `json:"premium"`      // 溢价
	AcceptPrice float64 `json:"accept_price"` // 可接受最低[高]单价
	MinLimit    uint32  `json:"min_limit"`    // 最小限额
	MaxLimit    uint32  `json:"max_limit"`    // 最大限额
	IsTwolevel  uint32  `json:"is_twolevel"`  // 是否要通过二级认证:0不通过 1通过
	Pays        string  `json:"pays"`         // 支付方式:以 , 分隔: 1,2,3
	Remarks     string  `json:"remarks"`      // 交易备注
	Reply       string  `json:"reply"`        // 自动回复问候语
	IsUsd       uint32  `json:"is_usd"`       // 是否美元支付:0否 1是
	States      uint32  `json:"states"`       // 状态:0下架 1上架
	Inventory   uint64  `json:"inventory"`    // 库存
	Fee         uint64  `json:"fee"`          // 手续费用
	CreatedTime string  `json:"created_time"` // 创建时间
	UpdatedTime string  `json:"updated_time"` // 修改时间
}

// 获取广告(买卖)
func (this *CurrencyGroup) GetAds(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 请求的数据结构
	req := struct {
		Id uint64 `form:"id" json:"id" binding:"required"` // 广告ID
		//Uid         uint64  `form:"uid" json:"uid"`         // 用户ID
		//TypeId      uint32  `form:"type_id" json:"type_id"` // 类型:1出售 2购买
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	// 调用 rpc 获取广告(买卖)
	data, err := rpc.InnerService.CurrencyService.CallGetAds(&proto.AdsGetRequest{
		Id: req.Id,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	if data.Id == 0 {
		ret.SetErrCode(ERRCODE_ADS_NOTEXIST)
		return
	}

	ret.SetErrCode(ERRCODE_SUCCESS)
	ret.SetDataSection(RET_DATA, CurrencyAds{
		Id:          data.Id,
		Uid:         data.Uid,
		TypeId:      data.TypeId,
		TokenId:     data.TokenId,
		TokenName:   data.TokenName,
		Price:       convert.Int64ToFloat64By8Bit(int64(data.Price)),
		Num:         convert.Int64ToFloat64By8Bit(int64(data.Num)),
		Premium:     data.Premium,
		AcceptPrice: convert.Int64ToFloat64By8Bit(int64(data.AcceptPrice)),
		MinLimit:    data.MinLimit,
		MaxLimit:    data.MaxLimit,
		IsTwolevel:  data.IsTwolevel,
		Pays:        data.Pays,
		Remarks:     data.Remarks,
		Reply:       data.Reply,
		IsUsd:       data.IsUsd,
		States:      data.States,
		Inventory:   data.Inventory,
		Fee:         data.Fee,
		CreatedTime: data.CreatedTime,
		UpdatedTime: data.UpdatedTime,
	})

}

// 新增广告(买卖)
func (this *CurrencyGroup) AddAds(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 请求的数据结构
	req := struct {
		Uid         uint64  `form:"uid" json:"uid" binding:"required"`           // 用户ID
		TypeId      uint32  `form:"type_id" json:"type_id" binding:"required"`   // 类型:1出售 2购买
		TokenId     uint32  `form:"token_id" json:"token_id" binding:"required"` // 货币类型
		TokenName   string  `form:"token_name" json:"token_name"`                // 货币名称
		Price       float64 `form:"price" json:"price" binding:"required"`       // 单价
		Num         float64 `form:"num" json:"num" binding:"required"`           // 数量
		Premium     int32   `form:"premium" json:"premium"`                      // 溢价
		AcceptPrice float64 `form:"accept_price" json:"accept_price"`            // 可接受最低[高]单价
		MinLimit    uint32  `form:"min_limit" json:"min_limit"`                  // 最小限额
		MaxLimit    uint32  `form:"max_limit" json:"max_limit"`                  // 最大限额
		IsTwolevel  uint32  `form:"is_twolevel" json:"is_twolevel"`              // 是否要通过二级认证:0不通过 1通过
		Pays        string  `form:"pays" json:"pays" binding:"required"`         // 支付方式:以 , 分隔: 1,2,3
		Remarks     string  `form:"remarks" json:"remarks"`                      // 交易备注
		Reply       string  `form:"reply" json:"reply"`                          // 自动回复问候语
		IsUsd       uint32  `form:"is_usd" json:"is_usd"`                        // 是否美元支付:0否 1是
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if req.Uid == 0 {
		ret.SetErrCode(ERRCODE_ACCOUNT_NOTEXIST)
		return
	}

	if req.TypeId == 0 || req.TypeId >= 3 {
		ret.SetErrCode(ERRCODE_ADS_TYPE_NOTEXIST)
		return
	}

	if req.TokenId == 0 {
		ret.SetErrCode(ERRCODE_TOKENS_NOTEXIST)
		return
	}

	if req.Pays == "" {
		ret.SetErrCode(ERRCODE_PAYS_NOTEXIST)
		return
	}

	if req.Price < 0 || req.Num < 0 {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	// 检证货币类型 ==========
	// 调用 rpc 获取货币类型
	tokenData, err := rpc.InnerService.CurrencyService.CallGetCurrencyTokens(&proto.CurrencyTokensRequest{
		Id: req.TokenId,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	if tokenData.Id == 0 {
		ret.SetErrCode(ERRCODE_TOKENS_NOTEXIST)
		return
	}

	// 检证支付方式 ==========
	paysList := strings.Split(req.Pays, ",")
	// 调用 rpc 获取支付方式列表
	paysData, err := rpc.InnerService.CurrencyService.CallCurrencyPaysList(&proto.CurrencyPaysRequest{})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	if len(paysData.Data) == 0 {
		ret.SetErrCode(ERRCODE_UNKNOWN, "获取支付方式列表失败")
		return
	}
	for _, plv := range paysList {

		isPays := false
		paysId, err := strconv.Atoi(plv)
		if err != nil {
			ret.SetErrCode(ERRCODE_PAYS_NOTEXIST)
			return
		}

		for _, pdv := range paysData.Data {
			if uint32(paysId) == pdv.Id {
				isPays = true
				break
			}
		}

		if !isPays {
			ret.SetErrCode(ERRCODE_PAYS_NOTEXIST)
			return
		}

	}

	// 数据过虑暂不做

	// 调用 rpc 新增广告(买卖)
	code, err := rpc.InnerService.CurrencyService.CallAddAds(&proto.AdsModel{
		Uid:         req.Uid,
		TypeId:      req.TypeId,
		TokenId:     req.TokenId,
		TokenName:   tokenData.Name,
		Price:       uint64(convert.Float64ToInt64By8Bit(req.Price)),
		Num:         uint64(convert.Float64ToInt64By8Bit(req.Num)),
		Premium:     req.Premium,
		AcceptPrice: uint64(convert.Float64ToInt64By8Bit(req.AcceptPrice)),
		MinLimit:    req.MinLimit,
		MaxLimit:    req.MaxLimit,
		IsTwolevel:  req.IsTwolevel,
		Pays:        req.Pays,
		Remarks:     req.Remarks,
		Reply:       req.Reply,
		IsUsd:       req.IsUsd,
		States:      1,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	if code != 0 {
		ret.SetErrCode(int32(code))
		return
	}

	ret.SetErrCode(ERRCODE_SUCCESS)

}

// 修改广告(买卖)
func (this *CurrencyGroup) UpdatedAds(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 请求的数据结构
	req := struct {
		Id          uint64  `form:"id" json:"id" binding:"required"`       // 广告ID
		Price       float64 `form:"price" json:"price" binding:"required"` // 单价
		Num         float64 `form:"num" json:"num" binding:"required"`     // 数量
		Premium     int32   `form:"premium" json:"premium"`                // 溢价
		AcceptPrice float64 `form:"accept_price" json:"accept_price"`      // 可接受最低[高]单价
		MinLimit    uint32  `form:"min_limit" json:"min_limit"`            // 最小限额
		MaxLimit    uint32  `form:"max_limit" json:"max_limit"`            // 最大限额
		IsTwolevel  uint32  `form:"is_twolevel" json:"is_twolevel"`        // 是否要通过二级认证:0不通过 1通过
		Pays        string  `form:"pays" json:"pays" binding:"required"`   // 支付方式:以 , 分隔: 1,2,3
		Remarks     string  `form:"remarks" json:"remarks"`                // 交易备注
		Reply       string  `form:"reply" json:"reply"`                    // 自动回复问候语
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if req.Id == 0 {
		ret.SetErrCode(ERRCODE_ADS_NOTEXIST)
		return
	}

	if req.Pays == "" {
		ret.SetErrCode(ERRCODE_PAYS_NOTEXIST)
		return
	}

	if req.Price < 0 || req.Num < 0 {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	// 检证支付方式 ==========
	paysList := strings.Split(req.Pays, ",")
	// 调用 rpc 获取支付方式列表
	paysData, err := rpc.InnerService.CurrencyService.CallCurrencyPaysList(&proto.CurrencyPaysRequest{})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	if len(paysData.Data) == 0 {
		ret.SetErrCode(ERRCODE_UNKNOWN, "获取支付方式列表失败")
		return
	}
	for _, plv := range paysList {

		isPays := false
		paysId, err := strconv.Atoi(plv)
		if err != nil {
			ret.SetErrCode(ERRCODE_PAYS_NOTEXIST)
			return
		}

		for _, pdv := range paysData.Data {
			if uint32(paysId) == pdv.Id {
				isPays = true
				break
			}
		}

		if !isPays {
			ret.SetErrCode(ERRCODE_PAYS_NOTEXIST)
			return
		}

	}

	// 数据过虑暂不做

	// 调用 rpc 修改广告(买卖)
	code, err := rpc.InnerService.CurrencyService.CallUpdatedAds(&proto.AdsModel{
		Id:          req.Id,
		Price:       uint64(convert.Float64ToInt64By8Bit(req.Price)),
		Num:         uint64(convert.Float64ToInt64By8Bit(req.Num)),
		Premium:     req.Premium,
		AcceptPrice: uint64(convert.Float64ToInt64By8Bit(req.AcceptPrice)),
		MinLimit:    req.MinLimit,
		MaxLimit:    req.MaxLimit,
		IsTwolevel:  req.IsTwolevel,
		Pays:        req.Pays,
		Remarks:     req.Remarks,
		Reply:       req.Reply,
	})

	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	if code != 0 {
		ret.SetErrCode(int32(code))
		return
	}

	ret.SetErrCode(ERRCODE_SUCCESS)

}

// 修改广告(买卖)状态
func (this *CurrencyGroup) UpdatedAdsStatus(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 请求的数据结构
	req := struct {
		Id       int `form:"id" json:"id" binding:"required"`               // 广告ID
		StatusId int `form:"status_id" json:"status_id" binding:"required"` // 状态: 1下架 2上架 3正常(不删除) 4删除
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if req.Id <= 0 || req.StatusId <= 0 || req.StatusId >= 5 {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	// 调用 rpc 修改广告(买卖)状态
	code, err := rpc.InnerService.CurrencyService.CallUpdatedAdsStatus(&proto.AdsStatusRequest{
		Id:       uint64(req.Id),
		StatusId: uint32(req.StatusId),
	})

	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	if code != 0 {
		ret.SetErrCode(int32(code))
		return
	}

	ret.SetErrCode(ERRCODE_SUCCESS)

}

// 法币交易列表 - 响应数据结构
type AdsListResponse struct {
	Page    uint32         `json:"page"`     // 指定第几页
	PageNum uint32         `json:"page_num"` // 指定每页的记录数
	Total   uint64         `json:"total"`    // 总记录数
	List    []AdsListsData `json:"list"`
}
type AdsListsData struct {
	Id          uint64  `json:"id"`           // 广告ID
	Uid         uint64  `json:"uid"`          // 用户ID
	Price       float64 `json:"price"`        // 单价
	Num         float64 `json:"num"`          // 数量
	MinLimit    uint32  `json:"min_limit"`    // 最小限额
	MaxLimit    uint32  `json:"max_limit"`    // 最大限额
	Pays        string  `json:"pays"`         // 支付方式:以 , 分隔: 1,2,3
	CreatedTime string  `json:"created_time"` // 创建时间
	UpdatedTime string  `json:"updated_time"` // 修改时间
	UserName    string  `json:"user_name"`    // 用户名
	UserFace    string  `json:"user_face"`    // 用户头像
	UserVolume  uint32  `json:"user_volume"`  // 用户成交量
	TypeId      uint32  `json:"type_id"`      // 类型:1出售 2购买
	TokenId     uint32  `json:"token_id"`     // 货币类型
	TokenName   string  `json:"token_name"`   // 货币名称
}

// 法币交易列表 - (广告(买卖))
func (this *CurrencyGroup) AdsList(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 请求的数据结构
	req := struct {
		TypeId       uint32 `form:"type_id" json:"type_id" binding:"required"` // 类型:1出售 2购买
		TokenId      uint32 `form:"token_id" json:"token_id"`                  // 货币类型
		TokenName    string `form:"token_name" json:"token_name"`              // 货币名称
		Page         int    `form:"page" json:"page"`                          // 指定第几页
		PageNum      int    `form:"page_num" json:"page_num"`                  // 指定每页的记录数
		FiatCurrency string `form:"fiat_currency" json:"fiat_currency"`        // 指定 CNY | USD
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if req.TypeId <= 0 || req.TypeId >= 3 {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	if req.TokenId <= 0 {
		req.TokenId = 2 // 比特币
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageNum <= 0 {
		req.PageNum = 9
	}

	// 调用 rpc 法币交易列表 - (广告(买卖))
	data, err := rpc.InnerService.CurrencyService.CallAdsList(&proto.AdsListRequest{
		TypeId:  req.TypeId,
		TokenId: req.TokenId,
		Page:    uint32(req.Page),
		PageNum: uint32(req.PageNum),
	})

	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	// 法币交易列表 - 响应数据结构
	reaList := AdsListResponse{Page: data.Page, PageNum: data.PageNum, Total: data.Total}
	for _, v := range data.Data {

		adsLists := AdsListsData{
			Id:          v.Id,
			Uid:         v.Uid,
			Price:       convert.Int64ToFloat64By8Bit(int64(v.Price)),
			Num:         convert.Int64ToFloat64By8Bit(int64(v.Num)),
			MinLimit:    v.MinLimit,
			MaxLimit:    v.MaxLimit,
			Pays:        v.Pays,
			CreatedTime: v.CreatedTime,
			UpdatedTime: v.UpdatedTime,
			UserName:    v.UserName,
			UserFace:    v.UserFace,
			UserVolume:  v.UserVolume,
			TypeId:      v.TypeId,
			TokenId:     v.TokenId,
			TokenName:   v.TokenName,
		}

		reaList.List = append(reaList.List, adsLists)
	}

	ret.SetDataSection(RET_DATA, reaList)
	ret.SetErrCode(ERRCODE_SUCCESS)

}

// 个人法币交易列表 - (广告(买卖))
func (this *CurrencyGroup) AdsUserList(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 请求的数据结构
	req := struct {
		TypeId  uint32 `form:"type_id" json:"type_id" binding:"required"` // 类型:1出售 2购买
		Page    int    `form:"page" json:"page"`                          // 指定第几页
		PageNum int    `form:"page_num" json:"page_num"`                  // 指定每页的记录数
		Uid     uint64 `form:"uid" json:"uid" binding:"required"`         // 用户ID
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if req.TypeId <= 0 || req.TypeId >= 3 || req.Uid == 0 {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageNum <= 0 {
		req.PageNum = 9
	}

	// 调用 rpc 个人法币交易列表 - (广告(买卖))
	data, err := rpc.InnerService.CurrencyService.CallAdsUserList(&proto.AdsListRequest{
		Uid:     req.Uid,
		TypeId:  req.TypeId,
		Page:    uint32(req.Page),
		PageNum: uint32(req.PageNum),
	})

	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	// 法币交易列表 - 响应数据结构
	reaList := AdsListResponse{Page: data.Page, PageNum: data.PageNum, Total: data.Total}
	for _, v := range data.Data {

		adsLists := AdsListsData{
			Id:          v.Id,
			Uid:         v.Uid,
			Price:       convert.Int64ToFloat64By8Bit(int64(v.Price)),
			Num:         convert.Int64ToFloat64By8Bit(int64(v.Num)),
			MinLimit:    v.MinLimit,
			MaxLimit:    v.MaxLimit,
			Pays:        v.Pays,
			CreatedTime: v.CreatedTime,
			UpdatedTime: v.UpdatedTime,
			TypeId:      v.TypeId,
			TokenId:     v.TokenId,
			TokenName:   v.TokenName,
		}

		reaList.List = append(reaList.List, adsLists)
	}

	ret.SetDataSection(RET_DATA, reaList)
	ret.SetErrCode(ERRCODE_SUCCESS)

}

// 获取货币类型
func (this *CurrencyGroup) GetTokens(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 请求的数据结构
	req := struct {
		Id uint32 `form:"id" json:"id" binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if req.Id == 0 {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	// 调用 rpc 获取货币类型
	data, err := rpc.InnerService.CurrencyService.CallGetCurrencyTokens(&proto.CurrencyTokensRequest{
		Id: req.Id,
	})

	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	if data.Id == 0 {
		ret.SetErrCode(ERRCODE_TOKENS_NOTEXIST)
		return
	}

	ret.SetDataSection(RET_DATA, data)
	ret.SetErrCode(ERRCODE_SUCCESS)

}

// 获取货币类型列表
func (this *CurrencyGroup) GetTokensList(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 调用 rpc 获取货币类型列表
	data, err := rpc.InnerService.CurrencyService.CallCurrencyTokensList(&proto.CurrencyTokensRequest{})

	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetDataSection(RET_DATA, data)
	ret.SetErrCode(ERRCODE_SUCCESS)
}

// 获取支付方式
func (this *CurrencyGroup) GetPays(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 请求的数据结构
	req := struct {
		Id uint32 `form:"id" json:"id" binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if req.Id == 0 {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	// 调用 rpc 获取支付方式
	data, err := rpc.InnerService.CurrencyService.CallGetCurrencyPays(&proto.CurrencyPaysRequest{
		Id: req.Id,
	})

	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	if data.Id == 0 {
		ret.SetErrCode(ERRCODE_PAYS_NOTEXIST)
		return
	}

	ret.SetDataSection(RET_DATA, data)
	ret.SetErrCode(ERRCODE_SUCCESS)

}

// 获取支付方式列表
func (this *CurrencyGroup) GetPaysList(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 调用 rpc 获取支付方式列表
	data, err := rpc.InnerService.CurrencyService.CallCurrencyPaysList(&proto.CurrencyPaysRequest{})

	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetDataSection(RET_DATA, data)
	ret.SetErrCode(ERRCODE_SUCCESS)
}

// 新增订单聊天
func (this *CurrencyGroup) AddChats(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 请求的数据结构
	req := struct {
		OrderId string `form:"order_id" json:"order_id" binding:"required"` // 订单ID
		Uid     uint64 `form:"uid" json:"uid" binding:"required"`           // 用户ID
		Content string `form:"content" json:"content" binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if req.OrderId == "" || req.Uid == 0 || req.Content == "" {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	// 验证订单和获取订单
	// 验证用户和获取用户名

	// 调用 rpc 新增订单聊天
	code, err := rpc.InnerService.CurrencyService.CallGetCurrencyChats(&proto.CurrencyChats{
		OrderId: req.OrderId,
		Uid:     req.Uid,
		Content: req.Content,
	})

	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	if code != 0 {
		ret.SetErrCode(int32(code))
		return
	}

	ret.SetErrCode(ERRCODE_SUCCESS)

}

// 获取订单聊天列表
func (this *CurrencyGroup) GetChatsList(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	// 请求的数据结构
	req := struct {
		OrderId string `form:"order_id" json:"order_id" binding:"required"` // 订单ID
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if req.OrderId == "" {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	// 调用 rpc 获取订单聊天列表
	data, err := rpc.InnerService.CurrencyService.CallCurrencyChatsList(&proto.CurrencyChats{
		OrderId: req.OrderId,
	})

	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetDataSection(RET_DATA, data)
	ret.SetErrCode(ERRCODE_SUCCESS)
}
