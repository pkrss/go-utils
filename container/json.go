package container

import (
	"container/list"
	"encoding/json"
)

func ListMarshalJSON(this *list.List) ([]byte, error) {

	if this == nil {
		return []byte("null"), nil
	}

	if this.Len() == 0 {
		return []byte("[]"), nil
	}

	c := "["

	i := 0
	b := this.Front()
	for e := this.Back(); true; b = b.Next() {
		if i > 0 {
			c += ","
		}

		b2, e2 := json.Marshal(b.Value)
		if e2 != nil {
			return nil, e2
		}

		c += string(b2)

		if b == e {
			break
		}
		i++
	}
	c += "]"
	return []byte(c), nil
}
