package model

import (
	"digicon/currency_service/dao"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"time"
	"digicon/currency_service/rpc/client"
	log "github.com/sirupsen/logrus"
	"fmt"
	. "digicon/common/constant"
)




type UserCurrencyWechatPay struct {
	Uid         uint64 `xorm:"not null pk default 0 comment('用户uid') INT(10)"     json:"uid"`
	Name        string `xorm:"not null default '' comment('用户姓名') VARCHAR(20)"  json:"name"`
	Wechat      string `xorm:"not null default '' comment('微信号码') VARCHAR(20)"  json:"wechat"`
	ReceiptCode string `xorm:"not null default '' comment('收款二维码图片路径') VARCHAR(100)"   json:"receipt_code"`
	CreateTime  string `xorm:"not null comment('创建时间') DATETIME"`
	UpdateTime  string `xorm:"not null comment('修改时间') DATETIME"`
}

func (w *UserCurrencyWechatPay) SetWechatPay(req *proto.WeChatPayRequest) (int32, error) {

	/////////////////  1.  验证  验证码 /////////////////////////
	fmt.Println("req:", req)

	rsp, err := client.InnerService.UserSevice.CallAuthVerify(&proto.AuthVerifyRequest{
		Uid:      req.Uid,
		Code:     req.Verify,
		AuthType: SMS_WECHAT_PAY, // 设置银行卡支付 8
	})
	fmt.Println(rsp, err)
	if err != nil {
		log.Println(rsp)
		log.Errorln(err)
		return ERRCODE_SMS_CODE_DIFF, err
	}


	if rsp != nil && rsp.Code != ERRCODE_SUCCESS {
		log.Println(rsp)
		return ERRCODE_SMS_CODE_DIFF, nil
	}



	//是否需要验证微信是否属于该账户
	//查询数据库是否存在
	engine := dao.DB.GetMysqlConn()
	has, err := engine.Exist(&UserCurrencyWechatPay{
		Uid: req.Uid,
	})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	current := time.Now().Format("2006-01-02 15:04:05")

	if has {
		w.UpdateTime = current
		_, err := engine.Where("uid=?", req.Uid).Update(w)
		if err != nil {
			return ERRCODE_UNKNOWN, err
		}
	} else {
		_, err := engine.InsertOne(&UserCurrencyWechatPay{
			Uid:         req.Uid,
			Name:        req.Name,
			Wechat:      req.Wechat,
			ReceiptCode: req.ReceiptCode,
			CreateTime:  current,
			UpdateTime:  current,
		})
		if err != nil {
			return ERRCODE_UNKNOWN, err
		}
	}
	return ERRCODE_SUCCESS, nil
}

func (w *UserCurrencyWechatPay) GetByUid(uid uint64) (err error) {
	engine := dao.DB.GetMysqlConn()
	_, err = engine.Where("uid =?", uid).Get(w)
	return
}
