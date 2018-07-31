package controller

import (
	//"digicon/gateway/log"
	"digicon/gateway/rpc"
	Err "digicon/proto/common"
	proto "digicon/proto/rpc"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContentManageGroup struct{}

func (c *ContentManageGroup) Router(r *gin.Engine) {
	cont := r.Group("/content")

	{
		cont.POST("/b", c.Banner)
		cont.POST("/addlink", c.AddFriendlyLink)
		cont.GET("/linklist", c.GetFriendlyLink)
		cont.GET("/bannerlist", c.GetBannerList)
	}
}

func (cm *ContentManageGroup) Banner(c *gin.Context) {
	ret := Err.NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	return
}

func (cm *ContentManageGroup) AddFriendlyLink(c *gin.Context) {
	ret := Err.NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		WebName   string `form:"web_name" json:"web_name" binding:"required"`
		LinkName  string `form:"link_name" json:"link_name" binding:"required"`
		Aorder    int32  `form:"order" json:"order" binding:"required"`
		LinkState int32  `form:"link_state" json:"link_state" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(Err.ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.PublicService.CallAddFriendlyLink(&proto.AddFriendlyLinkRequest{
		WebName:   req.WebName,
		LinkName:  req.LinkName,
		Aorder:    req.Aorder,
		LinkState: req.LinkState,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(Err.ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Code)
	return
}

func (cm *ContentManageGroup) GetFriendlyLink(c *gin.Context) {
	ret := Err.NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Page  int32 `form:"page" json:"page" binding:"required"`
		Count int32 `form:"count" json:"count" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(Err.ERRCODE_PARAM, err.Error())
		return
	}

	fmt.Println(req)
	rsp, err := rpc.InnerService.PublicService.CallGetFriendlyLink(&proto.FriendlyLinkRequest{})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(Err.ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Code)
	ret.SetDataSection("list", rsp.Friend)
	return
}

func (cm *ContentManageGroup) GetBannerList(c *gin.Context) {
	ret := Err.NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	fmt.Println("2121211111111111111111")
	rsp, err := rpc.InnerService.PublicService.CallGetBannerList(&proto.BannerRequest{})
	if err != nil {
		fmt.Println("cccccccccccccccccccccccccccccc")
		log.Errorf(err.Error())
		ret.SetErrCode(Err.ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Code)
	ret.SetDataSection("list", rsp.List)
	return
}
