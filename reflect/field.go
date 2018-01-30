package reflect

import (
	"reflect"
	"strings"

	"github.com/pkrss/go-utils/types"
)

func GetStructFieldToString(ob interface{}, field string, caseSensitive ...bool) string {
	v := GetStructField(ob, field, caseSensitive...)
	return types.GetValueString(v)
}

func GetStructField(v interface{}, field string, caseSensitive ...bool) interface{} {

	if field == "" {
		return nil
	}

	ks := strings.Split(field, ".")

	ret := v

	caseS := true
	if len(caseSensitive) > 0 {
		caseS = caseSensitive[0]
	}

	var i int
	for c := len(ks); i < c; i++ {
		k := ks[i]

		ret = GetStructFieldSimple(ret, k, caseS)
		if ret == nil {
			break
		}
	}

	return ret
}

func GetStructFieldSimple(v interface{}, field string, caseSensitive bool) interface{} {

	if v == nil || field == "" {
		return nil
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}

	if caseSensitive {
		v2 := val.FieldByName(field)
		if v2.IsValid() {
			return v2.Interface()
		}
	} else {
		field = strings.ToLower(field)
		c := val.NumField()
		for i := 0; i < c; i++ {
			valueField := val.Field(i)
			typeField := val.Type().Field(i)

			if strings.ToLower(typeField.Name) == field {
				return valueField.Interface()
			}

			// tag := typeField.Tag
			// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
		}
	}
	return nil
}

func GetStructFieldNames(v interface{}, structColsParams ...[]string) (ret []string) {
	n2v := GetStructFieldName2ValueMap(v, structColsParams...)
	ret = make([]string, 0)
	if n2v != nil {
		for k, _ := range n2v {
			ret = append(ret, k)
		}
	}
	return
}

func GetStructFieldName2ValueMap(v interface{}, structColsParams ...[]string) (ret map[string]interface{}) {
	m := GetStructFieldMap(v, structColsParams...)
	if m == nil {
		return nil
	}
	ret = make(map[string]interface{}, 0)
	for k, v := range m {
		ret[k] = v.Interface()
	}
	return ret
}

func GetStructFieldMap(v interface{}, structColsParams ...[]string) (ret map[string]reflect.Value) {
	return GetStructFieldMap2(v, nil, structColsParams...)
}

func GetStructFieldMap2(v interface{}, ret map[string]reflect.Value, structColsParams ...[]string) map[string]reflect.Value {

	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}

	ret = make(map[string]reflect.Value)

	var cols []string
	var excludeCols []string

	if len(structColsParams) > 1 {
		excludeCols = structColsParams[1]
	}

	if len(structColsParams) > 0 {
		cols = structColsParams[0]
	}

	c := val.NumField()
	colsCount := len(cols)
	excludeColsCount := len(excludeCols)
	var f bool
	for i := 0; i < c; i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		s := typeField.Name

		if colsCount > 0 {
			f = false
			for j := 0; j < colsCount; j++ {
				if cols[j] == s {
					f = true
					break
				}
			}
			if !f {
				continue
			}
		}

		if excludeColsCount > 0 {
			f = false
			for j := 0; j < excludeColsCount; j++ {
				if cols[j] == s {
					f = true
					break
				}
			}

			if f {
				continue
			}

		}

		if typeField.Anonymous {
			GetStructFieldMap2(valueField.Interface(), ret, structColsParams...)
			continue
		}

		ret[s] = valueField // .Interface()

		// tag := typeField.Tag
		// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}
	return ret
}

type FieldInfo struct {
	Tag reflect.StructTag
	Val reflect.Value
}

func GetStructFieldInfoMap(v interface{}) (ret map[string]FieldInfo) {
	return GetStructFieldInfoMap2(v, nil)
}

func GetStructFieldInfoMap2(v interface{}, ret map[string]FieldInfo) map[string]FieldInfo {

	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}

	ret = make(map[string]FieldInfo)

	c := val.NumField()
	for i := 0; i < c; i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		s := typeField.Name

		if typeField.Anonymous {
			GetStructFieldInfoMap2(valueField.Interface(), ret)
			continue
		}

		f := FieldInfo{}
		f.Tag = typeField.Tag
		f.Val = valueField
		ret[s] = f
	}
	return ret
}
