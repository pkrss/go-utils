package types

import (
	"strconv"
	"strings"
)

func Int64ArrayToString(a []int64) string {
	c := "["
	for i, b := range a {
		if i > 0 {
			c += ","
		}
		c += strconv.FormatInt(b, 10)
	}
	c += "]"
	return c
}

func Int64ArrayToInt32Array(a []int64) []int {
	var c []int
	for _, b := range a {
		c = append(c, int(b))
	}
	return c
}

func IntArrayToInt64Array(a []int) []int64 {
	var c []int64
	for _, b := range a {
		c = append(c, int64(b))
	}
	return c
}

func StringToInt64Array(a string) ([]int64, error) {
	var d []int64

	s := strings.Trim(a, "\"")
	s = strings.Replace(s, " ", "", -1)

	if s == "null" {
		return d, nil
	}
	f := strings.IndexAny(s, "[{")
	e := strings.LastIndexAny(s, "]}")
	if f < 0 || e < 0 || f >= e {
		return d, nil
	}
	s = s[f+1 : e]
	ss := strings.Split(s, ",")

	for _, sl := range ss {
		v, err := strconv.ParseInt(sl, 10, 64)
		if err == nil {
			d = append(d, v)
		} else {
			return d, err
		}
	}
	return d, nil
}

func StringToInt32Array(a string) ([]int, error) {
	b, e := StringToInt64Array(a)
	if e != nil {
		return []int{}, e
	}
	return Int64ArrayToInt32Array(b), e
}

func StringArraySub(a []string, b []string) []string {

	r := make([]string, 0)

	for _, key := range a {
		key2 := strings.ToLower(key)

		found := false
		for _, denyKey := range b {
			denyKey2 := strings.ToLower(denyKey)
			if denyKey2 == key2 {
				found = true
				break
			}
		}
		if found {
			continue
		}
		r = append(r, key)
	}
	return r
}
