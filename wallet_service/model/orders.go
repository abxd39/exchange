package models
//
//import (
//	. "digicon/wallet_service/utils"
//
//)
//type Orders struct {
//	Id       int    `xorm:"not null pk INT(11)"`
//	Txhash   string `xorm:"not null comment('交易hash') VARCHAR(200)"`
//	From     string `xorm:"not null comment('打款方') VARCHAR(42)"`
//	To       string `xorm:"not null comment('收款方') VARCHAR(42)"`
//	Amount   string `xorm:"not null comment('值') DECIMAL(64,8)"`
//	Value    string `xorm:"comment('原始16进制转账数据') VARCHAR(32)"`
//	Chainid  int    `xorm:"not null comment('链id') INT(11)"`
//	Contract string `xorm:"not null default '' comment('合约地址') VARCHAR(42)"`
//	Tokenid  int    `xorm:"not null comment('币种id') INT(11)"`
//
//}
//
//func (this *Orders)Txhash_exist(txhash string) (bool,error){
//
//
//	return false,nil
//}
//
//func (this *Orders)Insert(txhash ,from,to,amount,value string,contract string,chainid int,tokenid int) (int,error){
//	this.Id = 0
//	this.Txhash = txhash
//	this.From = from
//	this.To = to
//	this.Amount= 	amount
//	this.Value = value
//	this.Chainid = chainid
//	this.Contract = contract
//	this.Tokenid = tokenid
//	affacted,err := Engine.InsertOne(this)
//	return int(affacted),err
//}