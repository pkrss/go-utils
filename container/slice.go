package container

import (
	"reflect"
)

func CreateSlice(myType reflect.Type, len int, cap int) interface{} { // return is: *[]myType
	// Create a slice to begin with
	slice := reflect.MakeSlice(reflect.SliceOf(myType), len, cap)

	// Create a pointer to a slice value and set it to the slice
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)
	return x.Elem().Addr().Interface().(interface{})
}
