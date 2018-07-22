package main

import (
	cf "digicon/config_service/conf"
	"digicon/config_service/dao"
	. "digicon/config_service/log"
	//"digicon/config_service/rpc/client"
	"digicon/config_service/http"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()

	cf.Init()
	InitLog()
	Log.Infof("begin run server")

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
	Log.Infof("server close by sig %s", sig.String())
}
