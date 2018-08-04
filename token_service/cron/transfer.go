package cron

import (
	"digicon/common/constant"
	proto "digicon/proto/rpc"
	"digicon/token_service/dao"
	"digicon/token_service/model"
	"encoding/json"
	"time"
)

//划转到法币成功事件处理
func HandlerTransferToCurrencyDone() {
	rdsClient := dao.DB.GetRedisConn()
	userTokenMD := new(model.UserToken)

	for {
		msgBody, err := rdsClient.LPop(constant.RDS_TOKEN_TO_CURRENCY_DONE).Bytes()
		if err != nil {
			continue
		}

		msg := &proto.TransferToCurrencyDoneMessage{}
		err = json.Unmarshal(msgBody, msg)
		if err != nil {
			continue
		}

		err = userTokenMD.TransferToCurrencyDone(msg)
		if err != nil {
			continue
		}
	}
}

//划转到法币消息发送或处理失败，定时重新发送
func ResendTransferToCurrencyMsg() {
	rdsClient := dao.DB.GetRedisConn()
	transferRecordMD := new(model.TransferRecord)

	for {
		list, err := transferRecordMD.ListOverime()
		if err != nil {
			continue
		}

		for _, v := range list {
			msg, err := json.Marshal(proto.TransferToCurrencyTodoMessage{
				Id:         int64(v.Id),
				Uid:        int32(v.Uid),
				TokenId:    int32(v.TokenId),
				Num:        v.Num,
				CreateTime: v.CreateTime,
			})
			if err != nil {
				continue
			}

			rdsClient.RPush(constant.RDS_TOKEN_TO_CURRENCY_TODO, msg)
		}

		time.Sleep(10 * time.Second)
	}
}
