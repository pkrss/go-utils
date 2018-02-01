package simple

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"

	"github.com/pkrss/go-utils/orm"
)

type SimpleAuthRestHelper struct {
	Dao           orm.BaseDaoInterface
	C             SimpleAuthController
	OldCodeFormat bool
}

func CreateSimpleAuthRestHelper(c SimpleAuthUserInterface, v pqsql.BaseModelInterface) (ret SimpleAuthRestHelper) {
	d := orm.CreateBaseDao(v)
	ret = SimpleAuthRestHelper{C: c, Dao: d}
	return
}

func (this *SimpleAuthRestHelper) OnGetList(l *[]BaseModelInterface, selSql string, cb SelectListCallback) {

	pageable := this.GetPageableFromRequest()

	l, total, err := this.Dao.SelectSelSqlList(selSql, &pageable, this.C, cb)
	if err != nil {
		this.AjaxError(err.Error())
		return
	}

	this.C.AjaxDbList(pageable, l, len(l), total, this.OldCodeFormat)
}

func (this *SimpleAuthRestHelper) OnGetListWithPrivilege(requiredPrivilege interface{}, selSql string, cb SelectListCallback) {
	e := C.CheckUserPrivilege(requiredPrivilege)
	if e != nil {
		this.AjaxUnAuthorized(e.Error())
		return
	}

	this.OnGetList(selSql, cb)
}

func (this *SimpleAuthRestHelper) OnPost(structColsParams ...[]string) {

	ob := this.Dao.CreateModelObject()
	e := json.Unmarshal(this.Ctx.Input.RequestBody, ob)

	if e != nil {
		this.C.AjaxError(e)
		return
	}

	err := this.Dao.Insert(ob, structColsParams...)

	if err != nil {
		this.AjaxError(err.Error())
		return
	}

	this.C.AjaxDbRecord(ob, this.OldCodeFormat)
}

func (this *SimpleAuthRestHelper) OnPostWithPrivilege(requiredPrivilege interface{}, structColsParams ...[]string) {
	e := C.CheckUserPrivilege(requiredPrivilege)
	if e != nil {
		this.AjaxUnAuthorized(e.Error())
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
		this.C.AjaxError(e)
		return
	}

	ob, err := this.Dao.FindOneById(id)

	if err != nil {
		this.C.AjaxError(err.Error())
		return
	}

	this.C.AjaxDbRecord(ob, this.OldCodeFormat)
}

func (this *SimpleAuthRestHelper) OnGetOneWithPrivilege(requiredPrivilege interface{}, k string, t reflect.Kind) {

	e := C.CheckUserPrivilege(requiredPrivilege)
	if e != nil {
		this.AjaxUnAuthorized(e.Error())
		return
	}

	this.OnGetOne(k, t)
}

func (this *SimpleAuthRestHelper) OnPut(k string, t reflect.Kind, structColsParams ...[]string) {

	id, e := this.GetIdParam(k, t)

	if e != nil {
		this.C.AjaxError(e)
		return
	}

	ob := this.Dao.CreateModelObject()
	e = json.Unmarshal(this.Ctx.Input.RequestBody, ob)
	if e != nil {
		this.C.AjaxError(e)
		return
	}

	err := this.Dao.UpdateById(ob, id, structColsParams...)

	if err != nil {
		this.C.AjaxError(err.Error())
		return
	}

	this.C.AjaxDbRecord(ob, this.OldCodeFormat)
}

func (this *SimpleAuthRestHelper) OnPutWithPrivilege(requiredPrivilege interface{}, k string, t reflect.Kind, structColsParams ...[]string) {

	e := C.CheckUserPrivilege(requiredPrivilege)
	if e != nil {
		this.AjaxUnAuthorized(e.Error())
		return
	}

	this.OnPut(k, t, structColsParams...)
}

func (this *SimpleAuthRestHelper) OnDelete(k string, t reflect.Kind) {

	id, e := this.GetIdParam(k, t)

	if e != nil {
		this.C.AjaxError(e)
		return
	}

	ob, err := this.Dao.FindOneById(id)
	if err != nil {
		this.C.AjaxError(err.Error())
		return
	}

	err := this.Dao.DeleteOneById(id)

	if err != nil {
		this.C.AjaxError(err.Error())
		return
	}

	this.C.AjaxDbRecord(ob, this.OldCodeFormat)
}

func (this *SimpleAuthRestHelper) OnDeleteWithPrivilege(requiredPrivilege interface{}, k string, t reflect.Kind) {

	e := C.CheckUserPrivilege(requiredPrivilege)
	if e != nil {
		this.AjaxUnAuthorized(e.Error())
		return
	}

	this.OnDelete(k, t)
}
