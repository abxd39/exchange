package main

import (
	cf "digicon/user_service/conf"
	"digicon/user_service/dao"
	log "github.com/sirupsen/logrus"
	"digicon/user_service/rpc"
	"digicon/user_service/rpc/client"
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

	log.Infof("begin run server")
	dao.InitDao()
	go rpc.RPCServerInit()
	client.InitInnerService()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())
}
