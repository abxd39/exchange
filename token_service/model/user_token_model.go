package model

import (
	"digicon/common/errors"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	"fmt"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"

	"digicon/common/snowflake"
	"digicon/common/xtime"
	"time"
)

type UserToken struct {
	Id      int64
	Uid     uint64 `xorm:"unique(currency_uid) INT(11)"`
	TokenId int    `xorm:"comment('币种') unique(currency_uid) INT(11)"`
	Balance int64  `xorm:"comment('余额') BIGINT(20)"`
	Frozen  int64  `xorm:"comment('冻结余额') BIGINT(20)"`
	Version int    `xorm:"version"`
}

type UserTokenWithName struct {
	UserToken `xorm:"extends"`
	WorthCny  float64
	TokenName string
}

type UserTokenTotalMoney struct {
	TotalCny float64
	TotalUsd float64
}

func (*UserToken) TableName() string {
	return "user_token"
}

// 计算用户所有币的总金额，人民币、美元等
func (s *UserToken) CalcTotalMoney(uid uint64) (*UserTokenTotalMoney, error) {
	userTokenTotal := &UserTokenTotalMoney{}

	engine := DB.GetMysqlConn()
	_, err := engine.SQL(fmt.Sprintf("SELECT SUM(tmp.cny) AS total_cny,SUM(tmp.usd) AS total_usd FROM"+
		" (SELECT (ut.balance+ut.frozen)/100000000 * ctc.price/100000000 AS cny,"+
		" (ut.balance+ut.frozen)/100000000 * ctc.usd_price/100000000 AS usd"+
		" FROM %s ut LEFT JOIN %s ctc ON ctc.token_id=ut.token_id WHERE ut.uid=%d GROUP BY ut.token_id"+
		") tmp", s.TableName(), new(ConfigTokenCny).TableName(), uid)).Get(userTokenTotal)
	if err != nil {
		return nil, err
	}

	return userTokenTotal, nil
}

// 用户币币余额列表
func (s *UserToken) GetUserTokenList(filter map[string]interface{}) ([]UserTokenWithName, error) {
	engine := DB.GetMysqlConn()
	query := engine.Where("1=1")

	// 筛选
	if v, ok := filter["uid"]; ok {
		query.And("ut.uid=?", v)
	}
	if _, ok := filter["no_zero"]; ok {
		query.And("ut.balance!=0 OR ut.frozen!=0")
	}
	if v, ok := filter["token_id"]; ok {
		query.And("ut.token_id=?", v)
	}

	var list []UserTokenWithName
	err := query.
		Table(s).
		Alias("ut").
		Select("ut.*, t.mark as token_name, (ut.balance+ut.frozen)/100000000 * ctc.price/100000000 as worth_cny").
		Join("LEFT", []string{new(CommonTokens).TableName(), "t"}, "t.id=ut.token_id").
		Join("LEFT", []string{new(ConfigTokenCny).TableName(), "ctc"}, "ctc.token_id=ut.token_id").
		Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

//获取实体
func (s *UserToken) GetUserToken(uid uint64, token_id int) (err error) {
	var ok bool
	ok, err = DB.GetMysqlConn().Where("uid=? and token_id=?", uid, token_id).Get(s)
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	if !ok {
		s.Uid = uid
		s.TokenId = int(token_id)

		_, err = DB.GetMysqlConn().InsertOne(s)
		if err != nil {
			log.Errorln(err.Error())
			return
		}

		ok, err = DB.GetMysqlConn().Where("uid=? and token_id=?", uid, token_id).Get(s)
		if err != nil {
			log.Errorln(err.Error())
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
func (s *UserToken) AddMoney(session *xorm.Session, num int64, ukey string, ty proto.TOKEN_TYPE_OPERATOR) (err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"num":      num,
				"uid":      s.Uid,
				"token_id": s.TokenId,
				"balance":  s.Balance,
			}).Errorf("add  money  error %s", err.Error())
		}
	}()

	s.Balance += num
	_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("balance").Update(s)
	if err != nil {
		return
	}

	//交易流水
	err = InsertRecord(session, &MoneyRecord{
		Uid:     s.Uid,
		TokenId: s.TokenId,
		Ukey:    ukey,
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
		Type:    int(ty),
		Num:     num,
		Balance: s.Balance,
	})
	if err != nil {
		return
	}
	return
}

//加冻结代币数量，如：注册赠送代币默认放到冻结代币里
func (s *UserToken) AddFrozen(session *xorm.Session, num int64, ukey string, ty proto.TOKEN_TYPE_OPERATOR) (err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"num":      num,
				"uid":      s.Uid,
				"token_id": s.TokenId,
				"balance":  s.Balance,
			}).Errorf("add frozen money  error %s", err.Error())
		}
	}()

	s.Frozen += num
	_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("frozen").Update(s)
	if err != nil {
		return
	}

	_, err = session.Insert(&Frozen{
		Uid:     s.Uid,
		Ukey:    ukey,
		Num:     num,
		TokenId: s.TokenId,
		Type:    int(ty),
		Opt:     int(proto.TOKEN_OPT_TYPE_ADD),
	})
	if err != nil {
		return
	}
	return
}

//减少代币数量
func (s *UserToken) SubMoney(session *xorm.Session, num int64) (ret int32, err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"num":      num,
				"uid":      s.Uid,
				"token_id": s.TokenId,
				"balance":  s.Balance,
			}).Errorf("sub money info data error %s", err.Error())
		}
	}()
	if s.Balance < num {
		ret = ERR_TOKEN_LESS
	}

	s.Balance -= num
	_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("balance").Update(s)
	if err != nil {
		return
	}
	return
}

//减代币数量
/*
func (s *UserToken) SubMoney(session *xorm.Session, num int64, ukey string, ty int) (ret int32, err error) {
	m := &MoneyRecord{}
	ok, err := m.CheckExist(ukey, ty)
	if err != nil {
		log.Errorln(err.Error())
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
			log.Errorln(err.Error())
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
			log.Errorln(err.Error())
			session.Rollback()
			return
		}

		err = session.Commit()
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	} else {

		//开始事务入账处理
		_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Decr("balance", num).Update(&UserToken{})

		if err != nil {
			log.Errorln(err.Error())
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
			log.Errorln(err.Error())
			return
		}
	}

	return
}
*/

//冻结资金
func (s *UserToken) SubMoneyWithFronzen(sess *xorm.Session, num int64, entrust_id string, ty proto.TOKEN_TYPE_OPERATOR) (ret int32, err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"num":         num,
				"uid":         s.Uid,
				"token_id":    s.TokenId,
				"balance":     s.Balance,
				"entrusdt_id": entrust_id,
				"ty":          ty,
			}).Errorf("sub  money with fronzen error %s", err.Error())
		}
	}()
	var aff int64
	if s.Balance < num {
		ret = ERR_TOKEN_LESS
	}

	s.Balance -= num
	s.Frozen += num
	aff, err = sess.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("balance", "frozen").Update(s)
	if err != nil {
		ret = ERRCODE_UNKNOWN

		return
	}

	if aff == 0 {
		err = errors.New("update balance err version is wrong")
		ret = ERRCODE_UNKNOWN

		return
	}

	f := Frozen{
		Uid:     s.Uid,
		Ukey:    entrust_id,
		Num:     num,
		TokenId: s.TokenId,
		Type:    int(ty),
		Opt:     int(proto.TOKEN_OPT_TYPE_ADD),
	}

	_, err = sess.Insert(f)
	if err != nil {
		ret = ERRCODE_UNKNOWN

		return
	}

	//交易流水
	err = InsertRecord(sess, &MoneyRecord{
		Uid:     s.Uid,
		TokenId: s.TokenId,
		Ukey:    entrust_id,
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
		Type:    int(ty),
		Num:     num,
		Balance: s.Balance,
	})
	if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}

	return

}

//消耗冻结资金
func (s *UserToken) NotifyDelFronzen(sess *xorm.Session, num int64, ukey string, ty proto.TOKEN_TYPE_OPERATOR) (err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"num":         num,
				"uid":         s.Uid,
				"token_id":    s.TokenId,
				"balance":     s.Balance,
				"entrusdt_id": ukey,
				"ty":          ty,
			}).Errorf("notify  money del fronzen error %s", err.Error())
		}
	}()
	if s.Frozen < num {
		err = errors.New("please check why fronze num is less")
		return
	}

	var aff int64
	s.Frozen -= num

	aff, err = sess.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("frozen").Update(s)
	if err != nil {
		return
	}
	if aff == 0 {
		err = errors.New("update balance err version is wrong")
		return
	}

	f := Frozen{
		Uid:     s.Uid,
		Ukey:    ukey,
		Num:     num,
		TokenId: s.TokenId,
		Type:    int(ty),
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
	}

	_, err = sess.Insert(f)
	if err != nil {
		return
	}

	return
}

//返还冻结资金
func (s *UserToken) ReturnFronzen(sess *xorm.Session, num int64, entrust_id string, ty proto.TOKEN_TYPE_OPERATOR) (err error) {
	if s.Frozen < num {
		err = errors.New("please check fronze data because fronzn num is not enough")
		return
	}
	var aff int64
	s.Balance += num
	s.Frozen -= num
	aff, err = sess.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("frozen", "balance").Update(s)
	if err != nil {
		return
	}
	if aff == 0 {
		err = errors.New("update balance err version is wrong")
		return
	}

	_, err = sess.Insert(&Frozen{
		Uid:     s.Uid,
		Ukey:    entrust_id,
		Num:     num,
		TokenId: s.TokenId,
		Type:    int(ty),
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
	})
	if err != nil {
		return
	}

	err = InsertRecord(sess, &MoneyRecord{
		Uid:     s.Uid,
		TokenId: int(s.TokenId),
		Ukey:    entrust_id,
		Opt:     int(proto.TOKEN_OPT_TYPE_ADD),
		Type:    int(ty),
		Balance: s.Balance,
		Num:     num,
	})

	return
}

//获取个人资金明细
func (s *UserToken) GetAllToken(uid uint64) []*UserToken {
	r := make([]*UserToken, 0)
	err := DB.GetMysqlConn().Where("uid=?", uid).Find(&r)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	return r
}

//划转到法币
func (s *UserToken) TransferToCurrency(uid uint64, tokenId int, tokenName string, num int64) error {
	//检查代币是否足够
	err := s.GetUserToken(uid, tokenId)
	if err != nil {
		return err
	}
	if s.Balance < num {
		return errors.NewNormal("余额不足")
	}

	//整理数据
	now := time.Now().Unix()
	xid := snowflake.SnowflakeNode.Generate()

	//开始划转，XA事务
	tokenSession := DB.GetMysqlConn().NewSession()
	currencySession := DB.GetCurrencyMysqlConn().NewSession()
	defer func() {
		tokenSession.Close()
		currencySession.Close()
	}()

	//两个库指定同一个事务id，表明两个库的操作处于同一个事务中
	tokenSession.Exec(fmt.Sprintf("XA START '%d'", xid))
	currencySession.Exec(fmt.Sprintf("XA START '%d'", xid))

	//1.代币账户减
	newBalance := s.Balance - num
	_, err = tokenSession.Exec(fmt.Sprintf("UPDATE %s SET"+
		" balance=%d,"+
		" version=version+1"+
		" WHERE uid=%d"+
		" AND token_id=%d"+
		" AND version=%d", s.TableName(), newBalance, uid, tokenId, s.Version))
	if err != nil {
		tokenSession.Exec(fmt.Sprintf("XA END '%d'", xid))
		tokenSession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))
		currencySession.Exec(fmt.Sprintf("XA END '%d'", xid))
		currencySession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))

		tokenSession.Exec(fmt.Sprintf("XA ROLLBACK '%d'", xid))
		currencySession.Exec(fmt.Sprintf("XA ROLLBACK '%d'", xid))

		return errors.NewSys(err)
	}

	//2.代币流水
	_, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
		" (uid, token_id, ukey, type, opt, balance, num, created_time) VALUES"+
		" (%d, %d, %d, %d, %d, %d, %d, %d)", new(MoneyRecord).TableName(), uid, tokenId, xid, proto.TOKEN_TYPE_OPERATOR_TRANSFER_TO_CURRENCY, proto.TOKEN_OPT_TYPE_DEL, newBalance, num, now))
	if err != nil {
		tokenSession.Exec(fmt.Sprintf("XA END '%d'", xid))
		tokenSession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))
		currencySession.Exec(fmt.Sprintf("XA END '%d'", xid))
		currencySession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))

		tokenSession.Exec(fmt.Sprintf("XA ROLLBACK '%d'", xid))
		currencySession.Exec(fmt.Sprintf("XA ROLLBACK '%d'", xid))

		return errors.NewSys(err)
	}

	//3.法币账号加
	userCurrency := OutUserCurrency{}
	has, err := currencySession.Clone().Where("uid=?", uid).And("token_id=?", tokenId).Get(&userCurrency)
	if err != nil {
		return errors.NewSys(err)
	}
	newCurrBalance := userCurrency.Balance + num
	userCurrencyTableName := userCurrency.TableName()

	if !has {
		_, err = currencySession.Exec(fmt.Sprintf("INSERT INTO %s"+
			" (uid, token_id, token_name, balance, version)"+
			" VALUES (%d, %d, '%s', %d, 1)", userCurrencyTableName, uid, tokenId, tokenName, newCurrBalance))
		if err != nil {
			tokenSession.Exec(fmt.Sprintf("XA END '%d'", xid))
			tokenSession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))
			currencySession.Exec(fmt.Sprintf("XA END '%d'", xid))
			currencySession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))

			tokenSession.Exec(fmt.Sprintf("XA ROLLBACK '%d'", xid))
			currencySession.Exec(fmt.Sprintf("XA ROLLBACK '%d'", xid))

			return errors.NewSys(err)
		}
	} else {
		_, err = currencySession.Exec(fmt.Sprintf("UPDATE %s SET balance=%d,version=version+1 WHERE uid=%d AND token_id=%d AND version=%d", userCurrencyTableName, newCurrBalance, uid, tokenId, userCurrency.Version))
		if err != nil {
			tokenSession.Exec(fmt.Sprintf("XA END '%d'", xid))
			tokenSession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))
			currencySession.Exec(fmt.Sprintf("XA END '%d'", xid))
			currencySession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))

			tokenSession.Exec(fmt.Sprintf("XA ROLLBACK '%d'", xid))
			currencySession.Exec(fmt.Sprintf("XA ROLLBACK '%d'", xid))

			return errors.NewSys(err)
		}
	}

	//4.法币流水
	_, err = currencySession.Exec(fmt.Sprintf("INSERT INTO %s"+
		" (uid, trade_uid, order_id, token_id, num, surplus, operator, created_time)"+
		" VALUES (%d, %[2]d, '%d', %d, %d, %d, %d, '%s')", new(OutUserCurrencyHistory).TableName(), uid, xid, tokenId, num, newCurrBalance, 3, xtime.Unix2Date(now, xtime.LAYOUT_DATE_TIME)))
	if err != nil {
		tokenSession.Exec(fmt.Sprintf("XA END '%d'", xid))
		tokenSession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))
		currencySession.Exec(fmt.Sprintf("XA END '%d'", xid))
		currencySession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))

		tokenSession.Exec(fmt.Sprintf("XA ROLLBACK '%d'", xid))
		currencySession.Exec(fmt.Sprintf("XA ROLLBACK '%d'", xid))

		return errors.NewSys(err)
	}

	//提交
	tokenSession.Exec(fmt.Sprintf("XA END '%d'", xid))
	tokenSession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))
	currencySession.Exec(fmt.Sprintf("XA END '%d'", xid))
	currencySession.Exec(fmt.Sprintf("XA PREPARE '%d'", xid))

	tokenSession.Exec(fmt.Sprintf("XA COMMIT '%d'", xid))
	currencySession.Exec(fmt.Sprintf("XA COMMIT '%d'", xid))

	return nil
}
