package orm

import (
	"reflect"
	"strings"
)

func getObFields(ob interface{}) []interface{} {
	return getObFields2(ob, nil)

}

func getObFields2(ob interface{}, fields []interface{}) []interface{} {
	if ob == nil {
		return nil
	}

	val := reflect.ValueOf(ob)

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}

	if fields == nil {
		fields = make([]interface{}, 0)
	}

	c := val.NumField()
	for i := 0; i < c; i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		if typeField.Anonymous {
			getObFields2(valueField.Interface(), fields)
			continue
		}

		tag := typeField.Tag
		tagVal := tag.Get("orm")
		tagVal = strings.Trim(tagVal, " \t;")
		if tagVal == "-" {
			continue
		}

		var field interface{}

		switch valueField.Kind() {
		case reflect.Ptr:
			field = valueField.Elem().Interface()
		// case reflect.Struct, reflect.Slice, reflect.Map:
		// 	field = valueField.Interface()
		default:
			field = valueField.Addr().Interface()
		}

		if field == nil {
			continue
		}

		fields = append(fields, field)

		// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}

	return fields
}

func setObFields(ob interface{}, fields []interface{}) error {
	return nil
}
