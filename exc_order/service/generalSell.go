package service

import (
	"exc_order/utils"
	"fmt"
	"github.com/tidwall/gjson"
	"time"
)

type General struct{}
func NewGeneral() *General {
	return &General{}
}

//出售
func (p *General) Sell(local_symbol,symbol string) {
	fmt.Println("一般下单：",local_symbol,symbol)
	defer utils.PanicRecover()
	//local_symbol,symbol := utils.SymbolMap(symbol)
	if symbol == "" {
		return
	}
	//symbol := "eosusdt"
	sell_data := p.pullDataForHuoBi(symbol)
	if sell_data.Result == false {
		fmt.Println("拉去数据失败")
		return
	}
	//登录
	err,cfg := utils.GetGoConfigP()
	if err == false {
		fmt.Println("登录失败")
		return
	}
	user_data := NewLogin().Login(cfg.MustValue("sell_user","ukey"),cfg.MustValue("sell_user","pwd"),cfg.MustInt("sell_user","utype"))
	//开始委托下单
	opt := utils.RandOptPrice()
	res := p.entrustOrder(user_data,sell_data,local_symbol,opt)
	if res == false {
		fmt.Println("委托下单失败")
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
}

//购买
func (p *General) Buy(sell_data utils.PullData,symbol,local_symbol string,opt int) {
	defer utils.PanicRecover()
	//local_symbol,symbol := utils.SymbolMap("ETH/USDT")
	if symbol == "" {
		return
	}
	//登录
	err,cfg := utils.GetGoConfigP()
	if err == false {
		fmt.Println("购买失败")
		return
	}
	user_data := NewLogin().Login(cfg.MustValue("buy_user","ukey"),cfg.MustValue("sell_user","pwd"),cfg.MustInt("sell_user","utype"))
	//NewLogin().BalanceList(user_data.Uid,user_data.Token)
	//购买
	p.buyOrder(user_data,sell_data,local_symbol,opt)
}

//组合购买
func (p *General) buyOrder(user_data utils.LoginData,trade_data utils.PullData,symbol string,opt int) {
	defer utils.PanicRecover()
	params := `uid=%d&token=%s&symbol=%s&opt=%d&on_price=%f&type=%d&num=%f`
	params = fmt.Sprintf(
		params,
		user_data.Uid,
		user_data.Token,
		symbol,
		1,
		trade_data.Price,
		opt,
		trade_data.Price * trade_data.Amount)
	url := utils.GetApiUrl("entrust_order")
	result := utils.HttpPostRequest(url,params)
	if gjson.Get(result,"code").Int() != 0 {
		fmt.Println("buy失败",result)
	}
	fmt.Println("buy参数：",params,result)
}


//拉取数据
func (p *General) pullDataForHuoBi(symbol string) utils.PullData {
	defer utils.PanicRecover()
	url := utils.GetHuoBiTradeApiUrl("trade",symbol)
	data := utils.HttpGetRequest(url)
	fmt.Println(data)
	if gjson.Get(data,"status").String() == "error" {
		return utils.PullData{Result:false}
	}
	list_data := gjson.Get(data,"tick.data").Array()
	for _,v := range list_data {
		return utils.PullData{
			Result:true,
			Price:v.Get("price").Float(),
			Amount:v.Get("amount").Float()}
	}
	return utils.PullData{Result:false}
}

//拉取指定symbol价格
func (p *General) pullPriceBySymbol(symbol string) float64 {
	defer utils.PanicRecover()
	url := utils.GetHuoBiTradeApiUrl("trade",symbol)
	data := utils.HttpGetRequest(url)
	//fmt.Println(data)
	if gjson.Get(data,"status").String() == "error" {
		return 0.0
	}
	return gjson.Get(data,"tick.close").Float()
}

//组合委托下单
func (p *General) entrustOrder(user_data utils.LoginData,trade_data utils.PullData,symbol string,opt int) bool {
	defer utils.PanicRecover()
	params := `uid=%d&token=%s&symbol=%s&opt=%d&on_price=%f&type=%d&num=%f`
	params = fmt.Sprintf(
		params,
		user_data.Uid,
		user_data.Token,
		symbol,
		2,
		trade_data.Price,
		opt,
		trade_data.Price * trade_data.Amount)
	url := utils.GetApiUrl("entrust_order")
	result := utils.HttpPostRequest(url,params)
	if gjson.Get(result,"code").Int() != 0 {
		fmt.Println("委托下单数据：",result,params)
		return false
	}
	fmt.Println("sell参数：",params)
	return true
}
