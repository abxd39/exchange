package model

import (
	"digicon/common/errors"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	"fmt"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"

	"digicon/common/constant"
	"digicon/common/snowflake"
	"github.com/gin-gonic/gin/json"
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
		" (SELECT SUM(ut.balance+ut.frozen)/100000000 * ctc.price/100000000 AS cny,"+
		" SUM(ut.balance+ut.frozen)/100000000 * ctc.usd_price/100000000 AS usd"+
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
		Join("LEFT", []string{new(OutCommonTokens).TableName(), "t"}, "t.id=ut.token_id").
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
		Opt:     int(proto.TOKEN_OPT_TYPE_ADD),
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

	_, err = session.Insert(&FrozenHistory{
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
func (s *UserToken) SubMoney(session *xorm.Session, num int64, ukey string, ty proto.TOKEN_TYPE_OPERATOR) (ret int32, err error) {
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

	f := FrozenHistory{
		Uid:     s.Uid,
		Ukey:    entrust_id,
		Num:     num,
		TokenId: s.TokenId,
		Type:    int(ty),
		Opt:     int(proto.TOKEN_OPT_TYPE_ADD),
		Frozen:  s.Frozen,
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

	f := FrozenHistory{
		Uid:     s.Uid,
		Ukey:    ukey,
		Num:     num,
		TokenId: s.TokenId,
		Type:    int(ty),
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
		Frozen:  s.Frozen,
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

	_, err = sess.Insert(&FrozenHistory{
		Uid:     s.Uid,
		Ukey:    entrust_id,
		Num:     num,
		TokenId: s.TokenId,
		Type:    int(ty),
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
		Frozen:  s.Frozen,
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

//划出到法币
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
	transferId := snowflake.SnowflakeNode.Generate()
	now := time.Now().Unix()

	//开始划转
	tokenSession := DB.GetMysqlConn().NewSession()
	defer tokenSession.Close()

	//事务
	err = tokenSession.Begin()
	if err != nil {
		return errors.NewSys(err)
	}

	//1.代币账户
	newBalance := s.Balance - num
	_, err = tokenSession.Exec(fmt.Sprintf("UPDATE %s SET"+
		" balance=%d,"+
		" version=version+1"+
		" WHERE uid=%d"+
		" AND token_id=%d"+
		" AND version=%d", s.TableName(), newBalance, uid, tokenId, s.Version))
	if err != nil {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}

	//2.代币流水
	_, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
		" (uid, token_id, ukey, type, opt, balance, num, created_time) VALUES"+
		" (%d, %d, %d, %d, %d, %d, %d, %d)", new(MoneyRecord).TableName(), uid, tokenId, transferId, proto.TOKEN_TYPE_OPERATOR_TRANSFER_TO_CURRENCY, proto.TOKEN_OPT_TYPE_DEL, newBalance, num, now))
	if err != nil {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}

	//3.划转记录
	_, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
		" (id, uid, token_id, num, states, create_time) VALUES"+
		" (%d, %d, %d, %d, %d, %d)", new(TransferRecord).TableName(), transferId, uid, tokenId, num, 1, now))
	if err != nil {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}

	//提交
	err = tokenSession.Commit()
	if err != nil {
		return errors.NewSys(err)
	}

	// 发送划转消息事件给currency服
	go func() {
		msg, err := json.Marshal(proto.TransferToCurrencyTodoMessage{
			Id:         int64(transferId),
			Uid:        int32(uid),
			TokenId:    int32(tokenId),
			Num:        num,
			CreateTime: now,
		})
		if err != nil {
			return
		}

		rdsClient := DB.GetCommonRedisConn()
		rdsClient.RPush(constant.RDS_TOKEN_TO_CURRENCY_TODO, msg)
	}()

	return nil
}

//划出到法币成功，把划转记录标记为已完成
func (s *UserToken) TransferToCurrencyDone(msg *proto.TransferToCurrencyDoneMessage) error {
	engine := DB.GetMysqlConn()
	_, err := engine.Exec(fmt.Sprintf("UPDATE %s SET states=2,update_time=%d WHERE id=%d AND states=1", new(TransferRecord).TableName(), time.Now().Unix(), msg.Id))
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}

//从法币转入
func (s *UserToken) TransferFromCurrency(msg *proto.TransferToTokenTodoMessage) error {
	//!!!!重要，判断消息是否已处理过
	rdsClient := DB.GetCommonRedisConn()
	isDid, history, err := new(MoneyRecord).IsTransferFromCurrencyDid(msg.Id)
	if err != nil {
		return err
	}
	if isDid { //已处理过，返回done消息并直接退出!
		msg, err := json.Marshal(proto.TransferToTokenDoneMessage{
			Id:       msg.Id,
			DoneTime: history.CreatedTime,
		})
		if err != nil {
			return errors.NewSys(err)
		}

		rdsClient.RPush(constant.RDS_CURRENCY_TO_TOKEN_DONE, msg)
		return nil
	}

	//整理数据
	now := time.Now().Unix()

	//开始处理
	tokenSession := DB.GetMysqlConn().NewSession()
	sessionClone := tokenSession.Clone()
	defer func() {
		tokenSession.Close()
		sessionClone.Close()
	}()

	//判断用户代币账户是否存在
	userToken := UserToken{}
	has, err := sessionClone.Where("uid=?", msg.Uid).And("token_id=?", msg.TokenId).Get(&userToken)
	if err != nil {
		return errors.NewSys(err)
	}
	newTokenBalance := userToken.Balance + msg.Num

	//1.代币账户
	if !has { //代币账号不存在，新建
		_, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
			" (uid, token_id, token_name, balance, version)"+
			" VALUES (%d, %d, '%s', %d, 1)", s.TableName(), msg.Uid, msg.TokenId, "", newTokenBalance))
		if err != nil {
			tokenSession.Rollback()
			return errors.NewSys(err)
		}
	} else { //代币账号已存在，更新
		_, err = tokenSession.Exec(fmt.Sprintf("UPDATE %s SET balance=%d,version=version+1 WHERE uid=%d AND token_id=%d AND version=%d", s.TableName(), newTokenBalance, msg.Uid, msg.TokenId, userToken.Version))
		if err != nil {
			tokenSession.Rollback()
			return errors.NewSys(err)
		}
	}

	//2.代币流水
	_, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
		" (uid, token_id, ukey, type, opt, balance, num, created_time)"+
		" VALUES (%d, %d, '%d', %d, %d, %d, %d, %d)", new(MoneyRecord).TableName(), msg.Uid, msg.TokenId, msg.Id, proto.TOKEN_TYPE_OPERATOR_TRANSFER_FROM_CURRENCY, proto.TOKEN_OPT_TYPE_ADD, newTokenBalance, msg.Num, now))
	if err != nil {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}

	//提交
	err = tokenSession.Commit()
	if err != nil {
		return errors.NewSys(err)
	}

	//发送处理成功消息给currency服
	go func() {
		doneMsg, err := json.Marshal(proto.TransferToTokenDoneMessage{
			Id:       msg.Id,
			DoneTime: now,
		})
		if err != nil {
			return
		}
		rdsClient.RPush(constant.RDS_CURRENCY_TO_TOKEN_DONE, doneMsg)
	}()

	return nil
}
