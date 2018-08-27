package model

import (
	. "digicon/token_service/dao"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"time"
)

type TokenFreeHistory struct {
	Id          int64  `xorm:"pk default 1 BIGINT(18)"`
	TokenId     int    `xorm:"comment('代币ID') INT(11)"`
	Opt         int    `xorm:"comment('操作方向1加2减') INT(11)"`
	Type        int    `xorm:"comment('流水类型1区块 2委托 3注册奖励 4邀请奖励 5撤销委托 6交易入账 7冻结退回 8系统自动退小额余额 9 交易确认扣减冻结数量 10划转到法币 11划转到币币') INT(11)"`
	Num         int64  `xorm:"comment('数量') BIGINT(20)"`
	CreatedTime int64  `xorm:"comment('操作时间') created  BIGINT(20)"`
	Ukey        string `xorm:"comment('联合key') VARCHAR(128)"`
	//TradeId     int  `xorm:"comment('交易ID') INT(11)"`
}

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

func ReverseInfo(index int) {

	var err error
	//var ok bool
	g := make([]*TokenFreeHistory, 0)
	err = DB.GetMysqlConn().Where("id>? and id<=?", (index-1)*1000, index*1000).Find(&g)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if len(g) > 0 {
		for k, v := range g {

			if v.Id%2 == 1 {
				t := g[k+1]

				if v.Ukey != t.Ukey {
					log.Fatalln("err ")
				}
				trade_num := t.Num

				main_num := v.Num

				t.Num = main_num
				v.Num = trade_num

				_, err = DB.GetMysqlConn().Where("id=?", t.Id).Cols("num").Update(t)
				if err != nil {
					log.Fatalln(err.Error())
				}

				_, err = DB.GetMysqlConn().Where("id=?", v.Id).Cols("num").Update(v)
				if err != nil {
					log.Fatalln(err.Error())
				}
			}

		}
	} else {
		log.Infof("final process time %d ", time.Now().Unix())
		return
	}
	ReverseInfo(index + 1)
}
