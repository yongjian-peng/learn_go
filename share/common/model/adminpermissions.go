package model

import (
	"time"
)

// AdminPermissions Api权限表
type AdminPermissions struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:id" json:"id"`
	Pid        int       `gorm:"column:pid;type:int(11);comment:父级id;NOT NULL" json:"pid"`
	Name       string    `gorm:"column:name;type:varchar(255);comment:权限名称;NOT NULL" json:"name"`
	RouteUrl   string    `gorm:"column:route_url;type:varchar(255);comment:接口路由url;NOT NULL" json:"route_url"`
	Sort       int       `gorm:"column:sort;type:int(10);default:0;comment:排序值;NOT NULL" json:"sort"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;comment:创建时间" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;comment:更新时间" json:"update_time"`
}

func (m *AdminPermissions) TableName() string {
	return "admin_permissions"
}
