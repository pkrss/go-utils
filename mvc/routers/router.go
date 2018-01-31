package routers

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/pkrss/go-utils/mvc/controllers"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	c := controllers.Controller{W: w, R: r}
	c.AjaxError("Not implement!")
}

type Route struct {
	ContollerType reflect.Type
	MethodStr     string
}

func (this *Route) Handler(w http.ResponseWriter, r *http.Request, urlPathParameters map[string]string) {
	reqMethod := strings.ToLower(r.Method)

	var methodName string

	for f := true; f; f = false {
		if len(this.MethodStr) == 0 {
			break
		}

		methodStrAry := strings.Split(this.MethodStr, ":")
		if len(methodStrAry) > 1 {
			method := strings.ToLower(methodStrAry[0])
			if reqMethod != method {
				NotFoundHandler(w, r)
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

	objType := reflect.New(this.ContollerType)
	obj := objType.Elem().Addr().Interface().(controllers.ControllerInterface)
	obj.SetResponseWriter(w)
	obj.SetRequest(r)
	obj.SetUrlParameters(urlPathParameters)

	m := objType.MethodByName(methodName)

	if !m.IsValid() || m.IsNil() {
		NotFoundHandler(w, r)
		return
	}

	m.Call([]reflect.Value{})
}
