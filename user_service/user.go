package main

import (
	cf "digicon/user_service/conf"
	"digicon/user_service/dao"
	. "digicon/user_service/log"
	"digicon/user_service/rpc"
	"digicon/user_service/rpc/client"
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
	client.InitInnerService()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	//	model.SendSms("18664328365", "86",1)

	sig := <-quitChan
	Log.Infof("server close by sig %s", sig.String())
}
