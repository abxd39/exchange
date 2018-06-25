package tools

import (
	. "digicon/proto/common"
	cf "digicon/user_service/conf"
	//"digicon/user_service/model"
	"encoding/json"
	"fmt"
	"github.com/liudng/godump"
	"strconv"
	"bytes"
	"net/http"
	"io/ioutil"
	"unsafe"
)
/*
func Send253YunSms(phone, code string) (rcode int32, msg string) {

	content := fmt.Sprintf("【253云通讯】您好，您的验证码是%s", code)
	ret, err := sms.Send253Sms(phone, cf.SmsAccount, cf.SmsPwd, content, cf.SmsWebUrl)
	if err != nil {
		rcode = ERRCODE_UNKNOWN
		msg = err.Error()
		return
	}

	type SmsRet struct {
		Code      string `json:"code"`
		MessageId string `json:"msg_id"`
		Time      string `json:"time"`
		ErrorMsg  string `json:"error_msg"`
	}

	p := &SmsRet{}
	err = json.Unmarshal([]byte(ret), p)
	if err != nil {
		rcode = ERRCODE_UNKNOWN
		msg = err.Error()
		return
	}

	code_, _ := strconv.Atoi(p.Code)
	godump.Dump(p)
	switch int32(code_) {
	case 0:
		rcode = ERRCODE_SUCCESS
		return
	case 109:
		rcode = ERRCODE_SMS_MONEY_ENGOUGE
		return
	case 104:
		rcode = ERRCODE_SMS_SYS_BUSY
		return
	case 103:
		rcode = ERRCODE_SMS_COMMIT_QUICK
		return
	default:
		rcode = ERRCODE_UNKNOWN
		msg = p.ErrorMsg
		return
	}

	return
}

*/
func SendInter253YunSms(phone, code string) (rcode int32, msg string) {
	content := fmt.Sprintf("【253云通讯】您好，您的验证码是%s", code)
	ret, err := sendInter253Sms(phone, cf.SmsAccount, cf.SmsPwd, content, cf.SmsWebUrl)
	if err != nil {
		rcode = ERRCODE_UNKNOWN
		msg = err.Error()
		return
	}

	type SmsRet struct {
		Code      string `json:"code"`
		MessageId string `json:"msg_id"`
		ErrorMsg  string `json:"error_msg"`
	}

	p := &SmsRet{}
	err = json.Unmarshal([]byte(ret), p)
	if err != nil {
		rcode = ERRCODE_UNKNOWN
		msg = err.Error()
		return
	}

	code_, _ := strconv.Atoi(p.Code)
	godump.Dump(p)
	switch int32(code_) {
	case 0:
		rcode = ERRCODE_SUCCESS
		return
	case 110:
		rcode = ERRCODE_SMS_MONEY_ENGOUGE
		return
	case 104:
		rcode = ERRCODE_SMS_SYS_BUSY
		return
	case 103:
		rcode = ERRCODE_SMS_COMMIT_QUICK
		return
	default:
		rcode = ERRCODE_UNKNOWN
		msg = p.ErrorMsg
		return
	}

	return
}


func sendInter253Sms(phone, account, pwd, content, weburl string) (rsp string, err error) {
	params := make(map[string]interface{})
	params["account"] =  cf.SmsAccount
	params["password"] = cf.SmsPwd
	// 手机号码，格式(区号+手机号码)，例如：8615800000000，其中86为中国的区号
	params["mobile"] = phone
	//params["msg"] =url.QueryEscape(content)
	params["msg"] =content
	bytesData, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err.Error() )
		return
	}
	reader := bytes.NewReader(bytesData)
	url := "http://intapi.253.com/send/json"
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
	str := (*string)(unsafe.Pointer(&respBytes))
	rsp = *str
	return
}