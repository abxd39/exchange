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
//0.00001001
func StringToInt64By8Bit(s string)  (int64,error){
	d,err:=decimal.NewFromString(s)
	l:=d.Round(8).Coefficient().Int64()
	//g,_:=decimal.NewFromString("100000000")
	//l:=d.Mul(g)
	return l,err
}

func Float64ToInt64By8Bit(s float64)  (int64){
	d:=decimal.NewFromFloat(s)
	l:=d.Round(8).Coefficient().Int64()
	return l
}