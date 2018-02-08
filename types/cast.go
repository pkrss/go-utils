package types

import (
	"fmt"
	"reflect"
	"strconv"
)

func CastToInt64(o interface{}) (ret int64, err error) {
	if o == nil {
		return
	}
	t := reflect.ValueOf(o)
	if t.Kind() == reflect.Ptr {
		return CastToInt64(t.Elem().Interface())
	}
	switch v := o.(type) {
	case int64:
		ret = v
	case int:
		ret = int64(v)
	case float32:
		ret = int64(v)
	case float64:
		ret = int64(v)
	case string:
		ret, err = strconv.ParseInt(v, 10, 64)
		return
	default:
		err = fmt.Errorf("Unknown type to int64 error: %v %T", o, o)
	}
	return
}
