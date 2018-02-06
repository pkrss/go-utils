package fields

import (
	"fmt"
	"strings"
)

// beacuse beego orm string empty write to postgresql will throw error : - pq: invalid input syntax for uuid: ""

type UUID string

func MakeUUID(d string) UUID {
	return UUID(strings.Replace(d, "-", "", -1))
}

func UUID2String(jt *UUID) string {
	return strings.Replace(string(*jt), "-", "", -1)
}

func (this *UUID) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		this.Set("")
		return nil
	}
	this.Set(s)
	return nil
}

func (this *UUID) MarshalJSON() ([]byte, error) {
	t := UUID2String(this)

	if len(t) == 0 {
		return []byte("null"), nil
	}

	return []byte("\"" + t + "\""), nil
}

func (this UUID) Value() string {
	return UUID2String(&this)
}

func (this *UUID) Set(d string) {
	*this = MakeUUID(d)
}

func (this *UUID) String() string {
	return this.Value()
}

// for github.com/go-pg
func (this *UUID) Scan(src interface{}) error {
	switch d := src.(type) {
	case []byte:
		src = string(d)
	}
	return this.SetRaw(src)
}

// for github.com/go-pg
func (this UUID) AppendValue(b []byte, quote int) ([]byte, error) {
	s := string(this)
	if quote == 2 {
		b = append(b, '"')
	} else if quote == 1 {
		b = append(b, '\'')
	}

	for i := 0; i < len(s); i++ {
		c := s[i]

		if c == '\000' {
			continue
		}

		if quote >= 1 {
			if c == '\'' {
				b = append(b, '\'', '\'')
				continue
			}
		}

		if quote == 2 {
			if c == '"' {
				b = append(b, '\\', '"')
				continue
			}
			if c == '\\' {
				b = append(b, '\\', '\\')
				continue
			}
		}

		b = append(b, c)
	}

	if quote >= 2 {
		b = append(b, '"')
	} else if quote == 1 {
		b = append(b, '\'')
	}
	return b, nil
}

func (this *UUID) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case string:
		this.Set(d)
	case nil:
		return nil
	default:
		return fmt.Errorf("<UUID.SetRaw> unknown value `%v`", value)
	}
	return nil
}

func (this *UUID) RawValue() interface{} {
	s := this.Value()
	if len(s) == 0 {
		return nil
	}
	return s
}

func (this *UUID) IsNil() bool {
	t := this.Value()

	if len(t) == 0 {
		return true
	}
	return false
}
