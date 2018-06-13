package controller

import (
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NoticeControll struct{}

func (this *NoticeControll) Router(r *gin.Engine) {
	notice := r.Group("/notice")
	{
		notice.GET("/des/:id", Noticedescription)
		notice.GET("/list", NoticeList)

	}
}

type NoticeListParam struct {
	ID int32 `form:"id" binding:"required"`
}

func Noticedescription(c *gin.Context) {
	ret := NewErrorMessage()
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	var param NoticeListParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret[ErrCodeRet] = ERRCODE_PARAM
		ret[ErrCodeMessage] = err.Error()
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallNoticeDesc(param.ID)
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

func NoticeList(c *gin.Context) {
	ret := NewErrorMessage()
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	var param NoticeListParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret[ErrCodeRet] = ERRCODE_PARAM
		ret[ErrCodeMessage] = err.Error()
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallnoticeList()
	if err != nil {
		ret[ErrCodeRet] = ERRCODE_UNKNOWN
		ret[ErrCodeMessage] = err.Error()
		return
	}

	ret[ErrCodeRet] = rsp.Err
	ret[ErrCodeMessage] = GetErrorMessage(rsp.Err)

	d := ret[RetData].(map[string]interface{})
	d["Notice"] = rsp.Notice
}
