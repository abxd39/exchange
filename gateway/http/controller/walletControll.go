package controller

import (
	"digicon/gateway/rpc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WalletController struct {

}

func (this *WalletController) Router(router *gin.Engine){

	r := router.Group("wallet")
	r.GET("create",this.create)
	r.GET("update",this.update)
	r.GET("query",this.query)
	r.GET("delete",this.delete)
	r.GET("findOne",this.findOne)
}

func (this *WalletController)Index(ctx *gin.Context){

}
func (this *WalletController)create(ctx *gin.Context){
	rsp, err := rpc.InnerService.WalletSevice.CallCreateWallet(1,1)
	if err != nil {
		ctx.String(http.StatusOK, "err 0000 rsp")
		return
	}
	//var ret = NewErrorMessage()
	//ret["code"] = "0"
	//ret["msg"] = rsp.
	ctx.JSON(http.StatusOK, rsp)
}
func (this *WalletController)update(ctx *gin.Context){
	rsp, err := rpc.InnerService.WalletSevice.Callhello("eth")
	if err != nil {
		ctx.String(http.StatusOK, "err 0000 rsp")
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}
func (this *WalletController)query(ctx *gin.Context){

}
func (this *WalletController)delete(ctx *gin.Context){

}
func (this *WalletController)findOne(ctx *gin.Context){

}
func (this *WalletController)before() gin.HandlerFunc {
		return nil
}