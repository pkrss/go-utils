package pqsql

import (
	"reflect"
	"sync"

	pkReflect "github.com/pkrss/go-utils/reflect"
)

var typeName2TableName = make(map[string]string)
var typeName2TableNameMutex sync.RWMutex

func getTableName(ob interface{}) string {
	v := reflect.ValueOf(ob)
	typeName := v.Type().Name()

	typeName2TableNameMutex.RLock()

	tableName, ok := typeName2TableName[typeName]
	if ok {
		typeName2TableNameMutex.RUnlock()
		return tableName
	}

	// outV := reflect.ValueOf(ob).MethodByName("TableName").Call([]reflect.Value{})
	// tableName = outV[0].String()

	m := pkReflect.GetStructMethod(ob, "TableName")
	if m != nil {
		tmpV := reflect.ValueOf(m)
		outV := tmpV.Call([]reflect.Value{})
		tableName = outV[0].String()
	}

	typeName2TableNameMutex.Lock()
	typeName2TableName[typeName] = tableName
	typeName2TableNameMutex.Unlock()
	return tableName
}
