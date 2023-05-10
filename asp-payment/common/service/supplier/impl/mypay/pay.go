package mypay

import (
	"asp-payment/api-server/req"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/service/supplier/interfaces"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
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
	client, err := NewClient(channelDepartInfo, requestId, constant.MyPayLogFileName)
	if err != nil {
		return &scanData, appError.CodeSupplierInitClientCode
	}

	//return scanData
	amount := goutils.Fen2Yuan(cast.ToInt64(orderInfo.TotalFee))
	bm := make(model.BodyMap)
	// QrCode
	//bm.Set("token", client.AppId).
	//	Set("amount", amount).
	//	Set("pn", "MyPay").
	//	Set("vpaAdderess", "9818934245@ybl").
	//	Set("userId", 100569).
	//	Set("userTxnId", orderInfo.Sn).
	//	Set("userName", orderInfo.CustomerName)

	// DynamicQrCode
	bm.Set("token", client.AppId).
		Set("pn", "SURYA ENTERPRISES").
		//Set("pn", "MYPAY78868898739"). // MyPay Communication Private Limited
		Set("amount", amount)

	//{ "statuscode": "Success", "status": "Recharge Succesfull", "data": { "status": "Success", "rechargeAmount": 10, "operatorRef": "434347558", "operatorCode": "AT", "userTxnId": "9385fede-7e24-45cc-adf4-b404c7f8fb79", "mobileNo": "9799518927", "myPayTxnId": "b4501e309ca64029b79326d06587721a" } }
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
	if payRsq.CreateOrderBody.Success != CreateOrderSuccessCode {
		scanData.Code = payRsq.CreateOrderBody.Success
		scanData.Msg = payRsq.CreateOrderBody.Message
		return &scanData, appError.CodeSupplierChannelErrCode
	}
	// 赋值给统一返回结构体
	payRsq.GenerateCreateOrder(&scanData)
	scanData.CashFeeType = orderInfo.CashFeeType
	scanData.FinishTime = int64(orderInfo.FinishTime)
	scanData.TradeState = orderInfo.TradeState
	scanData.Rate = cast.ToInt(orderInfo.Rate)
	scanData.Qrcode = orderInfo.Qrcode
	scanData.SettlementTime = orderInfo.SettlementTime
	scanData.PaymentsURL = orderInfo.PaymentsUrl
	scanData.ReturnURL = orderInfo.ReturnUrl
	scanData.SettlementsURL = orderInfo.SettlementsUrl
	scanData.CashFeeType = orderInfo.CashFeeType
	scanData.BankType = orderInfo.BankType
	scanData.CashFee = cast.ToInt(orderInfo.TotalFee)

	return &scanData, nil

	return &scanData, nil
}

// WAPPAY First Pay WAPPAY 支付方式
// 入参 订单信息 渠道信息 项目信息
// 返回 支付的链接
func (p *PayImpl) WAPPAY(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (model.BodyMap, *appError.Error) {

	bm := make(model.BodyMap)
	client, err := NewClient(channelDepartInfo, requestId, constant.MyPayLogFileName)
	if err != nil {
		return bm, appError.CodeSupplierInitClientCode
	}
	//return scanData
	amount := goutils.Fen2Yuan(cast.ToInt64(orderInfo.TotalFee))

	// DynamicQrCode
	bm.Set("token", client.AppId).
		Set("pn", "SURYA ENTERPRISES").
		//Set("pn", "MYPAY78868898739"). // MyPay Communication Private Limited
		Set("amount", amount)

	//bm.Set("url", "upi://pay?pa=604751@icici&pn=SURYA ENTERPRISES&tr=EZV2023032112052390684001&am=100.0&cu=INR&mc=5411").
	//	Set("transactionId", "7085829230OHGLJ")
	//return bm, nil

	//{ "statuscode": "Success", "status": "Recharge Succesfull", "data": { "status": "Success", "rechargeAmount": 10, "operatorRef": "434347558", "operatorCode": "AT", "userTxnId": "9385fede-7e24-45cc-adf4-b404c7f8fb79", "mobileNo": "9799518927", "myPayTxnId": "b4501e309ca64029b79326d06587721a" } }
	// 组装参数 初始化请求的结构体
	payRsq, cErr := client.CreateOrder(bm)
	// 请求 curl

	// 获取响应结果 初始化响应的结构体
	if cErr != nil {
		return bm, appError.CodeSupplierHttpErrorCode
	}

	if payRsq.StatusCode != 200 {
		return bm, appError.CodeSupplierHttpCode
	}
	if payRsq.CreateOrderBody.Success != CreateOrderSuccessCode {
		//scanData.Code = payRsq.CreateOrderBody.Success
		//scanData.Msg = payRsq.CreateOrderBody.Message
		return bm, appError.CodeSupplierChannelErrCode
	}
	bm.Set("url", payRsq.CreateOrderBody.Url).
		Set("transactionId", payRsq.CreateOrderBody.MyPayTransactionId)
	return bm, nil
}

// Payout 78Pay web 支付方式
// 入参 订单信息 渠道信息 项目信息
// 返回 支付的链接
func (p *PayImpl) Payout(requestId string, channelDepartInfo *model.AspChannelDepartConfig, payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutCreateData, *appError.Error) {
	var payoutCreateData interfaces.ThirdPayoutCreateData
	payoutCreateData.Code = DataSuccessCode
	client, err := NewClient(channelDepartInfo, requestId, constant.MyPayLogFileName)
	if err != nil {
		return &payoutCreateData, appError.CodeSupplierInternalChannelErrCode
	}
	amount := goutils.Fen2Yuan(cast.ToInt64(payoutInfo.TotalFee))
	bm := make(model.BodyMap)
	bm.Set("token", client.AppId).
		Set("UserName", payoutInfo.CustomerName).
		Set("AccountNumber", payoutInfo.BankCard).
		Set("IfscCode", payoutInfo.Ifsc).
		Set("paymentMode", 1).
		Set("beneficiaryName", payoutInfo.CustomerName).
		Set("Amount", amount). //
		Set("MobileNo", payoutInfo.CustomerPhone).
		Set("Email", payoutInfo.CustomerEmail).
		Set("Remark", "").
		Set("TransactionId", payoutInfo.Sn)
	// 生成签名
	// fmt.Println("bm--------------", bm)
	// 组装参数 初始化请求的结构体
	// 请求 curl
	//{"statuscode":"TXN","status":"Transaction Successful","Data":{ "userId": 0, "userName": null, "token": "125562XXXXXXXXX45b5ea80", "ipAdderess": "61.2.243.152", "accountNumber": "427XXXXXX294", "ifscCode": "BAXXXXXXRO", "beneficiaryName": "NeXXXXXin", "amount": 10, "mobileNo": "97XXXXXXX27", "email": "user@example.com", "operatorId": "121XXXXXFLQQ", "transactionId": "123456789", "myPayTransactionId": "173XXXXXXXXQI","Status":"Success","StatusCode":"TXN","StatusMessage":"Transaction Successful","Remark":"B370D32F1F","PayoutId":81260,"BankRef":"127813690308","CreateDate":"2021-10-05T13:08:46.347"}}
	payRsq, pErr := client.CreatePayout(bm)
	// 获取响应结果 初始化响应的结构体
	if pErr != nil {
		return &payoutCreateData, appError.CodeSupplierHttpErrorCode
	}
	// fmt.Println("bm32132--------------", bm)
	if payRsq.StatusCode != 200 {
		return &payoutCreateData, appError.CodeSupplierHttpCode
	}

	if payRsq.CreatePayoutBody.StatusCode == CreatePayoutDataFailedCode {
		payoutCreateData.Code = payRsq.CreatePayoutBody.StatusCode
		payoutCreateData.Msg = payRsq.CreatePayoutBody.Status
		payoutCreateData.TradeState = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED
		return &payoutCreateData, appError.CodeSupplierChannelErrCode
	}

	payoutCreateData.Msg = payRsq.StatusMsg
	payoutCreateData.TransactionID = payRsq.CreatePayoutBody.CreatePayoutData.MyPayTransactionId
	payoutCreateData.CashFee = cast.ToInt(payoutInfo.TotalFee)
	payoutCreateData.CashFeeType = payoutInfo.CashFeeType
	payoutCreateData.FinishTime = int64(payoutInfo.FinishTime)
	payoutCreateData.TradeState = GetPayoutStatus(payRsq.CreatePayoutBody.CreatePayoutData.StatusCode)
	payoutCreateData.CashFeeType = payoutInfo.CashFeeType
	payoutCreateData.BankType = payoutInfo.BankType
	payoutCreateData.BankUtr = payRsq.CreatePayoutBody.CreatePayoutData.BankRef
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

	payClient, err := NewClient(channelDepartInfo, requestId, constant.MyPayLogFileName)
	if err != nil {
		return &thirdQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	amount := goutils.Fen2Yuan(cast.ToInt64(orderInfo.TotalFee))
	bm := make(model.BodyMap)
	bm.Set("Token", payClient.AppId).
		Set("DealerTxnId", orderInfo.Sn).
		Set("Amount", amount).
		//Set("BankRef", "").
		Set("TxnType", "DQR")
	fmt.Println("APP_ID", payClient.AppId)
	fmt.Println("ORDER_ID", orderInfo.Sn)
	// 组装参数 初始化请求的结构体
	payRsp, qErr := payClient.QueryOrder(bm)
	//fmt.Println("qErr:",qErr)
	// 获取响应结果 初始化响应的结构体
	if qErr != nil {
		return &thirdQueryData, appError.CodeSupplierHttpErrorCode
	}

	if payRsp.StatusCode != 200 {
		return &thirdQueryData, appError.CodeSupplierHttpCode
	}

	fmt.Println("payRsp.QueryOrderBody:", fmt.Sprintf("%+v", payRsp.QueryOrderBody))

	//有错误，设置错误消息
	if payRsp.QueryOrderBody.StatusCode != DataSuccessCode {
		thirdQueryData.Code = payRsp.QueryOrderBody.StatusCode
		// 错误的话，错误消息往上层传递
		thirdQueryData.Msg = payRsp.QueryOrderBody.Status
		return &thirdQueryData, nil
	}

	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if cast.ToUint(payRsp.QueryOrderBody.QueryOrderData.RechargeAmount) != orderInfo.TotalFee {
		logger.ApiWarn(payClient.LogFileName, payClient.RequestId, "payRsp.QueryPayoutBody.Amount != payoutInfo.TotalFee ", zap.Any("orderInfo", orderInfo))
		thirdQueryData.Code = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		// 错误的话，错误消息往上层传递
		thirdQueryData.Msg = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
		return &thirdQueryData, appError.CodeSupplierInternalChannelParamsFailedErrCode // 道内部参数错误，失败
	}

	//根据上游支付状态，转换为本系统支付状态
	upstreamStatus := GetPaymentTradeState(payRsp.QueryOrderBody.QueryOrderData.Status) // 字符串 success
	params := map[string]string{}
	//默认完成时间取数据库的值
	params["finish_time"] = cast.ToString(orderInfo.FinishTime)
	params["status"] = orderInfo.TradeState // 默认为数据库的状态

	//上游是非成功状态，并且数据库的状态不是成功状态，更新数据库
	if upstreamStatus != constant.ORDER_TRADE_STATE_SUCCESS && orderInfo.TradeState != constant.ORDER_TRADE_STATE_SUCCESS {
		params["status"] = upstreamStatus // 设置最新的状态，成功状态不可逆转
	}
	// 赋值给统一返回结构体
	thirdQueryData.TransactionID = payRsp.QueryOrderBody.QueryOrderData.MyPayTxnId
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

	payClient, err := NewClient(channelDepartInfo, requestId, constant.MyPayLogFileName)
	if err != nil {
		return &thirdPayoutQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	amount := goutils.Fen2Yuan(cast.ToInt64(payoutInfo.TotalFee))
	bm := make(model.BodyMap)
	bm.Set("Token", payClient.AppId).
		Set("DealerTxnId", payoutInfo.Sn).
		Set("Amount", amount).
		Set("BankRef", payoutInfo.BankUtr).
		Set("TxnType", "PYT")
	fmt.Println("APP_ID", payClient.AppId)
	// 组装参数 初始化请求的结构体
	// { "statuscode": "Success", "status": "Record Found", "data": { "transactionStatusCode": "TXN", "transactionStatus": "Success", "transactionStatusMessage": "Transaction Successful", "transactionAmount": 10, "transactionBankRef": "20xxxxxx18", "txnOrder": { "dealerTxnId": "3xxxxx0", "MyPayTransactionId": "2201252362SDGMK", "spKey": "IMPS", "account": "91xxxxx2", "optional1": "Pxxxx6", "optional2": "neeraj ", "optional3": "", "optional4": "", "optional5": "" } }, "txnType": "PYT", "timestamp": "202203291759465281", "MyPayGuid": "2beb18d2be1340eb88712d0b28b51d9c" }
	payRsp, qErr := payClient.QueryPayout(bm)
	// 获取响应结果 初始化响应的结构体
	if qErr != nil {
		return &thirdPayoutQueryData, appError.CodeSupplierHttpErrorCode
	}

	if payRsp.StatusCode != 200 {
		return &thirdPayoutQueryData, appError.CodeSupplierHttpCode
	}

	//有错误，设置错误消息
	if payRsp.QueryPayoutBody.Status != DataSuccessCode {
		thirdPayoutQueryData.Code = payRsp.QueryPayoutBody.StatusCode
		// 错误的话，错误消息往上层传递
		thirdPayoutQueryData.Msg = payRsp.QueryPayoutBody.Status
		return &thirdPayoutQueryData, nil
	}
	queryPayoutData := payRsp.QueryPayoutBody

	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if cast.ToUint(payRsp.QueryOrderBody.QueryOrderData.RechargeAmount) != payoutInfo.TotalFee {
		logger.ApiWarn(payClient.LogFileName, payClient.RequestId, "payRsp.QueryPayoutBody.Amount != payoutInfo.TotalFee ", zap.Any("payoutInfo", payoutInfo))
		thirdPayoutQueryData.Code = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		// 错误的话，错误消息往上层传递
		thirdPayoutQueryData.Msg = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
		return &thirdPayoutQueryData, appError.CodeSupplierInternalChannelParamsFailedErrCode // 道内部参数错误，失败
	}

	upstreamStatus := GetPayoutStatus(queryPayoutData.StatusCode)
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
		finishTime := goutils.GetCurTimeUnixSecond()
		params["finish_time"] = goutils.Int642String(finishTime)
		params["status"] = upstreamStatus // 设置最新的状态，成功状态不可逆转
	}

	// 赋值给统一返回结构体
	thirdPayoutQueryData.TransactionID = params["transaction_id"]
	thirdPayoutQueryData.CashFee = cast.ToInt(payoutInfo.TotalFee)
	thirdPayoutQueryData.CashFeeType = payoutInfo.FeeType
	thirdPayoutQueryData.FinishTime = cast.ToInt64(params["finish_time"])
	thirdPayoutQueryData.TradeState = params["status"]
	thirdPayoutQueryData.Msg = payRsp.QueryPayoutBody.Status

	return &thirdPayoutQueryData, nil
}

// GetDepartAccountInfo 获取商户账户信息
func (p *PayImpl) GetDepartAccountInfo(requestId string, channelDepartInfo *model.AspChannelDepartConfig) (*interfaces.ThirdMerchantAccountQueryData, *appError.Error) {
	var thirdMerchantAccountQueryData interfaces.ThirdMerchantAccountQueryData
	thirdMerchantAccountQueryData.Code = DataSuccessCode

	return &thirdMerchantAccountQueryData, nil

}

// AddBeneficiary 添加受益人 需要的client 中 channel_depart 中的 config 字段
func (p *PayImpl) AddBeneficiary(requestId string, clientIp string, channelDepartInfo *model.AspChannelDepartConfig, reqBeneficiary *req.AspBeneficiary) (*interfaces.ThirdAddBeneficiary, *appError.Error) {
	var thirdData interfaces.ThirdAddBeneficiary
	thirdData.Code = DataSuccessCode

	return &thirdData, nil
}

func (p *PayImpl) PayNotify() {
	//TODO implement me
	panic("implement me")
}

func (p *PayImpl) PayoutUpi(requestId string, channelDepartInfo *model.AspChannelDepartConfig, payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutCreateData, *appError.Error) {
	return nil, appError.NotImplementedErrCode
}

func (p *PayImpl) UpiValidate(requestId string, channelDepartInfo *model.AspChannelDepartConfig, upiValidateInfo *req.AspPayoutUpiValidate) (*interfaces.ThirdUpiValidate, *appError.Error) {
	return nil, appError.NotImplementedErrCode
}
