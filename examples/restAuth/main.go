package main

import (
	"log"
	"net/http"

	"github.com/pkrss/go-utils/examples/restAuth/auth"
	myControllers "github.com/pkrss/go-utils/examples/restAuth/controllers"
	pkControllers "github.com/pkrss/go-utils/mvc/controllers"
	pkRouters "github.com/pkrss/go-utils/mvc/routers"
	orm "github.com/pkrss/go-utils/orm"
	"github.com/pkrss/go-utils/orm/pqsql"
)

///////////////////////////////////////////////////////////////////////////////////

func main() {

	// create db connect
	pqsql.Db = pqsql.CreatePgSql()
	orm.DefaultOrmAdapter = &pqsql.PgSqlAdapter{}
	pkControllers.DefaultAuthImpl = &auth.MyAuthImpl{}

	// mvc sample 1: used parameter create rest service
	pkRouters.AddRouter("/oauth/apps/{id:\\d+}", myControllers.CreateOAuthAppRestController())
	pkRouters.AddRouterOptSlash("/oauth/apps", myControllers.CreateOAuthAppListRestController())

	// mvc sample 2: used controller create rest service
	pkRouters.AddRouter("/users/{id:\\.+}", &myControllers.UserIdController{})
	pkRouters.AddRouterOptSlash("/users", &myControllers.UserController{})
	pkRouters.AddRouterOptSlash("/users/login", &myControllers.UserController{}, "get:Login")

	// start server
	localAddr := "127.0.0.1:8080"
	log.Printf("Server bind in %s\n", localAddr)
	log.Fatal(http.ListenAndServe(localAddr, pkRouters.GetApp()))

}
