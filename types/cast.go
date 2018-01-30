package types

import (
	"strconv"
)

func CastToInt64(o interface{}) int64 {
	switch v := o.(type) {
	case int64:
		return v
	case *int64:
		return *v
	case int:
		return int64(v)
	case *int:
		return int64(*v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	case *string:
		i, _ := strconv.ParseInt(*v, 10, 64)
		return i
	}
	return 0
}
