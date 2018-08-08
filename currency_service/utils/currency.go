package utils

import (
	"fmt"
	"strconv"
	"math/rand"
)

func Round2(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}



var headList = []string{"https://sdun.oss-cn-shenzhen.aliyuncs.com/aa6b04d79c699fe229464dd3cd86ce88.png",
	"https://sdun.oss-cn-shenzhen.aliyuncs.com/fd4617859847ca447350cf82d403a943.png",
	"https://sdun.oss-cn-shenzhen.aliyuncs.com/f832383013c9c455b9304eaf36a87d26.png",
}


func GetRandHead() string {
	headListLen := len(headList)
	n := rand.Intn(headListLen)
	if n >= headListLen {
		return "https://sdun.oss-cn-shenzhen.aliyuncs.com/aa6b04d79c699fe229464dd3cd86ce88.png"
	}else{
		return headList[n]
	}
}

