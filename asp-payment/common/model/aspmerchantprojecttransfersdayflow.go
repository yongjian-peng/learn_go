package model

// AspMerchantProjectTransfersDayFlow 资金待结算余额转可用余额日结明细表
type AspMerchantProjectTransfersDayFlow struct {
	Id                   int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	MchId                int    `gorm:"column:mch_id;type:int(11);comment:cp ID;NOT NULL" json:"mch_id"`
	MchProjectId         int    `gorm:"column:mch_project_id;type:int(11);comment:cp 项目ID;NOT NULL" json:"mch_project_id"`
	MchProjectCurrencyId int    `gorm:"column:mch_project_currency_id;type:int(11);comment:账户币种ID;NOT NULL" json:"mch_project_currency_id"`
	Currency             string `gorm:"column:currency;type:varchar(45);comment:币种;NOT NULL" json:"currency"`
	TotalFee             int64  `gorm:"column:total_fee;type:bigint(20);comment:待结算余额;NOT NULL" json:"total_fee"`
	Day                  int64  `gorm:"column:day;type:bigint(20);comment:结算日期;NOT NULL" json:"day"`
	Remark               string `gorm:"column:remark;type:varchar(255);comment:备注;NOT NULL" json:"remark"`
	CreateTime           int64  `gorm:"column:create_time;type:bigint(20);comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime           int64  `gorm:"column:update_time;type:bigint(20);comment:更新时间;NOT NULL" json:"update_time"`
}

func (m *AspMerchantProjectTransfersDayFlow) TableName() string {
	return "asp_merchant_project_transfers_day_flow"
}
