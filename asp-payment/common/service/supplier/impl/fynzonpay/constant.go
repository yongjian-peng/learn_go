package fynzonpay

import "asp-payment/common/pkg/constant"

const (
	Success                    = 200
	DataSuccessCode            = "success" // code 成功状态码
	AddBeneficiarySuccessCode  = "0000"    // code 成功状态码
	QueryOrderSuccessCode      = "0000"    // code 成功状态码
	PayoutSuccessCode          = "0000"    // code 成功状态码
	BaseUrlProd                = "https://bo.fynzonpay.com"
	OrderCreateUrl             = "/directapi.do"             // directapi.do payment.do
	PayoutBeneficiaryCreateUrl = "/payout/addbeneficiary.do" // 添加受益人
	QueryOrderUrl              = "/validate.do"              // /validate.do?transaction_id=1951061135&api_token=QZfMTEyMF8yMDE4MDcyNzEyMjgyNMj
	QueryPayoutUrl             = "/payout/payoutdetail.do"
	PayOutCreateUrl            = "/payout/sendpayout.do"
	CardSend                   = "curl"     //
	CardSend2                  = "CHECKOUT" //
	Action                     = "product"  //
	Curr                       = "INR"      //
	ProductName                = "Product"  //
	Error                      = "error"
)

// GetPaymentTradeState 做上游返回的收款状态 和 系统收款状态的映射
// 输入 上游的状态 返回 系统收款的状态
func GetPaymentTradeState(state string) string {
	var MapStatus = make(map[string]string)
	MapStatus["0"] = constant.ORDER_TRADE_STATE_PENDING
	MapStatus["8"] = constant.ORDER_TRADE_STATE_PENDING
	MapStatus["2"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["3"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["5"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["7"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["10"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["11"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["12"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["13"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["14"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["15"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["16"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["17"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["20"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["21"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["22"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["23"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["24"] = constant.ORDER_TRADE_STATE_PAYERROR
	MapStatus["1"] = constant.ORDER_TRADE_STATE_SUCCESS
	MapStatus["9"] = constant.ORDER_TRADE_STATE_SUCCESS
	status := ""
	if payStatus, ok := MapStatus[state]; ok {
		status = payStatus
	} else {
		status = constant.ORDER_TRADE_STATE_PAYERROR
	}

	return status
}

// GetPayoutStatus 做上游返回的提现状态 和 系统提现状态的映射
// 输入 上游的状态 返回 系统提现的状态
func GetPayoutStatus(code string) string {
	var MapPaymentStatus = make(map[string]string)
	MapPaymentStatus["0"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	MapPaymentStatus["3"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	MapPaymentStatus["1"] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS
	MapPaymentStatus["2"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED
	MapPaymentStatus["10"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED
	status := ""
	if payStatus, ok := MapPaymentStatus[code]; ok {
		status = payStatus
	} else {
		status = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED
	}

	return status
}

func GetPayoutCallBackStatus(code string) string {
	var MapPayoutStatus = make(map[string]string, 4)
	MapPayoutStatus["0"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	MapPayoutStatus["3"] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	MapPayoutStatus["1"] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS
	MapPayoutStatus["2"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED
	MapPayoutStatus["10"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED
	status := ""
	if payStatus, ok := MapPayoutStatus[code]; ok {
		status = payStatus
	} else {
		status = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	}
	return status
}
