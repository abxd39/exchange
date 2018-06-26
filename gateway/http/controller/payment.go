package controller

import (
	"digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*CurrencyGroup) BankPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Id         uint64 `form:"uid" json:"uid" binding:"required"`
		Token      string `form:"token" json:"token" binding:"required"`
		Name       string `form:"name" json:"name" binding:"required"`
		Card_num   string `form:"card_num" json:"card_num" binding:"required"`
		Verify_num string `form:"verify_num" json:"verify_num" binding:"required"`
		Bank_name  string `form:"bank_name" json:"bank_name" binding:"required"`
		Bank_info  string `form:"bank_info" json:"bank_info" binding:"required"`
		Phone      string `form:"phone" json:"phone" binding:"required"`
		Verify     string `form:"verify" json:"verify" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	fmt.Printf("%#v\n", req)
	rsp, err := rpc.InnerService.CurrencyService.CallBankPay(&proto.BankPayRequest{
		Uid:       req.Id,
		Token:     req.Token,
		Name:      req.Name,
		CardNum:   req.Card_num,
		VerifyNum: req.Verify_num,
		BankName:  req.Bank_name,
		BankInfo:  req.Bank_info,
		Phone:     req.Phone,
		Verify:    req.Verify,
	})

	if err != nil {
		log.Log.Errorf(err.Error())
		return
	}
	ret.SetErrCode(rsp.Code)
	return

}

func (py *CurrencyGroup) Alipay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Id           uint64 `form:"uid" json:"uid" binding:"required"`
		Token        string `form:"token" json:"token" binding:"required"`
		Name         string `form:"name" json:"name" binding:"required"`
		Alipay       string `form:"alipay" json:"alipay" binding:"required"`
		Receipt_code string `form:"receipt_code" json:"receipt_code" binding:"required"`
		Phone        string `form:"phone" json:"phone" binding:"required"`
		Verify       string `form:"verify" json:"verify" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.CurrencyService.CallAliPay(&proto.AlipayRequest{
		Uid:         req.Id,
		Token:       req.Token,
		Name:        req.Name,
		Alipay:      req.Alipay,
		ReceiptCode: req.Receipt_code,
		Phone:       req.Phone,
		Verify:      req.Verify,
	})
	if err != nil {
		log.Log.Errorf(err.Error())
		return
	}
	ret.SetErrCode(rsp.Code)
	return
}

func (py *CurrencyGroup) Paypal(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Id     uint64 `form:"uid" json:"uid" binding:"required"`
		Token  string `form:"token" json:"token" binding:"required"`
		Paypal string `form:"Paypal" json:"Paypal" binding:"required"`
		Phone  string `form:"phone" json:"phone" binding:"required"`
		Verify string `form:"verify" json:"verify" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallPaypal(&proto.PaypalRequest{
		Uid:    req.Id,
		Token:  req.Token,
		Paypal: req.Paypal,
		Phone:  req.Phone,
		Verify: req.Verify,
	})
	if err != nil {
		log.Log.Errorf(err.Error())
		return
	}
	ret.SetErrCode(rsp.Code)
	return
}

func (py *CurrencyGroup) WeChatPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Id           uint64 `form:"uid" json:"uid" binding:"required"`
		Token        string `form:"token" json:"token" binding:"required"`
		Name         string `form:"name" json:"name" binding:"required"`
		Wechat       string `form:"wechat" json:"wechat" binding:"required"`
		Receipt_code string `form:"receipt_code" json:"receipt_code" binding:"required"`
		Phone        string `form:"phone" json:"phone" binding:"required"`
		Verify       string `form:"verify" json:"verify" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallWeChatPay(&proto.WeChatPayRequest{
		Uid:         req.Id,
		Token:       req.Token,
		Name:        req.Name,
		Wechat:      req.Wechat,
		ReceiptCode: req.Receipt_code,
		Phone:       req.Phone,
		Verify:      req.Verify,
	})
	if err != nil {
		log.Log.Errorf(err.Error())
		return
	}
	// ret.SetDataSection("data", rsp.Data)
	ret.SetErrCode(rsp.Code)
	return
}
