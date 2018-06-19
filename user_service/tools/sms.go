package tools

import (
	"digicon/common/sms"
	. "digicon/proto/common"
	cf "digicon/user_service/conf"
	//"digicon/user_service/model"
	"encoding/json"
	"fmt"
	"github.com/liudng/godump"
	"strconv"
)

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
	//msg, ok := CheckErrorMessage(int32(code_))
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
	/*
		if ok {
			rcode = ERRCODE_SUCCESS
			msg = msg
		} else {
			rcode = ERRCODE_UNKNOWN
			msg = p.ErrorMsg
		}
	*/
	return
}


