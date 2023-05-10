package mypay

import "asp-payment/common/pkg/constant"

const (
	Success                    = 200
	DataSuccessCode            = "Success" // code 成功状态码
	CreatePayoutDataFailedCode = "ERR"     // code 失败状态码
	CreateOrderSuccessCode     = "true"    // code 创建动态二维码 成功状态码
	QueryOrderSuccessCode      = "000"     // code 成功状态码
	BaseUrlProd                = "http://api.mypay.zone"
	OrderCreateUrl             = "/api/DynamicQrCode" // /api/QrCode /api/DynamicQrCode
	OrderCreateRechargeUrl     = "/api/Services/Recharge"
	QueryOrderUrl              = "/api/Reports/txnStatus" // /api/Services/StatusCheck /api/Reports/txnStatus
	QueryPayoutUrl             = "/api/Reports/txnStatus"
	PayOutCreateUrl            = "/api/Payout"
)

// GetPaymentTradeState 做上游返回的收款状态 和 系统收款状态的映射
// 输入 上游的状态 返回 系统收款的状态
func GetPaymentTradeState(state string) string {
	var mapPaymentStatus = make(map[string]string)
	mapPaymentStatus["Pending"] = constant.ORDER_TRADE_STATE_PENDING
	mapPaymentStatus["PENDING"] = constant.ORDER_TRADE_STATE_PENDING
	mapPaymentStatus["Success"] = constant.ORDER_TRADE_STATE_SUCCESS
	mapPaymentStatus["SUCCESS"] = constant.ORDER_TRADE_STATE_SUCCESS
	mapPaymentStatus["Failed"] = constant.ORDER_TRADE_STATE_FAILED
	mapPaymentStatus["FAILED"] = constant.ORDER_TRADE_STATE_FAILED
	status := ""
	if payStatus, ok := mapPaymentStatus[state]; ok {
		status = payStatus
	} else {
		status = constant.ORDER_TRADE_STATE_FAILED
	}

	return status
}

// GetPayoutStatus 做上游返回的提现状态 和 系统提现状态的映射
// 输入 上游的状态 返回 系统提现的状态
func GetPayoutStatus(code string) string {
	var mapPayoutStatus = make(map[string]string, 4)
	mapPayoutStatus["TXN"] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS // 成功 => 成功
	mapPayoutStatus["ERR"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED  // 失败 => 失败
	mapPayoutStatus["TUP"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING // 处理中 => 处理中
	mapPayoutStatus["UP"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING  // 待处理 => 处理中
	status := ""
	if payStatus, ok := mapPayoutStatus[code]; ok {
		status = payStatus
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
