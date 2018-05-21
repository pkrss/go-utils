package models

import (
	"time"

	"github.com/pkrss/go-utils/orm"
)

////////////////////////////////////////////////////////////////////////////

type UserContext struct {
	orm.BaseModel
	ID         int64     `json:"id"`
	UserId     int64     `json:"userId"`
	UserName   string    `json:"userName"`
	Ip         string    `json:"ip"`
	Role       int       `json:"role"`
	CreateTime time.Time `json:"createTime,int64" orm:"auto_now_add"`
	UpdateTime time.Time `json:"updateTime,int64" orm:"auto_now"`
}

func (this *UserContext) TableName() string {
	return "pkrss_user_context"
}
