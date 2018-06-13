package model

type SmsRet struct {
	Code      string `json:"code"`
	MessageId string `json:"msg_id"`
	Time      string `json:"time"`
	ErrorMsg  string `json:"error_msg"`
}
