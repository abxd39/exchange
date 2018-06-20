package model

type UserToken struct {
	Uid     int    `xorm:"unique(currency_uid) INT(11)"`
	TokenId int    `xorm:"comment('币种') unique(currency_uid) INT(11)"`
	Balance string `xorm:"comment('余额') DECIMAL(20,4)"`
	Version int32
}

func (s *UserToken) AddMoney(uid int,token_id int,num string,hash string)  {

}