package zpayimpl

import (
	"asp-payment/common/pkg/constant"
)

// 接口文档地址 https://www.yuque.com/docs/share/1b09478d-665e-4c08-85f3-a0be4d437299?#

const (
	Success         = 200                     // api code 成功状态码 为什么不使用
	DataSuccessCode = "0000"                  // code 成功状态码
	baseUrlProd     = "http://api.cxhd01.com" // URL 生产环境 http://api.cxhd01.com http://api.zpaypro.com

	orderCreate          = "/pay/order"         // 创建代收 POST
	orderPayout          = "/pay/withdraw"      // 创建代付
	queryOrder           = "/pay/queryOrder"    // 查询代收
	queryPayout          = "/pay/queryWithdraw" // 查询提现结果
	queryMerchantAccount = "/pay/getPartner"    // 查询账户详情

	SignType_HMAC_SHA256 = "HMAC-SHA256"
	SignType_MD5         = "MD5"

	StatusCreate     = 0 // 状态（0：创建订单；1：代收中；2：代收成功；3：代收失败）
	StatusUserPaying = 1
	StatusSuccess    = 2
	StatusFailed     = 3

	StatusPayoutCreate     = 0 // 状态（0：创建订单；1：代收中；2：代收成功；3：代收失败）
	StatusPayoutUserPaying = 1 //  代付中
	StatusPayoutSuccess    = 2 //  代付成功
	StatusPayoutFailed     = 3 //  代付失败
	ClientInitErr          = "ZPay 初始化失败"
	DoCurlErr              = "ZPay 请求上游失败"
	SaveOrderErr           = "ZPay 写入数据失败，请重试"
)

// 做上游返回的收款状态 和 系统收款状态的映射
// 输入 上游的状态 返回 系统收款的状态
func GetZPayPaymentStatus(code int) string {
	var MapZPayPaymentStatus = make(map[int]string)
	MapZPayPaymentStatus[0] = constant.ORDER_TRADE_STATE_PENDING
	MapZPayPaymentStatus[1] = constant.ORDER_TRADE_STATE_USERPAYING
	MapZPayPaymentStatus[2] = constant.ORDER_TRADE_STATE_SUCCESS
	MapZPayPaymentStatus[3] = constant.ORDER_TRADE_STATE_FAILED
	status := ""
	if ZPayStatus, ok := MapZPayPaymentStatus[code]; ok {
		status = ZPayStatus
	} else {
		status = constant.ORDER_TRADE_STATE_PAYERROR
	}

	return status
}

// 做上游返回的收款状态 和 系统收款状态的映射
// 输入 上游的状态 返回 系统收款的状态
func GetZPayPaymentCallBackStatus(code int) string {
	var MapZPayPaymentStatus = make(map[int]string)
	MapZPayPaymentStatus[0] = constant.ORDER_TRADE_STATE_FAILED
	MapZPayPaymentStatus[1] = constant.ORDER_TRADE_STATE_SUCCESS
	status := ""
	if ZPayStatus, ok := MapZPayPaymentStatus[code]; ok {
		status = ZPayStatus
	} else {
		status = constant.ORDER_TRADE_STATE_PAYERROR
	}

	return status
}

// 做上游返回的收款状态 和 系统收款状态的映射
// 输入 上游的状态 返回 系统收款的状态
func GetZPayPayoutCallBackStatus(code int) string {
	var MapZPayPaymentStatus = make(map[int]string)
	MapZPayPaymentStatus[0] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED
	MapZPayPaymentStatus[1] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS
	status := ""
	if ZPayStatus, ok := MapZPayPaymentStatus[code]; ok {
		status = ZPayStatus
	} else {
		status = constant.ORDER_TRADE_STATE_PAYERROR
	}

	return status
}

func GetZPayPaymentTradeState(trade_state string) (code int) {
	var MapZPayPaymentStatus = make(map[string]int)
	MapZPayPaymentStatus[constant.ORDER_TRADE_STATE_PENDING] = 0
	MapZPayPaymentStatus[constant.ORDER_TRADE_STATE_USERPAYING] = 1
	MapZPayPaymentStatus[constant.ORDER_TRADE_STATE_SUCCESS] = 2
	MapZPayPaymentStatus[constant.ORDER_TRADE_STATE_FAILED] = 3
	if ZPayStatus, ok := MapZPayPaymentStatus[trade_state]; ok {
		code = ZPayStatus
	} else {
		code = 0
	}
	return
}

// 做上游返回的提现状态 和 系统提现状态的映射
// 输入 上游的状态 返回 系统提现的状态
func GetZPayPayoutStatus(code int) string {
	var MapZPayPayoutStatus = make(map[int]string, 4)
	MapZPayPayoutStatus[0] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING // 代付中 == 创建代付订单
	MapZPayPayoutStatus[1] = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING // 代付中
	MapZPayPayoutStatus[2] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS // 待付成功
	MapZPayPayoutStatus[3] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED  // 代付失败
	status := ""
	if ZPayStatus, ok := MapZPayPayoutStatus[code]; ok {
		status = ZPayStatus
	} else {
		status = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	}

	return status
}

func GetZPayPayoutTradeState(trade_state string) (code int) {
	var MapZPayPayoutStatus = make(map[string]int)
	MapZPayPayoutStatus[constant.PAYOUT_TRADE_STATE_PENDING] = 0
	MapZPayPayoutStatus[constant.PAYOUT_TRADE_STATE_USERPAYING] = 1
	MapZPayPayoutStatus[constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS] = 2
	MapZPayPayoutStatus[constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED] = 3
	if ZPayStatus, ok := MapZPayPayoutStatus[trade_state]; ok {
		code = ZPayStatus
	} else {
		code = 0
	}
	return
}
