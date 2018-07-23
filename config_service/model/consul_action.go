package model

import (
	"digicon/config_service/confcommon"
	"fmt"
)

func PutToConsul(PutType int32) (err error) {
	if PutType == confcommon.ConfigQueneType {
		mconfig := new(ConfigQuenesModel)
		err = mconfig.PutToConsul()
	}
	return
}

func GetFromConsul(GetType int32) (result interface{}, err error) {
	fmt.Println("gettype: ", GetType)
	if GetType == confcommon.ConfigQueneType {
		mconfig := new(ConfigQuenesModel)
		return mconfig.GetQuenesFromConsul()
	}
	return
}
