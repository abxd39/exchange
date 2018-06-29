package model

import (
	"digicon/common/sms"
	"encoding/json"
	"github.com/liudng/godump"
	"strconv"
)

type KLineData struct {
	ID     int64   `json:"id"`     // K线ID
	Amount float64 `json:"amount"` // 成交量
	Count  int64   `json:"count"`  // 成交笔数
	Open   float64 `json:"open"`   // 开盘价
	Close  float64 `json:"close"`  // 收盘价, 当K线为最晚的一根时, 时最新成交价
	Low    float64 `json:"low"`    // 最低价
	High   float64 `json:"high"`   // 最高价
	Vol    float64 `json:"vol"`    // 成交额, 即SUM(每一笔成交价 * 该笔的成交数量)
}

type KLineReturn struct {
	Status  string      `json:"status"`   // 请求处理结果, "ok"、"error"
	Ts      int64       `json:"ts"`       // 响应生成时间点, 单位毫秒
	Data    []KLineData `json:"data"`     // KLine数据
	Ch      string      `json:"ch"`       // 数据所属的Channel, 格式: market.$symbol.kline.$period
	ErrCode string      `json:"err-code"` // 错误代码
	ErrMsg  string      `json:"err-msg"`  // 错误提示
}

// API请求地址, 不要带最后的/
const (
	MARKET_URL string = "https://api.huobi.pro"
	TRADE_URL  string = "https://api.huobi.pro"
)

func GetKLine(strSymbol, strPeriod string, nSize int) KLineReturn {
	kLineReturn := KLineReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["period"] = strPeriod
	mapParams["size"] = strconv.Itoa(nSize)

	strRequestUrl := "/market/history/kline"
	strUrl := MARKET_URL + strRequestUrl
	godump.Dump(mapParams)
	jsonKLineReturn := sms.HttpGetRequest(strUrl, mapParams)
	godump.Dump(jsonKLineReturn)
	json.Unmarshal([]byte(jsonKLineReturn), &kLineReturn)
	return kLineReturn
}
