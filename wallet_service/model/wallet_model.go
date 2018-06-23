package models

import (
	"time"
)

type Wallet struct {
	Id          int       `xorm:"INT(11)"`
	Uid         int       `xorm:"comment('用户id') INT(11)"`
	Tokenid     int       `xorm:"comment('币id') INT(11)"`
	TokenName   string    `xorm:"comment('币名称') VARCHAR(10)"`
	Address     string    `xorm:"comment('钱包地址') VARCHAR(42)"`
	Amount      string    `xorm:"comment('金额') DECIMAL(64,8)"`
	UpdatedTime time.Time `xorm:"not null default '0000-00-00 00:00:00' comment('更新时间') TIMESTAMP"`
	CreatedTime time.Time `xorm:"not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Freeze      string    `xorm:"comment('冻结金额') DECIMAL(64,8)"`
}
