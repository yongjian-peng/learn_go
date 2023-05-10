package model

// AspIpFilter 商户白名单的表
type AspIpFilter struct {
	Id         int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	Type       int    `gorm:"column:type;type:tinyint(4);default:1;comment:1:CP项目白名单;NOT NULL" json:"type"`
	TypeId     int    `gorm:"column:type_id;type:int(11);default:0;comment:类型对应的表主键的id;NOT NULL" json:"type_id"`
	IpAddr     string `gorm:"column:ip_addr;type:varchar(45);comment:IP地址，暂仅支持IPv4格式;NOT NULL" json:"ip_addr"`
	CreateTime uint64 `gorm:"column:create_time;type:bigint(20) unsigned;default:0;comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime uint64 `gorm:"column:update_time;type:bigint(20) unsigned;default:0;comment:更新时间;NOT NULL" json:"update_time"`
}

func (m *AspIpFilter) TableName() string {
	return "asp_ip_filter"
}
