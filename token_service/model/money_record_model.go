package model

import (
	. "digicon/token_service/dao"
	. "digicon/token_service/log"
	"github.com/go-xorm/xorm"
)

const (
	MONEY_UKEY_TYPE_HASH    = 1
	MONEY_UKEY_TYPE_ENTRUST = 2
)

type MoneyRecord struct {
	Id          int64  `xorm:"pk autoincr BIGINT(20)"`
	Uid         int    `xorm:"comment('用户ID') INT(11)"`
	TokenId     int    `xorm:"comment('代币ID') INT(11)"`
	Ukey        string `xorm:"comment('联合key') unique VARCHAR(128)"`
	Type        int    `xorm:"comment('流水类型1区块2委托') INT(11)"`
	Opt         int    `xorm:"comment('操作方向1加2减') TINYINT(4)"`
	Num         int64  `xorm:"comment('数量') BIGINT(20)"`
	CreatedTime int64  `xorm:"comment('操作时间') BIGINT(20)"`
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
