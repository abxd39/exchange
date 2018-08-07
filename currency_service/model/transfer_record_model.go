package model

import (
	"digicon/common/errors"
	"digicon/currency_service/dao"
	"time"
)

type TransferRecord struct {
	Id         int64 `xorm:"id pk" json:"id"'`
	Uid        int32 `xorm:"uid" json:"uid"`
	TokenId    int32 `xorm:"token_id" json:"token_id"`
	Num        int64 `xorm:"num" json:"num"`
	States     int32 `xorm:'states' json:"states"`
	CreateTime int64 `xorm:'create_time' json:"create_time"`
}

func (TransferRecord) TableName() string {
	return "transfer_record"
}

//超时未处理的消息
//每次取1000条
func (t *TransferRecord) ListOvertime(overSeconds int64) ([]*TransferRecord, error) {
	var list []*TransferRecord
	engine := dao.DB.GetMysqlConn()
	err := engine.Where("states=1").And("create_time<?", time.Now().Unix()-overSeconds).Limit(1000).OrderBy("create_time asc").Find(&list)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return list, nil
}
