package main

import (
	cf "digicon/token_service/conf"
	"digicon/token_service/dao"
	. "digicon/token_service/log"
	"digicon/token_service/model"
	"digicon/token_service/rpc"
	"digicon/token_service/rpc/client"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"digicon/token_service/exchange"
)

func main() {
	flag.Parse()
	cf.Init()
	InitLog()
	Log.Infof("begin run server")
	dao.InitDao()
	go rpc.RPCServerInit()
	client.InitInnerService()
	model.GetQueneMgr().Init()
	//model.GetKLine("BTC/USDT","1min",10)

	go exchange.InitExchange()
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	Log.Infof("server close by sig %s", sig.String())
}
