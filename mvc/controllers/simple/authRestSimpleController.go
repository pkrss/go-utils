package simple

import (
	"reflect"

	"github.com/pkrss/go-utils/mvc/controllers"
	"github.com/pkrss/go-utils/orm"
	pkReflect "github.com/pkrss/go-utils/reflect"
)

func CreateListRestController(params *ListRestParams) *ListAuthRestController {
	return &ListAuthRestController{Params: params}
}

func CreateItemRestController(params *ItemRestParams) *ItemAuthRestController {
	return &ItemAuthRestController{Params: params}
}

//////////////////////////////////////////////////////////////////////////

type ListRestParams struct {
	RecordModel orm.BaseModelInterface

	SelectListCbFun orm.SelectListCallback
	SelectSql       string
	SelectPrivilege interface{}

	PostPrivilege        interface{}
	PostStructColsParams *pkReflect.StructSelCols

	OnRestDbCbFun OnListRestDbCallback
}

type ListRestParamsGettor interface {
	GetParams() *ListRestParams
}

type ListAuthRestController struct {
	controllers.AuthController
	Model  orm.BaseModelInterface
	Helper *ListAuthRestHelper
	Params *ListRestParams
}

func (this *ListAuthRestController) GetParams() *ListRestParams {
	return this.Params
}

func (this *ListAuthRestController) OnPrepare() {
	this.Model = this.Params.RecordModel
	this.Helper = CreateListAuthRestHelper(this, this.Model)
}
func (this *ListAuthRestController) OnLeave() {
	this.Helper = nil
	this.Model = nil
}

func (this *ListAuthRestController) Get() {
	if this.Params.SelectPrivilege != nil {
		this.Helper.OnGetListWithPrivilege(this.Params.SelectPrivilege, this.Params.SelectSql, this.Params.SelectListCbFun)
	} else {
		this.Helper.OnGetList(this.Params.SelectSql, this.Params.SelectListCbFun)
	}
}

func (this *ListAuthRestController) Post() {
	if this.Params.PostPrivilege != nil {
		this.Helper.OnPostWithPrivilege(this.Params.PostPrivilege, this.Params.PostStructColsParams)
	} else {
		this.Helper.OnPost(this.Params.PostStructColsParams)
	}
}

//////////////////////////////////////////////////////////////////////////

type ItemRestParams struct {
	RecordModel orm.BaseModelInterface

	SelectPrivilege interface{}
	IdUrlParam      string
	IdType          reflect.Kind

	PutPrivilege        interface{}
	PutStructColsParams *pkReflect.StructSelCols

	DeletePrivilege interface{}
	OnRestDbCbFun   OnItemRestDbCallback
}

type ItemRestParamsGettor interface {
	GetParams() *ItemRestParams
}

type ItemAuthRestController struct {
	controllers.AuthController
	Model  orm.BaseModelInterface
	Helper *ItemAuthRestHelper
	Params *ItemRestParams
}

func (this *ItemAuthRestController) GetParams() *ItemRestParams {
	return this.Params
}

func (this *ItemAuthRestController) OnPrepare() {
	this.Model = this.Params.RecordModel
	this.Helper = CreateItemAuthRestHelper(this, this.Model)
}

func (this *ItemAuthRestController) OnLeave() {
	this.Helper = nil
	this.Model = nil
}

func (this *ItemAuthRestController) Get() {
	if this.Params.SelectPrivilege != nil {
		this.Helper.OnGetOneWithPrivilege(this.Params.SelectPrivilege, this.Params.IdUrlParam, this.Params.IdType)
	} else {
		this.Helper.OnGetOne(this.Params.IdUrlParam, this.Params.IdType)
	}
}

func (this *ItemAuthRestController) Put() {
	if this.Params.PutPrivilege != nil {
		this.Helper.OnPutWithPrivilege(this.Params.PutPrivilege, this.Params.IdUrlParam, this.Params.IdType, this.Params.PutStructColsParams)
	} else {
		this.Helper.OnPut(this.Params.IdUrlParam, this.Params.IdType, this.Params.PutStructColsParams)
	}
}

func (this *ItemAuthRestController) Patch() {
	this.Put()
}

func (this *ItemAuthRestController) Delete() {
	if this.Params.DeletePrivilege != nil {
		this.Helper.OnDeleteWithPrivilege(this.Params.DeletePrivilege, this.Params.IdUrlParam, this.Params.IdType)
	} else {
		this.Helper.OnDelete(this.Params.IdUrlParam, this.Params.IdType)
	}
}
