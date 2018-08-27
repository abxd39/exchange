package model

import (
	"digicon/common/convert"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	"github.com/golang/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"
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
var CnyPriceMap map[int32]*proto.CnyBaseData
var SymbolPriceMap map[string]*proto.SymbolPriceData

func InitCnyPrice() {
	CnyPriceMap = make(map[int32]*proto.CnyBaseData, 0)
	SymbolPriceMap = make(map[string]*proto.SymbolPriceData, 0)

	r, err := DB.GetRedisConn().Get("history.price.go.micro").Result()
	if err != nil {
		log.Fatal("please init price first")
	}

	out := &proto.CnyPriceResponse{}
	err = jsonpb.UnmarshalString(r, out)
	if err != nil {
		log.Fatal("please init price first err %s", err.Error())
		return
	}

	for _, v := range out.Data {
		CnyPriceMap[v.TokenId] = v
	}

	for _, v := range out.Symbols {
		SymbolPriceMap[v.Symbol] = v
	}
}

func GetCnyPrice(token_id int32) int64 {
	g, ok := CnyPriceMap[token_id]
	if ok {
		return g.CnyPriceInt
	}
	return 0
}

func GetUsdPrice(token_id int32) int64 {
	g, ok := CnyPriceMap[token_id]
	if ok {
		return g.UsdPriceInt
	}
	return 0
}

func GetOnPriceCnyPrice(symbol string, on_price int64) string {
	g, ok := SymbolPriceMap[symbol]
	if !ok {
		return ""
	}
	return convert.Int64MulInt64DivInt64By8Bit(g.CnyPriceInt, on_price, g.Price)
}

func GetOnPriceCnyPrice2(symbol string, on_price string) string {
	g, ok := SymbolPriceMap[symbol]
	if !ok {
		return ""
	}
	return convert.Int64MulInt64DivInt64By8Bit(g.CnyPriceInt, convert.StringToInt64(on_price), g.Price)
}
