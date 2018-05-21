package models

import (
	"time"

	"github.com/pkrss/go-utils/orm"
)

////////////////////////////////////////////////////////////////////////////

type User struct {
	orm.BaseModel

	ID                  int64     `json:"id"`
	UserName            string    `json:"userName"`
	Password            string    `json:"password"`
	Role                int       `json:"role"`
	Denied              bool      `json:"denied"`
	Email               string    `json:"email"`
	LastErrorLoginTime  time.Time `json:"last_error_login_time,int64"`
	LastErrorLoginTimer int       `json:"last_error_login_timer"`
	CreateTime          time.Time `json:"createTime,int64" orm:"auto_now_add"`
	UpdateTime          time.Time `json:"updateTime,int64" orm:"auto_now"`
}

func (this *User) TableName() string {
	return "pkrss_user"
}

func (this *User) FilterValue() {
	this.Password = "******"
}
