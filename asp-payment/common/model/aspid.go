package model

type AspId struct {
	Id         uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Type       string `gorm:"column:type;type:varchar(20);comment:资源类型,如merchant,store;NOT NULL" json:"type"`
	Key        string `gorm:"column:key;type:varchar(32);comment:签名的key;NOT NULL" json:"key"`
	Status     uint   `gorm:"column:status;type:tinyint(3) unsigned;default:0;comment:app状态;NOT NULL" json:"status"`
	CreateTime uint64 `gorm:"column:create_time;type:bigint(20) unsigned;default:0;NOT NULL" json:"create_time"`
	UpdateTime uint64 `gorm:"column:update_time;type:bigint(20) unsigned;default:0;NOT NULL" json:"update_time"`
}

func (m *AspId) TableName() string {
	return "asp_id"
}
