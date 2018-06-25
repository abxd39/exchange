package main

import (
	cf "digicon/user_service/conf"
	"digicon/user_service/dao"
	. "digicon/user_service/log"
	"digicon/user_service/rpc"
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
	)
	//model.SendSms("57002661", 10)

	sig := <-quitChan
	Log.Infof("server close by sig %s", sig.String())
}
