package cron

import (
	cf "digicon/currency_service/conf"
	"github.com/robfig/cron"
)

func InitCron() {
	//划入
	go HandlerTransferFromToken()

	//划出
	go HandlerTransferToTokenDone()

	//cron
	if cf.Cfg.MustBool("cron", "run", false) {
		c := cron.New()
		c.AddFunc("0 30 * * * *", ResendTransferToTokenMsg)
		c.Start()
	}


	// 定时任务统计
	go DailyStart()
}
