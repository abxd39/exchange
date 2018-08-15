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
	cf "digicon/wallet_service/conf"
	"digicon/wallet_service/rpc/client"
	"fmt"
	"digicon/wallet_service/utils"
	"digicon/wallet_service/cron"
	"digicon/wallet_service/rpc/watch"
)

func init() {
	cf.Init()
	fmt.Println("log .....")
	path := cf.Cfg.MustValue("log", "log_dir")
	name := cf.Cfg.MustValue("log", "log_name")
	level := cf.Cfg.MustValue("log", "log_level")
	xlog.InitLogger(path, name, level)
	fmt.Println("log start ...")
	utils.Init()
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

	///////////////////
	//  统计每天的币数
	go cron.DailyStart()
	//go new(cron.WalletDailyCountSheet).Run()

	/////////////////////


	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan

	log.Infof("server close by sig %s", sig.String())
}
