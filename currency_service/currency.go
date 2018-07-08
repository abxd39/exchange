package main

import (
	cf "digicon/currency_service/conf"
	"digicon/currency_service/dao"
	. "digicon/currency_service/log"
	"digicon/currency_service/rpc"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"digicon/currency_service/rpc/client"
)

func main() {
	flag.Parse()
	cf.Init()
	InitLog()
	Log.Infof("begin run server")
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
	Log.Infof("server close by sig %s", sig.String())
}
