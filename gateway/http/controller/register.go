package controller

import (
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserControll struct{}

func (this *UserControll) Router(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/register", RegisterController)
		user.POST("/login", LoginController)
		user.POST("/forget", ForgetPwdController)
		user.POST("/auth", AuthSecurityController)
		user.POST("/change_pwd", ForgetPwdController)
	}
}

type RegisterParam struct {
	Phone      string `form:"phone" binding:"required"`
	Pwd        string `form:"pwd" binding:"required"`
	Confirm    string `form:"confirm" binding:"required"`
	InviteCode string `form:"invite_code" binding:"required"`
	Country    int    `form:"country" binding:"required"`
}

func RegisterController(c *gin.Context) {
	ret := NewErrorMessage()
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	var param RegisterParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret[ErrCodeRet] = ERRCODE_PARAM
		ret[ErrCodeMessage] = err.Error()
		return
	}

	if param.Pwd != param.Confirm {
		ret[ErrCodeRet] = ERRCODE_PWD_COMFIRM
		ret[ErrCodeMessage] = GetErrorMessage(ERRCODE_PWD_COMFIRM)
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallRegister(param.Phone, param.Pwd, param.InviteCode, param.Country)
	if err != nil {
		ret[ErrCodeRet] = ERRCODE_UNKNOWN
		ret[ErrCodeMessage] = err.Error()
		return
	}

	ret[ErrCodeRet] = rsp.Err
	ret[ErrCodeMessage] = GetErrorMessage(rsp.Err)

}

type LoginParam struct {
	Phone string `form:"phone" binding:"required"`
	Pwd   string `form:"pwd" binding:"required"`
}

func LoginController(c *gin.Context) {
	ret := NewErrorMessage()
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	var param LoginParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret[ErrCodeRet] = ERRCODE_PARAM
		ret[ErrCodeMessage] = err.Error()
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallLogin(param.Phone, param.Pwd)
	if err != nil {
		ret[ErrCodeRet] = ERRCODE_UNKNOWN
		ret[ErrCodeMessage] = err.Error()
		return
	}
	ret[ErrCodeRet] = rsp.Err
	ret[ErrCodeMessage] = GetErrorMessage(rsp.Err)
}

type ForgetPwdParam struct {
	Phone string `form:"phone" binding:"required"`
}

func ForgetPwdController(c *gin.Context) {
	ret := NewErrorMessage()
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()

	var param ForgetPwdParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret[ErrCodeRet] = ERRCODE_PARAM
		ret[ErrCodeMessage] = err.Error()
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallForgetPwd(param.Phone)
	if err != nil {
		ret[ErrCodeRet] = ERRCODE_UNKNOWN
		ret[ErrCodeMessage] = err.Error()
		return
	}

	ret[ErrCodeRet] = rsp.Err
	ret[ErrCodeMessage] = GetErrorMessage(rsp.Err)
	d := ret[RetData].(map[string]interface{})
	d["phone"] = rsp.Phone
	d["email"] = rsp.Email
}

type AuthSecurityParam struct {
	Phone     string `form:"phone" binding:"required"`
	PhoneCode string `form:"phone_code" binding:"required"`
	EmailCode string `form:"email_code" binding:"required"`
}

func AuthSecurityController(c *gin.Context) {
	ret := NewErrorMessage()
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	var param AuthSecurityParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret[ErrCodeRet] = ERRCODE_PARAM
		ret[ErrCodeMessage] = err.Error()
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallAuthSecurity(param.Phone, param.PhoneCode, param.EmailCode)
	if err != nil {
		ret[ErrCodeRet] = ERRCODE_UNKNOWN
		ret[ErrCodeMessage] = err.Error()
		return
	}
	ret[ErrCodeRet] = rsp.Err
	ret[ErrCodeMessage] = GetErrorMessage(rsp.Err)
	d := ret[RetData].(map[string]interface{})
	d["security_key"] = rsp.SecurityKey
}

type PhoneParam struct {
	Phone string `form:"phone" binding:"required"`
}

func SendPhoneSMSController(c *gin.Context) {
	ret := NewErrorMessage()
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	var param PhoneParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret[ErrCodeRet] = ERRCODE_PARAM
		ret[ErrCodeMessage] = err.Error()
		return
	}

}
