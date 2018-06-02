package tmp

import (
	"sync"
)

var dataMapLocker *sync.RWMutex = &sync.RWMutex{}
var dataMap map[string]interface{} = make(map[string]interface{})

func DataSet(key string, obj interface{}) {
	dataMapLocker.Lock()

	dataMap[key] = obj

	dataMapLocker.Unlock()
}

func DataGet(key string) (obj interface{}, ok bool) {
	dataMapLocker.Lock()

	obj, ok = dataMap[key]

	dataMapLocker.Unlock()

	return
}
