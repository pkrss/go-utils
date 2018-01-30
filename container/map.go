package container

import (
	"encoding/json"
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

	for key, val := range s {
		d[key] = val
	}

	by, err := json.Marshal(&d)
	if err != nil {
		return ""
	}

	return string(by)
}
