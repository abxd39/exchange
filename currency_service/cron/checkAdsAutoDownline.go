package cron

import (
	log "github.com/sirupsen/logrus"
	"digicon/currency_service/model"
	"time"
)

func CheckAdsAutoDownline () {
	log.Println("timeNow:", time.Now())
	displayTokenIdList := new(model.CommonTokens).DisplayCurrencyTokens()

	pageNum := 20
	page := 1
	for _, token := range displayTokenIdList{
		typeId := 1
		adsList, total := new(model.Ads).AdsList(uint32(typeId), uint32(token.Id), uint32(page), uint32(pageNum))
		for _,ads := range adsList {
			model.AdsAutoDownline(ads.Id)
		}
		for int64(page * pageNum) < total {
			adsList, total = new(model.Ads).AdsList( uint32(typeId), uint32(token.Id), uint32(page), uint32(pageNum))
			for _,ads := range adsList {
				model.AdsAutoDownline(ads.Id)
			}
			if int64(page * pageNum) >= total {
				break
			}else{
				page += 1
			}
			//fmt.Println("typeId:", typeId, " page:", page)
			log.Println("typeId:", typeId, " page:", page)
			time.Sleep(time.Second * 3)
		}

	}
	pageNum = 20
	page = 1
	for _, token := range displayTokenIdList{
		typeId := 2
		adsList, total := new(model.Ads).AdsList(uint32(typeId), uint32(token.Id), uint32(page), uint32(pageNum))
		for _,ads := range adsList {
			model.AdsAutoDownline(ads.Id)
		}
		for int64(page * pageNum) < total {
			adsList, total = new(model.Ads).AdsList( uint32(typeId), uint32(token.Id), uint32(page), uint32(pageNum))
			for _,ads := range adsList {
				model.AdsAutoDownline(ads.Id)
			}
			if int64(page * pageNum) >= total {
				break
			}else{
				page += 1
			}
			//fmt.Println("typeId:", typeId, " page:", page)
			log.Println("typeId:", typeId, " page:", page)
			time.Sleep(time.Second * 3)
		}

	}

	log.Println("check end:", time.Now())
	//fmt.Println("check end:", time.Now())

}