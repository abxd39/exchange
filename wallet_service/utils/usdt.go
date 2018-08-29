package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

//var btcClient *rpcclient.Client
type UsdtClient struct {
}

// btc 客户端链接
//func (p *BtcClient) NewClient(host, user, pass string) (btcclient *rpcclient.Client, err error) {
//	connCfg := &rpcclient.ConnConfig{
//		Host:         host,
//		User:         user,
//		Pass:         pass,
//		HTTPPostMode: true,
//		DisableTLS:   true,
//	}
//
//	btcclient, err = rpcclient.New(connCfg, nil)
//	if err != nil {
//		log.Errorln(err.Error())
//	}
//	return
//}

/*
 curl  --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "sendtoaddress", "params": ["1M72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd", 0.1, "donation", "seans outpost"] }' -H 'content-type: text/plain;' http://user:pass@127.0.0.1:8332/
*/
func UsdtSendToAddressFunc(url string, from_address string,to_address string,propertyid int,mount string) (string, error) {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "omni_send"
	//param := []string{}
	params := make([]interface{}, 0, 2)
	params = append(params, from_address)
	params = append(params, to_address)
	params = append(params, propertyid)
	params = append(params, mount)
	data["params"] = params

	rsp, err := UsdtRpcPost(url, data)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	if gjson.Get(string(rsp),"error").String() != "" {
		msg := gjson.Get(string(rsp),"error.message").String()
		return "",errors.New(msg)
	}
	result := gjson.Get(string(rsp),"result").String()
	return result,nil
}

/*
  curl  --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "walletpassphrase", "params": ["my pass phrase", 60] }' -H 'content-type: text/plain;' http://user@pass127.0.0.1:18332/
  解锁钱包
*/
func UsdtWalletPhrase(url string, pass string, keepTime int64) error {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "walletpassphrase"
	params := make([]interface{}, 0, 2)
	params = append(params, pass)
	params = append(params, keepTime)
	data["params"] = params
	_, err := UsdtRpcPost(url, data)
	if err != nil {
		fmt.Println(err.Error())
		log.Errorf(err.Error())
		return err
	}
	return nil
}

/*
	btc get new address
	curl  --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "getnewaddress", "params": [] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
*/
func UsdtGetNewAddress(url string, account string) (address string, err error) {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "getnewaddress"
	params := []string{}
	data["params"] = params
	result, err := UsdtRpcPost(url, data)
	if err != nil {
		log.Errorln(err.Error())
		return "", err
	}
	address = gjson.Get(string(result), "result").String()
	//fmt.Println("创建比特币地址：", address, result)
	//address = string(result)
	return address, nil
}

/*
	btc dump privkey
	curl --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "dumpprivkey", "params": ["myaddress"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
*/
func UsdtDumpPrivKey(url string, myaddress string) (privateKey string, err error) {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "dumpprivkey"
	params := []string{}
	params = append(params, myaddress)
	data["params"] = params
	result, err := UsdtRpcPost(url, data)
	privateKey = gjson.Get(string(result), "result").String()
	if err != nil {
		log.Errorln(err.Error())
		return "", err
	}
	//privateKey = string(result)
	return privateKey, nil
}

//列出最近的交易记录
func UsdtListtransactions(url string) (error,string) {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "omni_listtransactions"

	params := make([]interface{}, 0, 2)
	params = append(params, "")
	params = append(params, 100)

	data["params"] = params

	result, err := BtcRpcPost(url, data)
	if err != nil {
		return err,""
	}
	if gjson.Get(string(result), "error").String() != "" {
		return errors.New("error"),""
	}
	return nil,gjson.Get(string(result), "result").String()
}

/*
	btc rpc
*/
func UsdtRpcPost(url string, send map[string]interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(send)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println("rpc post:", err.Error())
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	//fmt.Println("resp:", resp)
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

//检查usdt是否足够
func UsdtCheckBalance(uid int32, amount string) (bool, error) {
	fmt.Println(uid, amount)

	return true, nil
}

//列出最近的交易记录
func UsdtOmniGettransaction(url string,txhash string) (error,string) {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "omni_gettransaction"

	params := make([]interface{}, 0, 1)
	params = append(params, txhash)

	data["params"] = params

	result, err := USDTRpcPost(url, data)

	if err != nil {
		return err,""
	}
	if errinfo := gjson.Get(string(result), "error").String(); errinfo != "" {
		return errors.New(errinfo),""
	}
	return nil,gjson.Get(string(result), "result").String()
}

//获取钱包区块高度
func USDTGetBlockCount(url string) (error,int) {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "getblockcount"

	result, err := USDTRpcPost(url, data)
	if err != nil {
		return err,0
	}
	if errinfo := gjson.Get(string(result), "error").String(); errinfo != "" {
		return errors.New(errinfo),0
	}
	num := gjson.Get(string(result), "result").Int()
	return nil,int(num)
}

//根据区块号查询区块
func USDTOmniListblocktransactions(url string,block int) (error,[]string) {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "omni_listblocktransactions"

	params := make([]interface{}, 0, 1)
	params = append(params, block)

	data["params"] = params

	result, err := USDTRpcPost(url, data)
	if err != nil {
		return err,[]string{}
	}

	if errinfo := gjson.Get(string(result), "error").String(); errinfo != "" {
		return errors.New(errinfo),[]string{}
	}

	res := gjson.Get(string(result), "result").String()

	var tran []string
	err = json.Unmarshal([]byte(res),&tran)
	if err != nil {
		return err,[]string{}
	}

	return nil,tran
}


func USDTRpcPost(url string, send map[string]interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(send)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println("rpc post:", err.Error())
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	//fmt.Println("resp:", resp)
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
