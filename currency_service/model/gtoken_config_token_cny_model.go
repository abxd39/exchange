package model

import (
	"digicon/currency_service/dao"
)

//  货币价格汇率
type TokenConfigTokenCNy struct {
	TokenId int   `xorm:"not null pk comment(' 币类型') INT(10)"`
	Price   int64 `xorm:"comment('人民币价格') BIGINT(20)"`
}


func (this *TokenConfigTokenCNy) GetPrice(tokenid uint32) (err error){
	_, err = dao.DB.GetTokenMysqlConn().Table("config_token_cny").Where("token_id =?", tokenid).Get(this)
	return
}



func (this *TokenConfigTokenCNy) List()  (result []TokenConfigTokenCNy, err error){
	err = dao.DB.GetCommonMysqlConn().Table("config_token_cny").Find(&result)
	return
}