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
	default:
		err = fmt.Errorf("Unknown type to int64 error: %v %T", o, o)
	}
	return
}

func CastToInt(o interface{}) (ret int, err error) {
	if o == nil {
		return
	}
	t := reflect.ValueOf(o)
	if t.Kind() == reflect.Ptr {
		return CastToInt(t.Elem().Interface())
	}
	var v2 int64
	switch v := o.(type) {
	case int64:
		ret = int(v)
	case int:
		ret = v
	case float32:
		ret = int(v)
	case float64:
		ret = int(v)
	case string:
		v2, err = strconv.ParseInt(v, 10, 32)
		if err == nil {
			ret = int(v2)
		}
	default:
		err = fmt.Errorf("Unknown type to int error: %v %T", o, o)
	}
	return
}
