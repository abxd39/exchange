package model

import (
	"digicon/currency_service/dao"
	"digicon/currency_service/log"
	"digicon/currency_service/rpc/client"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"time"
)

type UserCurrencyWechatPay struct {
	Uid         uint64 `xorm:"not null pk default 0 comment('用户uid') INT(10)"`
	Name        string `xorm:"not null default '' comment('用户姓名') VARCHAR(20)"`
	Wechat      string `xorm:"not null default '' comment('微信号码') VARCHAR(20)"`
	ReceiptCode string `xorm:"not null default '' comment('收款二维码图片路径') VARCHAR(100)"`
	CreateTime  string `xorm:"not null comment('创建时间') DATETIME"`
	UpdateTime  string `xorm:"not null comment('修改时间') DATETIME"`
}

func (w *UserCurrencyWechatPay) SetWechatPay(req *proto.WeChatPayRequest) (int32, error) {

	/////////////////  1.  验证  验证码 /////////////////////////
	rsp, err := client.InnerService.UserSevice.CallAuthVerify(&proto.AuthVerifyRequest{
		Uid:      req.Uid,
		Code:     req.Verify,
		AuthType: 10, // 设置银行卡支付 10
	})
	if err != nil {
		log.Log.Errorln(err.Error())
		return ERRCODE_SMS_CODE_DIFF, err
	}
	if rsp.Code != ERRCODE_SUCCESS {
		log.Log.Errorln(err.Error())
		return ERRCODE_SMS_CODE_DIFF, err
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
