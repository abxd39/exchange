package cron

import (
	"digicon/common/app"
	cf "digicon/currency_service/conf"
	"github.com/robfig/cron"
)

var CronInstance *cron.Cron

func InitCron() {
	// 划入
	app.NewGoroutine(HandlerTransferFromToken)

	// 划出
	app.NewGoroutine(HandlerTransferToTokenDone)

	// 手动执行日汇总
	//app.NewGoroutine(func() { new(DailyCountSheet).RunByDays(1535644800) })

	// 定时任务
	if cf.Cfg.MustBool("cron", "run", false) {
		CronInstance = cron.New()
		CronInstance.AddJob("0 0 1 * * *", DailyCountSheet{})            // 凌晨1点
		CronInstance.AddFunc("0 0 */1 * * *", CheckAdsAutoDownline)      // 每一小时
		CronInstance.AddFunc("0 */30 * * * *", ResendTransferToTokenMsg) // 每半小时
		CronInstance.Start()
	}

}
