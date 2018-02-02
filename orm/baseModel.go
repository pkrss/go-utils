package orm

import "reflect"

type BaseModelInterface interface {
	TableName() string
	IdColumn() string
	IdType() reflect.Type
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

func (this *BaseModel) IdType() reflect.Type {
	return nil
}

func (this *BaseModel) FilterValue() {
}

func (this *BaseModel) SelSql() string {
	return ""
}
