package model

// 银行编码类目表
type AspBankCategory struct {
	CategoryId   uint   `gorm:"column:category_id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"category_id"`
	Status       int    `gorm:"column:status;type:tinyint(4);default:1;NOT NULL" json:"status"`
	Code         string `gorm:"column:code;type:varchar(64);default:0;comment:对应编码;NOT NULL" json:"code"`
	CategoryName string `gorm:"column:category_name;type:varchar(255);comment:对应的名称;NOT NULL" json:"category_name"`
	CreateTime   uint64 `gorm:"column:create_time;type:bigint(20) unsigned;default:0;NOT NULL" json:"create_time"`
	UpdateTime   uint64 `gorm:"column:update_time;type:bigint(20) unsigned;default:0;NOT NULL" json:"update_time"`
}

func (m *AspBankCategory) TableName() string {
	return "asp_bank_category"
}
