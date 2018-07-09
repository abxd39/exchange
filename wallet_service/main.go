package main

import (
	"digicon/wallet_service/rpc"
	"digicon/wallet_service/rpc/client"
	. "digicon/wallet_service/utils"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()

	go rpc.RPCServerInit()
	//new(client.Watch).Start("https://rinkeby.infura.io/mew")  // need ...
	go new(client.BTCWatch).Start("http://bitcoin:bitcoin@localhost:18332/")
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
