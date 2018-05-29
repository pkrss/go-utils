package container

import (
	"encoding/json"
	"strings"

	"github.com/pkrss/go-utils/types"
)

func MapFromString(s string) map[string]interface{} {
	var m map[string]interface{}
	if s == "" {
		return m
	}
	json.Unmarshal([]byte(s), &m)
	return m
}

func MapStringMerge(dest string, src string) string {
	d := MapFromString(dest)
	s := MapFromString(src)

	if len(s) == 0 {
		return dest
	}

	MapMerge(d, s)

	by, err := json.Marshal(&d)
	if err != nil {
		return ""
	}

	return string(by)
}

func MapMerge(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {

	if len(src) == 0 {
		return dest
	}

	for key, val := range src {
		dest[key] = val
	}

	return dest
}

func MapGetValue(m map[string]interface{}, field string) (ret interface{}) {

	if len(m) == 0 || field == "" {
		return ""
	}

	ks := strings.Split(field, ".")

	var i int
	for c := len(ks); i < c; i++ {
		k := ks[i]

		if k == "" {
			break
		}

		v, ok := m[k]
		if !ok {
			break
		}

		if i == c-1 {
			// ret = types.GetValueString(v)
			ret = v
			break
		}

		m = v.(map[string]interface{})
		if m == nil {
			break
		}

	}

	return ret
}

func MapGetStringValue(m map[string]interface{}, field string) (ret string) {
	v := MapGetValue(m, field)
	if v != nil {
		ret = types.GetValueString(v)
	}
	return
}

func MapKeys(m map[string]interface{}) []string {
	ret := make([]string, 0)
	for k, _ := range m {
		ret = append(ret, k)
	}
	return ret
}

func MapValues(m map[string]interface{}) []interface{} {
	ret := make([]interface{}, 0)
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}
