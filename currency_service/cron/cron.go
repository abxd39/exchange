package cron

import (
	"digicon/common/app"
	cf "digicon/currency_service/conf"
	"github.com/robfig/cron"
)

var CronInstance *cron.Cron

func InitCron() {
	//划入
	app.AsyncTask(HandlerTransferFromToken, true)

	//划出
	app.AsyncTask(HandlerTransferToTokenDone, true)

	// 手动执行日汇总
	//app.AsyncTask(func() { new(DailyCountSheet).RunByDays(1535644800) }, false)


	go new(DailyCountSheet).Run()
	//go new(DailyCountSheet).RunByDays(1535644800)
	//cron

	// 定时任务

	if cf.Cfg.MustBool("cron", "run", false) {
		CronInstance = cron.New()
		CronInstance.AddJob("0 0 1 * * *", DailyCountSheet{})          // 日汇总
		CronInstance.AddFunc("0 0 */1 * * *", CheckAdsAutoDownline)    // ads auto downline
		CronInstance.AddFunc("0 30 * * * *", ResendTransferToTokenMsg) // 划转到币币消息重发机制
		CronInstance.Start()
	}

}
