package main

import (
	"digicon/common/snowflake"
	"digicon/common/xlog"
	cf "digicon/currency_service/conf"
	"digicon/currency_service/cron"
	"digicon/currency_service/dao"
	"digicon/currency_service/rpc"
	"digicon/currency_service/rpc/client"
	"digicon/currency_service/rpc/handler"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	cf.Init()
	path := cf.Cfg.MustValue("log", "log_dir")
	name := cf.Cfg.MustValue("log", "log_name")
	level := cf.Cfg.MustValue("log", "log_level")
	xlog.InitLogger(path, name, level)
}

func main() {
	flag.Parse()
	log.Infof("begin run server")
	snowflake.Init()
	dao.InitDao()

	go rpc.RPCServerInit()
	go handler.InitCheckOrderStatus()

	client.InitInnerService()

	cron.InitCron()

	// 定时任务统计
	go cron.DailyStart()
	//go new(cron.DailyCountSheet).Run()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())
}
