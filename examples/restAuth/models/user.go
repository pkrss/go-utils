package models

import (
	"github.com/pkrss/go-utils/orm"
	custom "github.com/pkrss/go-utils/orm/fields"
)

////////////////////////////////////////////////////////////////////////////

type User struct {
	orm.BaseModel

	Id         custom.UUID     `json:"id"`
	Mobile     string          `json:"mobile"`
	NickName   string          `json:"nickName"`
	CreateTime custom.JsonTime `json:"createTime"`
	Password   string          `json:"password"`
	Role       int             `json:"role"`
}

func (this *User) TableName() string {
	return "myzc_user"
}

func (this *User) FilterValue() {
	this.Password = "******"
}
