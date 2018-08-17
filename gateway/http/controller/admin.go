package controller

import (
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"fmt"
)

type AdminGroup struct{}

const KEY = "hhhhhhhhhhhhhhhhhh"

func (s *AdminGroup) Router(r *gin.Engine) {
	action := r.Group("/admin")
	{
		action.POST("/refresh", s.Refresh)

		action.POST("/register_reward", s.RegisterReward)

		action.POST("/users_total", s.UserToatl)

		action.POST("/get_users_balances", s.GetUsersBalances)

	}
}

func (s *AdminGroup) Refresh(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	param := struct {
		Uid uint64 `form:"uid" json:"uid" binding:"required"`
		Key string `form:"key" json:"key" binding:"required"`
	}{}

	if err := c.ShouldBind(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Key != KEY {
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

func (s *AdminGroup) RegisterReward(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	param := &struct {
		Uid int64  `form:"uid" binding:"required"`
		Key string `form:"key" json:"key" binding:"required"`
	}{}

	if err := c.ShouldBind(param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Key != KEY {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallRegisterReward(&proto.RegisterRewardRequest{
		Uid: param.Uid,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}


func (s *AdminGroup) UserToatl(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	param:= &struct {
		Uids []uint64 `json:"uid" binding:"required"`
		Key string `form:"key" json:"key" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Key != KEY {
		ret.SetErrCode(ERRCODE_PARAM)
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallTokenBalanceCny(&proto.TokenBalanceCnyRequest{
		Uids: param.Uids,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list",rsp.Data)
}


func (s *AdminGroup) GetUsersBalances (c *gin.Context) {
	ret := NewPublciError()
	defer func(){
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	req := &struct {
		Key       string  `  json:"key" binding:"required"`
		Uid       int64   ` json:"uid"     `   //  当前用户 uid
		Uids      []uint64 ` json:"uids"     binding:"required"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(err)
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	if req.Key != KEY {
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	rsp, err := rpc.InnerService.CurrencyService.CallGetUserBalanceUids(&proto.GetUserBalanceUids{
		Uids:   req.Uids,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}

	ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
    ret.SetDataSection("list", rsp.Data)
}

