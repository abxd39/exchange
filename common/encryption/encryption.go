package encryption

import (
	"crypto/sha256"
	"fmt"
)

func Gensha256(phone string,nowtime int64,salt string)  []byte {
	s := fmt.Sprintf("%s%d%s",phone,nowtime,salt)

	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return bs
}
