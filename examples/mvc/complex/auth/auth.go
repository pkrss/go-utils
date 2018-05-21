package auth

import (
	"strconv"

	"github.com/pkrss/go-utils/examples/mvc/complex/models"
	"github.com/pkrss/go-utils/mvc/controllers"
)

type MyAuthImpl struct {
}

func (this *MyAuthImpl) TokenKey() string {
	return HTTP_HEAER_TOKEN
}

func (this *MyAuthImpl) LoadUserToken(c controllers.ControllerInterface) string {

	var token string

	tokenKey := this.TokenKey()

	for f := true; f; f = false {

		k := tokenKey

		if token != "" {
			break
		}

		token = c.Header(k)

		if token != "" {
			break
		}

		token = c.GetString(k)

		if token != "" {
			break
		}

		token = c.CookieValue(k)

		if token != "" {
			break
		}
	}

	return token
}

func (this *MyAuthImpl) LoadTokenObj(token string) interface{} {
	id, e := strconv.ParseInt(token, 10, 64)
	if e != nil {
		return nil
	}
	uc, e := LoadContextFromToken(id)
	if e != nil {
		return nil
	}
	return uc
}

func (this *MyAuthImpl) SaveTokenObj(token string, obj interface{}) {
	id, e := strconv.ParseInt(token, 10, 64)
	if e != nil {
		return
	}
	uc := obj.(*models.UserContext)
	uc.ID = id
	RefreshUserContext(uc)
}

func (this *MyAuthImpl) CheckUserPrivilege(userContext interface{}, requiredPrivilege interface{}) bool {
	return CheckUserPrivilege(userContext.(*models.UserContext), requiredPrivilege.(int))
}

func (this *MyAuthImpl) IsClientManagerOrSelf(userContext interface{}, targetUserId interface{}) bool {
	u2 := userContext.(*models.UserContext)
	t2 := targetUserId.(int64)
	return IsClientManagerOrSelf(u2, t2)
}
