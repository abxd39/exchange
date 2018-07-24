package model

import (
	proto "digicon/proto/rpc"
	///. "digicon/token_service/dao"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

import (
	. "digicon/token_service/dao"
)

type EntrustData struct {
	EntrustId  string `xorm:"not null pk comment('委托记录表（委托管理）') VARCHAR(64)"`
	Uid        uint64 `xorm:"not null comment('用户id') INT(32)"`
	SurplusNum int64  `xorm:"not null comment('剩余数量') BIGINT(20)"`
	Opt        proto.ENTRUST_OPT
	Type       proto.ENTRUST_TYPE
	OnPrice    int64 `xorm:"not null comment('委托价格(挂单价格全价格 卖出价格是扣除手续费的）') BIGINT(20)"`
	Fee        int64 `xorm:"not null comment('手续费比例') BIGINT(20)"`
	States     int   `xorm:"not null comment('状态0正常1撤单2成交') TINYINT(4)"`
}

type EntrustDetail struct {
	EntrustId   string `xorm:"not null pk comment('委托记录表（委托管理）') VARCHAR(64)"`
	Uid         uint64 `xorm:"not null comment('用户id') INT(32)"`
	Symbol      string `xorm:" VARCHAR(64)"`
	TokenId     int    `xorm:"not null comment('货币id') INT(32)"`
	AllNum      int64  `xorm:"not null comment('总数量') BIGINT(20)"`
	SurplusNum  int64  `xorm:"not null comment('剩余数量') BIGINT(20)"`
	Price       int64  `xorm:"not null comment('实际价格(卖出价格）') BIGINT(20)"`
	Mount       int64  `xorm:"not null comment('全部实际价值') BIGINT(20)"`
	Opt         int    `xorm:"not null comment('类型 买入单1 卖出单2 ') TINYINT(4)"`
	Type        int    `xorm:"not null comment('类型 市价委托1 还是限价委托2') TINYINT(4)"`
	OnPrice     int64  `xorm:"not null comment('委托价格(挂单价格全价格 卖出价格是扣除手续费的）') BIGINT(20)"`
	Fee         int64  `xorm:"not null comment('手续费比例') BIGINT(20)"`
	States      int    `xorm:"not null comment('0是挂单，1是部分成交,2成交， 3撤销') TINYINT(4)"`
	CreatedTime int64  `xorm:"not null comment('添加时间') created BIGINT(20)"`
}

func (s *EntrustDetail) Insert(sess *xorm.Session) error {
	_, err := sess.Insert(s)
	if err != nil {
		log.WithFields(logrus.Fields{
			"entrust_id":  s.EntrustId,
			"uid":         s.Uid,
			"all_num":     s.AllNum,
			"SurplusNum":  s.SurplusNum,
			"opt":         s.Opt,
			"on_price":    s.OnPrice,
			"states":      s.States,
			"create_time": s.CreatedTime,
			"fee":         s.Fee,
		}).Errorf("%s", err.Error())
		sess.Rollback()
		return err
	}
	return nil
}

func (s *EntrustDetail) GetHistory(uid uint64, limit, page int) []EntrustDetail {
	m := make([]EntrustDetail, 0)
	err := DB.GetMysqlConn().Where("uid=?", uid).Limit(limit, page-1).Find(&m)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	return m
}

func (s *EntrustDetail) GetList(uid uint64, limit, page int) []EntrustDetail {
	m := make([]EntrustDetail, 0)
	i := []int{0, 1}
	err := DB.GetMysqlConn().Where("uid=?", uid).In("states", i).Limit(limit, page-1).Find(&m)
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	return m
}

func (s *EntrustDetail) UpdateStates(sess *xorm.Session, entrust_id string, states int, deal_num int64) error {

	_, err := sess.Where("entrust_id=?", entrust_id).Cols("states", "surplus_num").Decr("surplus_num", deal_num).Update(&EntrustDetail{States: states})
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}
