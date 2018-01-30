package orm

import (
	"time"
)

// beacuse json unseriable time with string, i want to response with int64 timestamp format.

type JsonTime time.Time

var nilTime = time.Time{} // ).UnixNano()

// // Value - Implementation of valuer for database/sql
// func (this *JsonTime) Value() (driver.Value, error) {
// 	return time.Time(*this), nil
// }

func (this *JsonTime) Scan(src interface{}) error {
	if src == nil {
		*this = JsonTime(nilTime)
		return nil
	}

	b := src.([]byte)
	t, e := time.ParseInLocation("2006-01-02 15:04:05", string(b), time.Local)
	if e != nil {
		return e
	}
	*this = JsonTime(t)
	return nil
}

func (this *JsonTime) UnmarshalJSON(b []byte) (err error) {
	d := time.Time{}
	e := d.UnmarshalJSON(b)
	if e != nil {
		return e
	}
	*this = JsonTime(d)
	return e
}

func (this *JsonTime) MarshalJSON() ([]byte, error) {
	t := time.Time(*this)

	if t == nilTime {
		return []byte("null"), nil
	}

	return t.MarshalJSON()
}
