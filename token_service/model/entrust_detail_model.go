package model

import (
	. "digicon/user_service/dao"
	. "digicon/user_service/log"
	"github.com/sirupsen/logrus"
)

type EntrustDetail struct {
	EntrustId   string `xorm:"not null pk comment('委托记录表（委托管理）') VARCHAR(64)"`
	Uid         int    `xorm:"not null comment('用户id') INT(32)"`
	TokenId     int    `xorm:"not null comment('货币id') INT(32)"`
	AllNum      string `xorm:"not null comment('总数量') DECIMAL(20,8)"`
	SurplusNum  string `xorm:"not null comment('剩余数量') DECIMAL(20,8)"`
	Price       string `xorm:"not null comment('实际价格(卖出价格）') DECIMAL(20,8)"`
	CreatedTime int    `xorm:"not null comment('添加时间') INT(10)"`
	Opt         int    `xorm:"not null comment('类型 卖出单1 还是买入单0') TINYINT(4)"`
	OnPrice     string `xorm:"not null comment('委托价格(挂单价格全价格 卖出价格是扣除手续费的）') DECIMAL(20,8)"`
	Fee         string `xorm:"not null comment('手续费比例') DECIMAL(20,8)"`
	States      int    `xorm:"not null comment('状态0正常1撤单') TINYINT(4)"`
}

func (s *EntrustDetail) Save() error {
	_, err := DB.GetMysqlConn().Insert(s)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"entrust_id": s.EntrustId,
			"uid":     s.Uid,
			"all_num":    s.AllNum,
			"SurplusNum":s.SurplusNum,
			"opt":s.Opt,
			"on_price":s.OnPrice,
			"states":s.States,
			"create_time":s.CreatedTime,
			"fee":s.Fee,
		}).Errorf("%s",err.Error())
		return err
	}
	return nil
}
