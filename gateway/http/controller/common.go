package controller

import (
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonGroup struct{}

func (s *CommonGroup) Router(r *gin.Engine) {
	user := r.Group("/common")
	{
		user.GET("/token_list", s.GetToknesList)
		user.GET("/get_site_config", s.GetSiteConfig)
	}
}

func (s *CommonGroup) GetToknesList(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	rsp, err := rpc.InnerService.UserSevice.CallTokensList()
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err)
	ret.SetDataSection("list", rsp.Data)
}

func (s *CommonGroup) GetSiteConfig(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	rsp, err := rpc.InnerService.PublicService.CallGetSiteConfig()
	if err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetErrCode(rsp.Code, rsp.Message)

	if rsp.Code == ERRCODE_SUCCESS {
		site := struct {
			Name            string `json:"name"`
			EnglishName     string `json:"english_name"`
			Title           string `json:"title"`
			EnglishTitle    string `json:"english_title"`
			Logo            string `json:"logo"`
			Keyword         string `json:"keyword"`
			EnglishKeyword  string `json:"english_keyword"`
			Desc            string `json:"desc"`
			EnglishDesc     string `json:"english_desc"`
			Beian           string `json:"beian"`
			StatisticScript string `json:"statistic_script"`
		}{
			Name:            rsp.Data.Site.Name,
			EnglishName:     rsp.Data.Site.EnglishName,
			Title:           rsp.Data.Site.Title,
			EnglishTitle:    rsp.Data.Site.EnglishTitle,
			Logo:            rsp.Data.Site.Logo,
			Keyword:         rsp.Data.Site.Keyword,
			EnglishKeyword:  rsp.Data.Site.EnglishKeyword,
			Desc:            rsp.Data.Site.Desc,
			EnglishDesc:     rsp.Data.Site.EnglishDesc,
			Beian:           rsp.Data.Site.Beian,
			StatisticScript: rsp.Data.Site.StatisticScript,
		}

		kefu := struct {
			Phone   string `json:"phone"`
			Email   string `json:"email"`
			Address string `json:"address"`
			Dianbao string `json:"dianbao"`
		}{
			Phone:   rsp.Data.Kefu.Phone,
			Email:   rsp.Data.Kefu.Email,
			Address: rsp.Data.Kefu.Address,
			Dianbao: rsp.Data.Kefu.Dianbao,
		}

		ret.SetDataSection("site_config", site)
		ret.SetDataSection("kefu_config", kefu)
	}
}
