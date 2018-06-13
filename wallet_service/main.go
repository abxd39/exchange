package main

import (
	"digicon/wallet_service/rpc"
	"flag"
	"os"
	"os/signal"
	"syscall"
	. "digicon/wallet_service/utils"
)

func main() {
	flag.Parse()

	go rpc.RPCServerInit()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan

	Log.Infof("server close by sig %s", sig.String())
}
