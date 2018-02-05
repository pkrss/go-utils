package auth

import "github.com/pkrss/go-utils/mvc/controllers"

////////////////////////////////////////////////////////////////////////////////
// system auth part

const (
	User int = iota
	Admin
)

type UserContext struct {
	Role   int
	UserId string
}

type MyAuthImpl struct {
}

func (this *MyAuthImpl) TokenKey() string {
	return "X-PKRSS-SAMPLE"
}

func (this *MyAuthImpl) LoadUserToken(c controllers.ControllerInterface) string {

	token := this.Token

	if token != "" {
		return token
	}

	tokenKey := this.TokenKey()

	return c.CookieValue(k)
}

func (this *MyAuthImpl) LoadTokenObj(token string) interface{} {
	return &UserContext{Role: Admin, UserId: "1"} // json.Unmarshal(redis.Get(token))
}

func (this *MyAuthImpl) SaveTokenObj(token string, obj interface{}) {
	// redis.Set(token, json.Marshal(obj))
}

func (this *MyAuthImpl) CheckUserPrivilege(userContext interface{}, requiredPrivilege interface{}) bool {
	return userContext.(*UserContext).Role >= requiredPrivilege.(int)
}

func (this *MyAuthImpl) IsClientManagerOrSelf(userContext interface{}, targetUserId interface{}) bool {
	return userContext.(*UserContext).Role >= Admin || userContext.(*UserContext).UserId == targetUserId.(string)
}

// system auth part
////////////////////////////////////////////////////////////////////////////////
