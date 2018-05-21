package controllers

import (
	"errors"
	"reflect"
	"strings"

	"github.com/pkrss/go-utils/examples/mvc/complex/models"
	"github.com/pkrss/go-utils/mvc/controllers/simple"
	"github.com/pkrss/go-utils/orm"
	pkReflect "github.com/pkrss/go-utils/reflect"
)

type UserController struct {
	simple.ListAuthRestController
}

func (this *UserController) OnPrepare() {
	// postStructColsParams := &pkReflect.StructSelCols{
	// 	ExcludeCols: []string{"Id", "Password", "CreateTime"},
	// }
	params := simple.ListRestParams{
		RecordModel: &models.User{}, SelectListCbFun: userGetList,
		SelectPrivilege: models.EOperManagerStart, PostPrivilege: models.EOperManagerStart, PostStructColsParams: nil,
		OnRestDbCbFun: userListOnRestDbCallback,
	}

	this.ListAuthRestController.Params = &params
	this.ListAuthRestController.OnPrepare()
	this.ListAuthRestController.Helper.OldCodeFormat = true
}

func userGetList(listRawHelper *orm.ListRawHelper) error {
	pageable := listRawHelper.Pageable
	if len(pageable.Sort) == 0 {
		pageable.Sort = "-id"
	}

	listRawHelper.SetCondArrLike("q", "user_name")
	return nil
}

func userListOnRestDbCallback(op simple.ListOpType, ob interface{}, dao orm.BaseDaoInterface, c *simple.ListAuthRestController) error {
	switch op {
	case simple.BeforePost:
		record := ob.(*models.User)

		var newUserDomain *models.User

		newUserDomain = record

		_, err := dao.FindOneByFilter("user_name", newUserDomain.UserName, &pkReflect.StructSelCols{
			IncludeCols: []string{"id"},
		})
		if err == nil {
			return errors.New("user name is already exist!")
		}
	}
	return nil
}

func (this *UserController) Get() {
	this.ListAuthRestController.Get()
}

func (this *UserController) All() {
	this.Get()
}

func (this *UserController) Like() {
	this.Get()
}

func (this *UserController) Regs() {

	method := this.GetRequest().Method
	method = strings.ToLower(method)
	if method != "post" {
		this.AjaxError("only support post!")
		return
	}

	userName := this.GetString("userName")
	password := this.GetString("password")

	if userName == "" || password == "" {
		this.AjaxError("parameter is error")
		return
	}

	// ...
}

func (this *UserController) Post() {
	this.ListAuthRestController.Post()
}

func (this *UserController) Login() {
	this.RenderViewSimple("views/users/login.html")
}

type UserIdController struct {
	simple.ItemAuthRestController
}

func (this *UserIdController) OnPrepare() {
	putStructColsParams := &pkReflect.StructSelCols{
		ExcludeCols: []string{"Id", "Password", "CreateTime"},
	}
	params := simple.ItemRestParams{
		RecordModel: &models.User{}, IdUrlParam: ":id", IdType: reflect.String,
		PutStructColsParams: putStructColsParams, DeletePrivilege: models.EOperManagerStart,
		OnRestDbCbFun: userOnRestDbCallback,
	}

	this.ItemAuthRestController.Params = &params
	this.ItemAuthRestController.OnPrepare()
	this.ItemAuthRestController.Helper.OldCodeFormat = true
}

func (this *UserIdController) Get() {
	this.ItemAuthRestController.Get()
}

func (this *UserIdController) Detail() {
	this.Get()
}

func (this *UserIdController) Put() {
	this.ItemAuthRestController.Put()
}

func (this *UserIdController) Patch() {
	this.Put()
}

func (this *UserIdController) Delete() {
	this.AjaxError("deny this operator")
}

func userOnRestDbCallback(op simple.ItemOpType, ob interface{}, dao orm.BaseDaoInterface, c *simple.ItemAuthRestController) error {
	switch op {
	case simple.AfterGet:
		record := ob.(*models.User)

		if !c.CheckUserIsClientManagerOrSelf(record.ID) {
			return errors.New("not allow!")
		}

		// mobile := record.Mobile

		// models.UserFilterPublicFields(&record)

		// record.Mobile = mobile
	case simple.BeforePut:
		record := ob.(*models.User)

		putStructColsParams := c.GetParams().PutStructColsParams

		denyKeys := []string{"Password", "UpdateTime"}
		putStructColsParams.ExcludeCols = append(putStructColsParams.ExcludeCols, denyKeys...)

		userContext := c.LoadUserContext().(*models.UserContext)
		if models.CheckUserPrivilege(userContext, models.EOperManagerStart) {
			// ob.Password = ""
		} else if record.ID == userContext.UserId {
			userCanNotModifyKeys := []string{}
			putStructColsParams.ExcludeCols = append(putStructColsParams.ExcludeCols, userCanNotModifyKeys...)
		} else {
			return errors.New("not allowed")
		}
		// case simple.AfterPut:
		// record := ob.(*models.User)
		// refresh data to redis
	}
	return nil
}
