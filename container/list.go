package container

import (
	"container/list"
	// "github.com/pkrss/go-utils/beans"
)

func ListSub(data *list.List, offset int, limit int) *list.List {
	if data == nil || data.Len() == 0 {
		return data
	}

	l := data.Len()

	if limit == 0 {
		return list.New()
	}

	if limit < 0 {
		limit = l
	}

	if offset < 0 {
		offset = l + offset
	}

	if offset < 0 {
		offset = 0
	}

	if offset > l {
		return list.New()
	}

	var elem2 *list.Element
	elem := data.Front()
	end := data.Back()

	for offset > 0 {
		offset--

		elem2 = elem

		if elem != end {
			elem = elem.Next()
		}

		data.Remove(elem2)

		if elem2 == end {
			break
		}
	}

	for limit > 0 {
		limit--

		if elem != end {
			elem = elem.Next()
		}

		if elem == end {
			break
		}
	}

	for {
		elem2 = elem

		if elem != end {
			elem = elem.Next()
		}

		data.Remove(elem2)

		if elem2 == end {
			break
		}
	}

	return data
}

// func ListSubPage(data *list.List, pageable *beans.Pageable) *list.List {
// 	if data != nil && pageable != nil {
// 		limit := pageable.PageSize
// 		offset := (pageable.PageNumber - 1) * limit

// 		data = ListSub(data, offset, limit)
// 	}
// 	return data
// }

func List2Array(data *list.List) (ret []interface{}) {
	if data == nil {
		return
	}

	l := data.Len()
	ret = make([]interface{}, l)

	if l == 0 {
		return
	}

	b := data.Front()
	e := data.Back()
	i := 0
	for {
		ret[i] = b.Value
		i++

		if b == e {
			break
		}
		b = b.Next()
	}
	return
}
