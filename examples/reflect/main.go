package main

import (
	"log"

	"github.com/pkrss/go-utils/reflect"
)

type A struct {
	AVal string
}

type AA struct {
	A
	AAVal string
}

func testCopy() {
	var aa1 AA
	var aa2 AA
	aa2.A.AVal = "bbbb"
	reflect.CopyStruct(&aa1, aa2)
	log.Printf("aa1=%v\n", aa1)
}

func main() {
	testCopy()
}
