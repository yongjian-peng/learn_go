package model

// 资金流明细表-所有的资金出入都需要记录在这个表
type AspMerchantProjectCapitalFlow struct {
	Id                  int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	MchId               int    `gorm:"column:mch_id;type:int(11);comment:cp ID;NOT NULL" json:"mch_id"`
	MchProjectId        int    `gorm:"column:mch_project_id;type:int(11);comment:cp项目ID;NOT NULL" json:"mch_project_id"`
	AccountCurrencyId   int    `gorm:"column:account_currency_id;type:int(11);comment:账户币种ID;NOT NULL" json:"account_currency_id"`
	Currency            string `gorm:"column:currency;type:varchar(45);comment:申请出账的币种;NOT NULL" json:"currency"`
	TotalFee            int64  `gorm:"column:total_fee;type:bigint(20);comment:待结算金额(入VA为正，出VA为负);NOT NULL" json:"total_fee"`
	TotalFeeBefore      int64  `gorm:"column:total_fee_before;type:bigint(255);default:0;comment:当前资金变动前待结算余额;NOT NULL" json:"total_fee_before"`
	TotalFeeSurplus     int64  `gorm:"column:total_fee_surplus;type:bigint(20);comment:当前资金变动后待结算余额(转账成功后填入);NOT NULL" json:"total_fee_surplus"`
	FreezeFee           int64  `gorm:"column:freeze_fee;type:bigint(20);default:0;comment:冻结余额;NOT NULL" json:"freeze_fee"`
	FreezeFeeBefore     int64  `gorm:"column:freeze_fee_before;type:bigint(20);default:0;comment:当前资金变动前冻结余额;NOT NULL" json:"freeze_fee_before"`
	FreezeFeeSurplus    int64  `gorm:"column:freeze_fee_surplus;type:bigint(20);default:0;comment:当前资金变动后冻结余额;NOT NULL" json:"freeze_fee_surplus"`
	AvailableFee        int64  `gorm:"column:available_fee;type:bigint(20);default:0;comment:可用余额;NOT NULL" json:"available_fee"`
	AvailableFeeBefore  int64  `gorm:"column:available_fee_before;type:bigint(20);default:0;comment:当前资金变动前可用余额;NOT NULL" json:"available_fee_before"`
	AvailableFeeSurplus int64  `gorm:"column:available_fee_surplus;type:bigint(20);default:0;comment:当前资金变动后可用余额;NOT NULL" json:"available_fee_surplus"`
	BusinessType        int    `gorm:"column:business_type;type:tinyint(4);comment:业务类型 1:代收 2:代付冻结 3:充值 4:月结算 5:代付解冻 6:代收日结;NOT NULL" json:"business_type"`
	BusinessSourceId    int64  `gorm:"column:business_source_id;type:bigint(20);comment:具体业务类型来源表关联id (类型关联表 1: asp_order 2: asp_payout 4: asp_merchant_project_month_settlement 5:asp_payout 6:asp_merchant_project_transfers_day_flow );NOT NULL" json:"business_source_id"`
	Remark              string `gorm:"column:remark;type:varchar(50);comment:备注;NOT NULL" json:"remark"`
	CreateTime          int64  `gorm:"column:create_time;type:bigint(20);comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime          int64  `gorm:"column:update_time;type:bigint(20);comment:更新时间;NOT NULL" json:"update_time"`
}

func (m *AspMerchantProjectCapitalFlow) TableName() string {
	return "asp_merchant_project_capital_flow"
}
