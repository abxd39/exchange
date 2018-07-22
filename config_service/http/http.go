package http

import (
	cf "digicon/config_service/conf"
	"fmt"
	"github.com/gin-gonic/gin"
	"digicon/config_service/http/controller"
)

func initRouter() *gin.Engine {
	r := gin.Default()
	//new(controller.WebChatGroup).Router(r)
	new(controller.ConfigGroup).Router(r)

	return r
}

func InitHttpServer() {
	port := cf.Cfg.MustInt("http", "port")
	r := initRouter()
	r.Run(fmt.Sprintf(":%d", port))
}
