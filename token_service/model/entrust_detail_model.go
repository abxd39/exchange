package model

import (
	proto "digicon/proto/rpc"
	///. "digicon/token_service/dao"
	"digicon/common/model"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

import (
	. "digicon/token_service/dao"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

/*
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
*/

type EntrustDetail struct {
	EntrustId   string  `xorm:"not null pk comment('委托记录表（委托管理）') VARCHAR(64)"`
	Uid         uint64  `xorm:"not null comment('用户id') INT(32)"`
	Symbol      string  `xorm:" VARCHAR(64)"`
	TokenId     int     `xorm:"not null comment('货币id') INT(32)"`
	AllNum      int64   `xorm:"not null comment('总数量') BIGINT(20)"`
	SurplusNum  int64   `xorm:"not null comment('剩余数量') BIGINT(20)"`
	Price       int64   `xorm:"not null comment('实际价格(卖出价格）') BIGINT(20)"`
	Sum         int64   `xorm:"not null comment('委托总额') BIGINT(20)"`
	Opt         int     `xorm:"not null comment('类型 买入单1 卖出单2 ') TINYINT(4)"`
	Type        int     `xorm:"not null comment('类型 市价委托1 还是限价委托2') TINYINT(4)"`
	OnPrice     int64   `xorm:"not null comment('委托价格(挂单价格全价格 卖出价格是扣除手续费的）') BIGINT(20)"`
	FeePercent  float64 `xorm:"not null comment('手续费比例') BIGINT(20)"`
	IsFree      bool    `xorm:"not null comment('免手续费') TINYINT(1)"`
	States      int     `xorm:"not null comment('0是挂单，1是部分成交,2成交， 3撤销') TINYINT(4)"`
	CreatedTime int64   `xorm:"not null comment('添加时间') created BIGINT(20)"`
	TradeNum    int64   `xorm:"not null comment('成交数量')  BIGINT(20)"`
	Version     int     `xorm:"version"`
}

func Insert(sess *xorm.Session, s *EntrustDetail) error {
	_, err := sess.Insert(s)
	if err != nil {
		log.WithFields(log.Fields{
			"entrust_id":  s.EntrustId,
			"uid":         s.Uid,
			"all_num":     s.AllNum,
			"SurplusNum":  s.SurplusNum,
			"opt":         s.Opt,
			"on_price":    s.OnPrice,
			"states":      s.States,
			"create_time": s.CreatedTime,
			//"fee":         s.Fee,
		}).Errorf("%s", err.Error())
		return err
	}
	return nil
}

func (s *EntrustDetail) GetBibiHistory(uid int64, limit, page int, symbol string, opt, states, startTime, endTime int) (*model.ModelList, []*EntrustDetail, error) {
	//m := make([]EntrustDetail, 0)
	var statess []int
	if states == 0 {
		statess = []int{4, 1, 2, 3}
	} else if states == 1 {
		statess = []int{4}
	} else {
		statess = []int{1, 2, 3}
	}
	var optt []int
	if opt == 0 {
		optt = []int{1, 2}
	} else if opt == 1 {
		optt = []int{1}
	} else {
		optt = []int{2}
	}

	engine := DB.GetMysqlConn()
	query := engine.Where("uid = ?", uid)
	if symbol != "" {
		query.Where("symbol = ?", symbol)
	}
	if startTime != 0 {
		query.Where("created_time >= ?", startTime)
	}
	if endTime != 0 {
		query.Where("created_time <= ?", endTime)
	}
	query.In("states", statess)
	query.In("opt", optt)
	//query := engine.Where("symbol = ? and uid=? and created_time > ? and created_time <= ?", symbol,uid,startTime,endTime).In("states",statess).In("opt",optt)

	tempQuery := query.Clone()
	count, err := tempQuery.Count(s)

	if err != nil {
		return nil, nil, err
	}
	// 获取分页
	offset, modelList := model.Paging(page, limit, int(count))

	var list []*EntrustDetail

	err = query.Desc("created_time").Limit(modelList.PageSize, offset).Find(&list)

	if err != nil {
		log.Errorln(err.Error())
		return nil, nil, err
	}
	modelList.Items = list
	return modelList, list, nil
}

func (s *EntrustDetail) GetHistory(uid uint64, limit, page int) []EntrustDetail {
	m := make([]EntrustDetail, 0)
	err := DB.GetMysqlConn().Where("uid=?", uid).Desc("created_time").Limit(limit, page-1).Find(&m)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	return m
}

func (s *EntrustDetail) GetList(uid uint64, limit, page int) []EntrustDetail {
	m := make([]EntrustDetail, 0)
	i := []int{int(proto.TRADE_STATES_TRADE_UN), int(proto.TRADE_STATES_TRADE_PART)}
	err := DB.GetMysqlConn().Where("uid=?", uid).In("states", i).Desc("created_time").Limit(limit, page-1).Find(&m)
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	return m
}

/*
func (s *EntrustDetail) SubSurplusInCache(num int64) {
	s.SurplusNum -= num
	s.TradeNum += num
}
*/

//减少剩余数量
func (s *EntrustDetail) SubSurplus(sess *xorm.Session, deal_num int64) error {

	if s.SurplusNum > deal_num {
		s.States = int(proto.TRADE_STATES_TRADE_PART)
	} else if s.SurplusNum == deal_num {
		s.States = int(proto.TRADE_STATES_TRADE_ALL)
	} else {
		log.WithFields(log.Fields{
			"uid":        s.Uid,
			"entrust_id": s.EntrustId,
			"sulplus":    s.SurplusNum,
			"deal_num":   deal_num,
			"os_id":      os.Getpid(),
		}).Info("decr entrust_detail surplus is less deal_num info")
		return errors.New(fmt.Sprintf("decr entrust_detail surplus is less deal_num %d", deal_num))
	}

	log.WithFields(log.Fields{
		"uid":        s.Uid,
		"entrust_id": s.EntrustId,
		"sulplus":    s.SurplusNum,
		"deal_num":   deal_num,
		"os_id":      os.Getpid(),
	}).Info("just record entrust_detail surplus ")
	s.SurplusNum -= deal_num
	s.TradeNum += deal_num
	//_, err := sess.Where("entrust_id=?", s.EntrustId).Cols("states", "surplus_num", "price").Decr("surplus_num", deal_num).Incr("trade_num", deal_num).Update(s)

	aff, err := sess.Where("entrust_id=?", s.EntrustId).Cols("states", "surplus_num", "trade_num", "price").Update(s)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	if aff == 0 {
		err = errors.New("version is err")
		return err
	}
	return nil
}

func GetEntrust(entrust_id string) *EntrustDetail {
	e := &EntrustDetail{}
	ok, err := DB.GetMysqlConn().Where("entrust_id=?", entrust_id).Get(e)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	if ok {
		return e
	}
	return nil
}
