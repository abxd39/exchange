package controller

import (
	"digicon/gateway/rpc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloController(c *gin.Context) {

	//s, _ := c.Params.Get("name")

	//rsp, err := rpc.InnerService.UserSevice.CallGreet(s)
	rsp, err := rpc.InnerService.UserSevice.CallNoticeDesc(1)
	if err != nil {
		c.String(http.StatusOK, "err rsp")
		return
	}

	//c.String(http.StatusOK, rsp)
	c.JSON(http.StatusOK, rsp)
}
