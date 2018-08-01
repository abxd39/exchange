package service

import (
	"exc_order/utils"
	"fmt"
	"time"
)

//SDC交易
type Sdc struct {}

func NewSdc() *Sdc {
	return &Sdc{}
}

//委托下单
func (p *Sdc) Sell(local_symbol,symbol string,huobi_symbol string) {
	fmt.Println("sdc下单：",local_symbol,symbol)
	defer utils.PanicRecover()
	//local_symbol,symbol := utils.SymbolMap(symbol)
	if symbol == "" || local_symbol == "" {
		return
	}
	//替换成比特币
	r_local_symbol,r_symbol := utils.ReplaceSDCToBTC(local_symbol,symbol)
	fmt.Println("sdc替换：",r_local_symbol,r_symbol)
	if r_local_symbol == "" || r_symbol == "" {
		return
	}
	//symbol := "eosusdt"
	sell_data := NewGeneral().pullDataForHuoBi(huobi_symbol)//r_symbol
	if sell_data.Result == false {
		fmt.Println("sdc错误1")
		return
	}
	//按照比例格式化单价
	var res bool
	if huobi_symbol == "ethbtc" {
		sell_data.Price = utils.GetBTCToSDCPrice() * sell_data.Price
	} else if huobi_symbol == "btcusdt" {
		sell_data.Price = sell_data.Price / utils.GetBTCToSDCPrice()
		sell_data.Amount += 100
	}
	//sell_data.Price = utils.GetBTCToSDCPrice() * sell_data.Price //utils.FromatBTCToSDC(sell_data.Price)
	//if res != true {
	//	fmt.Println("sdc错误2")
	//	return
	//}
	//登录
	err,cfg := utils.GetGoConfigP()
	if err == false {
		fmt.Println("sdc错误3")
		return
	}
	user_data := NewLogin().Login(cfg.MustValue("sell_user","ukey"),cfg.MustValue("sell_user","pwd"),cfg.MustInt("sell_user","utype"))
	//开始委托下单
	opt := utils.RandOptPrice()
	res = NewGeneral().entrustOrder(user_data,sell_data,local_symbol,opt)
	if res == false {
		fmt.Println("sdc错误4",sell_data.Amount,sell_data.Price)
		return
	}
	go func() {
		for {
			select {
			case <- time.After(time.Duration(5) * time.Second):
				fmt.Println("延迟购买：")
				//购买
				p.Buy(sell_data,symbol,local_symbol,opt)
				return
			}
		}
	}()
	//购买
	//NewSdc().Buy(sell_data,symbol,local_symbol,opt)
}

//BTC/SDC委托下单处理
func (p *Sdc) BtcSdcSell(local_symbol,symbol string) {
	fmt.Println("sdc下单：",local_symbol,symbol)
	defer utils.PanicRecover()
	//local_symbol,symbol := utils.SymbolMap(symbol)
	if symbol == "" || local_symbol == "" {
		return
	}

	//拉取最新的比特币价格
	//bit_price := NewGeneral().pullPriceBySymbol("btcusdt")
	sell_data := utils.PullData{
		Result:true,
		Price:utils.GetBTCToSDCPrice(),
		Amount:utils.RandBTCToSDCNum()}

		fmt.Println("销售价格：",sell_data)

	//登录
	err,cfg := utils.GetGoConfigP()
	if err == false {
		return
	}
	user_data := NewLogin().Login(cfg.MustValue("sell_user","ukey"),cfg.MustValue("sell_user","pwd"),cfg.MustInt("sell_user","utype"))
	//开始委托下单
	opt := utils.RandOptPrice()
	res := NewGeneral().entrustOrder(user_data,sell_data,local_symbol,opt)
	if res == false {
		return
	}
	go func() {
		for {
			select {
			case <- time.After(time.Duration(5) * time.Second):
				fmt.Println("延迟购买：")
				//购买
				p.Buy(sell_data,symbol,local_symbol,opt)
				return
			}
		}
	}()
	//购买
	//NewSdc().Buy(sell_data,symbol,local_symbol,opt)
}

//购买
func (p *Sdc) Buy(sell_data utils.PullData,symbol,local_symbol string,opt int) {
	defer utils.PanicRecover()
	//local_symbol,symbol := utils.SymbolMap("ETH/USDT")
	//if symbol == "" {
	//	return
	//}
	//登录
	err,cfg := utils.GetGoConfigP()
	if err == false {
		return
	}
	user_data := NewLogin().Login(cfg.MustValue("buy_user","ukey"),cfg.MustValue("sell_user","pwd"),cfg.MustInt("sell_user","utype"))
	//NewLogin().BalanceList(user_data.Uid,user_data.Token)
	//购买
	NewGeneral().buyOrder(user_data,sell_data,local_symbol,opt)
}
