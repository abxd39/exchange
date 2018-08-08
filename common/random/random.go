package random

import (
	"fmt"
	"math/rand"
	"time"
)

func Random6dec() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
}

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}



var headList = []string{"https://sdun.oss-cn-shenzhen.aliyuncs.com/aa6b04d79c699fe229464dd3cd86ce88.png"}


func GetRandHead() string {
	headListLen := len(headList)
	n := rand.Intn(headListLen)
	if n >= headListLen {
		return "https://sdun.oss-cn-shenzhen.aliyuncs.com/aa6b04d79c699fe229464dd3cd86ce88.png"
	}else{
		return headList[n]
	}
}


var randHeadList = []string{"https://sdun.oss-cn-shenzhen.aliyuncs.com/aa6b04d79c699fe229464dd3cd86ce88.png",
	"https://sdun.oss-cn-shenzhen.aliyuncs.com/dbed8a73d3912c8fae53df635f98706c.png",
	"https://sdun.oss-cn-shenzhen.aliyuncs.com/6ab58203a1dc916432de00af83c1daca.png",
	"https://sdun.oss-cn-shenzhen.aliyuncs.com/2490df8d46315a2aaa3e6ef37a60e166.png",
}
func SetRegisterRandHeader() string {
	headListLen := len(randHeadList)
	n := rand.Intn(headListLen)
	if n >= headListLen {
		return "https://sdun.oss-cn-shenzhen.aliyuncs.com/aa6b04d79c699fe229464dd3cd86ce88.png"
	}else{
		return randHeadList[n]
	}
}