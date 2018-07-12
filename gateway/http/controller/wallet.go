package controller

import (
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type WalletGroup struct {
}

func (this *WalletGroup) Router(router *gin.Engine) {

	r := router.Group("wallet")
	r.POST("create", this.Create)       // 创建钱包
	r.POST("signtx", this.Signtx)       // 签名
	r.POST("sendrawtx", this.SendRawTx) // 广播
	r.POST("tibi", this.Tibi)           //
	r.GET("getvalue", this.GetValue)   // 查询链上余额
	r.POST("address/save", this.AddressSave)
	r.GET("address/list", this.AddressList)
	r.POST("address/delete", this.AddressDelete)

	//btc signtx
	r.POST("signtx_btc", this.BtcSigntx) // btc 签名
	r.POST("biti_btc", this.BtcTiBi)     // btc
}

///////////////////////// start btc ///////////////////////////

func (this *WalletGroup) BtcSigntx(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Uid     int32  `form:"uid"          json:"uid"        binding:"required"` // 用户 id
		TokenId int32  `form:"token_id"     json:"token_id"   binding:"required"` // 币种ID
		Address string `form:"address"      json:"address"    binding:"required"` // 要发送给的地址
		Amount  string `form:"amount"       json:"amount"     binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		fmt.Println(err.Error())
		Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.WalletSevice.CallBtcSigntx(&proto.BtcSigntxRequest{
		Uid:     param.Uid,
		Tokenid: param.TokenId,
		Address: param.Address,
		Amount:  param.Amount,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	} else {
		fmt.Println(rsp.Data)
		//ret.SetDataValue(""rsp.Data)
		ret.SetDataSection("txhash", rsp.Data)
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}

}

func (this *WalletGroup) BtcTiBi(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Uid     int32  `form:"uid"      json:"uid"       binding:"required"`
		TokenId int32  `form:"token_id" json:"token_id"  binding:"required"`
		To      string `form:"to"       json:"to"        binding:"required"`
		Amount  string `form:"amount"   json:"amount"    binding:"required"`
		//Gasprice int32  `form:"gasprice" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		fmt.Println(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.WalletSevice.CallBtcTibi(&proto.BtcTibiRequest{
		Uid:     param.Uid,
		Tokenid: param.TokenId,
		To:      param.To,
		Amount:  param.Amount,
	})
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetErrCode(int32(rsp.Code), GetErrorMessage(int32(rsp.Code)))
	ret.SetDataSection("txhash", rsp.Data)
	//ret.SetDataValue(rsp.Data)
	//ctx.JSON(http.StatusOK, rsp)
	return
}

///////////////////////// end btc ///////////////////////////

func (this *WalletGroup) Index(ctx *gin.Context) {

}
func (this *WalletGroup) Create(ctx *gin.Context) {
	userid, _ := strconv.Atoi(ctx.Query("uid"))
	tokenid, _ := strconv.Atoi(ctx.Query("token_id"))

	rsp, err := rpc.InnerService.WalletSevice.CallCreateWallet(userid, tokenid)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
func (this *WalletGroup) Signtx(ctx *gin.Context) {
	userid, err1 := strconv.Atoi(ctx.Query("uid"))
	tokenid, err2 := strconv.Atoi(ctx.Query("token_id"))
	//to := "0x8e430b7fc9c41736911e1699dbcb6d4753cbe3b6"
	to := ctx.Query("to")
	gasprice, err3 := strconv.ParseInt(ctx.Query("gasprice"), 10, 64)
	amount := ctx.Query("amount")
	if err1 != nil || err2 != nil || err3 != nil {
		ctx.String(http.StatusOK, "参数错误")
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallSigntx(userid, tokenid, to, gasprice, amount)
	if err != nil {
		ctx.String(http.StatusOK, "err rsp")
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
func (this *WalletGroup) Update(ctx *gin.Context) {
	rsp, err := rpc.InnerService.WalletSevice.Callhello("eth")
	if err != nil {
		ctx.String(http.StatusOK, "err 0000 rsp")
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (this *WalletGroup) SendRawTx(ctx *gin.Context) {
	ret := NewPublciError()
	type Param struct {
		TokenId int32  `form:"token_id" binding:"required"`
		Signtx  string `form:"signtx" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		ctx.JSON(http.StatusOK, ret)
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallSendRawTx(param.TokenId, param.Signtx)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rsp)
	return
}

//申请提币
func (this *WalletGroup) Tibi(ctx *gin.Context) {
	ret := NewPublciError()
	//defer func() {
	//	ctx.JSON(http.StatusOK, ret.GetResult())
	//}()
	type Param struct {
		Uid      int32  `form:"uid" binding:"required"`
		Token_id int32  `form:"token_id" binding:"required"`
		To       string `form:"to" binding:"required"`
		Amount   string `form:"amount" binding:"required"`
		Gasprice int32  `form:"gasprice" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		ctx.JSON(http.StatusOK, ret)
		return
	}
	rsp, err := rpc.InnerService.WalletSevice.CallTibi(param.Uid, param.Token_id, param.To, param.Gasprice, param.Amount)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rsp)
	return
}

//查询钱包链上额度
func (this *WalletGroup) GetValue(ctx *gin.Context) {
	ret := NewPublciError()
	type Param struct {
		Uid      int32 `form:"uid" binding:"required"`
		Token_id int32 `form:"token_id" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		ctx.JSON(http.StatusOK, ret)
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallGetValue(param.Uid, param.Token_id)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rsp)
	return
}

func (this *WalletGroup) AddressSave(ctx *gin.Context) {
	ret := NewPublciError()

	type Param struct {
		Uid      int32  `form:"uid" binding:"required"`
		Token_id int32  `form:"token_id" binding:"required"`
		Address  string `form:"address" binding:"required"`
		Mark     string `form:"mark" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {

		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		ctx.JSON(http.StatusOK, ret.GetResult())
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallAddressSave(param.Uid, param.Token_id, param.Address, param.Mark)
	if err != nil {

		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}
func (this *WalletGroup) AddressList(ctx *gin.Context) {
	ret := NewPublciError()
	type Param struct {
		Uid      int32  `form:"uid"      binding:"required"`
		//Token_id int32  `form:"token_id" binding:"required"`
		//Address  string `form:"address"  binding:"required"`
		//Mark     string `form:"mark"     binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		ctx.JSON(http.StatusOK, ret.GetResult())
		return
	}

	//rsp, err := rpc.InnerService.WalletSevice.CallAddressList(param.Uid, param.Token_id)
	rsp, err := rpc.InnerService.WalletSevice.CallAddressList(param.Uid)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}
func (this *WalletGroup) AddressDelete(ctx *gin.Context) {
	ret := NewPublciError()
	type Param struct {
		Uid int32 `form:"uid" binding:"required"`
		Id  int32 `form:"id" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		ctx.JSON(http.StatusOK, ret.GetResult())
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallAddressDelete(param.Uid, param.Id)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rsp)
	return
}
