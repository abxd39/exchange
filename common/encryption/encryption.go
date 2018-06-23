package encryption

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func Gensha256(phone string, nowtime int64, salt string) []byte {
	s := fmt.Sprintf("%s%d%s", phone, nowtime, salt)

	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return bs
}

func GenBase64(input string) []byte {
	return []byte(base64.StdEncoding.EncodeToString([]byte(input)))
}
