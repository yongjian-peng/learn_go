package model

// AspMerchant 外部商家表包含（cp 联运）
type AspMerchant struct {
	Id           int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	ParentId     int    `gorm:"column:parent_id;type:int(11);default:0;comment:外部商家&商家的父级id;NOT NULL" json:"parent_id"`
	MerchantType int    `gorm:"column:merchant_type;type:tinyint(1);default:1;comment:外部商家类型 1 商家代理 2商家;NOT NULL" json:"merchant_type"`
	Title        string `gorm:"column:title;type:varchar(255);comment:外部商家账户标题;NOT NULL" json:"title"`
	Name         string `gorm:"column:name;type:varchar(255);comment:外部商家名称;NOT NULL" json:"name"`
	Sort         int    `gorm:"column:sort;type:tinyint(255);default:0;comment:外部商家账户排序值 权重越大优先级越高;NOT NULL" json:"sort"`
	Status       int    `gorm:"column:status;type:tinyint(255);default:0;comment:外部商家状态;NOT NULL" json:"status"`
	Email        string `gorm:"column:email;type:varchar(255);comment:外部商家邮箱;NOT NULL" json:"email"`
	Remark       string `gorm:"column:remark;type:varchar(1024);comment:备注;NOT NULL" json:"remark"`
	CreateTime   uint64 `gorm:"column:create_time;type:bigint(20) unsigned;default:0;comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime   uint64 `gorm:"column:update_time;type:bigint(20) unsigned;default:0;comment:修改时间;NOT NULL" json:"update_time"`
}

func (m *AspMerchant) TableName() string {
	return "asp_merchant"
}
