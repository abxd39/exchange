package controller

import (
	"digicon/common/convert"
	"digicon/gateway/rpc"
	proto "digicon/proto/rpc"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"

	"digicon/gateway/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/gorilla/websocket"
	//"time"
	. "digicon/proto/common"
	"encoding/json"
	"digicon/common/random"
)

type CurrencyGroup struct{}

func (this *CurrencyGroup) NewRouter(r *gin.Engine) {
	NCurrency := r.Group("/currency")
	{
		NCurrency.GET("/tokens", this.GetTokens)                    // 获取货币类型
		NCurrency.GET("/otc_list", this.AdsList)                    // 法币交易列表 - (广告(买卖))
		NCurrency.GET("/otc_user_list", this.AdsUserList)           // 个人法币交易列表 - (广告(买卖))
		NCurrency.POST("/user_currency_rating", this.GetUserRating) // 获取用戶评级
		NCurrency.GET("/trade_history", this.GetTradeHistory)       // 获取历史交易
		NCurrency.GET("/recent_transaction_price", this.GetRecentTransactionPrice)
	}
}

func (this *CurrencyGroup) Router(r *gin.Engine) {
	Currency := r.Group("/currency", TokenVerify)
	{
		Currency.GET("/otc", this.GetAds)                           // 获取广告(买卖)
		Currency.POST("/created_otc", this.AddAds)                  // 新增广告(买卖)
		Currency.POST("/updated_otc", this.UpdatedAds)              // 修改广告(买卖)
		Currency.POST("/updated_otc_status", this.UpdatedAdsStatus) // 修改广告(买卖)状态

		Currency.GET("/tokens_list", this.GetTokensList) // 获取货币类型列表
		Currency.GET("/pays", this.GetPays)              // 获取支付方式
		Currency.GET("/pays_list", this.GetPaysList)     // 获取支付方式列表
		Currency.POST("/created_chats", this.AddChats)   // 新增订单聊天
		Currency.GET("/chats_list", this.GetChatsList)   // 获取订单聊天列表

		//// order ////
		Currency.GET("/orders", this.OrdersList)           // 获取订单列表
		Currency.POST("/add_order", this.AddOrder)         // 添加订单
		Currency.POST("/ready_order", this.ReadyOrder)     // 待放行
		Currency.POST("/confirm_order", this.ConfirmOrder) // 确认放行
		Currency.POST("/cancel_order", this.CancelOrder)   // 取消订单
		Currency.POST("/delete_order", this.CancelOrder)   // 删除订单

		Currency.GET("/trade_detail", this.TradeDetail) //获取订单付款信息

		////payment///
		Currency.POST("/bank_pay", this.BankPay)      // 添加 bank_pay
		Currency.GET("/bank_pay", this.GetBankPay)    // 获取 bank_pay
		Currency.PUT("/bank_pay", this.UpdateBankPay) // 更新 bank_pay

		Currency.POST("/alipay", this.Alipay)      // 添加 ali_pay
		Currency.GET("/alipay", this.GetAliPay)    // 获取 ali_pay
		Currency.PUT("/alipay", this.UpdateAliPay) // 更新 ali_pay

		Currency.POST("/wechatpay", this.WeChatPay)      // 添加 wechat_pay
		Currency.GET("/wechatpay", this.GetWeChatPay)    // 获取 wechat_pay
		Currency.PUT("/wechatpay", this.UpdateWeChatPay) // 更新 wechat_pay

		Currency.POST("/paypal", this.Paypal)      // 添加 paypal
		Currency.GET("/paypal", this.GetPaypal)    // 获取 paypal
		Currency.PUT("/paypal", this.UpdatePaypal) // 更新 paypal

		// 追加
		Currency.GET("/selling_price", this.GetSellingPrice)       // 售价
		Currency.GET("/currency_balance", this.GetCurrencyBalance) // 余额

		//
		//Currency.GET("/add_user_balance", this.AddUserBalance)
		Currency.GET("/get_user_currency_detail", this.GetUserCurrencyDetail)
		Currency.GET("/get_user_currency", this.GetUserCurrency) //  获取法币账户

		Currency.GET("/get_asset_detail", this.GetAssetDetail) //  获取法币资产明细

		Currency.POST("/transfer_to_token", this.TransferToToken) // 法币划转到代币

		Currency.GET("/get_pay_set", this.GetHasSetPay)
	}
}

// 买卖(广告)
type CurrencyAds struct {
	Id        uint64  `json:"id"`
	Uid       uint64  `json:"uid"`        // 用户ID
	TypeId    uint32  `json:"type_id"`    // 类型:1出售 2购买
	TokenId   uint32  `json:"token_id"`   // 货币类型
	TokenName string  `json:"token_name"` // 货币名称
	Price     float64 `json:"price"`      // 单价
	Num       float64 `json:"num"`        // 数量
	//Premium     int32   `json:"premium"`      // 溢价
	Premium     float64 `json:"premium"`      // 溢价
	AcceptPrice float64 `json:"accept_price"` // 可接受最低[高]单价
	MinLimit    uint32  `json:"min_limit"`    // 最小限额
	MaxLimit    uint32  `json:"max_limit"`    // 最大限额
	IsTwolevel  uint32  `json:"is_twolevel"`  // 是否要通过二级认证:0不通过 1通过
	Pays        string  `json:"pays"`         // 支付方式:以 , 分隔: 1,2,3
	Remarks     string  `json:"remarks"`      // 交易备注
	Reply       string  `json:"reply"`        // 自动回复问候语
	IsUsd       uint32  `json:"is_usd"`       // 是否美元支付:0否 1是
	States      uint32  `json:"states"`       // 状态:0下架 1上架
	CreatedTime string  `json:"created_time"` // 创建时间
	UpdatedTime string  `json:"updated_time"` // 修改时间
}

// 获取广告(买卖)
func (this *CurrencyGroup) GetAds(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetData())
	}()

	// 请求的数据结构
	req := struct {
		Id           uint64 `form:"id" json:"id" binding:"required"`    // 广告ID
		FiatCurrency string `form:"fiat_currency" json:"fiat_currency"` // 指定 CNY | USD
		//Uid         uint64  `form:"uid" json:"uid"`         // 用户ID
		//TypeId      uint32  `form:"type_id" json:"type_id"` // 类型:1出售 2购买
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
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
	ret.SetDataValue(CurrencyAds{
		Id:          data.Id,
		Uid:         data.Uid,
		TypeId:      data.TypeId,
		TokenId:     data.TokenId,
		TokenName:   data.TokenName,
		Price:       convert.Int64ToFloat64By8Bit(int64(data.Price)),
		Num:         convert.Int64ToFloat64By8Bit(int64(data.Num)),
		Premium:     convert.Int64ToFloat64By8Bit(int64(data.Premium)),
		AcceptPrice: convert.Int64ToFloat64By8Bit(int64(data.AcceptPrice)),
		MinLimit:    data.MinLimit,
		MaxLimit:    data.MaxLimit,
		IsTwolevel:  data.IsTwolevel,
		Pays:        data.Pays,
		Remarks:     data.Remarks,
		Reply:       data.Reply,
		IsUsd:       data.IsUsd,
		States:      data.States,
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
		Token     string  `form:"token" json:"token" binding:"required"`       // token验证
		Uid       uint64  `form:"uid" json:"uid" binding:"required"`           // 用户ID
		TypeId    uint32  `form:"type_id" json:"type_id" binding:"required"`   // 类型:1出售 2购买
		TokenId   uint32  `form:"token_id" json:"token_id" binding:"required"` // 货币类型
		TokenName string  `form:"token_name" json:"token_name"`                // 货币名称
		Price     float64 `form:"price" json:"price" binding:"required"`       // 单价
		Num       float64 `form:"num" json:"num" binding:"required"`           // 数量
		//Premium      int32   `form:"premium" json:"premium"`                      // 溢价
		Premium      float64 `form:"premium" json:"premium"`              // 溢价
		AcceptPrice  float64 `form:"accept_price" json:"accept_price"`    // 可接受最低[高]单价
		MinLimit     uint32  `form:"min_limit" json:"min_limit"`          // 最小限额
		MaxLimit     uint32  `form:"max_limit" json:"max_limit"`          // 最大限额
		IsTwolevel   uint32  `form:"is_twolevel" json:"is_twolevel"`      // 是否要通过二级认证:0不通过 1通过
		Pays         string  `form:"pays" json:"pays" binding:"required"` // 支付方式:以 , 分隔: 1,2,3
		Remarks      string  `form:"remarks" json:"remarks"`              // 交易备注
		Reply        string  `form:"reply" json:"reply"`                  // 自动回复问候语
		IsUsd        uint32  `form:"is_usd" json:"is_usd"`                // 是否美元:0否 1是 (无用)
		FiatCurrency string  `form:"fiat_currency" json:"fiat_currency"`  // 指定 CNY | USD
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
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
		Premium:     convert.Float64ToInt64By8Bit(req.Premium),
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
		Token   string  `form:"token" json:"token" binding:"required"` // token验证
		Uid     uint64  `form:"uid" json:"uid" binding:"required"`     // 用户ID
		Id      uint64  `form:"id" json:"id" binding:"required"`       // 广告ID
		Price   float64 `form:"price" json:"price" binding:"required"` // 单价
		Num     float64 `form:"num" json:"num" binding:"required"`     // 数量
		Premium float64 `form:"premium" json:"premium"`                // 溢价
		//Premium      int32   `form:"premium" json:"premium"`                // 溢价
		AcceptPrice  float64 `form:"accept_price" json:"accept_price"`    // 可接受最低[高]单价
		MinLimit     uint32  `form:"min_limit" json:"min_limit"`          // 最小限额
		MaxLimit     uint32  `form:"max_limit" json:"max_limit"`          // 最大限额
		IsTwolevel   uint32  `form:"is_twolevel" json:"is_twolevel"`      // 是否要通过二级认证:0不通过 1通过
		Pays         string  `form:"pays" json:"pays" binding:"required"` // 支付方式:以 , 分隔: 1,2,3
		Remarks      string  `form:"remarks" json:"remarks"`              // 交易备注
		Reply        string  `form:"reply" json:"reply"`                  // 自动回复问候语
		FiatCurrency string  `form:"fiat_currency" json:"fiat_currency"`  // 指定 CNY | USD
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
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
	//fmt.Println("req:", req, req.IsTwolevel )
	// 调用 rpc 修改广告(买卖)
	code, err := rpc.InnerService.CurrencyService.CallUpdatedAds(&proto.AdsModel{
		Id:    req.Id,
		Price: uint64(convert.Float64ToInt64By8Bit(req.Price)),
		Num:   uint64(convert.Float64ToInt64By8Bit(req.Num)),
		//Premium:     req.Premium,
		Premium:     convert.Float64ToInt64By8Bit(req.Premium),
		AcceptPrice: uint64(convert.Float64ToInt64By8Bit(req.AcceptPrice)),
		MinLimit:    req.MinLimit,
		MaxLimit:    req.MaxLimit,
		IsTwolevel:  req.IsTwolevel,
		Pays:        req.Pays,
		Remarks:     req.Remarks,
		Reply:       req.Reply,
	})

	if err != nil {
		log.Errorf(err.Error())
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
		Token    string `form:"token" json:"token" binding:"required"`         // token验证
		Uid      uint64 `form:"uid" json:"uid" binding:"required"`             // 用户ID
		Id       int    `form:"id" json:"id" binding:"required"`               // 广告ID
		StatusId int    `form:"status_id" json:"status_id" binding:"required"` // 状态: 1下架 2上架 3正常(不删除) 4删除
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if req.Id <= 0 || req.StatusId <= 0 || req.StatusId >= 5 {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	// 调用 rpc 修改广告(买卖)状态
	rsp, err := rpc.InnerService.CurrencyService.CallUpdatedAdsStatus(&proto.AdsStatusRequest{
		Id:       uint64(req.Id),
		StatusId: uint32(req.StatusId),
	})
	fmt.Println("code:", rsp, err)

	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}

	if rsp.Code != 0 {
		ret.SetErrCode(int32(rsp.Code), GetErrorMessage(int32(rsp.Code)))
		return
	} else {
		ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	}

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

	Premium float64 `json:"premium"` // 溢价
	States  uint32  `json:"states"`  // 状态:0下架 1上架
}

// 法币交易列表 - (广告(买卖))
func (this *CurrencyGroup) AdsList(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetData())
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
		log.Errorf(err.Error())
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
	if req.FiatCurrency == "" {
		req.FiatCurrency = "cny"
	} else {
		req.FiatCurrency = strings.ToLower(req.FiatCurrency)
	}
	if req.FiatCurrency != "cny" && req.FiatCurrency != "usd" {
		req.FiatCurrency = "cny"
	}

	// 调用 rpc 法币交易列表 - (广告(买卖))
	data, err := rpc.InnerService.CurrencyService.CallAdsList(&proto.AdsListRequest{
		TypeId:  req.TypeId,
		TokenId: req.TokenId,
		Page:    uint32(req.Page),
		PageNum: uint32(req.PageNum),
	})
	fmt.Println("result ....")

	if err != nil {
		fmt.Println("rpc error:", err)
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	dataLen := len(data.Data)
	//fmt.Println("dataLen:", dataLen)

	// 法币交易列表 - 响应数据结构
	reaList := AdsListResponse{
		Page:    data.Page,
		PageNum: data.PageNum,
		Total:   data.Total,
	}

	if dataLen != 0 {
		reaList.List = make([]AdsListsData, dataLen)

		// 收集用户id 和 移动数据
		userList := make([]uint64, 0, dataLen)
		for i := 0; i < dataLen; i++ {
			adsLists := AdsListsData{
				Id:    data.Data[i].Id,
				Uid:   data.Data[i].Uid,
				Price: utils.PriceFiat(int64(data.Data[i].Price), req.FiatCurrency),
				//Num:      utils.NumFiat(int64(data.Data[i].Num), data.Data[i].Balance),
				Num:      convert.Int64ToFloat64By8Bit(int64(data.Data[i].Num)),
				MinLimit: data.Data[i].MinLimit,
				//MaxLimit:    uint32(utils.PriceFiatMaxLimit(int64(data.Data[i].MaxLimit), data.Data[i].Balance, int32(data.Data[i].TypeId), req.FiatCurrency)),
				MaxLimit:    data.Data[i].MaxLimit,
				Pays:        data.Data[i].Pays,
				CreatedTime: data.Data[i].CreatedTime,
				UpdatedTime: data.Data[i].UpdatedTime,
				UserName:    data.Data[i].UserName,
				UserFace:    data.Data[i].UserFace,
				UserVolume:  data.Data[i].UserVolume,
				TypeId:      data.Data[i].TypeId,
				TokenId:     data.Data[i].TokenId,
				TokenName:   data.Data[i].TokenName,

				Premium: convert.Int64ToFloat64By8Bit(data.Data[i].Premium),
				States:  data.Data[i].States,
			}
			userList = append(userList, data.Data[i].Uid)
			reaList.List[i] = adsLists
		}

		// 调用 rpc 用户头像和昵称
		//fmt.Println(userList)
		ulist, err := rpc.InnerService.UserSevice.CallGetNickName(&proto.UserGetNickNameRequest{Uid: userList})
		//fmt.Println("ulist:", ulist.User)
		if err != nil {
			fmt.Println("get user name error!", err.Error())
			log.Errorf(err.Error())
			ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
			return
		}
		if len(ulist.User) == 0 {
			ret.SetErrCode(ERRCODE_UNKNOWN, "获取用户头像和昵称失败")
			return
		}

		// 添加 用户头像和昵称
		for l := 0; l < dataLen; l++ {
			for _, u := range ulist.User {
				if reaList.List[l].Uid == u.Uid {
					reaList.List[l].UserName = u.NickName
					if u.HeadSculpture == ""{
						reaList.List[l].UserFace = random.GetRandHead()
					}else{
						reaList.List[l].UserFace = u.HeadSculpture
					}
					break
				}
			}
		}
		//fmt.Println(userList)
	}
	ret.SetDataValue(reaList)
	ret.SetErrCode(ERRCODE_SUCCESS)
}

// 个人法币交易列表 - (广告(买卖))
func (this *CurrencyGroup) AdsUserList(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetData())
	}()

	// 请求的数据结构
	req := struct {
		TypeId       uint32 `form:"type_id" json:"type_id" binding:"required"` // 类型:1出售 2购买
		Page         int    `form:"page" json:"page"`                          // 指定第几页
		PageNum      int    `form:"page_num" json:"page_num"`                  // 指定每页的记录数
		Uid          uint64 `form:"uid" json:"uid" binding:"required"`         // 用户ID
		FiatCurrency string `form:"fiat_currency" json:"fiat_currency"`        // 指定 CNY | USD
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
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
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	//fmt.Println("data:", data)
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
			States:      v.States,
			Premium:     convert.Int64ToFloat64By8Bit(v.Premium),
		}
		reaList.List = append(reaList.List, adsLists)
	}

	ret.SetDataValue(reaList)
	ret.SetErrCode(ERRCODE_SUCCESS)

}

// 获取货币类型
func (this *CurrencyGroup) GetTokens(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetData())
	}()

	// 请求的数据结构
	req := struct {
		Id uint32 `form:"id" json:"id" binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
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
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	if data.Id == 0 {
		ret.SetErrCode(ERRCODE_TOKENS_NOTEXIST)
		return
	}

	ret.SetDataValue(data)
	ret.SetErrCode(ERRCODE_SUCCESS)

}

// 获取货币类型列表
func (this *CurrencyGroup) GetTokensList(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetData())
	}()

	// 调用 rpc 获取货币类型列表
	data, err := rpc.InnerService.CurrencyService.CallCurrencyTokensList(&proto.CurrencyTokensRequest{})

	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetDataValue(data.Data)
	ret.SetErrCode(ERRCODE_SUCCESS)
}

// 获取支付方式
func (this *CurrencyGroup) GetPays(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetData())
	}()

	// 请求的数据结构
	req := struct {
		Id uint32 `form:"id" json:"id" binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
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
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	if data.Id == 0 {
		ret.SetErrCode(ERRCODE_PAYS_NOTEXIST)
		return
	}

	ret.SetDataValue(data)
	ret.SetErrCode(ERRCODE_SUCCESS)

}

// 获取支付方式列表
func (this *CurrencyGroup) GetPaysList(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetData())
	}()

	// 调用 rpc 获取支付方式列表
	data, err := rpc.InnerService.CurrencyService.CallCurrencyPaysList(&proto.CurrencyPaysRequest{})

	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetDataValue(data.Data)
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
		Token   string `form:"token"     json:"token"         binding:"required"` // token验证
		OrderId string `form:"order_id"  json:"order_id"      binding:"required"` // 订单ID
		Uid     uint64 `form:"uid"       json:"uid"           binding:"required"` // 用户ID
		Content string `form:"content"   json:"content"       binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
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
		log.Errorf(err.Error())
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
		c.JSON(http.StatusOK, ret.GetData())
	}()

	// 请求的数据结构
	req := struct {
		OrderId string `form:"order_id" json:"order_id" binding:"required"` // 订单ID
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
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
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetDataValue(data.Data)
	ret.SetErrCode(ERRCODE_SUCCESS)
}

// get 售价

func (this *CurrencyGroup) GetSellingPrice(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		TokenId uint32 `form:"token_id"  json:"token_id" binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	if req.TokenId == 0 {
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	data, err := rpc.InnerService.CurrencyService.CallGetSellingPrice(&proto.SellingPriceRequest{
		TokenId: req.TokenId,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	//fmt.Println("price:", data)
	type respPrice struct {
		Cny    float64 `json:"cny"`
		MinCny float64 `json:min_cny`
	}
	var rPrce respPrice
	err = json.Unmarshal([]byte(data.Data), &rPrce)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetDataSection("price", rPrce.Cny)
	ret.SetDataSection("min_price", rPrce.MinCny)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

/*
  get user currency 获取法币账户
*/

func (this *CurrencyGroup) GetUserCurrency(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Id  uint64 `form:"id"         json:"id" `
		Uid uint64 `form:"uid"        json:"uid"       binding:"required"`
		//TokenId uint32 `form:"token_id"   json:"token_id" `
		NoZero bool `form:"no_zero"   json:"no_zero"  `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		fmt.Println(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	rsp, err := rpc.InnerService.CurrencyService.CallGetUserCurrency(&proto.UserCurrencyRequest{
		Uid:    req.Uid,
		NoZero: req.NoZero,
	})
	if err != nil {
		log.Errorln(err.Error())
		fmt.Println(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}

	type RespBalance struct {
		Id        uint64  `json:"id"`
		Uid       uint64  `json:"uid"`
		TokenId   uint32  `json:"token_id"`
		TokenName string  `json:"token_name"`
		Address   string  `json:"address"`
		Freeze    float64 `json:"freeze"`
		Balance   float64 `json:"balance"`
		Valuation float64 `json:"valuation"`
	}
	type RespData struct {
		UCurrencyList []RespBalance
		Sum           float64 `json:"sum"`
		SumCNY        float64 `json:"sum_cny"`
	}

	var respdata RespData
	if err = json.Unmarshal([]byte(rsp.Data), &respdata); err != nil {
		log.Errorln(err.Error())
		fmt.Println(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	//fmt.Println("respdata:", respdata)
	ret.SetDataSection("list", respdata.UCurrencyList)
	ret.SetDataSection("sum", respdata.Sum)
	ret.SetDataSection("sum_cny", utils.Round2(respdata.SumCNY, 2))
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
}

func (this *CurrencyGroup) GetUserCurrencyDetail(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Id      uint64 `form:"id"      json:"id" `
		Uid     uint64 `form:"uid"     json:"uid"       binding:"required"`
		TokenId uint32 `form:"token_id" json:"token_id" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetUserCurrencyDetail(&proto.UserCurrencyRequest{
		Id:      req.Id,
		Uid:     req.Uid,
		TokenId: req.TokenId,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
	}
	ret.SetDataSection("uid", rsp.Uid)
	ret.SetDataSection("token_id", rsp.TokenId)
	ret.SetDataSection("token_name", rsp.TokenName)
	ret.SetDataSection("balance", convert.Int64ToFloat64By8Bit(rsp.Balance))
	ret.SetDataSection("freeze", convert.Int64ToFloat64By8Bit(rsp.Freeze))
	ret.SetDataSection("address", rsp.Address)
	ret.SetDataSection("valuation", rsp.Valuation)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

// get this.GetCurrencyQuota)     // 余额
func (this *CurrencyGroup) GetCurrencyBalance(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid     uint64 `form:"uid" json:"uid" binding:"required"`
		TokenId uint32 `form:"token_id"  json:"token_id" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	if req.Uid == 0 {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	// 获取当前法币账户余额
	data, err := rpc.InnerService.CurrencyService.CallGetCurrencyBalance(&proto.GetCurrencyBalanceRequest{
		Uid:     req.Uid,
		TokenId: req.TokenId,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetDataSection("balance", data.Data)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	//ret.SetDataSection("msg", GetErrorMessage(ERRCODE_SUCCESS))
	return
}

// get GetUserRating
// 获取用戶评级

func (this *CurrencyGroup) GetUserRating(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid uint64 `form:"uid"   json:"uid" binding:"required"` //
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	if req.Uid == 0 {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.CurrencyService.CallGetUserRating(&proto.GetUserRatingRequest{
		Uid: req.Uid,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	type UserCurrencyCount struct {
		Uid           uint64 `json:"uid"`
		NickName      string `json:"nick_name"`
		HeadSculpture string `json:"head_scul"`
		CreatedTime   string `json:"created_time"`

		Cancel uint32  `json:"cancel"` // 取消
		Good   float64 `json:"good"`   // 好评率

		Orders       uint32  `json:"orders"`        // 订单数
		CompleteRate float64 `json:"complete_rate"` //  完成率
		MonthRate    int64   `json:"month_rate"`    // 30日成单
		Success      int64   `json:"success"`       // 成功订单
		Failure      int64   `json:"failure"`       // 失败
		AverageTo    int64   `json:"average_to"`    // 120 分钟

		EmailAuth    int32 `json:"email_auth"`     //
		PhoneAuth    int32 `json:"phone_auth"`     //
		RealName     int32 `json:"real_name"`      //
		TwoLevelAuth int32 `json:"two_level_auth"` //
	}
	//fmt.Println("data:", rsp.Data)
	var uCurrencyCount UserCurrencyCount
	err = json.Unmarshal([]byte(rsp.Data), &uCurrencyCount)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println("uCurrencyCount:", uCurrencyCount)
	ret.SetErrCode(rsp.Code, rsp.Message)
	ret.SetDataSection("nick_name", uCurrencyCount.NickName)
	ret.SetDataSection("head_scul", uCurrencyCount.HeadSculpture)
	ret.SetDataSection("created_time", uCurrencyCount.CreatedTime)
	ret.SetDataSection("orders", uCurrencyCount.Orders)

	//ret.SetDataSection("appeal", uCurrencyCount.Success+uCurrencyCount.Failure) // 申诉
	//ret.SetDataSection("success", uCurrencyCount.Success)
	ret.SetDataSection("appeal", 0)  // 申诉
	ret.SetDataSection("success", 0) //

	ret.SetDataSection("average_to", uCurrencyCount.AverageTo)
	ret.SetDataSection("month_rate", uCurrencyCount.MonthRate)
	ret.SetDataSection("complete_rate", uCurrencyCount.CompleteRate)

	ret.SetDataSection("email_auth", uCurrencyCount.EmailAuth)
	ret.SetDataSection("phone_auth", uCurrencyCount.PhoneAuth)
	ret.SetDataSection("real_name", uCurrencyCount.RealName)
	ret.SetDataSection("two_level_auth", uCurrencyCount.TwoLevelAuth)
	return
}

/*
	资产明细
*/
func (this *CurrencyGroup) GetAssetDetail(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid     uint64 `form:"uid"       json:"uid"    binding:"required"` //
		Page    uint32 `form:"page"      json:"page"`
		PageNum uint32 `form:"page_num"  json:"page_num"  `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	if req.Uid == 0 {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetAssetDetail(&proto.GetAssetDetailRequest{
		Uid:     req.Uid,
		Page:    req.Page,
		PageNum: req.PageNum,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}

	type NewUserCurrencyHisotry struct {
		Id          int     `json:"id"                  `
		Uid         int32   `json:"uid"               `
		TradeUid    int32   `json:"trade_uid"         `
		TokenId     int     `json:"token_id"            `
		TokenName   string  `json:"token_name"`
		Num         float64 `json:"num"                 `
		Operator    int     `json:"operator"            `
		CreatedTime string  `json:"created_time"        `
		TradeName   string  `json:"trade_name"         `
	}

	type OldUserTotalHistory struct {
		NewList []NewUserCurrencyHisotry
		Total   int64  `json:"total"`
		Page    uint32 `json:"page"`
		PageNum uint32 `json:"page_num"`
	}

	var oldData OldUserTotalHistory
	err = json.Unmarshal([]byte(rsp.Data), &oldData)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}

	ret.SetDataSection("list", oldData.NewList)
	ret.SetDataSection("total", oldData.Total)
	ret.SetDataSection("page", oldData.Page)
	ret.SetDataSection("page_num", oldData.PageNum)
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	return
}

/*
	法币划转到代币
*/
func (this *CurrencyGroup) TransferToToken(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Uid     uint64  `form:"uid"        json:"uid"           binding:"required"`
		TokenId uint32  `form:"token_id"   json:"token_id"      binding:"required"`
		Num     float64 `form:"num"        json:"num"           binding:"required"`
	}{}

	err := c.ShouldBind(&req)

	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	if req.Uid == 0 {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	//fmt.Println(req)
	rsp, err := rpc.InnerService.CurrencyService.CallTransferToToken(&proto.TransferToTokenRequest{
		Uid:     req.Uid,
		TokenId: req.TokenId,
		Num:     uint64(convert.Float64ToInt64By8Bit(req.Num)),
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetErrCode(rsp.Code, rsp.Message)
	return
}

///*
//	测试rpc添加用户余额
//*/
//func (this *CurrencyGroup) AddUserBalance(ctx *gin.Context) {
//	rsp, err := rpc.InnerService.CurrencyService.CallAddUserBalance(&proto.AddUserBalanceRequest{
//		Uid:     2,
//		TokenId: 2,
//		Amount:  "3.33",
//	})
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	fmt.Println("rsp:", rsp.Data, rsp.Code, rsp.Message)
//}
