package pqsql

import (
	"reflect"
	"sync"

	"github.com/pkrss/go-utils/reflect"
)

var typeName2TableName = make(map[string]string)
var typeName2TableNameMutex sync.Mutex

func getTableName(ob interface{}) string {
	v := reflect.ValueOf(ob)
	typeName := v.Type().Name

	typeName2TableNameMutex.Lock()
	defer typeName2TableNameMutex.Unlock()

	tableName, ok := typeName2TableName[typeName]
	if ok {
		return tableName
	}

	m := reflect.GetStructMethod("TableName")

	if m != nil {
		tableName = m.Call()
	}

	typeName2TableName[typeName] = tableName
	return tableName
}
