package model

//
//import (
//	"digicon/currency_service/dao"
//	log "github.com/sirupsen/logrus"
//)
//
//// 货币类型表
//type Tokens struct {
//	Id     uint32 `xorm:"not null pk autoincr INT(10)" json:"id"`
//	Name   string `xorm:"VARBINARY(20)" json:"name"`
//	CnName string `xorm:"VARBINARY(20)" json:"cn_name"`
//}
//
//// 获取货币类型
//func (this *Tokens) Get(id uint32, name string) *Tokens {
//
//	data := new(Tokens)
//	var isdata bool
//	var err error
//	if id > 0 {
//		isdata, err = dao.DB.GetMysqlConn().Id(id).Get(data)
//	} else {
//		isdata, err = dao.DB.GetMysqlConn().Where("name=?", name).Get(data)
//	}
//
//	if err != nil {
//		log.Errorln(err.Error())
//		return nil
//	}
//
//	if isdata == false {
//		return nil
//	}
//
//	return data
//}
//
//// 获取货币类型列表
//func (this *Tokens) List() []Tokens {
//
//	data := make([]Tokens, 0)
//	err := dao.DB.GetMysqlConn().Find(&data)
//	if err != nil {
//		log.Errorln(err.Error())
//		return nil
//	}
//
//	return data
//}
