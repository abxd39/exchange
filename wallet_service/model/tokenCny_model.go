package models

import "digicon/wallet_service/utils"

type ConfigTokenCny struct {
	TokenId          int       `xorm:"token_id not null pk INT(11)"`
	Price       int       `xorm:"price"`
	UsdPrice    int       `xorm:"usd_price"`
}

func (this *ConfigTokenCny) GetById(token_id int) (bool,error) {
	return utils.Engine_token.Where("token_id = ?",token_id).Get(this)
}
