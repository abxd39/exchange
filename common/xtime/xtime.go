package xtime

import (
	"database/sql/driver"
	"strconv"
	"time"
)

// Time be used to MySql timestamp converting.
type Time int64

func (jt *Time) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
		*jt = Time(sc.Unix())
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*jt = Time(i)
	}
	return
}

func (jt Time) Value() (driver.Value, error) {
	return time.Unix(int64(jt), 0), nil
}

func (jt Time) Time() time.Time {
	return time.Unix(int64(jt), 0)
}

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration time.Duration

func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

const (
	LAYOUT_DATE         = "2006-01-02"
	LAYOUT_DATE_TIME    = "2006-01-02 15:04:05"
	LAYOUT_DATE_TIME_12 = "2006-01-02 03:04:05"
)

// 日期转时间戳
func Date2Unix(date string, layout string) int64 {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(layout, date, loc)
	return theTime.Unix()
}

// 时间戳转日期
func Unix2Date(unix int64, layout string) string {
	return time.Unix(unix, 0).Format(layout)
}
