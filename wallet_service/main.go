package main

import (
	. "digicon/wallet_service/utils"
	"flag"
	"os"
	"os/signal"
	"syscall"
	//"digicon/wallet_service/rpc/client"
	"digicon/wallet_service/rpc"
)

func main() {
	flag.Parse()

	go rpc.RPCServerInit()
	 //new(client.Watch).Start("https://rinkeby.infura.io/mew")
	//return
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan

	Log.Infof("server close by sig %s", sig.String())
}
