package dao

var DB *Dao

type Dao struct {
	redis       *RedisCli
	mysql       *Mysql
	mysql2      *Mysql2
	commonMysql *MysqlCommon
}

func NewDao() (dao *Dao) {
	mysql := NewMysql()
	rediscli := NewRedisCli()
	dao = &Dao{
		redis:       rediscli,
		mysql:       mysql,
		mysql2:      NewMysql2(),
		commonMysql: NewMysqlCommon(),
	}
	return
}

func InitDao() {
	DB = NewDao()
}
