package model

import (
	. "digicon/user_service/dao"
	log "github.com/sirupsen/logrus"
)

/*
type TokenList struct {
	TokenId int    `xorm:"not null pk autoincr INT(11)"`
	Name    string `xorm:"comment('货币名称') VARCHAR(64)"`
	Detail  string `xorm:"comment('详情地址') VARCHAR(255)"`
	Mark    string `xorm:"comment('英文标识') CHAR(10)"`
	Logo    string `xorm:"comment('货币logo') VARCHAR(255)"`
}
*/
type Tokens struct {
	Id        int    `xorm:"not null pk autoincr INT(11)"`
	Name      string `xorm:"default '' comment('货币名称') VARCHAR(64)"`
	Detail    string `xorm:"default '' comment('详情地址') VARCHAR(255)"`
	Signature string `xorm:"default '' comment('签名方式(eth,eip155,btc)') VARCHAR(255)"`
	Chainid   int    `xorm:"default 0 comment('链id') INT(11)"`
	Github    string `xorm:"default '' comment('项目地址') VARCHAR(255)"`
	Web       string `xorm:"default '' VARCHAR(255)"`
	Mark      string `xorm:"comment('英文标识') CHAR(10)"`
	Logo      string `xorm:"default '' comment('货币logo') VARCHAR(255)"`
	Contract  string `xorm:"default '' comment('合约地址') VARCHAR(255)"`
	Node      string `xorm:"default '' comment('节点地址') VARCHAR(100)"`
	Decimal   int    `xorm:"not null default 1 comment('精度 1个eos最小精度的10的18次方') INT(11)"`
}

func (this *Tokens) GetTokens() []Tokens {
	g := make([]Tokens, 0)
	err := DB.GetMysqlConn().Find(&g)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	return g
}
