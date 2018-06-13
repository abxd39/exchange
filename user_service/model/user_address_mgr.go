package model

type UserAddressMgr struct {
	Uid     int    `xorm:"comment('用户ID') INT(11)"`
	TokenId int    `xorm:"comment('货币类型') INT(11)"`
	Address string `xorm:"comment('常用地址') VARCHAR(255)"`
}
