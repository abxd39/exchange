package handler

import (
	"digicon/currency_service/model"
	log "github.com/sirupsen/logrus"
	"fmt"
)


func InitCheckOrderStatus(){
	log.Println("check orders status ....")
	ods, err := new(model.Order).GetOrdersByStatus()
	if err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
	}else{
		log.Errorln(len(ods))
		for _,od := range ods{
			log.Println(od.Id, od.ExpiryTime)
			model.CheckOrderExiryTime(od.Id, od.ExpiryTime)
		}
	}
}