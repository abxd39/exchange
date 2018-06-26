package dao

import (
	"digicon/backstage_service/conf"
	"digicon/backstage_service/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Mysql struct {
	engi *xorm.Engine
}

func NewMysql() (mysql *Mysql) {
	dsource := conf.Cfg.MustValue("mysql", "conn")
	//root:current@tcp(47.106.136.96:3306)/rumi?charset=utf8
	engine, err := xorm.NewEngine("mysql", dsource)
	if err != nil {
		log.Log.Fatalf("db err is %s", err)
	}

	engine.ShowSQL(true)
	//cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Hour, 1000)
	//engine.SetDefaultCacher(cacher)

	err = engine.Ping()
	if err != nil {
		log.Log.Fatalf("db err is %s", err)
	}

	mysql = &Mysql{
		engi: engine,
	}
	return mysql
}

func (d *Dbs) GetMysqlInstance() *xorm.Engine {
	return d.sql.engi
}
