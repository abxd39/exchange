package main

import (
	cf "digicon/public_service/conf"
	"digicon/public_service/dao"
	//plog "digicon/public_service/log"
	"digicon/public_service/rpc"
	"flag"
	"os"
	"os/signal"
	"syscall"
	plog "github.com/sirupsen/logrus"
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
	//cf.Init()
	//plog.InitLog()
	plog.Infof("begin run server")
	dao.InitDao()
	go rpc.RPCServerInit()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	plog.Infof("server close by sig %s", sig.String())
}
