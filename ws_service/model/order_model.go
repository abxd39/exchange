package model

import "database/sql"

type Order struct {
	Id          uint64         `xorm:"not null pk autoincr comment('ID')  INT(10)"                 json:"id"`
	OrderId     string         `xorm:"not null pk comment('订单ID') INT(10)"                        json:"order_id"`
	AdId        uint64         `xorm:"not null default 0 comment('广告ID') index INT(10)"           json:"ad_id"`
	AdType      uint32         `xorm:"not null default 0 comment('广告类型:1出售 2购买') TINYINT(1)"  json:"ad_type"`
	Price       int64          `xorm:"not null default 0 comment('价格') BIGINT(64)"                json:"price"`
	Num         int64          `xorm:"not null default 0 comment('数量') BIGINT(64)"                json:"num"`
	TokenId     uint64         `xorm:"not null default 0 comment('货币类型') INT(10)"               json:"token_id"`
	PayId       string         `xorm:"not null default 0 comment('支付类型') VARCHAR(64)"           json:"pay_id"`
	SellId      uint64         `xorm:"not null default 0 comment('卖家id') INT(10)"                json:"sell_id"`
	SellName    string         `xorm:"not null default '' comment('卖家昵称') VARCHAR(64)"          json:"sell_name"`
	BuyId       uint64         `xorm:"not null default 0 comment('买家id') INT(10)"                json:"buy_id"`
	BuyName     string         `xorm:"not null default '' comment('买家昵称') VARCHAR(64)"         json:"buy_name"`
	Fee         int64          `xorm:"not null default 0 comment('手续费用') BIGINT(64)"           json:"fee"`
	States      uint32         `xorm:"not null default 0 comment('订单状态: 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消') TINYINT(1)"   json:"states"`
	PayStatus   uint32         `xorm:"not null default 0 comment('支付状态: 1待支付 2待放行(已支付) 3确认支付(已完成)') TINYINT(1)"  json:"pay_status"`
	CancelType  uint32         `xorm:"not null default 0 comment('取消类型: 1卖方 2 买方') TINYINT(1)"   json:"cancel_type"`
	CreatedTime string         `xorm:"not null comment('创建时间') DATETIME"                       json:"created_time"`
	UpdatedTime string         `xorm:"comment('修改时间')     DATETIME"                           json:"updated_time"`
	ConfirmTime sql.NullString `xorm:"default null comment('确认支付时间')  DATETIME"             json:"confirm_time"`
	ReleaseTime sql.NullString `xorm:"default null comment('放行时间')     DATETIME"              json:"release_time"`
	ExpiryTime  string         `xorm:"comment('过期时间')     DATETIME"                           json:"expiry_time"`
}
