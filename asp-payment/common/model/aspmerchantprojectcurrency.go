package model

// AspMerchantProjectCurrency 账户币种关联表
type AspMerchantProjectCurrency struct {
	Id                   int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	MchId                int    `gorm:"column:mch_id;type:int(11);comment:商户ID;NOT NULL" json:"mch_id"`
	MchProjectId         uint   `gorm:"column:mch_project_id;type:int(10) unsigned;default:0;comment:商户项目ID;NOT NULL" json:"mch_project_id"`
	CurrencyId           uint   `gorm:"column:currency_id;type:int(10) unsigned;comment:币种ID;NOT NULL" json:"currency_id"`
	Currency             string `gorm:"column:currency;type:varchar(255);comment:币种名称;NOT NULL" json:"currency"`
	TotalFee             int64  `gorm:"column:total_fee;type:bigint(20);comment:待结算账户余额（不包含可用余额）;NOT NULL" json:"total_fee"`
	FreezeFee            int64  `gorm:"column:freeze_fee;type:bigint(20);comment:冻结金额（包含订单投诉冻结和预扣款冻结）;NOT NULL" json:"freeze_fee"`
	SettlementInProgress int64  `gorm:"column:settlement_in_progress;type:bigint(20);default:0;comment:提现中余额;NOT NULL" json:"settlement_in_progress"`
	AvailableTotalFee    int64  `gorm:"column:available_total_fee;type:bigint(255);default:0;comment:可用余额;NOT NULL" json:"available_total_fee"`
	Status               uint   `gorm:"column:status;type:tinyint(3) unsigned;default:0;comment:状态（1正常，2关闭）;NOT NULL" json:"status"`
	CreateTime           int64  `gorm:"column:create_time;type:bigint(20);comment:添加时间;NOT NULL" json:"create_time"`
	UpdateTime           uint64 `gorm:"column:update_time;type:bigint(20) unsigned;default:0;comment:修改时间;NOT NULL" json:"update_time"`
}

func (m *AspMerchantProjectCurrency) TableName() string {
	return "asp_merchant_project_currency"
}
