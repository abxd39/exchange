package main

import (
	cf "digicon/price_service/conf"
	"digicon/price_service/dao"
	. "digicon/price_service/log"
	"digicon/price_service/model"
	"digicon/price_service/rpc"
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
	model.GetQueneMgr().Init()
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	Log.Infof("server close by sig %s", sig.String())
}
