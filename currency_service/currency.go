package main

import (
	"digicon/common/app"
	"digicon/common/snowflake"
	"digicon/common/xlog"
	cf "digicon/currency_service/conf"
	"digicon/currency_service/cron"
	"digicon/currency_service/dao"
	"digicon/currency_service/model"
	"digicon/currency_service/rpc"
	"digicon/currency_service/rpc/client"
	"digicon/currency_service/rpc/handler"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	cf.Init()
	path := cf.Cfg.MustValue("log", "log_dir")
	name := cf.Cfg.MustValue("log", "log_name")
	level := cf.Cfg.MustValue("log", "log_level")
	xlog.InitLogger(path, name, level)
	dao.InitDao()
	model.InitCnyPrice()
}

func main() {
	flag.Parse()
	log.Infof("begin run server")
	snowflake.Init()

	go rpc.RPCServerInit()
	go handler.InitCheckOrderStatus()

	client.InitInnerService()

	// 定时任务
	cron.InitCron()

	// 监听退出
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		//syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())

	// 标记程序退出
	app.IsAppExit = true

	// 关闭定时任务调度，正在运行的不影响
	if cron.CronInstance != nil {
		cron.CronInstance.Stop()
	}

	// 不立刻杀死进程
	time.Sleep(3 * 60 * time.Second)
}
