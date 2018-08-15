package utils

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	cf "digicon/wallet_service/conf"
)

var Engine_wallet *xorm.Engine
var Engine_token *xorm.Engine
var Engine_common *xorm.Engine

var EngineUserCurrency *xorm.Engine

var Redis *redis.Client

func Init() {
	var err error

	//mysql初始化
	dsource := cf.Cfg.MustValue("mysql", "wallet_conn")
	Engine_wallet, err = xorm.NewEngine("mysql", dsource)
	if err != nil {
		fmt.Println("mysql connect error!.")
		panic(err)
	}
	Engine_wallet.ShowSQL(false)
	err = Engine_wallet.Ping()
	if err != nil {
		panic(err)
	}

	dsource = cf.Cfg.MustValue("mysql", "token_conn")
	Engine_token, err = xorm.NewEngine("mysql", dsource)

	if err != nil {
		panic(err)
	}
	Engine_token.ShowSQL(false)
	err = Engine_token.Ping()
	if err != nil {
		panic(err)
	}

	dsource = cf.Cfg.MustValue("mysql", "common_conn")
	Engine_common, err = xorm.NewEngine("mysql", dsource)

	if err != nil {
		panic(err)
	}
	Engine_common.ShowSQL(false)
	err = Engine_common.Ping()
	if err != nil {
		panic(err)
	}

	// ------------- user currency ----------------
	dsource = cf.Cfg.MustValue("mysql", "currency_conn")
	EngineUserCurrency, err = xorm.NewEngine("mysql", dsource)
	if err != nil {
		fmt.Println("connect db currency error!")
		panic(err)
	}
	EngineUserCurrency.ShowSQL(false)
	err = EngineUserCurrency.Ping()
	if err != nil {
		panic(err)
	}

	////

	//redis初始化
	//Redis = nil

	addr := cf.Cfg.MustValue("redis", "addr")
	pass := cf.Cfg.MustValue("redis", "pass")
	num := cf.Cfg.MustInt("redis", "num")
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass, // no password set
		DB:       num,  // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("redis connect faild ")
	}
	fmt.Println(pong)
	_, err = client.ZRangeWithScores("token:1/2", 0, 1).Result()
	if err != nil {
		log.Fatalf("redis connect faild ")
	}

	Redis = client
}
