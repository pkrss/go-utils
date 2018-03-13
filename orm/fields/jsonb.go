package fields

import (
	"bytes"
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

// for github.com/go-pg
type Parser struct {
	b []byte
}

func NewParse(b []byte) *Parser {
	return &Parser{
		b: b,
	}
}

func (p *Parser) Valid() bool {
	return len(p.b) > 0
}
func (p *Parser) Advance() {
	p.b = p.b[1:]
}

func (p *Parser) Peek() byte {
	if p.Valid() {
		return p.b[0]
	}
	return 0
}
func (p *Parser) Skip(c byte) bool {
	if p.Peek() == c {
		p.Advance()
		return true
	}
	return false
}
func (p *Parser) Read() byte {
	if p.Valid() {
		c := p.b[0]
		p.Skip(c)
		return c
	}
	return 0
}
func (p *Parser) SkipBytes(b []byte) bool {
	if len(b) > len(p.b) {
		return false
	}
	if !bytes.Equal(p.b[:len(b)], b) {
		return false
	}
	p.b = p.b[len(b):]
	return true
}

// for github.com/go-pg
func (this JsonbField) AppendValue(b []byte, quote int) []byte {

	jsonb := []byte(string(this))
	if quote == 1 {
		b = append(b, '\'')
	}

	p := NewParse(jsonb)
	for p.Valid() {
		c := p.Read()
		switch c {
		case '\'':
			if quote == 1 {
				b = append(b, '\'', '\'')
			} else {
				b = append(b, '\'')
			}
		case '\000':
			continue
		case '\\':
			if p.SkipBytes([]byte("u0000")) {
				b = append(b, "\\\\u0000"...)
			} else {
				b = append(b, '\\')
				if p.Valid() {
					b = append(b, p.Read())
				}
			}
		default:
			b = append(b, c)
		}
	}

	if quote == 1 {
		b = append(b, '\'')
	}

	return b
}

// for github.com/go-pg
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
