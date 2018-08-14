package model

import (
	"digicon/common/errors"
	"digicon/public_service/dao"
)

const (
	SITE_CONFIG_NAME_SITE = "site"
	SITE_CONFIG_NAME_KEFU = "kefu"
	SITE_CONFIG_NAME_SMS  = "sms"
)

type ConfigModel struct {
	Name  string `xorm:"name"`
	Value string `xorm:"value"`
}

func (ConfigModel) TableName() string {
	return "config"
}

// 网站基础配置
type SiteConfig struct {
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
}

// 短信配置
type SmsConfig struct {
	Url      string `json:"url"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// 客服配置
type KefuConfig struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// 获取配置
func (c *ConfigModel) Get(name string) (*ConfigModel, error) {
	configModel := new(ConfigModel)

	session := dao.DB.GetMysqlConn().Where("1=1")
	has, err := session.And("name=?", name).Get(configModel)
	if err != nil {
		return nil, errors.NewSys(err)
	}
	if !has {
		return nil, errors.NewNormal("配置不存在")
	}

	return configModel, nil
}
