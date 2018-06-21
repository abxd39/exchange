package model

import (
	. "digicon/proto/common"
	. "digicon/user_service/dao"
	. "digicon/user_service/log"
)

type UserToken struct {
	Uid     int   `xorm:"unique(currency_uid) INT(11)"`
	TokenId int   `xorm:"comment('币种') unique(currency_uid) INT(11)"`
	Balance int64 `xorm:"comment('余额') BIGINT(20)"`
	Version int
}

func (s *UserToken) GetUserToken(uid, token_id int) (err error) {
	var ok bool
	ok, err = DB.GetMysqlConn().Where("uid=? and token_id=?", uid, token_id).Get(s)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	if !ok {
		s.Uid = uid
		s.TokenId = token_id

		_, err = DB.GetMysqlConn().InsertOne(s)
		if err != nil {
			Log.Errorln(err.Error())
			return
		}
	}

	return
}

func (s *UserToken) AddMoney(num int64, hash string) (ret int32, err error) {
	m := &MoneyRecord{}
	ok, err := m.CheckExist(hash)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	if ok {
		ret = ERR_TOKEN_REPEAT
		return
	}

	//开始事务入账处理
	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()

	_, err = session.Where("uid=? and token=?", s.Uid, s.TokenId).Incr("balance=balance+?", num).Update(&UserToken{})

	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		return
	}

	_, err = session.InsertOne(&MoneyRecord{
		Uid:     s.Uid,
		TokenId: s.TokenId,
		Num:     num,
		Hash:    hash,
		Opt:     1,
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
	return
}

func (s *UserToken) SubMoney(num int64, hash string) (ret int32, err error) {
	m := &MoneyRecord{}
	ok, err := m.CheckExist(hash)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	if ok {
		ret = ERR_TOKEN_REPEAT
		return
	}

	if s.Balance < num {
		ret = ERR_TOKEN_LESS
	}
	//开始事务入账处理
	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()

	_, err = session.Where("uid=? and token=?", s.Uid, s.TokenId).Decr("balance=balance-?", num).Update(&UserToken{})

	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		return
	}

	_, err = session.InsertOne(&MoneyRecord{
		Uid:     s.Uid,
		TokenId: s.TokenId,
		Num:     num,
		Hash:    hash,
		Opt:     0,
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
	return
}
