package model

// 代付表
type AspPayout struct {
	Id                 int     `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Sn                 string  `gorm:"column:sn;type:varchar(32);comment:订单Serial Numner;NOT NULL" json:"sn"`
	DepartId           int     `gorm:"column:depart_id;type:int(11);default:0;comment:asp（平台）商户id 关联 departs表id;NOT NULL" json:"depart_id"`
	ChannelId          uint    `gorm:"column:channel_id;type:int(11) unsigned;default:0;comment:渠道id;NOT NULL" json:"channel_id"`
	CurrencyId         int     `gorm:"column:currency_id;type:int(11);default:0;comment:币种ID;NOT NULL" json:"currency_id"`
	MchId              uint    `gorm:"column:mch_id;type:int(11) unsigned;default:0;comment:外部商家账户id;NOT NULL" json:"mch_id"`
	MchProjectId       int     `gorm:"column:mch_project_id;type:int(11);default:0;comment:外部商家项目ID;NOT NULL" json:"mch_project_id"`
	UserId             uint    `gorm:"column:user_id;type:int(11) unsigned;default:0;comment:用户Id;NOT NULL" json:"user_id"`
	Body               string  `gorm:"column:body;type:varchar(128);comment:退款简介;NOT NULL" json:"body"`
	TransactionId      string  `gorm:"column:transaction_id;type:varchar(45);comment:交易id，上游渠道返回;NOT NULL" json:"transaction_id"`
	AdapterTransId     string  `gorm:"column:adapter_trans_id;type:varchar(50);comment: 第三方交易id,如yedpay id;NOT NULL" json:"adapter_trans_id"`
	OutTradeNo         string  `gorm:"column:out_trade_no;type:varchar(45);comment:外部商家订单号;NOT NULL" json:"out_trade_no"`
	FeeType            string  `gorm:"column:fee_type;type:varchar(45);comment:币种;NOT NULL" json:"fee_type"`
	TotalFee           uint    `gorm:"column:total_fee;type:int(11) unsigned;default:0;comment:支付金额（标价金额，例如：卢币）;NOT NULL" json:"total_fee"`
	CashFeeType        string  `gorm:"column:cash_fee_type;type:varchar(45);comment:现金支付货币类型;NOT NULL" json:"cash_fee_type"`
	CashFee            int     `gorm:"column:cash_fee;type:int(11);default:0;comment:现金支付金额（实际支付金额）;NOT NULL" json:"cash_fee"`
	Rate               uint    `gorm:"column:rate;type:int(11) unsigned;default:0;comment:汇率;NOT NULL" json:"rate"`
	FeeRate            float64 `gorm:"column:fee_rate;type:decimal(5,5);comment:费率，成本价;NOT NULL" json:"fee_rate"`
	TotalChargeFee     int     `gorm:"column:total_charge_fee;type:int(11);default:0;comment:总计手续费（手续费+固定手续费）;NOT NULL" json:"total_charge_fee"`
	ChargeFee          int     `gorm:"column:charge_fee;type:int(11);default:0;comment:手续费:单位分 存的正数;NOT NULL" json:"charge_fee"`
	FixedAmount        int     `gorm:"column:fixed_amount;type:int(11);default:0;comment:固定手续费 存的正数;NOT NULL" json:"fixed_amount"`
	FixedCurrency      string  `gorm:"column:fixed_currency;type:varchar(64);comment:固定手续费币种;NOT NULL" json:"fixed_currency"`
	SpbillCreateIp     string  `gorm:"column:spbill_create_ip;type:varchar(45);comment:客户端ip;NOT NULL" json:"spbill_create_ip"`
	ClientIp           string  `gorm:"column:client_ip;type:varchar(32);comment:用户ip;NOT NULL" json:"client_ip"`
	NotifyUrl          string  `gorm:"column:notify_url;type:varchar(256);comment:通知回调url  交易类型为JSAPI、NATIVE必填;NOT NULL" json:"notify_url"`
	DeviceInfo         string  `gorm:"column:device_info;type:varchar(45);comment:设备号：由商户自定义;NOT NULL" json:"device_info"`
	Vpa                string  `gorm:"column:vpa;type:varchar(255);comment:vpa号,印度upi提现方式需要提供;NOT NULL" json:"vpa"`
	Ifsc               string  `gorm:"column:ifsc;type:varchar(64);comment:ifsc-只有印度银行卡提现需要. 注意：非印度银行卡，该字段不要赋值;NOT NULL" json:"ifsc"`
	CustomerName       string  `gorm:"column:customer_name;type:varchar(255);comment:客户相关信息客户名称;NOT NULL" json:"customer_name"`
	CustomerEmail      string  `gorm:"column:customer_email;type:varchar(255);comment:客户相关信息客户邮箱;NOT NULL" json:"customer_email"`
	CustomerPhone      string  `gorm:"column:customer_phone;type:varchar(255);comment:客户相关信息客户手机号;NOT NULL" json:"customer_phone"`
	OrderToken         string  `gorm:"column:order_token;type:varchar(255);comment:订单token;NOT NULL" json:"order_token"`
	Attach             string  `gorm:"column:attach;type:varchar(128);comment:附加数据;NOT NULL" json:"attach"`
	TimeStart          string  `gorm:"column:time_start;type:char(14);comment:交易开始时间;NOT NULL" json:"time_start"`
	TimeExpire         string  `gorm:"column:time_expire;type:char(14);comment:交易过期时间;NOT NULL" json:"time_expire"`
	TradeState         string  `gorm:"column:trade_state;type:varchar(45);comment:交易状态：APPLY 待审核 RETURN 拒绝审核  FREEZE_SUCCESS 冻结成功 CHANNEL_PENDING 请求上游成功 返回支付中 CHANNEL_FAILED 请求上游成功 返回失败 CHANNEL_SUCCESS 请求上游成功 返回成功 SUCCESS, 解冻成功 + 代付成功 FAILED 解冻成功 + 代付失败 REVOKED, 解冻成功 + 代付取消;NOT NULL" json:"trade_state"`
	PayType            string  `gorm:"column:pay_type;type:varchar(64);comment:提现类型(枚举值) bank:提现到银行卡(印度地区支持) clabe:墨西哥特有的方式，入前确认 PIX: 巴西特有方式 upi: 仅印度地区;NOT NULL" json:"pay_type"`
	TimeEnd            string  `gorm:"column:time_end;type:char(14);comment:交易结束时间;NOT NULL" json:"time_end"`
	BankType           string  `gorm:"column:bank_type;type:varchar(45);comment:支付收款银行;NOT NULL" json:"bank_type"`
	BankCard           string  `gorm:"column:bank_card;type:varchar(255);comment:提现目标银行卡id/pix账号/clabe卡号，其他提现方式可以忽略参数;NOT NULL" json:"bank_card"`
	BankCode           string  `gorm:"column:bank_code;type:varchar(255);comment:银行编码-墨西哥clabe提现需要;NOT NULL" json:"bank_code"`
	BankUtr            string  `gorm:"column:bank_utr;type:varchar(255);comment:Bank UTR No;NOT NULL" json:"bank_utr"`
	Address            string  `gorm:"column:address;type:varchar(255);comment:允许收款人地址、字母数字和空格（但脚本、HTML 标签会被清理或删除;NOT NULL" json:"address"`
	City               string  `gorm:"column:city;type:varchar(64);comment:收款城市，仅限字母和空格;NOT NULL" json:"city"`
	IsCallback         uint    `gorm:"column:is_callback;type:tinyint(4) unsigned;default:0;comment:当 notify_url 存在的时候，表示是否有通知商户。;NOT NULL" json:"is_callback"`
	IsCheckout         int     `gorm:"column:is_checkout;type:tinyint(4);default:0;comment:当审核订单的时候，审核状态（0已申请 1已完成 2进行中）;NOT NULL" json:"is_checkout"`
	Provider           string  `gorm:"column:provider;type:varchar(32);default:wechat;comment:支付提供商: paytm cashfree paypal;NOT NULL" json:"provider"`
	Adapter            string  `gorm:"column:adapter;type:varchar(32);default:wechat;comment:支付业务适配api提供方: firstpay zpay;NOT NULL" json:"adapter"`
	TradeType          string  `gorm:"column:trade_type;type:varchar(32);comment:交易类型（　1:H5，2:APP 3:PAYOUT);NOT NULL" json:"trade_type"`
	Note               string  `gorm:"column:note;type:varchar(255);comment:支付备注，如失败原因等;NOT NULL" json:"note"`
	BeneficiaryId      string  `gorm:"column:beneficiary_id;type:varchar(255);comment:受益人标识id;NOT NULL" json:"beneficiary_id"`
	SupplierReturnCode string  `gorm:"column:supplier_return_code;type:varchar(30);comment:供应商返回的错误码;NOT NULL" json:"supplier_return_code"`
	SupplierReturnMsg  string  `gorm:"column:supplier_return_msg;type:text;comment:供应商返回的错误信息" json:"supplier_return_msg"`
	CreateTime         uint64  `gorm:"column:create_time;type:bigint(20) unsigned;default:0;comment:记录创建时间;NOT NULL" json:"create_time"`
	FinishTime         uint64  `gorm:"column:finish_time;type:bigint(20) unsigned;default:0;comment:订单付款成功，更新状态时间。;NOT NULL" json:"finish_time"`
}

func (m *AspPayout) TableName() string {
	return "asp_payout"
}
