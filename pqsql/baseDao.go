package pqsql

import (
	"errors"
	"hx98/base/beans"
	"reflect"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"

	"github.com/go-pg/pg/types"
	pkReflect "github.com/pkrss/go-utils/reflect"
)

type BaseDaoInterface interface {
	FindOneById(id interface{}) (BaseModelInterface, error)
	FindOneByFilter(col string, val interface{}, selCols ...string) (BaseModelInterface, error)
	UpdateByFilter(ob BaseModelInterface, col string, val interface{}, structColsParams ...[]string) error
	UpdateById(ob BaseModelInterface, id interface{}, structColsParams ...[]string) error
	Insert(ob BaseModelInterface, structColsParams ...[]string) error
	SelectListByRawHelper(listRawHelper ListRawHelper, cb ...SelectListCallback) (int64, error)
	SelectList(resultListPointer interface{}, pageable *beans.Pageable, cb ...SelectListCallback) (int64, error)
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

	objType := reflect.New(this.ObjType)
	obj := objType.Elem().Addr().Interface().(BaseModelInterface)

	models := db.Model(obj).Where(col+" = ?", val)
	if len(selCols) > 0 {
		models = models.Column(selCols...)
	}

	err = models.Select(obj)
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

type SelectListCallback func(listRawHelper ListRawHelper)

func (this *BaseDao) SelectListByRawHelper(listRawHelper ListRawHelper, cb ...SelectListCallback) (int64, error) {

	if listRawHelper.Query == nil {

		db, err := this.GetDb()
		if err != nil {
			return 0, err
		}
		listRawHelper.DbQuery = db.Model(this.ObjModel)
	}

	if len(cb) > 0 {
		cb[0](listRawHelper)
	}

	return listRawHelper.Query()
}

func (this *BaseDao) SelectList(resultListPointer interface{}, pageable *beans.Pageable, cb ...SelectListCallback) (int64, error) {
	// v := reflect.ValueOf(resultListPointer)
	// switch v.Kind() {
	// case reflect.Ptr:
	// 	v = v.Elem()
	// }

	tableName := this.ObjModel.TableName()
	// if v.Kind() == reflect.Slice {
	// 	newv := reflect.MakeSlice(v.Type(), 1, 1)
	// 	tableName = getTableName(newv.Index(0))
	// }

	if tableName == "" {
		return 0, errors.New("tableName is empty")
	}

	listRawHelper := MakeListRawHelper(resultListPointer, pageable)

	return this.SelectListByRawHelper(*listRawHelper, cb...)
}
