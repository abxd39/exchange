package models

import (
	"digicon/common/convert"
	. "digicon/wallet_service/utils"
	"encoding/json"
	"errors"
	"fmt"
	//"github.com/micro/go-micro/errors"
)

// 创建btc 钱包
func NewBTC(userId int, tokenId int, password string, chainId int) (addr string, err error) {
	var walletTokenModel = WalletToken{
		Uid:      userId,
		Password: password,
		Tokenid:  tokenId,
		Type:     "btc",
		Chainid:  chainId,
	}

	url := "http://localhost:18332"
	user := "bitcoin"
	pass := "bitcoin"
	btcClient, err := new(BtcClient).NewClient(url, user, pass)
	if err != nil {
		return
	}



	err = btcClient.WalletPassphrase(password, 1*60*60)
	if err != nil {
		msg := "钱包解锁失败!"
		Log.Errorln(msg)
		fmt.Println(msg)
		return
	}

	//btcClient.ListUnspent()
	//btcClient.GetTransaction()
	//btcClient.DecodeRawTransaction()
	address, err := btcClient.GetNewAddress(string(userId))
	if err != nil {
		msg := "生成地址错误!"
		Log.Errorln(msg)
		fmt.Println(msg)
		return
	}
	privateKey, err := btcClient.DumpPrivKey(address)
	if err != nil {
		msg := "获取地址私钥错误"
		Log.Errorln("获取地址", address.String(), "私钥错误!")
		fmt.Println(msg)
		return
	}

	walletTokenModel.Address = address.String()
	walletTokenModel.Privatekey = privateKey.String()

	type mbtcWallet struct {
		Address    string
		Privatekey string
	}
	btcWallet := new(mbtcWallet)
	btcWallet.Address = address.String()
	btcWallet.Privatekey = privateKey.String()

	strBtcWallet, _ := json.Marshal(btcWallet)

	walletTokenModel.Keystore = string(strBtcWallet)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	err = walletTokenModel.Create()
	if err != nil {
		fmt.Println("create btc token error")
	}
	return walletTokenModel.Address, err
}

//
// btc send to address
//
func BtcSendToAddress(toAddress string, mount string, tokenId int32, uid int) (string, error) {
	wToken := new(WalletToken)
	wToken.GetByUid(uid)
	password := wToken.Password

	token := Tokens{}
	token.GetByid(int(tokenId))
	url := token.Node

	//fmt.Println("----------------------------")
	err := BtcWalletPhrase(url, password, 1*60*60)

	if err != nil {
		msg := "钱包解锁失败!"
		Log.Errorln(msg)
		fmt.Println(msg)
		return "", nil
	}

	enough, err := BtcCheckBalance(int32(uid), mount)
	if !enough {
		msg := "balance not enough!"
		err = errors.New(msg)
		Log.Errorln(msg)
		return "", err
	}
	//fmt.Println("btc send before ...")
	txHash, err := BtcSendToAddressFunc(url, toAddress, mount)
	if err != nil {
		fmt.Println(err.Error())
		Log.Errorf(err.Error())
		return "", err
	}
	amount, err := convert.StringToInt64By8Bit(mount)
	if err != nil {
		fmt.Println(err)
	}
	tio := new(TokenInout) //
	row, err := tio.BtcInsert(txHash, wToken.Address, toAddress, "BTC", amount,
		wToken.Chainid, int(tokenId), 0, int(uid),
	)

	if err != nil || row <= 0 {
		Log.Errorln(err.Error())
		fmt.Println(err.Error())
	}

	return txHash, err
}

/*
	btc tibi
*/
func BtcTiBiToAddress(toAddress string, mount string, TokenId int32, uid int32) (string, error) {
	//fmt.Println(toAddress, mount, TokenId, uid)
	txhash, err := BtcSendToAddress(toAddress, mount, TokenId, int(uid))
	return txhash, err
}

/*
	检查余额
*/
func BtcCheckBalance(uid int32, amount string) (bool, error) {
	fmt.Println(uid, amount)

	return true, nil
}
