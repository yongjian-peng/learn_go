package abcpay

import "asp-payment/common/pkg/constant"

const (
	Success               = 200
	DataSuccessCode       = "0" // code 成功状态码
	QueryOrderSuccessCode = "0" // code 成功状
	PayResponseJson       = "1" //代收银行编码
	BaseUrlProd           = "https://pay.abcpayapp.com"
	BasePayOutUrlProd     = "https://7d923247f37ba22a66fa01b1911bsdc4.abcpayapp.com"
	OrderCreateUrl        = "/index/recharge"
	QueryOrderUrl         = "/index/order"
	QueryPayoutUrl        = "/api/query/queryorder"
	PayOutCreateUrl       = "/index/api/payout"
)

// GetPaymentTradeState 做上游返回的收款状态 和 系统收款状态的映射
// 输入 上游的状态 返回 系统收款的状态
func GetPaymentTradeState(state string) string {
	var MapPaymentStatus = make(map[string]string)
	MapPaymentStatus["2"] = constant.ORDER_TRADE_STATE_PENDING
	MapPaymentStatus["1"] = constant.ORDER_TRADE_STATE_SUCCESS
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
	var MapPayoutStatus = make(map[string]string, 4)
	MapPayoutStatus["1"] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS // 成功 => 成功
	MapPayoutStatus["2"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED  // 失败 => 失败
	MapPayoutStatus["5"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING // 处理中 => 处理中
	status := ""
	if ZPayStatus, ok := MapPayoutStatus[code]; ok {
		status = ZPayStatus
	} else {
		status = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	}
	return status
}

func GetPayoutCallBackStatus(code string) string {
	var MapPayoutStatus = make(map[string]string, 4)
	MapPayoutStatus["1"] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS // 成功 => 成功
	MapPayoutStatus["2"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED  // 失败 => 失败
	MapPayoutStatus["5"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING // 处理中 => 处理中
	status := ""
	if ZPayStatus, ok := MapPayoutStatus[code]; ok {
		status = ZPayStatus
	} else {
		status = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	}
	return status
}
