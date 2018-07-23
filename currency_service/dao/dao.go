package dao

var DB *Dao

type Dao struct {
	redis *RedisCli
	mysql *Mysql
	tokenMysql *MysqlToken
}

func NewDao() (dao *Dao) {
	mysql := NewMysql()
	tkmysql := NewTokenMysql()
	rediscli := NewRedisCli()

	dao = &Dao{
		redis: rediscli,
		mysql: mysql,
		tokenMysql: tkmysql,
	}
	return
}

func InitDao() {
	DB = NewDao()
}
