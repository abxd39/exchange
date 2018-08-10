package model

import (
	. "digicon/price_service/dao"
	"log"
)

type ConfigTokenCny struct {
	TokenId  int   `xorm:"not null pk comment(' 币类型') INT(10)"`
	Price    int64 `xorm:"comment('人民币价格') BIGINT(20)"`
	UsdPrice int64 `xorm:"comment('美元价格') BIGINT(20)"`
}

var configTokenCnyData map[int32]*ConfigTokenCny

func InitConfigTokenCny() {
	configTokenCnyData = make(map[int32]*ConfigTokenCny, 0)
	err := DB.GetMysqlConn2().Find(&configTokenCnyData)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func GetTokenCnyPrice(token_id int32) *ConfigTokenCny {
	g, ok := configTokenCnyData[token_id]
	if ok {
		return g
	}
	log.Fatal("usdt cny price is null")
	return nil
}
