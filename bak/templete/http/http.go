package http

import (
	cf "digicon/wallet_service/conf"
	"digicon/wallet_service/http/controller"
	"fmt"
	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/test/:name", controller.HelloController)
	return r
}

func InitHttpServer() {
	port := cf.Cfg.MustInt("http", "port")
	r := initRouter()
	r.Run(fmt.Sprintf(":%d", port))
}
