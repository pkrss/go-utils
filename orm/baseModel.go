package orm

type BaseModelInterface interface {
	TableName() string
	IdColumn() string
	FilterValue()
	SelSql() string
}

type BaseModel struct {
}

func (this *BaseModel) TableName() string {
	return ""
}

func (this *BaseModel) IdColumn() string {
	return "id"
}

func (this *BaseModel) FilterValue() {
}

func (this *BaseModel) SelSql() string {
	return ""
}
