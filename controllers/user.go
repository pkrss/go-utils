package controllers

type ControllerUserInterface interface {
	TokenKey() string
	LoadTokenObj(token string) interface{}
	SaveTokenObj(token string, obj interface{}) string
}

var UserController ControllerUserInterface
