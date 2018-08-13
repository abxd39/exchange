package client

import (
	"digicon/wallet_service/model"
	"digicon/wallet_service/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"

	"digicon/common/convert"
	"fmt"
	"github.com/ouqiang/timewheel"
	"strconv"
	"strings"
	"time"
)

type ListUnspentResult struct {
	TxID          string  `json:"txid"`
	Vout          uint32  `json:"vout"`
	Address       string  `json:"address"`
	Account       string  `json:"account"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	RedeemScript  string  `json:"redeemScript,omitempty"`
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	Spendable     bool    `json:"spendable"`
}

// part of the GetTransactionResult.
type GetTransactionDetailsResult struct {
	Account           string   `json:"account"`
	Address           string   `json:"address,omitempty"`
	Amount            float64  `json:"amount"`
	Category          string   `json:"category"`
	InvolvesWatchOnly bool     `json:"involveswatchonly,omitempty"`
	Fee               *float64 `json:"fee,omitempty"`
	Vout              uint32   `json:"vout"`
}

// GetTransactionResult models the data from the gettransaction command.
type GetTransactionResult struct {
	Amount          float64                       `json:"amount"`
	Fee             float64                       `json:"fee,omitempty"`
	Confirmations   int64                         `json:"confirmations"`
	BlockHash       string                        `json:"blockhash"`
	BlockIndex      int64                         `json:"blockindex"`
	BlockTime       int64                         `json:"blocktime"`
	TxID            string                        `json:"txid"`
	WalletConflicts []string                      `json:"walletconflicts"`
	Time            int64                         `json:"time"`
	TimeReceived    int64                         `json:"timereceived"`
	Details         []GetTransactionDetailsResult `json:"details"`
	Hex             string                        `json:"hex"`
}

type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type Vin struct {
	Coinbase  string     `json:"coinbase"`
	Txid      string     `json:"txid"`
	Vout      uint32     `json:"vout"`
	ScriptSig *ScriptSig `json:"scriptSig"`
	Sequence  uint32     `json:"sequence"`
	Witness   []string   `json:"txinwitness"`
}

type ScriptPubKeyResult struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex,omitempty"`
	ReqSigs   int32    `json:"reqSigs,omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses,omitempty"`
}

type Vout struct {
	Value        float64            `json:"value"`
	N            uint32             `json:"n"`
	ScriptPubKey ScriptPubKeyResult `json:"scriptPubKey"`
}

// TxRawResult models the data from the getrawtransaction command.
type TxRawResult struct {
	Hex           string `json:"hex"`
	Txid          string `json:"txid"`
	Hash          string `json:"hash,omitempty"`
	Size          int32  `json:"size,omitempty"`
	Vsize         int32  `json:"vsize,omitempty"`
	Version       int32  `json:"version"`
	LockTime      uint32 `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	BlockHash     string `json:"blockhash,omitempty"`
	Confirmations uint64 `json:"confirmations,omitempty"`
	Time          int64  `json:"time,omitempty"`
	Blocktime     int64  `json:"blocktime,omitempty"`
}

/*
	查询btc 是否有转账
*/
type BTCWatch struct {
	Url string
	//WalletTokenModel *models.WalletToken     // 钱包详情
	TxModel        *models.TokenChainInout // 链上交易记录
	TokenInoutMode *models.TokenInout      // 平台交易记录
	TokenModel     *models.Tokens          // 币种类
}

/*

 */
func (this *BTCWatch) Start() {
	//url := "http://bitcoin:bitcoin@localhost:28332/"

	fmt.Println("btc start ....")

	this.TxModel = new(models.TokenChainInout)
	this.TokenInoutMode = new(models.TokenInout)
	this.TokenModel = new(models.Tokens)

	exists, err := this.TokenModel.GetByName("BTC")
	if err != nil {
		fmt.Println("err:", err)
	}
	if !exists {
		fmt.Println("token not exists btc ...")
	}
	//fmt.Println(this.TokenModel)
	this.Url = this.TokenModel.Node

	this.Work()
}

var btcTw *timewheel.TimeWheel

func (this *BTCWatch) Work() {

	btcTw = timewheel.New(1*time.Second, 3600, func(data timewheel.TaskData) {
		btcTw.AddTimer(60*time.Second, "btc", timewheel.TaskData{})
		fmt.Println("btc watch ...")
		this.Check()
	})
	btcTw.Start()
	btcTw.AddTimer(1*time.Second, "btc", timewheel.TaskData{})

	//for {
	//	fmt.Println("btc watch ...")
	//	this.Check()
	//	time.Sleep(60 * time.Second) //
	//}
}

/*
	检查数据库中是否已经存在
*/
func (this *BTCWatch) Check() {
	resultInterface, err := this.BtcListUnspent(6, 12)
	if err != nil {
		fmt.Println(err.Error())
	}
	result, err := json.Marshal(resultInterface)
	if err != nil {
		fmt.Println(err.Error())
	}

	var listunspentResult []ListUnspentResult
	err = json.Unmarshal(result, &listunspentResult)
	if err != nil {
		fmt.Println(err.Error())
	}

	tkChainInOut := new(models.TokenChainInout)
	for i := 0; i < len(listunspentResult); i++ {
		curUnspend := listunspentResult[i]
		exists, err := tkChainInOut.TxIDExist(curUnspend.TxID)
		fmt.Println(curUnspend.TxID, " exists: ", exists)
		if err != nil {
			continue
		}
		if exists {
			continue
		} else {
			this.InsertRecord(curUnspend) //
		}

	}

}

/*
  有入账，则记录到交易数据库中
*/
func (this *BTCWatch) InsertRecord(curUnspend ListUnspentResult) {
	resultInterface, err := this.BtcGetTransaction(curUnspend.TxID)
	if err != nil {
		fmt.Println(err)
	}
	result, err := json.Marshal(resultInterface)
	if err != nil {
		fmt.Println(err.Error())
		log.Errorln(err)
	}
	var tranDetail GetTransactionResult
	err = json.Unmarshal(result, &tranDetail)
	if err != nil {
		log.Errorln(err)
	}
	tranInOutInterface, err := this.BtcDecodeRawTransaction(tranDetail.Hex)
	if err != nil {
		fmt.Println(err)
		log.Errorln(err.Error())
	}
	tranResult, err := json.Marshal(tranInOutInterface)
	if err != nil {
		log.Errorln(err.Error())
	}
	var tranInOutResult TxRawResult
	err = json.Unmarshal(tranResult, &tranInOutResult)
	if err != nil {
		log.Errorln(err)
	}

	voutResult := tranInOutResult.Vout
	var from, to string
	var value float64
	for i := 0; i < len(voutResult); i++ {
		if voutResult[i].N == 0 {
			from = strings.Join(voutResult[i].ScriptPubKey.Addresses, ",")
		}
		if voutResult[i].N == 1 {
			to = strings.Join(voutResult[i].ScriptPubKey.Addresses, ",")
			value = voutResult[i].Value
		}
	}
	intValue := convert.Float64ToInt64By8Bit(value)
	strValue := strconv.FormatInt(intValue, 10)

	wToken := new(models.WalletToken)
	err = wToken.GetByAddress(to)
	if err != nil {
		log.Errorln(err.Error())
		fmt.Println(err.Error())
		return
	}
	chainId := wToken.Tokenid
	uId := wToken.Uid
	tokenId := wToken.Tokenid
	tokenName := wToken.TokenName

	txmodel := &models.TokenChainInout{
		Txhash:    curUnspend.TxID,
		From:      from,
		To:        to,
		Value:     strValue,
		Contract:  "",
		Chainid:   chainId,
		Type:      1,
		Tokenid:   tokenId,
		TokenName: tokenName,
		Uid:       uId,
	}
	row, err := txmodel.InsertThis()
	if row <= 0 || err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
	}
	//更新完成状态
	new(models.TokenInout).BteUpdateAppleDone(curUnspend.TxID)

}

/*
   func:decoderawtransaction
   curl --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "decoderawtransaction", "params": ["hexstring"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
*/
func (this *BTCWatch) BtcDecodeRawTransaction(hexstring string) (result interface{}, err error) {
	data := make(map[string]interface{})
	data["id"] = 1
	data["jsonrpc"] = "1.0"
	data["method"] = "decoderawtransaction"
	params := []string{}
	params = append(params, hexstring)
	data["params"] = params
	return this.BtcRpcGet(data)
}

/*
 curl  --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "listunspent", "params": [6, 9999999, [] , true, { "minimumAmount": 0.005 } ] }' -H 'content-type: text/plain;' http://user:pass@127.0.0.1:8332/
 func :列出所有最近交易的(根据节点确认来查询)
*/
func (this *BTCWatch) BtcListUnspent(confirmStart, confirmEnd int) (result interface{}, err error) {
	data := make(map[string]interface{})
	data["id"] = 1
	data["jsonrpc"] = "1.0"
	data["method"] = "listunspent"

	params := []int{}
	if confirmStart <= 0 && confirmEnd <= 0 {
		params = append(params)
	} else {
		params = append(params, confirmStart, confirmEnd)
	}
	data["params"] = params
	return this.BtcRpcGet(data)
}

/*
  获取交易详情
  curl --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "gettransaction", "params": ["1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
*/
func (this *BTCWatch) BtcGetTransaction(txid string) (result interface{}, err error) {
	data := make(map[string]interface{})
	data["id"] = 1
	data["jsonrpc"] = "1.0"
	data["method"] = "gettransaction"

	params := []string{}
	params = append(params, txid)
	data["params"] = params
	return this.BtcRpcGet(data)
}

/*
	get function
*/
func (this *BTCWatch) BtcRpcGet(data map[string]interface{}) (result interface{}, err error) {
	rsp, err := utils.BtcRpcPost(this.Url, data)
	if err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
		return
	}
	ret := make(map[string]interface{})
	err = json.Unmarshal(rsp, &ret)
	if err != nil {
		log.Errorln(err)
		return
	}
	result, ok := ret["result"]
	if !ok {
		log.Errorln("get result error!", ok)
		return
	}
	return
}
