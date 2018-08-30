package model

import (
	"digicon/common/convert"
	. "digicon/token_service/dao"
	"fmt"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"time"
)

type TokenDailySheet struct {
	Id           int64 `xorm:"not null pk autoincr comment('自增id')BIGINT(20)"`
	TokenId      int   `xorm:"not null comment('货币id') INT(11)"`
	FeeBuyCny    int64 `xorm:"not null comment('买手续费折合cny') BIGINT(20)"`
	FeeBuyTotal  int64 `xorm:"not null comment('买手续费总额') BIGINT(20)"`
	FeeSellCny   int64 `xorm:"not null comment('卖手续费折合cny') BIGINT(20)"`
	FeeSellTotal int64 `xorm:"not null comment('卖手续费总额') BIGINT(20)"`
	BuyTotal     int64 `xorm:"not null comment('买总额') BIGINT(20)"`
	BuyTotalCny  int64 `xorm:"not null comment('买总额折合') BIGINT(20)"`
	SellTotalCny int64 `xorm:"not null comment('卖总额折合') BIGINT(20)"`
	SellTotal    int64 `xorm:"not null comment('卖总额') BIGINT(20)"`
	Date         int64 `xorm:"not null comment('时间戳，精确到天') BIGINT(20)"`
	//Day          time.Time `xorm:"comment('那天') DATETIME"`
}

func (t *TokenDailySheet) TimingFunc(begin, end int64) {
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
			h := &TokenDailySheet{}
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

			h.TokenId = convert.BytesToIntAscii(t)
			h.BuyTotal = convert.BytesToInt64Ascii(a)
			h.FeeBuyTotal = convert.BytesToInt64Ascii(b)
			h.FeeBuyCny = convert.BytesToInt64Ascii(c)
			h.BuyTotalCny = convert.BytesToInt64Ascii(d)

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
			h := &TokenDailySheet{}
			t, _ := v["token_admission_id"]
			a, _ := v["a"]
			b, _ := v["b"]
			c, _ := v["c"]
			d, _ := v["d"]

			t_ := convert.BytesToIntAscii(t)
			m, ok := l[t_]
			if !ok {
				h.TokenId = convert.BytesToIntAscii(t)
				h.SellTotal = convert.BytesToInt64Ascii(a)
				h.FeeSellTotal = convert.BytesToInt64Ascii(b)
				h.FeeSellCny = convert.BytesToInt64Ascii(c)
				h.SellTotalCny = convert.BytesToInt64Ascii(d)
				l[h.TokenId] = h
			} else {
				m.SellTotal = convert.BytesToInt64Ascii(a)
				m.FeeSellTotal = convert.BytesToInt64Ascii(b)
				m.FeeSellCny = convert.BytesToInt64Ascii(c)
				m.SellTotalCny = convert.BytesToInt64Ascii(d)
			}
		}
	}

	for _, v := range l {
		p := time.Unix(begin, 0).Format("2006-01-02 ")
		log.Infof("insert into token_id %d,time %s", v.TokenId, p)
		v.Date = begin
		_, err = DB.GetMysqlConn().Cols("token_id", "fee_buy_cny", "fee_buy_total", "fee_sell_cny", "fee_sell_total", "buy_total", "sell_total_cny", "sell_total", "date").InsertOne(v)
		if err != nil {
			log.Errorln(err.Error())
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

	be := begin + 86400
	if be > time.Now().Unix() {
		return
	}
	t.TimingFunc(begin+86400, end+86400)
}

func (t *TokenDailySheet) Run() {
	toBeCharge := time.Now().Format("2006-01-02 ") + "00:00:00"
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
	unix := theTime.Unix()

	t.TimingFunc(unix-86400, unix)
}

//启动
func DailyStart() {
	fmt.Println("daily count start ....")
	log.Println("daily count start ....")

	i := 0
	c := cron.New()

	//AddFunc
	spec := "0 0 1 * *" // every day ...
	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})
	//AddJob方法
	c.AddJob(spec, &TokenDailySheet{})
	//启动计划任务
	c.Start()
	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}
}
