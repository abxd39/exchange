package cron

import (
	"digicon/common/constant"
	proto "digicon/proto/rpc"
	"digicon/token_service/dao"
	"digicon/token_service/model"
	"encoding/json"
)

//划入代币处理，来源：法币
func HandlerTransferFromCurrency() {
	rdsClient := dao.DB.GetCommonRedisConn()
	userTokenMD := new(model.UserToken)

	for {
		msgBody, err := rdsClient.LPop(constant.RDS_CURRENCY_TO_TOKEN_TODO).Bytes()
		if err != nil {
			continue
		}

		msg := &proto.TransferToTokenTodoMessage{}
		err = json.Unmarshal(msgBody, msg)
		if err != nil {
			continue
		}

		userTokenMD.TransferFromCurrency(msg)
	}
}
