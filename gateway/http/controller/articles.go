package controller

import (
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticlesGroup struct{}

func (this *ArticlesGroup) Router(r *gin.Engine) {
	Articles := r.Group("/articles")
	{

		Articles.GET("/des", this.ArticlesDetail)
		Articles.GET("/list", this.ArticlesList)

	}
}

func (this *ArticlesGroup) ArticlesDetail(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	type ArticlesDetailParam struct {
		Id int32 `form:"id" binding:"required"`
	}
	var param ArticlesDetailParam
	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.PublicService.CallArticlesDesc(param.Id)

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("articles", rsp.Data)
	fmt.Println("gatway ", rsp)
}

func (this *ArticlesGroup) ArticlesList(c *gin.Context) {

	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	type ArticlesListParam struct {
		ArticlesType int32 `form:"type" binding:"required"`
		Page         int32 `form:"page" binding:"required"`
		PageNum      int32 `form:"page_num" binding:""`
	}
	var param ArticlesListParam
	//fmt.Println("param1:", param)
	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	//fmt.Println("param2:", param)
	rsp, err := rpc.InnerService.PublicService.CallArticlesList(param.ArticlesType, param.Page, param.PageNum)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	//fmt.Println("gatway return value ", rsp.Articles)
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list", rsp.Articles)
}
