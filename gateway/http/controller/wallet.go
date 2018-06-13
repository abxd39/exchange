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
	r.GET("create",this.Create)
	r.GET("update",this.Update)
	r.GET("query",this.Query)
	r.GET("delete",this.Delete)
	r.GET("findOne",this.FindOne)
}

func (this *WalletController)Index(ctx *gin.Context){

}
func (this *WalletController)Create(ctx *gin.Context){

	rsp, err := rpc.InnerService.WalletSevice.CallCreateWallet(1,1)
	if err != nil {
		ctx.String(http.StatusOK, "err 0000 rsp")
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
func (this *WalletController)Update(ctx *gin.Context){
	rsp, err := rpc.InnerService.WalletSevice.Callhello("eth")
	if err != nil {
		ctx.String(http.StatusOK, "err 0000 rsp")
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}
func (this *WalletController)Query(ctx *gin.Context){

}
func (this *WalletController)Delete(ctx *gin.Context){

}
func (this *WalletController)FindOne(ctx *gin.Context){

}
func (this *WalletController)before() gin.HandlerFunc {
		return nil
}