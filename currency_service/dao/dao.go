package dao

var DB *Dao

type Dao struct {
	redis       *RedisCli
	commonRedis *RedisCliCommon
	mysql       *Mysql
	tokenMysql  *MysqlToken
	commonMysql *MysqlCommon
}

func NewDao() (dao *Dao) {
	mysql := NewMysql()
	tkmysql := NewTokenMysql()
	comMysql := NewCommonMysql()
	rediscli := NewRedisCli()

	dao = &Dao{
		redis:       rediscli,
		commonRedis: NewRedisCliCommon(),
		mysql:       mysql,
		tokenMysql:  tkmysql,
		commonMysql: comMysql,
	}
	return
}

func InitDao() {
	DB = NewDao()
}
