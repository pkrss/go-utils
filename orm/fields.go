package orm

import (
	"reflect"

	pkReflect "github.com/pkrss/go-utils/reflect"
	pkStrings "github.com/pkrss/go-utils/strings"
)

type StructColsStruct struct {
	SelCols   []string
	UnSelCols []string
}

func getObDbFieldsAndValues(ob interface{}, structColsParams ...*pkReflect.StructSelCols) map[string]interface{} {
	selCols := pkReflect.GetStructFieldMap(ob, structColsParams...)
	ret := make(map[string]interface{}, 0)
	for k, v := range selCols {
		dbKey := pkStrings.StringToCamelCase(k)
		switch v.Kind() {
		case reflect.Ptr:
			v = v.Elem()
		}
		ret[dbKey] = v.Interface()
	}
	return ret
}
