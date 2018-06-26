package main

import (
	"digicon/backstage_service/log"
	"digicon/backstage_service/rpc"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	log.Log.Infof("begin run backstage servce")
	go rpc.RPCServerInit()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Log.Infof("backstage server close by sig %s", sig.String())
}
