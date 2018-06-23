package models

type TibiAddress struct {
	Id      int    `xorm:"not null pk autoincr INT(11)"`
	Uid     int    `xorm:"not null comment('用户id') INT(11)"`
	TokenId int    `xorm:"not null comment('币种id') INT(11)"`
	Address string `xorm:"not null comment('地址') VARCHAR(60)"`
	Mark    string `xorm:"not null default '' comment('备注') VARCHAR(255)"`
}
