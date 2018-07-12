package controller

import (
	"digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"encoding/json"
)



type RspBankPay struct {
	Uid        uint64 `form:"uid"        json:"uid"         binding:"required"`
	Name       string `form:"name"       json:"name"        binding:"required"`
	Card_num   string `form:"card_num"   json:"card_num"    binding:"required"`
	Verify_num string `form:"verify_num" json:"verify_num"  binding:"required"`
	Bank_name  string `form:"bank_name"  json:"bank_name"   binding:"required"`
	Bank_info  string `form:"bank_info"  json:"bank_info"   binding:"required"`
}
type BankPay struct {
	RspBankPay
	Verify     string `form:"verify"     json:"verify"      binding:"required"`
}


type  RspAliPay struct {
	Uid           uint64 `form:"uid"          json:"uid"     binding:"required"`
	Name         string `form:"name"         json:"name"    binding:"required"`
	Alipay       string `form:"alipay"       json:"alipay"  binding:"required"`
	Receipt_code string `form:"receipt_code" json:"receipt_code" binding:"required"`	
}

type AliPay struct {
	RspAliPay
	Verify       string `form:"verify"       json:"verify"  binding:"required"`
}

type RspPaypalPay struct {
	Uid     uint64 `form:"uid"       json:"uid"     binding:"required"`
	Paypal string `form:"paypal"    json:"paypal"  binding:"required"`
}

type PaypalPay struct {
	RspPaypalPay
	Verify string `form:"verify"    json:"verify"  binding:"required"`
}

type RspWeChatPay struct {
	Uid          uint64 `form:"uid"          json:"uid"     binding:"required"`
	Name         string `form:"name"         json:"name"    binding:"required"`
	Wechat       string `form:"wechat"       json:"wechat"  binding:"required"`
	Receipt_code string `form:"receipt_code" json:"receipt_code" binding:"required"`	
}

type WeChatPay struct {
	RspWeChatPay 
	Verify       string `form:"verify"       json:"verify"  binding:"required"`
}



func (*CurrencyGroup) BankPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var req BankPay
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	fmt.Printf("%#v\n", req)
	rsp, err := rpc.InnerService.CurrencyService.CallBankPay(&proto.BankPayRequest{
		Uid:       req.Uid,
		//Token:     req.Token,
		Name:      req.Name,
		CardNum:   req.Card_num,
		VerifyNum: req.Verify_num,
		BankName:  req.Bank_name,
		BankInfo:  req.Bank_info,
		Verify:    req.Verify,
	})

	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Code)
	return

}

func (this *CurrencyGroup) GetBankPay (c *gin.Context) {
	ret := NewPublciError()
	defer func(){
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid      uint64 `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetBankPay(&proto.PayRequest{
		Uid:req.Uid,
	})
	if err != nil {
		log.Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	var bankpay RspBankPay
	err = json.Unmarshal([]byte(rsp.Data), &bankpay)
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	ret.SetDataSection("bank_pay", bankpay)
	return
}


func (this *CurrencyGroup) UpdateBankPay(c *gin.Context) {
	ret := NewPublciError()
	defer func(){
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var req BankPay
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallUpdateBankPay(&proto.BankPayRequest{
		Uid:       req.Uid,
		//Token:     req.Token,
		Name:      req.Name,
		CardNum:   req.Card_num,
		VerifyNum: req.Verify_num,
		BankName:  req.Bank_name,
		BankInfo:  req.Bank_info,
		Verify:    req.Verify,
	})
	if err != nil {
		log.Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}else{
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
		//ret.SetDataSection("")
	}
	return
}


////////////   bank end ///////////////



func (py *CurrencyGroup) Alipay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var req AliPay
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	rsp, err := rpc.InnerService.CurrencyService.CallAliPay(&proto.AlipayRequest{
		Uid:         req.Uid,
		Name:        req.Name,
		Alipay:      req.Alipay,
		ReceiptCode: req.Receipt_code,
	})
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
		return
	}
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	return
}


func (this *CurrencyGroup) GetAliPay(c *gin.Context) {
	ret := NewPublciError()
	defer func(){
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid      uint64 `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetAliPay(&proto.PayRequest{
		Uid:req.Uid,
	})
	var alipay RspAliPay
	err = json.Unmarshal([]byte(rsp.Data), &alipay)
	if err != nil {
		log.Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}else{
		ret.SetDataSection("ali_pay", alipay)
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))

	}

	return
}

func (this *CurrencyGroup) UpdateAliPay (c *gin.Context){
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var req AliPay
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallUpdateAliPay(&proto.AlipayRequest{
		Uid:         req.Uid,
		//Token:      req.Token,
		Name:        req.Name,
		Alipay:      req.Alipay,
		ReceiptCode: req.Receipt_code,
		Verify:      req.Verify,
	})
	if err != nil {
		log.Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}else{
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}

	return
}


///////////////////// ali pay  end /////////////////////



func (py *CurrencyGroup) Paypal(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var req PaypalPay
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallPaypal(&proto.PaypalRequest{
		Uid:    req.Uid,
		//Token:  req.Token,
		Paypal: req.Paypal,
		Verify: req.Verify,
	})
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetErrCode(rsp.Code)
	return
}


func (this *CurrencyGroup) GetPaypal(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid      uint64 `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetPaypal(&proto.PayRequest{
		Uid:req.Uid,
	})
	var paypay RspPaypalPay
	err = json.Unmarshal([]byte(rsp.Data), &paypay)
	if err != nil {
		log.Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
	}else{
		ret.SetDataSection("paypal_pay", paypay)
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}
	return
}

func (this *CurrencyGroup) UpdatePaypal(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var req PaypalPay
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallUpdatePaypal(&proto.PaypalRequest{
		Uid:    req.Uid,
		//Token:  req.Token,
		Paypal: req.Paypal,
		//Verify: req.Verify,
	})
	if err != nil {
		log.Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
	}else{
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))

	}
	return
}





//////////////////////// paypal end ////////////////////////////////////


func (py *CurrencyGroup) WeChatPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var req WeChatPay
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallWeChatPay(&proto.WeChatPayRequest{
		Uid:         req.Uid,
		Name:        req.Name,
		Wechat:      req.Wechat,
		ReceiptCode: req.Receipt_code,
	})
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	// ret.SetDataSection("data", rsp.Data)
	ret.SetErrCode(rsp.Code)
	return
}


func (this *CurrencyGroup) GetWeChatPay (c *gin.Context){
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid      uint64 `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetWeChatPay(&proto.PayRequest{
		Uid:req.Uid,
	})
	var wechatPay RspWeChatPay
	err = json.Unmarshal([]byte(rsp.Data), &wechatPay)
	if err != nil {
		log.Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}else{
		ret.SetDataSection("wechat_pay", wechatPay)
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}
	return
}



func (this *CurrencyGroup) UpdateWeChatPay (c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var req WeChatPay
	err := c.ShouldBind(&req)
	if err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallUpdateWeChatPay(&proto.WeChatPayRequest{
		Uid:         req.Uid,
		Name:        req.Name,
		Wechat:      req.Wechat,
		ReceiptCode: req.Receipt_code,
	})
	if err != nil {
		log.Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}else{
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}

	return
}
