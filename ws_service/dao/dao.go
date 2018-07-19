package dao

var DB *Dao

type Dao struct {
	redis *RedisCli
	mysql *Mysql
}


func NewDao()(dao *Dao) {
	mysql :=NewMysql()
	rediscli := NewRedisCli()
	dao = &Dao{
		redis: rediscli,
		mysql: mysql,
	}
	return
}

func InitDao(){
	DB = NewDao()
}