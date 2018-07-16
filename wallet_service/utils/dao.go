package utils

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var Engine_wallet *xorm.Engine
var Engine_token *xorm.Engine
var Engine_common *xorm.Engine

var EngineUserCurrency *xorm.Engine

var Redis *redis.Conn

func init() {
	var err error

	//mysql初始化
	dsource := Cfg.MustValue("mysql", "wallet_conn")
	Engine_wallet, err = xorm.NewEngine("mysql", dsource)
	if err != nil {
		panic(err)
	}
	err = Engine_wallet.Ping()
	if err != nil {
		panic(err)
	}

	dsource = Cfg.MustValue("mysql", "token_conn")
	Engine_token, err = xorm.NewEngine("mysql", dsource)

	if err != nil {
		panic(err)
	}
	err = Engine_token.Ping()
	if err != nil {
		panic(err)
	}

	dsource = Cfg.MustValue("mysql", "common_conn")
	Engine_common, err = xorm.NewEngine("mysql", dsource)

	if err != nil {
		panic(err)
	}
	err = Engine_token.Ping()
	if err != nil {
		panic(err)
	}

	// ------------- user currency ----------------
	dsource = Cfg.MustValue("mysql", "currency_conn")
	EngineUserCurrency, err = xorm.NewEngine("mysql", dsource)
	if err != nil {
		fmt.Println("connect db currency error!")
		panic(err)
	}
	err = EngineUserCurrency.Ping()
	if err != nil {
		panic(err)
	}

	////

	//redis初始化
	Redis = nil

}
