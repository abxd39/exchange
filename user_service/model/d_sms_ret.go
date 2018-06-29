package model

import (
	"digicon/common/random"
	//"github.com/sirupsen/logrus"
	//. "digicon/user_service/dao"
	"bytes"
	. "digicon/proto/common"
	cf "digicon/user_service/conf"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/liudng/godump"
	"github.com/pkg/errors"
)

const (
	SMS_REGISTER        = 1 //注册业务
	SMS_FORGET          = 2
	SMS_CHANGE_PWD      = 3
	SMS_RESET_GOOGLE = 4
	SMS_RESET_TRADE_PWD = 5
	SMS_MAX             = 6

)

//发送短信
func SendSms(phone, country string, ty int32) (ret int32, err error) {
	code := random.Random6dec()
	r := &RedisOp{}
	err = r.SetSmsCode(phone, code, ty)
	if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}
	g := fmt.Sprintf("%s%s", country, phone)
	ret, err = SendInterSms(g, code)

	return
}

//验证短信
func AuthSms(phone string, ty int32, code string) (ret int32, err error) {
	godump.Dump(phone)
	godump.Dump(code)
	godump.Dump(ty)
	r := RedisOp{}
	auth_code, err := r.GetSmsCode(phone, ty)
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

//短信通用处理
func ProcessSmsLogic(ty int32, phone, region string) (ret int32, err error) {
	switch ty {
	case SMS_REGISTER:
		//TODO判断
		u := User{}
		ret, err = u.CheckUserExist(phone, "phone")
		if err != nil {
			return
		}

		if ret != ERRCODE_SUCCESS {
			return
		}

		ret, err = SendSms(phone, region, ty)
	case SMS_FORGET:
		ret, err = SendSms(phone, region, ty)
	case SMS_CHANGE_PWD:
		ret, err = SendSms(phone, region, ty)
	case SMS_RESET_GOOGLE:
		ret, err = SendSms(phone, region, ty)
	default:
		return

	}
	return

}

//phone, cf.SmsAccount, cf.SmsPwd, content, cf.SmsWebUrl
func SendInterSms(phone, code string) (ret int32, err error) {
	params := make(map[string]interface{})
	params["account"] = cf.SmsAccount
	params["password"] = cf.SmsPwd
	// 手机号码，格式(区号+手机号码)，例如：8615800000000，其中86为中国的区号
	params["mobile"] = phone
	content := fmt.Sprintf("【爱来多科技】您的验证码是：%s", code)
	params["msg"] = content
	bytesData, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err.Error())
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
	godump.Dump(p)

	g, err := strconv.Atoi(p.Code)
	if err != nil {
		ret = int32(g)
		return
	}
	if g == 0 {
		ret = ERRCODE_SUCCESS
		return
	}
	ret = ERRCODE_UNKNOWN
	err = errors.New(p.ErrorMsg)
	return
}
