package main

import (
	cf "digicon/token_service/conf"
	"digicon/token_service/dao"
	. "digicon/token_service/log"
	"digicon/token_service/model"
	"digicon/token_service/rpc"
	"flag"
	"github.com/liudng/godump"
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
	g := model.GetKLine("bchbtc", "5min", 150)
	godump.Dump(g)
	sig := <-quitChan
	Log.Infof("server close by sig %s", sig.String())
}
