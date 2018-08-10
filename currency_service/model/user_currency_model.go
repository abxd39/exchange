package model

import (
	"database/sql"
	"digicon/common/constant"
	"digicon/common/errors"
	"digicon/common/snowflake"
	"digicon/common/xtime"
	"digicon/currency_service/dao"
	proto "digicon/proto/rpc"
	"fmt"
	"github.com/gin-gonic/gin/json"
	log "github.com/sirupsen/logrus"
	"time"
)

// 用户虚拟货币资产表
type UserCurrency struct {
	Id        uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid       uint64 `xorm:"INT(10)"     json:"uid"`                                          // 用户ID
	TokenId   uint32 `xorm:"INT(10)"     json:"token_id"`                                     // 虚拟货币类型
	TokenName string `xorm:"VARCHAR(36)" json:"token_name"`                                   // 虚拟货币名字
	Freeze    int64  `xorm:"BIGINT not null default 0"   json:"freeze"`                       // 冻结
	FreezeCny int64  `xorm:"BIGINT not null default 0"   json:"freeze_cny"`
	Balance   int64  `xorm:"not null default 0 comment('余额') BIGINT"   json:"balance"`        // 余额
	BalanceCny  int64 `xorm:"BIGINT not null default 0"                  json:"balance_cny"`
	Address   string `xorm:"not null default '' comment('充值地址') VARCHAR(255)" json:"address"` // 充值地址
	Version   int64  `xorm:"version"`
}

func (UserCurrency) TableName() string {
	return "user_currency"
}

func (this *UserCurrency) Get(id uint64, uid uint64, token_id uint32) *UserCurrency {

	data := new(UserCurrency)
	var isdata bool
	var err error

	if id > 0 {
		isdata, err = dao.DB.GetMysqlConn().Id(id).Get(data)
	} else {
		isdata, err = dao.DB.GetMysqlConn().Where("uid=? AND token_id=?", uid, token_id).Get(data)
	}

	if err != nil {
		log.Errorln(err.Error())
		return nil
	}

	if !isdata {
		return nil
	}

	return data
}

func (this *UserCurrency) GetUserCurrency(uid uint64, nozero bool) (uCurrenList []UserCurrency, err error) {
	engine := dao.DB.GetMysqlConn()
	if nozero {
		err = engine.Where("uid = ? AND balance > 0 ", uid).Find(&uCurrenList)
	} else {
		err = engine.Where("uid=?", uid).Find(&uCurrenList)
	}
	return
}

func (this *UserCurrency) GetBalance(uid uint64, token_id uint32) (data UserCurrency, err error) {
	//data := new(UserCurrency)
	_, err = dao.DB.GetMysqlConn().Where("uid=? AND token_id=?", uid, token_id).Get(&data)
	return

}

func (this *UserCurrency) SetBalance(uid uint64, token_id uint32, amount int64) (err error) {
	engine := dao.DB.GetMysqlConn()
	sql := "UPDATE user_currency SET   balance= balance + ?, version = version + 1 WHERE uid = ? AND token_id = ? AND version = ?"
	sqlRest, err := engine.Exec(sql, amount, uid, token_id, this.Version)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	if rst, _ := sqlRest.RowsAffected(); rst == 0 {
		log.Errorln("添加余额失败")
		err = errors.New("添加余额失败!")
		return err
	}
	return
}

func (this *UserCurrency) TransferToToken(uid uint64, tokenId int, num int64) error {
	//检查法币是否足够
	userCurrency, err := this.GetBalance(uid, uint32(tokenId))
	if err != nil {
		return err
	}
	if userCurrency.Balance < num {
		return errors.NewNormal("余额不足")
	}

	//整理数据
	transferId := snowflake.SnowflakeNode.Generate()
	now := time.Now().Unix()

	//开始划转
	currencySession := dao.DB.GetMysqlConn().NewSession()
	defer currencySession.Close()

	//事务
	err = currencySession.Begin()
	if err != nil {
		return errors.NewSys(err)
	}

	//1.法币账户
	newBalance := userCurrency.Balance - num
	result, err := currencySession.Exec(fmt.Sprintf("UPDATE %s SET"+
		" balance=%d,"+
		" version=version+1"+
		" WHERE uid=%d"+
		" AND token_id=%d"+
		" AND version=%d", this.TableName(), newBalance, uid, tokenId, userCurrency.Version))
	if err != nil {
		currencySession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 { //无影响行数
		currencySession.Rollback()
		return errors.NewSys(err)
	}

	//2.法币流水
	result, err = currencySession.Exec(fmt.Sprintf("INSERT INTO %s"+
		" (uid, trade_uid, order_id, token_id, num, surplus, operator, created_time)"+
		" VALUES (%d, %[2]d, '%d', %d, %d, %d, %d, '%s')", new(UserCurrencyHistory).TableName(), uid, transferId, tokenId, num, newBalance, 4, xtime.Unix2Date(now, xtime.LAYOUT_DATE_TIME)))
	if err != nil {
		currencySession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 { //无影响行数
		currencySession.Rollback()
		return errors.NewSys(err)
	}

	//3.划转记录
	result, err = currencySession.Exec(fmt.Sprintf("INSERT INTO %s"+
		" (id, uid, token_id, token_name, num, states, create_time) VALUES"+
		" (%d, %d, %d, '%s', %d, %d, %d)", new(TransferRecord).TableName(), transferId, uid, tokenId, userCurrency.TokenName, num, 1, now))
	if err != nil {
		currencySession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 { //无影响行数
		currencySession.Rollback()
		return errors.NewSys(err)
	}

	//提交
	err = currencySession.Commit()
	if err != nil {
		return errors.NewSys(err)
	}

	// 发送划转消息给token服
	go func() {
		msg, err := json.Marshal(proto.TransferToTokenTodoMessage{
			Id:         int64(transferId),
			Uid:        int32(uid),
			TokenId:    int32(tokenId),
			Num:        num,
			CreateTime: now,
		})
		if err != nil {
			return
		}

		rdsClient := dao.DB.GetCommonRedisConn()
		rdsClient.RPush(constant.RDS_CURRENCY_TO_TOKEN_TODO, msg)
	}()

	return nil
}

//从代币转入，同一个消息只能处理一次（消息重发机制可能导致同一个消息发送多次）
func (this *UserCurrency) TransferFromToken(msg *proto.TransferToCurrencyTodoMessage) error {
	//!!!!重要，判断消息是否已处理过
	rdsClient := dao.DB.GetCommonRedisConn()
	isHandled, history, err := new(UserCurrencyHistory).IsTransferFromTokenHandled(msg.Id)
	if err != nil {
		return err
	}
	if isHandled { //已处理过，返回done消息并直接退出!
		msg, err := json.Marshal(proto.TransferToCurrencyDoneMessage{
			Id:       msg.Id,
			DoneTime: xtime.Date2Unix(history.CreatedTime, xtime.LAYOUT_DATE_TIME),
		})
		if err != nil {
			return errors.NewSys(err)
		}

		rdsClient.RPush(constant.RDS_TOKEN_TO_CURRENCY_DONE, msg)
		return nil
	}

	//整理数据
	now := time.Now().Unix()

	//开始处理
	currencySession := dao.DB.GetMysqlConn().NewSession()
	sessionClone := currencySession.Clone()
	defer func() {
		currencySession.Close()
		sessionClone.Close()
	}()

	//判断用户法币账户是否存在
	userCurrency := UserCurrency{}
	has, err := sessionClone.Where("uid=?", msg.Uid).And("token_id=?", msg.TokenId).Get(&userCurrency)
	if err != nil {
		return errors.NewSys(err)
	}
	newCurrBalance := userCurrency.Balance + msg.Num

	// 开始事务
	err = currencySession.Begin()
	if err != nil {
		return errors.NewSys(err)
	}

	//1.法币账号
	var result sql.Result
	if !has { //法币账户不存在，新建
		//获取token_name
		tokensModel := new(CommonTokens)
		token := tokensModel.Get(uint32(msg.TokenId), "")
		if token == nil {
			currencySession.Rollback()
			return errors.NewNormal("获取token信息失败")
		}

		//新建法币账户
		result, err = currencySession.Exec(fmt.Sprintf("INSERT INTO %s"+
			" (uid, token_id, token_name, balance, version)"+
			" VALUES (%d, %d, '%s', %d, 1)", this.TableName(), msg.Uid, msg.TokenId, token.Mark, newCurrBalance))
	} else { //法币账户已存在，更新
		result, err = currencySession.Exec(fmt.Sprintf("UPDATE %s SET balance=%d,version=version+1 WHERE uid=%d AND token_id=%d AND version=%d", this.TableName(), newCurrBalance, msg.Uid, msg.TokenId, userCurrency.Version))
	}
	if err != nil {
		currencySession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 { //无影响行数
		currencySession.Rollback()
		return errors.NewSys(err)
	}

	//2.法币流水，再次检查消息是否已处理!!!（order_id是否已存在）
	result, err = currencySession.Exec(fmt.Sprintf("INSERT INTO %s"+
		" (uid, trade_uid, order_id, token_id, num, surplus, operator, created_time)"+
		" SELECT %d, %[2]d, '%d', %d, %d, %d, %d, '%s'"+
		" FROM DUAL"+
		" WHERE NOT EXISTS (SELECT order_id FROM %[1]s WHERE order_id='%[3]d')", new(UserCurrencyHistory).TableName(), msg.Uid, msg.Id, msg.TokenId, msg.Num, newCurrBalance, 3, xtime.Unix2Date(now, xtime.LAYOUT_DATE_TIME)))
	if err != nil {
		currencySession.Rollback()
		return errors.NewSys(err)
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 { //无影响行数
		currencySession.Rollback()
		return errors.NewSys(err)
	}

	//提交
	err = currencySession.Commit()
	if err != nil {
		return errors.NewSys(err)
	}

	//发送处理成功消息给token服
	go func() {
		doneMsg, err := json.Marshal(proto.TransferToCurrencyDoneMessage{
			Id:       msg.Id,
			DoneTime: now,
		})
		if err != nil {
			return
		}
		rdsClient.RPush(constant.RDS_TOKEN_TO_CURRENCY_DONE, doneMsg)
	}()

	return nil
}

//划转到代币成功，把划转记录标记为已完成
func (this *UserCurrency) TransferToTokenDone(msg *proto.TransferToTokenDoneMessage) error {
	engine := dao.DB.GetMysqlConn()
	//判断states为1才更新（消息可能重复发送）
	_, err := engine.Exec(fmt.Sprintf("UPDATE %s SET states=2,update_time=%d,done_time=%d WHERE id=%d AND states=1", new(TransferRecord).TableName(), time.Now().Unix(), msg.DoneTime, msg.Id))
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}
