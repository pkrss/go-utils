package controllers

import (
	"log"

	"github.com/pkrss/go-utils/redis"
)

///////////////////////////////////////////////////////////

type ControllerUserInterface interface {
	TokenKey() string
	LoadTokenObj(token string) interface{}
	SaveTokenObj(token string, obj interface{})
}

var DefaultUserInterface ControllerUserInterface

///////////////////////////////////////////////////////////

type AuthControllerInterface interface {
}

type AuthController struct {
	Controller

	UserContext   interface{}
	Token         string
	UserInterface ControllerUserInterface
}

func (this *AuthController) GetUserController() ControllerUserInterface {
	if this.UserInterface != nil {
		return this.UserInterface
	}
	return DefaultUserInterface
}

func (this *AuthController) LoadUserToken() string {

	token := this.Token

	if token != "" {
		return token
	}

	u := this.GetUserController()
	if u == nil {
		return ""
	}
	tokenKey := u.TokenKey()

	for f := true; f; f = false {

		k := tokenKey

		if token != "" {
			break
		}

		token = this.Header(k)

		if token != "" {
			break
		}

		token = this.Query(k)

		if token != "" {
			break
		}

		token = this.CookieValue(k)

		if token != "" {
			break
		}
	}

	this.Token = token

	return token
}

//读取登录状态
func (this *AuthController) LoadUserContext() interface{} {
	if this.UserContext != nil {
		return this.UserContext
	}

	token := this.LoadUserToken()

	if token == "" {
		return nil
	}

	u := this.GetUserController()
	if u == nil {
		return nil
	}

	this.UserContext = u.LoadTokenObj(token)

	return this.UserContext
}

//登录状态验证
func (this *BaseController) CheckUserPrivilege(requiredPrivilege int) bool {
	userContext := this.LoadUserContext()

	if hxToken.CheckUserPrivilege(userContext, requiredPrivilege) {
		return true
	}

	log.Printf("401 token=%s userContext=%v str=%s \n", this.token, userContext, redis.GetCache(hxToken.CACHE_PREFIX+this.token))

	this.Abort("401")

	return false
}
func (this *BaseController) CheckUserIsClientManagerOrSelf(targetUserId string) bool {
	userContext := this.LoadUserContext()

	if hxToken.IsClientManagerOrSelf(userContext, targetUserId) {
		return true
	}

	log.Printf("401 token=%s userContext=%v str=%s \n", this.token, userContext, redis.GetCache(hxToken.CACHE_PREFIX+this.token))

	this.Abort("401")

	return false
}
