package model

type QuenesConfig struct {
	TokenId      int `xorm:"unique(union_quene_id) INT(11)"`
	TokenTradeId int `xorm:"unique(union_quene_id) INT(11)"`
	Switch       int `xorm:"TINYINT(4)"`
}
