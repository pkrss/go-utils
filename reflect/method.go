package reflect

import (
	"reflect"
	"strings"
)

func GetStructMethod(v interface{}, field string) interface{} {

	if field == "" {
		return nil
	}

	ks := strings.Split(field, ".")

	ret := v

	var i int
	for c := len(ks); i < c; i++ {
		k := ks[i]

		ret = GetStructMethodSimple(ret, k)
		if ret == nil {
			break
		}
	}

	return ret
}

func GetStructMethodSimple(v interface{}, method string) interface{} {

	if v == nil || method == "" {
		return nil
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}

	v2 := val.MethodByName(method)
	if v2.IsValid() {
		return v2.Interface()
	}

	// c := val.NumMethod()
	// for i := 0; i < c; i++ {
	// 	valueMethod := val.Method(i)
	// 	typeMethod := val.Type().Method(i)

	// 	if typeMethod.Name == method {
	// 		return valueMethod.Interface()
	// 	}
	// }
	return nil
}
