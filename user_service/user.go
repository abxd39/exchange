package main

import (
	"digicon/common/xlog"
	cf "digicon/user_service/conf"
	"digicon/user_service/cron"
	"digicon/user_service/dao"
	"digicon/user_service/rpc"
	"digicon/user_service/rpc/client"
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
	dao.InitDao()
	go rpc.RPCServerInit()
	client.InitInnerService()

	// 定时脚本
	cron.InitCron()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())
}
