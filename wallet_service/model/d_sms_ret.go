package models

import (
	. "digicon/proto/common"
	"github.com/go-redis/redis"
	"fmt"
	"github.com/apex/log"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"errors"
)

//验证短信
func AuthSms(phone string, ty int32, code string) (ret int32, err error) {
	r := RedisOp{}
	auth_code, err := r.GetSmsCode(phone, ty)
	fmt.Println("验证码：",phone,ty,code,auth_code,err)
	if err == redis.Nil {
		ret = ERRCODE_SMS_CODE_NIL
		return
	} else if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}

	if code != auth_code {
		ret = ERRCODE_SMS_CODE_DIFF
		return
	}

	ret = ERRCODE_SUCCESS
	return
}

//phone, cf.SmsAccount, cf.SmsPwd, content, cf.SmsWebUrl
func SendInterSms(phone, content string) (ret int32, err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"phone": phone,
				"code":  content,
			}).Errorf("SendSms error %s", err.Error())
		}
	}()

	params := make(map[string]interface{})
	params["account"] = "I1757342"
	params["password"] = "i1PYZXVaWt2de6"
	// 手机号码，格式(区号+手机号码)，例如：8615800000000，其中86为中国的区号
	params["mobile"] = phone
	params["msg"] = content
	bytesData, err := json.Marshal(params)
	if err != nil {
		return
	}

	reader := bytes.NewReader(bytesData)
	url := "http://intapi.253.com/send/json"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	fmt.Println("-------------发送结果啊：",resp,err)
	if err != nil {
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	type SmsRet struct {
		Code      string `json:"code"`
		MessageId string `json:"msg_id"`
		Time      string `json:"time"`
		ErrorMsg  string `json:"error_msg"`
	}

	p := &SmsRet{}
	err = json.Unmarshal([]byte(respBytes), p)
	if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}

	g, err := strconv.Atoi(p.Code)
	if err != nil {
		ret = int32(g)
		return
	}

	switch g {
	case 0:
		ret = ERRCODE_SUCCESS
		return
	case 108:
		ret = ERRCODE_SMS_PHONE_FORMAT
		return
	default:
		ret = ERRCODE_UNKNOWN
		err = errors.New(fmt.Sprintf("code:%s,msg=%s", p.Code, p.ErrorMsg))
		return
	}

	return
}
