package reflect

import (
	"reflect"
	"strings"

	"github.com/pkrss/go-utils/types"
)

type StructSelCols struct {
	IncludeCols []string
	ExcludeCols []string
}

func GetStructFieldToString(ob interface{}, field string, caseSensitive ...bool) string {
	v := GetStructFieldValue(ob, field, caseSensitive...)
	return types.GetValueString(v)
}

func GetStructFieldValue(v interface{}, field string, caseSensitive ...bool) interface{} {

	val := GetStructField(v, field, caseSensitive...)
	if !val.IsValid() {
		return nil
	}
	return val.Interface()
}

func GetStructFieldValueSimple(v interface{}, field string, caseSensitive bool) interface{} {

	val := GetStructFieldSimple(reflect.ValueOf(v), field, caseSensitive)
	if !val.IsValid() {
		return nil
	}
	return val.Interface()
}

func GetStructField(v interface{}, field string, caseSensitive ...bool) reflect.Value {

	if field == "" {
		return reflect.Value{}
	}

	ks := strings.Split(field, ".")

	ret := reflect.ValueOf(v)

	caseS := true
	if len(caseSensitive) > 0 {
		caseS = caseSensitive[0]
	}

	var i int
	for c := len(ks); i < c; i++ {
		k := ks[i]

		ret = GetStructFieldSimple(ret, k, caseS)
		if !ret.IsValid() {
			break
		}

		if i == c-1 {
			return ret
		}
	}

	return reflect.Value{}
}

func GetStructFieldSimple(val reflect.Value, field string, caseSensitive bool) reflect.Value {

	if !val.IsValid() || field == "" {
		return reflect.Value{}
	}

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}

	if caseSensitive {
		v2 := val.FieldByName(field)
		if v2.IsValid() {
			return v2
		}
	} else {
		field = strings.ToLower(field)
		c := val.NumField()
		for i := 0; i < c; i++ {
			valueField := val.Field(i)
			typeField := val.Type().Field(i)

			if strings.ToLower(typeField.Name) == field {
				return valueField
			}

			// tag := typeField.Tag
			// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
		}
	}
	return reflect.Value{}
}

func GetStructFieldNames(v interface{}, structColsParams ...*StructSelCols) (ret []string) {
	n2v := GetStructFieldName2ValueMap(v, structColsParams...)
	ret = make([]string, 0)
	if n2v != nil {
		for k, _ := range n2v {
			ret = append(ret, k)
		}
	}
	return
}

func GetStructFieldName2ValueMap(v interface{}, structColsParams ...*StructSelCols) (ret map[string]interface{}) {
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

func GetStructFieldMap(v interface{}, structColsParams ...*StructSelCols) (ret map[string]reflect.Value) {
	m := GetStructFieldInfoMap(v, structColsParams...)
	if m == nil {
		return nil
	}
	ret = make(map[string]reflect.Value, 0)
	for k, v := range m {
		ret[k] = v.Val
	}
	return ret
}

type FieldInfo struct {
	Tag reflect.StructTag
	Val reflect.Value
}

func GetStructFieldInfoMap(v interface{}, structColsParams ...*StructSelCols) (ret map[string]FieldInfo) {
	return GetStructFieldInfoMap2(v, nil, structColsParams...)
}

func GetStructFieldInfoMap2(v interface{}, ret map[string]FieldInfo, structColsParams ...*StructSelCols) map[string]FieldInfo {

	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}

	ret = make(map[string]FieldInfo)

	var cols []string
	var excludeCols []string

	if len(structColsParams) > 0 && structColsParams[0] != nil {
		cols = structColsParams[0].IncludeCols
		excludeCols = structColsParams[0].ExcludeCols
	}

	colsCount := len(cols)
	excludeColsCount := len(excludeCols)
	var f bool

	c := val.NumField()
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
				if excludeCols[j] == s {
					f = true
					break
				}
			}

			if f {
				continue
			}

		}

		if typeField.Anonymous {
			GetStructFieldInfoMap2(valueField.Interface(), ret, structColsParams...)
			continue
		}

		f := FieldInfo{}
		f.Tag = typeField.Tag
		f.Val = valueField
		ret[s] = f
	}
	return ret
}

func SetStructFieldValue(obj interface{}, field string, v interface{}, caseSensitive ...bool) bool {
	val := GetStructField(obj, field, caseSensitive...)
	if !val.IsValid() || !val.CanSet() {
		return false
	}

	val.Set(reflect.ValueOf(v))
	return true
}
