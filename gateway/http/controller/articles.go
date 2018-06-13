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
		fmt.Println("ccc")
		Articles.GET("/des/:id", this.ArticlesDetail)
		Articles.GET("/list", this.ArticlesList)

	}
}

type ArticlesDetailParam struct {
	ID int32 `form:"id" binding:"required"`
}

func (this *ArticlesGroup) ArticlesDetail(c *gin.Context) {
	fmt.Println("bbb")
	ret := NewErrorMessage()
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	var param ArticlesDetailParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret[ErrCodeRet] = ERRCODE_PARAM
		ret[ErrCodeMessage] = err.Error()
		return
	}

	rsp, err := rpc.InnerService.PublicService.CallArticlesDesc(param.ID)
	fmt.Println("aaa")
	if err != nil {
		ret[ErrCodeRet] = ERRCODE_UNKNOWN
		ret[ErrCodeMessage] = err.Error()
		return
	}

	ret[ErrCodeRet] = rsp.Err
	ret[ErrCodeMessage] = GetErrorMessage(rsp.Err)
	d := ret[RetData].(map[string]interface{})
	d["id"] = rsp.Id
	d["Title"] = rsp.Title
	d["Description"] = rsp.Description
	d["Content"] = rsp.Content
	d["Covers"] = rsp.Covers
	d["ContentImages"] = rsp.ContentImages
	d["Type"] = rsp.Type
	d["TypeName"] = rsp.TypeName
	d["Author"] = rsp.Author
	d["Weight"] = rsp.Weight
	d["Shares"] = rsp.Shares
	d["Hits"] = rsp.Hits
	d["Comments"] = rsp.Comments
	d["DisplayMark"] = rsp.DisplayMark
	d["CreateTime"] = rsp.CreateTime
	d["UpdateTime"] = rsp.UpdateTime
	d["AdminId"] = rsp.AdminId
	d["AdminNickname"] = rsp.AdminNickname
}

type ArticlesListParam struct {
	ArticlesType int32 `form:"ArticlesType" binding:"required"`
	StartRow     int32 `form:"StartRow" binding:"required"`
	EndRow       int32 `form:"EndRow" binding:"required"`
}

func (this *ArticlesGroup) ArticlesList(c *gin.Context) {
	ret := NewErrorMessage()
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	var param ArticlesListParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret[ErrCodeRet] = ERRCODE_PARAM
		ret[ErrCodeMessage] = err.Error()
		return
	}

	rsp, err := rpc.InnerService.PublicService.CallArticlesList()
	if err != nil {
		ret[ErrCodeRet] = ERRCODE_UNKNOWN
		ret[ErrCodeMessage] = err.Error()
		return
	}

	ret[ErrCodeRet] = rsp.Err
	ret[ErrCodeMessage] = GetErrorMessage(rsp.Err)

	d := ret[RetData].(map[string]interface{})
	d["Articles"] = rsp.Articles
}
