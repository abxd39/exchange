package model

type Trade struct {
	TradeId      int    `xorm:"not null pk comment('交易表的id') INT(11)"`
	TradeNo      string `xorm:"comment('订单号') VARCHAR(32)"`
	Uid          int    `xorm:"comment('买家uid') index INT(11)"`
	TokenId      int    `xorm:"comment('主货币id') index INT(11)"`
	TokenTradeId int    `xorm:"comment('交易币种') INT(11)"`
	Price        string `xorm:"comment('价格') DECIMAL(20,8)"`
	Num          string `xorm:"comment('数量') DECIMAL(20,8)"`
	Money        string `xorm:"DECIMAL(20,8)"`
	Fee          string `xorm:"comment('手续费') DECIMAL(20,8)"`
	Type         int    `xorm:"comment(' buy  0或sell 1') index TINYINT(255)"`
	DealTime     int    `xorm:"comment('成交时间') INT(11)"`
	Status       string `xorm:"comment('0是挂单，1是部分成交,2成交， -1撤销') VARCHAR(255)"`
}
