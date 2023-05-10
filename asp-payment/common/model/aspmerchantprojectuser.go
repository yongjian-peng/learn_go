package model

// AspMerchantProjectUser 外部商家项目-用户信息表
type AspMerchantProjectUser struct {
	Id           int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	MchId        int    `gorm:"column:mch_id;type:int(11);default:0;comment:外部商家id;NOT NULL" json:"mch_id"`
	MchProjectId int    `gorm:"column:mch_project_id;type:int(11);default:1;comment:外部商家项目id;NOT NULL" json:"mch_project_id"`
	Uid          string `gorm:"column:uid;type:varchar(255);comment:外部商家项目玩家id;NOT NULL" json:"uid"`
	Phone        string `gorm:"column:phone;type:varchar(32);comment:外部商家项目玩家手机号;NOT NULL" json:"phone"`
	VipLevel     int    `gorm:"column:vip_level;type:tinyint(4);default:0;comment:外部商家项目玩家 模型等级;NOT NULL" json:"vip_level"`
	Status       int    `gorm:"column:status;type:tinyint(255);default:0;comment:外部商家状态;NOT NULL" json:"status"`
	Email        string `gorm:"column:email;type:varchar(255);comment:外部商家邮箱;NOT NULL" json:"email"`
	Remark       string `gorm:"column:remark;type:varchar(1024);comment:备注;NOT NULL" json:"remark"`
	CreateTime   uint64 `gorm:"column:create_time;type:bigint(20) unsigned;default:0;comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime   uint64 `gorm:"column:update_time;type:bigint(20) unsigned;default:0;comment:修改时间;NOT NULL" json:"update_time"`
}

func (m *AspMerchantProjectUser) TableName() string {
	return "asp_merchant_project_user"
}
