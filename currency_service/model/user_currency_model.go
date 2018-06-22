package model

import (
	"digicon/currency_service/dao"
	. "digicon/currency_service/log"
)

// 用户虚拟货币资产表
type UserCurrency struct {
	Id        uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid       uint64 `xorm:"INT(10)"     json:"uid"`                                          // 用户ID
	TokenId   uint32 `xorm:"INT(10)"     json:"token_id"`                                     // 虚拟货币类型
	TokenName string `xorm:"VARCHAR(36)" json:"token_name"`                                   // 虚拟货币名字
	Freeze    uint64 `xorm:"BIGINT not null default 0"   json:"freeze"`                       // 冻结
	Balance   uint64 `xorm:"not null default 0 comment('余额') BIGINT"   json:"balance"`        // 余额
	Address   string `xorm:"not null default '' comment('充值地址') VARCHAR(255)" json:"address"` // 充值地址
}

func (this *UserCurrency) Get(uid uint64, token_id uint32) *UserCurrency {

	data := new(UserCurrency)
	isdata, err := dao.DB.GetMysqlConn().Where("uid=? AND token_id=?", uid, token_id).Get(data)
	if err != nil {
		Log.Errorln(err.Error())
		return nil
	}

	if !isdata {
		return nil
	}

	return data
}
