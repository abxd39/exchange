package models

import (
	. "digicon/wallet_service/utils"
)

type Blocks struct {
	Id     int    `xorm:"not null pk INT(11)"`
	Number int    `xorm:"comment('区块高度') INT(11)"`
	Txs    string `xorm:"comment('包含的交易') TEXT"`
	Hash   string `xorm:"comment('区块hash') VARCHAR(48)"`
}

func (this *Blocks) Max_number() (int, error) {
	_, err := Engine.Select("max(number) number").Get(this)
	if err != nil {
		return 0, nil
	}
	//fmt.Println(ok,err,this)
	return this.Number, nil
}
