package dao

import (
	"digicon/kline_service/conf"
	log "github.com/sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"time"
)

type Mysql struct {
	im *xorm.Engine
}

func NewMysql() (mysql *Mysql) {
	dsource := conf.Cfg.MustValue("mysql", "conn")

	engine, err := xorm.NewEngine("mysql", dsource)
	if err != nil {
		log.Fatalf("db err is %s", err)
	}
	engine.ShowSQL(true)
	cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Hour, 1000)
	engine.SetDefaultCacher(cacher)

	err = engine.Ping()
	if err != nil {
		log.Fatalf("db err is %s", err)
	}

	mysql = &Mysql{
		im: engine,
	}
	return mysql
}
