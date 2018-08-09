package utils

import (
	"digicon/common/convert"
	"fmt"
	"strconv"
)

// 以下是试例 - 等待真正的接口

func PriceFiat(price int64, fiat string) float64 {

	pe := convert.Int64ToFloat64By8Bit(price)

	// 人民币汇率
	usd := 6.3

	if fiat == "cny" {
		return pe
	}
	if fiat == "usd" {
		return float64(int((pe/usd)*100)) / 100
	}

	return pe
}
func PriceFiatMaxLimit(max int64, balance int64, type_id int32, fiat string) int64 {
	if max == 0 {
		return max
	}

	// 虚拟货币当前价格和人民币汇率
	var btc float64 = 5120
	usd := 6.3

	price := convert.Int64ToFloat64By8Bit(balance) * btc

	if fiat == "cny" {

		if price > float64(max) {
			return max
		}
		return int64(price)
	}

	if fiat == "usd" {

		price = price * usd

		if price > float64(max) {
			return max
		}
		return int64(price)
	}

	return max
}

func NumFiat(num int64, balance int64) float64 {

	if num > balance {
		return convert.Int64ToFloat64By8Bit(balance)
	}

	return convert.Int64ToFloat64By8Bit(num)
}




func Round2(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}



