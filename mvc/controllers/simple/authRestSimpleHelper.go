package simple

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/pkrss/go-utils/orm"
	pkReflect "github.com/pkrss/go-utils/reflect"
)

type ListOpType int

const (
	_ ListOpType = iota
	BeforeGetList
	AfterGetList
	BeforePost
	AfterPost
)

type OnListRestDbCallback func(op ListOpType, ob interface{}, dao orm.BaseDaoInterface, c *ListAuthRestController) error

type ListAuthRestHelper struct {
	Dao           orm.BaseDaoInterface
	C             *ListAuthRestController
	OldCodeFormat bool
}

func CreateListAuthRestHelper(c *ListAuthRestController, v orm.BaseModelInterface) (ret *ListAuthRestHelper) {
	d := orm.CreateBaseDao(v)
	ret = &ListAuthRestHelper{C: c, Dao: d}
	return
}

func (this *ListAuthRestHelper) OnGetList(selSql string, cb orm.SelectListCallback) {

	pageable := this.C.GetPageableFromRequest()

	if this.C.Params.OnRestDbCbFun != nil {
		e := this.C.Params.OnRestDbCbFun(BeforeGetList, nil, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(e.Error())
			return
		}
	}

	l, total, err := this.Dao.SelectSelSqlList(selSql, pageable, this.C, cb)
	if err != nil {
		this.C.AjaxError(err.Error())
		return
	}

	if err == nil && this.C.Params.OnRestDbCbFun != nil {
		e := this.C.Params.OnRestDbCbFun(AfterGetList, l, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(err.Error())
			return
		}
	}

	c := 0
	if l != nil {
		c = reflect.ValueOf(l).Elem().Len()
	}
	this.C.AjaxDbList(pageable, l, c, total, this.OldCodeFormat)
}

func (this *ListAuthRestHelper) OnGetListWithPrivilege(requiredPrivilege interface{}, selSql string, cb orm.SelectListCallback) {
	ok := this.C.CheckUserPrivilege(requiredPrivilege)
	if !ok {
		this.C.AjaxError("Privilege error")
		return
	}

	this.OnGetList(selSql, cb)
}

func (this *ListAuthRestHelper) OnPost(structColsParams ...*pkReflect.StructSelCols) {

	ob := this.Dao.CreateModelObject()

	e := this.C.RequestBodyToJsonObject(ob)

	if e != nil {
		this.C.AjaxError(e.Error())
		return
	}

	if this.C.Params.OnRestDbCbFun != nil {
		e := this.C.Params.OnRestDbCbFun(BeforePost, ob, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(e.Error())
			return
		}
	}

	err := this.Dao.Insert(ob, structColsParams...)

	if err == nil && this.C.Params.OnRestDbCbFun != nil {
		e := this.C.Params.OnRestDbCbFun(AfterPost, ob, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(err.Error())
			return
		}
	}

	if err != nil {
		this.C.AjaxError(err.Error())
		return
	}

	this.C.AjaxDbRecord(ob, this.OldCodeFormat)
}

func (this *ListAuthRestHelper) OnPostWithPrivilege(requiredPrivilege interface{}, structColsParams ...*pkReflect.StructSelCols) {
	ok := this.C.CheckUserPrivilege(requiredPrivilege)
	if !ok {
		this.C.AjaxError("Privilege error")
		return
	}

	this.OnPost(structColsParams...)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ItemOpType int

const (
	_ ItemOpType = iota

	BeforeGet
	AfterGet
	BeforePut
	AfterPut
	BeforeDelete
	AfterDelete
)

type OnItemRestDbCallback func(op ItemOpType, ob interface{}, dao orm.BaseDaoInterface, c *ItemAuthRestController) error

type ItemAuthRestHelper struct {
	Dao           orm.BaseDaoInterface
	C             *ItemAuthRestController
	OnRestDbCbFun OnItemRestDbCallback
	OldCodeFormat bool
}

func CreateItemAuthRestHelper(c *ItemAuthRestController, v orm.BaseModelInterface) (ret *ItemAuthRestHelper) {
	d := orm.CreateBaseDao(v)
	ret = &ItemAuthRestHelper{C: c, Dao: d}
	return
}

func (this *ItemAuthRestHelper) GetIdParam(k string, t reflect.Kind) (id interface{}, e error) {

	s := this.C.Query(k)
	if s == "" {
		e = errors.New("id is empty")
		return
	}
	switch t {
	case reflect.Int64:
		id, e = strconv.ParseInt(s, 10, 64)
		if e != nil {
			return
		}
	case reflect.Int32:
		id, e = strconv.ParseInt(s, 10, 32)
		if e != nil {
			return
		}
	case reflect.String:
		id = s
	}

	if id == nil {
		e = errors.New("id is nil")
	}

	return
}

func (this *ItemAuthRestHelper) OnGetOne(k string, t reflect.Kind) {

	id, e := this.GetIdParam(k, t)

	if e != nil {
		this.C.AjaxError(e.Error())
		return
	}

	if this.C.Params.OnRestDbCbFun != nil {
		e := this.C.Params.OnRestDbCbFun(BeforeGet, nil, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(e.Error())
			return
		}
	}

	ob, err := this.Dao.FindOneById(id)

	if err == nil && this.C.Params.OnRestDbCbFun != nil {
		e := this.C.Params.OnRestDbCbFun(AfterGet, ob, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(e.Error())
			return
		}
	}

	if err != nil {
		this.C.AjaxError(err.Error())
		return
	}

	this.C.AjaxDbRecord(ob, this.OldCodeFormat)
}

func (this *ItemAuthRestHelper) OnGetOneWithPrivilege(requiredPrivilege interface{}, k string, t reflect.Kind) {

	ok := this.C.CheckUserPrivilege(requiredPrivilege)
	if !ok {
		this.C.AjaxError("Privilege error")
		return
	}

	this.OnGetOne(k, t)
}

func (this *ItemAuthRestHelper) OnPut(k string, t reflect.Kind, structColsParams ...*pkReflect.StructSelCols) {

	id, e := this.GetIdParam(k, t)

	if e != nil {
		this.C.AjaxError(e.Error())
		return
	}

	ob := this.Dao.CreateModelObject()

	e = this.C.RequestBodyToJsonObject(ob)
	if e != nil {
		this.C.AjaxError(e.Error())
		return
	}

	if this.C.Params.OnRestDbCbFun != nil {
		pkReflect.SetStructFieldValue(ob, ob.IdColumn(), id)
		e := this.C.Params.OnRestDbCbFun(BeforePut, ob, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(e.Error())
			return
		}
	}

	err := this.Dao.UpdateById(ob, id, structColsParams...)

	if err == nil && this.C.Params.OnRestDbCbFun != nil {
		e := this.C.Params.OnRestDbCbFun(AfterPut, ob, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(e.Error())
			return
		}
	}

	if err != nil {
		this.C.AjaxError(err.Error())
		return
	}

	this.C.AjaxDbRecord(ob, this.OldCodeFormat)
}

func (this *ItemAuthRestHelper) OnPutWithPrivilege(requiredPrivilege interface{}, k string, t reflect.Kind, structColsParams ...*pkReflect.StructSelCols) {

	ok := this.C.CheckUserPrivilege(requiredPrivilege)
	if !ok {
		this.C.AjaxUnAuthorized("Privilege error")
		return
	}

	this.OnPut(k, t, structColsParams...)
}

func (this *ItemAuthRestHelper) OnDelete(k string, t reflect.Kind) {

	id, e := this.GetIdParam(k, t)

	if e != nil {
		this.C.AjaxError(e.Error())
		return
	}

	ob, err := this.Dao.FindOneById(id)
	if err != nil {
		this.C.AjaxError(err.Error())
		return
	}

	if this.C.Params.OnRestDbCbFun != nil {
		pkReflect.SetStructFieldValue(ob, ob.IdColumn(), id)
		e := this.C.Params.OnRestDbCbFun(BeforeDelete, ob, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(err.Error())
			return
		}
	}

	err = this.Dao.DeleteOneById(id)

	if this.C.Params.OnRestDbCbFun != nil {
		e := this.C.Params.OnRestDbCbFun(AfterDelete, ob, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(e.Error())
			return
		}
	}

	if err != nil {
		this.C.AjaxError(err.Error())
		return
	}

	this.C.AjaxDbRecord(ob, this.OldCodeFormat)
}

func (this *ItemAuthRestHelper) OnDeleteWithPrivilege(requiredPrivilege interface{}, k string, t reflect.Kind) {

	ok := this.C.CheckUserPrivilege(requiredPrivilege)
	if !ok {
		this.C.AjaxError("Privilege error")
		return
	}

	this.OnDelete(k, t)
}
