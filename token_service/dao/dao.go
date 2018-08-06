package dao

var DB *Dao

type Dao struct {
	redis         *RedisCli
	commonRedis   *RedisCliCommon
	mysql         *Mysql
	mysql2        *Mysql2
	commonMysql   *MysqlCommon
	currencyMysql *MysqlCurrency
}

func NewDao() (dao *Dao) {
	mysql := NewMysql()
	rediscli := NewRedisCli()
	dao = &Dao{
		redis:         rediscli,
		commonRedis:   NewRedisCliCommon(),
		mysql:         mysql,
		mysql2:        NewMysql2(),
		commonMysql:   NewMysqlCommon(),
		currencyMysql: NewMysqlCurrency(),
	}
	return
}

func InitDao() {
	DB = NewDao()
}
