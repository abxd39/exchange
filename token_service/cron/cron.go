package cron

import (
	"digicon/common/convert"
	cf "digicon/token_service/conf"
	"digicon/token_service/dao"
	"digicon/token_service/model"
	"fmt"
	"github.com/robfig/cron"
	"time"

	"digicon/common/xtime"
	proto "digicon/proto/rpc"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/log"
)

var CronInstance *cron.Cron

func InitCron() {
	//划入
	go HandlerTransferFromCurrency()

	//划出
	go HandlerTransferToCurrencyDone()

	//cron
	if cf.Cfg.MustBool("cron", "run", false) {
		CronInstance = cron.New()
		CronInstance.AddFunc("0 30 * * * *", ResendTransferToCurrencyMsg) // 每半小时
		CronInstance.AddFunc("0 0 0 * * *", ReleaseRegisterReward)        // 每天凌晨0点
		CronInstance.Start()
	}
}

// 释放注册奖励
func ReleaseRegisterReward() {
	// 可释放的用户
	type canReleaseUser struct {
		Uid                 int32 `xorm:"uid"`
		RegisterRewardTotal int64 `xorm:"register_reward_total"` // 总奖励数
		SurplusReward       int64 `xorm:"surplus_reward"`        // 剩余可释放数
	}

	// 整理数据
	rewardTokenId := 4
	pageIndex := 1
	pageSize := 200

	// 表名
	userTokenTable := new(model.UserToken).TableName()
	frozenHistoryTable := new(model.FrozenHistory).TableName()
	moneyRecordTable := new(model.MoneyRecord).TableName()

	// 汇总数据
	var totalUser, noAuthUser, notNormaUser, releaseSuccess, releaseFail int64
	noAuthUidList := make([]int32, 0)
	notNormalUidList := make([]int32, 0)
	releaseFailUidList := make([]int32, 0)

	// 开始释放
	for {
		offset := (pageIndex - 1) * pageSize

		// 获取可释放的用户
		var canReleaseUserList []*canReleaseUser
		err := dao.DB.GetMysqlConn().SQL(fmt.Sprintf("SELECT a.uid, a.register_reward_total, a.register_reward_total-IFNULL(b.released_total, 0) surplus_reward"+
			" FROM (SELECT SUM(num) register_reward_total, uid FROM %s WHERE token_id=%d AND type IN (%d, %d) GROUP BY uid) a"+ // 注册奖励总数
			" LEFT JOIN (SELECT SUM(num) released_total, uid FROM %s WHERE token_id=%d AND type=%d GROUP BY uid) b ON a.uid=b.uid"+ // 已释放总数
			" WHERE IFNULL(b.released_total, 0)<a.register_reward_total"+ // 已释放总数 小于 注册奖励总数
			" ORDER BY a.uid ASC LIMIT %d, %d", frozenHistoryTable, rewardTokenId, proto.TOKEN_TYPE_OPERATOR_HISTORY_REGISTER, proto.TOKEN_TYPE_OPERATOR_HISTORY_IVITE, frozenHistoryTable, rewardTokenId, proto.TOKEN_TYPE_OPERATOR_HISTORY_RELEASE, offset, pageSize)).
			Find(&canReleaseUserList)
		if err != nil {
			log.Error("【释放注册奖励】获取可释放的用户出错, err:", err.Error())
			break
		}
		if len(canReleaseUserList) == 0 { // !!!!无数据，表示分页结束，退出
			break
		}

		// 循环列表进行处理
		for _, v := range canReleaseUserList {
			totalUser++

			// 判断用户状态、是否通过二级认证
			user := model.OutUser{}
			has, err := dao.DB.GetCommonMysqlConn().Select("*").Where("uid=?", v.Uid).And("security_auth&4=4").Get(&user)
			if err != nil {
				log.Errorf("【释放注册奖励】判断用户二级认证出错, uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				continue
			}
			if !has { // 未通过二级认证
				noAuthUser++
				noAuthUidList = append(noAuthUidList, v.Uid)

				continue
			}
			if user.Status != 1 { // 状态不正常
				notNormaUser++
				notNormalUidList = append(notNormalUidList, v.Uid)

				continue
			}

			// 开始释放
			//事务
			tokenSession := dao.DB.GetMysqlConn().NewSession()
			err = tokenSession.Begin()
			if err != nil {
				log.Errorf("【释放注册奖励】开启事务出错，uid: %d, err: ", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Close()
				continue
			}

			// 1. 获取释放数量
			releaseNum := convert.Int64MulFloat64(v.RegisterRewardTotal, 0.001) // 默认释放注册总奖励数的千分之1

			// 判断是否加速
			userToken := model.UserToken{}
			_, err = tokenSession.SQL(fmt.Sprintf("SELECT balance, frozen FROM %s WHERE uid=%d AND token_id=%d FOR UPDATE", userTokenTable, v.Uid, rewardTokenId)).Get(&userToken)
			if err != nil {
				log.Errorf("【释放注册奖励】获取用户余额出错, uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				tokenSession.Close()
				continue
			}
			if userToken.Balance >= 1000*100000000 { // 余额大于等于1000，加速释放，多释放余额的千分之1
				releaseNum += convert.Int64MulFloat64(userToken.Balance, 0.001)
			}

			// 确定最终释放数量
			if releaseNum > v.SurplusReward { // 释放数量以可释放数量为准
				releaseNum = v.SurplusReward
			}

			// 2. 写表
			// 2.1.user_token表
			result, err := tokenSession.Exec(fmt.Sprintf("UPDATE %s SET balance=balance+%d, frozen=frozen-%d WHERE uid=%d AND token_id=%d AND frozen>=%d",
				userTokenTable, releaseNum, releaseNum, v.Uid, rewardTokenId, releaseNum))
			if err != nil {
				log.Errorf("【释放注册奖励】更新user_token表出错，uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				tokenSession.Close()
				continue
			}

			// 判断影响行数
			if affected, err := result.RowsAffected(); err != nil { // 获取影响行数错误
				log.Errorf("【释放注册奖励】更新user_token表影响行数出错，uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				tokenSession.Close()
				continue
			} else if affected != 1 { // 影响行数必须为1，为0或大于1均出错
				log.Errorf("【释放注册奖励】更新user_token表影响行数出错，uid: %d，affected: %d", v.Uid, affected)

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				tokenSession.Close()
				continue
			}

			// 2.2.frozen_history表
			nowTime := time.Now()
			now := nowTime.Unix()
			uKey := fmt.Sprintf("release_%d_%s", v.Uid, nowTime.Format(xtime.LAYOUT_DATE))
			rType := proto.TOKEN_TYPE_OPERATOR_HISTORY_RELEASE

			_, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
				" (uid, ukey, opt, token_id, num, type, create_time, frozen)"+
				" VALUES"+
				" (%d, '%s', %d, %d, %d, %d, %d, %d)",
				frozenHistoryTable, v.Uid, uKey, proto.TOKEN_OPT_TYPE_DEL, rewardTokenId, releaseNum,
				rType, now, userToken.Frozen-releaseNum))
			if err != nil {
				log.Errorf("【释放注册奖励】插入frozen_history表出错，uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				tokenSession.Close()
				continue
			}

			// 2.3.money_record表
			_, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
				" (uid, token_id, ukey, type, opt, balance, num, created_time)"+
				" VALUES"+
				" (%d, %d, '%s', %d, %d, %d, %d, %d)",
				moneyRecordTable, v.Uid, rewardTokenId, uKey, rType, proto.TOKEN_OPT_TYPE_ADD, userToken.Balance+releaseNum, releaseNum, now))
			if err != nil {
				log.Errorf("【释放注册奖励】插入money_record表出错，uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				tokenSession.Close()
				continue
			}

			// 2.4.提交事务
			err = tokenSession.Commit()
			if err != nil {
				log.Errorf("【释放注册奖励】提交事务出错，uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				tokenSession.Close()
				continue
			}

			tokenSession.Close()
			releaseSuccess++
		}

		pageIndex++

		time.Sleep(1 * time.Second)
	}

	// 汇总数据
	log.Errorf("【释放注册奖励汇总】用户总数: %d, 未通过认证人数: %d, 非正常人数: %d, 释放成功人数: %d, 释放失败人数: %d, 未通过认证UID: %v, 非正常UID: %v, 释放失败UID: %v",
		totalUser, noAuthUser, notNormaUser, releaseSuccess, releaseFail, noAuthUidList, notNormalUidList, releaseFailUidList)
}
