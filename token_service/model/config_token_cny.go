package model

import (
	proto "digicon/proto/rpc"
)

/*
type ConfigTokenCny struct {
	TokenId  int   `xorm:"not null pk comment(' 币类型') INT(10)"`
	Price    int64 `xorm:"comment('人民币价格') BIGINT(20)"`
	UsdPrice int64 `xorm:"comment('美元价格') BIGINT(20)"`
}

func (*ConfigTokenCny) TableName() string {
	return "config_token_cny"

}

var configTokenCnyData map[int]*ConfigTokenCny

func InitConfigTokenCny() {
	configTokenCnyData = make(map[int]*ConfigTokenCny, 0)
	err := DB.GetMysqlConn().Find(&configTokenCnyData)
	if err != nil {
		log.Fatalln(err.Error())
	}
	
	
}

func GetTokenCnyPrice(token_id int) int64 {
	g, ok := configTokenCnyData[token_id]
	if ok {
		return g.Price
	}
	return 0
}

func GetCnyData() map[int]*ConfigTokenCny {
	return configTokenCnyData
}

func GetTokenUsdPrice(token_id int) int64 {
	g, ok := configTokenCnyData[token_id]
	if ok {
		return g.UsdPrice
	}
	return 0
}
*/
var CnyPriceMap   map[int32]*proto.CnyBaseData

func InitCnyPrice()  {
	CnyPriceMap =make(map[int32]*proto.CnyBaseData,0)
}

func GetCnyPrice(token_id int32) int64 {
	g,ok:=CnyPriceMap[token_id]
	if ok {
		return g.CnyPriceInt
	}
	return 0
}

func GetUsdPrice(token_id int32) int64 {
	g,ok:=CnyPriceMap[token_id]
	if ok {
		return g.UsdPriceInt
	}
	return 0
}
