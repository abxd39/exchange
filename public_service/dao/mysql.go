package dao

import (
	. "digicon/proto/common"
	"digicon/public_service/conf"
	. "digicon/public_service/log"
	"digicon/user_service/model"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Mysql struct {
	im *xorm.Engine
}

func NewMysql() (mysql *Mysql) {
	dsource := conf.Cfg.MustValue("mysql", "conn")

	engine, err := xorm.NewEngine("mysql", dsource)
	if err != nil {
		Log.Fatalf("db err is %s", err)
	}
	engine.ShowSQL(true)
	cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Hour, 1000)
	engine.SetDefaultCacher(cacher)

	err = engine.Ping()
	if err != nil {
		Log.Fatalf("db err is %s", err)
	}

	mysql = &Mysql{
		im: engine,
	}
	return mysql
}

func (s *Dao) NoticeList(tp, startRow, endRow int32, u *[]model.NoticeStruct) int32 {
	//err := s.mysql.im.Find(&u)
	total, err := s.mysql.im.Where("type =?", tp).Count(&u)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	if startRow > endRow || startRow > int32(total) {
		Log.Errorln("查询的其实行列不合法")
		return ERRCODE_UNKNOWN
	}
	s.mysql.im.Where("type=?", tp).Limit(int(startRow), int(endRow)).Find(&u)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	return ERRCODE_SUCCESS
}

func (s *Dao) NoticeDescription(Id int32, u *model.NoticeDetailStruct) int32 {
	u = &model.NoticeDetailStruct{}
	ok, err := s.mysql.im.Where("ID=?", Id).Get(u)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	if ok {
		return ERRCODE_SUCCESS

	}

	return ERRCODE_ACCOUNT_NOTEXIST

}
