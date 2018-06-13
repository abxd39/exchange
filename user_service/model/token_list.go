package model

type TokenList struct {
	TokenId int    `xorm:"not null pk autoincr INT(11)"`
	Name    string `xorm:"comment('货币名称') VARCHAR(64)"`
	Detail  string `xorm:"comment('详情地址') VARCHAR(255)"`
	Mark    string `xorm:"comment('英文标识') CHAR(10)"`
	Logo    string `xorm:"comment('货币logo') VARCHAR(255)"`
}
