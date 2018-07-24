package main

import (
	cf "digicon/kline_service/conf"
	"digicon/kline_service/dao"
	log "github.com/sirupsen/logrus"
	"digicon/kline_service/rpc"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"digicon/common/xlog"
)

func init()  {
	cf.Init()
	path := cf.Cfg.MustValue("log", "log_dir")
	name := cf.Cfg.MustValue("log", "log_name")
	level := cf.Cfg.MustValue("log", "log_level")
	xlog.InitLogger(path,name,level)
}

func main() {
	flag.Parse()
	cf.Init()
	//InitLog()
	log.Infof("begin run server")
	dao.InitDao()
	go rpc.RPCServerInit()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())
}
