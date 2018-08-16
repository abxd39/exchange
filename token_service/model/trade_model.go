package model

import (
	"digicon/common/convert"
	"digicon/common/model"
	. "digicon/token_service/dao"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/liudng/godump"
	log "github.com/sirupsen/logrus"
	"time"
)

/*
const (
	TRADE_STATES_PART = 1 //部分成交
	TRADE_STATES_ALL  = 2 //全部成交
	TRADE_STATES_DEL  = 3 //撤单
)
*/
type Trade struct {
	TradeId          int    `xorm:"not null pk autoincr comment('交易表的id') INT(11)"`
	TradeNo          string `xorm:"comment('订单号') unique(uni_reade_no) VARCHAR(32)"`
	Uid              uint64 `xorm:"comment('买家uid') index unique(uni_reade_no) BIGINT(11)"`
	TokenId          int    `xorm:"comment('主货币id') index INT(11)"`
	TokenTradeId     int    `xorm:"comment('交易币种') INT(11)"`
	TokenAdmissionId int    `xorm:"comment('入账token_id') index INT(11)"`
	Symbol           string `xorm:"comment('交易队') VARCHAR(32)"`
	Price            int64  `xorm:"comment('价格') BIGINT(20)"`
	Num              int64  `xorm:"comment('数量') BIGINT(20)"`
	//Balance      int64  `xorm:"BIGINT(20)"`
	EntrustId string `xorm:"comment('委托ID')  VARCHAR(32)"`
	Fee       int64  `xorm:"comment('手续费数量') BIGINT(20)"`
	Opt       int    `xorm:"comment(' buy  1或sell 2') index TINYINT(4)"`
	DealTime  int64  `xorm:"comment('成交时间') BIGINT(20)"`
	//States    int    `xorm:"comment('0是挂单，1是部分成交,2成交， -1撤销') INT(11)"`
	FeeCny   int64 `xorm:"comment('手续费人民币') BIGINT(20)"`
	TotalCny int64 `xorm:"comment('成交额人民币') BIGINT(20)"`
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

	query := engine.Where("uid=?", uid).Desc("deal_time", "trade_id")
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

func GetUserTradeByEntrustId(entrust_id string) (g []*Trade, err error) {
	g = make([]*Trade, 0)
	err = DB.GetMysqlConn().Where("entrust_id=?", entrust_id).Find(&g)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func CaluateAvgPrice(t []*Trade) int64 {
	var amount, sum int64

	for _, v := range t {
		amount += convert.Int64MulInt64By8Bit(v.Num, v.Price)
		sum += v.Num
	}
	tt := convert.Int64DivInt64By8Bit(amount, sum)

	godump.Dump(tt)
	return convert.Int64DivInt64By8Bit(amount, sum)
}

func Test2(beid, endid int64) {
	log.Infof("begin id=%d,endid=%d", beid, endid)
	g := make([]*Trade, 0)
	err := DB.GetMysqlConn().Where("trade_id>=? and trade_id<?", beid, endid).Find(&g)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	if len(g) == 0 {
		return
	}

	for _, v := range g {
		log.Infof("process id %d", v.TradeId)
		if v.Opt == 1 {
			v.TokenAdmissionId = v.TokenTradeId
			_, err = DB.GetMysqlConn().Where("trade_id=?", v.TradeId).Cols("token_admission_id").Update(v)
			if err != nil {
				log.Fatalln(err.Error())
				return
			}
		} else {
			v.TokenAdmissionId = v.TokenId
			_, err = DB.GetMysqlConn().Where("trade_id=?", v.TradeId).Cols("token_admission_id").Update(v)
			if err != nil {
				log.Fatalln(err.Error())
				return
			}
		}
	}

	Test2(beid+1000, endid+1000)
	/*
		lastt := stime + 43200

		log.Infof("beigin %d",time.Now().Unix())
		_,err:=DB.GetMysqlConn().Exec("call statisticss_daily_fee()")
		if err!=nil {
			log.Fatalln(err.Error())
		}

		log.Infof("end %d",time.Now().Unix())

	*/
}

func Test3(begin, end int64) {
	//g:=make([]*Trade,0)
	//buy
	sql := fmt.Sprintf("select sum(num) as a,sum(fee) as b ,sum(fee_cny) as c ,sum(total_cny) as d,token_admission_id  from trade where deal_time>=%d and deal_time<%d  and opt=1 group by token_admission_id", begin, end)
	r, err := DB.GetMysqlConn().Query(sql)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	l := make(map[int]*TokenDailySheet)

	if len(r) > 0 {
		for _, v := range r {
			h:=&TokenDailySheet{}
			t, ok := v["token_admission_id"]
			if !ok {
				log.Fatal("ok u")
			}

			a, ok := v["a"]
			if !ok {
				log.Fatal("ok a")
			}
			b, ok := v["b"]
			if !ok {
				log.Fatal("ok b")
			}
			c, ok := v["c"]
			if !ok {
				log.Fatal("ok c")
			}
			d, ok := v["d"]
			if !ok {
				log.Fatal("ok d")
			}

			h.TokenId  = convert.BytesToIntAscii(t)
			h.BuyTotal= convert.BytesToInt64Ascii(a)
			h.FeeBuyTotal= convert.BytesToInt64Ascii(b)
			h.FeeBuyCny = convert.BytesToInt64Ascii(c)
			h.BuyTotalCny= convert.BytesToInt64Ascii(d)

			l[h.TokenId] = h
		}

	}


	sql = fmt.Sprintf("select sum(num) as a,sum(fee) as b ,sum(fee_cny) as c ,sum(total_cny) as d,token_admission_id  from trade where deal_time>=%d and deal_time<%d  and opt=2 group by token_admission_id", begin, end)
	r, err = DB.GetMysqlConn().Query(sql)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	if len(r) > 0 {
		for _, v := range r {
			h:=&TokenDailySheet{}
			t, _ := v["token_admission_id"]
			a, _ := v["a"]
			b, _ := v["b"]
			c, _ := v["c"]
			d, _ := v["d"]

			t_ := convert.BytesToIntAscii(t)
			m, ok := l[t_]
			if !ok {
				h.TokenId  = convert.BytesToIntAscii(t)
				h.SellTotal= convert.BytesToInt64Ascii(a)
				h.FeeSellTotal = convert.BytesToInt64Ascii(b)
				h.FeeSellCny = convert.BytesToInt64Ascii(c)
				h.SellTotalCny = convert.BytesToInt64Ascii(d)
				l[h.TokenId] = h
			} else {
				m.SellTotal  = convert.BytesToInt64Ascii(a)
				m.FeeSellTotal = convert.BytesToInt64Ascii(b)
				m.FeeSellCny = convert.BytesToInt64Ascii(c)
				m.SellTotalCny = convert.BytesToInt64Ascii(d)
			}
		}
	}

	for _,v:=range l  {
		p:=time.Unix(begin, 0).Format("2006-01-02 ")
		log.Infof("insert into token_id %d,time %s",v.TokenId,p)
		v.Date=begin
		_,err = DB.GetMysqlConn().Cols("token_id","fee_buy_cny","fee_buy_total","fee_sell_cny","fee_sell_total","buy_total","sell_total_cny","sell_total","date").InsertOne(v)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
	}
	//sql := fmt.Sprintf("insert into TokenDailySheet (`token_id`,`FeeBuyCny`,`FeeBuyTotal`,`FeeSellCny`,`FeeSellTotal`,`BuyTotal`,`BuyTotalCny`,`SellTotalCny`,`SellTotal`)  values(20001,0,1) on  DUPLICATE key update num=num+values(num)")
/*
	_,err = DB.GetMysqlConn().Insert(l)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
*/

	be:=begin+86400
	if be>time.Now().Unix() {
		return
	}
	Test3(begin+86400,end+86400)
}
