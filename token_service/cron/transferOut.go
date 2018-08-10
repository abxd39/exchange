package cron

import (
	"digicon/common/constant"
	proto "digicon/proto/rpc"
	"digicon/token_service/dao"
	"digicon/token_service/model"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

//划出到法币成功，标记消息状态为已完成
func HandlerTransferToCurrencyDone() {
	rdsClient := dao.DB.GetCommonRedisConn()
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

//消息重发机制，防止发送失败或远程处理失败导致消息丢失
func ResendTransferToCurrencyMsg() {
	rdsClient := dao.DB.GetCommonRedisConn()
	transferRecordMD := new(model.TransferRecord)

	list, err := transferRecordMD.ListOvertime(10)
	log.Info("划转到法币消息重发，overtime_list：", len(list), ", error：", err)
	if err != nil {
		return
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
}
