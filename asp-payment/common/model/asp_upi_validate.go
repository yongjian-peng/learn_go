package model

// 验证upi合法表
type AspUpiValidate struct {
	Id                 int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	MchId              int    `gorm:"column:mch_id;type:int(11);comment:外部商户id;NOT NULL" json:"mch_id"`
	DepartId           int    `gorm:"column:depart_id;type:int(11);default:0;comment:ASP（平台）商户id 关联departs表id;NOT NULL" json:"depart_id"`
	ChannelId          uint   `gorm:"column:channel_id;type:int(11) unsigned;default:0;comment:渠道id;NOT NULL" json:"channel_id"`
	ChannelMchId       string `gorm:"column:channel_mch_id;type:varchar(45);comment:渠道商户上游给的id;NOT NULL" json:"channel_mch_id"`
	CurrencyId         int    `gorm:"column:currency_id;type:int(11);default:0;comment:币种ID;NOT NULL" json:"currency_id"`
	MchProjectId       int    `gorm:"column:mch_project_id;type:int(11);default:0;comment:外部商家项目ID;NOT NULL" json:"mch_project_id"`
	SpbillCreateIp     string `gorm:"column:spbill_create_ip;type:varchar(45);comment:客户端ip;NOT NULL" json:"spbill_create_ip"`
	CustomerId         string `gorm:"column:customer_id;type:varchar(64);comment:客户相关信息客户ID;NOT NULL" json:"customer_id"`
	CustomerName       string `gorm:"column:customer_name;type:varchar(255);comment:客户相关信息客户名称;NOT NULL" json:"customer_name"`
	CustomerEmail      string `gorm:"column:customer_email;type:varchar(255);comment:客户相关信息客户邮箱;NOT NULL" json:"customer_email"`
	CustomerPhone      string `gorm:"column:customer_phone;type:varchar(255);comment:客户相关信息客户手机号;NOT NULL" json:"customer_phone"`
	Vpa                string `gorm:"column:vpa;type:varchar(64);comment:vpa号,印度upi提现方式需要提供;NOT NULL" json:"vpa"`
	TradeState         string `gorm:"column:trade_state;type:varchar(45);comment:交易状态：PENDING 验证中 SUCCESS, 验证成功 ERROR,验证异常  FAILED 验证失败;NOT NULL" json:"trade_state"`
	Provider           string `gorm:"column:provider;type:varchar(32);comment:支付提供商: paytm cashfree paypal;NOT NULL" json:"provider"`
	Adapter            string `gorm:"column:adapter;type:varchar(32);comment:支付业务适配api提供方: firstpay zpay;NOT NULL" json:"adapter"`
	TradeType          string `gorm:"column:trade_type;type:varchar(32);comment:交易类型（　1:H5，2:APP 3:PAYOUT);NOT NULL" json:"trade_type"`
	CreateTime         uint64 `gorm:"column:create_time;type:bigint(20) unsigned;default:0;comment:记录创建时间（订单已北京时区记录 东8区）;NOT NULL" json:"create_time"`
	FinishTime         uint64 `gorm:"column:finish_time;type:bigint(20) unsigned;default:0;comment:订单支付成功，更新状态时间。;NOT NULL" json:"finish_time"`
	UnionpayAppend     string `gorm:"column:unionpay_append;type:text;comment:支付附加数据;NOT NULL" json:"unionpay_append"`
	SupplierReturnCode string `gorm:"column:supplier_return_code;type:varchar(30);comment:供应商返回的错误码;NOT NULL" json:"supplier_return_code"`
	SupplierReturnMsg  string `gorm:"column:supplier_return_msg;type:text;comment:供应商返回的信息;NOT NULL" json:"supplier_return_msg"`
}

func (m *AspUpiValidate) TableName() string {
	return "asp_upi_validate"
}
