package model

import (
	"digicon/currency_service/dao"
	//log "github.com/sirupsen/logrus"
	"digicon/currency_service/rpc/client"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type UserCurrencyBankPay struct {
	Uid        uint64 `xorm:"not null pk default 0 comment('用户uid') INT(10)"     json:"uid"`
	Name       string `xorm:"not null default '' comment('用户姓名') VARCHAR(20)"  json:"name"`
	CardNum    string `xorm:"not null default '' comment('银行卡号') VARCHAR(20)"  json:"card_num"`
	BankName   string `xorm:"not null default '' comment('银行名称') VARCHAR(20)"  json:"bank_name"` 
	BankInfo   string `xorm:"not null default '' comment('支行名称') VARCHAR(20)"  json:"bank_info"`
	CreateTime string `xorm:"not null comment('创建时间') DATETIME"`
	UpdataTime string `xorm:"not null comment('修改时间') DATETIME"`
}

func (p *UserCurrencyBankPay) SetBankPay(req *proto.BankPayRequest) (int32, error) {
	//比较两次输入的银行卡号是否匹配
	if b := strings.Compare(req.CardNum, req.VerifyNum); b != 0 {
		return ERRCODE_ACCOUNT_BANK_CARD_NUMBER_MISMATCH, errors.New("bankcard number with verify bankcard number mismatching")
	}
	/////////////////  1.  验证  验证码 /////////////////////////
	rsp, err := client.InnerService.UserSevice.CallAuthVerify(&proto.AuthVerifyRequest{
		Uid:      req.Uid,
		Code:     req.Verify,
		AuthType: 7, // 设置银行卡支付 7
	})
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_SMS_CODE_DIFF, err
	}
	if rsp.Code != ERRCODE_SUCCESS {
		log.Errorln(err.Error())
		return ERRCODE_SMS_CODE_DIFF, err
	}

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
	//fmt.Printf("\n\n SetBankPay    %#v \n\n", req)
	if has {
		p.UpdataTime = current
		_, err := engine.Where("uid =?", req.Uid).Update(p)
		if err != nil {
			fmt.Println("update error:", err.Error())
			return ERRCODE_UNKNOWN, err
		}
		//return ERRCODE_ACCOUNT_EXIST, errors.New("account already exist!!")
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
			fmt.Println(err.Error())
			return ERRCODE_UNKNOWN, err
		}
	}
	//
	return ERRCODE_SUCCESS, nil
}

func (p *UserCurrencyBankPay) GetByUid(uid uint64) (err error) {
	engine := dao.DB.GetMysqlConn()
	_, err = engine.Where("uid =?", uid).Get(p)
	return
}
