package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"unsafe"
)

func Send253Sms(phone, account, pwd, content, weburl string) (rsp string, err error) {
	params := make(map[string]interface{})
	params["account"] = account //"N2562426"
	params["password"] = pwd    //"rSLFN2Io96772f"
	params["phone"] = phone
	//params["msg"] =url.QueryEscape("【253云通讯】您好，您的验证码是999999")
	params["msg"] = url.QueryEscape(content)
	params["report"] = "false"
	bytesData, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	//weburl := "http://smssh1.253.com/msg/send/json"  //请求地址请参考253云通讯自助通平台查看或者询问您的商务负责人获取
	request, err := http.NewRequest("POST", weburl, reader)
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
	str := (*string)(unsafe.Pointer(&respBytes))
	rsp = *str
	return
}
