package models
//
//import(
//	. "digicon/wallet_service/utils"
//)
//type Transations struct {
//	Id       int    `xorm:"INT(11)"`
//	Txhash   string `xorm:"comment('交易hash') VARCHAR(58)"`
//	From     string `xorm:"comment('打款地址') VARCHAR(40)"`
//	To       string `xorm:"comment('付款地址') VARCHAR(40)"`
//	Value    string `xorm:"comment('金额') VARCHAR(30)"`
//	Contract string `xorm:"comment('合约地址') VARCHAR(20)"`
//	Chainid  int    `xorm:"comment('链id') INT(11)"`
//}
//func (this *Transations)Insert(txhash,from,to,value,contract string,chainid int) (int,error){
//	this.Id=0
//	this.Txhash=txhash
//	this.From = from
//	this.To = to
//	this.Value =value
//	this.Contract=contract
//	this.Chainid=chainid
//	affected,err:=Engine_wallet.InsertOne(this)
//	return int(affected),err
//}