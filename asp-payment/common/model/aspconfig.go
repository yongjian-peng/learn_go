package model

// 系统配置表
type AspConfig struct {
	Id         uint   `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Name       string `gorm:"column:name;type:varchar(50);comment:config name;NOT NULL" json:"name"`
	Data       string `gorm:"column:data;type:text;comment:config data;NOT NULL" json:"data"`
	Status     int    `gorm:"column:status;type:tinyint(4);default:1;comment:config status 1: normal 0: disable;NOT NULL" json:"status"`
	Remark     string `gorm:"column:remark;type:varchar(255);comment:备注;NOT NULL" json:"remark"`
	CreateTime int64  `gorm:"column:create_time;type:bigint(20)" json:"create_time"`
	UpdateTime int64  `gorm:"column:update_time;type:bigint(20)" json:"update_time"`
}

func (m *AspConfig) TableName() string {
	return "asp_config"
}
