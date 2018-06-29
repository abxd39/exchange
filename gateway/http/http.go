package http

import (
	cf "digicon/gateway/conf"
	"digicon/gateway/http/controller"
	"fmt"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	r := gin.Default()
	//gin.SetMode(gin.ReleaseMode)

	new(controller.UserGroup).Router(r)
	new(controller.WalletGroup).Router(r)
	new(controller.ArticleGroup).Router(r)
	new(controller.ActionGroup).Router(r)
	new(controller.CurrencyGroup).Router(r)
	new(controller.KineGroup).Router(r)
	new(controller.TokenGroup).Router(r)
<<<<<<< HEAD
	new(controller.MarketGroup).Router(r)
=======
	new(controller.ContentManageGroup).Router(r)
>>>>>>> b022b9ed479667f457bd75362cd5a5f6e7555df3
	return r
}

func InitHttpServer() {
	port := cf.Cfg.MustInt("http", "port")
	r := initRouter()
	r.Run(fmt.Sprintf(":%d", port))
}
