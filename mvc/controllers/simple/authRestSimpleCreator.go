package simple

import (
	"reflect"

	"github.com/pkrss/go-utils/orm"
	pkReflect "github.com/pkrss/go-utils/reflect"
)

type SimpleAuthRestListCreateParams struct {
	RecordModel orm.BaseModelInterface

	SelectListCbFun orm.SelectListCallback
	SelectSql       string
	SelectPrivilege interface{}

	PostPrivilege        interface{}
	PostStructColsParams *pkReflect.StructSelCols

	OnRestDbCbFun OnRestDbCallback
}

type SimpleAuthRestListCreateController struct {
	SimpleAuthRestController
	Params *SimpleAuthRestListCreateParams
}

func (this *SimpleAuthRestListCreateController) OnPrepare() {
	this.SimpleAuthRestController.Model = this.Params.RecordModel
	this.SimpleAuthRestController.OnPrepare()
	this.SimpleAuthRestController.Helper.OnRestDbCbFun = this.Params.OnRestDbCbFun
}

func (this *SimpleAuthRestListCreateController) Get() {
	if this.Params.SelectPrivilege != nil {
		this.SimpleAuthRestController.Helper.OnGetListWithPrivilege(this.Params.SelectPrivilege, this.Params.SelectSql, this.Params.SelectListCbFun)
	} else {
		this.SimpleAuthRestController.Helper.OnGetList(this.Params.SelectSql, this.Params.SelectListCbFun)
	}
}

func (this *SimpleAuthRestListCreateController) Post() {
	if this.Params.PostPrivilege != nil {
		this.SimpleAuthRestController.Helper.OnPostWithPrivilege(this.Params.PostPrivilege, this.Params.PostStructColsParams)
	} else {
		this.SimpleAuthRestController.Helper.OnPost(this.Params.PostStructColsParams)
	}
}

func CreateSimpleListRestController(params *SimpleAuthRestListCreateParams) *SimpleAuthRestListCreateController {
	return &SimpleAuthRestListCreateController{Params: params}
}

type SimpleAuthRestCreateParams struct {
	RecordModel orm.BaseModelInterface

	SelectPrivilege interface{}
	IdUrlParam      string
	IdType          reflect.Kind

	PutPrivilege        interface{}
	PutStructColsParams *pkReflect.StructSelCols

	DeletePrivilege interface{}
	OnRestDbCbFun   OnRestDbCallback
}

type SimpleAuthRestCreateController struct {
	SimpleAuthRestController
	Params *SimpleAuthRestCreateParams
}

func (this *SimpleAuthRestCreateController) GetParams() *SimpleAuthRestCreateParams {
	return this.Params
}

func (this *SimpleAuthRestCreateController) OnPrepare() {
	this.Model = this.Params.RecordModel
	this.SimpleAuthRestController.OnPrepare()
}

func (this *SimpleAuthRestCreateController) Get() {
	if this.Params.SelectPrivilege != nil {
		this.SimpleAuthRestController.Helper.OnGetOneWithPrivilege(this.Params.SelectPrivilege, this.Params.IdUrlParam, this.Params.IdType)
	} else {
		this.SimpleAuthRestController.Helper.OnGetOne(this.Params.IdUrlParam, this.Params.IdType)
	}
}

func (this *SimpleAuthRestCreateController) Put() {
	if this.Params.PutPrivilege != nil {
		this.SimpleAuthRestController.Helper.OnPutWithPrivilege(this.Params.PutPrivilege, this.Params.IdUrlParam, this.Params.IdType, this.Params.PutStructColsParams)
	} else {
		this.SimpleAuthRestController.Helper.OnPut(this.Params.IdUrlParam, this.Params.IdType, this.Params.PutStructColsParams)
	}
}

func (this *SimpleAuthRestCreateController) Patch() {
	this.Put()
}

func (this *SimpleAuthRestCreateController) Delete() {
	if this.Params.DeletePrivilege != nil {
		this.SimpleAuthRestController.Helper.OnDeleteWithPrivilege(this.Params.DeletePrivilege, this.Params.IdUrlParam, this.Params.IdType)
	} else {
		this.SimpleAuthRestController.Helper.OnDelete(this.Params.IdUrlParam, this.Params.IdType)
	}
}

func CreateSimpleRestController(params *SimpleAuthRestCreateParams) *SimpleAuthRestCreateController {
	return &SimpleAuthRestCreateController{Params: params}
}
