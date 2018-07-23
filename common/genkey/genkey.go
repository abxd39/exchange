package genkey

import (
	"fmt"
	"time"
)

func GetUnionKey(id, subid int) string {
	return fmt.Sprintf("%s_%s", id, subid)
}

func GetTimeUnionKey(id int64) string {
	return fmt.Sprintf("%d_%d", time.Now().Unix(), id)
}

func GetSymbol(a, b string) string {
	return fmt.Sprintf("%s/%s", a, b)
}

func GetPulishKey(symbol string) string {
	return fmt.Sprintf("%s:channel", symbol)
}
