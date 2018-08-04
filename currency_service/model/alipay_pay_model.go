package model

import (
	"digicon/currency_service/dao"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"time"

	log "github.com/sirupsen/logrus"
	"digicon/currency_service/rpc/client"
)

type UserCurrencyAlipayPay struct {
	Uid         uint64 `xorm:"not null pk default 0 comment('用户uid') INT(10)"       json:"uid"`
	Name        string `xorm:"not null default '' comment('用户姓名') VARCHAR(20)"     json:"name"`
	Alipay      string `xorm:"not null default '' comment('支付宝账号') VARCHAR(20)"   json:"alipay"`
	ReceiptCode string `xorm:"not null default '' comment('支付宝收款二维码图片路径') VARCHAR(100)"  json:"receipt_code"`
	CreateTime  string `xorm:"not null comment('创建时间') DATETIME"`
	UpdataTime  string `xorm:"not null comment('修改时间') DATETIME"`
}

func (ali *UserCurrencyAlipayPay) SetAlipay(req *proto.AlipayRequest) (int32, error) {
	//func (ali *UserCurrencyAlipayPay) SetAlipay(req) (int32, error) {

	//////////////////  1.  验证  验证码 /////////////////////////
	rsp, err := client.InnerService.UserSevice.CallAuthVerify(&proto.AuthVerifyRequest{
		Uid:      req.Uid,
		Code:     req.Verify,
		AuthType: 9, // 设置支付宝支付 9
	})
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_SMS_CODE_DIFF, err
	}
	if rsp.Code != ERRCODE_SUCCESS {
		log.Errorln(err.Error())
		return ERRCODE_SMS_CODE_DIFF, err
	}

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
		ali.UpdataTime = current
		//fmt.Println("ali: ", ali)
		_, err := engine.Where("uid=?", req.Uid).Update(ali)
		if err != nil {
			//fmt.Println("err: ", err.Error())
			return ERRCODE_UNKNOWN, err
		}
	} else {
		ali.CreateTime = current
		ali.UpdataTime = current
		//fmt.Println("ali:", ali)
		_, err := engine.InsertOne(ali)
		if err != nil {
			return ERRCODE_UNKNOWN, err
		}
	}
	return ERRCODE_SUCCESS, nil
}

func (ali *UserCurrencyAlipayPay) GetByUid(uid uint64) (err error) {
	engine := dao.DB.GetMysqlConn()
	_, err = engine.Where("uid =?", uid).Get(ali)
	return
}
