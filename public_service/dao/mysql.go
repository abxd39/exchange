package dao

import (
	. "digicon/proto/common"
	"digicon/public_service/conf"
	. "digicon/public_service/log"
	"digicon/public_service/model"
	"log"
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

func (s *Dao) ArticlesList(tp, page, page_num int32, u *[]model.Articles_list) int32 {
	//err := s.mysql.im.Find(&u)
	default_page := int32(10)
	list := new(model.Articles_list)
	if page_num <= 0 {
		page_num = default_page
	}
	total, err := s.mysql.im.Where("type =?", tp).Count(list)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	var count int32
	if int32(total) <= default_page || page == 1 {
		count = 1
	} else {
		count = page * page_num
	}
	err = s.mysql.im.Where("type=?", tp).Limit(int(page_num), int(count)).Find(u)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return ERRCODE_SUCCESS

}

func (s *Dao) ArticlesDescription(Id int32, u *model.ArticlesCopy1) int32 {
	//u = &model.ArticlesCopy1{}
	ok, err := s.mysql.im.Where("id=?", Id).Get(u)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	if ok {
		return ERRCODE_SUCCESS

	}

	return ERRCODE_ACCOUNT_NOTEXIST

}
