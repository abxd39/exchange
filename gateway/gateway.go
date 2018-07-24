package main

import (
	"digicon/gateway/http"
	log "github.com/sirupsen/logrus"
	"digicon/gateway/rpc"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"digicon/common/xlog"
	cf "digicon/gateway/conf"
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
	//InitLog()
	go http.InitHttpServer()
	go rpc.InitInnerService()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())
}
