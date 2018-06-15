package controller

import (
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ActionGroup struct{}

func (s *ActionGroup) Router(r *gin.Engine) {
	action := r.Group("/action")
	{
		action.POST("/get_google_key", s.GetGoogleAuthCode)
		action.POST("/auth_google_code", s.AuthGoogleCode)
	}
}

func (s *ActionGroup) GetGoogleAuthCode(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type GoogleAuthCodeParam struct {
		Uid int32 `form:"uid" binding:"required"`
	}
	var param GoogleAuthCodeParam

	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallGoogleSecretKey(param.Uid)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("key", rsp.SecretKey)
	ret.SetDataSection("url", rsp.Url)
}

func (s *ActionGroup) AuthGoogleCode(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type GoogleAuthCodeParam struct {
		Uid  int32  `form:"uid" binding:"required"`
		Code uint32 `form:"code" binding:"required"`
	}
	var param GoogleAuthCodeParam

	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallAuthGoogleSecretKey(param.Uid, param.Code)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}
