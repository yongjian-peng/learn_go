package model

import (
	"time"
)

type AdminUserRole struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	AdminId    int       `gorm:"column:admin_id;type:int(11)" json:"admin_id"`
	RoleId     int       `gorm:"column:role_id;type:int(11)" json:"role_id"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

func (m *AdminUserRole) TableName() string {
	return "admin_user_role"
}
