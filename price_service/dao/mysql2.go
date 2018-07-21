package dao

import (
	"digicon/price_service/conf"
	. "digicon/price_service/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"time"
)

type Mysql2 struct {
	im *xorm.Engine
}

func NewMysql2() (mysql *Mysql2) {
	dsource := conf.Cfg.MustValue("mysql", "conn2")

	engine, err := xorm.NewEngine("mysql", dsource)
	if err != nil {
		Log.Fatalf("db err is %s", err)
	}
	engine.ShowSQL(false)
	cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Hour, 1000)
	engine.SetDefaultCacher(cacher)

	err = engine.Ping()
	if err != nil {
		Log.Fatalf("db err is %s", err)
	}

	mysql = &Mysql2{
		im: engine,
	}
	return mysql
}

func (s *Dao) GetMysqlConn2() *xorm.Engine {
	return s.mysql2.im
}
