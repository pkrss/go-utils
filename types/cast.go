package types

import (
	"encoding/json"
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
	var i int64
	i, err = CastToInt64(o)
	if err != nil {
		return
	}
	ret = int(i)
	return
}

func CastToFloat64(o interface{}) (ret float64, err error) {
	if o == nil {
		return
	}
	t := reflect.ValueOf(o)
	if t.Kind() == reflect.Ptr {
		return CastToFloat64(t.Elem().Interface())
	}
	switch v := o.(type) {
	case int64:
		ret = float64(v)
	case int:
		ret = float64(v)
	case float32:
		ret = float64(v)
	case float64:
		ret = v
	case string:
		ret, err = strconv.ParseFloat(v, 64)
	default:
		err = fmt.Errorf("Unknown type to float64 error: %v %T", o, o)
	}
	return
}

func CastTofloat32(o interface{}) (ret float32, err error) {
	var i float64
	i, err = CastToFloat64(o)
	if err != nil {
		return
	}
	ret = float32(i)
	return
}

func CastToString(o interface{}) (ret string, err error) {
	if o == nil {
		return
	}
	t := reflect.ValueOf(o)
	if t.Kind() == reflect.Ptr {
		return CastToString(t.Elem().Interface())
	}
	switch v := o.(type) {
	case int64:
		ret = strconv.FormatInt(v, 10)
	case int:
		ret = strconv.FormatInt(int64(v), 10)
	case float32:
		ret = strconv.FormatFloat(float64(v), 'f', 8, 32)
	case float64:
		ret = strconv.FormatFloat(v, 'f', 8, 64)
	case string:
		ret = v
	default:
		var by []byte
		by, err = json.Marshal(&v)
		if err != nil {
			return
		}
		ret = string(by)
	}
	return
}
