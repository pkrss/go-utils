package profile

import (
	"os"
	"strconv"
)

type MyGetString func(string) string

var gMyGetString MyGetString

func SetMyGetString(myGetString MyGetString) {
	gMyGetString = myGetString
}

func ProfileReadString(key string, def ...string) string {
	v := os.Getenv(key)

	if v == "" && gMyGetString != nil {
		v = gMyGetString(key)
	}

	if v == "" {
		if len(def) > 0 {
			v = def[0]
		}
	}
	return v
}

func ProfileReadInt(key string, def ...int) int {

	v := 0
	if len(def) > 0 {
		v = def[0]
	}

	s := ProfileReadString(key)
	if s != "" {
		v2, err := strconv.Atoi(s)
		if err == nil {
			v = v2
		}
	}

	return v
}

func ProfileReadBool(key string, def ...bool) bool {

	v := false
	if len(def) > 0 {
		v = def[0]
	}

	s := ProfileReadString(key)
	if s != "" {
		v2, err := strconv.ParseBool(s)
		if err == nil {
			v = v2
		}
	}

	return v
}
