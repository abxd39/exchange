package model

import (
	"digicon/currency_service/dao"
	log "github.com/sirupsen/logrus"
)

// 支付方式表
type Pays struct {
	Id     uint32 `xorm:"not null pk autoincr INT(10)" json:"id"`
	TypeId uint32 `xorm:"TINYINT(2)" json:"type_id"`
	ZhPay  string `xorm:"VARBINARY(20)" json:"zh_pay"`
	EnPay  string `xorm:"VARBINARY(20)" json:"en_pay"`
	States uint32 `xorm:"TINYINT(2)" json:"states"`
}

// 获取支付方式
func (this *Pays) Get(id uint32, name string) *Pays {

	data := new(Pays)
	var isdata bool
	var err error
	if id > 0 {
		isdata, err = dao.DB.GetMysqlConn().Id(id).Get(data)
	} else {
		isdata, err = dao.DB.GetMysqlConn().Where("en_pay=?", name).Get(data)
	}

	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	if isdata == false {
		return nil
	}

	return data
}

// 获取支付方式列表
func (this *Pays) List() []Pays {

	data := make([]Pays, 0)
	err := dao.DB.GetMysqlConn().Find(&data)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}

	return data
}
