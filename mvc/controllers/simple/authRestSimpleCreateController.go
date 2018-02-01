package simple

import (
	"hx98/base/constants"
	"reflect"
	"sx98/admin/models"

	"github.com/pkrss/go-utils/pqsql"
)

type SimpleAuthRestListCreateParams struct {
	RecordList      *[]pqsql.BaseModelInterface
	RecordModel     pqsql.BaseModelInterface
	SelectListCbFun SelectListCallback
	SelectSql       string
}

type SimpleAuthRestListCreateController struct {
	base.SimpleAuthRestController
	Params *SimpleAuthRestCreateParams
}

func (this *SimpleAuthRestListCreateController) OnPrepare() {
	this.Model = this.Params.RecordModel
	this.SimpleAuthRestController.OnPrepare()
}

func (this *SimpleAuthRestListCreateController) Get() {
	v := reflect.ValueOf(this.Params.RecordList)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := reflect.TypeOf(v.Interface())
	reflect.MakeSlice()
	newObj := reflect.New(t)
	obj := newObj.Elem().Addr().Interface().(BaseModelInterface)

	var l []models.Menu
	this.SimpleAuthRestController.Helper.OnGetListWithPrivilege(constants.EClientAccount, &l, this.Params.SelectSql, menuGetList)
}

func (this *SimpleAuthRestListCreateController) Post() {
	this.SimpleAuthRestController.Helper.OnPostWithPrivilege(constants.EClientAccount, nil, []string{"Id", "CreateTime"})
}

func CreateSimpleRestController(params *SimpleAuthRestCreateParams) *SimpleAuthRestCreateController {
	ret := SimpleAuthRestCreateController{Params: params}

	return ret
}
