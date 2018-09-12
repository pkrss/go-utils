package tmp

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/pkrss/go-utils/types"
)

var dataMapLocker = &sync.RWMutex{}
var dataMap = make(map[string]interface{})

var dataMapInvalidLocker = &sync.RWMutex{}
var dataMapInvalid = make(map[string]int64)

type invalidCbFun func(interface{})

var dataMapInvalidCb = make(map[string]invalidCbFun)

// DataSet ...
func DataSet(key string, obj interface{}) {
	dataMapLocker.Lock()

	dataMap[key] = obj

	dataMapLocker.Unlock()
}

// DataGet ...
func DataGet(key string) (obj interface{}, ok bool) {
	dataMapLocker.Lock()

	obj, ok = dataMap[key]

	dataMapLocker.Unlock()

	return
}

// DataDelete ...
func DataDelete(key string) {
	dataMapLocker.Lock()

	delete(dataMap, "tmpTime-"+key)
	delete(dataMap, "tmp-"+key)
	delete(dataMapInvalid, key)

	dataMapLocker.Unlock()

	return
}

func TmpDataSet(key string, obj interface{}) {
	TmpDataSet2(key, obj, strconv.FormatInt(time.Now().Unix(), 10))
}

func TmpDataSet2(key string, obj interface{}, t string) {
	DataSet("tmpTime-"+key, t)
	DataSet("tmp-"+key, obj)
}

func TmpPeriodDataSet(key string, obj interface{}, invalidSeconds int64) {
	DataSet("tmp-"+key, obj)

	dataMapInvalidLocker.Lock()
	dataMapInvalid[key] = time.Now().Unix() + invalidSeconds
	dataMapInvalidLocker.Unlock()
}

func TmpPeriodDataSetWithInvalidCb(key string, obj interface{}, invalidSeconds int64, invalidCb invalidCbFun) {
	DataSet("tmp-"+key, obj)

	dataMapInvalidLocker.Lock()
	dataMapInvalid[key] = time.Now().Unix() + invalidSeconds
	dataMapInvalidCb[key] = invalidCb
	dataMapInvalidLocker.Unlock()
}

func TmpRunOnce(invalidSeconds int64) {

	invalidCbMap := make(map[interface{}]interface{})

	dataMapInvalidLocker.Lock()
	d := make([]string, 0)
	now := time.Now().Unix()
	for k, v := range dataMapInvalid {
		if now > v {
			d = append(d, k)
		}
	}

	for _, v := range d {
		cb, ok := dataMapInvalidCb[v]
		if ok {
			v2, ok2 := dataMapInvalid[v]
			if ok2 {
				invalidCbMap[cb] = v2
			}
		}

		delete(dataMapInvalidCb, v)
		delete(dataMapInvalid, v)
		delete(dataMap, "tmp-"+v)
	}

	dataMapInvalidLocker.Unlock()

	for cb, v := range invalidCbMap {
		cb.(invalidCbFun)(v)
	}
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

func FileDataGetObject(key string, ret interface{}) (e error) {

	s, e2 := FileDataGetString(key)
	if e2 != nil {
		return
	}
	if s == "" {
		e = errors.New("FileDataGetObject() fail, not exist: " + key)
		return
	}

	return json.Unmarshal([]byte(s), ret)
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

func FileTmpDataGetObject(key string, invalidSeconds int64, ret interface{}) (obj interface{}, ok bool) {
	obj, ok = TmpDataGet(key, invalidSeconds)
	if ok {
		return
	}

	s, e := FileDataGetString("tmpTime-" + key)
	if e == nil {
		i, e2 := strconv.ParseInt(s, 10, 64)
		if e2 == nil {
			now := time.Now().Unix()
			if now-i < invalidSeconds {
				e = FileDataGetObject(key, ret)
				if e == nil {
					TmpDataSet2(key, ret, strconv.FormatInt(i, 10))
					obj = ret
					ok = true
					return
				} else {
					log.Println(e.Error())
				}
			}
		}
	}

	return
}

func FileTmpDataSetObject(key string, ret interface{}) {
	TmpDataSet(key, ret)
	e := FileDataSet(key, ret)
	if e != nil {
		log.Println(e.Error())
	}
	e = FileDataSet("tmpTime-"+key, strconv.FormatInt(time.Now().Unix(), 10))
	if e != nil {
		log.Println(e.Error())
	}
}
