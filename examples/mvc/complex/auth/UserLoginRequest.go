package auth

import (
	"github.com/pkrss/go-utils/mvc/controllers"
)

type UserLoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	ClientIp string `json:"clientIp"`
	Ctl      controllers.ControllerInterface
}
