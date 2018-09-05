package app

import (
	log "github.com/sirupsen/logrus"
)

// 标记app是否已退出
var IsAppExit = false

// 启动一个异步执行任务
func AsyncTask(f func(), recoverPullAgain bool) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Error("协程panic，err：", e)
			}

			// 判断是否重新拉起
			if recoverPullAgain {
				AsyncTask(f, recoverPullAgain)
			}
		}()

		f()
	}()
}
