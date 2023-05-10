package zpayimpl

import (
	"asp-payment/api-server/req"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/service/supplier/interfaces"
	"context"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type PayImpl struct{}

func NewZPayImpl() *PayImpl {
	return &PayImpl{}
}

var (
	ctx = context.Background()
)

func initClient(channelDepartInfo *model.AspChannelDepartConfig, request_id string) (*Client, error) {

	// fmt.Println("channelDepartInfo------------------------", channelDepartInfo.Config)
	// config 转 struct 字符串转 struct
	var channelConfigInfo model.AspChannelDepartConfigInfo
	_ = goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo)
	//fmt.Println("channelConfigInfo------------------------", channelConfigInfo.Zpay)
	// 初始化 client 需要参数 appid Signature
	client, err := NewClient(request_id, channelConfigInfo.PartnerId, channelConfigInfo.ApplicationId, channelConfigInfo.Signature, channelConfigInfo.PayoutSignature)

	if err != nil {
		return client, err
	}
	client.RequestId = request_id
	client.LogFileName = constant.ZPayLogFileName
	return client, nil
}

// H5 First Pay web 支付方式
// 入参 订单信息 渠道信息 项目信息
// 返回 支付的链接
func (p *PayImpl) H5(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (*interfaces.ScanData, *appError.Error) {

	var scanData interfaces.ScanData

	scanData.Code = DataSuccessCode
	client, err := initClient(channelDepartInfo, requestId)

	if err != nil {
		return &scanData, appError.CodeSupplierInitClientCode
	}
	// return scanData

	extraMap := make(map[string]interface{}, 3)
	extraMap["userName"] = orderInfo.CustomerName
	extraMap["userEmail"] = orderInfo.CustomerEmail
	extraMap["userPhone"] = orderInfo.CustomerPhone
	extraJson := goutils.ConvertToString(extraMap)

	// callbackUrl 参数待定 看怎么使用
	// payWay 固定 2 paytm 支付方式
	bm := make(model.BodyMap)
	bm.Set("partnerId", client.PartnerId).
		Set("applicationId", client.ApplicationId).
		Set("payWay", 2).
		Set("partnerOrderNo", orderInfo.Sn).
		Set("amount", orderInfo.TotalFee).
		Set("currency", orderInfo.FeeType).
		Set("clientIp", orderInfo.SpbillCreateIp).
		Set("notifyUrl", client.ZpayH5Notify).
		Set("subject", "play game").
		Set("body", "play game").
		Set("extra", extraJson).
		Set("version", "1.0")

	if orderInfo.ReturnUrl != "" {
		bm.Set("callbackUrl", orderInfo.ReturnUrl)
	}

	//scanData.PaymentLink = "https://www.mobilkwik.in/gateway/cashier/P221599591847519457280/a819fa6fb54e473da27244af69ddf31f"
	//scanData.TransactionID = "orderInfo.TransactionId"
	//scanData.CashFee = orderInfo.CashFee
	//scanData.CashFeeType = orderInfo.CashFeeType
	//scanData.FinishTime = int64(orderInfo.FinishTime)
	//scanData.TradeState = GetZPayPaymentStatus(2)
	//scanData.Rate = cast.ToInt(orderInfo.Rate)
	//scanData.Qrcode = orderInfo.Qrcode
	//scanData.SettlementTime = orderInfo.SettlementTime
	//scanData.PaymentsURL = orderInfo.PaymentsUrl
	//scanData.ReturnURL = orderInfo.ReturnUrl
	//scanData.SettlementsURL = orderInfo.SettlementsUrl
	//scanData.BankType = orderInfo.BankType
	//return &scanData, nil
	// 生成签名
	// return scanData
	// 组装参数 初始化请求的结构体
	zpRsq, errRes := client.CreateOrder(ctx, bm)
	// 获取响应结果 初始化响应的结构体
	if errRes != nil {
		return &scanData, appError.CodeSupplierHttpErrorCode
	}

	if zpRsq.Code != 200 {
		return &scanData, appError.CodeSupplierHttpCode
	}

	if zpRsq.Response.Code != DataSuccessCode { // 返回的 code 码 成功 0000
		scanData.Code = zpRsq.Response.Code
		// 错误的话，错误消息往上层传递
		scanData.Msg = zpRsq.Response.Message
		return &scanData, appError.CodeSupplierChannelErrCode
	}

	// 赋值给统一返回结构体
	zpRsq.Generate(&scanData)
	scanData.TransactionID = orderInfo.TransactionId
	scanData.CashFee = orderInfo.CashFee
	scanData.CashFeeType = orderInfo.CashFeeType
	scanData.FinishTime = int64(orderInfo.FinishTime)
	scanData.TradeState = orderInfo.TradeState
	scanData.Rate = cast.ToInt(orderInfo.Rate)
	scanData.Qrcode = orderInfo.Qrcode
	scanData.SettlementTime = orderInfo.SettlementTime
	scanData.PaymentsURL = orderInfo.PaymentsUrl
	scanData.ReturnURL = orderInfo.ReturnUrl
	scanData.SettlementsURL = orderInfo.SettlementsUrl
	scanData.BankType = orderInfo.BankType
	return &scanData, nil
}

// WAPPAY First Pay WAPPAY 支付方式
// 入参 订单信息 渠道信息 项目信息
// 返回 支付的链接
func (p *PayImpl) WAPPAY(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (model.BodyMap, *appError.Error) {
	return nil, nil
}

// Payout First Pay web 支付方式
// 入参 订单信息 渠道信息 项目信息
// 返回 支付的链接
func (p *PayImpl) Payout(requestId string, channelDepartInfo *model.AspChannelDepartConfig, payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutCreateData, *appError.Error) {
	// fmt.Println("payoutInfo--------------", payoutInfo)
	var payoutCreateData interfaces.ThirdPayoutCreateData

	payoutCreateData.Code = DataSuccessCode
	client, err := initClient(channelDepartInfo, requestId)
	if err != nil {
		//payoutCreateData.Status = constant.SupplierInitClientCode
		//payoutCreateData.Msg = ClientInitErr
		return &payoutCreateData, appError.CodeSupplierInternalChannelErrCode
	}
	//payoutCreateData.TransactionID = payoutInfo.TransactionId
	//payoutCreateData.CashFee = payoutInfo.CashFee
	//payoutCreateData.CashFeeType = payoutInfo.CashFeeType
	//payoutCreateData.FinishTime = int64(payoutInfo.FinishTime)
	//payoutCreateData.TradeState = GetZPayPayoutStatus(3)
	//payoutCreateData.CashFeeType = payoutInfo.CashFeeType
	//payoutCreateData.BankType = payoutInfo.BankType
	//return &payoutCreateData, nil

	// 收款方式 receiptMode 0：UPI 1：IMPS 2：Card 3：PIX
	// accountExtra1 特殊参数1 (当receiptMode 为1和2时，必传)
	// 特殊参数1    当参数receiptMode = 1（IMPS） IFSC CODE
	// 特殊参数1    当参数receiptMode = 2（Card） 分行代码
	// 特殊参数1    当参数receiptMode = 3（PIX）巴西税号（必须真实）
	// accountExtra2 特殊参数2 (当receiptMode 为1和2时，必传)
	// 特殊参数2    当参数receiptMode = 1（IMPS） 银行编码
	// 特殊参数2    当参数receiptMode = 2（Card） 银行编码
	// accountExtra3 特殊参数3 (当receiptMode 为2和3时，必传)
	// 特殊参数3    当参数receiptMode = 2（Card） 详细收款方式：CPF、CNPJ
	// 特殊参数3    当参数receiptMode = 3（PIX） 详细收款方式：CPF、CNPJ、PHONE、EMAIL、EVP

	bm := make(model.BodyMap)
	bm.Set("partnerId", client.PartnerId).
		Set("partnerWithdrawNo", payoutInfo.Sn).
		Set("amount", payoutInfo.TotalFee).
		Set("currency", payoutInfo.FeeType).
		Set("notifyUrl", client.ZpayPayoutNotify).
		Set("receiptMode", 1).
		Set("accountNumber", payoutInfo.BankCard). //
		Set("accountName", payoutInfo.CustomerName).
		Set("accountPhone", payoutInfo.CustomerPhone).
		Set("accountEmail", payoutInfo.CustomerEmail).
		Set("accountExtra1", payoutInfo.Ifsc).
		Set("accountExtra2", payoutInfo.BankCode). // ICICINBBXXX 银行编码
		// Set("accountExtra3", "CPF").
		Set("version", "1.0")
	// 生成签名
	// fmt.Println("bm--------------", bm)
	// 组装参数 初始化请求的结构体
	// 请求 curl
	zpRsq, errRes := client.CreatePayout(ctx, bm)
	// 获取响应结果 初始化响应的结构体
	// {"response": "{"code":"0000","message":"SUCCESS","data":null}
	if errRes != nil {
		return &payoutCreateData, appError.CodeSupplierHttpErrorCode
	}
	// fmt.Println("bm32132--------------", bm)
	if zpRsq.Code != 200 {
		return &payoutCreateData, appError.CodeSupplierHttpCode
	}
	if zpRsq.Response.Code != DataSuccessCode {
		payoutCreateData.Code = zpRsq.Response.Code
		// 错误的话，错误消息往上层传递
		payoutCreateData.Msg = zpRsq.Response.Message
		return &payoutCreateData, appError.CodeSupplierChannelErrCode
	}
	upstreamStatus := goutils.String2Int(zpRsq.Response.Data)

	zpRsq.Generate(&payoutCreateData)
	payoutCreateData.TransactionID = payoutInfo.TransactionId
	payoutCreateData.CashFee = payoutInfo.CashFee
	payoutCreateData.CashFeeType = payoutInfo.CashFeeType
	payoutCreateData.FinishTime = int64(payoutInfo.FinishTime)
	payoutCreateData.TradeState = GetZPayPayoutStatus(upstreamStatus)
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
	thirdQueryData.Code = DataSuccessCode
	client, err := initClient(channelDepartInfo, requestId)
	if err != nil {
		return &thirdQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	// fmt.Println("thirdQueryData------------", thirdQueryData)
	//return &thirdQueryData
	bm := make(model.BodyMap)
	bm.Set("partnerId", client.PartnerId).
		Set("applicationId", client.ApplicationId).
		Set("partnerOrderNo", orderInfo.Sn).
		Set("version", "1.0")
	//thirdQueryData.CashFee = cast.ToInt(orderInfo.TotalFee)
	//thirdQueryData.CashFeeType = orderInfo.CashFeeType
	//thirdQueryData.FinishTime = cast.ToInt64(orderInfo.FinishTime)
	//thirdQueryData.TradeState = GetZPayPaymentStatus(2)
	//return &thirdQueryData, nil
	// 组装参数 初始化请求的结构体
	zpRsq, errRes := client.QueryOrder(ctx, bm)
	// 请求 curl

	// 获取响应结果 初始化响应的结构体
	if errRes != nil {
		return &thirdQueryData, appError.CodeSupplierHttpErrorCode
	}

	if zpRsq.Code != 200 {
		return &thirdQueryData, appError.CodeSupplierHttpCode
	}
	if zpRsq.Response.Code != DataSuccessCode {
		thirdQueryData.Code = zpRsq.Response.Code
		// 错误的话，错误消息往上层传递
		thirdQueryData.Msg = zpRsq.Response.Message
		return &thirdQueryData, nil
	}
	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if cast.ToUint(zpRsq.Response.QueryOrderData.Amount) != orderInfo.TotalFee {
		logger.ApiWarn(client.LogFileName, client.RequestId, "zpRsq.Response.QueryOrderData.Amount != payoutInfo.TotalFee ", zap.Any("orderInfo", orderInfo))
		thirdQueryData.Code = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		// 错误的话，错误消息往上层传递
		thirdQueryData.Msg = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
		return &thirdQueryData, appError.CodeSupplierInternalChannelParamsFailedErrCode // 道内部参数错误，失败
	}
	upstreamStatus := GetZPayPaymentStatus(zpRsq.Response.QueryOrderData.Status) // 字符串 success

	params := map[string]string{}
	params["finish_time"] = cast.ToString(orderInfo.FinishTime)
	params["status"] = orderInfo.TradeState
	// 中间有 代收中的状态 状态更改是不可逆序的 需要注意
	if orderInfo.TradeState != upstreamStatus {
		upstreamCode := zpRsq.Response.QueryOrderData.Status        // 上游的状态 数字
		orderCode := GetZPayPaymentTradeState(orderInfo.TradeState) // 订单的状态 转成数字
		if orderCode < upstreamCode {

			params["transaction_id"] = zpRsq.Response.QueryOrderData.OrderNo
			params["cash_fee"] = cast.ToString(orderInfo.TotalFee)
			params["payment_link"] = orderInfo.PaymentLink
			params["finish_time"] = cast.ToString(orderInfo.FinishTime)
			params["status"] = upstreamStatus
			// 如果是成功 必须提现状态是 顺序的增长的，不能状态逆序
			// 如果上游返回成功 提现状态是 已申请 则修改
			if upstreamStatus == constant.ORDER_TRADE_STATE_SUCCESS {
				// params["finish_time"] = goutils.Int642String(goutils.GetDateTimeUnix())
				finish_time := goutils.GetTimesTampToUnix(zpRsq.Response.QueryOrderData.SuccessTime)
				params["finish_time"] = cast.ToString(finish_time)
			}
		} else {
			// 如果发生了这种情况 则 记录错误日志
			logger.ApiError(client.LogFileName, client.RequestId, "payoutInfo.TradeState != upstreamStatus error ", zap.Any("payoutInfo", orderInfo))
		}
	}

	// 赋值给统一返回结构体
	zpRsq.Generate(&thirdQueryData)
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

	client, err := initClient(channelDepartInfo, requestId)
	if err != nil {
		//thirdPayoutQueryData.Status = constant.SupplierInitClientCode
		//thirdPayoutQueryData.Msg = ClientInitErr
		return &thirdPayoutQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	//thirdPayoutQueryData.TransactionID = payoutInfo.TransactionId
	//thirdPayoutQueryData.CashFee = payoutInfo.CashFee
	//thirdPayoutQueryData.CashFeeType = payoutInfo.CashFeeType
	//thirdPayoutQueryData.FinishTime = cast.ToInt64(payoutInfo.FinishTime)
	//thirdPayoutQueryData.TradeState = GetZPayPayoutStatus(2)
	//
	//return &thirdPayoutQueryData, nil

	// fmt.Println("thirdPayoutQueryData----------------", thirdPayoutQueryData)
	// fmt.Println("client--------------", client)
	bm := make(model.BodyMap)
	bm.Set("partnerId", client.PartnerId).
		Set("partnerWithdrawNo", payoutInfo.Sn).
		Set("version", "1.0")
	// 组装参数 初始化请求的结构体
	zpRsq, errRes := client.QueryPayout(ctx, bm)
	// 请求 curl

	// 获取响应结果 初始化响应的结构体
	if errRes != nil {
		return &thirdPayoutQueryData, appError.CodeSupplierHttpErrorCode
	}
	if zpRsq.Code != 200 {
		return &thirdPayoutQueryData, appError.CodeSupplierHttpCode
	}
	// fmt.Println("zpRsq----------------", zpRsq)
	if zpRsq.Response.Code != DataSuccessCode {
		thirdPayoutQueryData.Code = zpRsq.Response.Code
		// 错误的话，错误消息往上层传递
		thirdPayoutQueryData.Msg = zpRsq.Msg
		return &thirdPayoutQueryData, nil

	}
	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if cast.ToUint(zpRsq.Response.QueryPayoutData.Amount) != payoutInfo.TotalFee {
		logger.ApiWarn(client.LogFileName, client.RequestId, "zpRsq.Response.QueryOrderData.Amount != payoutInfo.TotalFee ", zap.Any("payoutInfo", payoutInfo))
		thirdPayoutQueryData.Code = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		// 错误的话，错误消息往上层传递
		thirdPayoutQueryData.Msg = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
		return &thirdPayoutQueryData, appError.CodeSupplierInternalChannelParamsFailedErrCode // 道内部参数错误，失败
	}

	QueryPayoutData := zpRsq.Response.QueryPayoutData

	upstreamStatus := GetZPayPayoutStatus(QueryPayoutData.Status)
	params := map[string]string{}
	params["finish_time"] = cast.ToString(payoutInfo.FinishTime)
	params["transaction_id"] = payoutInfo.TransactionId
	if payoutInfo.TradeState != upstreamStatus {
		upstreamCode := QueryPayoutData.Status                      // 上游的状态 数字
		orderCode := GetZPayPayoutTradeState(payoutInfo.TradeState) // 订单的状态 转成数字
		if orderCode < upstreamCode {
			params["transaction_id"] = QueryPayoutData.WithdrawNo
			params["status"] = upstreamStatus // 默认值
			params["finish_time"] = cast.ToString(payoutInfo.FinishTime)
			// 如果是成功 必须提现状态是 顺序的增长的，不能状态逆序
			// 如果上游返回成功 提现状态是 已申请 则修改
			if upstreamStatus == constant.PAYOUT_TRADE_STATE_SUCCESS {
				finish_time := goutils.GetTimesTampToUnix(QueryPayoutData.SuccessTime)
				params["finish_time"] = goutils.Int642String(finish_time)
			}
		} else {
			// 如果发生了这种情况 则 记录错误日志
			logger.ApiError(client.LogFileName, client.RequestId, "payoutInfo.TradeState != upstreamStatus error ", zap.Any("payoutInfo", payoutInfo))
		}
	} else if payoutInfo.TransactionId == "" {
		params["transaction_id"] = QueryPayoutData.WithdrawNo
		params["status"] = upstreamStatus // 默认值
		params["finish_time"] = cast.ToString(payoutInfo.FinishTime)
	}
	zpRsq.Generate(&thirdPayoutQueryData)
	// 赋值给统一返回结构体
	thirdPayoutQueryData.TransactionID = params["transaction_id"]
	thirdPayoutQueryData.CashFee = QueryPayoutData.Amount
	thirdPayoutQueryData.CashFeeType = payoutInfo.FeeType
	thirdPayoutQueryData.FinishTime = cast.ToInt64(params["finish_time"])
	thirdPayoutQueryData.TradeState = upstreamStatus

	// fmt.Println("thirdPayoutQueryData----------------", thirdPayoutQueryData)
	return &thirdPayoutQueryData, nil
}

// GetDepartAccountInfo 获取商户账户信息
func (p *PayImpl) GetDepartAccountInfo(requestId string, channelDepartInfo *model.AspChannelDepartConfig) (*interfaces.ThirdMerchantAccountQueryData, *appError.Error) {
	var thirdMerchantAccountQueryData interfaces.ThirdMerchantAccountQueryData
	thirdMerchantAccountQueryData.Code = DataSuccessCode

	client, err := initClient(channelDepartInfo, requestId)
	if err != nil {
		return &thirdMerchantAccountQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	// fmt.Println("initClient--------------------", client)
	// return thirdMerchantAccountQueryData
	bm := make(model.BodyMap)
	bm.Set("partnerId", client.PartnerId).
		Set("type", "default").
		Set("version", "1.0")
	// 组装参数 初始化请求的结构体
	zpRsq, errRes := client.QueryMerchantAccount(ctx, bm)
	if errRes != nil {
		logger.ApiWarn(client.LogFileName, client.RequestId, "client.QueryMerchantAccount ", zap.String("request_id", client.RequestId), zap.Error(err))
		return &thirdMerchantAccountQueryData, appError.CodeSupplierHttpErrorCode
	}
	thirdMerchantAccountQueryData.Msg = zpRsq.Msg
	if zpRsq.Code != 200 {
		return &thirdMerchantAccountQueryData, appError.CodeSupplierHttpCode
	}
	if zpRsq.Response.Code != DataSuccessCode {
		thirdMerchantAccountQueryData.Code = zpRsq.Response.Code
		thirdMerchantAccountQueryData.Msg = zpRsq.Response.Message
		return &thirdMerchantAccountQueryData, nil
	}

	zpRsq.Generate(&thirdMerchantAccountQueryData)

	return &thirdMerchantAccountQueryData, nil

}

// AddBeneficiary 添加受益人 需要的client 中 channel_depart 中的 config 字段
func (p *PayImpl) AddBeneficiary(requestId string, clientIp string, channelDepartInfo *model.AspChannelDepartConfig, reqBeneficiary *req.AspBeneficiary) (*interfaces.ThirdAddBeneficiary, *appError.Error) {
	var thirdData interfaces.ThirdAddBeneficiary
	thirdData.Code = DataSuccessCode

	return &thirdData, nil
}

func (p *PayImpl) PayoutUpi(requestId string, channelDepartInfo *model.AspChannelDepartConfig, payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutCreateData, *appError.Error) {
	return nil, appError.NotImplementedErrCode
}

func (p *PayImpl) UpiValidate(requestId string, channelDepartInfo *model.AspChannelDepartConfig, upiValidateInfo *req.AspPayoutUpiValidate) (*interfaces.ThirdUpiValidate, *appError.Error) {
	return nil, appError.NotImplementedErrCode
}
