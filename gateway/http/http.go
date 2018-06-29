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

	new(controller.ContentManageGroup).Router(r)

=======
	new(controller.ContentManageGroup).Router(r)
>>>>>>> 24f28d93fbf1067e4b0682efae084fa17d0e5fee
	return r
}

func InitHttpServer() {
	port := cf.Cfg.MustInt("http", "port")
	r := initRouter()
	r.Run(fmt.Sprintf(":%d", port))
}
