package controller

import (
	"digicon/gateway/log"
	"digicon/gateway/rpc"
	x "digicon/proto/common"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticlesGroup struct{}

func (this *ArticlesGroup) Router(r *gin.Engine) {
	Articles := r.Group("/articles")
	{
		log.Log.Printf("Router func")
		Articles.GET("/des", this.ArticlesDetail)
		Articles.GET("/list", this.ArticlesList)

	}
}

func (this *ArticlesGroup) ArticlesDetail(c *gin.Context) {
	ret := x.NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	type ArticlesDetailParam struct {
		Id int32 `form:"id" binding:"required"`
	}
	var param ArticlesDetailParam
	if err := c.ShouldBindQuery(&param); err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(x.ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.PublicService.CallArticlesDesc(param.Id)

	if err != nil {
		ret.SetErrCode(x.ERRCODE_UNKNOWN, err.Error())
		return
	}
	type ArticlesCopy1 struct {
		Id            int    `json:"Id"`
		Title         string `json:"Title"`
		Description   string `json:"Description"`
		Content       string `json:"Content"`
		Covers        string `json:"Covers"`
		ContentImages string `json:"ContentImages"`
		Type          int    `json:"Type"`
		TypeName      string `json:"TypeName"`
		Author        string `json:"Author"`
		Weight        int    `json:"Weight"`
		Shares        int    `json:"Shares"`
		Hits          int    `json:"Hits"`
		Comments      int    `json:"Comments"`
		DisplayMark   int    `json:"DisplayMark"`
		AdminId       int    `json:"AdminId"`
		CreateTime    string `json:"CreateTime"`
		UpdateTime    string `json:"UpdateTime"`
		AdminNickname string `json:"AdminNickname"`
	}
	articles := &ArticlesCopy1{}
	if err = json.Unmarshal([]byte(rsp.Data), articles); err != nil {
		log.Log.Errorf(err.Error())
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("articles", &articles)
}

func (this *ArticlesGroup) ArticlesList(c *gin.Context) {

	ret := x.NewPublciError()
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
		log.Log.Errorf(err.Error())
		ret.SetErrCode(x.ERRCODE_PARAM, err.Error())
		return
	}
	//fmt.Println("param2:", param)
	rsp, err := rpc.InnerService.PublicService.CallArticlesList(param.ArticlesType, param.Page, param.PageNum)
	if err != nil {
		ret.SetErrCode(x.ERRCODE_UNKNOWN, err.Error())
		return
	}
	//fmt.Println("gatway return value ", rsp.Articles)
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list", rsp.Articles)
}
