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

var dataMapInvalidLocker *sync.RWMutex = &sync.RWMutex{}
var dataMapInvalid map[string]int64 = make(map[string]int64)

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

func DataDelete(key string) {
	dataMapLocker.Lock()

	delete(dataMap, "tmpTime-"+key)
	delete(dataMap, "tmp-"+key)
	delete(dataMapInvalid, key)

	dataMapLocker.Unlock()

	return
}

func TmpDataSet(key string, obj interface{}) {
	DataSet("tmpTime-"+key, strconv.FormatInt(time.Now().Unix(), 10))
	DataSet("tmp-"+key, obj)
}

func TmpPeriodDataSet(key string, obj interface{}, invalidSeconds int64) {
	DataSet("tmp-"+key, obj)

	dataMapInvalidLocker.Lock()
	dataMapInvalid[key] = time.Now().Unix() + invalidSeconds
	dataMapInvalidLocker.Unlock()
}

func TmpRunOnce(key string, obj interface{}, invalidSeconds int64) {

	dataMapInvalidLocker.Lock()
	d := make([]string, 0)
	now := time.Now().Unix()
	for k, v := range dataMapInvalid {
		if now > v {
			d = append(d, k)
		}
	}

	for _, v := range d {
		delete(dataMapInvalid, v)
		delete(dataMap, "tmp-"+v)
	}

	dataMapInvalidLocker.Unlock()
}

func TmpDataGet(key string, invalidSeconds int64) (obj interface{}, ok bool) {
	i, _ := DataGet("tmpTime-" + key)
	if i == nil {
		return nil, false
	}

	n, _ := strconv.ParseInt(i.(string), 10, 64)
	now := time.Now().Unix()
	if now-n > invalidSeconds {
		DataDelete(key)
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
