package pqsql

import (
	"errors"
	"reflect"

	"github.com/pkrss/go-utils/beans"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"

	"github.com/go-pg/pg/types"
	pkReflect "github.com/pkrss/go-utils/reflect"
)

type BaseDaoInterface interface {
	CreateModelObject() BaseModelInterface
	// create type is: *[]BaseModel
	CreateModelSlice(len int, cap int) interface{}
	FindOneById(id interface{}) (BaseModelInterface, error)
	FindOneByFilter(col string, val interface{}, selCols ...string) (BaseModelInterface, error)
	UpdateByFilter(ob BaseModelInterface, col string, val interface{}, structColsParams ...[]string) error
	UpdateById(ob BaseModelInterface, id interface{}, structColsParams ...[]string) error
	Insert(ob BaseModelInterface, structColsParams ...[]string) error
	SelectSelSqlList(partSql string, pageable *beans.Pageable, userData interface{}, cb SelectListCallback) (resultListPointer interface{}, total int64, e error)
	DeleteOneById(id interface{}) error
	DeleteByFilter(col string, val interface{}) error
}

type BaseDao struct {
	ObjModel BaseModelInterface
	ObjType  reflect.Type
	Db       *pg.DB
}

func CreateBaseDao(v BaseModelInterface) (dao BaseDaoInterface) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}
	ret := BaseDao{}
	ret.ObjModel = v
	ret.ObjType = reflect.TypeOf(val.Interface())

	dbTable := orm.Tables.Get(ret.ObjType)
	dbTable.Name = types.Q(v.TableName())

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

func (this *BaseDao) GetDb() (*pg.DB, error) {
	db := this.Db
	if db == nil {
		db = Db
	}
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return db, nil
}

func (this *BaseDao) FindOneById(id interface{}) (BaseModelInterface, error) {
	return this.FindOneByFilter(this.ObjModel.IdColumn(), id)
}

func (this *BaseDao) FindOneByFilter(col string, val interface{}, selCols ...string) (BaseModelInterface, error) {

	db, err := this.GetDb()
	if err != nil {
		return nil, err
	}

	obj := this.CreateModelObject()

	selSql := this.ObjModel.SelSql()
	if selSql == "" {

		models := db.Model(obj).Where(col+" = ?", val)
		if len(selCols) > 0 {
			models = models.Column(selCols...)
		}

		err = models.Select(obj)
	} else {
		_, err = db.QueryOne(obj, selSql+" WHERE "+col+" = ?", val)
	}

	if err != nil {
		return nil, err
	}

	obj.FilterValue()

	return obj, err
}

func (this *BaseDao) UpdateByFilter(ob BaseModelInterface, col string, val interface{}, structColsParams ...[]string) error {
	db, err := this.GetDb()
	if err != nil {
		return err
	}

	models := db.Model(ob).Where(col+" = ?", val)
	selCols := pkReflect.GetStructFieldNames(ob, structColsParams...)
	if len(selCols) > 0 {
		models = models.Column(selCols...)
	}

	_, err = models.Update(ob)
	if err != nil {
		return err
	}

	return err
}

func (this *BaseDao) DeleteOneById(id interface{}) error {
	return this.DeleteByFilter(this.ObjModel.IdColumn(), id)
}

func (this *BaseDao) DeleteByFilter(col string, val interface{}) error {
	db, err := this.GetDb()
	if err != nil {
		return err
	}

	models := db.Model(this.ObjModel).Where(col+" = ?", val)

	_, err = models.Delete()

	return err
}

func (this *BaseDao) UpdateById(ob BaseModelInterface, id interface{}, structColsParams ...[]string) error {
	return this.UpdateByFilter(ob, ob.IdColumn(), id, structColsParams...)
}

func (this *BaseDao) Insert(ob BaseModelInterface, structColsParams ...[]string) error {
	db, err := this.GetDb()
	if err != nil {
		return err
	}

	models := db.Model(ob)
	selCols := pkReflect.GetStructFieldNames(ob, structColsParams...)
	if len(selCols) > 0 {
		models = models.Column(selCols...)
	}

	_, err = models.Insert(ob)
	if err != nil {
		return err
	}

	return err
}

type SelectListCallback func(listRawHelper *ListRawHelper)

func (this *BaseDao) SelectSelSqlList(partSql string, pageable *beans.Pageable, userData interface{}, cb SelectListCallback) (resultListPointer interface{}, total int64, e error) {

	if this.ObjModel.TableName() == "" {
		e = errors.New("tableName is empty")
		return
	}

	resultListPointer = this.CreateModelSlice(0, 0)
	db, err := this.GetDb()
	if err != nil {
		e = err
		return
	}

	listRawHelper := ListRawHelper{}
	listRawHelper.Pageable = pageable
	listRawHelper.WhereArgs = make([]interface{}, 0)
	listRawHelper.ObjModel = this.ObjModel
	listRawHelper.Db = db
	listRawHelper.UserData = userData

	if cb != nil {
		cb(&listRawHelper)
	}

	total, e = listRawHelper.SelSqlListQuery(partSql, resultListPointer)

	return

}
