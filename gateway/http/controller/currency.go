package controller

import (
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CurrencyGroup struct{}

func (this *CurrencyGroup) Router(r *gin.Engine) {
	Currency := r.Group("/currency")
	{
		Currency.POST("/add_ads", this.AddAds)                      // 新增广告(买卖)
		Currency.POST("/updated_ads", this.UpdatedAds)              // 修改广告(买卖)
		Currency.POST("/updated_ads_status", this.UpdatedAdsStatus) // 修改广告(买卖)状态

	}
}

// 新增广告(买卖)
func (this *CurrencyGroup) AddAds(c *gin.Context) {

	ret := NewErrorMessage()

	// 请求的数据结构
	req := struct {
		Uid         uint64  `form:"uid" json:"uid" binding:"required"`               // 用户ID
		TypeId      uint32  `form:"type_id" json:"type_id" binding:"required"`       // 类型:1出售 2购买
		TokenId     uint32  `form:"token_id" json:"token_id" binding:"required"`     // 货币类型
		TokenName   string  `form:"token_name" json:"token_name" binding:"required"` // 货币名称
		Price       float64 `form:"price" json:"price" binding:"required"`           // 单价
		Num         float64 `form:"num" json:"num" binding:"required"`               // 数量
		Premium     int32   `form:"premium" json:"premium"`                          // 溢价
		AcceptPrice float64 `form:"accept_price" json:"accept_price"`                // 可接受最低[高]单价
		MinLimit    uint32  `form:"min_limit" json:"min_limit"`                      // 最小限额
		MaxLimit    uint32  `form:"max_limit" json:"max_limit"`                      // 最大限额
		IsTwolevel  uint32  `form:"is_twolevel" json:"is_twolevel"`                  // 是否要通过二级认证:0不通过 1通过
		Pays        string  `form:"pays" json:"pays" binding:"required"`             // 支付方式:以 , 分隔: 1,2,3
		Remarks     string  `form:"remarks" json:"remarks"`                          // 交易备注
		Reply       string  `form:"reply" json:"reply"`                              // 自动回复问候语
		IsUsd       uint32  `form:"is_usd" json:"is_usd"`                            // 是否美元支付:0否 1是
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		ret[ERR_CODE_RET] = ERRCODE_PARAM
		ret[ERR_CODE_MESSAGE] = GetErrorMessage(ERRCODE_PARAM)
		c.JSON(http.StatusOK, ret)
		return
	}

	// 数据过虑暂不做

	// 调用 rpc 新增广告(买卖)
	code, err := rpc.InnerService.CurrencyService.CallAddAds(&proto.AdsRequest{
		Uid:         req.Uid,
		TypeId:      req.TypeId,
		TokenId:     req.TokenId,
		TokenName:   req.TokenName,
		Price:       req.Price,
		Num:         req.Num,
		Premium:     req.Premium,
		AcceptPrice: req.AcceptPrice,
		MinLimit:    req.MinLimit,
		MaxLimit:    req.MaxLimit,
		IsTwolevel:  req.IsTwolevel,
		Pays:        req.Pays,
		Remarks:     req.Remarks,
		Reply:       req.Reply,
		IsUsd:       req.IsUsd,
		States:      1,
	})

	if err != nil || code != 0 {
		ret[ERR_CODE_RET] = ERRCODE_PARAM
		ret[ERR_CODE_MESSAGE] = GetErrorMessage(ERRCODE_PARAM)
		c.JSON(http.StatusOK, ret)
		return
	}

	ret[ERR_CODE_RET] = ERRCODE_SUCCESS
	ret[ERR_CODE_MESSAGE] = GetErrorMessage(ERRCODE_SUCCESS)
	c.JSON(http.StatusOK, ret)

}

// 修改广告(买卖)

func (this *CurrencyGroup) UpdatedAds(c *gin.Context) {

	ret := NewErrorMessage()

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
	if err != nil || req.Id == 0 {
		ret[ERR_CODE_RET] = ERRCODE_PARAM
		ret[ERR_CODE_MESSAGE] = GetErrorMessage(ERRCODE_PARAM)
		c.JSON(http.StatusOK, ret)
		return
	}

	// 数据过虑暂不做

	// 调用 rpc 修改广告(买卖)
	code, err := rpc.InnerService.CurrencyService.CallUpdatedAds(&proto.AdsRequest{
		Id:          req.Id,
		Price:       req.Price,
		Num:         req.Num,
		Premium:     req.Premium,
		AcceptPrice: req.AcceptPrice,
		MinLimit:    req.MinLimit,
		MaxLimit:    req.MaxLimit,
		IsTwolevel:  req.IsTwolevel,
		Pays:        req.Pays,
		Remarks:     req.Remarks,
		Reply:       req.Reply,
	})

	if err != nil || code != 0 {
		ret[ERR_CODE_RET] = code
		ret[ERR_CODE_MESSAGE] = GetErrorMessage(int32(code))
		c.JSON(http.StatusOK, ret)
		return
	}

	ret[ERR_CODE_RET] = ERRCODE_SUCCESS
	ret[ERR_CODE_MESSAGE] = GetErrorMessage(ERRCODE_SUCCESS)
	c.JSON(http.StatusOK, ret)

}

// 修改广告(买卖)状态
func (this *CurrencyGroup) UpdatedAdsStatus(c *gin.Context) {

	ret := NewErrorMessage()

	// 请求的数据结构
	req := struct {
		Id       uint64 `form:"id" json:"id" binding:"required"`               // 广告ID
		StatusId uint32 `form:"status_id" json:"status_id" binding:"required"` // 状态: 1下架 2上架 3正常(不删除) 4删除
	}{}

	err := c.ShouldBind(&req)
	if err != nil || req.Id == 0 || req.StatusId == 0 || req.StatusId >= 5 {
		ret[ERR_CODE_RET] = ERRCODE_PARAM
		ret[ERR_CODE_MESSAGE] = GetErrorMessage(ERRCODE_PARAM)
		c.JSON(http.StatusOK, ret)
		return
	}

	// 调用 rpc 修改广告(买卖)状态
	code, err := rpc.InnerService.CurrencyService.CallUpdatedAdsStatus(&proto.AdsStatusRequest{
		Id:       req.Id,
		StatusId: req.StatusId,
	})

	if err != nil || code != 0 {
		ret[ERR_CODE_RET] = code
		ret[ERR_CODE_MESSAGE] = GetErrorMessage(int32(code))
		c.JSON(http.StatusOK, ret)
		return
	}

	ret[ERR_CODE_RET] = ERRCODE_SUCCESS
	ret[ERR_CODE_MESSAGE] = GetErrorMessage(ERRCODE_SUCCESS)
	c.JSON(http.StatusOK, ret)

}
