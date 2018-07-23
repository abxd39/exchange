package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type USDTClient struct {
}

/*
	usdt new address
	exp:curl --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "getnewaddress", "params": [] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
*/
func USDTNewAddress(url string) (address string, err error) {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "getnewaddress"
	params := []string{}
	data["params"] = params
	result, err := UsdtRpcPost(url, data)
	if err != nil {
		return "", err
	}
	address, ok := result.(string)
	if !ok {
		msg := "get new address error!"
		err = errors.New(msg)
		return "", err
	}
	return address, nil
}

/*
	func: USDT wallet Phrase usdt 解锁
*/
func USDTWalletPhrase(url string, pass string, keepTime int64) error {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	//data["method"] =

	return nil
}

/*
	usdt send to address
	exp: curl  --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "sendtoaddress", "params": ["1M72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd", 0.1, "donation", "seans outpost"] }' -H 'content-type: text/plain;' http://127.0.0.1:8332/
*/
func USDTSendToAddress(url string, address string, mount string) (string, error) {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "sendtoaddress"
	params := make([]interface{}, 0, 2)
	params = append(params, address)
	params = append(params, mount)
	data["params"] = params

	result, err := UsdtRpcPost(url, data)
	if err != nil {
		return "", err
	}
	txHash, ok := result.(string)
	if !ok {
		msg := "get transaction txid error!"
		err = errors.New(msg)
		return "", err
	}
	return txHash, nil
}

/*
	usdt rpc
*/
func UsdtRpcPost(url string, send map[string]interface{}) (result interface{}, err error) {
	bytesData, err := json.Marshal(send)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println("rpc post:", err.Error())
		return "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	ret := make(map[string]interface{})
	err = json.Unmarshal(respBytes, &ret)
	if err != nil {
		Log.Errorln(err)
		return "", err
	}
	result, ok := ret["result"]
	if !ok {
		msg := "get result error!"
		Log.Errorln(msg, ok)
		err = errors.New(msg)
		return
	}
	return
}
