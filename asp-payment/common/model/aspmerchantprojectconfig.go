package model

// 外部商家项目配置表（游戏表）
type AspMerchantProjectConfig struct {
	Id                 int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	MchId              int    `gorm:"column:mch_id;type:int(11);default:0;comment:关联外部商家id;NOT NULL" json:"mch_id"`
	MchProjectId       string `gorm:"column:mch_project_id;type:varchar(255);comment:cp游戏名称;NOT NULL" json:"mch_project_id"`
	InFeeRate          string `gorm:"column:in_fee_rate;type:varchar(10);comment:费率：收款费率，保存两位小数 范围 0.00 ~ 0.99 之间;NOT NULL" json:"in_fee_rate"`
	OutFeeRate         string `gorm:"column:out_fee_rate;type:varchar(10);comment:费率：提现支出费率 保存两位小数 范围 0.00 ~ 0.99 之间 ;NOT NULL" json:"out_fee_rate"`
	InUpperLimit       int    `gorm:"column:in_upper_limit;type:int(11);default:0;comment:每笔收款上限金额 单位分 默认 100万卢币;NOT NULL" json:"in_upper_limit"`
	InLowerLimit       int    `gorm:"column:in_lower_limit;type:int(11);default:0;comment:每笔收款下限金额 单位分 默认 1百卢币;NOT NULL" json:"in_lower_limit"`
	OutUpperLimit      int    `gorm:"column:out_upper_limit;type:int(11);default:0;comment:每笔提现上限金额 单位分 默认 1万卢币;NOT NULL" json:"out_upper_limit"`
	OutAuditUpperLimit int    `gorm:"column:out_audit_upper_limit;type:int(11);default:0;comment:每笔收款审核上线金额 单位分 默认 5千卢币;NOT NULL" json:"out_audit_upper_limit"`
	OutLowerLimit      int    `gorm:"column:out_lower_limit;type:int(11);default:0;comment:每笔提现下限金额 单位分 默认 1百卢币;NOT NULL" json:"out_lower_limit"`
	OutDayUpperLimit   int    `gorm:"column:out_day_upper_limit;type:int(11);default:0;comment:每天代收上限金额;NOT NULL" json:"out_day_upper_limit"`
	OutDayUpperNum     int    `gorm:"column:out_day_upper_num;type:int(11);default:0;comment:每天代收上限笔数;NOT NULL" json:"out_day_upper_num"`
	FixedInAmount      int    `gorm:"column:fixed_in_amount;type:int(11);default:0;comment:固定手续费代收(分);NOT NULL" json:"fixed_in_amount"`
	FixedOutAmount     int    `gorm:"column:fixed_out_amount;type:int(11);default:0;comment:固定手续费代付(分);NOT NULL" json:"fixed_out_amount"`
	FixedCurrency      string `gorm:"column:fixed_currency;type:varchar(64);comment:固定手续费币种;NOT NULL" json:"fixed_currency"`
	Remark             string `gorm:"column:remark;type:varchar(1024);comment:备注;NOT NULL" json:"remark"`
	CreateTime         uint64 `gorm:"column:create_time;type:bigint(20) unsigned;default:0;comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime         uint64 `gorm:"column:update_time;type:bigint(20) unsigned;default:0;comment:修改时间;NOT NULL" json:"update_time"`
}

func (m *AspMerchantProjectConfig) TableName() string {
	return "asp_merchant_project_config"
}
