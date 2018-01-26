package pqsql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/pkrss/go-utils/reflect"
	pkStrings "github.com/pkrss/go-utils/strings"
)

type BaseModelInterface interface {
	TableName() string
	FetchById(id interface{}) error
	UpdateById(id interface{}) error
	Insert() (interface{}, error)
}

type BaseModel struct {
	Db *sql.DB
}

func (this *BaseModel) TableName() string {
	return ""
}

func (this *BaseModel) FetchById(id interface{}) error {
	sql := `SELECT * FROM ` + this.TableName() + ` WHERE id = $1`
	return Db.QueryRow(sql, id).Scan(this)
}

func (this *BaseModel) UpdateById(id interface{}, cols ...string) error {
	n2v := reflect.GetStructFieldName2ValueMap(this)
	if len(n2v) == 0 {
		return errors.New("GetStructFieldName2ValueMap return nil")
	}

	ns := make([]string, 0)
	vs := make([]interface{}, 0)
	ts := make([]string, len(n2v), '?')

	sql := `UPDATE ` + this.TableName() + ` SET`

	for k, v := range n2v {
		n := pkStrings.StringToCamelCase(k)
		sql += n + `=?,`
		vs = append(vs, v)
	}

	sql = sql[0:-1]

	Db.QueryRow(sql, vs)

	return nil
}

func (this *BaseModel) Insert(cols ...string) (interface{}, error) {

	n2v := reflect.GetStructFieldName2ValueMap(this)
	if n2v == nil {
		return nil, errors.New("GetStructFieldName2ValueMap return nil")
	}

	ns := make([]string, 0)
	vs := make([]interface{}, 0)
	ts := make([]string, len(n2v), '?')
	for k, v := range n2v {
		n := pkStrings.StringToCamelCase(k)
		ns = append(ns, n)
		vs = append(vs, v)
	}

	sql := `INSERT INTO ` + this.TableName() + `(` + strings.Join(ns, ',') + `) VALUES(` + strings.Join(ts, ',') + `) RETURNING id`

	var id int
	err := Db.QueryRow(sql, vs).Scan(&id)
	return id, err
}
