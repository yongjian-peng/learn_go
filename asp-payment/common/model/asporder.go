package model

// 订单表openid: 表示用户在微信系统中对当前渠道公众号（channel_config.appid）的唯一身份标识，mp_appid:　商户公众号appid。 商户请求接口时选填。sub_openid：如果mp_appid 有值，表示用户在微信系统中对当前公众号（mp_appid）的唯一身份标识。
type AspOrder struct {
	Id             int     `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Sn             string  `gorm:"column:sn;type:varchar(32);comment:订单Serial Numner;NOT NULL" json:"sn"`
	MchId          int     `gorm:"column:mch_id;type:int(11);comment:外部商户id;NOT NULL" json:"mch_id"`
	DepartId       int     `gorm:"column:depart_id;type:int(11);default:0;comment:ASP（平台）商户id 关联departs表id;NOT NULL" json:"depart_id"`
	ChannelId      uint    `gorm:"column:channel_id;type:int(11) unsigned;default:0;comment:渠道id;NOT NULL" json:"channel_id"`
	ChannelMchId   string  `gorm:"column:channel_mch_id;type:varchar(45);comment:渠道商户上游给的id;NOT NULL" json:"channel_mch_id"`
	CurrencyId     int     `gorm:"column:currency_id;type:int(11);default:0;comment:币种ID;NOT NULL" json:"currency_id"`
	MchProjectId   int     `gorm:"column:mch_project_id;type:int(11);default:0;comment:外部商家项目ID;NOT NULL" json:"mch_project_id"`
	UserId         uint    `gorm:"column:user_id;type:int(11) unsigned;default:0;comment:用户Id;NOT NULL" json:"user_id"`
	Body           string  `gorm:"column:body;type:varchar(128);comment:商品介绍;NOT NULL" json:"body"`
	TransactionId  string  `gorm:"column:transaction_id;type:varchar(45);comment:交易id，渠道返回;NOT NULL" json:"transaction_id"`
	AdapterTransId string  `gorm:"column:adapter_trans_id;type:varchar(50);comment: 第三方交易id,如yedpay id;NOT NULL" json:"adapter_trans_id"`
	OutTradeNo     string  `gorm:"column:out_trade_no;type:varchar(45);comment:外部商家订单号;NOT NULL" json:"out_trade_no"`
	FeeType        string  `gorm:"column:fee_type;type:varchar(45);comment:币种;NOT NULL" json:"fee_type"`
	TotalFee       uint    `gorm:"column:total_fee;type:int(11) unsigned;default:0;comment:支付金额（标价金额，例如：卢币）;NOT NULL" json:"total_fee"`
	Discount       uint    `gorm:"column:discount;type:int(11) unsigned;default:0;comment:优惠金额;NOT NULL" json:"discount"`
	CashFeeType    string  `gorm:"column:cash_fee_type;type:varchar(45);comment:现金支付货币类型;NOT NULL" json:"cash_fee_type"`
	CashFee        int     `gorm:"column:cash_fee;type:int(11);default:0;comment:现金支付金额（实际支付金额）;NOT NULL" json:"cash_fee"`
	Rate           uint    `gorm:"column:rate;type:int(11) unsigned;default:0;comment:汇率;NOT NULL" json:"rate"`
	FeeRate        float64 `gorm:"column:fee_rate;type:decimal(5,5);comment:费率，成本价;NOT NULL" json:"fee_rate"`
	TotalChargeFee int     `gorm:"column:total_charge_fee;type:int(11);default:0;comment:总计手续费（手续费+固定手续费）;NOT NULL" json:"total_charge_fee"`
	ChargeFee      int     `gorm:"column:charge_fee;type:int(11);default:0;comment:手续费:单位分 存的正数;NOT NULL" json:"charge_fee"`
	FixedAmount    int     `gorm:"column:fixed_amount;type:int(11);default:0;comment:固定手续费 存的正数;NOT NULL" json:"fixed_amount"`
	FixedCurrency  string  `gorm:"column:fixed_currency;type:varchar(64);comment:固定手续费币种;NOT NULL" json:"fixed_currency"`
	SpbillCreateIp string  `gorm:"column:spbill_create_ip;type:varchar(45);comment:客户端ip;NOT NULL" json:"spbill_create_ip"`
	ClientIp       string  `gorm:"column:client_ip;type:varchar(32);comment:用户ip;NOT NULL" json:"client_ip"`
	NotifyUrl      string  `gorm:"column:notify_url;type:varchar(256);comment:通知回调url  交易类型为JSAPI、NATIVE必填;NOT NULL" json:"notify_url"`
	AuthCode       string  `gorm:"column:auth_code;type:varchar(512);comment:pay授权码：交易类型为 MICROPAY 必填;NOT NULL" json:"auth_code"`
	DeviceInfo     string  `gorm:"column:device_info;type:varchar(45);comment:设备号：由商户自定义，微信提供。;NOT NULL" json:"device_info"`
	CustomerId     string  `gorm:"column:customer_id;type:varchar(64);comment:客户相关信息客户ID;NOT NULL" json:"customer_id"`
	CustomerName   string  `gorm:"column:customer_name;type:varchar(255);comment:客户相关信息客户名称;NOT NULL" json:"customer_name"`
	CustomerEmail  string  `gorm:"column:customer_email;type:varchar(255);comment:客户相关信息客户邮箱;NOT NULL" json:"customer_email"`
	CustomerPhone  string  `gorm:"column:customer_phone;type:varchar(255);comment:客户相关信息客户手机号;NOT NULL" json:"customer_phone"`
	OrderToken     string  `gorm:"column:order_token;type:varchar(255);comment:订单token;NOT NULL" json:"order_token"`
	PaymentLink    string  `gorm:"column:payment_link;type:varchar(1024);comment:支付地址;NOT NULL" json:"payment_link"`
	PaymentsUrl    string  `gorm:"column:payments_url;type:varchar(1024);comment:付款列表地址;NOT NULL" json:"payments_url"`
	SettlementsUrl string  `gorm:"column:settlements_url;type:varchar(1024);comment:结算信息地址;NOT NULL" json:"settlements_url"`
	Attach         string  `gorm:"column:attach;type:varchar(128);comment:附加数据;NOT NULL" json:"attach"`
	TimeStart      string  `gorm:"column:time_start;type:char(14);comment:交易开始时间;NOT NULL" json:"time_start"`
	TimeExpire     string  `gorm:"column:time_expire;type:char(14);comment:交易过期时间;NOT NULL" json:"time_expire"`
	TimeEnd        string  `gorm:"column:time_end;type:char(14);comment:交易结束时间;NOT NULL" json:"time_end"`
	GoodsTag       string  `gorm:"column:goods_tag;type:varchar(45);comment:商品标签;NOT NULL" json:"goods_tag"`
	ProductId      string  `gorm:"column:product_id;type:varchar(45);comment:产品id;NOT NULL" json:"product_id"`
	SceneInfo      string  `gorm:"column:scene_info;type:varchar(45);comment:场景信息，选填字段，详见接口文档。;NOT NULL" json:"scene_info"`
	TradeState     string  `gorm:"column:trade_state;type:varchar(45);comment:交易状态：PENDING 未支付 SUCCESS, 支付成功 CANCELLED, 支付取消 USERPAYING, 支付中 PAYERROR,支付异常  FAILED 支付失败;NOT NULL" json:"trade_state"`
	BankType       string  `gorm:"column:bank_type;type:varchar(45);comment:支付扣款银行;NOT NULL" json:"bank_type"`
	IsCallback     uint    `gorm:"column:is_callback;type:tinyint(4) unsigned;default:0;comment:当 notify_url 存在的时候，表示是否有通知商户。;NOT NULL" json:"is_callback"`
	CreateTime     uint64  `gorm:"column:create_time;type:bigint(20) unsigned;default:0;comment:记录创建时间（订单已北京时区记录 东8区）;NOT NULL" json:"create_time"`
	FinishTime     uint64  `gorm:"column:finish_time;type:bigint(20) unsigned;default:0;comment:订单支付成功，更新状态时间。;NOT NULL" json:"finish_time"`
	SettlementTime uint64  `gorm:"column:settlement_time;type:bigint(20) unsigned;default:0;comment:订单结算时间;NOT NULL" json:"settlement_time"`
	IsBill         uint    `gorm:"column:is_bill;type:tinyint(4) unsigned;default:0;comment:是否已对账（同微信等平台）;NOT NULL" json:"is_bill"`
	Wallet         string  `gorm:"column:wallet;type:varchar(20);default:cn;comment:钱包地区:cn,hk;NOT NULL" json:"wallet"`
	Region         string  `gorm:"column:region;type:varchar(20);default:hk;comment:业务地区:hk,au,us,hk;NOT NULL" json:"region"`
	Provider       string  `gorm:"column:provider;type:varchar(32);default:wechat;comment:支付提供商: paytm cashfree paypal;NOT NULL" json:"provider"`
	Adapter        string  `gorm:"column:adapter;type:varchar(32);default:wechat;comment:支付业务适配api提供方: wechat,yedpay,ecpay;NOT NULL" json:"adapter"`
	TradeType      string  `gorm:"column:trade_type;type:varchar(45);comment:交易类型（　1:H5　2: PAYOUT）
;NOT NULL" json:"trade_type"`
	H5Type             string `gorm:"column:h5_type;type:varchar(32);comment:H5 支付类型（H5 WAPPAY）;NOT NULL" json:"h5_type"`
	Qrcode             string `gorm:"column:qrcode;type:text;comment:二维码(此字段可能很长);NOT NULL" json:"qrcode"`
	Note               string `gorm:"column:note;type:varchar(255);comment:支付备注，如失败原因等;NOT NULL" json:"note"`
	ReturnUrl          string `gorm:"column:return_url;type:varchar(255);comment:支付完成跳转URL;NOT NULL" json:"return_url"`
	Utr                string `gorm:"column:utr;type:varchar(64);comment:UTR 流水编号;NOT NULL" json:"utr"`
	UnionpayAppend     string `gorm:"column:unionpay_append;type:text;comment:支付附加数据;NOT NULL" json:"unionpay_append"`
	SupplierReturnCode string `gorm:"column:supplier_return_code;type:varchar(30);comment:供应商返回的错误码;NOT NULL" json:"supplier_return_code"`
	SupplierReturnMsg  string `gorm:"column:supplier_return_msg;type:text;comment:供应商返回的信息;NOT NULL" json:"supplier_return_msg"`
}

func (m *AspOrder) TableName() string {
	return "asp_order"
}
