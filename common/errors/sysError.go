package errors

import (
	"fmt"

	. "digicon/proto/common"
)

// 系统错误接口
type SysErrorInterface interface {
	Error() string
	String() string
	Status() int
}

// 系统错误实现
type SysError struct {
	status    int
	simpleMsg string
	fullMsg   string
}

func (e *SysError) Error() string {
	return e.simpleMsg
}

func (e *SysError) String() string {
	return e.fullMsg
}

func (e *SysError) Status() int {
	return e.status
}

// 创建系统错误
func NewSys(options ...interface{}) error {
	var (
		status    int
		simpleMsg string
		fullMsg   string
	)

	for _, v := range options {
		switch opt := v.(type) {
		default:
			status = ERRCODE_UNKNOWN
			simpleMsg = "系统错误"
			fullMsg = fmt.Sprintf("系统错误: %v", opt)
		}
	}

	return &SysError{
		status:    status,
		simpleMsg: simpleMsg,
		fullMsg:   fullMsg,
	}
}
