package errors

import (
	. "digicon/proto/common"
	goErrors "errors"
)

// 兼容golang errors对象
func New(text string) error {
	return goErrors.New(text)
}

func GetErrStatus(err interface{}) int {
	switch v := err.(type) {
	case SysErrorInterface:
		return v.Status()
	case NormalErrorInterface:
		return v.Status()
	default:
		return ERRCODE_NORMAL_ERROR
	}
}

func GetErrMsg(err interface{}) string {
	switch v := err.(type) {
	case SysErrorInterface:
		return v.Error()
		//return v.String() // todo 根据API_ENV显示错误
	case error:
		return v.Error()
	case NormalErrorInterface:
		return v.Error()
	default:
		return ""
	}
}
