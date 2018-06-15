package dao

import (
	. "digicon/proto/common"
	"digicon/public_service/conf"
	. "digicon/public_service/log"
	"digicon/public_service/model"
	"fmt"
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

func (s *Dao) ArticlesList(tp, page, page_num int32, u *[]model.Articles_list) (int, int32) {
	//err := s.mysql.im.Find(&u)
	//default_page_count := int(10)
	var total_page int
	pagenum := int(page_num)
	list := new(model.Articles_list)
	if page <= 0 {
		page = 1
	}
	if page_num <= 0 {
		pagenum = 10
	}
	//没有指定 每页的行数
	var star_row int
	star_row = (int(page) - 1) * int(pagenum)

	total, err := s.mysql.im.Where("type =?", tp).Count(list)
	if err != nil {
		Log.Errorln(err.Error())
		return 0, ERRCODE_UNKNOWN
	}
	// var count int32
	// if int32(total) <=  || page == 1 {
	// 	count = 0
	// } else {
	// 	count = (page-1)*page_num - 1
	// }

	fmt.Println("total=", total, "type=", tp, "page=", page, "起始行star_row=", star_row, "page_num=", page_num)
	err = s.mysql.im.Where("type=?", tp).Limit(int(pagenum), int(star_row)).Find(u)
	if err != nil {
		log.Fatalf(err.Error())
	}

	total_page = int(total)
	total_page = total_page / pagenum
	fmt.Println("total_page=", total_page)
	return total_page, ERRCODE_SUCCESS

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
