package main

import (
	"digicon/wallet_service/rpc"
	"digicon/wallet_service/rpc/client"
	"flag"
	"os"
	"os/signal"
	"syscall"
	log "github.com/sirupsen/logrus"
	//cf "digicon/currency_service/conf"
	cf "digicon/wallet_service/utils"
	"digicon/common/xlog"
)

func init()  {
	//cf.Init()
	path := cf.Cfg.MustValue("log", "log_dir")
	name := cf.Cfg.MustValue("log", "log_name")
	level := cf.Cfg.MustValue("log", "log_level")
	xlog.InitLogger(path,name,level)
}



func main() {
	flag.Parse()

	go rpc.RPCServerInit()
	//new(client.Watch).Start("https://rinkeby.infura.io/mew")  // need ...
	//go new(client.BTCWatch).Start()
	//go new(client.BTCWatch).Start()
	go new(client.Watch).Start()
	//return
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan

	log.Infof("server close by sig %s", sig.String())
}
