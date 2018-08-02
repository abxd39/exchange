//介绍用户基本操作的功能注册登录流程
package controller

import (
	"digicon/common/check"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"

	"bytes"
	"crypto/md5"
	"digicon/common/ip"
	"digicon/gateway/rpc/client"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
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
		//user.POST("/Modify_pwd", s.ModifyPwdcontroller)
		user.POST("/send_sms", s.SendPhoneSMSController)
		user.POST("/send_email", s.SendEmailController)
		user.POST("/modify_login_pwd", TokenVerify, s.ModifyLoginPwd)
		user.POST("/modify_phone", TokenVerify, s.ModifyPhone1)
		user.POST("/set_new_phone", TokenVerify, s.ModifyPhone2)
		user.POST("/modify_trade_pwd", TokenVerify, s.ResetTradePwd)
		user.GET("/get_auth_method", s.GetCheckAuthMethod)

		// bind user email
		user.POST("/bind_email", TokenVerify, s.BindUserEmail)
		user.POST("/bind_phone", TokenVerify, s.BindUserPhone)
		user.POST("/set_nickname", TokenVerify, s.SetNickName)
		//
		//user.POST("/unbind_email", s.UnBindUserEmail)
		//user.POST("/unbind_phone", s.UnBindUserPhone)

		user.GET("/api1", Geeteam)
		user.POST("/api2", Geeteam2)
		//一级实名认证
		user.POST("/first_verify",TokenVerify,s.FirstVerify)
		//二级认证
		user.POST("/second_verify",TokenVerify,s.SecondVerify)
		//获取认证的次数
		user.GET("/verify_count",s.GetVerifyCount)
	}
}

func (s *UserGroup)GetVerifyCount(c*gin.Context)  {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	p:= &struct {
		Uid uint64 `form:"uid" json:"uid" binding:"required"`
	}{}
	if err := c.ShouldBind(p); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp,err:=rpc.InnerService.UserSevice.CallGetVerifyCount(p.Uid)
	if err!=nil{
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		log.Errorf(err.Error())
		return
	}
	ret.SetErrCode(rsp.Code)
	ret.SetDataSection("first", rsp.FirstCount)
	ret.SetDataSection("second", rsp.SecondCount)
	return
}

func (s *UserGroup)SecondVerify(c*gin.Context)  {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	p:= &struct {
		Uid int32 `form:"uid" json:"uid" binding:"required"`
		Fp string `form:"front_path" json:"front_path" binding:"required"`
		Rp string `form:"reverse_path" json:"front_path" binding:"required"`
		Hp string `form:"head_path" json:"front_path" binding:"required"`
		Number string `form:"number" json:"number" binding:"required"`
		VPath string `form:"video_path" json:"video_path" binding:"required"`
	}{}
	if err := c.ShouldBind(p); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	fmt.Println("---------------------->123")
	rsp,err:=rpc.InnerService.UserSevice.CallSecondVerify(&proto.SecondRequest{
		Uid:p.Uid,
		FrontPath:p.Fp,
		ReversePath:p.Rp,
		HeadPath:p.Hp,
		Number:p.Number,
		VideoPath:p.VPath,
		})
	if err!=nil{
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		log.Errorf(err.Error())
		return
	}
	ret.SetErrCode(rsp.Code)
	return
}

func (s *UserGroup)FirstVerify(c*gin.Context)  {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	param := &struct {
		Uid int32 `form:"uid" json:"uid" binding:"required"`
		Code string `form:"code" json:"code" binding:"required"` // 电话验证码
		Id string `form:"id" json:"id" binding:"required"`//身份证号码
		RealName string `form:"name" json:"name" binding:"required"`
		Gcode uint32 `form:"gcode" json:"gcode"`
	}{}
	if err := c.ShouldBind(param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp,err:=rpc.InnerService.UserSevice.CallFirstVerify(&proto.FirstVerifyRequest{
		Uid:uint64(param.Uid),
		RealName:param.RealName,
		GoogleCode:param.Gcode,
		PhoneCode:param.Code,
		IdCode:param.Id,

	})
	if err!=nil{
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		log.Errorf(err.Error())
		return
	}
	ret.SetErrCode(rsp.Code)
	return
}

func Geeteam(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	param := struct {
		Phone string `form:"phone" json:"phone" binding:"required"` // 用户ID
	}{}

	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallApi1(param.Phone)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("data", rsp.Data)
}

func Geeteam2(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	/*
		param := struct {
			Phone   	string `form:"phone" json:"phone" binding:"required"`     // 用户ID
			Challenge   string `form:"geetest_challenge" json:"uid" binding:"required"`     // 用户ID
			Validate   string `form:"geetest_validate" json:"uid" binding:"required"`     // 用户ID
			SecCode   string `form:"geetest_seccode" json:"uid" binding:"required"`     // 用户ID
			Status   int32 `form:"gt_server_status" json:"uid" binding:"required"`     // 用户ID
		}{}
	*/

	param := struct {
		Phone     string `form:"phone" json:"phone" binding:"required"`   // 用户ID
		Challenge string `form:"challenge" json:"uid" binding:"required"` // 用户ID
		Validate  string `form:"validate" json:"uid" binding:"required"`  // 用户ID
		SecCode   string `form:"seccode" json:"uid" binding:"required"`   // 用户ID
		Status    int32  `form:"status" json:"uid" binding:"required"`    // 用户ID
	}{}
	if err := c.ShouldBind(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallApi2(&proto.Api2Request{
		Challenge: param.Challenge,
		Validate:  param.Validate,
		Seccode:   param.SecCode,
		Phone:     param.Phone,
		Status:    param.Status,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)

}

func TokenVerify(c *gin.Context) {

	ret := NewPublciError()

	param := struct {
		Token string `form:"token" json:"token" binding:"required"` // token验证
		Uid   uint64 `form:"uid" json:"uid" binding:"required"`     // 用户ID
	}{}

	if c.Request.Method == "POST" {
		if err := c.ShouldBind(&param); err != nil {
			log.Errorf(err.Error())
			ret.SetErrCode(ERRCODE_PARAM, err.Error())
			c.AbortWithStatusJSON(http.StatusOK, ret.GetResult())
			return
		}
	} else if c.Request.Method == "GET" {
		if err := c.ShouldBindQuery(&param); err != nil {
			log.Errorf(err.Error())
			ret.SetErrCode(ERRCODE_PARAM, err.Error())
			c.AbortWithStatusJSON(http.StatusOK, ret.GetResult())
			return
		}
	} else {
		if err := c.ShouldBind(&param); err != nil {
			log.Errorf(err.Error())
			ret.SetErrCode(ERRCODE_PARAM, err.Error())
			c.AbortWithStatusJSON(http.StatusOK, ret.GetResult())
			return
		}
	}

	rsp, err := rpc.InnerService.UserSevice.CallTokenVerify(param.Uid, []byte(param.Token))
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		c.AbortWithStatusJSON(http.StatusOK, ret.GetResult())
		return
	}
	if rsp.Err == ERRCODE_SUCCESS {
		c.Next()
	} else {
		ret.SetErrCode(rsp.Err, rsp.Message)
		c.AbortWithStatusJSON(http.StatusOK, ret.GetResult())

		return
	}

}

//用户注册
func (s *UserGroup) RegisterController(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type RegisterParam struct {
		Ukey       string `form:"ukey" binding:"required"`
		Pwd        string `form:"pwd" binding:"required"`
		Confirm    string `form:"confirm" binding:"required"`
		InviteCode string `form:"invite_code" `
		Country    string `form:"country" `
		Code       string `form:"code" binding:"required"`
		Type       int32  `form:"type" binding:"required"`
	}

	var param RegisterParam
	if err := c.ShouldBind(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Type == 1 {
		if param.Country == "" {

			ret.SetErrCode(ERRCODE_PARAM)
			return
		}
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
		log.Errorf(err.Error())
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
		log.Errorf(err.Error())
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
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallLogin(param.Ukey, param.Pwd, param.Type, ip.RemoteIp(c.Request))
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	if rsp.Err != ERRCODE_SUCCESS {
		return
	}
	ret.SetDataSection(RET_DATA, client.LoginUserBaseData{
		Uid:   rsp.Data.Uid,
		Token: string(rsp.Data.Token),
	})
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
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallForgetPwd(param.Ukey, param.Pwd, param.Code, param.Type)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}

//安全认证
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
		log.Errorf(err.Error())
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
		Region string `form:"region" binding:"required"`
		Phone  string `form:"phone" binding:"required"`
		Type   int32  `form:"type" binding:"required"`
	}

	var param PhoneParam
	if err := c.ShouldBind(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallSendSms(param.Phone, param.Region, param.Type)
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
		log.Errorf(err.Error())
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
		Type  int32  `form:"type" binding:"required"`
	}

	var param EamilParam
	if err := c.ShouldBind(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if ok := check.CheckEmail(param.Email); !ok {
		ret.SetErrCode(ERRCODE_SMS_PHONE_FORMAT)
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallSendEmail(param.Email, param.Type)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetErrCode(rsp.Err, rsp.Message)
}

func (s *UserGroup) ModifyLoginPwd(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Uid        uint64 `form:"uid" binding:"required"`
		Token      string `form:"token" binding:"required"`
		OldPwd     string `form:"old_pwd" binding:"required"`
		NewPwd     string `form:"new_pwd" binding:"required"`
		ConfirmPwd string `form:"confirm_pwd" binding:"required"`
		Verify     string `form:"verify" binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.UserSevice.CallModifyUserLoginPwd(&proto.UserModifyLoginPwdRequest{
		Uid:        req.Uid,
		Token:      req.Token,
		OldPwd:     req.OldPwd,
		NewPwd:     req.NewPwd,
		ConfirmPwd: req.ConfirmPwd,
		Verify:     req.Verify,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err)
	return
}

func (s *UserGroup) ModifyPhone1(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid    uint64 `form:"uid" binding:"required"`
		Token  string `form:"token" binding:"required"`
		Verify string `form:"verify" binding:"required"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	fmt.Println("000000000000000000000000000000000000")
	fmt.Println(req)
	rsp, err := rpc.InnerService.UserSevice.CallModifyPhone1(&proto.UserModifyPhoneRequest{
		Uid:    req.Uid,
		Token:  req.Token,
		Verify: req.Verify,
	})

	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err)
	return
}

func (s *UserGroup) ModifyPhone2(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		Country string `form:"country" binding:"required"`
		Phone   string `form:"phone" binding:"required"`
		Verify  string `form:"verify" binding:"required"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.UserSevice.CallModifyPhone2(&proto.UserSetNewPhoneRequest{
		Uid:     req.Uid,
		Token:   req.Token,
		Country: req.Country,
		Phone:   req.Phone,
		Verify:  req.Verify,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err)
}

func (s *UserGroup) ResetTradePwd(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Uid        uint64 `form:"uid" binding:"required"`
		Token      string `form:"token" binding:"required"`
		NewPwd     string `form:"new_pwd" binding:"required"`
		ConfirmPwd string `form:"confirm_pwd" binding:"required"`
		Verify     string `form:"verify" binding:"required"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.UserSevice.CallModifyTradePwd(&proto.UserModifyTradePwdRequest{
		Uid:        req.Uid,
		Token:      req.Token,
		NewPwd:     req.NewPwd,
		ConfirmPwd: req.ConfirmPwd,
		Verify:     req.Verify,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err)
}

func (s *UserGroup) SetNickName(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Uid      uint64 `form:"uid" binding:"required" json:"uid"`
		Token    string `form:"token" binding:"required" json:"token"`
		NickName string `form:"nick_name" json:"nick_name" binding:"required"`
		Url      string `form:"file" json:"file" binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	url, err := s.upload_picture(req.Url)
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UOPLOA_FAILED, err.Error())
		return
	}
	if url == `` {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UOPLOA_FAILED, err.Error())
		return
	}
	rsp, err := rpc.InnerService.UserSevice.CallSetNickName(&proto.UserSetNickNameRequest{
		Uid:           req.Uid,
		Token:         req.Token,
		NickName:      req.NickName,
		HeadSculpture: url,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err)
}

func (s *UserGroup) GetCheckAuthMethod(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := struct {
		Ukey string `form:"ukey" binding:"required"`
		Type int32  `form:"type" binding:"required"`
	}{}
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	rsp, err := rpc.InnerService.UserSevice.CallCheckAuthSecurity(&proto.CheckSecurityRequest{
		Ukey: req.Ukey,
		Type: req.Type,
	})
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err)
	ret.SetDataSection("auth", rsp.Auth)
	ret.SetDataSection("region", rsp.Region)
}

/*
	// bind user email
*/
func (s *UserGroup) BindUserEmail(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid        uint64 `form:"uid"          json:"uid"          binding:"required"`
		Email      string `form:"email"        json:"email"        binding:"required"`
		EmailCode  string `form:"email_code"   json:"email_code"   binding:"required"`
		VerifyCode string `form:"verify_code"  json:"verify_code"  binding:"required"`
		VerifyType uint64 `form:"verify_type"  json:"verify_type"  binding:"required"` // 验证类型 (1: 短信验证, 2 谷歌验证, )
	}{}

	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("bind error:", err.Error())
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	fmt.Println("req:", req)
	rsp, err := rpc.InnerService.UserSevice.CallBindEmail(&proto.BindEmailRequest{
		Uid:        req.Uid,
		Email:      req.Email,
		EmailCode:  req.EmailCode,
		VerifyCode: req.VerifyCode,
		VerifyType: req.VerifyType,
	})
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetDataSection("data", rsp.Data)
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	return

}

/*
	func: bind user phone
*/
func (s *UserGroup) BindUserPhone(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		Uid        uint64 `form:"uid"          json:"uid"          binding:"required"`
		Phone      string `form:"phone"        json:"phone"        binding:"required"`
		PhoneCode  string `form:"phone_code"   json:"phone_code"   binding:"required"`
		VerifyCode string `form:"verify_code"  json:"verify_code"  binding:"required"`
		VerifyType uint64 `form:"verify_type"  json:"verify_type"  binding:"required"` // 验证类型 ( 1邮箱验证, 2谷歌验证 )
	}{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.UserSevice.CallBindPhone(&proto.BindPhoneRequest{
		Uid:        req.Uid,
		Phone:      req.Phone,
		PhoneCode:  req.PhoneCode,
		VerifyCode: req.VerifyCode,
		VerifyType: req.VerifyType,
	})
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetDataSection("data", rsp.Data)
	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	return

}

/*
func: UnBindUserEmail
*/
func (s *UserGroup) UnBindUserEmail(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		//Uid :=
	}{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

/*
func: UnBindUserPhone
*/
func (s *UserGroup) UnBindUserPhone(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		//
	}{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

//上传Ali coud
func (a *UserGroup) upload_picture(file string) (string, error) {
	var remoteurl string = "https://sdun.oss-cn-shenzhen.aliyuncs.com/"
	client, err := oss.New("http://oss-cn-shenzhen.aliyuncs.com", "LTAIcJgRedhxruPq", "d7p6tWRfy0B2QaRXk7q4mb5seLROtb")
	if err != nil {
		// HandleError(err)
		panic(err)
	}
	bucket, err := client.Bucket("sdun")
	if err != nil {
		return "", err
	}
	//查找base64
	fmt.Println("base34-1")
	base := strings.Index(file, ";base64,")
	if base < 0 {
		fmt.Println("base34-3")
		// 是远程的oss 文件路径
		return file, nil
	}
	fmt.Println("base34-2")
	fmt.Println(file)

	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	subm := strings.IndexByte(file, ',')
	if subm < 0 {
		return "", errors.New("find failed!!")
	}
	substr := file[:subm]
	subb := strings.IndexByte(substr, '/')
	sube := strings.IndexByte(substr, ';')
	if subb < 0 || sube < 0 {
		return "", errors.New("find fail!!")
	}
	fmt.Println(subb, sube, subm)
	fSuffix := substr[subb+1 : sube]
	value := file[subm+1:]
	h := md5.New()
	tempValue := value
	tempValue += timestamp
	h.Write([]byte(tempValue)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	okey := hex.EncodeToString(cipherStr)
	//fmt.Println(okey)
	okey += "."
	okey += fSuffix
	//fmt.Printf("%#v\n", okey)
	ddd, _ := base64.StdEncoding.DecodeString(value)
	err = bucket.PutObject(okey, bytes.NewReader(ddd))
	if err != nil {
		//fmt.Println(filePath)
		return "", err
	}
	fmt.Println(remoteurl + okey)
	return remoteurl + okey, nil
}
