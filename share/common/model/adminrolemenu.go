package model

import (
	"time"
)

// AdminRoleMenu 角色菜单表
type AdminRoleMenu struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:id" json:"id"`
	RoleId     int       `gorm:"column:role_id;type:int(11);comment:角色id" json:"role_id"`
	MenuId     int       `gorm:"column:menu_id;type:int(11);comment:角色菜单id" json:"menu_id"`
	Btn        string    `gorm:"column:btn;type:text;comment:角色菜单按钮json数组"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;comment:创建时间" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;comment:更新时间" json:"update_time"`
}

func (m *AdminRoleMenu) TableName() string {
	return "admin_role_menu"
}
