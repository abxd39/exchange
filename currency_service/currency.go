package main

import (
	cf "digicon/currency_service/conf"
	"digicon/currency_service/dao"
	//log "github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"digicon/currency_service/rpc"
	"digicon/currency_service/rpc/client"
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
	//xlog.Initlog.er(path,name,level)
}

func main() {
	flag.Parse()
	cf.Init()
	//Initlog.)
	log.Infof("begin run server")
	dao.InitDao()
	go rpc.RPCServerInit()

	client.InitInnerService()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())
}
