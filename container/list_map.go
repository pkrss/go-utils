package container

// ListMap ...
type ListMap struct {
	IDList []string `json:"idList"`
	idSet  map[string]int
}

// Init ...
func (c ListMap) Init(cnt int) {

	if c.IDList == nil {
		c.IDList = make([]string, 0, cnt)
	}
	if c.idSet == nil {
		c.idSet = make(map[string]int)
	}
	for _, v := range c.IDList {
		c.idSet[v] = 1
	}
}
