package controller

import (
	"github.com/gin-gonic/gin"
)

type KineGroup struct{}

func (this *KineGroup) Router(r *gin.Engine) {
	Kine := r.Group("/kine")
	{

		//Currency.POST("/add_ads", this.AddAds)         // 新增广告(买卖)
		//Currency.POST("/updated_ads", this.UpdatedAds) // 修改广告(买卖)

		Kine.POST("/hello", this.Hline) // 欢迎K线条
		Kine.POST("/", this.Hline)      // 欢迎K线条

	}
}

func (this *KineGroup) Hline(c *gin.Context) {
	/*
		ret := NewErrorMessage()
		// 调用 rpc K线
		code, err := rpc.InnerService.KineService.CallHline(&proto.KineRequest{
			value: "歡迎",
		})

		if err != nil || code != 0 {
			ret[ERR_CODE_RET] = ERRCODE_PARAM
			ret[ERR_CODE_MESSAGE] = GetErrorMessage(ERRCODE_PARAM)
			c.JSON(http.StatusOK, ret)
			return
		}

		ret[ERR_CODE_RET] = ERRCODE_SUCCESS
		ret[ERR_CODE_MESSAGE] = GetErrorMessage(ERRCODE_SUCCESS)
		c.JSON(http.StatusOK, ret)
	*/
}
