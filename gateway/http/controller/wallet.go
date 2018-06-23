package controller

import (
	"digicon/gateway/rpc"
	"github.com/gin-gonic/gin"
	. "digicon/proto/common"
	"net/http"
	"strconv"
)

type WalletGroup struct {

}

func (this *WalletGroup) Router(router *gin.Engine){

	r := router.Group("wallet")
	r.GET("create",this.Create)
	r.GET("signtx",this.Signtx)
	r.GET("tibi",this.Tibi)
	r.GET("getvalue",this.GetValue)
	r.GET("address/save",this.AddressSave)
	r.GET("address/list",this.AddressList)
	r.GET("address/delete",this.AddressDelete)

}

func (this *WalletGroup)Index(ctx *gin.Context){

}
func (this *WalletGroup)Create(ctx *gin.Context){
	userid ,_	:= strconv.Atoi(  ctx.Query("uid"))
	tokenid ,_	:= strconv.Atoi(   ctx.Query("token_id"))

	rsp, err := rpc.InnerService.WalletSevice.CallCreateWallet(userid,tokenid)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
func (this *WalletGroup)Signtx(ctx *gin.Context){
	userid ,err1	:= strconv.Atoi(  ctx.Query("uid"))
	tokenid ,err2	:= strconv.Atoi(   ctx.Query("token_id"))
	//to := "0x8e430b7fc9c41736911e1699dbcb6d4753cbe3b6"
	to := ctx.Query("to")
	gasprice,err3 := strconv.Atoi( ctx.Query("gasprice"))
	amount := ctx.Query("amount")
	if err1 != nil || err2!=nil || err3!=nil {
		ctx.String(http.StatusOK, "参数错误")
		return
	}


	rsp, err := rpc.InnerService.WalletSevice.CallSigntx(userid,tokenid,to,gasprice,amount)
	if err != nil {
		ctx.String(http.StatusOK, "err rsp")
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
func (this *WalletGroup)Update(ctx *gin.Context){
	rsp, err := rpc.InnerService.WalletSevice.Callhello("eth")
	if err != nil {
		ctx.String(http.StatusOK, "err 0000 rsp")
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}


func (this *WalletGroup)Tibi(ctx *gin.Context){
	ret := NewPublciError()
	//defer func() {
	//	ctx.JSON(http.StatusOK, ret.GetResult())
	//}()
	type Param struct {
		Uid       int32 `form:"uid" binding:"required"`
		Token_id  int32 `form:"token_id" binding:"required"`
		To		  string `form:"to" binding:"required"`
		Amount    string `form:"amount" binding:"required"`
		Gasprice  int32 `form:"gasprice" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.WalletSevice.CallTibi(param.Uid,param.Token_id,param.To,param.Gasprice,param.Amount)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rsp)
	return
}

func (this *WalletGroup)GetValue(ctx *gin.Context){
	ret := NewPublciError()
	type Param struct {
		Uid       int32 `form:"uid" binding:"required"`
		Token_id  int32 `form:"token_id" binding:"required"`

	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		ctx.JSON(http.StatusOK, ret)
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallGetValue(param.Uid,param.Token_id)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rsp)
	return
}
func (this *WalletGroup)AddressSave(ctx *gin.Context){
	ret := NewPublciError()

	type Param struct {
		Uid       int32 `form:"uid" binding:"required"`
		Token_id  int32 `form:"token_id" binding:"required"`
		Address   string `form:"address" binding:"required"`
		Mark  	  string `form:"mark" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {

		ret.SetErrCode(ERRCODE_PARAM,err.Error())
		ctx.JSON(http.StatusOK, ret.GetResult())
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallAddressSave(param.Uid,param.Token_id,param.Address,param.Mark)
	if err != nil {

		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}
func (this *WalletGroup)AddressList(ctx *gin.Context){
	ret := NewPublciError()
	type Param struct {
		Uid       int32 `form:"uid" binding:"required"`
		Token_id  int32 `form:"token_id" binding:"required"`
		Address   string `form:"address" binding:"required"`
		Mark  	  string `form:"mark" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		ctx.JSON(http.StatusOK, ret.GetResult())
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallAddressSave(param.Uid,param.Token_id,param.Address,param.Mark)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rsp)
	return
}
func (this *WalletGroup)AddressDelete(ctx *gin.Context)  {
	ret := NewPublciError()
	type Param struct {
		Uid       int32 `form:"uid" binding:"required"`
		Id  int32 `form:"Id" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		ctx.JSON(http.StatusOK, ret.GetResult())
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallAddressDelete(param.Uid,param.Id)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rsp)
	return
}
