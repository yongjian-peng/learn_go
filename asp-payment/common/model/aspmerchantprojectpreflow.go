package model

// AspMerchantProjectPreFlow 资金锁定明细表
type AspMerchantProjectPreFlow struct {
	Id                   int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	MchId                int    `gorm:"column:mch_id;type:int(11);comment:cp ID;NOT NULL" json:"mch_id"`
	MchProjectId         int    `gorm:"column:mch_project_id;type:int(11);comment:cp 项目ID;NOT NULL" json:"mch_project_id"`
	MchProjectCurrencyId int    `gorm:"column:mch_project_currency_id;type:int(11);comment:账户币种ID;NOT NULL" json:"mch_project_currency_id"`
	Currency             string `gorm:"column:currency;type:varchar(45);comment:币种;NOT NULL" json:"currency"`
	PreTotalFee          int64  `gorm:"column:pre_total_fee;type:bigint(20);comment:冻结金额(负);NOT NULL" json:"pre_total_fee"`
	BusinessType         int    `gorm:"column:business_type;type:tinyint(4);comment:业务类型（1代收，2代付）;NOT NULL" json:"business_type"`
	BusinessSourceId     int64  `gorm:"column:business_source_id;type:bigint(20);comment:产品对应的业务关联ID;NOT NULL" json:"business_source_id"`
	Status               int    `gorm:"column:status;type:tinyint(4);comment:锁定的状态(1.锁定中2.已解锁);NOT NULL" json:"status"`
	Remark               string `gorm:"column:remark;type:varchar(255);comment:备注;NOT NULL" json:"remark"`
	CreateTime           int64  `gorm:"column:create_time;type:bigint(20);comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime           int64  `gorm:"column:update_time;type:bigint(20);comment:更新时间;NOT NULL" json:"update_time"`
}

func (m *AspMerchantProjectPreFlow) TableName() string {
	return "asp_merchant_project_pre_flow"
}
