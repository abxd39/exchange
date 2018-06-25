package encryption

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
	"strconv"
	"bytes"
)

func Gensha256(phone string, nowtime int64, salt string) []byte {
	s := fmt.Sprintf("%s%d%s", phone, nowtime, salt)

	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return bs
}

func GenBase64(input string) []byte {
	s := base64.StdEncoding.EncodeToString([]byte(input))
	return []byte(s)
}


// 产生订单 ID
//  uid, 时间秒,
func CreateOrderId(userId int32, tokenId uint64) (orderId string) {
	tn := time.Now()
	tnn := tn.UnixNano()
	tnnStr := strconv.FormatInt(tnn, 10) // 获取微秒时间
	tnnStrLen := len(tnnStr)
	tStr := tn.Format("2006-01-02 15:04:05.34234")

	var buffer bytes.Buffer
	buffer.WriteString(strconv.FormatInt(int64(userId), 10))
	buffer.WriteString(tStr[2:4])    // 年
	buffer.WriteString(tStr[5:7])    // 月
	buffer.WriteString(tStr[8:10])   // 日
	buffer.WriteString(tStr[17:19])  // 秒
	//buffer.WriteString(tStr[22:24])
	buffer.WriteString(tnnStr[tnnStrLen-5 : tnnStrLen-2]) // 微秒
	orderId = buffer.String()
	return
}
