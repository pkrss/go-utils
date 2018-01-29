package pqsql

type BaseModelInterface interface {
	TableName() string
	IdColumn() string
	FilterValue()
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
