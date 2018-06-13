package utils
import (
_ "github.com/go-sql-driver/mysql"
"github.com/go-xorm/xorm"
"github.com/garyburd/redigo/redis"
)


var Engine *xorm.Engine
var Redis *redis.Conn

func init() {
	var err error
	println("dao 初始化")

	//mysql初始化
	dsource := Cfg.MustValue("mysql", "conn")
	Engine, err = xorm.NewEngine("mysql", dsource)


	if err != nil{
		panic( err)
	}
	err = Engine.Ping()
	if err != nil{
		panic(err)
	}
	//redis初始化
	Redis = nil

}
