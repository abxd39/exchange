package model

type Pays struct {
	Id     int    `xorm:"not null pk INT(11)"`
	ZhPay  string `xorm:"not null comment('中文名字') VARCHAR(10)"`
	EnPay  string `xorm:"default '' comment('英文名字') VARCHAR(10)"`
	Short  string `xorm:"default '' comment('首字母简写,排序字段') VARCHAR(10)"`
	Status int    `xorm:"default 1 comment('1 开启 0 关闭') TINYINT(4)"`
}
