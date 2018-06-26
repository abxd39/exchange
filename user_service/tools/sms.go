package tools

import (
	"digicon/common/sms"
	. "digicon/proto/common"
	cf "digicon/user_service/conf"
	//"digicon/user_service/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liudng/godump"
	"strconv"
)

func Send253YunSms(phone, code string) (rcode int32, err error) {
	content := fmt.Sprintf("【253云通讯】您好，您的验证码是%s", code)
	ret, err := sms.Send253Sms(phone, cf.SmsAccount, cf.SmsPwd, content, cf.SmsWebUrl)
	if err != nil {
		rcode = ERRCODE_UNKNOWN

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
		return
	}

	code_, _ := strconv.Atoi(p.Code)
	godump.Dump(p)
	switch int32(code_) {
	case 0:
		rcode = ERRCODE_SUCCESS
		return
	default:
		rcode = ERRCODE_UNKNOWN
		err = errors.New(p.ErrorMsg)
		return
	}

	return
}
