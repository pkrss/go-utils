package controllers

import (
	"net/http"
	"reflect"
	"strings"
)

func AddRouter(pattern string, c ControllerInterface, methodStr ...string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		reqMethod := strings.ToLower(r.Method)

		var methodName string

		for f := true; f; f = false {
			if len(methodStr) == 0 {
				break
			}

			methodStrAry := strings.Split(methodStr[0], ":")
			if len(methodStrAry) > 1 {
				method := strings.ToLower(methodStrAry[0])
				if reqMethod != method {
					c.AjaxError("Not implement method:" + method)
					return
				}
				methodName = methodStrAry[1]
			} else if len(methodStrAry) == 1 {
				methodName = methodStrAry[0]
			}

			if methodName == "" {
				break
			}
		}

		if methodName == "" {
			methodName = string(reqMethod[0]-'a'+'A') + reqMethod[1:]
		}

		val := reflect.ValueOf(c)
		switch val.Kind() {
		case reflect.Ptr:
			val = val.Elem()
		}
		objType := reflect.New(reflect.TypeOf(val.Interface()))
		obj := objType.Elem().Addr().Interface().(ControllerInterface)
		obj.SetResponseWriter(w)
		obj.SetRequest(r)
		// objType = reflect.ValueOf(obj)
		// m := pkReflect.GetStructMethod(obj, methodName)

		m := objType.MethodByName(methodName)

		if !m.IsValid() || m.IsNil() {
			c.AjaxError("Not implement method:" + methodName)
			return
		}

		m.Call([]reflect.Value{})
	})
}

func SetStaticPath(urlPattern string, fileLocalDir string) {
	fsh := http.FileServer(http.Dir(fileLocalDir))
	fsh = http.StripPrefix(urlPattern, fsh)
	http.Handle(urlPattern, fsh)
}
