package main

import (
	cf "digicon/wallet_service/conf"
	"digicon/wallet_service/http"
	"digicon/wallet_service/rpc"
	"github.com/golang/glog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cf.InitConf()
	go http.InitHttpServer()
	go rpc.InitInnerService()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	glog.Infof("server close by sig %s", sig.String())
}
