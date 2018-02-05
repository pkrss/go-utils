package main

import (
	"log"
	"net/http"

	"github.com/pkrss/go-utils/conf"
	pkControllers "github.com/pkrss/go-utils/mvc/controllers"
	pkRouters "github.com/pkrss/go-utils/mvc/routers"
	"github.com/pkrss/go-utils/profile"
)

////////////////////////////////////////////////////////////////////////////////
// system parse token

type MyControllerUser struct {
}

func (this *MyControllerUser) TokenKey() string {
	return "X-PKRSS-SAMPLE"
}

func (this *MyControllerUser) LoadTokenObj(token string) interface{} {
	return nil // json.Unmarshal(redis.Get(token))
}

func (this *MyControllerUser) SaveTokenObj(token string, obj interface{}) {
	// redis.Set(token, json.Marshal(obj))
}

func (this *MyControllerUser) CheckUserPrivilege(userContext interface{}, requiredPrivilege interface{}) bool {
	return false
}

func (this *MyControllerUser) IsClientManagerOrSelf(userContext interface{}, targetUserId interface{}) bool {
	return false
}

////////////////////////////////////////////////////////////////////////////////
// one controller

type MainController struct {
	pkControllers.Controller
}

// access: http://localhost:8080/
func (this *MainController) Get() {
	this.JsonResult("hello")
}

// access: http://localhost:8080/1/test1
func (this *MainController) Test1() {
	id := this.GetString(":id")
	this.JsonResult("Test1 id:" + id)
}

// access: http://localhost:8080/conf?q=runmode
// http://localhost:8080/conf?q=sys.smscode.hack
func (this *MainController) Conf() {
	q := this.GetString("q")
	this.JsonResult("conf [" + q + "]=" + profile.ProfileReadString(q))
}

////////////////////////////////////////////////////////////////////////////////
// below is sample code, your can umcomment below comment when used.

func main() {

	conf.InitByConfigFile()

	profile.SetMyGetString(conf.GetString)

	pkRouters.AddRouter("/{id:\\d+}/test1", &MainController{}, "get:Test1")
	r := pkRouters.AddRouter("/", &MainController{})
	pkRouters.AddRouter("/conf", &MainController{}, "get:Conf")

	// access: http://localhost:8080/s/
	pkRouters.SetStaticPath("/s", "s")

	pkControllers.DefaultUserInterface = &MyControllerUser{}

	port := profile.ProfileReadString("httpport", "8080")
	localAddr := "127.0.0.1:" + port // 0.0.0.0
	log.Printf("Server bind in %s\n", localAddr)
	log.Fatal(http.ListenAndServe(localAddr, r))
}
