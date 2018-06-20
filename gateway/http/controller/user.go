//介绍用户基本操作的功能注册登录流程
package controller

import (
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	"net/http"

	"digicon/common/check"
	proto "digicon/proto/rpc"

	"github.com/gin-gonic/gin"
	"github.com/liudng/godump"
)

type UserGroup struct{}

func (s *UserGroup) Router(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/register", s.RegisterController)
		//user.POST("/register_phone", RegisterPhoneController)
		//user.GET("/register_email", RegisterEmailController)
		user.POST("/login", s.LoginController)
		user.POST("/forget", s.ForgetPwdController)
		user.POST("/auth", s.AuthSecurityController)
		user.POST("/change_pwd", s.ChangePwdcontroller)
		user.POST("/send_sms", s.SendPhoneSMSController)
		user.POST("/send_email", s.SendEmailController)
		user.POST("/change_login_pwd", s.ChangeLoginPwd)
		user.POST("/change_phone", s.ChangePhone1)
		user.POST("/set_new_phon", s.ChangePhone2)
		user.POST("change_trade_pwd", s.ResetTradePwd)
	}
}

func (s *UserGroup) RegisterController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type RegisterParam struct {
		Ukey       string `form:"ukey" binding:"required"`
		Pwd        string `form:"pwd" binding:"required"`
		Confirm    string `form:"confirm" binding:"required"`
		InviteCode string `form:"invite_code" binding:"required"`
		Country    int32  `form:"country" binding:"required"`
		Code       string `form:"code" binding:"required"`
		Type       int32  `form:"type" binding:"required"`
	}

	var param RegisterParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Pwd != param.Confirm {
		ret.SetErrCode(ERRCODE_PWD_COMFIRM)
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallRegister(param.Ukey, param.Pwd, param.InviteCode, param.Country, param.Code, param.Type)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	/*
		ty, ok := c.GetPostForm("type")
		if !ok {
			ret := NewErrorMessage()
			defer func() {
				c.JSON(http.StatusOK, ret)
			}()
			ret[ERR_CODE_RET] = ERRCODE_PARAM
			ret[ERR_CODE_MESSAGE] = GetErrorMessage(ERRCODE_PARAM)
			return
		}



		if ty == "1" { //1电话注册
			s.RegisterPhoneController(c)
		} else if ty == "2" { //邮箱注册
			s.RegisterEmailController(c)
		}
	*/
}

//用户注册by phone

/*
func (s *UserGroup) RegisterPhoneController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type RegisterByPhoneParam struct {
		Phone      string `form:"phone" binding:"required"`
		Pwd        string `form:"pwd" binding:"required"`
		Confirm    string `form:"confirm" binding:"required"`
		InviteCode string `form:"invite_code" binding:"required"`
		Country    int    `form:"country" binding:"required"`
		Code       string `form:"code" binding:"required"`
	}

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

//用户注册by email
func (s *UserGroup) RegisterEmailController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type RegisterByEmailParam struct {
		Email      string `form:"email" binding:"required"`
		Pwd        string `form:"pwd" binding:"required"`
		Confirm    string `form:"confirm" binding:"required"`
		InviteCode string `form:"invite_code" binding:"required"`
		Country    int    `form:"country" binding:"required"`
		Code       string `form:"code" binding:"required"`
	}

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
*/
//用户登陆
func (s *UserGroup) LoginController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type LoginParam struct {
		Ukey string `form:"ukey" binding:"required"`
		Pwd  string `form:"pwd" binding:"required"`
		Type int32  `form:"type" binding:"required"`
	}
	var param LoginParam

	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallLogin(param.Ukey, param.Pwd, param.Type)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("data", rsp.Data)
}

//忘记密码
func (s *UserGroup) ForgetPwdController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type ForgetPwdParam struct {
		Ukey string `form:"ukey" binding:"required"`
		Type int32  `form:"type" binding:"required"`
		Code string `form:"code" binding:"required"`
		Pwd  string `form:"pwd" binding:"required"`
	}

	var param ForgetPwdParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallForgetPwd(param.Ukey, param.Pwd, param.Code, param.Type)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	godump.Dump(rsp)
	ret.SetErrCode(rsp.Err, rsp.Message)
}

//提交手机验证
func (s *UserGroup) AuthSecurityController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type AuthSecurityParam struct {
		Phone     string `form:"phone" binding:"required"`
		PhoneCode string `form:"phone_code" binding:"required"`
		EmailCode string `form:"email_code" binding:"required"`
	}

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

//发送短信
func (s *UserGroup) SendPhoneSMSController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type PhoneParam struct {
		Phone string `form:"phone" binding:"required"`
		Type  int32  `form:"type" binding:"required"`
	}

	var param PhoneParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if ok := check.CheckPhone(param.Phone); !ok {
		ret.SetErrCode(ERRCODE_SMS_PHONE_FORMAT)
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallSendSms(param.Phone, param.Type)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}

//改密码
func (s *UserGroup) ChangePwdcontroller(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type ChangePwdParam struct {
		SecurityKey string `form:"security_key" binding:"required"`
		Phone       string `form:"phone" binding:"required"`
		Pwd         string `form:"pwd" binding:"required"`
	}

	var param ChangePwdParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

}

//发生邮箱验证
func (s *UserGroup) SendEmailController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type EamilParam struct {
		Email string `form:"email" binding:"required"`
	}

	var param EamilParam
	if err := c.ShouldBind(&param); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if ok := check.CheckEmail(param.Email); !ok {
		ret.SetErrCode(ERRCODE_SMS_PHONE_FORMAT)
		return
	}
}

func (s *UserGroup) ChangeLoginPwd(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Uid        int32  `form:"id" binding:"required"`
		Token      string `form:"token" binding:"required"`
		Phone      string `form:"phone" binding:"required"`
		OldPwd     string `form:"old_pwd" binding:"required"`
		NewPwd     string `form:"new_pwd" binding:"required"`
		ConfirmPwd string `form:"confirm" binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.UserSevice.CallChangeUserLoginPwd(&proto.UserChangeLoginPwdRequest{
		Uid:        req.Uid,
		Phone:      req.Phone,
		Token:      req.Token,
		OldPwd:     req.OldPwd,
		NewPwd:     req.NewPwd,
		ConfirmPwd: req.ConfirmPwd,
	})
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	ret.SetDataSection("data", rsp.Data)
	ret.SetErrCode(ERRCODE_SUCCESS, err.Error())
	return
}

func (s *UserGroup) ChangePhone1(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid    int32  `form:"id" binding:"required"`
		Token  string `form:"token" binding:"required"`
		Phone  string `form:"phone" binding:"required"`
		verify string `form:"verify" binding:"required"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.UserSevice.CallChangePhone1(&proto.UserChangePhoneRequest{
		Uid:    req.Uid,
		Token:  req.Token,
		Phone:  req.Phone,
		Verify: req.verify,
	})
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	ret.SetDataSection("data", rsp.Data)
	ret.SetErrCode(ERRCODE_SUCCESS, err.Error())
	return
}

func (s *UserGroup) ChangePhone2(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid     int32  `form:"id" binding:"required"`
		Token   string `form:"token" binding:"required"`
		Country string `form:"country" binding:"required"`
		Phone   string `form:"phone" binding:"required"`
		verify  string `form:"verify" binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.UserSevice.CallChangePhone2(&proto.UserSetNewPhoneRequest{
		Uid:     req.Uid,
		Token:   req.Token,
		Country: req.Country,
		Phone:   req.Phone,
		Verify:  req.verify,
	})
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	ret.SetDataSection("data", rsp.Data)
	ret.SetErrCode(ERRCODE_SUCCESS, err.Error())
}

func (s *UserGroup) ResetTradePwd(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid        int32  `form:"id" binding:"required"`
		Token      string `form:"token" binding:"required"`
		NewPwd     string `form:"new_pwd" binding:"required"`
		ConfirmPwd string `form:"confirm_pwd" binding:"required"`
		Verify     string `form:"verify" binding:"required"`
	}{}
	rsp, err := rpc.InnerService.UserSevice.CallChangeTradePwd(&proto.UserChangeTradePwdRequest{
		Uid:        req.Uid,
		Token:      req.Token,
		NewPwd:     req.NewPwd,
		ConfirmPwd: req.ConfirmPwd,
		Verify:     req.Verify,
	})
	if err != nil {
		Log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	ret.SetDataSection("data", rsp.Data)
	ret.SetErrCode(ERRCODE_SUCCESS, err.Error())
}
