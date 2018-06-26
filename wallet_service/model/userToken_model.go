package models

import "digicon/wallet_service/utils"

type UserToken struct {
	Uid     int    `xorm:"unique(currency_uid) INT(11)"`
	TokenId int    `xorm:"comment('币种') unique(currency_uid) INT(11)"`
	Balance string `xorm:"comment('余额') DECIMAL(64,8)"`
}

func (this *UserToken) Add(amount string, uid, tokenid int) {
	utils.Engine_token.Incr("balance", amount).Where("uid=? and token_id", uid, tokenid)
}
