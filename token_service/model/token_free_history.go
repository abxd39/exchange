package model

import (
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

type TokenFreeHistory struct {
	Id          int64  `xorm:"pk default 1 BIGINT(18)"`
	TokenId     int    `xorm:"comment('代币ID') INT(11)"`
	Opt         int    `xorm:"comment('操作方向1加2减') INT(11)"`
	Type        int    `xorm:"comment('流水类型1区块 2委托 3注册奖励 4邀请奖励 5撤销委托 6交易入账 7冻结退回 8系统自动退小额余额 9 交易确认扣减冻结数量 10划转到法币 11划转到币币') INT(11)"`
	Num         int64  `xorm:"comment('数量') BIGINT(20)"`
	CreatedTime int64  `xorm:"comment('操作时间') created  BIGINT(20)"`
	Ukey        string `xorm:"comment('联合key') VARCHAR(128)"`
}

/*
type Finance struct {
	Balance int64 `xorm:"default 0 BIGINT(20)"`
	TokenId int   `xorm:"not null pk INT(20)"`
}
*/

func InsertIntoTokenFreeHistory(sess *xorm.Session, g ...*TokenFreeHistory) (err error) {
	_, err = sess.Insert(g)
	if err != nil {
		log.Error(err)
		return
	}
	return nil
}

/*
func GetFinance(token_id int) *Finance  {
	d:=&Finance{}
	ok,err :=DB.GetMysqlConn().Where("token_id=?",token_id).Get(d)
	if err!=nil {
		log.Error(err)
		return nil
	}
	if ok {
		log.Fatalln(errors.New(fmt.Sprintf("please config finance token =%d",token_id)))
		return d
	}
	return nil
}
*/
