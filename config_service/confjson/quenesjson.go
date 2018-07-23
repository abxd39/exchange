package confjson

/*
表： config_quenes
 */

type ConfigQuenes struct {
	Id                   int64  `json:"id"`
	TokenId              int    `json:"token_id"`
	TokenTradeId         int    `json:"token_trade_id"`
	BitCount             int    `json:"bit_count"`
	MinOrderUnm          int64  `json:"min_order_unm"`
	Switch               int    `json:"switch"`
	Price                int64  `json:"price"`
	Name                 string `json:"name"`
	Scope                string `json:"scope"`
	Low                  int64  `json:"low"`
	High                 int64  `json:"high"`
	Amount               int64  `json:"amount"`
	SellPoundage         int64  `json:"sell_poundage"`
	BuyPoundage          int64  `json:"buy_poundage"`
	BuyMinimumPrice      int64  `json:"buy_minimum_price"`
	BuyMaxmunPrice       int64  `json:"buy_maxmun_price"`
	SellMinimumPrice     int64  `json:"sell_minimum_price"`
	SellMaxmumPrice      int64  `json:"sell_maxmum_price"`
	MinimumTradingVolume int64  `json:"minimum_trading_volume"`
	MaxmumTradingVolume  int64  `json:"maxmum_trading_volume"`
	BeginTime            int    `json:"begin_time"`
	EndTime              int    `json:"end_time"`
	SaturdaySwitch       int    `json:"saturday_switch"`
	SundaySwitch         int    `json:"sunday_switch"`
}




//
//type ConfigQuenes struct {
//	Id           int64  `json:"id"`
//	TokenId      int    `json:"token_id"`
//	TokenTradeId int    `json:"token_trade_id"`
//	Switch       int    `json:"switch"`
//	Name         string `json:"name"`
//	/*
//		Price        int64  `xorm:"comment('初始价格') BIGINT(20)"`
//		Scope        string `xorm:"comment('振幅') DECIMAL(6,2)"`
//		Low          int64  `xorm:"comment('最低价') BIGINT(20)"`
//		High         int64  `xorm:"comment('最高价') BIGINT(20)"`
//		Amount       int64  `xorm:"comment('成交量') BIGINT(20)"`
//	*/
//}


/*

*/

