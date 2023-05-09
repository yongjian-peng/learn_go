package model

import (
	"time"
)

// User 用户表
type User struct {
	Id            int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:用户id" json:"id"`
	Appid         int       `gorm:"column:appid;type:int(11);comment:应用id;NOT NULL" json:"appid"`
	AppUid        int64     `gorm:"column:app_uid;type:bigint(20);comment:第三方应用的用户uid;NOT NULL" json:"app_uid"`
	Phone         string    `gorm:"column:phone;type:varchar(30);comment:手机号;NOT NULL" json:"phone"`
	InviteCode    string    `gorm:"column:invite_code;type:varchar(30);comment:用户的邀请码;NOT NULL" json:"invite_code"`
	InviteUid     int64     `gorm:"column:invite_uid;type:bigint(20);default:0;comment:邀请人uid(上级uid);NOT NULL" json:"invite_uid"`
	DeviceId      string    `gorm:"column:device_id;type:varchar(100);comment:设备号;NOT NULL" json:"device_id"`
	Status        int       `gorm:"column:status;type:tinyint(4);default:1;comment:用户状态 1:启用 0:禁用;NOT NULL" json:"status"`
	OnlineStatus  int       `gorm:"column:online_status;type:tinyint(4);default:0;comment:在线状态 1:在线 0:离线;NOT NULL" json:"online_status"`
	RegisterType  string    `gorm:"column:register_type;type:varchar(30);comment:注册方式;NOT NULL" json:"register_type"`
	RegisterTime  time.Time `gorm:"column:register_time;type:datetime;comment:注册时间" json:"register_time"`
	Level         int       `gorm:"column:level;type:int(11);default:1;comment:用户等级" json:"level"`
	LastLoginTime time.Time `gorm:"column:last_login_time;type:datetime;comment:最后登录时间" json:"last_login_time"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;comment:创建时间" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;comment:更新时间" json:"update_time"`
}

func (m *User) TableName() string {
	return "user"
}
