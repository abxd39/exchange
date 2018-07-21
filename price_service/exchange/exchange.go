package exchange

/*
import (
	"digicon/price_service/model"
	"digicon/price_service/rpc/client"
	"time"
)

func LoadCacheQuene() error {
	var flag = true
	for flag {
		rsp, err := client.InnerService.TokenSevice.CallGetConfigQuene()
		if err != nil {
			time.Sleep(1*time.Second)
			continue
		}

		for _, v := range rsp.Data {
			model.ConfigQuenes[v.Name] = v
		}

		for _, v := range rsp.CnyData {
			model.ConfigCny[v.TokenId] = v
		}
		flag = false
	}
	return nil
}
*/
