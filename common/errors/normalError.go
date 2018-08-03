package errors

import (
	. "digicon/proto/common"
	"fmt"
)

// 普通错误接口
type NormalErrorInterface interface {
	Error() string
	Status() int
}

// 普通错误实现
type NormalError struct {
	status int
	msg    string
}

func (e *NormalError) Error() string {
	return e.msg
}

func (e *NormalError) Status() int {
	return e.status
}

// 创建普通错误
func NewNormal(options ...interface{}) error {
	var (
		status int
		msg    string
	)

	for _, v := range options {
		switch opt := v.(type) {
		case int:
			status = opt
		default:
			msg = fmt.Sprintf("%v", opt)
		}
	}

	// 未指定错误码，使用默认错误码
	if status == 0 {
		status = ERRCODE_NORMAL_ERROR
	}

	return &NormalError{
		status: status,
		msg:    msg,
	}
}
