package cron

import (
	"digicon/common/app"
	"digicon/common/constant"
	"digicon/currency_service/dao"
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"encoding/json"
)

//划入法币处理，来源：代币
func HandlerTransferFromToken() {
	rdsClient := dao.DB.GetCommonRedisConn()
	userCurrencyMD := new(model.UserCurrency)

	for {
		//监听app是否退出
		if app.IsAppExit {
			return
		}

		//处理消息
		msgBody, err := rdsClient.LPop(constant.RDS_TOKEN_TO_CURRENCY_TODO).Bytes()
		if err != nil {
			continue
		}

		msg := &proto.TransferToCurrencyTodoMessage{}
		err = json.Unmarshal(msgBody, msg)
		if err != nil {
			continue
		}

		userCurrencyMD.TransferFromToken(msg)
	}
}
