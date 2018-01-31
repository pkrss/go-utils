package simple

import (
	"errors"
	"fmt"

	"github.com/pkrss/go-utils/mvc/controllers"
)

type SimpleAuthUserInterface interface {
	CheckUserPrivilege(userContext interface{}, requiredPrivilege interface{}) bool
	IsManagerOrSelf(userContext interface{}, userId interface{}) bool
}

var DefaultSimpleAuthUserInterface SimpleAuthUserInterface

type SimpleAuthController struct {
	controllers.AuthController
	AuthUserInterface SimpleAuthUserInterface
}

func (this *SimpleAuthController) GetAuthUserInterface() SimpleAuthUserInterface {
	if this.AuthUserInterface != nil {
		return this.AuthUserInterface
	}
	return DefaultSimpleAuthUserInterface
}

func (this *SimpleAuthController) CheckUserPrivilege(requiredPrivilege interface{}) error {
	u := this.GetAuthUserInterface()
	if u == nil {
		return errors.New("AuthUserInterface is nil")
	}
	userContext := this.LoadUserContext()
	if userContext == nil {
		return errors.New("User is not logined")
	}

	if u.CheckUserPrivilege(userContext, requiredPrivilege) {
		return nil
	}

	return fmt.Errorf("User is not authorized with privilege:%v", requiredPrivilege)
}
func (this *SimpleAuthController) CheckUserIsClientManagerOrSelf(targetUserId interface{}) error {
	u := this.GetAuthUserInterface()
	if u == nil {
		return errors.New("AuthUserInterface is nil")
	}
	userContext := this.LoadUserContext()
	if userContext == nil {
		return errors.New("User is not logined")
	}

	if u.IsClientManagerOrSelf(userContext, targetUserId) {
		return true
	}

	return errors.New("User is not authorized with manager and not self")
}
