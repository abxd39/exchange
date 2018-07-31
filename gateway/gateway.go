package main

import (
	"digicon/common/xlog"
	cf "digicon/gateway/conf"
	"digicon/gateway/http"
	"digicon/gateway/rpc"
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
