package model

import (
	. "digicon/token_service/dao"
	. "digicon/token_service/log"
	"github.com/go-xorm/xorm"
)

const (
	MONEY_UKEY_TYPE_HASH    = 1
	MONEY_UKEY_TYPE_ENTRUST = 2
	MONEY_UKEY_TYPE_TRADE   = 3
)

type MoneyRecord struct {
	Id          uint64 `xorm:"pk autoincr BIGINT(20)"`
	Uid         uint64 `xorm:"comment('用户ID') unique(hash_index)  INT(11)"`
	TokenId     int    `xorm:"comment('代币ID') INT(11)"`
	Ukey        string `xorm:"comment('联合key') unique(hash_index) VARCHAR(128)"`
	Type        int    `xorm:"comment('流水类型1区块2委托') INT(11)"`
	Opt         int    `xorm:"comment('操作方向1加2减') unique(hash_index) TINYINT(4)"`
	Num         int64  `xorm:"comment('数量') BIGINT(20)"`
	Surplus     int64  `xorm:"comment('余额') BIGINT(20)"`
	CreatedTime int64  `xorm:"comment('操作时间')  created BIGINT(20)"`
}

//检查流水记录是否存在
func (s *MoneyRecord) CheckExist(ukey string, ty int32) (ok bool, err error) {
	ok, err = DB.GetMysqlConn().Where("ukey=? and type=?", ukey, ty).Exist(&MoneyRecord{})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

//新增一条流水
func (s *MoneyRecord) InsertRecord(session *xorm.Session, p *MoneyRecord) (err error) {
	_, err = session.InsertOne(p)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}
