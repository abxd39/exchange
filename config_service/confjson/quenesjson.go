package confjson

/*
表： config_quenes
*/
type ConfigQuenes struct {
	Id           int64  `json:"id"`
	TokenId      int    `json:"token_id"`
	TokenTradeId int    `json:"token_trade_id"`
	Switch       int    `json:"switch"`
	Name         string `json:"name"`
	/*
		Price        int64  `xorm:"comment('初始价格') BIGINT(20)"`
		Scope        string `xorm:"comment('振幅') DECIMAL(6,2)"`
		Low          int64  `xorm:"comment('最低价') BIGINT(20)"`
		High         int64  `xorm:"comment('最高价') BIGINT(20)"`
		Amount       int64  `xorm:"comment('成交量') BIGINT(20)"`
	*/
}

/*

 */
