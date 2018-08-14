package main

import (
	"digicon/wallet_service/rpc"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	//cf "digicon/currency_service/conf"
	"digicon/common/xlog"
	cf "digicon/wallet_service/utils"
	"digicon/wallet_service/rpc/watch"
	"digicon/wallet_service/rpc/client"
)

func init() {
	//cf.Init()
	path := cf.Cfg.MustValue("log", "log_dir")
	name := cf.Cfg.MustValue("log", "log_name")
	level := cf.Cfg.MustValue("log", "log_level")
	xlog.InitLogger(path, name, level)
}

func main() {
	flag.Parse()

	//比特币充币提币监控
	go watch.StartBtcWatch()
	//以太币、ERC20代币提币检查
	go watch.StartEthCheckNew()
	//以太币、ERC20代币充币检查
	go watch.StartEthCBiWatch()

	go rpc.RPCServerInit()
	go client.InitInnerService()
	//new(client.Watch).Start("https://rinkeby.infura.io/mew")  // need ...
	//go new(client.BTCWatch).Start()
	//go new(client.BTCWatch).Start()
	//go new(client.Watch).Start()
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
