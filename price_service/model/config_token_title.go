package model

import (
	. "digicon/price_service/dao"
	log "github.com/sirupsen/logrus"
)

type ConfigTokenTitle struct {
	TokenId int    `xorm:"not null pk comment('代币id') INT(11)"`
	Mark    string `xorm:"comment('代币名字') VARCHAR(32)"`
	Weight  int    `xorm:"comment('排序权重') INT(11)"`
}

var ConfigTitles []*ConfigTokenTitle

func InitConfigTitle() {
	ConfigTitles = make([]*ConfigTokenTitle, 0)
	err := DB.GetMysqlConn2().Asc("weight").Find(&ConfigTitles)
	if err != nil {
		log.Fatal(err.Error())
	}
}
