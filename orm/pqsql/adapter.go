package pqsql

import (
	"errors"
	"reflect"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/types"
	pkOrm "github.com/pkrss/go-utils/orm"
)

type PgSqlAdapter struct {
	Db *pg.DB
}

func (this *PgSqlAdapter) RegModel(m pkOrm.BaseModelInterface) {
	val := reflect.ValueOf(m)
	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
	}
	objType := reflect.TypeOf(val.Interface())

	dbTable := orm.Tables.Get(objType)
	dbTable.Name = types.Q(m.TableName())
}

func (this *PgSqlAdapter) GetDb() (*pg.DB, error) {
	db := this.Db
	if db == nil {
		db = Db
	}
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return db, nil
}

func (this *PgSqlAdapter) FindOneByCond(m pkOrm.BaseModelInterface, cond string, val []interface{}, selCols []string) error {

	db, err := this.GetDb()
	if err != nil {
		return err
	}

	models := db.Model(m).Where(cond, val...)
	if len(selCols) > 0 {
		models = models.Column(selCols...)
	}

	err = models.Select(m)
	return err
}

func (this *PgSqlAdapter) FindOneBySql(m pkOrm.BaseModelInterface, sql string, val ...interface{}) error {

	db, err := this.GetDb()
	if err != nil {
		return err
	}

	_, err = db.QueryOne(m, sql, val...)
	return err
}

func (this *PgSqlAdapter) UpdateByCond(m pkOrm.BaseModelInterface, cond string, val []interface{}, selCols []string) error {
	db, err := this.GetDb()
	if err != nil {
		return err
	}

	models := db.Model(m).Where(cond, val...)

	if len(selCols) > 0 {
		models = models.Column(selCols...)
	}

	_, err = models.Update(m)

	return err
}

func (this *PgSqlAdapter) DeleteByCond(m pkOrm.BaseModelInterface, cond string, val ...interface{}) error {
	db, err := this.GetDb()
	if err != nil {
		return err
	}

	models := db.Model(m).Where(cond, val...)

	_, err = models.Delete()

	return err
}

func (this *PgSqlAdapter) Insert(m pkOrm.BaseModelInterface, selCols ...string) error {
	db, err := this.GetDb()
	if err != nil {
		return err
	}

	models := db.Model(m)
	if len(selCols) > 0 {
		models = models.Column(selCols...)
	}

	_, err = models.Insert(m)

	return err
}

func (this *PgSqlAdapter) LimitSqlStyle() string {
	return "LIMIT {limit} OFFSET {offset}"
}
func (this *PgSqlAdapter) QueryOneBySql(recordPointer interface{}, sql string, val ...interface{}) error {
	db, err := this.GetDb()
	if err != nil {
		return err
	}

	_, e := db.QueryOne(pg.Scan(recordPointer), sql, val...)
	return e
}

func (this *PgSqlAdapter) QueryBySql(recordPointer interface{}, sql string, val ...interface{}) error {
	db, err := this.GetDb()
	if err != nil {
		return err
	}

	_, e := db.Query(recordPointer, sql, val...)
	return e
}

func (this *PgSqlAdapter) InArg(arg interface{}) interface{} {
	return pg.In(arg)
}
