package firstpayimpl

import "asp-payment/common/pkg/constant"

const (
	Success         = 200
	DataSuccessCode = "200" // code 成功状态码
	HeaderSignature = "Signature"
	HeaderAppId     = "AppId"
	// URL
	baseUrlProd = "https://payment.rummybank.com" // 生产环境

	OrderCreateUrl       = "/v1/platform/collect_order"                     // 创建订单 POST
	orderPayoutURL       = "/v1/platform/payout"                            // 创建提现订单
	retryCollectCallback = "/v1/platform/retry_collect_callback?order_id="  // 重试充值回调
	retryPayoutCallback  = "/v1/platform/retry_payout_callback?order_id="   // 重试提现回调
	queryPayout          = "/v1/platform/inquiry_payout_status?order_id=%s" // 查询提现结果
	queryMerchantAccount = "/v1/platform/inquiry_account"                   // 查询账户详情

	SignType_HMAC_SHA256 = "HMAC-SHA256"

	StatusUserPaying = 0
	StatusSuccess    = 1
	StatusFailed     = 2

	StatusPayoutUserPaying = 0 // 提现状态 支付中
	StatusPayoutSuccess    = 1 // 提现状态 成功
	StatusPayoutFailed     = 2 // 提现状态 失败
	ClientInitErr          = "FirstPay 初始化失败"
	DoCurlErr              = "FirstPay 请求上游失败"
	SaveOrderErr           = "FirstPay 写入数据失败，请重试"
)

// 做上游返回的收款状态 和 系统收款状态的映射
// 输入 上游的状态 返回 系统收款的状态
func GetFirstPayPaymentStatus(code int) string {
	var MapFirstPayPaymentStatus = make(map[int]string)
	MapFirstPayPaymentStatus[0] = constant.ORDER_TRADE_STATE_PENDING
	MapFirstPayPaymentStatus[1] = constant.ORDER_TRADE_STATE_SUCCESS
	MapFirstPayPaymentStatus[2] = constant.ORDER_TRADE_STATE_FAILED
	status := ""
	if firstpayStatus, ok := MapFirstPayPaymentStatus[code]; ok {
		status = firstpayStatus
	} else {
		status = constant.ORDER_TRADE_STATE_PAYERROR
	}

	return status
}

// 做上游返回的提现状态 和 系统提现状态的映射
// 输入 上游的状态 返回 系统提现的状态
func GetFirstPayPayoutStatus(code int) string {
	var MapFirstPayPayoutStatus = make(map[int]string)
	MapFirstPayPayoutStatus[0] = constant.PAYOUT_TRADE_STATE_PENDING
	MapFirstPayPayoutStatus[1] = constant.PAYOUT_TRADE_STATE_SUCCESS
	MapFirstPayPayoutStatus[2] = constant.PAYOUT_TRADE_STATE_FAILED
	status := ""
	if firstpayStatus, ok := MapFirstPayPayoutStatus[code]; ok {
		status = firstpayStatus
	} else {
		status = constant.PAYOUT_TRADE_STATE_PENDING
	}

	return status
}
