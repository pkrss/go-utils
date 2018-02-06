package inner

import (
	"reflect"
	"strings"
	"time"

	pkReflect "github.com/pkrss/go-utils/reflect"
	pkStrings "github.com/pkrss/go-utils/strings"
)

type IsNullInterface interface {
	IsNil() bool
}

func GetStructDbFieldsAndValues(ob interface{}, writeMode bool, structColsParams ...*pkReflect.StructSelCols) map[string]interface{} {
	selCols := pkReflect.GetStructFieldInfoMap(ob, structColsParams...)
	ret := make(map[string]interface{}, 0)
	for k, v := range selCols {
		dbKey := pkStrings.StringToCamelCase(k)

		val := v.Val

		ormStr := v.Tag.Get("orm")
		if ormStr != "" {
			ormStr = strings.ToLower(ormStr)
			ss := strings.Split(ormStr, ";")

			isNil := false
			valTmp := val
			if valTmp.CanAddr() {
				valTmp = valTmp.Addr()
			}
			switch v := valTmp.Interface().(type) {
			case IsNullInterface:
				isNil = v.IsNil()
			}
			noInsertCol := false
			for _, s := range ss {
				if s == "ro" {
					if writeMode {
						noInsertCol = true
					}
				} else if s == "null" {
					if isNil {
						ret[dbKey] = nil
						noInsertCol = true
					}
				} else if s == "auto_now_add" {
					if isNil {
						val = reflect.ValueOf(time.Now())
					}
				} else if s == "auto_now" {
					val = reflect.ValueOf(time.Now())
				}
			}
			if noInsertCol {
				continue
			}
		}

		switch val.Kind() {
		case reflect.Ptr:
			val = val.Elem()
		}

		ret[dbKey] = val.Interface()
	}
	return ret
}
