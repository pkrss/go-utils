package fields

import (
	"fmt"
	"strings"
)

// beacuse beego orm read server null value will throw error.

// JsonbField postgres json field.
type JsonbField string

// Value return JsonbField value
func (j JsonbField) Value() string {
	return string(j)
}

// Set the JsonbField value
func (j *JsonbField) Set(d string) {
	*j = JsonbField(d)
}

// String convert JsonbField to string
func (j *JsonbField) String() string {
	return j.Value()
}

func (this *JsonbField) Scan(src interface{}) error {
	switch d := src.(type) {
	case []byte:
		src = string(d)
	}
	return this.SetRaw(src)
}

// SetRaw convert interface string to string
func (j *JsonbField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case string:
		j.Set(d)
	case nil: // only add this Row
		j.Set("")
	default:
		return fmt.Errorf("<JsonbField.SetRaw> unknown value `%s`", value)
	}
	return nil
}

// RawValue return JsonbField value
func (j *JsonbField) RawValue() interface{} {
	return j.Value()
}

func (this *JsonbField) IsNil() bool {
	t := this.Value()

	if t == "" {
		return true
	}
	return false
}

func (this *JsonbField) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		this.Set("")
		return nil
	}

	this.Set(s)
	return nil
}

func (this *JsonbField) MarshalJSON() ([]byte, error) {
	s := this.Value()

	if s == "" {
		return []byte("null"), nil
	}

	return []byte(s), nil
}
