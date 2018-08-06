package controller

import (
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type AdminGroup struct{}


func (s *AdminGroup) Router(r *gin.Engine) {
	action := r.Group("/admin")
	{
		action.POST("/refresh", s.Refresh)

	}
}


func (s *AdminGroup) Refresh(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()


	param:=struct {
		Uid uint64 `form:"uid" json:"uid" binding:"required"`
		Key string  `form:"key" json:"key" binding:"required"`
	}{}

	if err := c.ShouldBind(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Key!="hhhhhhhhhhhhhhhhhh" {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}
	rsp, err := rpc.InnerService.UserSevice.CallRefresh(param.Uid)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}