package service

import (
	"github.com/ouqiang/timewheel"
	"time"
	"exc_order/utils"
	"fmt"
)

var tw *timewheel.TimeWheel

//轮训指定的币
var loopConfig = make(map[string]int)
func init() {
	var data = utils.SymbolConfig
	for k,_ := range data {
		loopConfig[k] = 0
	}
}

//获取一个最小的
func getMinConfig() (string,string) {
	defer utils.PanicRecover()
	var minVal = 1
	var minK = ""
	var i = 0
	for k,v := range loopConfig {
		if i == 0 {
			minVal = v
			minK = k
		}
		if v < minVal {
			minVal = v
			minK = k
		}
		i++
	}
	loopConfig[minK] = minVal + 1
	return minK,getSymbolValueByKey(minK)
}

//根据key，查询value
func getSymbolValueByKey(key string) string {
	for k,v := range utils.SymbolConfig {
		if k == key {
			return v
		}
	}
	return ""
}

//开启定时器
func Start() {
	defer utils.PanicRecover()
	tw = timewheel.New(1 * time.Second, 3600, func(data timewheel.TaskData) {
		tw.AddTimer(3 * time.Second, "bitcoin", timewheel.TaskData{})
		local_symbol,symbol := getMinConfig()
		if local_symbol == "" || symbol == "" {
			return
		}
		if local_symbol == "BTC/SDC" {
			NewSdc().BtcSdcSell(local_symbol,symbol)
			return
		}
		if utils.IsSDC(local_symbol) == false {
			fmt.Println("一般定时器")
			NewGeneral().Sell(local_symbol,symbol)
			return
		}
		if symbol == "ethsdc" {
			NewSdc().Sell(local_symbol,symbol,"ethbtc")
		} else if symbol == "sdcusdt" {
			NewSdc().Sell(local_symbol,symbol,"btcusdt")
		}
		fmt.Println("SDC定时器")
		//NewSdc().Sell(local_symbol,symbol)
	})
	tw.Start()
	tw.AddTimer(1 * time.Second, "bitcoin", timewheel.TaskData{})
}