package controller

import (
	"digicon/gateway/rpc"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelloController(c *gin.Context) {

	s, _ := c.Params.Get("name")

	rsp, err := rpc.InnerService.UserSevice.CallGreet(s)
	if err != nil {
		c.String(http.StatusOK, "err rsp")
		return
	}

	c.String(http.StatusOK, rsp.Greeting)
}
