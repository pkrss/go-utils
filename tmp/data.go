package tmp

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/pkrss/go-utils/types"
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

func TmpDataSet(key string, obj interface{}) {
	DataSet("tmpTime-"+key, strconv.FormatInt(time.Now().Unix(), 10))
	DataSet("tmp-"+key, obj)
}

func TmpDataGet(key string, invalidSeconds int64) (obj interface{}, ok bool) {
	i, _ := DataGet("tmpTime-" + key)
	if i == nil {
		return nil, false
	}

	n, _ := strconv.ParseInt(i.(string), 10, 64)
	now := time.Now().Unix()
	if now-n > invalidSeconds {
		return nil, false
	}
	return DataGet("tmp-" + key)
}

type FuncFileWrite func(string, string) error
type FuncFileRead func(string) (string, error)

var gFuncFileWrite FuncFileWrite
var gFuncFileRead FuncFileRead

func SetFuncFile(funcFileRead FuncFileRead, funcFileWrite FuncFileWrite) {
	gFuncFileWrite = funcFileWrite
	gFuncFileRead = funcFileRead
}

func FileDataGetString(key string) (ret string, e error) {

	if gFuncFileRead == nil {
		e = errors.New("gFuncFileRead is nil")
		return
	}
	ret, e = gFuncFileRead(key)

	return
}

func FileDataSet(key string, obj interface{}) (e error) {
	if gFuncFileWrite == nil {
		e = errors.New("gFuncFileWrite is nil")
		return
	}

	var s string
	s, e = types.CastToString(obj)
	if e != nil {
		return
	}

	e = gFuncFileWrite(key, s)

	return
}
