package dao

var globalDb Dbs

type Dbs struct {
	sql   *Mysql
	redis *RedisCli
}

func NewDbs(dbs *Dbs) {
	sql := NewMysql()
	redis := NewRedisCli()
	dbs = &Dbs{
		sql:   sql,
		redis: redis,
	}
	return
}

func init() {
	NewDbs(&globalDb)
}
