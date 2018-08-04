package cron

import (
	"digicon/common/constant"
	"digicon/currency_service/dao"
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"encoding/json"
)

//划出到代币成功，标记消息状态为已完成
func HandlerTransferToTokenDone() {
	rdsClient := dao.DB.GetRedisConn()
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
