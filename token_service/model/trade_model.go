package model

import (
	"digicon/common/model"
	. "digicon/token_service/dao"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

const (
	TRADE_STATES_PART = 1 //部分成交
	TRADE_STATES_ALL  = 2 //全部成交
	TRADE_STATES_DEL  = 3 //撤单
)

type Trade struct {
	TradeId      int    `xorm:"not null pk autoincr comment('交易表的id') INT(11)"`
	TradeNo      string `xorm:"comment('订单号') unique(uni_reade_no) VARCHAR(32)"`
	Uid          uint64 `xorm:"comment('买家uid') index unique(uni_reade_no) BIGINT(11)"`
	TokenId      int    `xorm:"comment('主货币id') index INT(11)"`
	TokenTradeId int    `xorm:"comment('交易币种') INT(11)"`
	TokenName    string `xorm:"comment('交易队') VARCHAR(32)"`
	Price        int64  `xorm:"comment('价格') BIGINT(20)"`
	Num          int64  `xorm:"comment('数量') BIGINT(20)"`
	//Balance      int64  `xorm:"BIGINT(20)"`
	Fee          int64  `xorm:"comment('手续费') BIGINT(20)"`
	Opt          int    `xorm:"comment(' buy  1或sell 2') index TINYINT(4)"`
	DealTime     int64  `xorm:"comment('成交时间') BIGINT(20)"`
	States       int    `xorm:"comment('0是挂单，1是部分成交,2成交， -1撤销') INT(11)"`
}

func (s *Trade) Insert(session *xorm.Session, t ...*Trade) (err error) {
	defer func() {
		if err != nil {
			for _, v := range t {
				log.WithFields(log.Fields{
					"uid":      v.Uid,
					"opt":      v.Opt,
					"token_id": v.TokenId,
					"price":    v.Price,
					"fee":      v.Fee,
					"trade_no": v.TradeNo,
				}).Errorf("inset  money record error %s", err.Error())
			}
		}
	}()
	_, err = session.Insert(t)
	return
}

func (s *Trade) GetUserTradeList(pageIndex, pageSize int, uid uint64) (*model.ModelList, []*Trade, error) {
	engine := DB.GetMysqlConn()

	query := engine.Where("uid=?", uid).Desc("deal_time")
	tempQuery := *query

	count, err := tempQuery.Count(s)
	if err != nil {
		return nil, nil, err
	}

	// 获取分页
	offset, modelList := model.Paging(pageIndex, pageSize, int(count))

	var list []*Trade
	err = query.Select("*").Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, nil, err
	}
	modelList.Items = list

	return modelList, list, nil
}
