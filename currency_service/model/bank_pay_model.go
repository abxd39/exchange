package models

import (
	"time"
)

type UserCurrencyBankPay struct {
	Uid        int       `xorm:"not null pk default 0 comment('用户uid') INT(10)"`
	Name       string    `xorm:"not null default '' comment('用户姓名') VARCHAR(20)"`
	CardNum    string    `xorm:"not null default '' comment('银行卡号') VARCHAR(20)"`
	BankName   string    `xorm:"not null default '' comment('银行名称') VARCHAR(20)"`
	BankInfo   string    `xorm:"not null default '' comment('支行名称') VARCHAR(20)"`
	CreateTime time.Time `xorm:"not null comment('创建时间') DATETIME"`
	UpdataTime time.Time `xorm:"not null comment('修改时间') DATETIME"`
}
