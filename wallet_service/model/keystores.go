package models

//
//import (
//	"time"
//	. "digicon/wallet_service/utils"
//	. "github.com/ethereum/go-ethereum/cmd/wallet"
//)
//
//type Keystores struct {
//	Id         int       `xorm:"not null pk autoincr INT(11)"`
//	Userid     int       `xorm:"not null comment('用户id') INT(11)"`
//	Tokenid    int       `xorm:"not null comment('币id') INT(11)"`
//	Keystore   string    `xorm:"not null comment('钱包') VARCHAR(1024)"`
//	Password   string    `xorm:"not null comment('解锁密码') CHAR(20)"`
//	Privatekey string    `xorm:"not null comment('私钥') VARCHAR(120)"`
//	Address    string    `xorm:"not null comment('钱包地址') CHAR(100)"`
//	Updatetime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
//	Createtime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
//	Type       string    `xorm:"not null comment('钱包类型（eth,btc）') CHAR(20)"`
//}
//
//func (this *Keystores)Create() error{
//	this.Updatetime = time.Now()
//	this.Createtime = time.Now()
//	_,err := Engine.Insert(this)
//	return err
//}
//
//func (this *Keystores)AddrExist(addr string)(bool,error){
//	//Engine.ShowSQL(true)
//	return Engine.Where("address=?",addr).Get(this)
//}
//func (this *Keystores)Signtx(nonce int,to string,mount int,gasprice int)( []byte,error){
//	//func Signtx(key *keystore.Key,nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int) ([]byte,error)
//	key,err :=Unlock_keystore([]byte(this.Keystore),this.Password)
//	if err != nil {
//		return nil,err
//	}
//	token := &Tokens{Id:this.Tokenid}
//	ok,err:=token.GetByid()
//	if !ok{
//		return nil,err
//	}
//	var chainid int
//	if token.Eip155>0{
//		chainid=token.Chainid
//	}
//	gaslimit :=235600
//	var data string
//	return Signtx(key,nonce,to,mount,gaslimit,gasprice,data,chainid)
//}
////创建以太坊钱包
//func Neweth(userid int,tokenid int,password string)(addr string,err error){
//	var keystoreModel = Keystores{Userid:userid,Password:password,Tokenid:tokenid,Type:"eth"}
//
//	keystoreModel.Address,keystoreModel.Keystore,keystoreModel.Privatekey,err =New_keystore(password)
//	if err != nil {
//		return "",err
//	}
//	err = keystoreModel.Create()
//	return keystoreModel.Address,err
//}
//
