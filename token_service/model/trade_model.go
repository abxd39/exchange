package model

const (
	TRADE_STATES_PART=1
	TRADE_STATES_ALL=2
	TRADE_STATES_DEL=3
)
type Trade struct {
	TradeId      int    `xorm:"not null pk comment('交易表的id') INT(11)"`
	TradeNo      string `xorm:"comment('订单号') VARCHAR(32)"`
	Uid          int    `xorm:"comment('买家uid') index INT(11)"`
	TokenId      int    `xorm:"comment('主货币id') index INT(11)"`
	TokenTradeId int    `xorm:"comment('交易币种') INT(11)"`
	Price        int64  `xorm:"comment('价格') BIGINT(20)"`
	Num          int64  `xorm:"comment('数量') BIGINT(20)"`
	Money        int64  `xorm:"BIGINT(20)"`
	Fee          int64  `xorm:"comment('手续费') BIGINT(20)"`
	Type         int    `xorm:"comment(' buy  0或sell 1') index TINYINT(255)"`
	DealTime     int64    `xorm:"comment('成交时间') BIGINT(20)"`
	States       int 	`xorm:"comment('1是挂单，2是部分成交,3成交， 4撤销') TINYINT(4)"`
}
