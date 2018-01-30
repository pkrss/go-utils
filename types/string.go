package types

import "strconv"

func GetValueString(v interface{}) string {
	ret := ""

	if v == nil {
		return ret
	}

	switch v2 := v.(type) {
	case string:
		ret = v2
	case bool:
		ret = strconv.FormatBool(v2)
	case int:
		ret = strconv.Itoa(v2)
	case int64:
		ret = strconv.FormatInt(v2, 10)
	case float32:
		ret = strconv.FormatFloat(float64(v2), 'g', 30, 32)
	case float64:
		ret = strconv.FormatFloat(v2, 'g', 30, 64)
	default:
		ret = v.(string)
	}

	return ret
}
