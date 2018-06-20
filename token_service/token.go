package main

import (
	cf "digicon/token_service/conf"
	"digicon/token_service/dao"
	. "digicon/token_service/log"
	"digicon/token_service/model"
	"digicon/token_service/rpc"
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
