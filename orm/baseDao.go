package orm

import (
	"errors"
	"reflect"
	"strings"

	"github.com/pkrss/go-utils/beans"
	pkContainer "github.com/pkrss/go-utils/container"
	"github.com/pkrss/go-utils/orm/inner"
	pkReflect "github.com/pkrss/go-utils/reflect"
	pkStrings "github.com/pkrss/go-utils/strings"
)

type BaseDaoInterface interface {
	CreateModelObject() BaseModelInterface

	CreateModelSlice(len int, cap int) interface{} // create type is: *[]BaseModel
	FindOneById(id interface{}) (BaseModelInterface, error)
	FindOneByFilter(col string, val interface{}, structColsParams ...*pkReflect.StructSelCols) (BaseModelInterface, error)
	FindOneByFilters(colVals map[string]interface{}, structColsParams ...*pkReflect.StructSelCols) (BaseModelInterface, error)
	FindOneBySql(selSql string, val ...interface{}) (BaseModelInterface, error)
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
	return pkContainer.CreateSlice(this.ObjType, len, cap)
}

func (this *BaseDao) FindOneById(id interface{}) (BaseModelInterface, error) {
	return this.FindOneByFilter(this.ObjModel.IdColumn(), id)
}

func (this *BaseDao) FindOneByFilter(col string, val interface{}, structColsParams ...*pkReflect.StructSelCols) (BaseModelInterface, error) {
	kv := make(map[string]interface{})
	kv[col] = val
	return this.FindOneByFilters(kv, structColsParams...)
}

func (this *BaseDao) FindOneByFilters(colVals map[string]interface{}, structColsParams ...*pkReflect.StructSelCols) (BaseModelInterface, error) {
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

	sql := selSql + " WHERE "

	first := true
	vals := make([]interface{}, 0)
	for k, v := range colVals {
		if first {
			first = false
		} else {
			sql += ","
		}
		col := pkStrings.StringToCamelCase(k)
		sql += col + " = ?"
		vals = append(vals, v)
	}

	err = this.OrmAdapter.QueryOneBySql(obj, sql, vals...)

	if err == nil && obj == nil {
		err = errors.New("query one record is nil")
	}

	if err != nil {
		return nil, err
	}

	return obj, err
}

func (this *BaseDao) FindOneBySql(selSql string, val ...interface{}) (ret BaseModelInterface, e error) {
	ret = this.CreateModelObject()

	e = this.OrmAdapter.QueryOneBySql(ret, selSql, val...)

	if e != nil {
		return
	}

	return

}

func (this *BaseDao) UpdateByFilter(ob BaseModelInterface, col string, val interface{}, structColsParams ...*pkReflect.StructSelCols) error {
	idCol := this.getRealIdCol(ob.IdColumn())

	dbField2Values := inner.GetStructDbFieldsAndValues(ob, idCol, true, structColsParams...)
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

	col = pkStrings.StringToCamelCase(col)

	sql += " WHERE " + col + "=" + "?"
	values = append(values, val)

	return this.OrmAdapter.ExecSql(sql, values...)
}

func (this *BaseDao) DeleteOneById(id interface{}) error {
	idCol := this.getRealIdCol(this.ObjModel.IdColumn())
	return this.DeleteByFilter(idCol, id)
}

func (this *BaseDao) DeleteByFilter(col string, val interface{}) error {
	col = pkStrings.StringToCamelCase(col)
	sql := "DELETE FROM " + this.ObjModel.TableName() + " WHERE " + col + " = ?"
	return this.OrmAdapter.ExecSql(sql, val)
}

func (this *BaseDao) UpdateById(ob BaseModelInterface, id interface{}, structColsParams ...*pkReflect.StructSelCols) error {
	idCol := this.getRealIdCol(this.ObjModel.IdColumn())
	return this.UpdateByFilter(ob, idCol, id, structColsParams...)
}
func (this *BaseDao) getRealIdCol(idColumn string) string {
	if strings.Contains(idColumn, ".") {
		ss := strings.Split(idColumn, ".")
		c := len(ss)
		if c > 0 {
			idColumn = ss[c-1]
		}
	}
	return idColumn
}
func (this *BaseDao) Insert(ob BaseModelInterface, structColsParams ...*pkReflect.StructSelCols) error {
	idCol := this.getRealIdCol(this.ObjModel.IdColumn())
	dbField2Values := inner.GetStructDbFieldsAndValues(ob, idCol, true, structColsParams...)
	c := len(dbField2Values)
	if c == 0 {
		return errors.New("No fields need insert!")
	}

	sqlKeys := ""
	sqlKeys2 := ""
	values := make([]interface{}, c)
	i := 0
	for k, v := range dbField2Values {
		sqlKeys += k + ","
		sqlKeys2 += "?,"
		values[i] = v
		i++
	}

	if strings.HasSuffix(sqlKeys, ",") {
		sqlKeys = sqlKeys[0 : len(sqlKeys)-1]
		sqlKeys2 = sqlKeys2[0 : len(sqlKeys2)-1]
	}

	sql := "INSERT INTO " + ob.TableName() + " (" + sqlKeys + ") VALUES(" + sqlKeys2 + ")"

	if idCol != "" {

		returnSql := this.OrmAdapter.SqlReturnSql()
		idVal := pkReflect.GetStructField(ob, idCol, false)

		if returnSql != "" && idVal.IsValid() {
			returnSql = strings.Replace(returnSql, "{id}", idCol, -1)

			sql += returnSql

			e := this.OrmAdapter.QueryOneBySql(idVal.Addr().Interface(), sql, values...)
			// if e == nil {
			// 	pkReflect.SetStructFieldValue(ob, idCol, idVal.Interface())
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
