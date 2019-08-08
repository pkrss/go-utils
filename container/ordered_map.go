package container

import "sync"

// OrderedMap ...
type OrderedMap struct {
	IDList     []string `json:"idList"`
	idMap      map[string]interface{}
	limitCount int
	locker     *sync.RWMutex
}

// Init ...
func (c *OrderedMap) Init(limitCount int) {

	if limitCount < 20 {
		limitCount = 20
	}

	if c.IDList == nil {
		c.IDList = make([]string, 0, limitCount+1)
	}
	if c.idMap == nil {
		c.idMap = make(map[string]interface{})
	}
	for _, v := range c.IDList {
		c.idMap[v] = 1
	}
	c.locker = &sync.RWMutex{}
}

// Exist ...
func (c *OrderedMap) Exist(k string) bool {
	c.locker.RLock()
	defer c.locker.RUnlock()
	_, ok := c.idMap[k]
	return ok
}

// Get ...
func (c *OrderedMap) Get(k string) (v interface{}, ok bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	v, ok = c.idMap[k]
	return v, ok
}

// Put true: new push, false: update old value
func (c *OrderedMap) Put(k string, v interface{}) bool {

	c.locker.Lock()
	defer c.locker.Unlock()

	_, ok := c.idMap[k]
	c.idMap[k] = v

	if ok {
		return false
	}

	c.IDList = append(c.IDList, k)

	if cnt := len(c.IDList); cnt > c.limitCount {
		for i, l := 0, cnt-c.limitCount; i < l; i++ {
			delete(c.idMap, c.IDList[i])
		}
		c.IDList = c.IDList[cnt-c.limitCount:]
	}

	return true
}

// Info ...
func (c *OrderedMap) Info() map[string]interface{} {

	c.locker.RLock()
	defer c.locker.RUnlock()

	ret := make(map[string]interface{})
	idList := make([]string, len(c.IDList))
	for k, v := range c.IDList {
		idList[k] = v
	}
	idMap := make(map[string]interface{})
	for k, v := range c.idMap {
		idMap[k] = v
	}
	ret["idList"] = idList
	ret["idMap"] = idMap

	return ret
}
