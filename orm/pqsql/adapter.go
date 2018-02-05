package pqsql

import (
	"errors"
	"reflect"

	"github.com/go-pg/pg"
	// "github.com/go-pg/pg/orm"
	// "github.com/go-pg/pg/types"
	pkOrm "github.com/pkrss/go-utils/orm"
)

type PgSqlAdapter struct {
	Db *pg.DB
}

func (this *PgSqlAdapter) RegModel(m pkOrm.BaseModelInterface) {
	// val := reflect.ValueOf(m)
	// switch val.Kind() {
	// case reflect.Ptr:
	// 	val = val.Elem()
	// }
	// objType := reflect.TypeOf(val.Interface())

	// dbTable := orm.Tables.Get(objType)
	// dbTable.Name = types.Q(m.TableName())
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

func (this *PgSqlAdapter) ExecSql(sql string, val ...interface{}) error {
	db, err := this.GetDb()
	if err != nil {
		return err
	}

	_, e := db.Exec(sql, val...)
	return e
}
func (this *PgSqlAdapter) QueryOneBySql(outputRecord interface{}, sql string, val ...interface{}) error {
	if outputRecord == nil {
		return this.ExecSql(sql, val...)
	}

	db, err := this.GetDb()
	if err != nil {
		return err
	}

	vv := reflect.ValueOf(outputRecord)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	if vv.Kind() != reflect.Struct {
		outputRecord = pg.Scan(outputRecord)
	}

	_, e := db.QueryOne(outputRecord, sql, val...)
	return e
}

func (this *PgSqlAdapter) QueryBySql(outputRecords interface{}, sql string, val ...interface{}) error {
	if outputRecords == nil {
		return this.ExecSql(sql, val...)
	}

	db, err := this.GetDb()
	if err != nil {
		return err
	}

	_, e := db.Query(outputRecords, sql, val...)
	return e
}

func (this *PgSqlAdapter) SqlInArg(arg interface{}) interface{} {
	return pg.In(arg)
}

func (this *PgSqlAdapter) SqlReturnSql() string {
	return " RETURNING {id}"
}

func (this *PgSqlAdapter) SqlLimitStyle() string {
	return "LIMIT {limit} OFFSET {offset}"
}
