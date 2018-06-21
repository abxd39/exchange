package model

import (
	. "digicon/user_service/log"
	"github.com/shopspring/decimal"
	. "digicon/user_service/dao"
)

type UserToken struct {
	Uid     int    `xorm:"unique(currency_uid) INT(11)"`
	TokenId int    `xorm:"comment('币种') unique(currency_uid) INT(11)"`
	Balance string `xorm:"comment('余额') DECIMAL(20,4)"`
	Version int32
}

func (s *UserToken) AddMoney(uid int, token_id int, num string, hash string) (ret int32,err error) {
	m := &MoneyRecord{}
	ok, err := m.CheckExist(hash)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	if ok {
		return
	}
/*
	u := &UserToken{}
	ok, err = DB.GetMysqlConn().Where("uid=? and token_id=?", uid, token_id).Get(u)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	if !ok {
		u.Uid = uid
		u.TokenId = token_id

		_, err = DB.GetMysqlConn().InsertOne(u)
		if err != nil {
			Log.Errorln(err.Error())
			return
		}
	}
*/
	b, err := decimal.NewFromString(s.Balance)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	n, err := decimal.NewFromString(num)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	b.Add(n)

	//开始事务入账处理
	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()

	_, err = session.Cols("balance").Update(&UserToken{
		Balance: b.String(),
	})

	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		return
	}

	_, err = session.InsertOne(&MoneyRecord{
		Uid:     uid,
		TokenId: token_id,
		Num:     num,
		Hash:    hash,
		Opt:     true,
	})

	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		return
	}

	err = session.Commit()
	if err != nil {
		return
	}

}
