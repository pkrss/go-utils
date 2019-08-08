package container

import "sync"

// OrderedMap ...
type OrderedMap struct {
	IDList     []string               `json:"idList"`
	IDMap      map[string]interface{} `json:"idMap"`
	LimitCount int                    `json:"LimitCount"`
	locker     *sync.RWMutex
}

// Init ...
func (c *OrderedMap) Init(limitCount int) {

	if limitCount < 20 {
		limitCount = 20
	}
	c.LimitCount = limitCount

	if c.IDList == nil {
		c.IDList = make([]string, 0, limitCount+1)
	}
	if c.IDMap == nil {
		c.IDMap = make(map[string]interface{})
	}
	for _, v := range c.IDList {
		c.IDMap[v] = 1
	}
	c.locker = &sync.RWMutex{}
}

// Exist ...
func (c *OrderedMap) Exist(k string) bool {
	c.locker.RLock()
	defer c.locker.RUnlock()
	_, ok := c.IDMap[k]
	return ok
}

// Get ...
func (c *OrderedMap) Get(k string) (v interface{}, ok bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	v, ok = c.IDMap[k]
	return v, ok
}

// Put true: new push, false: update old value
func (c *OrderedMap) Put(k string, v interface{}) bool {

	c.locker.Lock()
	defer c.locker.Unlock()

	_, ok := c.IDMap[k]
	c.IDMap[k] = v

	if ok {
		return false
	}

	c.IDList = append(c.IDList, k)

	if cnt := len(c.IDList); cnt > c.LimitCount {
		for i, l := 0, cnt-c.LimitCount; i < l; i++ {
			delete(c.IDMap, c.IDList[i])
		}
		c.IDList = c.IDList[cnt-c.LimitCount:]
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
	for k, v := range c.IDMap {
		idMap[k] = v
	}
	ret["idList"] = idList
	ret["idMap"] = idMap

	return ret
}
