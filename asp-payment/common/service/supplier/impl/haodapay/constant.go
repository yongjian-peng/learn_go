package haodapay

import "asp-payment/common/pkg/constant"

const (
	Success               = 200
	DataSuccessCode       = "200"  // code 成功状态码
	QueryOrderSuccessCode = "200"  // code 成功状态码
	OrderCurrency         = "INR"  // 付款币种
	PayResponseJson       = "json" //代收银行编码
	BaseUrlProd           = "https://api.haodapayments.com"
	BasePayoutUrlProd     = "https://kepler.haodapayments.com"
	BasePayoutProxyUrl    = "http://52.66.205.254:13128"
	OrderCreateUrl        = "/api/v3/collection"
	UpiValidate           = "/api/v1/upi/validate"
	QueryOrderUrl         = "/api/v3/collection/status"
	PayOutCreateUrl       = "/api/v1/payout/initiate"
	PayOutCreateUpiUrl    = "/api/v1/upi/payout/initiate"
	QueryPayoutUrl        = "/api/v1/payout/checkstatus"
)

// GetPaymentTradeState 做上游返回的收款状态 和 系统收款状态的映射
// 输入 上游的状态 返回 系统收款的状态
func GetPaymentTradeState(state string) string {
	var MapPaymentStatus = make(map[string]string)
	MapPaymentStatus["pending"] = constant.ORDER_TRADE_STATE_PENDING
	MapPaymentStatus["PENDING"] = constant.ORDER_TRADE_STATE_PENDING
	MapPaymentStatus["success"] = constant.ORDER_TRADE_STATE_SUCCESS
	MapPaymentStatus["SUCCESS"] = constant.ORDER_TRADE_STATE_SUCCESS
	status := ""
	if PayStatus, ok := MapPaymentStatus[state]; ok {
		status = PayStatus
	} else {
		status = constant.ORDER_TRADE_STATE_PAYERROR
	}

	return status
}

// GetPayoutStatus 做上游返回的提现状态 和 系统提现状态的映射
// 输入 上游的状态 返回 系统提现的状态
func GetPayoutStatus(code string) string {
	var payoutStatus = make(map[string]string, 4)
	payoutStatus["Credited"] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS   // 成功 => 成功
	payoutStatus["Failed"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED      // 失败 => 失败
	payoutStatus["Processing"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING // 处理中 => 处理中
	payoutStatus["Pending"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING    // 待处理 => 处理中
	status := ""
	if payStatus, ok := payoutStatus[code]; ok {
		status = payStatus
	} else {
		status = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	}
	return status
}

func GetPayoutCallBackStatus(code string) string {
	var payoutStatus = make(map[string]string, 4)
	payoutStatus["success"] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS // 成功 => 成功
	payoutStatus["SUCCESS"] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS // 成功 => 成功
	payoutStatus["failed"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED   // 失败 => 失败
	status := ""
	if payStatus, ok := payoutStatus[code]; ok {
		status = payStatus
	} else {
		status = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	}
	return status
}
