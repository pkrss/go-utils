package controllers

import (
	"errors"
	"fmt"
	"reflect"
	"sx98/sys/constants"
	"time"

	"github.com/pkrss/go-utils/examples/restAuth/models"
	base "github.com/pkrss/go-utils/mvc/controllers"
	"github.com/pkrss/go-utils/mvc/controllers/simple"
	"github.com/pkrss/go-utils/orm"
	pkReflect "github.com/pkrss/go-utils/reflect"
	"github.com/pkrss/go-utils/uuid"
)

type UserController struct {
	simple.SimpleAuthRestListCreateController
}

func (this *UserController) OnPrepare() {
	postStructColsParams := &pkReflect.StructSelCols{
		ExcludeCols: []string{"Id", "Password", "CreateTime"},
	}
	params := simple.SimpleAuthRestListCreateParams{
		RecordModel: &models.User{}, SelectListCbFun: userGetList,
		SelectPrivilege: Admin, PostPrivilege: Admin, PostStructColsParams: postStructColsParams,
		OnRestDbCbFun: userListOnRestDbCallback,
	}

	this.SimpleAuthRestListCreateController.Params = &params
	this.SimpleAuthRestListCreateController.OnPrepare()
	this.SimpleAuthRestController.Helper.OldCodeFormat = true
}

func userGetList(listRawHelper *orm.ListRawHelper) error {
	pageable := listRawHelper.Pageable
	if len(pageable.Sort) == 0 {
		pageable.Sort = "-id"
	}

	listRawHelper.SetCondArrLike("q", "mobile", "nick_name")
	listRawHelper.SetCondArrEqu("role", orm.Int, "role")
	return nil
}

func userListOnRestDbCallback(op simple.OpType, ob interface{}, dao orm.BaseDaoInterface, c base.ControllerInterface) error {
	switch op {
	case simple.BeforePost:
		record := ob.(*models.User)

		var newUserDomain *models.User

		if "quickly" == c.GetString("reg") {
			i := 1
			for i < 10000 {
				mobile := fmt.Sprintf(constants.VestMobile_Prefix+"%08d", i)

				_, err := dao.FindOneByFilter("mobile", newUserDomain.Mobile, &pkReflect.StructSelCols{
					IncludeCols: []string{"id"},
				})
				if err != nil {
					newUserDomain = &models.User{}
					newUserDomain.Mobile = mobile
					newUserDomain.Password = constants.VestMd5Password
					break
				}
				i = i + 1
			}
		} else {
			newUserDomain := record

			_, err := dao.FindOneByFilter("mobile", newUserDomain.Mobile, &pkReflect.StructSelCols{
				IncludeCols: []string{"id"},
			})
			if err == nil {
				return errors.New("手机号已存在，请更换注册手机号")
			}
		}

		userContext := c.(base.AuthControllerInterface).LoadUserContext().(*UserContext)
		if newUserDomain.Role > userContext.Role {
			return errors.New("越权操作")
		}

		if newUserDomain.Id.IsNil() {
			newUserDomain.Id.Set(uuid.UuidCreate())
		}

		now := time.Now()

		if newUserDomain.CreateTime.IsNil() {
			newUserDomain.CreateTime.Set(now)
		}

	}
	return nil
}

func (this *UserController) Get() {
	this.SimpleAuthRestListCreateController.Get()
}

func (this *UserController) Post() {
	this.SimpleAuthRestListCreateController.Post()
}

func (this *UserController) Login() {
	this.RenderViewSimple("users/login.html")
}

type UserIdController struct {
	simple.SimpleAuthRestCreateController
}

func (this *UserIdController) OnPrepare() {
	putStructColsParams := &pkReflect.StructSelCols{
		ExcludeCols: []string{"Id", "Password", "CreateTime"},
	}
	params := simple.SimpleAuthRestCreateParams{
		RecordModel: &models.User{}, IdUrlParam: ":id", IdType: reflect.String,
		PutStructColsParams: putStructColsParams, DeletePrivilege: Admin,
		OnRestDbCbFun: userOnRestDbCallback,
	}

	this.SimpleAuthRestCreateController.Params = &params
	this.SimpleAuthRestCreateController.OnPrepare()
	this.SimpleAuthRestController.Helper.OldCodeFormat = true
}

func (this *UserIdController) Get() {
	this.SimpleAuthRestCreateController.Get()
}

func (this *UserIdController) Detail() {
	this.Get()
}

func (this *UserIdController) Put() {
	this.Put()
}

func (this *UserIdController) Login() {
	authImpl := this.GetUserAuthImpl()
	if authImpl == nil {
		return ""
	}

	this.SetCookieValue(authImpl.TokenKey(), "your created token value")
}

func (this *UserIdController) Delete() {
	this.AjaxError("禁止该操作")
}

func userOnRestDbCallback(op simple.OpType, ob interface{}, dao orm.BaseDaoInterface, c base.ControllerInterface) error {
	switch op {
	case simple.AfterGet:
		record := ob.(*models.User)

		if !c.(base.AuthControllerInterface).CheckUserIsClientManagerOrSelf(record.Id.String()) {
			return errors.New("无权操作!")
		}

		// mobile := record.Mobile

		// models.UserFilterPublicFields(&record)

		// record.Mobile = mobile
	case simple.BeforePut:
		record := ob.(*models.User)

		putStructColsParams := c.(*UserIdController).GetParams().PutStructColsParams

		denyKeys := []string{"Password", "UpdateTime"}
		putStructColsParams.ExcludeCols = append(putStructColsParams.ExcludeCols, denyKeys...)

		userContext := c.(base.AuthControllerInterface).LoadUserContext().(*UserContext)
		if userContext.Role >= Admin {
			// ob.Password = ""
		} else if record.Id.String() == userContext.UserId {
			userCanNotModifyKeys := []string{"InviteCode", "Denied", "Role", "Vip", "FansCount", "FollowsCount"}
			putStructColsParams.ExcludeCols = append(putStructColsParams.ExcludeCols, userCanNotModifyKeys...)
		} else {
			return errors.New("禁止跨用户操作")
		}
	}
	return nil
}
