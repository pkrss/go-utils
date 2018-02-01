package orm

import (
	"errors"
	"reflect"

	"github.com/pkrss/go-utils/beans"

	pkReflect "github.com/pkrss/go-utils/reflect"
)

type BaseDaoInterface interface {
	CreateModelObject() BaseModelInterface
	// create type is: *[]BaseModel
	CreateModelSlice(len int, cap int) interface{}
	FindOneById(id interface{}) (BaseModelInterface, error)
	FindOneByFilter(col string, val interface{}, structColsParams ...[]string) (BaseModelInterface, error)
	UpdateByFilter(ob BaseModelInterface, col string, val interface{}, structColsParams ...[]string) error
	UpdateById(ob BaseModelInterface, id interface{}, structColsParams ...[]string) error
	Insert(ob BaseModelInterface, structColsParams ...[]string) error
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

func (this *BaseDao) FindOneByFilter(col string, val interface{}, structColsParams ...[]string) (BaseModelInterface, error) {

	obj := this.CreateModelObject()

	selSql := this.ObjModel.SelSql()

	var selCols []string
	if len(structColsParams) > 0 {
		selCols = pkReflect.GetStructFieldNames(obj, structColsParams...)
	}

	var err error
	if selSql == "" {
		err = this.OrmAdapter.FindOneByCond(obj, col+" = ?", []interface{}{val}, selCols)
	} else {
		err = this.OrmAdapter.FindOneBySql(obj, selSql+" WHERE "+col+" = ?", val)
	}

	if err != nil {
		return nil, err
	}

	obj.FilterValue()

	return obj, err
}

func (this *BaseDao) UpdateByFilter(ob BaseModelInterface, col string, val interface{}, structColsParams ...[]string) error {
	var selCols []string
	if len(structColsParams) > 0 {
		selCols = pkReflect.GetStructFieldNames(ob, structColsParams...)
	}

	return this.OrmAdapter.UpdateByCond(ob, col+" = ?", []interface{}{val}, selCols)
}

func (this *BaseDao) DeleteOneById(id interface{}) error {
	return this.DeleteByFilter(this.ObjModel.IdColumn(), id)
}

func (this *BaseDao) DeleteByFilter(col string, val interface{}) error {
	return this.OrmAdapter.DeleteByCond(this.ObjModel, col+" = ?", val)
}

func (this *BaseDao) UpdateById(ob BaseModelInterface, id interface{}, structColsParams ...[]string) error {
	return this.UpdateByFilter(ob, ob.IdColumn(), id, structColsParams...)
}

func (this *BaseDao) Insert(ob BaseModelInterface, structColsParams ...[]string) error {
	var selCols []string
	if len(structColsParams) > 0 {
		selCols = pkReflect.GetStructFieldNames(ob, structColsParams...)
	}

	return this.OrmAdapter.Insert(ob, selCols...)
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
