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
	FindOneById(id interface{}) (BaseModelInterface, error)
	FindOneByFilter(col string, val interface{}, selCols ...string) (BaseModelInterface, error)
	UpdateByFilter(ob BaseModelInterface, col string, val interface{}, structColsParams ...[]string) error
	UpdateById(ob BaseModelInterface, id interface{}, structColsParams ...[]string) error
	Insert(ob BaseModelInterface, structColsParams ...[]string) error
	SelectListByRawHelper(listRawHelper *ListRawHelper, userData interface{}, cb ...SelectListCallback) (int64, error)
	SelectList(resultListPointer interface{}, pageable *beans.Pageable, userData interface{}, cb ...SelectListCallback) (int64, error)
	SelectSelSqlList(partSql string, resultListPointer []BaseModelInterface, pageable *beans.Pageable, userData interface{}, cb ...SelectListCallback) (int64, error)
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

func (this *BaseDao) SelectListByRawHelper(listRawHelper *ListRawHelper, userData interface{}, cb ...SelectListCallback) (int64, error) {

	if listRawHelper.DbQuery == nil {

		db, err := this.GetDb()
		if err != nil {
			return 0, err
		}
		listRawHelper.DbQuery = db.Model(this.ObjModel)
	}

	listRawHelper.UserData = userData

	if len(cb) > 0 {
		cb[0](listRawHelper)
	}

	return listRawHelper.Query()
}

func (this *BaseDao) SelectList(resultListPointer interface{}, pageable *beans.Pageable, userData interface{}, cb ...SelectListCallback) (int64, error) {

	if this.ObjModel.TableName() == "" {
		return 0, errors.New("tableName is empty")
	}

	listRawHelper := MakeListRawHelper(resultListPointer, pageable)

	return this.SelectListByRawHelper(listRawHelper, userData, cb...)
}

func (this *BaseDao) SelectSelSqlList(partSql string, resultListPointer []BaseModelInterface, pageable *beans.Pageable, userData interface{}, cb ...SelectListCallback) (int64, error) {

	if this.ObjModel.TableName() == "" {
		return 0, errors.New("tableName is empty")
	}

	listRawHelper := MakeListRawHelper(resultListPointer, pageable)
	listRawHelper.ObjModel = this.ObjModel
	listRawHelper.Db = this.Db
	listRawHelper.UserData = userData

	if listRawHelper.DbQuery == nil {

		db, err := this.GetDb()
		if err != nil {
			return 0, err
		}
		listRawHelper.DbQuery = db.Model(this.ObjModel)
	}

	if len(cb) > 0 {
		cb[0](listRawHelper)
	}

	return listRawHelper.SelSqlListQuery(partSql)

}
