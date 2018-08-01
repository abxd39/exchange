package utils

import (
	"github.com/gin-gonic/gin"
	"unicode"
	"strings"
	"math/rand"
	"time"
)

var SymbolConfig = make(map[string]string)
func init() {
	//usdt
	SymbolConfig["VEN/USDT"] = "venusdt"
	SymbolConfig["XRP/USDT"] = "xrpusdt"
	SymbolConfig["NEO/USDT"] = "neousdt"
	SymbolConfig["SDC/USDT"] = "sdcusdt"
	SymbolConfig["LTC/USDT"] = "ltcusdt"
	SymbolConfig["TRX/USDT"] = "trxusdt"
	SymbolConfig["ONT/USDT"] = "ontusdt"
	SymbolConfig["BTC/USDT"] = "btcusdt"
	SymbolConfig["ETH/USDT"] = "ethusdt"
	SymbolConfig["EOS/USDT"] = "eosusdt"
	//
	////btc
	SymbolConfig["ETH/BTC"] = "ethbtc"
	//
	////eth
	SymbolConfig["EOS/ETH"] = "eoseth"
	//
	////sds
	SymbolConfig["ETH/SDC"] = "ethsdc"
	SymbolConfig["BTC/SDC"] = "btcsdc"
	//以sdc为主货币，卖BTC，买SDC
}

//公共返回错误结果
func GetCommonError(msg string) gin.H {
	data := gin.H{}
	data["ret"] = 0
	data["msg"] = msg
	return data
}

//判断是否包含中文
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

func SymbolMap(key string) (string,string) {
	if _,ok := SymbolConfig[key];!ok {
		return "",""
	}
	return key,SymbolConfig[key]
}

//判断key是否SDC/USDT
func IsSDC(key string) bool {
	if strings.Contains(key, "SDC") {
		return true
	}
	return false
}

//把带sdc的替换成btc
func ReplaceSDCToBTC(local_symbol,symbol string) (string,string) {
	return strings.Replace(local_symbol,"SDC","BTC",-1),strings.Replace(symbol,"sdc","btc",-1)
}

//格式化sdc对比特币单价
func FromatBTCToSDC(new_bit_price float64) (bool,float64) {
	res,cfg := GetGoConfigP()
	if res != true {
		return false,0.0
	}
	bit_price := cfg.MustFloat64("default","bit_price")
	sdc_price := cfg.MustFloat64("default","sdc_price")
	new_sdc_price := sdc_price * new_bit_price / bit_price
	return true,new_sdc_price
}

//随机市价或限价
func RandOptPrice() int {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(100)
	if num <50 {
		return 1
	}
	return 2
}

//获取BTC对SDC单价
func GetBTCToSDCPrice() float64 {
	res,cfg := GetGoConfigP()
	if res != true {
		return 0.0
	}
	bit_price := cfg.MustFloat64("default","bit_price")
	sdc_price := cfg.MustFloat64("default","sdc_price")
	return bit_price / sdc_price
}

//随机数
func RandBTCToSDCNum() float64 {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(100)
	return float64(num) / float64(100)
}


