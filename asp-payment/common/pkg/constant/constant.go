package constant

type CurrencyActionType string

const (
	RegisterTypeMiniProgram = 1

	ProviderSunny = "sunny"

	TradeTypeFistPay       = "firstpay"
	TradeTypeZPay          = "zpay"
	TradeTypeSevenEight    = "78pay"
	TradeTypeFynzonPay     = "fynzonpay"
	TradeTypeAmarquickPay  = "amarquickpay"
	TradeTypeMyPay         = "mypay"
	TradeTypeHaoDaPay      = "haodapay"
	SupplierInitClientCode = "initcode"      // 支付渠道初始化 client 对应的code码
	SupplierHttpCode       = "httpcode"      // 支付渠道http响应 对应的code码
	SupplierHttpErrorCode  = "httperrorcode" // 支付渠道http 有错误 响应 对应的code码

	// 对于金额的变动操作
	SEND_TYPE_ORDER_SUCCESS        = "order_success"        // 代收成功
	SEND_TYPE_PAYOUT_APPLY         = "payout_apply"         // 代付申请
	SEND_TYPE_PAYOUT_AUDIT_SUCCESS = "payout_audit_success" // 代付审核成功
	SEND_TYPE_PAYOUT_AUDIT_FAILED  = "payout_audit_failed"  // 代付审核失败
	SEND_TYPE_PAYOUT_SUCCESS       = "payout_success"       // 代付上游成功
	SEND_TYPE_PAYOUT_FAILED        = "payout_failed"        // 代付上游失败
	SEND_TYPE_WITHDRAWAL_APPLY     = "withdrawal_apply"     // 提现申请
	SEND_TYPE_WITHDRAWAL_SUCCESS   = "withdrawal_success"   // 提现成功
	SEND_TYPE_WITHDRAWAL_FAILED    = "withdrawal_failed"    // 提现失败

	//支付方式
	TradeType_H5     = "H5"
	TradeType_PAYOUT = "PAYOUT"
	H5Type_H5        = "H5"
	H5Type_WAPPAY    = "WAPPAY"

	ORDER_STATISTIC_DATA_NAME_ALL = "ALL"
	ORDER_STATISTIC_DATA_NAME_MCH = "MCH"

	// 测试环境忽略的 projectid 1000012 1000030
	IGNORE_MERCHANT_PROJECT_ID = "1000012"

	// 代收 代付 同步状态，脚本请求每次处理的最大数量
	OrderCrondMaxUpdateCount = 300
	// 代付 审核 脚本请求每次处理最大数量
	PayoutAuditCrondMaxCount = 300
)
