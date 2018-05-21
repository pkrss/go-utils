package services

import (
	"errors"
	"strings"
	"time"

	"github.com/pkrss/go-utils/examples/mvc/complex/models"

	"github.com/pkrss/go-utils/crypto"
	"github.com/pkrss/go-utils/orm"
	pkTime "github.com/pkrss/go-utils/time"
)

func Authenticate(rqt *models.UserLoginRequest) (ret map[string]interface{}, e error) {
	var user *models.User

	if rqt == nil || rqt.UserName == "" || rqt.Password == "" {
		e = errors.New("parameter error")
		return
	}

	dao := orm.CreateBaseDao(&models.User{})
	record, err := dao.FindOneByFilter("user_name", rqt.UserName)
	if err != nil {
		e = err
		return
	}

	user = record.(*models.User)

	password := rqt.Password

	now := time.Now()
	if pkTime.CheckSamePeriod("1m", user.LastErrorLoginTime, now) {
		if user.LastErrorLoginTimer > 3 {
			e = errors.New("password error to much timer,please retry after 1 miniute!")
			return
		}
	}

	password = strings.ToLower(password)
	psw := strings.ToLower(user.Password)

	if password != psw {
		password2 := crypto.Md5(password)

		if password2 != psw {

			user.LastErrorLoginTime = now
			user.LastErrorLoginTimer++

			dao.UpdateByFilter(user, "last_error_login_timer", user.LastErrorLoginTimer)
			dao.UpdateByFilter(user, "last_error_login_time", user.LastErrorLoginTime)

			e = errors.New("password is error!")
			return
		}
	}

	if user.Denied {
		e = errors.New(rqt.UserName + " is denied")
		return
	}

	ext := make(map[string]interface{})
	userContext, err := models.CreateTokenContext(user, rqt, ext)
	if err != nil {
		e = err
		return
	}

	ret = make(map[string]interface{})

	ret["user"] = user
	ret["token"] = userContext.ID

	user.FilterValue()

	return
}
