package http

import (
	cf "digicon/gateway/conf"
	"digicon/gateway/http/controller"
	"fmt"
	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	new(controller.UserControll).Router(r)
	/*
		user := r.Group("/user")
		{
			user.POST("/register", controller.RegisterController)
			user.POST("/login", controller.LoginController)
			user.POST("/forget",controller.ForgetPwdController)
			user.POST("/auth",controller.AuthSecurityController)
			user.POST("/change_pwd",controller.ForgetPwdController)
		}
	*/
	return r
}

func InitHttpServer() {
	port := cf.Cfg.MustInt("http", "port")
	r := initRouter()
	r.Run(fmt.Sprintf(":%d", port))
}
