package convert

import (
	"github.com/shopspring/decimal"
	"bytes"
	"encoding/binary"
	"fmt"
	//"math/big"
)

func ByteToInt32(b []byte) (x uint32) {
	b_buf := bytes.NewBuffer(b)
	b_buf = bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, x)
	return
}

func ByteToInt64(b []byte) (x int64) {
	b_buf := bytes.NewBuffer(b)
	b_buf = bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, x)
	return
}

func Int64ToInt64By8Bit(b int64) int64 {

	a := decimal.New(b, 0)
	r := a.Mul(decimal.New(100000000, 0))
	return  r.IntPart()
}

func Int64ToFloat64By8Bit(b int64) (x float64) {
	a := decimal.New(b, -8)
	x, _ = a.Float64()
	return
}

func Int64ToStringBy8Bit(b int64) string {
	a := decimal.New(b, 0)
	r := a.Div(decimal.New(100000000, 0))
	return r.String()
}

//0.00001001
func StringToInt64By8Bit(s string) (int64, error) {
	d, err := decimal.NewFromString(s)
	l := d.Round(8).Coefficient().Int64()
	//g,_:=decimal.NewFromString("100000000")
	//l:=d.Mul(g)
	return l, err
}

func Float64ToInt64By8Bit(s float64) int64 {
	d := decimal.NewFromFloat(s)
	l := d.Round(8).Coefficient().Int64()
	return l
}

//请确保返回值不会越界否则调用下面的返回string类型
func Int64MulInt64By8Bit(a int64, b int64) int64 {
	dd := decimal.New(a, 0)
	dp := decimal.New(b, 0)
	m := dd.Mul(dp)
	d := decimal.New(100000000, 0)
	n := m.Div(d)
	return n.IntPart()
}

func Int64MulInt64MulInt64By16Bit(a int64, b, c int64) int64 {
	da := decimal.New(a, 0)
	db := decimal.New(b, 0)
	dc := decimal.New(c, 0)
	m := da.Mul(db).Mul(dc)
	d := decimal.New(100000000, 0)
	n := m.Div(d).Div(d)
	return n.IntPart()
}

func Int64MulInt64By8BitString(a int64, b int64) string {
	dd := decimal.New(a, 0)
	dp := decimal.New(b, 0)
	m := dd.Mul(dp)
	d := decimal.New(100000000, 0)
	n := m.Div(d)

	r := n.Div(decimal.New(100000000, 0))
	return r.String()
}

func Int64MulFloat64(a int64, b float64) int64 {
	dd := decimal.New(a, 0)
	dp := decimal.NewFromFloat(b)
	m := dd.Mul(dp)
	return m.IntPart()
}

//两数相除保持8位
func Int64DivInt64By8Bit(a int64, b int64) int64 {
	dd := decimal.New(a, 0)
	dp := decimal.New(b, 0)
	d := decimal.New(100000000, 0)

	//num := dd.Div(dp).Round(8).Coefficient().Int64()

	num := dd.Div(dp).Mul(d).IntPart()
	return num
}

//两数相除保持8位
func Int64DivInt64By8BitString(a int64, b int64) string {
	dd := decimal.New(a, 0)
	dp := decimal.New(b, 0)
	d := decimal.New(100000000, 0)
	num := dd.Div(dp).Mul(d).String()
	//num := dd.Div(dp).Round(8).Coefficient().String()
	return num
}

//两数相除保持2位
func Int64DivInt64StringPercent(a int64, b int64) string {
	dd := decimal.New(a, 0)
	dp := decimal.New(b, 0)
	d := decimal.New(100, 0)

	t := dd.Div(dp).Mul(d)
	k, _ := t.Float64()
	s := fmt.Sprintf("%.2f", k)

	return s
}


//两数相加保持3位
func Int64AddInt64Float64Percent(a int64, b int64) string {
	dd := decimal.New(a, 0)
	dp := decimal.New(b, 0)
	d := decimal.New(100000000, 0)


	t := dd.Add(dp).Div(d)
	k, _ := t.Float64()
	s := fmt.Sprintf("%.3f", k)

	return s
}

