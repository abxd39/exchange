package controller

import (
	//"digicon/gateway/log"
	"digicon/gateway/rpc"
	Err "digicon/proto/common"
	proto "digicon/proto/rpc"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"fmt"
)

type ArticleGroup struct{}

func (this *ArticleGroup) Router(r *gin.Engine) {
	article := r.Group("/article")
	{
		log.Printf("Router func")
		article.GET("/des", this.Article)
		article.GET("/list", this.ArticleList)
		article.GET("/type_list", this.ArticleTypeList)

	}
}
func (this *ArticleGroup) ArticleTypeList(c *gin.Context) {
	ret := Err.NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	rsp, err := rpc.InnerService.PublicService.CallGetArticleTypeList(&proto.ArticleTypeRequest{})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(Err.ERRCODE_UNKNOWN, err.Error())
		return
	}
	//fmt.Println("gatway return value ", rsp.Article)
	ret.SetErrCode(rsp.Code)
	ret.SetDataSection("list", rsp.Type)
	return
}

func (this *ArticleGroup) Article(c *gin.Context) {
	ret := Err.NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	type ArticleParam struct {
		Id int32 `form:"id" binding:"required"`
	}
	var param ArticleParam
	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(Err.ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.PublicService.CallArticle(param.Id)

	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(Err.ERRCODE_UNKNOWN, err.Error())
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
		log.Errorf(err.Error())
		ret.SetErrCode(Err.ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Code)
	ret.SetDataSection("article", &arti)
	return
}

func (this *ArticleGroup) ArticleList(c *gin.Context) {
	//fmt.Println("噢噢噢噢噢噢噢噢噢噢噢噢噢噢噢噢哦哦哦")
	ret := Err.NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	type ArticleListParam struct {
		ArticleType int32 `form:"type" binding:"required"`
		Page        int32 `form:"page" binding:"required"`
		Rows     int32 `form:"rows" `
	}
	var param ArticleListParam
	//fmt.Println("param1:", param)
	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(Err.ERRCODE_PARAM, err.Error())
		return
	}
	fmt.Printf("%#v\n", param)
	rsp, err := rpc.InnerService.PublicService.CallArticleList(param.ArticleType, param.Page, param.Rows)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(Err.ERRCODE_UNKNOWN, err.Error())
		return
	}
	//fmt.Println("gatway return value ", rsp.Article)
	ret.SetErrCode(rsp.Code)
	ret.SetDataSection("list", rsp.Article)
	return
}
