package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

//捕获异常
func PanicRecover1(c *gin.Context) {
	if err := recover(); err != nil {
		c.JSON(http.StatusOK,gin.H{
			"ret":0,
			"msg":err})
	}
}

func PanicRecover() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}