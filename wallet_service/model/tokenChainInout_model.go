package models

import (
	"digicon/wallet_service/utils"
	"time"
)

// 链上，充值，提币
type TokenChainInout struct {
	Id          int       `xorm:"INT(11)"`
	Txhash      string    `xorm:"comment('交易hash') VARCHAR(255)"`
	From        string    `xorm:"comment('打款地址') VARCHAR(42)"`
	To          string    `xorm:"comment('付款地址') VARCHAR(42)"`
	Value       string    `xorm:"comment('金额') VARCHAR(30)"`
	Contract    string    `xorm:"comment('合约地址') VARCHAR(42)"`
	Chainid     int       `xorm:"comment('链id') INT(11)"`
	Type        int       `xorm:"not null comment('平台转出:0,充值:1') INT(11)"`
	Signtx      string    `xorm:"comment('平台转出记录交易签名') VARCHAR(1024)"`
	Tokenid     int       `xorm:"not null comment('币种id') INT(11)"`
	TokenName   string    `xorm:"not null comment('币名称') VARCHAR(10)"`
	Uid         int       `xorm:"not null comment('用户id') INT(11)"`
	CreatedTime time.Time `xorm:"default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	Address        string    `xorm:"comment('交易账号(2018-8-13新增)') VARCHAR(42)"`
	Gas       int64    `xorm:"comment('gas数量') VARCHAR(30)"`
	Gas_price       int64    `xorm:"comment('gas价格,单位：wei') VARCHAR(30)"`
	Real_fee       int64    `xorm:"comment('实际消耗手续费') VARCHAR(30)"`
}

func (this *TokenChainInout) Insert(txhash, from, to, value, contract string, chainid int, uid int, tokenid int, tokenname string,opt int,gas int64,gasprice int64,realfee int64) (int, error) {
	this.Id = 0
	this.Txhash = txhash
	this.From = from
	this.To = to
	this.Value = value
	this.Contract = contract
	this.Chainid = chainid
	this.Type = opt

	this.Tokenid = tokenid
	this.TokenName = tokenname
	this.Uid = uid
	this.Gas = gas
	this.Gas_price = gasprice
	this.Real_fee = realfee
	//utils.Engine_wallet.ShowSQL(true)
	affected, err := utils.Engine_wallet.InsertOne(this)
	//fmt.Println("aaaa",uid,err)
	return int(affected), err
}
func (this *TokenChainInout) TxhashExist(hash string, chainid int) (bool, error) {
	utils.Engine_wallet.ShowSQL(false)
	//return utils.Engine_wallet.Where("txhash=? and chainid=?", hash, chainid).Get(this)
	return utils.Engine_wallet.Where("txhash=?", hash).Get(this)

}

// ==================================   BTC 添加  =============================================

func (this *TokenChainInout) InsertRecord(txhash, from, to, value, contract string, insertType int, chainid int, uid int, tokenid int, tokenname string) (int, error) {
	this.Id = 0
	this.Txhash = txhash
	this.From = from
	this.To = to
	this.Value = value
	this.Contract = contract
	this.Chainid = chainid
	this.Type = insertType

	this.Tokenid = tokenid
	this.TokenName = tokenname
	this.Uid = uid
	affected, err := utils.Engine_wallet.InsertOne(this)
	return int(affected), err
}

func (this *TokenChainInout) InsertThis() (int, error) {
	affected, err := utils.Engine_wallet.InsertOne(this)
	return int(affected), err
}

func (this *TokenChainInout) TxIDExist(txhash string) (bool, error) {
	tk := &TokenChainInout{Txhash: txhash}
	return utils.Engine_wallet.Exist(tk) //如果仅仅判断某条记录是否存在，则使用Exist方法，Exist的执行效率要比Get更高。
}
