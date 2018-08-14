package model

import (
	"digicon/common/errors"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	"fmt"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"

	"database/sql"
	"digicon/common/constant"
	"digicon/common/snowflake"
	"github.com/gin-gonic/gin/json"
	"time"
)

type UserToken struct {
	Uid        uint64 `xorm:"unique(currency_uid) INT(11)"`
	TokenId    int    `xorm:"comment('币种') unique(currency_uid) INT(11)"`
	TokenName  string `xorm:"token_name"`
	Balance    int64  `xorm:"comment('余额') BIGINT(20)"`
	Frozen     int64  `xorm:"comment('冻结余额') BIGINT(20)"`
	Version    int    `xorm:"version"`
	FrozenCny  int64  `xorm:"default 0 BIGINT(20)"`
	BalanceCny int64  `xorm:"default 0 BIGINT(20)"`
}

type UserTokenWithBalance struct {
	UserToken    `xorm:"extends"`
	TotalBalance int64
}

type UserTokenTotal struct {
	TokenId      int   `xorm:"token_id"`
	TotalBalance int64 `xorm:"total_balance"`
}

func (*UserToken) TableName() string {
	return "user_token"
}

// 计算用户所有币的总额
func (s *UserToken) CalcTotal(uid uint64) ([]*UserTokenTotal, error) {
	var userTokenTotal []*UserTokenTotal

	engine := DB.GetMysqlConn()
	err := engine.SQL(fmt.Sprintf("SELECT token_id, SUM(balance+frozen) AS total_balance"+
		" FROM %s WHERE uid=%d GROUP BY token_id",
		s.TableName(), uid)).Find(&userTokenTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return userTokenTotal, nil
}

// 用户币币余额列表
func (s *UserToken) GetUserTokenList(filter map[string]interface{}) ([]UserTokenWithBalance, error) {
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

	var list []UserTokenWithBalance
	err := query.
		Table(s).
		Alias("ut").
		Select("ut.*, (ut.balance+ut.frozen) as total_balance").
		Join("LEFT", []string{new(ConfigTokenCny).TableName(), "ctc"}, "ctc.token_id=ut.token_id").
		Find(&list)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return list, nil

	return nil, nil
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
		var token *OutCommonTokens
		token, err = new(OutCommonTokens).Get(uint32(token_id))
		if err != nil {
			return
		}

		s.Uid = uid
		s.TokenName = token.Mark
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

//获取实体
func (s *UserToken) GetUserTokenInSession(session *xorm.Session, uid uint64, token_id int, tokenName string) (userToken *UserToken, err error) {
	var ok bool
	userToken = &UserToken{}
	ok, err = session.Where("uid=? and token_id=?", uid, token_id).Get(userToken)
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	if !ok {
		s.Uid = uid
		s.TokenName = tokenName
		s.TokenId = int(token_id)

		_, err = session.InsertOne(s)
		if err != nil {
			log.Errorln(err.Error())
			return
		}

		ok, err = session.Where("uid=? and token_id=?", uid, token_id).Get(userToken)
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

//注册奖励
func (s *UserToken) RegisterReward(uid, rewardTokenId, rewardNum int64) error {
	if rewardTokenId == 0 || rewardNum == 0 { //无奖励
		return errors.NewNormal("请配置奖励币种、数量")
	}

	//判断是否已领取过奖励
	has, err := DB.GetMysqlConn().Where("uid=?", uid).And("type=?", proto.TOKEN_TYPE_OPERATOR_HISTORY_REGISTER).Exist(new(FrozenHistory))
	if err != nil {
		return errors.NewSys(err)
	}
	if has {
		return nil
	}

	//tokenName
	token, err := new(OutCommonTokens).Get(uint32(rewardTokenId))
	if err != nil {
		return err
	}

	//整理数据
	userTokenTableName := s.TableName()
	FrozenHistoryTableName := new(FrozenHistory).TableName()
	now := time.Now().Unix()

	//开始
	tokenSession := DB.GetMysqlConn().NewSession()
	defer tokenSession.Close()

	//事务
	err = tokenSession.Begin()
	if err != nil {
		return errors.NewSys(err)
	}

	//1.自己
	my, err := s.GetUserTokenInSession(tokenSession, uint64(uid), int(rewardTokenId), token.Mark)
	newFrozen := my.Frozen + rewardNum

	var result sql.Result
	result, err = tokenSession.Exec(fmt.Sprintf("UPDATE %s SET frozen=frozen+%d,version=version+1 WHERE uid=%d AND token_id=%d AND version=%d",
		userTokenTableName, rewardNum, uid, rewardTokenId, my.Version))
	if err != nil {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}

	//流水
	result, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s (uid, ukey, num, token_id, type, create_time, opt, frozen)"+
		" VALUES (%d, '%d', %d, %d, %d, %d, %d, %d)", FrozenHistoryTableName, uid, uid, rewardNum, rewardTokenId, proto.TOKEN_TYPE_OPERATOR_HISTORY_REGISTER, now, proto.TOKEN_OPT_TYPE_ADD, newFrozen))
	if err != nil {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}

	//2.推荐人
	//2.1一级推荐人
	userExMD := new(OutUserEx)
	userEx, err := userExMD.Get(uid)
	if err != nil {
		tokenSession.Rollback()
		return err
	}
	if inviteId := userEx.InviteId; inviteId > 0 { //有推荐人
		invite, err := s.GetUserTokenInSession(tokenSession, uint64(inviteId), int(rewardTokenId), token.Mark)
		if err != nil {
			tokenSession.Rollback()
			return errors.NewSys(err)
		}
		newFrozen = invite.Frozen + rewardNum

		result, err = tokenSession.Exec(fmt.Sprintf("UPDATE %s SET frozen=frozen+%d,version=version+1 WHERE uid=%d AND token_id=%d AND version=%d",
			userTokenTableName, rewardNum, inviteId, rewardTokenId, invite.Version))
		if err != nil {
			tokenSession.Rollback()
			return errors.NewSys(err)
		}
		if affected, err := result.RowsAffected(); err != nil || affected == 0 {
			tokenSession.Rollback()
			return errors.NewSys(err)
		}

		//流水
		result, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s (uid, ukey, num, token_id, type, create_time, opt, frozen)"+
			" VALUES (%d, '%s', %d, %d, %d, %d, %d, %d)", FrozenHistoryTableName, inviteId, fmt.Sprintf("%d-%d", uid, inviteId), rewardNum, rewardTokenId, proto.TOKEN_TYPE_OPERATOR_HISTORY_IVITE, now, proto.TOKEN_OPT_TYPE_ADD, newFrozen))
		if err != nil {
			tokenSession.Rollback()
			return errors.NewSys(err)
		}
		if affected, err := result.RowsAffected(); err != nil || affected == 0 {
			tokenSession.Rollback()
			return errors.NewSys(err)
		}

		//2.2二级推荐人
		inviteUserEx, err := userExMD.Get(inviteId)
		if err != nil {
			tokenSession.Rollback()
			return err
		}
		if secInviteId := inviteUserEx.InviteId; secInviteId > 0 { //有推荐人
			secInvite, err := s.GetUserTokenInSession(tokenSession, uint64(secInviteId), int(rewardTokenId), token.Mark)
			if err != nil {
				tokenSession.Rollback()
				return errors.NewSys(err)
			}
			newFrozen = secInvite.Frozen + rewardNum

			result, err = tokenSession.Exec(fmt.Sprintf("UPDATE %s SET frozen=frozen+%d,version=version+1 WHERE uid=%d AND token_id=%d AND version=%d",
				userTokenTableName, rewardNum, secInviteId, rewardTokenId, secInvite.Version))
			if err != nil {
				tokenSession.Rollback()
				return errors.NewSys(err)
			}
			if affected, err := result.RowsAffected(); err != nil || affected == 0 {
				tokenSession.Rollback()
				return errors.NewSys(err)
			}

			//流水
			result, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s (uid, ukey, num, token_id, type, create_time, opt, frozen)"+
				" VALUES (%d, '%s', %d, %d, %d, %d, %d, %d)", FrozenHistoryTableName, secInviteId, fmt.Sprintf("%d-%d-%d", uid, inviteId, secInviteId), rewardNum, rewardTokenId, proto.TOKEN_TYPE_OPERATOR_HISTORY_IVITE, now, proto.TOKEN_OPT_TYPE_ADD, newFrozen))
			if err != nil {
				tokenSession.Rollback()
				return errors.NewSys(err)
			}
			if affected, err := result.RowsAffected(); err != nil || affected == 0 {
				tokenSession.Rollback()
				return errors.NewSys(err)
			}
		}
	}

	//提交
	err = tokenSession.Commit()
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
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
func (s *UserToken) SubMoneyWithFronzen(sess *xorm.Session, num int64, ukey string, ty proto.TOKEN_TYPE_OPERATOR) (ret int32, err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"num":         num,
				"uid":         s.Uid,
				"token_id":    s.TokenId,
				"balance":     s.Balance,
				"entrusdt_id": ukey,
				"ty":          ty,
			}).Errorf("sub  money with fronzen error %s", err.Error())
		}
	}()
	var aff int64
	if s.Balance < num {
		ret = ERR_TOKEN_LESS
		return
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
		Ukey:    ukey,
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
		Ukey:    ukey,
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
func (s *UserToken) ReturnFronzen(sess *xorm.Session, num int64, ukey string, ty proto.TOKEN_TYPE_OPERATOR) (err error) {
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
		Ukey:    ukey,
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
		Ukey:    ukey,
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
func (s *UserToken) TransferToCurrency(uid uint64, tokenId int, num int64) error {
	//检查代币是否足够
	err := s.GetUserToken(uid, tokenId)
	if err != nil {
		return err
	}
	if s.Balance < num {
		return errors.NewNormal("余额不足")
	}
	tokenName := s.TokenName

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
	result, err := tokenSession.Exec(fmt.Sprintf("UPDATE %s SET"+
		" balance=%d,"+
		" version=version+1"+
		" WHERE uid=%d"+
		" AND token_id=%d"+
		" AND version=%d", s.TableName(), newBalance, uid, tokenId, s.Version))
	if err != nil {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 { //无影响行数
		tokenSession.Rollback()
		return errors.NewSys(err)
	}

	//2.代币流水
	result, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
		" (uid, token_id, ukey, type, opt, balance, num, created_time, transfer_time) VALUES"+
		" (%d, %d, %d, %d, %d, %d, %d, %d, %d)", new(MoneyRecord).TableName(), uid, tokenId, transferId, proto.TOKEN_TYPE_OPERATOR_TRANSFER_TO_CURRENCY, proto.TOKEN_OPT_TYPE_DEL, newBalance, num, now, now))
	if err != nil {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 { //无影响行数
		tokenSession.Rollback()
		return errors.NewSys(err)
	}

	//3.划转记录
	result, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
		" (id, uid, token_id, token_name, num, states, create_time) VALUES"+
		" (%d, %d, %d, '%s', %d, %d, %d)", new(TransferRecord).TableName(), transferId, uid, tokenId, tokenName, num, 1, now))
	if err != nil {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 { //无影响行数
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
	//判断states为1才更新（消息可能重复发送）
	_, err := engine.Exec(fmt.Sprintf("UPDATE %s SET states=2,update_time=%d,done_time=%d WHERE id=%d AND states=1", new(TransferRecord).TableName(), time.Now().Unix(), msg.DoneTime, msg.Id))
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}

//从法币转入，同一个消息只能处理一次（消息重发机制可能导致同一个消息发送多次）
func (s *UserToken) TransferFromCurrency(msg *proto.TransferToTokenTodoMessage) error {
	//!!!!重要，判断消息是否已处理过
	rdsClient := DB.GetCommonRedisConn()
	isHandled, history, err := new(MoneyRecord).IsTransferFromCurrencyHandled(msg.Id)
	if err != nil {
		return err
	}
	if isHandled { //已处理过，返回done消息并直接退出!
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
	var result sql.Result
	if !has { //代币账号不存在，新建
		result, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
			" (uid, token_id, balance, version)"+
			" VALUES (%d, %d, %d, 1)", s.TableName(), msg.Uid, msg.TokenId, newTokenBalance))
	} else { //代币账号已存在，更新
		result, err = tokenSession.Exec(fmt.Sprintf("UPDATE %s SET balance=%d,version=version+1 WHERE uid=%d AND token_id=%d AND version=%d", s.TableName(), newTokenBalance, msg.Uid, msg.TokenId, userToken.Version))
	}
	if err != nil {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 { //无影响行数
		tokenSession.Rollback()
		return errors.NewSys(err)
	}

	//2.代币流水，再次检查消息是否已处理!!!（判断ukey是否已存在）
	result, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
		" (uid, token_id, ukey, type, opt, balance, num, created_time, transfer_time)"+
		" SELECT %d, %d, '%d', %d, %d, %d, %d, %d, %d"+
		" FROM DUAL"+
		" WHERE NOT EXISTS (SELECT ukey FROM %[1]s WHERE ukey='%[4]d')", new(MoneyRecord).TableName(), msg.Uid, msg.TokenId, msg.Id, proto.TOKEN_TYPE_OPERATOR_TRANSFER_FROM_CURRENCY, proto.TOKEN_OPT_TYPE_ADD, newTokenBalance, msg.Num, now, msg.CreateTime))
	if err != nil {
		tokenSession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 { //无影响行数
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
