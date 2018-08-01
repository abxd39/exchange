package main

import (
	"exc_order/service"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	//
	//router.GET("/", func(c *gin.Context) {
	//	c.String(http.StatusOK, "Hello World")
	//})
	//router.GET("/wallet/create",service.CreateWallet)
	//router.GET("/wallet/list_wallets",service.ListWallets)
	//router.GET("/wallet/lock_wallet",service.LockWallet)
	//router.GET("/chain/get_info",service.GetInfo)
	//router.GET("/chain/get_block",service.GetBlock)
	//router.GET("/chain/get_account",service.GetAccount)
	//
	//router.GET("/wallet/create_key",service.CreateKey)
	//router.GET("/wallet/unlock_wallet",service.UnlockWallet)
	//router.GET("/wallet/list_keys",service.ListKeys)
	//router.GET("/chain/create_account",service.CreateAccount)
	//router.GET("/chain/transfer",service.TransferDeal)

	go service.Start()

	router.GET("/user/balance",service.NewLogin().BalanceList)
	//
	router.Run(":8000")

	//user_p := service.NewLogin()
	//user_p.Login("18665967060","byuwang@2016",1)
	select{}
}