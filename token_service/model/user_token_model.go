package model

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	. "digicon/token_service/log"
	"errors"
	"github.com/go-xorm/xorm"
)

type UserToken struct {
	Id      int64
	Uid     uint64 `xorm:"unique(currency_uid) INT(11)"`
	TokenId int    `xorm:"comment('币种') unique(currency_uid) INT(11)"`
	Balance int64  `xorm:"comment('余额') BIGINT(20)"`
	Frozen  int64  `xorm:"comment('冻结余额') BIGINT(20)"`
	Version int    `xorm:"version"`
}

//获取实体
func (s *UserToken) GetUserToken(uid uint64, token_id int) (err error) {
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

		ok, err = DB.GetMysqlConn().Where("uid=? and token_id=?", uid, token_id).Get(s)
		if err != nil {
			Log.Errorln(err.Error())
			return
		}

		if !ok {
			errors.New("insert user token err")
			return
		}
	}

	return
}

//加代币数量
func (s *UserToken) AddMoney(session *xorm.Session, num int64) (err error) {

	s.Balance += num
	_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("balance").Update(s)
	//_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Decr("balance", num).Update(&UserToken{})

	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	s.Version += 1
	return
}

/*
func (s *UserToken) AddMoney(num int64, ukey string, ty int) (ret int32, err error) {
	m := &MoneyRecord{}
	ok, err := m.CheckExist(ukey, ty)
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

	_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Incr("balance", num).Update(&UserToken{})

	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		return
	}

	_, err = session.InsertOne(&MoneyRecord{
		Uid:     s.Uid,
		TokenId: s.TokenId,
		Num:     num,
		Ukey:    ukey,
		Opt:     int(proto.TOKEN_OPT_TYPE_ADD),
		Type:    ty,
	})

	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		return
	}

	err = session.Commit()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}
*/

func (s *UserToken) SubMoney(session *xorm.Session, num int64) (ret int32, err error) {
	if s.Balance < num {
		ret = ERR_TOKEN_LESS
	}

	s.Balance -= num
	_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("balance").Decr("balance", num).Update(s)
	//_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Decr("balance", num).Update(s)

	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	s.Version += 1

	return
}

//减代币数量
/*
func (s *UserToken) SubMoney(session *xorm.Session, num int64, ukey string, ty int) (ret int32, err error) {
	m := &MoneyRecord{}
	ok, err := m.CheckExist(ukey, ty)
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
	if session == nil {
		//开始事务入账处理
		session = DB.GetMysqlConn().NewSession()
		defer session.Close()
		err = session.Begin()

		_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Decr("balance", num).Update(&UserToken{})

		if err != nil {
			Log.Errorln(err.Error())
			session.Rollback()
			return
		}

		_, err = session.InsertOne(&MoneyRecord{
			Uid:     s.Uid,
			TokenId: s.TokenId,
			Ukey:    ukey,
			Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
			Type:    ty,
		})

		if err != nil {
			Log.Errorln(err.Error())
			session.Rollback()
			return
		}

		err = session.Commit()
		if err != nil {
			Log.Errorln(err.Error())
			return
		}
	} else {

		//开始事务入账处理
		_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Decr("balance", num).Update(&UserToken{})

		if err != nil {
			Log.Errorln(err.Error())
			return
		}

		_, err = session.InsertOne(&MoneyRecord{
			Uid:     s.Uid,
			TokenId: s.TokenId,
			Ukey:    ukey,
			Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
			Type:    ty,
		})

		if err != nil {
			Log.Errorln(err.Error())
			return
		}
	}

	return
}
*/

//冻结资金
func (s *UserToken) SubMoneyWithFronzen(sess *xorm.Session, num int64, entrust_id string, ty int) (ret int32, err error) {
	var aff int64
	if s.Balance >= num {
		s.Balance -= num
		s.Frozen += num
		aff, err = sess.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("balance", "frozen").Update(s)
		if err != nil {
			Log.Errorln(err.Error())
			ret = ERRCODE_UNKNOWN
			return
		}

		if aff == 0 {
			err = errors.New("update balance err version is wrong")
			ret = ERRCODE_UNKNOWN
			return
		}

		s.Version += 1

		f := Frozen{
			Uid:     s.Uid,
			Ukey:    entrust_id,
			Num:     num,
			TokenId: s.TokenId,
			Type:    ty,
			Opt:     int(proto.TOKEN_OPT_TYPE_ADD),
		}

		_, err = sess.Insert(f)
		if err != nil {
			Log.Errorln(err.Error())
			ret = ERRCODE_UNKNOWN
			return
		}

		return
	}

	ret = ERR_TOKEN_LESS
	return
}

//消耗冻结资金
func (s *UserToken) NotifyDelFronzen(sess *xorm.Session, num int64, entrust_id string, ty int) (ret int32, err error) {
	if s.Frozen < num {
		ret = ERR_TOKEN_LESS
		return
	}

	var aff int64
	s.Frozen -= num

	aff, err = sess.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("frozen").Update(s)
	if err != nil {

		Log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}
	if aff == 0 {
		err = errors.New("update balance err version is wrong")
		ret = ERRCODE_UNKNOWN
		return
	}
	s.Version += 1

	f := Frozen{
		Uid:     s.Uid,
		Ukey:    entrust_id,
		Num:     num,
		TokenId: s.TokenId,
		Type:    ty,
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
	}

	_, err = sess.Insert(f)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	return
}

//返还冻结资金
func (s *UserToken) ReturnFronzen(sess *xorm.Session, num int64, entrust_id string) (ret int32, err error) {
	return
}
