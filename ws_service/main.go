package main

import (
	cf "digicon/ws_service/conf"
	"digicon/ws_service/dao"
	log "github.com/sirupsen/logrus"
	//"digicon/ws_service/rpc"
	//"digicon/ws_service/rpc/client"
	"digicon/ws_service/http"
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
	go http.InitHttpServer()
	//go rpc.RPCServerInit()

	//client.InitInnerService()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())
}
