package model

type TokenDailySheet struct {
	Id           int64   `xorm:"not null pk autoincr comment('自增id')BIGINT(20)"`
	TokenId      int   `xorm:"not null comment('货币id') INT(11)"`
	FeeBuyCny    int64 `xorm:"not null comment('买手续费折合cny') BIGINT(20)"`
	FeeBuyTotal  int64 `xorm:"not null comment('买手续费总额') BIGINT(20)"`
	FeeSellCny   int64 `xorm:"not null comment('卖手续费折合cny') BIGINT(20)"`
	FeeSellTotal int64 `xorm:"not null comment('卖手续费总额') BIGINT(20)"`
	BuyTotal     int64 `xorm:"not null comment('买总额') BIGINT(20)"`
	BuyTotalCny  int64 `xorm:"not null comment('买总额折合') BIGINT(20)"`
	SellTotalCny int64 `xorm:"not null comment('卖总额折合') BIGINT(20)"`
	SellTotal    int64 `xorm:"not null comment('卖总额') BIGINT(20)"`
	Date         int64 `xorm:"not null comment('时间戳，精确到天') BIGINT(20)"`
	//Day          time.Time `xorm:"comment('那天') DATETIME"`
}
