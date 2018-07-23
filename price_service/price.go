package main

import (
	cf "digicon/price_service/conf"
	"digicon/price_service/dao"
	. "digicon/price_service/log"
	"digicon/price_service/model"
	"digicon/price_service/rpc"
	"digicon/price_service/rpc/client"
	"flag"
	"os"
	"os/signal"
	"syscall"
	log "github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	cf.Init()
	InitLogger()
	log.Infof("begin run server")

	dao.InitDao()

	go rpc.RPCServerInit()
	client.InitInnerService()
	//exchange.LoadCacheQuene()
	model.GetQueneMgr().Init()
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())
}
