package models

import (
	. "digicon/wallet_service/utils"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
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

	tkModel := new(Tokens)
	tkModel.GetByName("btc")
	url := tkModel.Node

	err = BtcWalletPhrase(url, password, 1*60*60)
	if err != nil {
		msg := "钱包解锁失败!"
		log.Errorln(msg)
		fmt.Println(msg)
		return
	}

	address, err := BtcGetNewAddress(url, string(userId))
	if err != nil {
		msg := "生成地址错误!"
		log.Errorln(msg)
		fmt.Println(msg)
		return
	}

	privateKey, err := BtcDumpPrivKey(url, address)
	if err != nil {
		msg := "获取地址私钥错误"
		log.Errorln("获取地址", address, "私钥错误!")
		fmt.Println(msg)
		return
	}

	walletTokenModel.Address = address
	walletTokenModel.Privatekey = privateKey

	type mbtcWallet struct {
		Address    string
		Privatekey string
	}

	btcWallet := new(mbtcWallet)
	btcWallet.Address = address
	btcWallet.Privatekey = privateKey
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
	walletToken := new(WalletToken)
	walletToken.GetByUidTokenid(userId, tokenId)
	return walletToken.Address, err
}

//
// btc send to address
//
func BtcSendToAddress(toAddress string, mount string, tokenId int32, uid int, applyid int) (string, error) {
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
		log.Errorln(msg)
		fmt.Println(msg)
		return "", nil
	}

	enough, err := BtcCheckBalance(int32(uid), mount)
	if !enough {
		msg := "balance not enough!"
		err = errors.New(msg)
		log.Errorln(msg)
		return "", err
	}
	//fmt.Println("btc send before ...")
	txHash, err := BtcSendToAddressFunc(url, toAddress, mount)
	if err != nil {
		fmt.Println(err.Error())
		log.Errorf(err.Error())
		return "", err
	}
	//amount, err := convert.StringToInt64By8Bit(mount)
	//if err != nil {
	//	fmt.Println(err)
	//}
	tio := new(TokenInout) //
	//更新
	row, err := tio.UpdateApplyTiBi(applyid, txHash)
	//row, err := tio.BtcInsert(txHash, wToken.Address, toAddress, "BTC", amount,
	//	wToken.Chainid, int(tokenId), 0, int(uid),
	//)

	if err != nil || row <= 0 {
		log.Errorln(err.Error())
		fmt.Println(err.Error())
	}

	return txHash, err
}

/*
	btc tibi
*/
func BtcTiBiToAddress(toAddress string, mount string, TokenId int32, uid int32, applyid int) (string, error) {
	//fmt.Println(toAddress, mount, TokenId, uid)
	txhash, err := BtcSendToAddress(toAddress, mount, TokenId, int(uid), applyid)
	return txhash, err
}

/*
	检查余额
*/
func BtcCheckBalance(uid int32, amount string) (bool, error) {
	fmt.Println(uid, amount)

	return true, nil
}
