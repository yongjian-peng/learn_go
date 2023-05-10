package model

type AspOrderStatistic struct {
	Id                int     `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Period            uint    `gorm:"column:period;type:tinyint(3) unsigned;default:1;comment:统计的周期: 1 daily 天,2 monthly 月;NOT NULL" json:"period"`
	Date              uint64  `gorm:"column:date;type:bigint(20) unsigned;default:0;comment:统计数据日期;NOT NULL" json:"date"`
	DataName          string  `gorm:"column:data_name;type:char(32);comment:数据类别名称比如ALL、MCH、CH;NOT NULL" json:"data_name"`
	ChannelId         int     `gorm:"column:channel_id;type:int(11);default:0;comment:渠道id;NOT NULL" json:"channel_id"`
	DepartId          uint    `gorm:"column:depart_id;type:int(10) unsigned;default:0;comment:平台商户ID;NOT NULL" json:"depart_id"`
	MchId             int     `gorm:"column:mch_id;type:int(10);default:0;comment:外部商户ID;NOT NULL" json:"mch_id"`
	CurrencyId        int     `gorm:"column:currency_id;type:int(10);comment:币种ID;NOT NULL" json:"currency_id"`
	MchProjectId      int     `gorm:"column:mch_project_id;type:int(11);default:0;comment:外部商户项目id;NOT NULL" json:"mch_project_id"`
	Adapter           string  `gorm:"column:adapter;type:varchar(32);comment:支付服务提供方，接口实现方 zpay firstpay;NOT NULL" json:"adapter"`
	TradeType         string  `gorm:"column:trade_type;type:char(18);comment:支付方式 H5 PAYOUT等详情 参考 config.dictionary.trade_type;NOT NULL" json:"trade_type"`
	TotalFee          int64   `gorm:"column:total_fee;type:bigint(20);default:0;comment:成交总金额（代收）;NOT NULL" json:"total_fee"`
	TotalNum          int     `gorm:"column:total_num;type:int(11);default:0;comment:成交笔数（代收）;NOT NULL" json:"total_num"`
	ProfitFee         int64   `gorm:"column:profit_fee;type:bigint(11);default:0;comment:成交总手续费（代收）;NOT NULL" json:"profit_fee"`
	SuccessRate       int     `gorm:"column:success_rate;type:int(11);default:0;comment:代收成功率（0-100 之间的值）;NOT NULL" json:"success_rate"`
	PayoutTotalFee    int64   `gorm:"column:payout_total_fee;type:bigint(20);default:0;comment:代付总金额;NOT NULL" json:"payout_total_fee"`
	PayoutNum         int     `gorm:"column:payout_num;type:int(11);default:0;comment:代付笔数;NOT NULL" json:"payout_num"`
	PayoutProfitFee   int64   `gorm:"column:payout_profit_fee;type:bigint(11);default:0;comment:代付手续费;NOT NULL" json:"payout_profit_fee"`
	PayoutSuccessRate int     `gorm:"column:payout_success_rate;type:int(10);default:0;comment:代付成功率（0-100 之间的值）;NOT NULL" json:"payout_success_rate"`
	NetProfitFee      int64   `gorm:"column:net_profit_fee;type:bigint(20);default:0;comment:实际应扣净手续费（成交总手续费 + 代收手续费）;NOT NULL" json:"net_profit_fee"`
	NetTotalFee       int64   `gorm:"column:net_total_fee;type:bigint(20);default:0;comment:实际应手动结算金额（（成交代收总金额－代收总金额）－　实际应扣净手续费）;NOT NULL" json:"net_total_fee"`
	StatisticObject   string  `gorm:"column:statistic_object;type:varchar(25);comment:数据统计对象描述：例如：平台数据，商户数据。;NOT NULL" json:"statistic_object"`
	FeeRate           float64 `gorm:"column:fee_rate;type:decimal(6,4);comment:参考费率;NOT NULL" json:"fee_rate"`
	Apt               int     `gorm:"column:apt;type:int(11);default:0;comment:Amount Per Transation 付款笔均金额 单位分;NOT NULL" json:"apt"`
	Other             string  `gorm:"column:other;type:varchar(255);comment:其他数据，以json格式;NOT NULL" json:"other"`
	CreateTime        uint64  `gorm:"column:create_time;type:bigint(20) unsigned;default:0;comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime        uint64  `gorm:"column:update_time;type:bigint(20) unsigned;default:0;comment:更新时间;NOT NULL" json:"update_time"`
	VolChangeRate     float64 `gorm:"column:vol_change_rate;type:decimal(10,4);comment:交易量环比变化率(待定);NOT NULL" json:"vol_change_rate"`
	AmountChangeRate  float64 `gorm:"column:amount_change_rate;type:decimal(10,4);comment:交易金额环比变化率（待定）;NOT NULL" json:"amount_change_rate"`
	AptChangeRate     float64 `gorm:"column:apt_change_rate;type:decimal(10,4);comment:笔均金额 环比变化率（待定）;NOT NULL" json:"apt_change_rate"`
}

func (m *AspOrderStatistic) TableName() string {
	return "asp_order_statistic"
}
