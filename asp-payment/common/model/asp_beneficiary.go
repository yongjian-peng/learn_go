package model

// 商户代付受益人表
type AspBeneficiary struct {
	Id                  int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	MchProjectBenefiary string `gorm:"column:mch_project_benefiary;type:varchar(32);comment:商户的标识映射 受益人id;NOT NULL" json:"mch_project_benefiary"`
	MchId               int    `gorm:"column:mch_id;type:int(11);comment:外部商户id;NOT NULL" json:"mch_id"`
	DepartId            int    `gorm:"column:depart_id;type:int(11);default:0;comment:ASP（平台）商户id 关联departs表id;NOT NULL" json:"depart_id"`
	ChannelId           uint   `gorm:"column:channel_id;type:int(10) unsigned;default:0;comment:渠道id;NOT NULL" json:"channel_id"`
	ChannelMchId        string `gorm:"column:channel_mch_id;type:varchar(45);comment:渠道商户上游给的id;NOT NULL" json:"channel_mch_id"`
	CurrencyId          int    `gorm:"column:currency_id;type:int(10);default:0;comment:币种ID;NOT NULL" json:"currency_id"`
	MchProjectId        int    `gorm:"column:mch_project_id;type:int(11);default:0;comment:外部商家项目ID;NOT NULL" json:"mch_project_id"`
	SpbillCreateIp      string `gorm:"column:spbill_create_ip;type:varchar(45);comment:客户端ip;NOT NULL" json:"spbill_create_ip"`
	NotifyUrl           string `gorm:"column:notify_url;type:varchar(256);comment:通知回调url  交易类型为JSAPI、NATIVE必填;NOT NULL" json:"notify_url"`
	CustomerId          string `gorm:"column:customer_id;type:varchar(64);comment:客户相关信息客户ID;NOT NULL" json:"customer_id"`
	CustomerName        string `gorm:"column:customer_name;type:varchar(255);comment:客户相关信息客户名称;NOT NULL" json:"customer_name"`
	CustomerEmail       string `gorm:"column:customer_email;type:varchar(255);comment:客户相关信息客户邮箱;NOT NULL" json:"customer_email"`
	CustomerPhone       string `gorm:"column:customer_phone;type:varchar(255);comment:客户相关信息客户手机号;NOT NULL" json:"customer_phone"`
	Ifsc                string `gorm:"column:ifsc;type:varchar(64);comment:ifsc-只有印度银行卡提现需要. 注意：非印度银行卡，该字段不要赋值;NOT NULL" json:"ifsc"`
	BankCard            string `gorm:"column:bank_card;type:varchar(255);comment:提现目标银行卡id/pix账号/clabe卡号，其他提现方式可以忽略参数;NOT NULL" json:"bank_card"`
	BankCode            string `gorm:"column:bank_code;type:varchar(255);comment:银行编码-墨西哥clabe提现需要;NOT NULL" json:"bank_code"`
	TradeState          string `gorm:"column:trade_state;type:varchar(45);comment:交易状态：PENDING 未支付 SUCCESS, 支付成功 CANCELLED, 支付取消 USERPAYING, 支付中 PAYERROR,支付异常  FAILED 支付失败;NOT NULL" json:"trade_state"`
	Provider            string `gorm:"column:provider;type:varchar(32);comment:支付提供商: paytm cashfree paypal;NOT NULL" json:"provider"`
	Adapter             string `gorm:"column:adapter;type:varchar(32);comment:支付业务适配api提供方: firstpay zpay;NOT NULL" json:"adapter"`
	TradeType           string `gorm:"column:trade_type;type:varchar(32);comment:交易类型（　1:H5，2:APP 3:PAYOUT);NOT NULL" json:"trade_type"`
	BenefiaryId         string `gorm:"column:benefiary_id;type:varchar(255);comment:上游渠道受益人标识id;NOT NULL" json:"benefiary_id"`
	CreateTime          uint64 `gorm:"column:create_time;type:bigint(10) unsigned;default:0;comment:记录创建时间（订单已北京时区记录 东8区）;NOT NULL" json:"create_time"`
	FinishTime          uint64 `gorm:"column:finish_time;type:bigint(10) unsigned;default:0;comment:订单支付成功，更新状态时间。;NOT NULL" json:"finish_time"`
	UnionpayAppend      string `gorm:"column:unionpay_append;type:text;comment:支付附加数据;NOT NULL" json:"unionpay_append"`
	SupplierReturnCode  string `gorm:"column:supplier_return_code;type:varchar(30);comment:供应商返回的错误码;NOT NULL" json:"supplier_return_code"`
	SupplierReturnMsg   string `gorm:"column:supplier_return_msg;type:text;comment:供应商返回的信息;NOT NULL" json:"supplier_return_msg"`
}

func (m *AspBeneficiary) TableName() string {
	return "asp_beneficiary"
}
