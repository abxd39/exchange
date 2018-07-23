package main

import (
	cf "digicon/ws_service/conf"
	"digicon/ws_service/dao"
	. "digicon/ws_service/log"
	//"digicon/ws_service/rpc"
	//"digicon/ws_service/rpc/client"
	"digicon/ws_service/http"
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
