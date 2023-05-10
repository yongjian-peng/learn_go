package amarquickpay

import "asp-payment/common/pkg/constant"

const (
	Success               = 200
	DataSuccessCode       = "success" // code 成功状态码
	QueryOrderSuccessCode = "000"     // code 成功状态码
	BaseUrlProd           = "https://pg.amarquick.com"
	QueryOrderUrl         = "/v1/services/paymentServices/StatusApi"
	QueryPayoutUrl        = "/Payment_Dfpay_query.html"
	PayOutCreateUrl       = "/Payment_Dfpay_add.html"
)

// GetPaymentTradeState 做上游返回的收款状态 和 系统收款状态的映射
// 输入 上游的状态 返回 系统收款的状态
func GetPaymentTradeState(state string) string {
	var MapPaymentStatus = make(map[string]string)
	MapPaymentStatus["NOTPAY"] = constant.ORDER_TRADE_STATE_PENDING
	MapPaymentStatus["Captured"] = constant.ORDER_TRADE_STATE_SUCCESS
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
	MapPayoutStatus["3"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING // 处理中 => 处理中
	MapPayoutStatus["4"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING // 待处理 => 处理中
	MapPayoutStatus["5"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED  // 审核驳回 => 失败
	MapPayoutStatus["6"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING // 待审核 => 处理中
	MapPayoutStatus["7"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED  // 交易不存在 => 失败
	MapPayoutStatus["8"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING // 未知状态 => 处理中
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
	MapPayoutStatus["2"] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS // 成功 => 成功
	MapPayoutStatus["3"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED  // 失败 => 失败
	status := ""
	if ZPayStatus, ok := MapPayoutStatus[code]; ok {
		status = ZPayStatus
	} else {
		status = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	}
	return status
}
