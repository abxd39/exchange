package model

import (
	"digicon/currency_service/dao"
	_ "digicon/proto/common"
	proto "digicon/proto/rpc"
	"errors"
	"time"
)

type UserCurrencyPaypalPay struct {
	Uid        int    `xorm:"not null pk default 0 comment('用户uid') INT(10)"`
	Paypal     string `xorm:"not null default '' comment('paypal 账号') VARCHAR(20)"`
	CreateTime string `xorm:"not null comment('创建时间') DATETIME"`
	UpdateTime string `xorm:"not null comment('修改时间') DATETIME"`
}

func (pal *UserCurrencyPaypalPay) SetPaypal(req *proto.PaypalRequest) (int32, error) {
	//验证token
	//调用实名接口

	//检查数据库是否存在该条记录
	engine := dao.DB.GetMysqlConn()
	has, err := engine.Exist(&UserCurrencyPaypalPay{
		Uid: int(req.Uid),
	})
	if err != nil {
		return 1, err
	}
	current := time.Now()
	if has {
		//默认修改
		_, err := engine.ID(req.Uid).Update(&UserCurrencyPaypalPay{
			Paypal:     req.Paypal,
			UpdateTime: current.String(),
		})
		if err != nil {
			return 1, errors.New("modify paypal fail!!")
		}
	} else {
		_, err := engine.InsertOne(&UserCurrencyPaypalPay{
			Uid:        int(req.Uid),
			Paypal:     req.Paypal,
			CreateTime: current.String(),
			UpdateTime: current.String(),
		})
		if err != nil {
			return 1, err
		}
	}
	return 0, nil
}
