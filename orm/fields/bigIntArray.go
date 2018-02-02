package fields

import (
	"fmt"

	"github.com/pkrss/go-utils/types"
)

// because beego orm not support array field

type BigIntArray []int64

var nilBigIntArray = []int64{}

func MakeBigIntArray(d []int64) BigIntArray {
	return BigIntArray(d)
}

func BigIntArray2IntArray(jt *BigIntArray) []int64 {
	return []int64(*jt)
}

func (this *BigIntArray) UnmarshalJSON(b []byte) error {
	d, e := types.StringToInt64Array(string(b))
	if e != nil {
		return e
	}
	this.Set(d)
	return nil
}

func (this *BigIntArray) MarshalJSON() ([]byte, error) {
	t := BigIntArray2IntArray(this)

	if len(t) == 0 {
		return []byte("null"), nil
	}

	a := BigIntArray2IntArray(this)
	c := types.Int64ArrayToString(a)

	return []byte(c), nil
}

func (this BigIntArray) Value() []int64 {
	return BigIntArray2IntArray(&this)
}

func (this *BigIntArray) Set(d []int64) {
	*this = MakeBigIntArray(d)
}

func (this *BigIntArray) String() string {
	return types.Int64ArrayToString(BigIntArray2IntArray(this))
}

func (this *BigIntArray) Scan(src interface{}) error {
	return this.SetRaw(src)
}

func (this *BigIntArray) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case string:
		if len(d) > 0 {
			d2, e := types.StringToInt64Array(d)
			if e != nil {
				return fmt.Errorf("<BigIntArray.SetRaw> error string `%v`", value)
			}
			this.Set(d2)
		}
	case []int:
		d2 := types.IntArrayToInt64Array(d)
		this.Set(d2)
	case []int64:
		this.Set(d)
	case nil:
		return nil
	default:
		return fmt.Errorf("<BigIntArray.SetRaw> unknown value `%v`", value)
	}
	return nil
}

func (this *BigIntArray) RawValue() interface{} {
	if len(this.Value()) == 0 {
		return nil
	}
	return "array" + this.String()
}

func (this *BigIntArray) IsNil() bool {
	t := this.Value()

	if len(t) == 0 {
		return true
	}
	return false
}
