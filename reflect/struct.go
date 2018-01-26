package reflect

import (
	"reflect"
)

func CopyStruct(target interface{}, src interface{}) {
	if target == nil || src == nil {
		return
	}

	targetV := reflect.ValueOf(target)
	srcV := reflect.ValueOf(src)

	switch targetV.Kind() {
	case reflect.Ptr:
		targetV = targetV.Elem()
	}
	switch srcV.Kind() {
	case reflect.Ptr:
		srcV = srcV.Elem()
	}

	c := srcV.NumField()
	for i := 0; i < c; i++ {
		srcValueField := srcV.Field(i)
		srcTypeField := srcV.Type().Field(i)

		targetValueField := targetV.FieldByName(srcTypeField.Name)
		if !targetValueField.IsValid() {
			continue
		}
		targetValueField.Set(srcValueField)
	}
}
