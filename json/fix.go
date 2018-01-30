package json

import (
	"log"

	"github.com/robertkrimen/otto"
)

// var regFixJsonKey *regexp.Regexp

// func fixJsonKey2(data string) string {
// 	if regFixJsonKey == nil {
// 		regFixJsonKey = regexp.MustCompile(`(?:\s*['"]*)?([a-zA-Z0-9]+)(?:['"]*\s*)?:\s*[n'"\{\[1-9]`)
// 	}

// 	ret := regFixJsonKey.ReplaceAllStringFunc(data, func(a string) string {
// 		return `"` + a + `"`
// 	})

// 	fmt.Println(ret)

// 	return ret
// }

var vm *otto.Otto

func FixJsonKey(data string) string {
	if vm == nil {
		vm = otto.New()
	}
	value, err := vm.Run("JSON.stringify(" + data + ")")
	if err != nil {
		log.Println(err)
		return data
	}

	ret := value.String()

	return ret
}
