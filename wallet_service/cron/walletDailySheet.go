package cron

import (
	"fmt"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"digicon/wallet_service/model"
	"time"
	"admin/utils/convert"
)

type WalletDailyCountSheet struct {
}

func (this WalletDailyCountSheet) Run(){
	fmt.Println("Daily CountSheet ...")

	tokens, err := new(models.Tokens).ListTokens()
	if err != nil {
		log.Errorln(err)
		fmt.Println(err)
		return
	}

	nTime := time.Now()

	yesDayTime := nTime.AddDate(0, 0, -1)   // 统计昨天的

	yesDate := yesDayTime.Format("2006-01-02")

	startTime:= fmt.Sprintf("%s 00:00:00", yesDate)
	endTime := fmt.Sprintf("%s 23:59:59", yesDate)


	fmt.Println("startTime:", startTime, " endTime:",endTime, "nowtime:", nTime.Unix())

	for _, tk := range tokens{
		fmt.Println(tk.Id)
		tkId := tk.Id
		fmt.Println("count token id:",tkId)
		fmt.Println("startTime:", startTime, " endTime:", endTime)

		/// check exists ....
		mcuds := new(models.TokenInoutDailySheet)
		has, err := mcuds.CheckOneDay(uint32(tkId), startTime)
		if err != nil {
			log.Errorln(err)
			continue
		}
		fmt.Println("has:", has)

		if has.Id != 0 {
			log.Errorf("tokenid: %v 今天已经统计了", tkId)
			continue
		}else{
			fmt.Println("checkResult.Id:")
		}
		//////



		////
		tkinoutModel := new(models.TokenInout)
		tkinouts, err :=  tkinoutModel.GetInOutByTokenIdByTime(uint32(tkId), startTime, endTime)
		if err != nil {
			log.Errorln(err)
			fmt.Println(err)
			continue
		}

		var total_day_num int64
		var total_day_cny int64

		var total   int64
		var total_day_num_fee int64
		var total_fee int64
		var total_day_fee_cny   int64
		var total_put  int64
		var total_day_put int64
		var total_day_put_cny int64

		var tokenName string

		tokenName = tk.Mark

		for _, tkinout := range tkinouts{
			fmt.Println(tkinout)
			fmt.Println(tkinout.Amount)
			fmt.Println(tkinout.Fee)
			//token_name := tkinout.TokenName

			if tkinout.Opt == 1 {
				// 充币 (充币没有手续费)
				total_day_put += tkinout.Amount
				total_day_put_cny += tkinout.AmountCny

			}else if tkinout.Opt == 2{
				//  提币
				total_day_num += tkinout.Amount
				total_day_cny += tkinout.AmountCny

				total_day_num_fee += tkinout.Fee
				total_day_fee_cny += tkinout.FeeCny
			}
		}

		outtotal, err := tkinoutModel.GetOutSumByTokenId(uint32(tkId), endTime)
		if err != nil {
			log.Errorln(err)
			fmt.Println(err)
		}
		total = outtotal.Total
		total_fee = outtotal.TotalFee

		intotal, err := tkinoutModel.GetInSumByTokenId(uint32(tkId), endTime)
		if err != nil {
			log.Errorln(err)
			fmt.Println(err)
		}
		total_put = intotal.TotalPut


		onedayInOutModel := &models.TokenInoutDailySheet{
			TokenId:     tkId,
			TokenName:   tokenName,
			TotalDayNum: total_day_num,
			TotalDayCny: total_day_cny,
			TotalDayNumFee:  total_day_num_fee,
			TotalDayFeeCny:  total_day_fee_cny,
			TotalDayPut:     total_day_put,
			TotalDayPutCny:  total_day_put_cny,
			Total:           total,
			TotalFee:        total_fee,
			TotalPut:        total_put,
			Date:            yesDate,
		}
		//if has.Id  != 0  {
		//	err = onedayInOutModel.UpdateOneDayTotal(int64(has.Id))
		//	fmt.Println(onedayInOutModel.Id)
		//}else{
			err = onedayInOutModel.InsertOneDayTotal()
			if err!=nil{
				log.Infof("插入数据失败!!",err.Error())
				fmt.Println("插入数据失败!!",err.Error())
				continue
			}
			fmt.Println(onedayInOutModel.Id)
		//}


		if err != nil {
			log.Errorln("wallet统计失败", err)
			fmt.Println(err)
		}else{
			log.Println("统计成功!", tkId)
			log.Println("统计成功!", tkId)
		}

		fmt.Println(" insert ......")
	}

}

func (this WalletDailyCountSheet) CountData(begin int64){
	fmt.Println("工具走起")
	now :=time.Now().Unix()
	for begin< 1535385600{
		beginStr:=time.Unix(begin,0).Format("2006-01-02 15:04:05")
		this.Tool(beginStr)
		fmt.Println(beginStr)
		begin+=86400
	}
}

func (this WalletDailyCountSheet) Tool(beginStr string){
	fmt.Println("Daily CountSheet ...")

	tokens, err := new(models.Tokens).ListTokens()
	if err != nil {
		log.Errorln(err)
		fmt.Println(err)
	}

	nTime := time.Now()

	//yesDayTime := nTime.AddDate(0, 0, -1)   // 统计昨天的

	//yesDate := yesDayTime.Format("2006-01-02")
	//
	startTime:= fmt.Sprintf("%s 00:00:00", beginStr[:11])
	endTime := fmt.Sprintf("%s 23:59:59", beginStr[:11])


	fmt.Println("startTime:", startTime, " endTime:",endTime, "nowtime:", nTime.Unix())

	for _, tk := range tokens{
		fmt.Println(tk.Id)
		tkId := tk.Id
		fmt.Println("count token id:",tkId)
		fmt.Println("startTime:", startTime, " endTime:", endTime)

		/// check exists ....
		mcuds := new(models.TokenInoutDailySheet)
		has, err := mcuds.CheckOneDay(uint32(tkId), startTime)
		if err != nil {
			log.Errorln(err)
		}
		fmt.Println("has:", has)

		if has.Id != 0 {
			log.Errorf("tokenid: %v 今天已经统计了", tkId)
			continue
		}else{
			fmt.Println("checkResult.Id:")
		}
		//////



		////
		tkinoutModel := new(models.TokenInout)
		tkinouts, err :=  tkinoutModel.GetInOutByTokenIdByTime(uint32(tkId), startTime, endTime)
		if err != nil {
			log.Errorln(err)
			fmt.Println(err)
			continue
		}

		var total_day_num int64
		var total_day_cny int64

		var total   int64
		var total_day_num_fee int64
		var total_fee int64
		var total_day_fee_cny   int64
		var total_put  int64
		var total_day_put int64
		var total_day_put_cny int64

		var tokenName string

		tokenName = tk.Mark

		for _, tkinout := range tkinouts{
			fmt.Println("tkinout=",tkinout)
			fmt.Println("tkinout.Amount=",tkinout.Amount)
			fmt.Println("tkinout.Fee=",tkinout.Fee)
			//token_name := tkinout.TokenName

			if tkinout.Opt == 1 {
				// 充币 (充币没有手续费)
				//total_day_put += tkinout.Amount
				total_day_put,err=convert.Int64AddInt64(total_day_put,tkinout.Amount)
				fmt.Println("total_day_put",total_day_put)
				if err!=nil{
					fmt.Println("error=",err.Error())
				}

				total_day_put_cny ,err= convert.Int64AddInt64(total_day_put_cny,tkinout.AmountCny)
				if err!=nil{
					fmt.Println(err.Error())
				}

			}else if tkinout.Opt == 2{
				//  提币
				total_day_num ,_= convert.Int64AddInt64(total_day_num, tkinout.Amount)
				total_day_cny ,_ = convert.Int64AddInt64 (total_day_cny,tkinout.AmountCny)

				total_day_num_fee ,_= convert.Int64AddInt64(total_day_num_fee,tkinout.Fee)
				total_day_fee_cny ,_= convert.Int64AddInt64(total_day_fee_cny,tkinout.FeeCny)
			}
		}

		outtotal, err := tkinoutModel.GetOutSumByTokenId(uint32(tkId), endTime)
		if err != nil {
			log.Errorln(err)
			fmt.Println(err)
		}
		total = outtotal.Total
		total_fee = outtotal.TotalFee

		intotal, err := tkinoutModel.GetInSumByTokenId(uint32(tkId), endTime)
		if err != nil {
			log.Errorln(err)
			fmt.Println(err)
		}
		total_put = intotal.TotalPut


		onedayInOutModel := &models.TokenInoutDailySheet{
			TokenId:     tkId,
			TokenName:   tokenName,
			TotalDayNum: total_day_num,
			TotalDayCny: total_day_cny,
			TotalDayNumFee:  total_day_num_fee,
			TotalDayFeeCny:  total_day_fee_cny,
			TotalDayPut:     total_day_put,
			TotalDayPutCny:  total_day_put_cny,
			Total:           total,
			TotalFee:        total_fee,
			TotalPut:        total_put,
			Date:            startTime,
		}
		//if has.Id  != 0  {
		//	err = onedayInOutModel.UpdateOneDayTotal(int64(has.Id))
		//	fmt.Println(onedayInOutModel.Id)
		//}else{
		//	err = onedayInOutModel.InsertOneDayTotal()
		//	fmt.Println(onedayInOutModel.Id)
		//}
		fmt.Println("赋值后的结果为==",onedayInOutModel)
		err = onedayInOutModel.InsertOneDayTotal()
		if err!=nil{
			log.Infof("插入数据失败!!",err.Error())
			fmt.Println("插入数据失败!!",err.Error())
			continue
		}

		if err != nil {
			log.Errorln("wallet统计失败", err)
			fmt.Println("wallet统计失败",err)
		}else{
			log.Println("统计成功!", tkId)
			fmt.Println("统计成功!", tkId)
		}

		fmt.Println(" insert ......")
	}

}


//启动
func DailyStart() {
	fmt.Println("daily count start ....")
	log.Println("daily count start ....")

	i := 0
	c := cron.New()

	//AddFunc
	spec := "0 0 1 * * *" // every day ...
	specTwo := "0 0 4 * * *" // every day ...

	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})
	//AddJob方法
	c.AddJob(spec, WalletDailyCountSheet{})
	c.AddJob(specTwo, WalletDailyCountSheet{})
	//启动计划任务
	c.Start()
	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}
}
