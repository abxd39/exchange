package encryption

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

func GenMd5AndReverse(pwd string) string  {
	h := sha256.New()
	h.Write([]byte(pwd))
	bs := h.Sum(nil)
	s:= fmt.Sprintf("%x", bs)
	return ReverseString(s)
}

func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}


func Gensha256(phone string, nowtime int64, salt string) string {
	s := fmt.Sprintf("%s%d%s", phone, nowtime, salt)
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func GenBase64(input string) []byte {
	s := base64.StdEncoding.EncodeToString([]byte(input))
	return []byte(s)
}

// 产生订单 ID
//  uid, 时间秒,
func CreateOrderId(userId uint64, tokenId int32) (orderId string) {
	tn := time.Now()
	tnn := tn.UnixNano()
	tnnStr := strconv.FormatInt(tnn, 10) // 获取微秒时间
	tnnStrLen := len(tnnStr)
	tStr := tn.Format("2006-01-02 15:04:05.34234")

	var buffer bytes.Buffer
	buffer.WriteString(strconv.FormatInt(int64(userId), 10))
	buffer.WriteString(tStr[2:4])   // 年
	buffer.WriteString(tStr[5:7])   // 月
	buffer.WriteString(tStr[8:10])  // 日
	buffer.WriteString(tStr[17:19]) // 秒
	//buffer.WriteString(tStr[22:24])
	buffer.WriteString(tnnStr[tnnStrLen-5 : tnnStrLen-2]) // 微秒
	orderId = buffer.String()
	return
}
