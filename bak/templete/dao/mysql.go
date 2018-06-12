package dao

import (
	"digicon/wallet_service/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	"time"
)

type Mysql struct {
	im *xorm.Engine
}

func NewMysql() (mysql *Mysql) {
	dsource := conf.Cfg.MustValue("mysql", "conn")

	//root:current@tcp(47.106.136.96:3306)/rumi?charset=utf8
	engine, err := xorm.NewEngine("mysql", dsource)
	if err != nil {
		glog.Fatalf("db err is %s", err)
	}

	engine.ShowSQL(true)
	cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Hour, 1000)
	engine.SetDefaultCacher(cacher)

	err = engine.Ping()
	if err != nil {
		glog.Fatalf("db err is %s", err)
	}

	mysql = &Mysql{
		im: engine,
	}
	return mysql
}
