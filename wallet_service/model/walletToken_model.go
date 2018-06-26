package models

import (
	"time"
	"errors"
	. "digicon/wallet_service/utils"
	//. "github.com/ethereum/go-ethereum/cmd/wallet"
	"math/big"
	"fmt"
)

type WalletToken struct {
	Id          int       `xorm:"not null pk autoincr INT(11)"`
	Uid         int       `xorm:"not null comment('用户id') unique(user2token) INT(11)"`
	Tokenid     int       `xorm:"not null comment('币id') unique(user2token) INT(11)"`
	TokenName   string    `xorm:"not null comment('币种名称') VARCHAR(20)"`
	Chainid     int       `xorm:"not null default 0 comment('链id') INT(11)"`
	Contract    string    `xorm:"not null default '' comment('合约地址') VARCHAR(42)"`
	Keystore    string    `xorm:"not null comment('钱包') VARCHAR(1024)"`
	Type        string    `xorm:"not null comment('钱包类型（eth,btc）') CHAR(20)"`
	Nonce       int       `xorm:"not null default 1 comment('交易高度') INT(11)"`
	Password    string    `xorm:"not null comment('解锁密码') CHAR(20)"`
	Privatekey  string    `xorm:"not null comment('私钥') VARCHAR(120)"`
	Address     string    `xorm:"not null comment('钱包地址') CHAR(100)"`
	UpdatedTime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	CreatedTime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
}



func (this *WalletToken)Create() error{
	this.UpdatedTime = time.Now()
	this.CreatedTime = time.Now()
	_,err := Engine_wallet.Insert(this)
	return err
}

func (this *WalletToken)AddrExist(addr string,chainid int,contract string)(bool,error){
	//Engine_wallet.ShowSQL(true)
	//fmt.Println(addr,chainid,contract)
	return Engine_wallet.Where("address=? and chainid=? and contract=?",addr,chainid,contract).Get(this)
}
func (this *WalletToken)Signtx(to string,mount *big.Int,gasprice int64)( []byte,error){
	//func Signtx(key *keystore.Key,nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int) ([]byte,error)
	key,err :=Unlock_keystore([]byte(this.Keystore),this.Password)
	if err != nil {
		return nil,err
	}
	token := &Tokens{Id:this.Tokenid}
	ok,err:=token.GetByid(this.Tokenid)
	if !ok{
		return nil,err
	}
	nonce := this.Nonce
	this.NonceIncr(this.Id)
	switch token.Signature {
	case "eip155":
		gaslimit :=60000
		return Signtx(key,nonce,to,mount,gaslimit,gasprice,token.Contract,this.Chainid)
	case "eth":
		gaslimit :=60000
		return Signtx(key,nonce,to,mount,gaslimit,gasprice,token.Contract,0)
	default:
		return nil, errors.New("unknow type")
	}

}
func (this *WalletToken)NonceIncr(id int){
	Engine_wallet.Exec("update wallet_token set nonce=nonce+1 where id=?",id)
	fmt.Println("update wallet_token set nonce=nonce+1 where id=?",id)
}
//创建以太坊钱包
func Neweth(userid int,tokenid int,password string,chainid int)(addr string,err error){
	var walletTokenModel = WalletToken{Uid:userid,Password:password,Tokenid:tokenid,Type:"eth",Chainid:chainid}

	walletTokenModel.Address,walletTokenModel.Keystore,walletTokenModel.Privatekey,err =New_keystore(password)
	if err != nil {
		return "",err
	}
	err = walletTokenModel.Create()
	return walletTokenModel.Address,err
}


