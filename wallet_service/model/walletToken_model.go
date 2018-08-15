package models

import (
	. "digicon/wallet_service/utils"
	"errors"
	"math/big"
	"time"
)

// 钱包
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

//////////////// btc /////

func (this *WalletToken) GetByUid(uid int) error {
	_, err := Engine_wallet.Where("uid =?", uid).Get(this)
	if err != nil {
		return err
	}
	return nil
}

/*
 根据地址获取 钱包
*/
func (this *WalletToken) GetByAddress(address string) error {
	_, err := Engine_wallet.Where("address =?", address).Get(this)
	if err != nil {
		return err
	}
	return nil
}

/// ////////////////

func (this *WalletToken) Create() error {
	this.UpdatedTime = time.Now()
	this.CreatedTime = time.Now()
	_, err := Engine_wallet.Insert(this)
	return err
}

func (this *WalletToken) AddrExist(addr string, chainid int, contract string) (bool, error) {
	//Engine_wallet.ShowSQL(true)
	//fmt.Println(addr,chainid,contract)
	return Engine_wallet.Where("address=? and chainid=? and contract=?", addr, chainid, contract).Get(this)
}
func (this *WalletToken) Signtx(to string, mount *big.Int, gasprice int64, nonce int) ([]byte, error) {
	//func Signtx(key *keystore.Key,nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int) ([]byte,error)
	key, err := Unlock_keystore([]byte(this.Keystore), this.Password)
	if err != nil {
		return nil, err
	}
	token := &Tokens{Id: this.Tokenid}
	ok, err := token.GetByid(this.Tokenid)
	if !ok {
		return nil, err
	}
	//nonce := this.Nonce
	this.NonceIncr(this.Id)
	switch token.Signature {
	case "eip155":
		gaslimit := 60000
		return Signtx(key, nonce, to, mount, gaslimit, gasprice, token.Contract, this.Chainid)
	case "eth":
		gaslimit := 60000
		return Signtx(key, nonce, to, mount, gaslimit, gasprice, token.Contract, 0)
	default:
		return nil, errors.New("unknow type")
	}

}
func (this *WalletToken) NonceIncr(id int) {
	Engine_wallet.Exec("update wallet_token set nonce=nonce+1 where id=?", id)
	//fmt.Println("update wallet_token set nonce=nonce+1 where id=?", id)
}

//创建以太坊钱包
func Neweth(userid int, tokenid int, password string, chainid int) (addr string, err error) {
	var walletTokenModel = WalletToken{Uid: userid, Password: password, Tokenid: tokenid, Type: "eth", Chainid: chainid}

	walletTokenModel.Address, walletTokenModel.Keystore, walletTokenModel.Privatekey, err = New_keystore(password)
	if err != nil {
		return "", err
	}
	err = walletTokenModel.Create()
	return walletTokenModel.Address, err
}

func (this *WalletToken) GetByUidTokenid(uid int, tokenid int) error {
	_, err := Engine_wallet.Where("uid = ? and tokenid = ?", uid, tokenid).Get(this)
	if err != nil {
		return err
	}
	return nil
}

func (this *WalletToken) WalletTokenExist(uid int, tokenid int) (bool, string, string) {
	tk := &WalletToken{Uid: uid, Tokenid: tokenid}
	boo, err := Engine_wallet.Exist(tk)
	if err != nil {
		return false, "", ""
	}
	if boo == true {
		//存在
		walletToken := new(WalletToken)
		errr := walletToken.GetByUidTokenid(uid, tokenid)
		if errr != nil {
			return false, "", ""
		}
		tokenModel := &Tokens{Id: tokenid}
		_, err = tokenModel.GetByid(tokenid)
		return true, walletToken.Address, tokenModel.Signature
	}
	return false, "", ""
}

//查询所有比特币地址
func (this *WalletToken) GetAllAddress() (err error,data []WalletToken) {
	data = make([]WalletToken,0)
	sql := `select created_time,address,g_common.tokens.contract contract from wallet_token left join g_common.tokens on wallet_token.tokenid = g_common.tokens.id`
	err = Engine_wallet.SQL(sql).Find(&data)
	//err = Engine_wallet.Select("chainid,address,contract,created_time").Find(&data)
	return
}

//查询注册时间大于某个点的地址
func (this *WalletToken) GetAddressByTime(time string) (err error,data []WalletToken) {
	data = make([]WalletToken,0)
	err = Engine_wallet.Select("chainid,address,contract,created_time").Where("created_time > ?",time).Find(&data)
	return
}

//根据地址获取uid
func (this *WalletToken) GetUidByAddress(address string) error {
	_, err := Engine_wallet.Select("uid").Where("address =?", address).Get(this)
	if err != nil {
		return err
	}
	return nil
}
