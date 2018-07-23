package model

type Tokens struct {
	Id     int64  `xorm:"id pk autoincr"`
	Name   string `xorm:"name"`
	CnName string `xorm:"cn_name"`
}

func (*Tokens) TableName() string {
	return "tokens"
}
