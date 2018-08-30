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

func InitCron() {
	//划入
	go HandlerTransferFromCurrency()

	//划出
	go HandlerTransferToCurrencyDone()

	//cron
	if cf.Cfg.MustBool("cron", "run", false) {
		c := cron.New()
		c.AddFunc("0 30 * * * *", ResendTransferToCurrencyMsg) // 每半小时
		c.AddFunc("0 0 1 * * *", ReleaseRegisterReward)        // 每天凌晨1点
		c.Start()
	}
}

// 释放注册奖励
func ReleaseRegisterReward() {
	// 可释放的用户
	type canReleaseUser struct {
		Uid           int32 `xorm:"uid"`
		SurplusReward int64 `xorm:"surplus_reward"` // 剩余可释放数
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
	var totalUser, noAuthUser, releaseSuccess, releaseFail int64
	noAuthUidList := make([]int32, 0)
	releaseFailUidList := make([]int32, 0)

	// 开始释放
	tokenSession := dao.DB.GetMysqlConn().NewSession()
	defer tokenSession.Close()

	for {
		offset := (pageIndex - 1) * pageSize

		// 获取可释放的用户
		var canReleaseUserList []*canReleaseUser
		err := dao.DB.GetMysqlConn().SQL(fmt.Sprintf("SELECT a.uid, a.register_reward_total-IFNULL(b.released_total, 0) surplus_reward"+
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

			// 1.判断用户是否通过二级认证
			has, err := dao.DB.GetCommonMysqlConn().Select("*").Where("uid=?", v.Uid).And("security_auth&4=4").Exist(new(model.OutUser))
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

			// 2.根据用户余额确定释放数量
			releaseNum := 0.001 // 默认释放千分之1
			userToken := model.UserToken{}
			_, err = dao.DB.GetMysqlConn().Select("balance, frozen").Where("uid=?", v.Uid).And("token_id=?", rewardTokenId).Get(&userToken)
			if err != nil {
				log.Errorf("【释放注册奖励】获取用户余额出错, uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				continue
			}
			if userToken.Balance >= 1000*100000000 { // 余额大于等于1000，释放千分之2
				releaseNum = 0.002
			}

			// 3.继续确定释放数量
			releaseNumBy8Bit := convert.Float64ToInt64By8Bit(releaseNum)
			if v.SurplusReward < releaseNumBy8Bit { // 用户剩余可释放数量不足，释放数量改为剩余数量
				releaseNumBy8Bit = v.SurplusReward
			}

			// 4.开始释放
			//事务
			err = tokenSession.Begin()
			if err != nil {
				log.Errorf("【释放注册奖励】开启事务出错，uid: %d, err: ", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				continue
			}

			// 4.1.user_token表
			result, err := tokenSession.Exec(fmt.Sprintf("UPDATE %s SET balance=balance+%d, frozen=frozen-%d WHERE uid=%d AND frozen>=%d LIMIT 1",
				userTokenTable, releaseNumBy8Bit, releaseNumBy8Bit, v.Uid, releaseNumBy8Bit))
			if err != nil {
				log.Errorf("【释放注册奖励】更新user_token表出错，uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				continue
			}
			if affected, err := result.RowsAffected(); err != nil { // 获取影响行数错误
				log.Errorf("【释放注册奖励】更新user_token表影响行数出错，uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				continue
			} else if affected != 1 { // 影响行数必须为1
				log.Errorf("【释放注册奖励】更新user_token表影响行数出错，uid: %d，affected: %d", v.Uid, affected)

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				continue
			}

			// 4.2.frozen_history表
			nowTime := time.Now()
			now := nowTime.Unix()
			uKey := fmt.Sprintf("release_%d_%s", v.Uid, nowTime.Format(xtime.LAYOUT_DATE))
			rType := proto.TOKEN_TYPE_OPERATOR_HISTORY_RELEASE

			_, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
				" (uid, ukey, opt, token_id, num, type, create_time, frozen)"+
				" VALUES"+
				" (%d, '%s', %d, %d, %d, %d, %d, %d)",
				frozenHistoryTable, v.Uid, uKey, proto.TOKEN_OPT_TYPE_DEL, rewardTokenId, releaseNumBy8Bit,
				rType, now, userToken.Frozen-releaseNumBy8Bit))
			if err != nil {
				log.Errorf("【释放注册奖励】插入frozen_history表出错，uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				continue
			}

			// 4.3.money_record表
			_, err = tokenSession.Exec(fmt.Sprintf("INSERT INTO %s"+
				" (uid, token_id, ukey, type, opt, balance, num, created_time)"+
				" VALUES"+
				" (%d, %d, '%s', %d, %d, %d, %d, %d)",
				moneyRecordTable, v.Uid, rewardTokenId, uKey, rType, proto.TOKEN_OPT_TYPE_ADD, userToken.Balance+releaseNumBy8Bit, releaseNumBy8Bit, now))
			if err != nil {
				log.Errorf("【释放注册奖励】插入money_record表出错，uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				continue
			}

			// 4.4.提交事务
			err = tokenSession.Commit()
			if err != nil {
				log.Errorf("【释放注册奖励】提交事务出错，uid: %d, err: %s", v.Uid, err.Error())

				releaseFail++
				releaseFailUidList = append(releaseFailUidList, v.Uid)

				tokenSession.Rollback()
				continue
			}

			releaseSuccess++
		}

		pageIndex++

		time.Sleep(1 * time.Second)
	}

	// 汇总数据s
	log.Errorf("【释放注册奖励汇总】用户总数: %d, 未通过认证人数: %d, 释放成功人数: %d, 释放失败人数: %d, 未通过认证UID: %v, 释放失败UID: %v",
		totalUser, noAuthUser, releaseSuccess, releaseFail, noAuthUidList, releaseFailUidList)
}
