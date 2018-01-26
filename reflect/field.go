package reflect

import (
	"reflect"
	"strconv"
	"strings"
)

func GetStructFieldToString(ob interface{}, field string) string {
	v := GetStructField(ob, field)
	return GetValueString(v)
}

func GetValueString(v interface{}) string {
	ret := ""

	if v == nil {
		return ret
	}

	switch v2 := v.(type) {
	case string:
		ret = v2
	case bool:
		ret = strconv.FormatBool(v2)
	case int:
		ret = strconv.Itoa(v2)
	case int64:
		ret = strconv.FormatInt(v2, 10)
	case float32:
		ret = strconv.FormatFloat(float64(v2), 'g', 30, 32)
	case float64:
		ret = strconv.FormatFloat(v2, 'g', 30, 64)
	default:
		ret = v.(string)
	}

	return ret
}

func GetStructField(v interface{}, field string) interface{} {

	if field == "" {
		return nil
	}

	ks := strings.Split(field, ".")

	ret := v

	var i int
	for c := len(ks); i < c; i++ {
		k := ks[i]

		ret = GetStructFieldSimple(ret, k)
		if ret == nil {
			break
		}
	}

	return ret
}

func GetStructFieldSimple(v interface{}, field string) interface{} {

	if v == nil || field == "" {
		return nil
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}

	v2 := val.FieldByName(field)
	if v2.IsValid() {
		return v2.Interface()
	}

	// c := val.NumField()
	// for i := 0; i < c; i++ {
	// 	valueField := val.Field(i)
	// 	typeField := val.Type().Field(i)

	// 	if typeField.Name == field {
	// 		return valueField.Interface()
	// 	}

	// 	// tag := typeField.Tag
	// 	// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	// }
	return nil
}

func GetStructFieldName2ValueMap(v interface{}, cols ...string) (ret map[string]interface{}) {

	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}

	ret = make(map[string]interface{})

	c := val.NumField()
	c2 := len(cols)
	for i := 0; i < c; i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		s := typeField.Name
		f := false

		if c2 != 0 {
			for j := 0; j < c2; j++ {
				if cols[j] == s {
					f = true
					break
				}
			}
		}

		if f {
			continue
		}

		ret[s] = valueField.Interface()

		// tag := typeField.Tag
		// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}
	return
}
