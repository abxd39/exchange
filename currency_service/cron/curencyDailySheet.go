package cron

import (
	"digicon/common/convert"
	"digicon/currency_service/model"
	"fmt"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"time"
)

type DailyCountSheet struct {
}

func (this DailyCountSheet) Run() {
	fmt.Println("DailyCountSheet...")

	tokens := new(model.CommonTokens).List()
	mod := new(model.Order)

	nTime := time.Now()
	yesTime := nTime.AddDate(0, 0, -1)     // 统计昨天的
	yesDate := yesTime.Format("2006-01-02")
	startTime:= fmt.Sprintf("%s 00:00:00", yesDate)
	endTime := fmt.Sprintf("%s 23:59:59", yesDate)

	//endTime := nTime.Format("2006-01-02 15:04:05")
	//yesday := yesTime.Format("2006-01-02")


	log.Println("startTime:", startTime)
	log.Println("endTime:", endTime)

	for _, tk := range tokens {
		tkId := tk.Id
		fmt.Println("count token id: ", tkId)

		/// check exists ....
		mcuds := new(model.CurrencyDailySheet)
		checkResult, err := mcuds.GetOneDay(tkId, yesDate)
		if err != nil {
			log.Errorln(err)
		}
		fmt.Println("checkResult.Id:", checkResult.Id)
		if checkResult.Id != 0 {
			log.Printf("tokenid: %v 今天已经统计了", tkId)
			continue
		}
		//////

		orders, err := mod.GetOrderByTokenIdByTime(tkId, startTime, endTime)
		if err != nil {
			log.Errorln(err)
			continue
		}
		var sellTotal int64
		var buyTotal int64
		var sellTotalCny int64
		var buyTotaoCny int64
		var feeSellTotal int64
		var feeSellCny int64
		var feeBuyTotal int64
		var feeBuyCny int64

		for _, od := range orders {
			if od.AdType == model.BuyType {
				buyTotal += od.Num
				buyTotaoCny += convert.Int64MulInt64By8Bit(od.Num, od.Price)
				feeBuyTotal += od.Fee
				feeBuyCny += convert.Int64MulInt64By8Bit(od.Fee, od.Price)
			}
			if od.AdType == model.SellType {
				sellTotal += od.Num
				sellTotalCny += convert.Int64MulInt64By8Bit(od.Num, od.Price)
				feeSellTotal += od.Fee
				feeSellCny += convert.Int64MulInt64By8Bit(od.Fee, od.Price)
			}
		}

		sumBuy, sumSell, err := mod.GetTotalSum(tkId, endTime)


		if err != nil {
			log.Errorln(err)
		}
		mcds := model.CurrencyDailySheet{
			TokenId:      int32(tkId),
			SellTotal:    sellTotal,
			SellCny:      sellTotalCny,
			BuyTotal:     buyTotal,
			BuyCny:       buyTotaoCny,
			FeeSellTotal: feeSellTotal,
			FeeSellCny:   feeSellCny,
			FeeBuyTotal:  feeBuyTotal,
			FeeBuyCny:    feeBuyCny,

			BuyTotalAll:    sumBuy.BuyTotalAll,
			BuyTotalAllCny: sumBuy.BuyTotalAllCny,

			SellTotalAll:    sumSell.SellTotalAll,
			SellTotalAllCny: sumSell.SellTotalAllCny,

			Total:    sumBuy.BuyTotalAll + sumSell.SellTotalAll,
			TotalCny: sumBuy.BuyTotalAllCny + sumSell.SellTotalAllCny,
			Date:     yesDate,
		}

		err = mcds.InsertOneDay()
		if err != nil {
			log.Errorln(err)
		}else{
			log.Println("统计成功!", tkId)
		}
	}

}

//启动多个任务
func DailyStart() {
	fmt.Println("daily count start ....")
	log.Println("daily count start ....")

	i := 0
	c := cron.New()

	//AddFunc
	spec := "0 0 1 * *" // every day ...
	specTwo := "0 0 4 * *" // every day ...

	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})
	//AddJob方法
	c.AddJob(spec, DailyCountSheet{})
	c.AddJob(specTwo, DailyCountSheet{})


	//启动计划任务
	c.Start()
	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}
}
