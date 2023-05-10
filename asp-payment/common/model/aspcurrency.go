package model

// AspCurrency 系统支持的币种
type AspCurrency struct {
	Id           int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Currency     string `gorm:"column:currency;type:varchar(45);comment:币种识别码;NOT NULL" json:"currency"`
	CurrencyName string `gorm:"column:currency_name;type:varchar(64);comment:币种名称-中文;NOT NULL" json:"currency_name"`
	Img          string `gorm:"column:img;type:varchar(255);comment:币种图片地址;NOT NULL" json:"img"`
	Status       int    `gorm:"column:status;type:tinyint(4);default:1;comment:状态 1.正常;NOT NULL" json:"status"`
	CreateTime   int64  `gorm:"column:create_time;type:bigint(20);NOT NULL" json:"create_time"`
	UpdateTime   int64  `gorm:"column:update_time;type:bigint(20);NOT NULL" json:"update_time"`
}

func (m *AspCurrency) TableName() string {
	return "asp_currency"
}
