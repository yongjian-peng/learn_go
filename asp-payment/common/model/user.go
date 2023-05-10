package model

import (
	"database/sql"
)

type User struct {
	Id           uint64       `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	OpenId       string       `gorm:"column:open_id;type:varchar(255);NOT NULL" json:"open_id"`
	UnionId      string       `gorm:"column:union_id;type:varchar(255);NOT NULL" json:"union_id"`
	AvatarUrl    string       `gorm:"column:avatarUrl;type:varchar(255);comment:头像地址;NOT NULL" json:"avatarUrl"`
	NickName     string       `gorm:"column:nick_name;type:varchar(255);comment:昵称;NOT NULL" json:"nick_name"`
	Gender       int          `gorm:"column:gender;type:tinyint(4);comment:性别;NOT NULL" json:"gender"`
	RegisterType int          `gorm:"column:register_type;type:tinyint(4);comment:注册方式;NOT NULL" json:"register_type"`
	CreatedAt    sql.NullTime `gorm:"column:created_at;type:timestamp" json:"created_at"`
	UpdatedAt    sql.NullTime `gorm:"column:updated_at;type:timestamp" json:"updated_at"`
}

func (m *User) TableName() string {
	return "user"
}
