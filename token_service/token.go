package main

import (
	"digicon/common/app"
	"digicon/common/snowflake"
	"digicon/common/xlog"
	cf "digicon/token_service/conf"
	"digicon/token_service/cron"
	"digicon/token_service/dao"
	"digicon/token_service/model"
	"digicon/token_service/rpc"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	cf.Init()
	path := cf.Cfg.MustValue("log", "log_dir")
	name := cf.Cfg.MustValue("log", "log_name")
	level := cf.Cfg.MustValue("log", "log_level")
	xlog.InitLogger(path, name, level)
}

func main() {
	log.Infof("begin run server")

	flag.Parse()

	// 初始化snowflake
	snowflake.Init()

	// 初始化dao
	dao.InitDao()

	//model.Test2(1,1000)
	//model.Test3(1533139200,1533225600)
	//model.Testg()
	//a:=[5]int{100001, 100002, 100003}
	//model.GetAllBalanceCny(a)
	go rpc.RPCServerInit()

	model.GetQueneMgr().Init()
	//model.Test9(1535299200,1535385600)
	model.Test10()
	//model.Testu()
	fmt.Println("model get ...")
	//model.GetKLine("BTC/USDT","1min",10)
	//model.Test()
	//go exchange.InitExchange()

	// 定时任务
	cron.InitCron()

	// 监听退出
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		//syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Infof("server close by sig %s", sig.String())

	// 标记程序已退出
	app.IsAppExit = true

	// 关闭定时任务调度，正在运行的不影响
	if cron.CronInstance != nil {
		cron.CronInstance.Stop()
	}

	// 不立刻杀死进程
	time.Sleep(3 * 60 * time.Second)
}
