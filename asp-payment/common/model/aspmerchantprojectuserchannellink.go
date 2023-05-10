package model

// AspMerchantProjectUserChannelLink 外部商家项目-用户信息-支付渠道关联表
type AspMerchantProjectUserChannelLink struct {
	Id         int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	MchId      int    `gorm:"column:mch_id;type:int(11);comment:CP商户id;NOT NULL" json:"mch_id"`
	ProjectId  int    `gorm:"column:project_id;type:int(11);default:1;comment:外部商家项目id;NOT NULL" json:"project_id"`
	ChannelId  int    `gorm:"column:channel_id;type:int(11);default:0;comment:支付渠道ID;NOT NULL" json:"channel_id"`
	Uid        string `gorm:"column:uid;type:varchar(255);comment:外部商家项目玩家id;NOT NULL" json:"uid"`
	SuccessNum int    `gorm:"column:success_num;type:int(11);default:0;comment:成功次数;NOT NULL" json:"success_num"`
	FailedNum  int    `gorm:"column:failed_num;type:int(11);default:0;comment:失败次数;NOT NULL" json:"failed_num"`
	VipLevel   int    `gorm:"column:vip_level;type:tinyint(4);default:0;comment:外部商家项目玩家 模型等级;NOT NULL" json:"vip_level"`
	Status     int    `gorm:"column:status;type:tinyint(255);default:0;comment:外部商家状态;NOT NULL" json:"status"`
	Remark     string `gorm:"column:remark;type:varchar(1024);comment:备注;NOT NULL" json:"remark"`
	PayTime    int64  `gorm:"column:pay_time;type:bigint(20);default:0;comment:支付的时间;NOT NULL" json:"pay_time"`
	CreateTime uint64 `gorm:"column:create_time;type:bigint(20) unsigned;default:0;comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime uint64 `gorm:"column:update_time;type:bigint(20) unsigned;default:0;comment:修改时间;NOT NULL" json:"update_time"`
}

func (m *AspMerchantProjectUserChannelLink) TableName() string {
	return "asp_merchant_project_user_channel_link"
}
