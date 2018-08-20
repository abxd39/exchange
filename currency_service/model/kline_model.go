package model

import (
	"time"
	"fmt"
	"math/rand"
	"digicon/common/convert"
)



func GenerateKline() (uOrderHistoryList []Order){
	ctk := new(CommonTokens)
	fctk := ctk.Get(0, "BTC")
	tokenId := fctk.Id
	tctcy := new(TokenConfigTokenCNy)
	tctcy.GetPrice(uint32(tokenId))
	now := time.Now()
	for i:=0;  i< 20 ;  i++  {
		mm, _ := time.ParseDuration(fmt.Sprintf("-%dm", 5 * i)) // 过期时间15分钟
		createtime := now.Add(mm).Format("2006-01-02 15:04:05")
		rand.Seed(time.Now().UnixNano())
		n := rand.Float64()
		price := convert.Int64MulInt64By8Bit(tctcy.Price, convert.Float64ToInt64By8Bit(n))
		uOrderHistoryList = append(uOrderHistoryList, Order{Price: price, CreatedTime:createtime})
	}
	return
}


