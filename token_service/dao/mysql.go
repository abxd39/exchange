package dao

import (
	"digicon/token_service/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
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
	/*
	cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Hour, 1000)

	engine.SetDefaultCacher(cacher)
*/
	err = engine.Ping()
	if err != nil {
		log.Fatalf("db err is %s", err)
	}

	mysql = &Mysql{
		im: engine,
	}
	return mysql
}

func (s *Dao) GetMysqlConn() *xorm.Engine {
	return s.mysql.im
}

// g_common库
type MysqlCommon struct {
	im *xorm.Engine
}

func NewMysqlCommon() (mysql *MysqlCommon) {
	dsource := conf.Cfg.MustValue("mysql", "common")

	engine, err := xorm.NewEngine("mysql", dsource)
	if err != nil {
		log.Fatalf("db err is %s", err)
	}
	engine.ShowSQL(true)
	//cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Hour, 1000)
	//engine.SetDefaultCacher(cacher)

	err = engine.Ping()
	if err != nil {
		log.Fatalf("db err is %s", err)
	}

	mysql = &MysqlCommon{
		im: engine,
	}
	return mysql
}

func (s *Dao) GetCommonMysqlConn() *xorm.Engine {
	return s.commonMysql.im
}
