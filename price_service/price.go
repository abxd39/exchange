package main

import (
	cf "digicon/price_service/conf"
	"digicon/price_service/dao"

	"digicon/common/xlog"
	"digicon/price_service/model"
	"digicon/price_service/rpc"
	"digicon/price_service/rpc/client"
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
	model.Test3()
	return
	dao.InitDao()


	go rpc.RPCServerInit()
	client.InitInnerService()
	model.GetQueneMgr().Init()
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())
}
