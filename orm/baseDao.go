package orm

import (
	"errors"
	"reflect"
	"strings"

	"github.com/pkrss/go-utils/beans"

	pkContainer "github.com/pkrss/go-utils/container"
	pkReflect "github.com/pkrss/go-utils/reflect"
)

type BaseDaoInterface interface {
	CreateModelObject() BaseModelInterface

	CreateModelSlice(len int, cap int) interface{} // create type is: *[]BaseModel
	FindOneById(id interface{}) (BaseModelInterface, error)
	FindOneByFilter(col string, val interface{}, structColsParams ...*pkReflect.StructSelCols) (BaseModelInterface, error)
	UpdateByFilter(ob BaseModelInterface, col string, val interface{}, structColsParams ...*pkReflect.StructSelCols) error
	UpdateById(ob BaseModelInterface, id interface{}, structColsParams ...*pkReflect.StructSelCols) error
	Insert(ob BaseModelInterface, structColsParams ...*pkReflect.StructSelCols) error
	SelectSelSqlList(partSql string, pageable *beans.Pageable, userData interface{}, cb SelectListCallback) (resultListPointer interface{}, total int64, e error)
	DeleteOneById(id interface{}) error
	DeleteByFilter(col string, val interface{}) error
}

type BaseDao struct {
	ObjModel   BaseModelInterface
	ObjType    reflect.Type
	OrmAdapter OrmAdapterInterface
}

func CreateBaseDao(v BaseModelInterface, ormAdapters ...OrmAdapterInterface) (dao BaseDaoInterface) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}
	ret := BaseDao{}
	ret.ObjModel = v
	ret.ObjType = reflect.TypeOf(val.Interface())

	var ormAdapter OrmAdapterInterface
	if len(ormAdapters) > 0 {
		ormAdapter = ormAdapters[0]
	} else {
		ormAdapter = DefaultOrmAdapter
	}
	ret.OrmAdapter = ormAdapter

	ormAdapter.RegModel(v)

	return &ret
}
func (this *BaseDao) CreateModelObject() BaseModelInterface {
	objType := reflect.New(this.ObjType)
	obj := objType.Elem().Addr().Interface().(BaseModelInterface)
	return obj
}
func (this *BaseDao) CreateModelSlice(len int, cap int) interface{} {
	// Create a slice to begin with
	myType := this.ObjType // reflect.TypeOf(this.ObjModel)
	slice := reflect.MakeSlice(reflect.SliceOf(myType), len, cap)

	// Create a pointer to a slice value and set it to the slice
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)
	return x.Elem().Addr().Interface().(interface{})
}

func (this *BaseDao) FindOneById(id interface{}) (BaseModelInterface, error) {
	return this.FindOneByFilter(this.ObjModel.IdColumn(), id)
}

func (this *BaseDao) FindOneByFilter(col string, val interface{}, structColsParams ...*pkReflect.StructSelCols) (BaseModelInterface, error) {

	obj := this.CreateModelObject()

	selSql := this.ObjModel.SelSql()

	var selCols []string
	if len(structColsParams) > 0 {
		selCols = pkReflect.GetStructFieldNames(obj, structColsParams...)
	}

	var err error

	if selSql == "" {
		selSql = "SELECT "

		if len(selCols) == 0 {
			selSql += "*"
		} else {
			selSql += strings.Join(selCols, ",")
		}

		selSql += " FROM " + obj.TableName()
	}

	sql := selSql + " WHERE " + col + " = ?"

	err = this.OrmAdapter.QueryOneBySql(obj, sql, val)

	if err == nil && obj == nil {
		err = errors.New("query one record is nil")
	}

	if err != nil {
		return nil, err
	}

	obj.FilterValue()

	return obj, err
}

func (this *BaseDao) UpdateByFilter(ob BaseModelInterface, col string, val interface{}, structColsParams ...*pkReflect.StructSelCols) error {
	dbField2Values := getObDbFieldsAndValues(ob, structColsParams...)
	c := len(dbField2Values)
	if c == 0 {
		return errors.New("No fields need update!")
	}

	values := make([]interface{}, c)

	sql := "UPDATE " + ob.TableName() + " SET "
	i := 0
	for k, v := range dbField2Values {
		sql += k + "=?"
		sql += ","
		values[i] = v
		i++
	}

	if strings.HasSuffix(sql, ",") {
		sql = sql[0 : len(sql)-1]
	}

	sql += " WHERE " + col + "=" + "?"
	values = append(values, val)

	return this.OrmAdapter.ExecSql(sql, values...)
}

func (this *BaseDao) DeleteOneById(id interface{}) error {
	return this.DeleteByFilter(this.ObjModel.IdColumn(), id)
}

func (this *BaseDao) DeleteByFilter(col string, val interface{}) error {
	sql := "DELETE " + this.ObjModel.TableName() + " WHERE " + col + " = ?"
	return this.OrmAdapter.ExecSql(sql, val)
}

func (this *BaseDao) UpdateById(ob BaseModelInterface, id interface{}, structColsParams ...*pkReflect.StructSelCols) error {
	return this.UpdateByFilter(ob, ob.IdColumn(), id, structColsParams...)
}
func (this *BaseDao) Insert(ob BaseModelInterface, structColsParams ...*pkReflect.StructSelCols) error {

	dbField2Values := getObDbFieldsAndValues(ob, structColsParams...)
	c := len(dbField2Values)
	if c == 0 {
		return errors.New("No fields need insert!")
	}

	keys := pkContainer.MapKeys(dbField2Values)
	values := pkContainer.MapValues(dbField2Values)
	t := make([]string, c)
	for i := 0; i < c; i++ {
		t[i] = "?"
	}

	sql := "INSERT INTO " + ob.TableName() + " (" + strings.Join(keys, ",") + ") VALUES(" + strings.Join(t, ",") + ")"

	if ob.IdColumn() != "" {

		returnSql := this.OrmAdapter.SqlReturnSql()
		idVal := pkReflect.GetStructField(ob, ob.IdColumn(), false)

		if returnSql != "" && idVal.IsValid() {
			returnSql = strings.Replace(returnSql, "{id}", ob.IdColumn(), -1)

			sql += returnSql

			e := this.OrmAdapter.QueryOneBySql(idVal.Addr().Interface(), sql, values...)
			// if e == nil {
			// 	pkReflect.SetStructFieldValue(ob, ob.IdColumn(), idVal.Interface())
			// }
			return e

		}
	}

	return this.OrmAdapter.ExecSql(sql, values...)
}

type SelectListCallback func(listRawHelper *ListRawHelper) error

func (this *BaseDao) SelectSelSqlList(partSql string, pageable *beans.Pageable, userData interface{}, cb SelectListCallback) (resultListPointer interface{}, total int64, e error) {

	if this.ObjModel.TableName() == "" {
		e = errors.New("tableName is empty")
		return
	}

	resultListPointer = this.CreateModelSlice(0, 0)

	listRawHelper := ListRawHelper{}
	listRawHelper.Pageable = pageable
	listRawHelper.WhereArgs = make([]interface{}, 0)
	listRawHelper.ObjModel = this.ObjModel
	listRawHelper.OrmAdapter = this.OrmAdapter
	listRawHelper.UserData = userData

	if cb != nil {
		e = cb(&listRawHelper)
		if e != nil {
			return
		}
	}

	total, e = listRawHelper.SelSqlListQuery(partSql, resultListPointer)

	return

}
