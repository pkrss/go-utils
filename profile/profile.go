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

	var v string
	if gMyGetString != nil {
		v = gMyGetString(key)
	} else {
		v, _ = os.LookupEnv(key)
	}

	if v == "" {
		if len(def) > 0 {
			v = def[0]
		}
	}
	return v
}

func ProfileReadFloat64(key string, def ...float64) float64 {
	v := 0.0
	if len(def) > 0 {
		v = def[0]
	}

	s := ProfileReadString(key)
	if s != "" {
		v2, err := strconv.ParseFloat(s, 64)
		if err == nil {
			v = v2
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
