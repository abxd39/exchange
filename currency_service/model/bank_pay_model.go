package model

import (
	"digicon/currency_service/dao"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"errors"
	"fmt"
	"strings"
	"time"
)

type UserCurrencyBankPay struct {
	Uid        uint64 `xorm:"not null pk default 0 comment('用户uid') INT(10)"`
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
		Uid: req.Uid,
	})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	current := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("SetBankPay%#v\n", req)
	if has {
		return ERRCODE_ACCOUNT_EXIST, errors.New("account already exist!!")
	} else {
		//插入新的纪录

		_, err := engine.InsertOne(&UserCurrencyBankPay{
			Uid:        req.Uid,
			Name:       req.Name,
			CardNum:    req.CardNum,
			BankName:   req.BankName,
			BankInfo:   req.BankInfo,
			CreateTime: current,
			UpdataTime: current,
		})
		if err != nil {
			return ERRCODE_UNKNOWN, err
		}
	}
	//
	return ERRCODE_SUCCESS, nil
}
