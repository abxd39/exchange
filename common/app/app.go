package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"runtime"
)

// 标记app是否已退出
var IsAppExit = false

// 启动一个goroutine
func NewGoroutine(f func()) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				errorInfo := "协程panic，故障堆栈："
				for i := 1; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					} else {
						errorInfo += "\n"
					}
					errorInfo += fmt.Sprintf("%v %v", file, line)
				}

				log.Errorf(errorInfo)
			}
		}()

		f()
	}()
}
