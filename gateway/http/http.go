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
/*
	r.GET("/test/:name", controller.HelloController)
	r.POST("/user/register", controller.RegisterController)
	r.POST("/user/login", controller.LoginController)
	r.POST("/user/forget",controller.ForgetPwdController)
	r.POST("/user/auth",controller.AuthSecurityController)
	r.POST("/user/change_pwd",controller.ForgetPwdController)
*/
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
