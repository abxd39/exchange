package dao

var DB *Dao
//var CDao *ConsuleDao

type Dao struct {
	redis *RedisCli
	mysql *Mysql
	consul *ConsulCli
}



func NewDao()(dao *Dao) {
	mysql :=NewMysql()
	rediscli := NewRedisCli()
	consulcli := NewConsulCli()
	dao = &Dao{
		redis: rediscli,
		mysql: mysql,
		consul: consulcli,
	}
	return
}


func InitDao(){
	DB = NewDao()
}