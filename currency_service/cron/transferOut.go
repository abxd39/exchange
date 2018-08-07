package cron

import (
	"digicon/common/constant"
	"digicon/currency_service/dao"
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

//划出到代币成功，标记消息状态为已完成
func HandlerTransferToTokenDone() {
	rdsClient := dao.DB.GetCommonRedisConn()
	userCurrencyMD := new(model.UserCurrency)

	for {
		msgBody, err := rdsClient.LPop(constant.RDS_CURRENCY_TO_TOKEN_DONE).Bytes()
		if err != nil {
			continue
		}

		msg := &proto.TransferToTokenDoneMessage{}
		err = json.Unmarshal(msgBody, msg)
		if err != nil {
			continue
		}

		err = userCurrencyMD.TransferToTokenDone(msg)
		if err != nil {
			continue
		}
	}
}

//消息重发机制，防止发送失败或远程处理失败导致消息丢失
func ResendTransferToTokenMsg() {
	rdsClient := dao.DB.GetCommonRedisConn()
	transferRecordMD := new(model.TransferRecord)
	var overSeconds int64 = 10

	for {
		list, err := transferRecordMD.ListOvertime(overSeconds)
		log.Info("划转到代币消息重发，overtime_list：", len(list), ", error：", err)
		if err != nil {
			continue
		}

		for _, v := range list {
			log.Info("划转到代币消息重发，msg：", v, ", redis_list_name：", constant.RDS_CURRENCY_TO_TOKEN_TODO)
			msg, err := json.Marshal(proto.TransferToTokenTodoMessage{
				Id:         int64(v.Id),
				Uid:        int32(v.Uid),
				TokenId:    int32(v.TokenId),
				Num:        v.Num,
				CreateTime: v.CreateTime,
			})
			if err != nil {
				continue
			}

			cmd := rdsClient.RPush(constant.RDS_CURRENCY_TO_TOKEN_TODO, msg)
			_, err = cmd.Result()
			log.Info("划转到代币消息重发", err, cmd.Err())
		}

		time.Sleep(time.Duration(overSeconds) * time.Second)
	}
}
