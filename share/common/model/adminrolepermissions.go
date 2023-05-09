package model

import (
	"time"
)

// AdminRolePermissions 角色权限表
type AdminRolePermissions struct {
	Id            int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:id" json:"id"`
	RoleId        int       `gorm:"column:role_id;type:int(11);comment:角色id;NOT NULL" json:"role_id"`
	PermissionsId int       `gorm:"column:permissions_id;type:int(11);comment:权限id;NOT NULL" json:"permissions_id"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;comment:更新时间;NOT NULL" json:"update_time"`
}

func (m *AdminRolePermissions) TableName() string {
	return "admin_role_permissions"
}
