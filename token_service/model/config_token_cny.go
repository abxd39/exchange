package model

import (
	. "digicon/token_service/dao"
	. "digicon/token_service/log"
)

type ConfigTokenCny struct {
	TokenId int   `xorm:"not null pk comment(' 币类型') INT(10)"`
	Price   int64 `xorm:"comment('人民币价格') BIGINT(20)"`
}

var configTokenCnyData map[int]*ConfigTokenCny

func InitConfigTokenCny() {
	configTokenCnyData = make(map[int]*ConfigTokenCny, 0)
	err := DB.GetMysqlConn().Find(&configTokenCnyData)
	if err != nil {
		Log.Fatalln(err.Error())
	}
}

func GetTokenCnyPrice(token_id int) int64 {
	g, ok := configTokenCnyData[token_id]
	if ok {
		return g.Price
	}
	return 0
}
