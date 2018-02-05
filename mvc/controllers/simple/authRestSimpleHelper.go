package simple

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/pkrss/go-utils/mvc/controllers"
	"github.com/pkrss/go-utils/orm"
	pkReflect "github.com/pkrss/go-utils/reflect"
)

type OpType int

const (
	_ OpType = iota
	BeforeGetList
	AfterGetList
	BeforePost
	AfterPost

	BeforeGet
	AfterGet
	BeforePut
	AfterPut
	BeforeDelete
	AfterDelete
)

type OnRestDbCallback func(op OpType, ob interface{}, dao orm.BaseDaoInterface, c controllers.ControllerInterface) error

type SimpleAuthRestHelper struct {
	Dao           orm.BaseDaoInterface
	C             *SimpleAuthRestController
	OnRestDbCbFun OnRestDbCallback
	OldCodeFormat bool
}

func CreateSimpleAuthRestHelper(c *SimpleAuthRestController, v orm.BaseModelInterface) (ret *SimpleAuthRestHelper) {
	d := orm.CreateBaseDao(v)
	ret = &SimpleAuthRestHelper{C: c, Dao: d}
	return
}

func (this *SimpleAuthRestHelper) OnGetList(selSql string, cb orm.SelectListCallback) {

	pageable := this.C.GetPageableFromRequest()

	if this.OnRestDbCbFun != nil {
		e := this.OnRestDbCbFun(BeforeGetList, nil, this.Dao, this.C)
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

	if err == nil && this.OnRestDbCbFun != nil {
		e := this.OnRestDbCbFun(AfterGetList, l, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(err.Error())
			return
		}
	}

	c := 0
	if l != nil {
		c = reflect.ValueOf(l).Len()
	}
	this.C.AjaxDbList(pageable, l, c, total, this.OldCodeFormat)
}

func (this *SimpleAuthRestHelper) OnGetListWithPrivilege(requiredPrivilege interface{}, selSql string, cb orm.SelectListCallback) {
	ok := this.C.CheckUserPrivilege(requiredPrivilege)
	if !ok {
		this.C.AjaxError("权限不足")
		return
	}

	this.OnGetList(selSql, cb)
}

func (this *SimpleAuthRestHelper) OnPost(structColsParams ...*pkReflect.StructSelCols) {

	ob := this.Dao.CreateModelObject()

	e := this.C.RequestBodyToJsonObject(ob)

	if e != nil {
		this.C.AjaxError(e.Error())
		return
	}

	if this.OnRestDbCbFun != nil {
		e := this.OnRestDbCbFun(BeforePost, ob, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(e.Error())
			return
		}
	}

	err := this.Dao.Insert(ob, structColsParams...)

	if err == nil && this.OnRestDbCbFun != nil {
		e := this.OnRestDbCbFun(AfterPost, ob, this.Dao, this.C)
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

func (this *SimpleAuthRestHelper) OnPostWithPrivilege(requiredPrivilege interface{}, structColsParams ...*pkReflect.StructSelCols) {
	ok := this.C.CheckUserPrivilege(requiredPrivilege)
	if !ok {
		this.C.AjaxError("权限不足")
		return
	}

	this.OnPost(structColsParams...)
}

func (this *SimpleAuthRestHelper) GetIdParam(k string, t reflect.Kind) (id interface{}, e error) {

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

func (this *SimpleAuthRestHelper) OnGetOne(k string, t reflect.Kind) {

	id, e := this.GetIdParam(k, t)

	if e != nil {
		this.C.AjaxError(e.Error())
		return
	}

	if this.OnRestDbCbFun != nil {
		e := this.OnRestDbCbFun(BeforeGet, nil, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(e.Error())
			return
		}
	}

	ob, err := this.Dao.FindOneById(id)

	if err == nil && this.OnRestDbCbFun != nil {
		e := this.OnRestDbCbFun(AfterGet, ob, this.Dao, this.C)
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

func (this *SimpleAuthRestHelper) OnGetOneWithPrivilege(requiredPrivilege interface{}, k string, t reflect.Kind) {

	ok := this.C.CheckUserPrivilege(requiredPrivilege)
	if !ok {
		this.C.AjaxError("权限不足")
		return
	}

	this.OnGetOne(k, t)
}

func (this *SimpleAuthRestHelper) OnPut(k string, t reflect.Kind, structColsParams ...*pkReflect.StructSelCols) {

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

	if this.OnRestDbCbFun != nil {
		pkReflect.SetStructFieldValue(ob, ob.IdColumn(), id)
		e := this.OnRestDbCbFun(BeforePut, ob, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(e.Error())
			return
		}
	}

	err := this.Dao.UpdateById(ob, id, structColsParams...)

	if err == nil && this.OnRestDbCbFun != nil {
		e := this.OnRestDbCbFun(AfterPut, ob, this.Dao, this.C)
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

func (this *SimpleAuthRestHelper) OnPutWithPrivilege(requiredPrivilege interface{}, k string, t reflect.Kind, structColsParams ...*pkReflect.StructSelCols) {

	ok := this.C.CheckUserPrivilege(requiredPrivilege)
	if !ok {
		this.C.AjaxUnAuthorized("权限不足")
		return
	}

	this.OnPut(k, t, structColsParams...)
}

func (this *SimpleAuthRestHelper) OnDelete(k string, t reflect.Kind) {

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

	if this.OnRestDbCbFun != nil {
		pkReflect.SetStructFieldValue(ob, ob.IdColumn(), id)
		e := this.OnRestDbCbFun(BeforeDelete, ob, this.Dao, this.C)
		if e != nil {
			this.C.AjaxError(err.Error())
			return
		}
	}

	err = this.Dao.DeleteOneById(id)

	if this.OnRestDbCbFun != nil {
		e := this.OnRestDbCbFun(AfterDelete, ob, this.Dao, this.C)
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

func (this *SimpleAuthRestHelper) OnDeleteWithPrivilege(requiredPrivilege interface{}, k string, t reflect.Kind) {

	ok := this.C.CheckUserPrivilege(requiredPrivilege)
	if !ok {
		this.C.AjaxError("权限不足")
		return
	}

	this.OnDelete(k, t)
}
