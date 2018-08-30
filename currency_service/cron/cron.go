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

	//
	go new(DailyCountSheet).Run()

	//cron
	if cf.Cfg.MustBool("cron", "run", false) {
		c := cron.New()

		//AddFunc
		spec := "0 0 1 * * *"    // every day ...
		specTwo := "0 0 4 * * *" // every day ...

		// 定时任务统计
		c.AddJob(spec, DailyCountSheet{})
		c.AddJob(specTwo, DailyCountSheet{})

		c.AddFunc("0 30 * * * *", ResendTransferToTokenMsg)

		// ads auto downline. every one hour check
		c.AddFunc("0 0 */1 * * *", CheckAdsAutoDownline)

		c.Start()
	}

}
