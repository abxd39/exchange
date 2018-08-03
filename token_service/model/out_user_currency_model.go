package model

type OutUserCurrency struct {
	Id        uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid       uint64 `xorm:"INT(10)"     json:"uid"`                                          // 用户ID
	TokenId   uint32 `xorm:"INT(10)"     json:"token_id"`                                     // 虚拟货币类型
	TokenName string `xorm:"VARCHAR(36)" json:"token_name"`                                   // 虚拟货币名字
	Freeze    int64  `xorm:"BIGINT not null default 0"   json:"freeze"`                       // 冻结
	Balance   int64  `xorm:"not null default 0 comment('余额') BIGINT"   json:"balance"`        // 余额
	Address   string `xorm:"not null default '' comment('充值地址') VARCHAR(255)" json:"address"` // 充值地址
	Version   int64  `xorm:"version"`
}

func (OutUserCurrency) TableName() string {
	return "user_currency"
}
