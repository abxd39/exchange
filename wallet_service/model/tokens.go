package models

import(
	. "digicon/wallet_service/utils"
)

type Tokens struct {
	Id      int    `xorm:"not null pk INT(11)"`
	Name    string `xorm:"comment('货币名称') VARCHAR(180)"`
	Mark    string `xorm:"comment('英文标识') CHAR(30)"`
	Chainid int    `xorm:"comment('network id') INT(20)"`
	Eip155  int    `xorm:"default 1 comment('是否eip155加密') TINYINT(1)"`
	Type    string `xorm:"comment('钱包生成方式(eth,btc)') CHAR(9)"`
	Data    string `xorm:"comment('合约地址') CHAR(60)"`
}

func (this *Tokens) GetByid()(bool, error){
	return Engine.Id(this.Id).Get(this)
}