package errors

import (
	. "digicon/proto/common"
	goErrors "errors"
	"os"
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
		if os.Getenv("API_ENV") == "prod" { //生产环境只显示概略错误信息
			return v.Error()
		} else {
			return v.String()
		}
	case error:
		return v.Error()
	case NormalErrorInterface:
		return v.Error()
	default:
		return ""
	}
}
