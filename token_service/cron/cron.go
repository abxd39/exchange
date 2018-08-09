package cron

import (
	cf "digicon/currency_service/conf"
	"github.com/robfig/cron"
)

func InitCron() {
	//划入
	go HandlerTransferFromCurrency()

	//划出
	go HandlerTransferToCurrencyDone()

	//cron
	if cf.Cfg.MustBool("cron", "run", false) {
		c := cron.New()
		c.AddFunc("0 30 * * * *", ResendTransferToCurrencyMsg)
		c.Start()
	}
}
