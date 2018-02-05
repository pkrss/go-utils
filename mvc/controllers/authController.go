package controllers

///////////////////////////////////////////////////////////

type AuthImplInterface interface {
	TokenKey() string
	LoadUserToken(c ControllerInterface) string
	LoadTokenObj(token string) interface{}
	SaveTokenObj(token string, obj interface{})
	CheckUserPrivilege(userContext interface{}, requiredPrivilege interface{}) bool
	IsClientManagerOrSelf(userContext interface{}, targetUserId interface{}) bool
}

var DefaultAuthImpl AuthImplInterface

///////////////////////////////////////////////////////////

type AuthControllerInterface interface {
	GetUserAuthImpl() AuthImplInterface
	LoadUserToken() string
	LoadUserContext() interface{}
	CheckUserPrivilege(requiredPrivilege interface{}) bool
	CheckUserIsClientManagerOrSelf(targetUserId interface{}) bool
}

type AuthController struct {
	Controller

	UserContext interface{}
	Token       string
	UserAuthObj AuthImplInterface
}

func (this *AuthController) GetUserAuthImpl() AuthImplInterface {
	if this.UserAuthObj != nil {
		return this.UserAuthObj
	}
	return DefaultAuthImpl
}

func (this *AuthController) LoadUserToken() string {

	token := this.Token

	if token != "" {
		return token
	}

	authImpl := this.GetUserAuthImpl()
	if authImpl == nil {
		return ""
	}

	token = authImpl.LoadUserToken(this)

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

	authImpl := this.GetUserAuthImpl()
	if authImpl == nil {
		return nil
	}

	this.UserContext = authImpl.LoadTokenObj(token)

	return this.UserContext
}

//登录状态验证
func (this *AuthController) CheckUserPrivilege(requiredPrivilege interface{}) bool {
	userContext := this.LoadUserContext()

	if userContext == nil {
		return false
	}

	authImpl := this.GetUserAuthImpl()
	if authImpl == nil {
		return false
	}

	if authImpl.CheckUserPrivilege(userContext, requiredPrivilege) {
		return true
	}

	// log.Printf("401 token=%s userContext=%v str=%s \n", this.token, userContext, redis.GetCache(hxToken.CACHE_PREFIX+this.token))

	// this.Abort("401")

	return false
}
func (this *AuthController) CheckUserIsClientManagerOrSelf(targetUserId interface{}) bool {
	userContext := this.LoadUserContext()

	if userContext == nil {
		return false
	}

	authImpl := this.GetUserAuthImpl()
	if authImpl == nil {
		return false
	}

	if authImpl.IsClientManagerOrSelf(userContext, targetUserId) {
		return true
	}

	// log.Printf("401 token=%s userContext=%v str=%s \n", this.token, userContext, redis.GetCache(hxToken.CACHE_PREFIX+this.token))

	// this.Abort("401")

	return false
}
