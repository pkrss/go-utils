package auth

import (
	"strconv"
	"time"

	"github.com/pkrss/go-utils/examples/mvc/complex/models"
	"github.com/pkrss/go-utils/mvc/controllers"
	"github.com/pkrss/go-utils/orm"
)

func CreateTokenContext(user *models.User, rqt *UserLoginRequest, ext ...map[string]interface{}) (ret *models.UserContext, e error) {
	ret = &models.UserContext{}
	ret.UserId = user.ID
	ret.UserName = user.UserName
	ret.Ip = rqt.ClientIp
	ret.CreateTime = time.Now()
	ret.UpdateTime = ret.CreateTime

	dao := orm.CreateBaseDao(ret)
	e = dao.Insert(ret)
	if e != nil {
		return
	}

	if rqt.Ctl != nil {
		rqt.Ctl.SetCookieValue(HTTP_HEAER_TOKEN, strconv.FormatInt(ret.ID, 10), 3*24*60*60)
	}

	return
}

func LoadContextFromToken(id int64) (ret *models.UserContext, e error) {
	ret = &models.UserContext{}
	dao := orm.CreateBaseDao(ret)
	i, err := dao.FindOneById(id)
	if err != nil {
		e = err
		return
	}

	ret = i.(*models.UserContext)
	return
}

func DeleteContext(uc *models.UserContext, Ctl controllers.ControllerInterface) {
	Ctl.SetCookieValue(HTTP_HEAER_TOKEN, "", -1)

	dao := orm.CreateBaseDao(uc)
	dao.DeleteOneById(uc.ID)
}

func RefreshUserContext(uc *models.UserContext) {
	dao := orm.CreateBaseDao(uc)
	dao.UpdateByFilter(uc, "update_time", time.Now())
}
