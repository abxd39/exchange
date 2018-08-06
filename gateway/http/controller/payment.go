package controller

import (
	//"digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"

	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type RspBankPay struct {
	Uid       uint64 `form:"uid"        json:"uid"         binding:"required"`
	Name      string `form:"name"       json:"name"        binding:"required"`
	CardNum   string `form:"card_num"   json:"card_num"    binding:"required"`
	VerifyNum string `form:"verify_num" json:"verify_num"  binding:"required"`
	BankName  string `form:"bank_name"  json:"bank_name"   binding:"required"`
	BankInfo  string `form:"bank_info"  json:"bank_info"   binding:"required"`
}
type BankPay struct {
	RspBankPay
	Verify string `form:"verify"     json:"verify"      binding:"required"`
}

type RspAliPay struct {
	Uid          uint64 `form:"uid"          json:"uid"     binding:"required"`
	Name         string `form:"name"         json:"name"    binding:"required"`
	Alipay       string `form:"alipay"       json:"alipay"  binding:"required"`
	Receipt_code string `form:"receipt_code" json:"receipt_code" binding:"required"`
}

type AliPay struct {
	RspAliPay
	Verify string `form:"verify"       json:"verify"  binding:"required"`
}

type RspPaypalPay struct {
	Uid    uint64 `form:"uid"       json:"uid"     binding:"required"`
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
	Verify string `form:"verify"       json:"verify"  binding:"required"`
}

func (*CurrencyGroup) BankPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var req BankPay
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	fmt.Printf("%#v\n", req)
	rsp, err := rpc.InnerService.CurrencyService.CallBankPay(&proto.BankPayRequest{
		Uid: req.Uid,
		//Token:     req.Token,
		Name:      req.Name,
		CardNum:   req.CardNum,
		VerifyNum: req.Verify,
		BankName:  req.BankName,
		BankInfo:  req.BankInfo,
		Verify:    req.Verify,
	})

	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	type RRRBankPay struct {
		Uid      uint64 `form:"uid"        json:"uid"         `
		Name     string `form:"name"       json:"name"        `
		CardNum  string `form:"card_num"   json:"card_num"    `
		BankName string `form:"bank_name"  json:"bank_name"   `
		BankInfo string `form:"bank_info"  json:"bank_info"   `
	}
	rrrrbankpay := new(RRRBankPay)
	rrrrbankpay.Uid = req.Uid
	rrrrbankpay.Name = req.Name
	rrrrbankpay.BankInfo = req.BankInfo
	rrrrbankpay.BankName = req.BankName
	ret.SetDataSection("bank_pay", rrrrbankpay)
	ret.SetErrCode(rsp.Code)
	return

}

func (this *CurrencyGroup) GetBankPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid uint64 `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetBankPay(&proto.PayRequest{
		Uid: req.Uid,
	})
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	type RRRBankPay struct {
		Uid      uint64 `form:"uid"        json:"uid"         `
		Name     string `form:"name"       json:"name"        `
		CardNum  string `form:"card_num"   json:"card_num"    `
		BankName string `form:"bank_name"  json:"bank_name"   `
		BankInfo string `form:"bank_info"  json:"bank_info"   `
	}
	rrrrbankpay := new(RRRBankPay)
	err = json.Unmarshal([]byte(rsp.Data), rrrrbankpay)
	if err != nil {
		fmt.Println(err.Error())
	}
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	ret.SetDataSection("bank_pay", rrrrbankpay)
	return
}

func (this *CurrencyGroup) UpdateBankPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var req BankPay
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallUpdateBankPay(&proto.BankPayRequest{
		Uid: req.Uid,
		//Token:     req.Token,
		Name:      req.Name,
		CardNum:   req.CardNum,
		VerifyNum: req.Verify,
		BankName:  req.BankName,
		BankInfo:  req.BankInfo,
		Verify:    req.Verify,
	})
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	} else {
		type RRRBankPay struct {
			Uid      uint64 `form:"uid"        json:"uid"         `
			Name     string `form:"name"       json:"name"        `
			CardNum  string `form:"card_num"   json:"card_num"    `
			BankName string `form:"bank_name"  json:"bank_name"   `
			BankInfo string `form:"bank_info"  json:"bank_info"   `
		}
		rrrrbankpay := new(RRRBankPay)
		rrrrbankpay.Uid = req.Uid
		rrrrbankpay.Name = req.Name
		rrrrbankpay.BankInfo = req.BankInfo
		rrrrbankpay.BankName = req.BankName
		ret.SetDataSection("bank_pay", rrrrbankpay)
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}
	return
}

////////////   bank end ///////////////

func (this *CurrencyGroup) Alipay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Uid         uint64 `form:"uid"          json:"uid"     binding:"required"`
		Name        string `form:"name"         json:"name"    binding:"required"`
		Alipay      string `form:"alipay"       json:"alipay"  binding:"required"`
		Verify      string `form:"verify"       json:"verify"  binding:"required"`
		ReceiptCode string `form:"receipt_code"         json:"receipt_code" binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	url, err := this.upload_picture(req.ReceiptCode)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UOPLOA_FAILED, err.Error())
		return
	}
	if url == `` {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UOPLOA_FAILED, err.Error())
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallAliPay(&proto.AlipayRequest{
		Uid:         req.Uid,
		Name:        req.Name,
		Alipay:      req.Alipay,
		ReceiptCode: url,
		Verify:      req.Verify,
	})
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	} else {
		var alipay RspAliPay
		alipay.Uid = req.Uid
		alipay.Name = req.Name
		alipay.Alipay = req.Alipay
		alipay.Receipt_code = url
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
		//fmt.Println(alipay)
		ret.SetDataSection("ali_pay", alipay)
		return
	}
}

func (this *CurrencyGroup) GetAliPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid uint64 `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetAliPay(&proto.PayRequest{
		Uid: req.Uid,
	})
	var alipay RspAliPay
	err = json.Unmarshal([]byte(rsp.Data), &alipay)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	} else {
		ret.SetDataSection("ali_pay", alipay)
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}

	return
}

func (this *CurrencyGroup) UpdateAliPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Uid         uint64 `form:"uid"          json:"uid"     binding:"required"`
		Name        string `form:"name"         json:"name"    binding:"required"`
		Alipay      string `form:"alipay"       json:"alipay"  binding:"required"`
		Verify      string `form:"verify"       json:"verify"  binding:"required"`
		ReceiptCode string `form:"receipt_code"         json:"receipt_code" binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	url, err := this.upload_picture(req.ReceiptCode)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UOPLOA_FAILED, err.Error())
		return
	}
	if url == `` {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UOPLOA_FAILED, err.Error())
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallUpdateAliPay(&proto.AlipayRequest{
		Uid: req.Uid,
		//Token:      req.Token,
		Name:        req.Name,
		Alipay:      req.Alipay,
		ReceiptCode: url,
		Verify:      req.Verify,
	})
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	} else {
		var alipay RspAliPay
		alipay.Uid = req.Uid
		alipay.Name = req.Name
		alipay.Alipay = req.Alipay
		alipay.Receipt_code = url
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
		//fmt.Println(alipay)
		ret.SetDataSection("ali_pay", alipay)
		return
	}

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
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallPaypal(&proto.PaypalRequest{
		Uid: req.Uid,
		//Token:  req.Token,
		Paypal: req.Paypal,
		Verify: req.Verify,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}

	var paypay RspPaypalPay
	paypay.Paypal = req.Paypal
	paypay.Uid = req.Uid
	ret.SetDataSection("paypal_pay", paypay)
	ret.SetErrCode(rsp.Code)
	return
}

func (this *CurrencyGroup) GetPaypal(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid uint64 `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetPaypal(&proto.PayRequest{
		Uid: req.Uid,
	})
	var paypay RspPaypalPay
	err = json.Unmarshal([]byte(rsp.Data), &paypay)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
	} else {
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
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallUpdatePaypal(&proto.PaypalRequest{
		Uid: req.Uid,
		//Token:  req.Token,
		Paypal: req.Paypal,
		//Verify: req.Verify,
	})
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
	} else {
		var paypay RspPaypalPay
		paypay.Uid = req.Uid
		paypay.Paypal = req.Paypal
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
		ret.SetDataSection("paypal_pay", paypay)
	}
	return
}

//////////////////////// paypal end ////////////////////////////////////

func (this *CurrencyGroup) WeChatPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Uid         uint64 `form:"uid"          json:"uid"     binding:"required"`
		Name        string `form:"name"         json:"name"    binding:"required"`
		Wechat      string `form:"wechat"       json:"wechat"  binding:"required"`
		Verify      string `form:"verify"       json:"verify"  binding:"required"`
		ReceiptCode string `form:"receipt_code"         json:"receipt_code" binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	url, err := this.upload_picture(req.ReceiptCode)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UOPLOA_FAILED, err.Error())
		return
	}
	if url == `` {
		ret.SetErrCode(ERRCODE_UOPLOA_FAILED, err.Error())
		return
	}

	rsp, err := rpc.InnerService.CurrencyService.CallWeChatPay(&proto.WeChatPayRequest{
		Uid:         req.Uid,
		Name:        req.Name,
		Wechat:      req.Wechat,
		ReceiptCode: url,
		Verify:      req.Verify,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	if rsp.Code == ERRCODE_SUCCESS {
		var wechatPay RspWeChatPay
		wechatPay.Receipt_code = url
		wechatPay.Name = req.Name
		wechatPay.Uid = req.Uid
		wechatPay.Wechat = req.Wechat
		ret.SetDataSection("wechat_pay", wechatPay)
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	} else {
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}
	return
}

func (this *CurrencyGroup) GetWeChatPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid uint64 `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallGetWeChatPay(&proto.PayRequest{
		Uid: req.Uid,
	})
	var wechatPay RspWeChatPay
	err = json.Unmarshal([]byte(rsp.Data), &wechatPay)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	} else {
		ret.SetDataSection("wechat_pay", wechatPay)
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}
	return
}

func (this *CurrencyGroup) UpdateWeChatPay(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Uid         uint64 `form:"uid"          json:"uid"     binding:"required"`
		Name        string `form:"name"         json:"name"    binding:"required"`
		Wechat      string `form:"wechat"       json:"wechat"  binding:"required"`
		Verify      string `form:"verify"       json:"verify"  binding:"required"`
		ReceiptCode string `form:"receipt_code"         json:"receipt_code" binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	url, err := this.upload_picture(req.ReceiptCode)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UOPLOA_FAILED, err.Error())
		return
	}
	if url == `` {
		ret.SetErrCode(ERRCODE_UOPLOA_FAILED, err.Error())
		return
	}
	rsp, err := rpc.InnerService.CurrencyService.CallUpdateWeChatPay(&proto.WeChatPayRequest{
		Uid:         req.Uid,
		Name:        req.Name,
		Wechat:      req.Wechat,
		ReceiptCode: url,
		Verify:      req.Verify,
	})
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	if rsp.Code == ERRCODE_SUCCESS {
		var wechatPay RspWeChatPay
		wechatPay.Receipt_code = url
		wechatPay.Name = req.Name
		wechatPay.Uid = req.Uid
		wechatPay.Wechat = req.Wechat
		ret.SetDataSection("wechat_pay", wechatPay)
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	} else {
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}

	return
}

//上传Ali coud
func (a *CurrencyGroup) upload_picture(file string) (string, error) {
	var remoteurl string = "https://sdun.oss-cn-shenzhen.aliyuncs.com/"
	client, err := oss.New("http://oss-cn-shenzhen.aliyuncs.com", "LTAIcJgRedhxruPq", "d7p6tWRfy0B2QaRXk7q4mb5seLROtb")
	if err != nil {
		// HandleError(err)
		panic(err)
	}
	bucket, err := client.Bucket("sdun")
	if err != nil {
		return "", err
	}
	//查找base64
	fmt.Println("base34-1")
	base := strings.Index(file, ";base64,")
	if base < 0 {
		fmt.Println("base34-3")
		// 是远程的oss 文件路径
		return file, nil
	}
	fmt.Println("base34-2")
	//fmt.Println(file)
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	subm := strings.IndexByte(file, ',')
	if subm < 0 {
		return "", errors.New("find failed!!")
	}
	substr := file[:subm]
	subb := strings.IndexByte(substr, '/')
	sube := strings.IndexByte(substr, ';')
	if subb < 0 || sube < 0 {
		return "", errors.New("find fail!!")
	}
	fmt.Println(subb, sube, subm)
	fSuffix := substr[subb+1 : sube]
	value := file[subm+1:]
	h := md5.New()
	tempValue := value
	tempValue += timestamp
	h.Write([]byte(tempValue)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	okey := hex.EncodeToString(cipherStr)
	//fmt.Println(okey)
	okey += "."
	okey += fSuffix
	//fmt.Printf("%#v\n", okey)
	ddd, _ := base64.StdEncoding.DecodeString(value)
	err = bucket.PutObject(okey, bytes.NewReader(ddd))
	if err != nil {
		//fmt.Println(filePath)
		return "", err
	}
	fmt.Println(remoteurl + okey)
	return remoteurl + okey, nil
}
