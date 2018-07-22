package dao

import (
	"digicon/config_service/conf"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
)






type Mysql struct {
	//im *xorm.Engine
	Engine_wallet *xorm.Engine
    Engine_token *xorm.Engine
	Engine_common *xorm.Engine
}






func NewMysql()(mysql *Mysql){
	var err error
	//
	//var Engine_wallet *xorm.Engine
	//var Engine_token *xorm.Engine
	//var Engine_common *xorm.Engine
	//mysql初始化
	dsource := conf.Cfg.MustValue("mysql", "wallet_conn")
	Engine_wallet, err := xorm.NewEngine("mysql", dsource)
	if err != nil {
		panic(err)
	}
	err = Engine_wallet.Ping()
	if err != nil {
		panic(err)
	}

	dsource = conf.Cfg.MustValue("mysql", "token_conn")
	Engine_token, err := xorm.NewEngine("mysql", dsource)

	if err != nil {
		panic(err)
	}
	err = Engine_token.Ping()
	if err != nil {
		panic(err)
	}

	dsource = conf.Cfg.MustValue("mysql", "common_conn")
	Engine_common, err := xorm.NewEngine("mysql", dsource)

	if err != nil {
		panic(err)
	}
	err = Engine_token.Ping()
	if err != nil {
		panic(err)
	}


	mysql = &Mysql{
		//im: engine,
		Engine_common:Engine_common,
		Engine_wallet:Engine_common,
		Engine_token: Engine_token,
	}
	return mysql
}

func (s *Dao) GetMysqlTokenConn() *xorm.Engine {
	return s.mysql.Engine_token
}

func (s *Dao) GetMysqlWalletConn() *xorm.Engine{
	return s.mysql.Engine_wallet
}

func (s *Dao) GetMysqlCommonConn() *xorm.Engine{
	return s.mysql.Engine_common
}





