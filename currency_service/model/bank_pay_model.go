package model

import (
	"digicon/currency_service/dao"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"errors"
	"strings"
	"time"
)

type UserCurrencyBankPay struct {
	Uid        int    `xorm:"not null pk default 0 comment('用户uid') INT(10)"`
	Name       string `xorm:"not null default '' comment('用户姓名') VARCHAR(20)"`
	CardNum    string `xorm:"not null default '' comment('银行卡号') VARCHAR(20)"`
	BankName   string `xorm:"not null default '' comment('银行名称') VARCHAR(20)"`
	BankInfo   string `xorm:"not null default '' comment('支行名称') VARCHAR(20)"`
	CreateTime string `xorm:"not null comment('创建时间') DATETIME"`
	UpdataTime string `xorm:"not null comment('修改时间') DATETIME"`
}

func (*UserCurrencyBankPay) SetBankPay(req *proto.BankPayRequest) (int32, error) {
	//比较两次输入的银行卡号是否匹配
	if b := strings.Compare(req.CardNum, req.VerifyNum); b != 0 {
		return ERRCODE_ACCOUNT_BANK_CARD_NUMBER_MISMATCH, errors.New("bankcard number with verify bankcard number mismatching")
	}
	//验证token

	//实名
	//银行卡实名 电话与银行卡实名
	//调用接口

	engine := dao.DB.GetMysqlConn()
	//检查数据库是否已经存在uid
	has, err := engine.Exist(&UserCurrencyBankPay{
		Uid: int(req.Uid),
	})
	if err != nil {
		return 1, err
	}
	current := time.Now()
	if has {
		//默认为修改
		_, err := engine.ID(req.Uid).Update(&UserCurrencyBankPay{
			Name:       req.Name,
			CardNum:    req.CardNum,
			BankName:   req.BankName,
			BankInfo:   req.BankInfo,
			UpdataTime: current.String(),
		})
		if err != nil {
			return 0, errors.New("modify bank card number fail!! ")
		}
		//errors.New("bank card number already exist")
	} else {
		//插入新的纪录
		_, err := engine.InsertOne(&UserCurrencyBankPay{
			Uid:        int(req.Uid),
			Name:       req.Name,
			CardNum:    req.CardNum,
			BankName:   req.BankName,
			BankInfo:   req.BankInfo,
			CreateTime: current.String(),
			UpdataTime: current.String(),
		})
		if err != nil {
			return 1, errors.New("Set bank card number fail")
		}
	}
	//
	return 0, nil
}
