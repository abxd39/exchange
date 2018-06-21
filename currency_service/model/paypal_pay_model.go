package model

import (
	proto "digicon/proto/rpc"
	"time"
)

type UserCurrencyPaypalPay struct {
	Uid        int       `xorm:"not null pk default 0 comment('用户uid') INT(10)"`
	Paypal     string    `xorm:"not null default '' comment('paypal 账号') VARCHAR(20)"`
	CreateTime time.Time `xorm:"not null comment('创建时间') DATETIME"`
	UpdateTime time.Time `xorm:"not null comment('修改时间') DATETIME"`
}

func (pal *UserCurrencyPaypalPay) SetPaypal(req *proto.PaypalRequest) (int32, error) {
	return 0, nil
}
