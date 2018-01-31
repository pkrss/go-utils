package main

import (
	"log"
	"net/http"

	"github.com/pkrss/go-utils/conf"
	pkControllers "github.com/pkrss/go-utils/controllers"
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

////////////////////////////////////////////////////////////////////////////////
// one controller

type MainController struct {
	pkControllers.Controller
}

func (this *MainController) Get() {
	this.JsonResult("hello")
}

func (this *MainController) Test1() {
	id := this.GetString(":Id")
	this.JsonResult("Test1:" + id)
}

////////////////////////////////////////////////////////////////////////////////
// below is sample code, your can umcomment below comment when used.

func main() {

	conf.InitByConfigFile()

	profile.SetMyGetString(conf.GetString)

	pkControllers.AddRouter("/:Id/test1", &MainController{}, "get:Test1")
	pkControllers.AddRouter("/", &MainController{})

	// pkControllers.SetStaticPath("/s", "s")

	// pkRedis.InitRedis()

	pkControllers.UserController = &MyControllerUser{}

	port := profile.ProfileReadString("httpport", "8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
