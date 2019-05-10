package http

import (
	cf "exchange/ws_service/conf"
	"exchange/ws_service/http/controller"
	"fmt"
	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	r := gin.Default()
	new(controller.WebChatGroup).Router(r)
	return r
}

func InitHttpServer() {
	port := cf.Cfg.MustInt("http", "port")
	r := initRouter()
	r.Run(fmt.Sprintf(":%d", port))
}
