package client

import (
	"bytes"
	. "digicon/wallet_service/model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"unsafe"
	"github.com/ouqiang/timewheel"
	"time"
)

//查询methodid 0x70a08231
//data:"0x70a0823100000000000000000000000000c097f24ae4dd09359f87d85bc883a72a5a46c7"
//to:"0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0"

//转账MethodID: 0xa9059cbb
//[0]:  00000000000000000000000000c097f24ae4dd09359f87d85bc883a72a5a46c7
//[1]:  0000000000000000000000000000000000000000000000000de0b6b3a7640000
type Watch struct {
	Url         string
	Chainid     int
	BlockNumber int //当前处理到的高度

	//BlockModel *Blocks
	//KeystoreModel *Keystores
	WalletTokenModel *WalletToken     //钱包详情
	TxModel          *TokenChainInout //链上交易记录
	TokenInoutModel  *TokenInout      //平台交易记录
	TokenModel       *Tokens          //币种类
	ContextModel     *Context         //处理上下文
}

func (this *Watch) Start() {

	//查询ETH节点
	var data = new(Tokens)
	bool,er := data.GetByName("ETH")
	if bool != true || er != nil {
		fmt.Println("start fail")
		return
	}

	this.Url = data.Node

	//model初始化
	//this.WalletToken = new(Blocks)
	this.WalletTokenModel = new(WalletToken)
	this.TxModel = new(TokenChainInout)
	this.TokenInoutModel = new(TokenInout)
	this.TokenModel = new(Tokens)
	this.ContextModel = new(Context)
	//获取chainid
	var err error
	this.Chainid, err = this.getChainid()
	if err != nil {
		fmt.Println(err)
		return
	}

	this.BlockNumber, _ = this.ContextModel.MaxNumber(this.Url, this.Chainid)
	//this.BlockNumber=2464711
	this.Work()

}

var ethTw *timewheel.TimeWheel

func (this *Watch) Work() {
	ethTw = timewheel.New(1 * time.Second, 3600, func(data timewheel.TaskData) {
		ethTw.AddTimer(60 * time.Second, "eth", timewheel.TaskData{})
		fmt.Println("eth watch ...")
		this.WorkDone()
	})
	ethTw.Start()
	ethTw.AddTimer(1 * time.Second, "eth", timewheel.TaskData{})
}

func (this *Watch) WorkDone() {
	//num,err:=this.Get_balance("0x8e430b7fc9c41736911e1699dbcb6d4753cbe3b6")
	//当前最高块
	temp, err := this.Get_blockNumber()
	hight := int(temp)
	if err != nil {
		fmt.Println(err)
		return
	}
	if this.BlockNumber <= 0 {
		this.BlockNumber = hight - 10
	}

	if this.BlockNumber < hight-6 {
		for i := this.BlockNumber + 1; i <= hight-6; i++ {
			this.Blockhander(i)
			//记录当前进度
			this.ContextModel.Save(this.Url, this.Chainid, i)

		}
	}

	fmt.Println("watch.work end", hight, err)

}

func (this *Watch) Get_tokenBlance(addr string, contract string) {

}
func (this *Watch) Blockhander(num int) error {
	ret, err := this.GetblockBynumber(num)
	if err != nil {
		return err
	}
	var block map[string]interface{}
	//fmt.Println(string(ret))
	json.Unmarshal(ret, &block)
	txs := block["result"].(map[string]interface{})["transactions"].([]interface{})
	//fmt.Println(txs)
	for i := 0; i < len(txs); i++ {
		tx := txs[i].(map[string]interface{})
		if tx["to"] == nil { //部署合约交易直接跳过
			continue
		}

		//检查eth转账
		ext, err := this.WalletTokenModel.AddrExist(tx["to"].(string), this.Chainid, "")

		if err != nil {
			fmt.Println(err)
			return err
		}
		if ext {
			fmt.Println("发现一个eth转账")
			//TODO:
			this.newOrder(this.WalletTokenModel.Uid, tx["from"].(string), tx["to"].(string), this.Chainid, "", tx["value"].(string), tx["hash"].(string))

			continue
		}

		input := tx["input"].(string)
		//不是token转账跳过
		if strings.Count(input, "") < 138 || strings.Compare(input[0:10], "0xa9059cbb") != 0 {
			//fmt.Println(input)
			continue
		}

		ext, err = this.WalletTokenModel.AddrExist(fmt.Sprintf("0x%s", input[34:74]), this.Chainid, tx["to"].(string))
		//fmt.Println(ext,err,this.WalletTokenModel)
		if !ext {
			continue
		}
		var vstart int
		for i := 74; i < 138; i++ {
			if input[i:i+1] != "0" {
				vstart = i
				break
			}
		}
		if vstart == 0 {
			continue
		}
		fmt.Println("发现一个token转账")

		ok, err := this.newOrder(this.WalletTokenModel.Uid, tx["from"].(string), fmt.Sprintf("0x%s", input[34:74]), this.Chainid, tx["to"].(string), fmt.Sprintf("0x%s", input[vstart:138]), tx["hash"].(string))
		fmt.Println(ok, err)
		continue

	}
	return nil
	//fmt.Println(num,block["result"].(map[string]interface{})["transactions"].([]interface{})[1].(map[string]interface{})["to"],err)
}
func (this *Watch) newOrder(uid int, from string, to string, chainid int, contract string, value string, txhash string) (bool, error) {

	//交易是否已经收录
	exist, err := this.TxModel.TxhashExist(txhash, this.Chainid)

	if err != nil {
		return false, err
	}
	if exist {
		return false, errors.New("tx already exists")
	}

	//
	tokenid, err := this.TokenModel.GetidByContract(contract, this.Chainid)
	deci, _ := this.TokenModel.GetDecimal(tokenid)
	fmt.Println(tokenid, err)
	if err != nil {

		return false, err
	}
	if tokenid == 0 {
		return false, errors.New("token not exist")
	}
	fmt.Println("this.TxModel.Insert")
	this.TxModel.Insert(txhash, from, to, value, contract, chainid, uid, this.TokenModel.Id, this.TokenModel.Mark)

	this.TokenInoutModel.Insert(txhash, from, to, value, contract, chainid, uid, this.TokenModel.Id, this.TokenModel.Mark, deci)
	return true, nil
}
func (this *Watch) GetblockBynumber(num int) ([]byte, error) {
	send := make(map[string]interface{})
	send["jsonrpc"] = "2.0"
	send["method"] = "eth_getBlockByNumber"
	strconv.FormatInt(int64(num), 16)
	//str:=fmt.Sprintf("0x%s",strconv.FormatInt(int64(num),16))
	fmt.Println(num, fmt.Sprintf("0x%s", strconv.FormatInt(int64(num), 16)))
	send["params"] = []interface{}{fmt.Sprintf("0x%s", strconv.FormatInt(int64(num), 16)), true}
	send["id"] = this.Chainid
	rsp, err := this.post(send)
	//str := (*string)(unsafe.Pointer(&rsp))
	//fmt.Println(*str)

	return rsp, err
}

//获取区块高度
func (this *Watch) Get_blockNumber() (int64, error) {
	send := make(map[string]interface{})
	send["jsonrpc"] = "2.0"
	send["method"] = "eth_blockNumber"
	send["params"] = []string{}
	send["id"] = this.Chainid
	rsp, err := this.post(send)
	str := (*string)(unsafe.Pointer(&rsp))
	fmt.Println(*str)
	if err != nil {
		return 0, err
	}
	//
	data := make(map[string]interface{})
	err = json.Unmarshal(rsp, &data)
	if err != nil {
		return 0, err
	}
	result, ok := data["result"]
	if !ok {
		return 0, nil
	}
	var balance string
	balance, ok = result.(string)
	number, err := strconv.ParseInt(balance, 0, 64)
	//fmt.Println(data["result"],err)

	return number, nil
}

//余额查询
func (this *Watch) Get_balance(address string) (int64, error) {
	send := make(map[string]interface{})
	send["jsonrpc"] = "2.0"
	send["method"] = "eth_getBalance"
	send["params"] = []string{address, "latest"}
	send["id"] = 1
	rsp, err := this.post(send)
	str := (*string)(unsafe.Pointer(&rsp))
	fmt.Println(*str)
	if err != nil {
		return 0, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(rsp, &data)
	if err != nil {
		return 0, err
	}
	result, ok := data["result"]
	if !ok {
		return 0, nil
	}
	var balance string
	balance, ok = result.(string)
	number, err := strconv.ParseInt(balance, 0, 64)
	//fmt.Println(data["result"],err)

	return number, nil
	//return int(number/int64(1 000 000 000 000 000 000)), nil
}

func (this *Watch) getChainid() (int, error) {
	send := make(map[string]interface{})
	send["jsonrpc"] = "2.0"
	send["method"] = "net_version"
	send["params"] = []string{}
	send["id"] = 67
	rsp, err := this.post(send)
	str := (*string)(unsafe.Pointer(&rsp))
	fmt.Println(*str)
	if err != nil {
		return 0, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(rsp, &data)
	if err != nil {
		return 0, err
	}
	result, ok := data["result"]
	if !ok {
		return 0, nil
	}
	var balance string
	balance, ok = result.(string)
	number, err := strconv.ParseInt(balance, 0, 64)
	//fmt.Println(data["result"],err)

	return int(number), nil

}

//post RPC数据
func (this *Watch) post(send map[string]interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(send)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	url := this.Url
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	//byte数组直接转成string，优化内存
	return respBytes, nil
	//str := (*string)(unsafe.Pointer(&respBytes))
	//fmt.Println(*str)
}

// curl -X POST https://rinkeby.infura.io/mew --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x8e430b7fc9c41736911e1699dbcb6d4753cbe3b6", "latest"],"id":1}'
func (this *Watch) SamplePost() {
	send := make(map[string]interface{})
	send["jsonrpc"] = "2.0"
	send["method"] = "eth_getBalance"
	send["params"] = []string{"0x8e430b7fc9c41736911e1699dbcb6d4753cbe3b6", "latest"}
	send["id"] = 1
	bytesData, err := json.Marshal(send)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	url := "https://rinkeby.infura.io/mew"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)
}
