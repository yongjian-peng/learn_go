package model

import (
	"time"
)

// AdminUser 管理员表
type AdminUser struct {
	Id         int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:id" json:"id"`
	Username   string    `gorm:"column:username;type:varchar(50);comment:用户名;NOT NULL" json:"username"`
	Password   string    `gorm:"column:password;type:char(32);comment:密码;NOT NULL" json:"password"`
	Salt       string    `gorm:"column:salt;type:varchar(30);comment:密码加盐key;NOT NULL" json:"salt"`
	Avatar     string    `gorm:"column:avatar;type:varchar(255);comment:用户头像;NOT NULL" json:"avatar"`
	Status     int       `gorm:"column:status;type:tinyint(4);default:1;comment:用户状态 0:禁用 1:启用;NOT NULL" json:"status"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;comment:创建时间" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;comment:更新时间" json:"update_time"`
}

func (m *AdminUser) TableName() string {
	return "admin_user"
}
