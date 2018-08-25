package models

import (
	"fmt"
	. "digicon/wallet_service/utils"
	log "github.com/sirupsen/logrus"
	"encoding/json"
)

/*
	create usdt wallet
*/

// 创建btc 钱包
func NewUSDT(userId int, tokenId int, password string, chainId int) (addr string, err error) {
	var walletTokenModel = WalletToken{
		Uid:      userId,
		Password: password,
		Tokenid:  tokenId,
		Type:     "omni",
		Chainid:  chainId,
	}

	tkModel := new(Tokens)
	tkModel.GetByName("omni")
	url := tkModel.Node

	err = UsdtWalletPhrase(url, password, 1*60*60)
	if err != nil {
		msg := "钱包解锁失败!"
		log.Errorln(msg)
		fmt.Println(msg)
		return
	}

	address, err := UsdtGetNewAddress(url, string(userId))
	if err != nil {
		msg := "生成地址错误!"
		log.Errorln(msg)
		fmt.Println(msg)
		return
	}

	privateKey, err := UsdtDumpPrivKey(url, address)
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
