package cron

import (
	"github.com/robfig/cron"
)

func InitCron() {
	c := cron.New()

	//...
	//c.AddFunc("0 0 3 * * *", RegisterNoReward)

	c.Start()
}
