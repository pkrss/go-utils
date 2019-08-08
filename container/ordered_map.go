package container

// OrderedMap ...
type OrderedMap struct {
	IDList     []string `json:"idList"`
	idMap      map[string]interface{}
	limitCount int
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
}

// Exist ...
func (c *OrderedMap) Exist(k string) bool {
	_, ok := c.idMap[k]
	return ok
}

// Get ...
func (c *OrderedMap) Get(k string) (v interface{}, ok bool) {
	v, ok = c.idMap[k]
	return v, ok
}

// Push ...
func (c *OrderedMap) Push(k string, v interface{}) bool {
	// if _, ok := c.idMap[k]; ok {
	// 	return false
	// }
	c.idMap[k] = v
	c.IDList = append(c.IDList, k)

	if cnt := len(c.IDList); cnt > c.limitCount {
		for i, l := 0, cnt-c.limitCount; i < l; i++ {
			delete(c.idMap, c.IDList[i])
		}
		c.IDList = c.IDList[cnt-c.limitCount:]
	}

	return true
}
