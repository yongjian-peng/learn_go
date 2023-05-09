package model

import (
	"time"
)

type AdminRole struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;comment:id" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(255);comment:角色名称;NOT NULL" json:"name"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

func (m *AdminRole) TableName() string {
	return "admin_role"
}
