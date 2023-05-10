package haodapay

import (
	"asp-payment/api-server/req"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/service/supplier/interfaces"
	"fmt"
	"github.com/spf13/cast"
)

type PayImpl struct{}

func NewPayImpl() *PayImpl {
	return &PayImpl{}
}

// H5
// 入参 订单信息 渠道信息 项目信息
// 返回 支付的链接
func (p *PayImpl) H5(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (*interfaces.ScanData, *appError.Error) {
	var scanData interfaces.ScanData
	scanData.Code = DataSuccessCode
	client, err := NewClient(channelDepartInfo, requestId, constant.HaodaPayLogFileName)
	if err != nil {
		return &scanData, appError.CodeSupplierInitClientCode
	}
	amount := goutils.Fen2Yuan(cast.ToInt64(orderInfo.TotalFee))

	bm := make(model.BodyMap)
	bm.Set("order_id", orderInfo.Sn).
		Set("order_amount", amount).
		Set("order_currency", OrderCurrency)

	// 组装参数 初始化请求的结构体
	payRsq, cErr := client.CreateOrder(bm)
	// 请求 curl

	// 获取响应结果 初始化响应的结构体
	if cErr != nil {
		return &scanData, appError.CodeSupplierHttpErrorCode
	}

	if payRsq.StatusCode != 200 {
		return &scanData, appError.CodeSupplierHttpCode
	}
	if payRsq.CreateOrderBody.StatusCode != DataSuccessCode {
		scanData.Code = payRsq.CreateOrderBody.StatusCode
		scanData.Msg = payRsq.CreateOrderBody.Message
		return &scanData, appError.CodeSupplierChannelErrCode
	}
	// 赋值给统一返回结构体
	payRsq.GenerateCreateOrder(&scanData)
	paymentsUrl, _ := goutils.JsonEncode(payRsq.CreateOrderBody.Data)
	scanData.CashFeeType = orderInfo.CashFeeType
	scanData.FinishTime = int64(orderInfo.FinishTime)
	scanData.TradeState = orderInfo.TradeState
	scanData.Rate = cast.ToInt(orderInfo.Rate)
	scanData.Qrcode = orderInfo.Qrcode
	scanData.SettlementTime = orderInfo.SettlementTime
	scanData.PaymentsURL = paymentsUrl
	scanData.ReturnURL = orderInfo.ReturnUrl
	scanData.SettlementsURL = orderInfo.SettlementsUrl
	scanData.CashFeeType = orderInfo.CashFeeType
	scanData.BankType = orderInfo.BankType
	scanData.CashFee = cast.ToInt(orderInfo.TotalFee)

	return &scanData, nil
}

// WAPPAY First Pay WAPPAY 支付方式
// 入参 订单信息 渠道信息 项目信息
// 返回 支付的链接
func (p *PayImpl) WAPPAY(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (model.BodyMap, *appError.Error) {
	bm := make(model.BodyMap)
	client, err := NewClient(channelDepartInfo, requestId, constant.HaodaPayLogFileName)
	if err != nil {
		return bm, appError.CodeSupplierInitClientCode
	}

	amount := goutils.Fen2Yuan(cast.ToInt64(orderInfo.TotalFee))

	bm.Set("order_id", orderInfo.Sn).
		Set("order_amount", amount).
		Set("order_currency", OrderCurrency)

	// 组装参数 初始化请求的结构体
	payRsq, cErr := client.CreateOrder(bm)
	// 获取响应结果 初始化响应的结构体
	if cErr != nil {
		return bm, appError.CodeSupplierHttpErrorCode
	}

	if payRsq.StatusCode != 200 {
		return bm, appError.CodeSupplierHttpCode
	}
	if payRsq.CreateOrderBody.StatusCode != DataSuccessCode {
		return bm, appError.CodeSupplierChannelErrCode
	}

	paymentsUrl, _ := goutils.JsonEncode(payRsq.CreateOrderBody.Data)
	bm.Set("url", paymentsUrl).
		Set("transactionId", payRsq.CreateOrderBody.CreateOrderData.Reference)
	return bm, nil
}

// Payout 78Pay web 支付方式
// 入参 订单信息 渠道信息 项目信息
// 返回 支付的链接
func (p *PayImpl) Payout(requestId string, channelDepartInfo *model.AspChannelDepartConfig, payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutCreateData, *appError.Error) {

	var payoutCreateData interfaces.ThirdPayoutCreateData
	payoutCreateData.Code = DataSuccessCode
	client, err := NewClient(channelDepartInfo, requestId, constant.HaodaPayLogFileName)
	if err != nil {
		return &payoutCreateData, appError.CodeSupplierInternalChannelErrCode
	}

	amount := goutils.Fen2Yuan(cast.ToInt64(payoutInfo.TotalFee))
	bm := make(model.BodyMap)
	bm.Set("account_number", payoutInfo.BankCard).
		Set("account_ifsc", payoutInfo.Ifsc).
		Set("bankname", payoutInfo.BankCode).
		Set("confirm_acc_number", payoutInfo.BankCard).
		Set("requesttype", "IMPS").
		Set("beneficiary_name", payoutInfo.CustomerName). //
		Set("amount", amount).
		Set("narration", payoutInfo.Body).
		Set("reference", payoutInfo.Sn)
	// 生成签名
	// fmt.Println("bm--------------", bm)
	// 组装参数 初始化请求的结构体
	// 请求 curl

	payRsq, pErr := client.CreatePayout(bm)
	// 获取响应结果 初始化响应的结构体
	if pErr != nil {
		return &payoutCreateData, appError.CodeSupplierHttpErrorCode
	}
	// fmt.Println("bm32132--------------", bm)
	if payRsq.StatusCode != 200 {
		return &payoutCreateData, appError.CodeSupplierHttpCode
	}

	if payRsq.CreatePayoutBody.StatusCode != DataSuccessCode {
		payoutCreateData.Code = payRsq.CreatePayoutBody.StatusCode
		payoutCreateData.Msg = payRsq.CreatePayoutBody.Message
		payoutCreateData.TradeState = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED
		return &payoutCreateData, appError.CodeSupplierChannelErrCode
	}

	payoutCreateData.Msg = payRsq.StatusMsg
	payoutCreateData.TransactionID = payRsq.CreatePayoutBody.PayoutID
	payoutCreateData.CashFee = cast.ToInt(payoutInfo.TotalFee)
	payoutCreateData.CashFeeType = payoutInfo.CashFeeType
	payoutCreateData.FinishTime = int64(payoutInfo.FinishTime)
	payoutCreateData.TradeState = constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING
	payoutCreateData.CashFeeType = payoutInfo.CashFeeType
	payoutCreateData.BankType = payoutInfo.BankType
	return &payoutCreateData, nil

}

func (p *PayImpl) Web(requestId string, AspChannelDepart *model.AspChannelDepartConfig, orderInfo *model.AspOrder, channelDepartTradeType *req.DeptTradeTypeInfo, merchantInfo *model.AspMerchantProject) bool {
	var scanData interfaces.ScanData
	scanData.Code = "01"
	return true
}

// PayQuery 查询上游订单状态 需要的client 中 channel_depart 中的 config 字段
// 上游提供的是异步查询的方法，发送回调的消息，成功则等待回调通知订单结果
func (p *PayImpl) PayQuery(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (*interfaces.ThirdQueryData, *appError.Error) {

	var thirdQueryData interfaces.ThirdQueryData
	thirdQueryData.Code = DataSuccessCode //默认成功code

	payClient, err := NewClient(channelDepartInfo, requestId, constant.HaodaPayLogFileName)
	if err != nil {
		return &thirdQueryData, appError.CodeSupplierInternalChannelErrCode
	}

	// 如果没有传入 TransactionId 则返回最近的一条数据 会有错误的更新
	if orderInfo.TransactionId == "" {
		return &thirdQueryData, appError.CodeInvalidParamErrCode
	}

	bm := make(model.BodyMap)
	bm.Set("order_id", orderInfo.Sn)
	// 组装参数 初始化请求的结构体
	payRsp, qErr := payClient.QueryOrder(bm)
	// 获取响应结果 初始化响应的结构体
	if qErr != nil {
		return &thirdQueryData, appError.CodeSupplierHttpErrorCode
	}

	if payRsp.StatusCode != 200 {
		return &thirdQueryData, appError.CodeSupplierHttpCode
	}

	fmt.Println("payRsp.QueryOrderBody:", fmt.Sprintf("%+v", payRsp.QueryOrderBody))

	//有错误，设置错误消息
	if payRsp.QueryOrderBody.StatusCode != QueryOrderSuccessCode {
		thirdQueryData.Code = payRsp.QueryOrderBody.StatusCode
		// 错误的话，错误消息往上层传递
		thirdQueryData.Msg = payRsp.QueryOrderBody.Status
		return &thirdQueryData, nil
	}

	//根据上游支付状态，转换为本系统支付状态
	upstreamStatus := GetPaymentTradeState(payRsp.QueryOrderBody.Status) // 字符串 success
	params := map[string]string{}
	//默认完成时间取数据库的值
	params["finish_time"] = cast.ToString(orderInfo.FinishTime)
	params["status"] = orderInfo.TradeState // 默认为数据库的状态

	//上游是非成功状态，并且数据库的状态不是成功状态，更新数据库
	if upstreamStatus != constant.ORDER_TRADE_STATE_SUCCESS && orderInfo.TradeState != constant.ORDER_TRADE_STATE_SUCCESS {
		params["status"] = upstreamStatus // 设置最新的状态，成功状态不可逆转
	}

	//成功状态，获取上游完成时间
	if upstreamStatus == constant.ORDER_TRADE_STATE_SUCCESS {
		params["finish_time"] = goutils.Int642String(goutils.GetDateTimeUnix())
		params["status"] = upstreamStatus
		thirdQueryData.Utr = payRsp.QueryOrderBody.QueryOrderData.UTR
	}

	// 赋值给统一返回结构体
	thirdQueryData.TransactionID = orderInfo.TransactionId
	thirdQueryData.Msg = payRsp.QueryOrderBody.Status
	thirdQueryData.CashFee = cast.ToInt(orderInfo.TotalFee)
	thirdQueryData.CashFeeType = orderInfo.CashFeeType
	thirdQueryData.FinishTime = cast.ToInt64(params["finish_time"])
	thirdQueryData.TradeState = params["status"]
	return &thirdQueryData, nil
}

// PayoutQuery 查询上游订单状态 需要的client 中 channel_depart 中的 config 字段
// 上游提供的是异步查询的方法，发送回调的消息，成功则等待回调通知订单结果
func (p *PayImpl) PayoutQuery(requestId string, channelDepartInfo *model.AspChannelDepartConfig, payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutQueryData, *appError.Error) {

	var thirdPayoutQueryData interfaces.ThirdPayoutQueryData
	thirdPayoutQueryData.Code = DataSuccessCode

	payClient, err := NewClient(channelDepartInfo, requestId, constant.HaodaPayLogFileName)
	if err != nil {
		return &thirdPayoutQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	bm := make(model.BodyMap)
	bm.Set("payout_id", payoutInfo.TransactionId)
	// 组装参数 初始化请求的结构体
	//{"status":"success","msg":"apply successful","transaction_id":"L1208711582746222","balance":"97.00"}
	payRsp, qErr := payClient.QueryPayout(bm)
	// 获取响应结果 初始化响应的结构体
	if qErr != nil {
		return &thirdPayoutQueryData, appError.CodeSupplierHttpErrorCode
	}

	if payRsp.StatusCode != Success {
		return &thirdPayoutQueryData, appError.CodeSupplierHttpCode
	}

	//有错误，设置错误消息
	if payRsp.QueryPayoutBody.StatusCode != DataSuccessCode {
		thirdPayoutQueryData.Code = payRsp.QueryPayoutBody.StatusCode
		// 错误的话，错误消息往上层传递
		thirdPayoutQueryData.Msg = payRsp.QueryPayoutBody.Message
		return &thirdPayoutQueryData, nil
	}
	queryPayoutData := payRsp.QueryPayoutBody

	// 上游没有返回金额 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况

	upstreamStatus := GetPayoutStatus(queryPayoutData.QueryPayoutData.Status)
	params := map[string]string{}

	params["finish_time"] = cast.ToString(payoutInfo.FinishTime)
	params["transaction_id"] = payoutInfo.TransactionId
	params["status"] = payoutInfo.TradeState // 默认为数据库的状态

	//上游是非成功状态，并且数据库的状态不是成功状态，更新数据库
	fmt.Println("upstreamStatus:", upstreamStatus)
	fmt.Println("payoutInfo.TradeState:", payoutInfo.TradeState)
	if upstreamStatus != constant.PAYOUT_TRADE_STATE_SUCCESS && payoutInfo.TradeState != constant.PAYOUT_TRADE_STATE_SUCCESS {
		params["status"] = upstreamStatus // 设置最新的状态，成功状态不可逆转
	}

	//如果上游是成功状态，更新完成时间
	if upstreamStatus == constant.PAYOUT_TRADE_STATE_SUCCESS {
		params["finish_time"] = goutils.Int642String(goutils.GetDateTimeUnix())
		params["transaction_id"] = payoutInfo.TransactionId
		params["status"] = upstreamStatus // 设置最新的状态，成功状态不可逆转
		thirdPayoutQueryData.BankUtr = queryPayoutData.QueryPayoutData.UTR
	}

	// 赋值给统一返回结构体
	thirdPayoutQueryData.TransactionID = params["transaction_id"]
	thirdPayoutQueryData.CashFee = cast.ToInt(payoutInfo.TotalFee)
	thirdPayoutQueryData.CashFeeType = payoutInfo.FeeType
	thirdPayoutQueryData.FinishTime = cast.ToInt64(params["finish_time"])
	thirdPayoutQueryData.TradeState = params["status"]
	thirdPayoutQueryData.Msg = queryPayoutData.Status

	return &thirdPayoutQueryData, nil
}

// GetDepartAccountInfo 获取商户账户信息
func (p *PayImpl) GetDepartAccountInfo(requestId string, channelDepartInfo *model.AspChannelDepartConfig) (*interfaces.ThirdMerchantAccountQueryData, *appError.Error) {
	var thirdMerchantAccountQueryData interfaces.ThirdMerchantAccountQueryData
	thirdMerchantAccountQueryData.Code = DataSuccessCode

	return &thirdMerchantAccountQueryData, nil

}

func (p *PayImpl) PayNotify() {
	//TODO implement me
	panic("implement me")
}

// AddBeneficiary 添加受益人 需要的client 中 channel_depart 中的 config 字段
func (p *PayImpl) AddBeneficiary(requestId string, clientIp string, channelDepartInfo *model.AspChannelDepartConfig, reqBeneficiary *req.AspBeneficiary) (*interfaces.ThirdAddBeneficiary, *appError.Error) {
	return nil, appError.NotImplementedErrCode
}

// PayoutUpi 使用 upi 方式付款
func (p *PayImpl) PayoutUpi(requestId string, channelDepartInfo *model.AspChannelDepartConfig, payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutCreateData, *appError.Error) {
	var thirdData interfaces.ThirdPayoutCreateData
	thirdData.Code = DataSuccessCode //默认成功code

	client, err := NewClient(channelDepartInfo, requestId, constant.HaodaPayLogFileName)
	if err != nil {
		return &thirdData, appError.CodeSupplierInternalChannelErrCode
	}

	amount := goutils.Fen2Yuan(cast.ToInt64(payoutInfo.TotalFee))
	bm := make(model.BodyMap)
	bm.Set("vpa", payoutInfo.Vpa).
		Set("beneficiary_name", payoutInfo.CustomerName).
		Set("amount", amount).
		Set("narration", payoutInfo.Note).
		Set("reference", payoutInfo.Sn)

	// 组装参数 初始化请求的结构体
	payRsq, cErr := client.CreatePayoutUpi(bm)
	// 请求 curl

	// 获取响应结果 初始化响应的结构体
	if cErr != nil {
		return &thirdData, appError.CodeSupplierHttpErrorCode
	}

	if payRsq.StatusCode != 200 {
		return &thirdData, appError.CodeSupplierHttpCode
	}
	if payRsq.CreateOrderBody.StatusCode != DataSuccessCode {
		thirdData.Code = payRsq.CreatePayoutUpiBody.StatusCode
		thirdData.Msg = payRsq.CreatePayoutUpiBody.Message
		return &thirdData, appError.CodeSupplierChannelErrCode
	}
	// 赋值给统一返回结构体
	thirdData.Msg = payRsq.StatusMsg
	thirdData.TransactionID = payRsq.CreatePayoutUpiBody.PayoutID
	thirdData.CashFee = cast.ToInt(payoutInfo.TotalFee)
	thirdData.CashFeeType = payoutInfo.CashFeeType
	thirdData.FinishTime = int64(payoutInfo.FinishTime)
	thirdData.TradeState = GetPayoutStatus(payRsq.CreatePayoutUpiBody.Status)
	thirdData.CashFeeType = payoutInfo.CashFeeType
	thirdData.BankType = payoutInfo.BankType
	return &thirdData, nil
}

// UpiValidate 验证 upi 是否合法
func (p *PayImpl) UpiValidate(requestId string, channelDepartInfo *model.AspChannelDepartConfig, upiValidateInfo *req.AspPayoutUpiValidate) (*interfaces.ThirdUpiValidate, *appError.Error) {
	var thirdData interfaces.ThirdUpiValidate
	thirdData.Code = DataSuccessCode //默认成功code

	client, err := NewClient(channelDepartInfo, requestId, constant.HaodaPayLogFileName)
	if err != nil {
		return &thirdData, appError.CodeSupplierInternalChannelErrCode
	}

	bm := make(model.BodyMap)
	bm.Set("vpa", upiValidateInfo.Vpa)

	// 组装参数 初始化请求的结构体
	payRsq, cErr := client.UpiValidate(bm)
	// 请求 curl

	// 获取响应结果 初始化响应的结构体
	if cErr != nil {
		return &thirdData, appError.CodeSupplierHttpErrorCode
	}

	if payRsq.StatusCode != 200 {
		return &thirdData, appError.CodeSupplierHttpCode
	}
	if payRsq.CreateOrderBody.StatusCode != DataSuccessCode {
		thirdData.Code = payRsq.UpiValidateBody.StatusCode
		thirdData.Msg = payRsq.UpiValidateBody.Message
		return &thirdData, appError.CodeSupplierChannelErrCode
	}
	// 赋值给统一返回结构体
	thirdData.Msg = payRsq.UpiValidateBody.Status
	return &thirdData, nil
}
