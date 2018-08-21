package main

import (
	"digicon/common/snowflake"
	"digicon/common/xlog"
	cf "digicon/currency_service/conf"
	"digicon/currency_service/cron"
	"digicon/currency_service/rpc"
	"digicon/currency_service/rpc/client"
	"digicon/currency_service/rpc/handler"
	"digicon/currency_service/dao"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"digicon/currency_service/model"
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

	cron.InitCron()


	cron.CheckAdsAutoDownline()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())
}




