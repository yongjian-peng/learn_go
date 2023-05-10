package model

// 渠道商户支付方式关系表
type AspChannelDepartTradeType struct {
	Id                 int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	DepartId           uint   `gorm:"column:depart_id;type:int(10) unsigned;default:0;comment:蓝海（平台）商户或代理商id，对应 admin_departs 主键 ID;NOT NULL" json:"depart_id"`
	ChannelId          int    `gorm:"column:channel_id;type:int(11);default:0;comment:渠道id(y2p)　对应　channel_config 主键 ID;NOT NULL" json:"channel_id"`
	CurrencyId         int    `gorm:"column:currency_id;type:int(11);comment:币种id 取渠道的币种;NOT NULL" json:"currency_id"`
	Provider           string `gorm:"column:provider;type:varchar(50);comment:支付渠道大类（例如 paytm cashfree）筛选;NOT NULL" json:"provider"`
	Payment            string `gorm:"column:payment;type:varchar(50);default:wechat.wechat;comment:firstpay.h5 zpay.app 筛选;NOT NULL" json:"payment"`
	TradeType          string `gorm:"column:trade_type;type:varchar(20);comment:允许的业务类型 H5 筛选 APP ;NOT NULL" json:"trade_type"`
	H5Type             string `gorm:"column:h5_type;type:varchar(255);comment:H5 支付类型（H5 WAPPAY）;NOT NULL" json:"h5_type"`
	InFeeRate          string `gorm:"column:in_fee_rate;type:varchar(10);comment:费率：收款费率，即本级（商户）的手续费，下级（代理）的成本价。保存 0.001 ~ 0.999 之间的值;NOT NULL" json:"in_fee_rate"`
	OutFeeRate         string `gorm:"column:out_fee_rate;type:varchar(10);comment:费率 代付支出费率 保存 0.001 ~ 0.999 之间的值;NOT NULL" json:"out_fee_rate"`
	DayUpperLimit      int    `gorm:"column:day_upper_limit;type:int(11);default:0;comment:每天收款或者提现上限金额 单位分;NOT NULL" json:"day_upper_limit"`
	UpperLimit         int    `gorm:"column:upper_limit;type:int(11);default:0;comment:每笔收款或者提现上限金额 单位分;NOT NULL" json:"upper_limit"`
	LowerLimit         int    `gorm:"column:lower_limit;type:int(11);default:0;comment:每笔收款或者提现下限金额 单位分;NOT NULL" json:"lower_limit"`
	FixedAmount        int    `gorm:"column:fixed_amount;type:int(11);default:0;comment:固定手续费 单位分;NOT NULL" json:"fixed_amount"`
	FixedCurrency      string `gorm:"column:fixed_currency;type:varchar(64);comment:固定手续费币种;NOT NULL" json:"fixed_currency"`
	InFeeRateUpdating  string `gorm:"column:in_fee_rate_updating;type:varchar(10);default:0;comment:更新收款的费率 预留;NOT NULL" json:"in_fee_rate_updating"`
	OutFeeRateUpdating string `gorm:"column:out_fee_rate_updating;type:varchar(10);default:0;comment:更新提现费率 预留;NOT NULL" json:"out_fee_rate_updating"`
	Status             int    `gorm:"column:status;type:int(11);default:0;comment:0：提交申请 1：申请成功，2：退回，3：更新;NOT NULL" json:"status"`
}

func (m *AspChannelDepartTradeType) TableName() string {
	return "asp_channel_depart_tradetype"
}
