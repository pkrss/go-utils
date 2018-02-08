package fields

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// beacuse json unseriable time with string, i want to response with int64 timestamp format.

type JsonTime time.Time

var nilTime = time.Time{} // ).UnixNano()

func MakeJsonTimeNow() JsonTime {
	return MakeJsonTime(time.Now())
}

func MakeJsonTime(d time.Time) JsonTime {
	return JsonTime(d)
}

func JsonTime2Date(jt *JsonTime) time.Time {
	return time.Time(*jt)
}

func (ct *JsonTime) UnmarshalJSON(b []byte) (err error) {
	// s := strings.Trim(string(b), "\"")
	// if s == "null" {
	//     ct.Set(nilTime)
	//     return
	// }
	// v,err := strconv.ParseInt(s, 10, 64)
	// if err == nil {
	//     if v > 100000000000 {
	//         v = v / 1000
	//     }
	//     ct.Set(time.Unix(v, 0))
	// }else{
	//     ct.Set(nilTime)
	// }
	// return err
	d := time.Time{}
	e := d.UnmarshalJSON(b)
	if e != nil {
		return e
	}
	ct.Set(d)
	return e
}

func (ct *JsonTime) MarshalJSON() ([]byte, error) {
	t := JsonTime2Date(ct)

	if t == nilTime {
		return []byte("null"), nil
	}

	return t.MarshalJSON()
	// return []byte(fmt.Sprintf("%d", t.Unix()*1000)), nil
}

func (e JsonTime) Value() time.Time {
	return JsonTime2Date(&e)
}

func (e *JsonTime) Set(d time.Time) {
	*e = MakeJsonTime(d)
}

func (e *JsonTime) String() string {
	t := JsonTime2Date(e)
	return t.Format(time.RFC3339)
}

// for github.com/go-pg
func (this *JsonTime) Scan(value interface{}) error {
	switch d := value.(type) {
	case []byte:
		value = string(d)
	}
	switch d := value.(type) {
	case string:
		if len(d) > 0 {
			// t, _ := time.ParseInLocation("2006-01-02 15:04:05", d, time.Local)
			fmt := "2006-01-02 15:04:05"
			if len(d) > 12 && strings.ContainsAny(d[10:], "+-") {
				fmt = "2006-01-02 15:04:05Z07"
			}
			t, e := time.Parse(fmt, d)
			if e != nil {
				log.Printf("JSONTime parse [%s] error:%s\n", d, e.Error())
			}
			this.Set(t)
		}
		break
	case time.Time:
		*this = MakeJsonTime(d)
	case int64:
		*this = MakeJsonTime(time.Unix(d, 0))
	case nil:
		*this = JsonTime(nilTime)
	default:
		return fmt.Errorf("<JsonTime.Scan> unknown value `%v`", value)
	}
	return nil
}

// for github.com/go-pg
func (this JsonTime) AppendValue(b []byte, quote int) ([]byte, error) {
	tm := time.Time(this)

	if quote == 1 {
		b = append(b, '\'')
	}
	b = tm.UTC().AppendFormat(b, "2006-01-02 15:04:05")
	if quote == 1 {
		b = append(b, '\'')
	}
	return b, nil
}

func (e *JsonTime) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case string:
		if len(d) > 0 {
			t, _ := time.Parse(time.RFC3339, d) // time.RFC3339, 2017-10-13 16:11:11.283338 +0800 CST
			e.Set(t)
		}
		break
	case time.Time:
		*e = MakeJsonTime(d)
	case int64:
		*e = MakeJsonTime(time.Unix(d, 0))
	case nil:
		return nil
	default:
		return fmt.Errorf("<JsonTime.SetRaw> unknown value `%v`", value)
	}
	return nil
}

func (e *JsonTime) RawValue() interface{} {
	if e.IsNil() {
		return nil
	}
	return e.String()
}

func (this *JsonTime) IsNil() bool {
	t := this.Value().Unix()

	if t == 0 || t == nilTime.Unix() {
		return true
	}
	return false
}
