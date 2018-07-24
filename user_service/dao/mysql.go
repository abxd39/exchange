package dao

import (
	"digicon/user_service/conf"
	log "github.com/sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Mysql struct {
	im *xorm.Engine
}

func NewMysql() (mysql *Mysql) {
	dsource := conf.Cfg.MustValue("mysql", "conn")
	//root:current@tcp(47.106.136.96:3306)/rumi?charset=utf8
	engine, err := xorm.NewEngine("mysql", dsource)
	if err != nil {
		log.Fatalf("db err is %s", err)
	}

	engine.ShowSQL(true)
	//cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Hour, 1000)
	//engine.SetDefaultCacher(cacher)

	err = engine.Ping()
	if err != nil {
		log.Fatalf("db err is %s", err)
	}

	mysql = &Mysql{
		im: engine,
	}
	return mysql
}

func (s *Dao) GetMysqlConn() *xorm.Engine {
	return s.mysql.im
}

/*
func (s *Dao) RegisterByPhone(req *proto.RegisterPhoneRequest) int32 {
	if ret := s.CheckUserExist(req.Phone, "phone"); ret != ERRCODE_SUCCESS {
		return ret
	}

	e := &model.User{
		Pwd:   req.Pwd,
		Phone: req.Phone,
	}
	_, err := s.mysql.im.Cols("pwd", "phone").Insert(e)
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	_, err = s.mysql.im.Where("phone=?", req.Phone).Get(e)
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	m := &model.UserEx{
		Uid:          e.Uid,
		RegisterTime: time.Now().Unix(),
		InviteCode:   req.InviteCode,
	}

	_, err = s.mysql.im.Insert(m)
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	return ERRCODE_SUCCESS
}

func (s *Dao) RegisterByEmail(req *proto.RegisterEmailRequest) int32 {
	if ret := s.CheckUserExist(req.Email, "email"); ret != ERRCODE_SUCCESS {
		return ret
	}

	e := &model.User{
		Pwd:   req.Pwd,
		Email: req.Email,
	}
	_, err := s.mysql.im.Cols("pwd", "email").Insert(e)
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	_, err = s.mysql.im.Where("email=?", req.Email).Get(e)
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	m := &model.UserEx{
		Uid:          e.Uid,
		RegisterTime: time.Now().Unix(),
		InviteCode:   req.InviteCode,
	}

	_, err = s.mysql.im.Insert(m)
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	return ERRCODE_SUCCESS
}

func (s *Dao) CheckUserExist(param string, col string) int32 {
	sql := fmt.Sprintf("%s=?", col)
	ok, err := s.mysql.im.Where(sql, param).Get(&model.User{})
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	if ok {
		return ERRCODE_ACCOUNT_EXIST
	}
	return ERRCODE_SUCCESS
}

func (s *Dao) Login(phone, pwd string) int32 {
	m := &model.User{}
	ok, err := s.mysql.im.Where("phone=?", phone).Get(m)
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	if ok {
		if m.Pwd == pwd {
			return ERRCODE_SUCCESS
		}
		return ERRCODE_PWD
	}
	return ERRCODE_ACCOUNT_NOTEXIST
}


func (s Dao) GetUserByPhone(phone string) (u *model.User, ret int32) {
	u = &model.User{}
	ok, err := s.mysql.im.Where("phone=?", phone).Get(u)
	if err != nil {
		log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}

	if ok {
		ret = ERRCODE_SUCCESS
		return
	}
	ret = ERRCODE_ACCOUNT_NOTEXIST
	return
}

func (s *Dao) GetUserExByPhone(phone string) (u *model.UserEx, ret int32) {
	u = &model.UserEx{}
	ok, err := s.mysql.im.Where("phone=?", phone).Get(u)
	if err != nil {
		log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}

	if ok {
		ret = ERRCODE_SUCCESS
		return
	}

	ret = ERRCODE_ACCOUNT_NOTEXIST
	return
}

func (s *Dao) ModifyPwd(phone string, pwd string) (ret int32) {
	return 0
}
*/
