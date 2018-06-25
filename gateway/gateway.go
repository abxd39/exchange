package main

import (
	cf "digicon/gateway/conf"
	"digicon/gateway/http"
	. "digicon/gateway/log"
	"digicon/gateway/rpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cf.Init()
	InitLog()
	go http.InitHttpServer()
	go rpc.InitInnerService()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	sig := <-quitChan
	Log.Infof("server close by sig %s", sig.String())
}
