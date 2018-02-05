package routers

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkrss/go-utils/mvc/controllers"
)

func addHandle(pattern string, c controllers.ControllerInterface, methodStr string) func(w http.ResponseWriter, r *http.Request) {

	val := reflect.ValueOf(c)
	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}
	contollerType := reflect.TypeOf(val.Interface())

	route := Route{ContollerType: contollerType, MethodStr: methodStr, ContollerObj: c}

	return func(w http.ResponseWriter, r *http.Request) {
		var p map[string]string
		if strings.Contains(pattern, ":") {
			p = mux.Vars(r)
		}
		route.Handler(w, r, p)
	}
}

var app *mux.Router

func getApp() *mux.Router {
	if app == nil {
		app = mux.NewRouter()
	}
	return app
}

/*
	pattern: "/user/{name:[a-z]+}/profile"
*/
func AddRouter(pattern string, c controllers.ControllerInterface, methodStrs ...string) http.Handler {

	app := getApp()

	methodStr := ""
	if len(methodStrs) > 0 {
		methodStr = methodStrs[0]
	}

	h := addHandle(pattern, c, methodStr)

	r := app.HandleFunc(pattern, h)

	if methodStrAry := strings.Split(methodStr, ":"); len(methodStrAry) > 1 {
		m := methodStrAry[0]
		m = strings.ToUpper(m)
		r.Methods(strings.Split(m, ",")...)
	}

	return app
}

func SetStaticPath(urlPattern string, fileLocalDir string) http.Handler {

	app := getApp()

	fsh := http.FileServer(http.Dir(fileLocalDir))
	fsh = http.StripPrefix(urlPattern, fsh)
	app.PathPrefix(urlPattern).Handler(fsh)

	return app
}
