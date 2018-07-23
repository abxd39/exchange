package model

import (
	"digicon/config_service/confcommon"
	"digicon/config_service/confjson"
	"digicon/config_service/consulclient"
	"digicon/config_service/dao"
	. "digicon/config_service/log"
	"encoding/json"
	"fmt"
)

type ConfigQuenesModel struct {
}

func (c *ConfigQuenesModel) GetAllQuenes() []confjson.ConfigQuenes {
	t := make([]confjson.ConfigQuenes, 0)
	err := dao.DB.GetMysqlTokenConn().Where("switch=1").Find(&t)
	if err != nil {
		Log.Errorln(err.Error())
		return nil
	}
	return t
}

func (c *ConfigQuenesModel) PutToConsul() (err error) {
	result := c.GetAllQuenes()
	quenesJson, err := json.Marshal(result)
	if err != nil {
		fmt.Println(" quenes marshal json error:", quenesJson)
	}
	client := consulclient.NewClient()
	err = client.Put(confcommon.ConfigQuenesKey, quenesJson)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (c *ConfigQuenesModel) GetQuenesFromConsul() (result []confjson.ConfigQuenes, err error) {

	client := consulclient.NewClient()
	entries, _, _ := client.List(confcommon.ConfigQuenesKey)

	var quenesList []confjson.ConfigQuenes
	for _, enpair := range entries {
		var cquenes []confjson.ConfigQuenes
		err = json.Unmarshal(enpair.Value, &cquenes)
		if err != nil {
			fmt.Println("json unmarshal error!", err.Error())
			continue
		} else {
			for _, quenes := range cquenes {
				quenesList = append(quenesList, quenes)
			}
		}
	}
	if err != nil {
		fmt.Println("json unmarshal error!", err.Error())
		return
	}
	result = quenesList
	return
}
