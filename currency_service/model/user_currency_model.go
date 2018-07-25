package model

import (
	"digicon/currency_service/dao"
	"errors"
	log "github.com/sirupsen/logrus"
)

// 用户虚拟货币资产表
type UserCurrency struct {
	Id        uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid       uint64 `xorm:"INT(10)"     json:"uid"`                                          // 用户ID
	TokenId   uint32 `xorm:"INT(10)"     json:"token_id"`                                     // 虚拟货币类型
	TokenName string `xorm:"VARCHAR(36)" json:"token_name"`                                   // 虚拟货币名字
	Freeze    int64  `xorm:"BIGINT not null default 0"   json:"freeze"`                       // 冻结
	Balance   int64  `xorm:"not null default 0 comment('余额') BIGINT"   json:"balance"`        // 余额
	Address   string `xorm:"not null default '' comment('充值地址') VARCHAR(255)" json:"address"` // 充值地址
	Version   int64  `xorm:"version"`
}

func (this *UserCurrency) Get(id uint64, uid uint64, token_id uint32) *UserCurrency {

	data := new(UserCurrency)
	var isdata bool
	var err error

	if id > 0 {
		isdata, err = dao.DB.GetMysqlConn().Id(id).Get(data)
	} else {
		isdata, err = dao.DB.GetMysqlConn().Where("uid=? AND token_id=?", uid, token_id).Get(data)
	}

	if err != nil {
		log.Errorln(err.Error())
		return nil
	}

	if !isdata {
		return nil
	}

	return data
}

func (this *UserCurrency) GetUserCurrency(uid uint64, nozero bool) (uCurrenList []UserCurrency, err error) {
	engine := dao.DB.GetMysqlConn()
	if nozero {
		err = engine.Where("uid = ?", uid).Find(&uCurrenList)
	} else {
		err = engine.Where("uid=? AND balance=?", uid, 0).Find(&uCurrenList)
	}
	return
}

func (this *UserCurrency) GetBalance(uid uint64, token_id uint32) (data UserCurrency, err error) {
	//data := new(UserCurrency)
	_, err = dao.DB.GetMysqlConn().Where("uid=? AND token_id=?", uid, token_id).Get(&data)
	return

}

func (this *UserCurrency) SetBalance(uid uint64, token_id uint32, amount int64) (err error) {
	engine := dao.DB.GetMysqlConn()
	sql := "UPDATE user_currency SET   balance= balance + ?, version = version + 1 WHERE uid = ? AND token_id = ? AND version = ?"
	sqlRest, err := engine.Exec(sql, amount, uid, token_id, this.Version)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	if rst, _ := sqlRest.RowsAffected(); rst == 0 {
		log.Errorln("添加余额失败")
		err = errors.New("添加余额失败!")
		return err
	}
	return
}
