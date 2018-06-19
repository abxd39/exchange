package main

import (
	cf "digicon/kline_service/conf"
	"digicon/kline_service/dao"
	. "digicon/kline_service/log"
	"digicon/kline_service/rpc"
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
	go rpc.RPCServerInit()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	Log.Infof("server close by sig %s", sig.String())
}
