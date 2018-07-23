package dao

import (
	"digicon/currency_service/conf"
	. "digicon/currency_service/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Mysql struct {
	im *xorm.Engine
}

type MysqlToken struct {
	tim *xorm.Engine
}


func NewMysql() (mysql *Mysql) {
	dsource := conf.Cfg.MustValue("mysql", "conn")

	engine, err := xorm.NewEngine("mysql", dsource)
	if err != nil {
		Log.Fatalf("db err is %s", err)
	}
	engine.ShowSQL(true)
	//cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Hour, 1000)
	//engine.SetDefaultCacher(cacher)

	err = engine.Ping()
	if err != nil {
		Log.Fatalf("db err is %s", err)
	}
	mysql = &Mysql{
		im: engine,
	}
	return mysql
}


func (s *Dao) GetMysqlConn() *xorm.Engine {
	return s.mysql.im
}


func NewTokenMysql() (tkmysql *MysqlToken) {
	dsource := conf.Cfg.MustValue("mysql", "token_conn")

	engine, err := xorm.NewEngine("mysql", dsource)
	if err != nil {
		Log.Fatalf("db err is %s", err)
	}
	engine.ShowSQL(true)
	err = engine.Ping()
	if err != nil {
		Log.Fatalf("db err is %s", err)
	}
	tkmysql = &MysqlToken{
		tim: engine,
	}
	return tkmysql
}

func (s *Dao) GetTokenMysqlConn() *xorm.Engine {
	return  s.tokenMysql.tim
}