package http

import  (
	"github.com/gin-gonic/gin"
	cf "digicon/ws_service/conf"
	"digicon/ws_service/http/controller"
	"fmt"
)

func initRouter() *gin.Engine{
	r := gin.Default()
	new(controller.WebChatGroup).Router(r)
	return r
}


func InitHttpServer(){
	port := cf.Cfg.MustInt("http", "port")
	r := initRouter()
	r.Run(fmt.Sprintf(":%d", port ))
}