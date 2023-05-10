package model

// AspMerchantProjectChannelDepartTradeTypeLink
// cp项目支付渠道内部商户支付方式关联表-一个cp项目可以使用多个内部商户去支付
type AspMerchantProjectChannelDepartTradeTypeLink struct {
	Id            uint   `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	DepartId      int    `gorm:"column:depart_id;type:int(11);default:0;comment:内部商户ID;NOT NULL" json:"depart_id"`
	MchId         int    `gorm:"column:mch_id;type:int(11);default:0;comment:cp的id;NOT NULL" json:"mch_id"`
	MchProjectId  int    `gorm:"column:mch_project_id;type:int(11);default:0;comment:cp项目id;NOT NULL" json:"mch_project_id"`
	Status        int    `gorm:"column:status;type:tinyint(2);default:0;comment:状态 0 否 1 是;NOT NULL" json:"status"`
	ChannelId     int    `gorm:"column:channel_id;type:int(11);comment:渠道id" json:"channel_id"`
	ChannelStatus int    `gorm:"column:channel_status;type:tinyint(4);comment:渠道状态" json:"channel_status"`
	TradeStatus   int    `gorm:"column:trade_status;type:tinyint(4);comment:支付配置状态" json:"trade_status"`
	TradeType     string `gorm:"column:trade_type;type:varchar(30);comment:支付类型" json:"trade_type"`
	Sort          int    `gorm:"column:sort;type:int(10);default:0;comment:排序 权重值越大 优先级越高;NOT NULL" json:"sort"`
	CreateTime    uint64 `gorm:"column:create_time;type:bigint(20) unsigned;default:0;comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime    uint64 `gorm:"column:update_time;type:bigint(20) unsigned;default:0;comment:修改时间;NOT NULL" json:"update_time"`
}

func (m *AspMerchantProjectChannelDepartTradeTypeLink) TableName() string {
	return "asp_merchant_project_channel_depart_trade_type_link"
}
