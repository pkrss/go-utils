package pqsql

type BaseDao struct {
}

func (this *BaseDao) FetchById(id interface{}, ouputOb interface{}) error {
	sql := `SELECT * FROM ` + getTableName(ouputOb) + ` WHERE id = $1 LIMIT 1`
	return Db.Query(sql, id).Scan(ouputOb)
}

func (this *BaseDao) UpdateById(id interface{}) error {
	sql := `SELECT * FROM ` + getTableName(ouputOb) + ` WHERE id = $1 LIMIT 1`
	return Db.Query(sql, id).Scan(this)
}
