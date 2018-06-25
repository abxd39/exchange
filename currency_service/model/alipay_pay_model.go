package model

import (
	"digicon/currency_service/dao"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"errors"
	"time"
)

type UserCurrencyAlipayPay struct {
	Uid         uint64 `xorm:"not null pk default 0 comment('用户uid') INT(10)"`
	Name        string `xorm:"not null default '' comment('用户姓名') VARCHAR(20)"`
	Alipay      string `xorm:"not null default '' comment('支付宝账号') VARCHAR(20)"`
	ReceiptCode string `xorm:"not null default '' comment('支付宝收款二维码图片路径') VARCHAR(100)"`
	CreateTime  string `xorm:"not null comment('创建时间') DATETIME"`
	UpdataTime  string `xorm:"not null comment('修改时间') DATETIME"`
}

func (ali *UserCurrencyAlipayPay) SetAlipay(req *proto.AlipayRequest) (int32, error) {
	//验证token
	//是否需要验证支付宝是否属于该账户
	//查询数据库是否存在
	engine := dao.DB.GetMysqlConn()
	has, err := engine.Exist(&UserCurrencyAlipayPay{
		Uid: req.Uid,
	})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	current := time.Now().Format("2006-01-02 15:04:05")
	if has {
		return ERRCODE_ACCOUNT_EXIST, errors.New("account already exist!!")
	} else {
		_, err := engine.InsertOne(&UserCurrencyAlipayPay{
			Uid:         req.Uid,
			Name:        req.Name,
			Alipay:      req.Alipay,
			ReceiptCode: req.ReceiptCode,
			CreateTime:  current,
			UpdataTime:  current,
		})
		if err != nil {
			return ERRCODE_UNKNOWN, err
		}
	}
	return ERRCODE_SUCCESS, nil
}
