package controller

import (
	"digicon/gateway/log"
	"digicon/gateway/rpc"
	x "digicon/proto/common"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleGroup struct{}

func (this *ArticleGroup) Router(r *gin.Engine) {
	article := r.Group("/article")
	{
		log.Log.Printf("Router func")
		article.GET("/des", this.Article)
		article.GET("/list", this.ArticleList)

	}
}

func (this *ArticleGroup) Article(c *gin.Context) {
	ret := x.NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	type ArticleParam struct {
		Id int32 `form:"id" binding:"required"`
	}
	var param ArticleParam
	if err := c.ShouldBindQuery(&param); err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(x.ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.PublicService.CallArticle(param.Id)

	if err != nil {
		ret.SetErrCode(x.ERRCODE_UNKNOWN, err.Error())
		return
	}
	type Article struct {
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
		Astatus       int    `json:"Astatus"`
		AdminId       int    `json:"AdminId"`
		CreateTime    string `json:"CreateTime"`
		UpdateTime    string `json:"UpdateTime"`
		AdminNickname string `json:"AdminNickname"`
	}
	arti := &Article{}
	if err = json.Unmarshal([]byte(rsp.Data), arti); err != nil {
		log.Log.Errorf(err.Error())
	}
	ret.SetErrCode(rsp.Err)
	ret.SetDataSection("article", &arti)
	return
}

func (this *ArticleGroup) ArticleList(c *gin.Context) {

	ret := x.NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	type ArticleListParam struct {
		ArticleType int32 `form:"type" binding:"required"`
		Page        int32 `form:"page" binding:"required"`
		PageNum     int32 `form:"page_num" binding:""`
	}
	var param ArticleListParam
	//fmt.Println("param1:", param)
	if err := c.ShouldBindQuery(&param); err != nil {
		log.Log.Errorf(err.Error())
		ret.SetErrCode(x.ERRCODE_PARAM, err.Error())
		return
	}
	//fmt.Println("param2:", param)
	rsp, err := rpc.InnerService.PublicService.CallArticleList(param.ArticleType, param.Page, param.PageNum)
	if err != nil {
		ret.SetErrCode(x.ERRCODE_UNKNOWN, err.Error())
		return
	}
	//fmt.Println("gatway return value ", rsp.Article)
	ret.SetErrCode(rsp.Err)
	ret.SetDataSection("list", rsp.Article)
	return
}
