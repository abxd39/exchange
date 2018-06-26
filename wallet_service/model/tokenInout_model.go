package models

import (
	"digicon/wallet_service/utils"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

type TokenInout struct {
	Id          int       `xorm:"not null pk INT(11)"`
	Txhash      string    `xorm:"not null comment('交易hash') VARCHAR(200)"`
	From        string    `xorm:"not null comment('打款方') VARCHAR(42)"`
	To          string    `xorm:"not null comment('收款方') VARCHAR(42)"`
	Amount      int64     `xorm:"not null comment('金额') BIGINT(20)"`
	Value       string    `xorm:"comment('原始16进制转账数据') VARCHAR(32)"`
	Chainid     int       `xorm:"not null comment('链id') INT(11)"`
	Contract    string    `xorm:"not null default '' comment('合约地址') VARCHAR(42)"`
	Tokenid     int       `xorm:"not null comment('币种id') INT(11)"`
	States      int       `xorm:"comment('未处理0，已处理1') TINYINT(1)"`
	TokenName   string    `xorm:"comment('币种名称') VARCHAR(10)"`
	Uid         int       `xorm:"comment('用户id') INT(11)"`
	CreatedTime time.Time `xorm:"comment('创建时间') TIMESTAMP"`
}

func (this *TokenInout) Insert(txhash, from, to, value, contract string, chainid int, uid int, tokenid int, tokenname string, decim int) (int, error) {
	this.Id = 0
	this.Txhash = txhash
	this.From = from
	this.To = to
	this.Value = value
	temp, _ := new(big.Int).SetString(value[2:], 16)
	amount := decimal.NewFromBigInt(temp, int32(8-decim)).IntPart()
	this.Amount = amount
	this.Contract = contract
	this.Chainid = chainid
	this.Tokenid = tokenid
	this.TokenName = tokenname
	this.Uid = uid
	affected, err := utils.Engine_wallet.InsertOne(this)
	return int(affected), err
}
