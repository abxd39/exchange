package convert

import (
	"bytes"
	"encoding/binary"
	"github.com/shopspring/decimal"
)

func ByteToInt32(b []byte) (x uint32) {
	b_buf := bytes.NewBuffer(b)
	b_buf = bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, x)
	return
}


func Int64ToFloat64By8Bit(b int64) (x float64) {
	a:=decimal.New(b,-8)
	x,_ =a.Float64()
	return
}