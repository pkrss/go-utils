package pqsql

type BaseModel struct {
}

func (this *BaseModel) TableName() string {
	return ""
}

func (this *BaseModel) FetchById(id interface{}) error {
	sql := `SELECT * FROM ` + TableName() + ` WHERE id = $1 LIMIT 1`
	return Db.Query(sql, id).Scan(this)
}

func (this *BaseModel) UpdateById(id interface{}) error {
	sql := `SELECT * FROM ` + TableName() + ` WHERE id = $1 LIMIT 1`
	return Db.Query(sql, id).Scan(this)
}


func (this *BaseModel) Insert() id, error {
	sql := `INSERT INTO ` + TableName() + `(name, favorite_fruit, age) VALUES('beatrice', 'starfruit', 93) RETURNING id`

	var id int
	err := db.QueryRow(sql).Scan(&id)
	return id, err
}
