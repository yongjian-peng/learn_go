package model

import (
	"time"
)

// App 应用信息
type App struct {
	Id          int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:应用id" json:"id"`
	Secret      string    `gorm:"column:secret;type:char(32);comment:应用秘钥" json:"secret"`
	Name        string    `gorm:"column:name;type:varchar(50);comment:应用名称" json:"name"`
	PackageName string    `gorm:"column:package_name;type:varchar(255);comment:应用报名" json:"package_name"`
	Salt        string    `gorm:"column:salt;type:char(6);comment:应用秘钥盐" json:"salt"`
	Status      int       `gorm:"column:status;type:tinyint(4);comment:状态 1:启用 0:禁用" json:"status"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;comment:更新时间" json:"update_time"`
}

func (m *App) TableName() string {
	return "app"
}
