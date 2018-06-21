package models

import (
	"time"
)

type UserCurrencyAlipayPay struct {
	Uid        int       `xorm:"not null pk default 0 comment('用户uid') INT(10)"`
	Name       string    `xorm:"not null default '' comment('用户姓名') VARCHAR(20)"`
	Alipay     string    `xorm:"not null default '' comment('支付宝账号') VARCHAR(20)"`
	ReciptCode string    `xorm:"not null default '' comment('支付宝收款二维码图片路径') VARCHAR(100)"`
	CreateTime time.Time `xorm:"not null comment('创建时间') DATETIME"`
	UpdataTime time.Time `xorm:"not null comment('修改时间') DATETIME"`
}
