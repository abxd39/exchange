package controller

import (
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserGroup struct{}

func (this *UserGroup) Router(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/register", RegisterController)
		//user.POST("/register_phone", RegisterPhoneController)
		//user.GET("/register_email", RegisterEmailController)
		user.POST("/login", LoginController)
		user.POST("/forget", ForgetPwdController)
		user.POST("/auth", AuthSecurityController)
		user.POST("/change_pwd", ChangePwdcontroller)
		user.POST("/send_sms", SendPhoneSMSController)
		user.POST("/send_email", SendEmailController)
	}
}

func RegisterController(c *gin.Context) {
	ty, ok := c.GetPostForm("type")
	if !ok {
		ret := NewErrorMessage()
		defer func() {
			c.JSON(http.StatusOK, ret)
		}()
		ret[ErrCodeRet] = ERRCODE_PARAM
		ret[ErrCodeMessage] = GetErrorMessage(ERRCODE_PARAM)
		return
	}
	if ty == "1" { //1电话注册
		RegisterPhoneController(c)
	} else if ty == "2" { //邮箱注册
		RegisterEmailController(c)
	}
}

type RegisterByPhoneParam struct {
	Phone      string `form:"phone" binding:"required"`
	Pwd        string `form:"pwd" binding:"required"`
	Confirm    string `form:"confirm" binding:"required"`
	InviteCode string `form:"invite_code" binding:"required"`
	Country    int    `form:"country" binding:"required"`
	Code       string `form:"code" binding:"required"`
}

//用户注册by phone
func RegisterPhoneController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var param RegisterByPhoneParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Pwd != param.Confirm {
		ret.SetErrCode(ERRCODE_PWD_COMFIRM)
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallRegisterByPhone(param.Phone, param.Pwd, param.InviteCode, param.Country, param.Code)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}

type RegisterByEmailParam struct {
	Email      string `form:"email" binding:"required"`
	Pwd        string `form:"pwd" binding:"required"`
	Confirm    string `form:"confirm" binding:"required"`
	InviteCode string `form:"invite_code" binding:"required"`
	Country    int    `form:"country" binding:"required"`
	Code       string `form:"code" binding:"required"`
}

//用户注册by email
func RegisterEmailController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var param RegisterByEmailParam
	if err := c.ShouldBindQuery(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Pwd != param.Confirm {
		ret.SetErrCode(ERRCODE_PWD_COMFIRM)
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallRegisterByEmail(param.Email, param.Pwd, param.InviteCode, param.Country)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}

type LoginParam struct {
	Phone string `form:"phone" binding:"required"`
	Pwd   string `form:"pwd" binding:"required"`
}

//用户登陆
func LoginController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var param LoginParam

	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallLogin(param.Phone, param.Pwd)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}

type ForgetPwdParam struct {
	Phone string `form:"phone" binding:"required"`
}

//忘记密码
func ForgetPwdController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()

	var param ForgetPwdParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallForgetPwd(param.Phone)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("phone", rsp.Phone)
	ret.SetDataSection("email", rsp.Email)
}

type AuthSecurityParam struct {
	Phone     string `form:"phone" binding:"required"`
	PhoneCode string `form:"phone_code" binding:"required"`
	EmailCode string `form:"email_code" binding:"required"`
}

//提交手机验证
func AuthSecurityController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var param AuthSecurityParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallAuthSecurity(param.Phone, param.PhoneCode, param.EmailCode)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}

type PhoneParam struct {
	Phone string `form:"phone" binding:"required"`
	Type  int32  `form:"type" binding:"required"`
}

//发生短信
func SendPhoneSMSController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var param PhoneParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallSendSms(param.Phone, param.Type)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}

type ChangePwdParam struct {
	SecurityKey string `form:"security_key" binding:"required"`
	Phone       string `form:"phone" binding:"required"`
	Pwd         string `form:"pwd" binding:"required"`
}

//改密码
func ChangePwdcontroller(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var param ChangePwdParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

}

type EamilParam struct {
	Phone string `form:"phone" binding:"required"`
}

//发生邮箱验证
func SendEmailController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	var param EamilParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
}
