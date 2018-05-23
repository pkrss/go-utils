package routers

import (
	"net/http"

	"github.com/pkrss/go-utils/examples/mvc/complex/controllers"

	pkRouters "github.com/pkrss/go-utils/mvc/routers"
)

func AddRouter(contextPath string) http.Handler {
	pkRouters.AddRouterOptSlash(contextPath+"/users/", &controllers.UserController{})
	pkRouters.AddRouter(contextPath+"/users/all", &controllers.UserController{}, "get:All")
	pkRouters.AddRouter(contextPath+"/users/like", &controllers.UserController{}, "get:Like")
	pkRouters.AddRouter(contextPath+"/users/login", &controllers.UserController{}, "get:Login")
	pkRouters.AddRouter(contextPath+"/users/{id:.+}", &controllers.UserIdController{})
	pkRouters.AddRouter(contextPath+"/users/{id:.+}/detail", &controllers.UserIdController{}, "get:Detail")

	pkRouters.AddRouter(contextPath+"/authenticates", &controllers.AuthenticatesController{}, "post:Login")
	return pkRouters.AddRouter(contextPath+"/authenticates/pings", &controllers.AuthenticatesController{}, "post:Pings")
}
