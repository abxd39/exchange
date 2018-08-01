package service

import (
	"fmt"
	"exc_order/utils"
	"github.com/tidwall/gjson"
	"strings"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Login struct{}

func NewLogin() *Login {
	return &Login{}
}

//regist
func (p *Login) regist() {

}

//login
func (p *Login) Login(ukey,pwd string,utype int) utils.LoginData {
	defer utils.PanicRecover()
	params := fmt.Sprintf("ukey=%s&pwd=%s&type=%d",ukey,pwd,utype)
	url := utils.GetApiUrl("login")
	data := utils.HttpPostRequest(url,params)
	if gjson.Get(url,"code").Int() != 0 {
		return utils.LoginData{Result:false}
	}
	return utils.LoginData{
		Result:true,
		Uid:gjson.Get(data,"data.data.uid").Int(),
		Token:gjson.Get(data,"data.data.token").String()}
}

//余额查询
func (p *Login) BalanceList(c *gin.Context) {
	err,cfg := utils.GetGoConfigP()
	if err == false {
		return
	}
	user_data := NewLogin().Login(cfg.MustValue("sell_user","ukey"),cfg.MustValue("sell_user","pwd"),cfg.MustInt("sell_user","utype"))
	defer utils.PanicRecover()
	url := utils.GetApiUrl("balancelist")
	params := "uid=%d&token=%s"
	params = fmt.Sprintf(params,user_data.Uid,user_data.Token)
	url = strings.Join([]string{url,"?",params},"")
	result := utils.HttpGetRequest(url)
	fmt.Println("列表：",result)
	c.JSON(http.StatusOK,gin.H{
		"data":result})
}
