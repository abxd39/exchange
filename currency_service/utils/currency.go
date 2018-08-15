package utils

import (
	"fmt"
	"strconv"
	"github.com/shopspring/decimal"
	"strings"
	log "github.com/sirupsen/logrus"
)

func Round2(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}



/*
保留两位小数
 */
func Int64ToStringBy8Bit(b int64) string {
	a := decimal.New(b, 0)
	r := a.Div(decimal.New(100000000, 0))
	s  := r.String()
	splitResult := strings.Split(s, ".")
	log.Println("splitResult s:", splitResult)
	var result string
	if len(splitResult) >= 2 {
		//fmt.Println("splitResult:", splitResult)
		if len(splitResult[1]) >= 3{
			result = splitResult[0] + "." + splitResult[1][:3]
		}else{
			result = splitResult[0] + "." + splitResult[1][:]
		}
	}else{
		result = splitResult[0]
	}
	return result
}

//t := dd.Div(dp).Mul(d)
//k, _ := t.Float64()
//s := fmt.Sprintf("%.2f", k)