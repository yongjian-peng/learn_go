package model

// AspMerchantProject 外部商家项目表（游戏表）
type AspMerchantProject struct {
	Id              int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	MchId           int    `gorm:"column:mch_id;type:int(11);default:0;comment:关联外部商家id;NOT NULL" json:"mch_id"`
	ProjectName     string `gorm:"column:project_name;type:varchar(255);comment:cp游戏名称;NOT NULL" json:"project_name"`
	ProjectType     string `gorm:"column:project_type;type:varchar(255);comment:cp游戏类型 slots tp;NOT NULL" json:"project_type"`
	ProjectCodeType string `gorm:"column:project_code_type;type:varchar(255);comment:cp游戏编码类型 native:原生 h5:h5;NOT NULL" json:"project_code_type"`
	Cooperation     string `gorm:"column:cooperation;type:varchar(255);comment:合作模式 united:联运 alone:独代;NOT NULL" json:"cooperation"`
	Status          int    `gorm:"column:status;type:tinyint(255);default:0;comment:cp游戏状态 0 1;NOT NULL" json:"status"`
	Remark          string `gorm:"column:remark;type:varchar(1024);comment:备注;NOT NULL" json:"remark"`
	CreateTime      uint64 `gorm:"column:create_time;type:bigint(20) unsigned;default:0;comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime      uint64 `gorm:"column:update_time;type:bigint(20) unsigned;default:0;comment:修改时间;NOT NULL" json:"update_time"`
}

func (m *AspMerchantProject) TableName() string {
	return "asp_merchant_project"
}
