package controllers

import (
	"github.com/pkrss/go-utils/examples/mvc/complex/auth"
	"github.com/pkrss/go-utils/examples/mvc/complex/models"
	"github.com/pkrss/go-utils/examples/mvc/complex/services"

	"github.com/pkrss/go-utils/mvc/controllers"
)

type AuthenticatesController struct {
	controllers.AuthController
}

func (this *AuthenticatesController) Login() {
	var rqt auth.UserLoginRequest
	e := this.RequestBodyToJsonObject(&rqt)
	if e != nil {
		this.AjaxError(e.Error())
		return
	}
	rqt.ClientIp = this.GetClientIpAddr()
	rqt.Ctl = this
	rsp, err := services.Authenticate(&rqt)
	if err != nil {
		models.SyswarnAdd("login error:"+err.Error(), rqt.ClientIp)
		this.AjaxError(err.Error())
		return
	}

	this.AjaxDbRecord(rsp, true)
}

func (this *AuthenticatesController) Pings() {
	ok := this.CheckUserPrivilege(auth.EGuest)
	if !ok {
		this.AjaxError("not logined!")
		return
	}

	c := this.LoadUserContext().(*models.UserContext)

	auth.RefreshUserContext(c)

	this.AjaxDbRecord(nil, true)
}
