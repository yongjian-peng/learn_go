package fynzonpay

import (
	"asp-payment/api-server/req"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/service/supplier/interfaces"
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

	bm := make(model.BodyMap)
	client, err := NewClient(channelDepartInfo, requestId, constant.FynzonpayLogFileName)
	if err != nil {
		return nil, appError.CodeSupplierInitClientCode
	}

	productName := goutils.IfString(orderInfo.Body == "", orderInfo.Body, ProductName)

	//return scanData
	amount := goutils.Fen2Yuan(cast.ToInt64(orderInfo.TotalFee))
	bm.Set("api_token", client.SecretKey).
		Set("store_id", client.AppId).
		Set("cardsend", CardSend).
		Set("client_ip", orderInfo.ClientIp).
		Set("action", Action).
		Set("source", "Host-Redirect-Card-Payment (Core PHP)").
		Set("source_url", "https://apiprod.sunnypay.net").
		Set("price", amount).
		Set("curr", Curr).
		Set("product_name", productName).
		Set("fullname", orderInfo.CustomerName).
		Set("email", orderInfo.CustomerEmail).
		Set("bill_street_1", "A97B North Block").
		Set("bill_street_2", "West Vinod Nagar").
		Set("bill_city", "New Delhi").
		Set("bill_state", "DL").
		Set("bill_country", "IND").
		Set("bill_zip", "110092").
		Set("bill_phone", orderInfo.CustomerPhone).
		Set("id_order", orderInfo.Sn).
		Set("notify_url", config.AppConfig.Urls.FynzonpayWappayNotifyUrl).
		Set("success_url", config.AppConfig.Urls.FynzonpayOrderReturnUrl).
		Set("error_url", config.AppConfig.Urls.FynzonpayWappayErrorUrl)
	//Set("store_id", client.AppId).
	//Set("cardsend", CardSend).
	//Set("client_ip", orderInfo.ClientIp).
	//Set("action", Action).
	//Set("source", "Curl-Direct-Card-Payment").
	//Set("source_url", "https://apiprod.sunnypay.net").
	//Set("price", amount).
	//Set("curr", Curr).
	//Set("product_name", ProductName).
	//Set("fullname", orderInfo.CustomerName).
	//Set("email", orderInfo.CustomerEmail).
	//Set("bill_street_1", "A97B North Block").
	//Set("bill_street_2", "West Vinod Nagar").
	//Set("bill_city", "New Delhi").
	//Set("bill_state", "DL").
	//Set("bill_country", "IND").
	//Set("bill_zip", "110092").
	//Set("bill_phone", orderInfo.CustomerPhone).
	//Set("id_order", orderInfo.Sn).
	//Set("notify_url", config.AppConfig.Urls.FynzonpayWappayNotifyUrl).
	//Set("success_url", config.AppConfig.Urls.FynzonpayOrderReturnUrl).
	//Set("error_url", config.AppConfig.Urls.FynzonpayWappayErrorUrl).
	////Set("ccno", "4242424242424242").
	////Set("ccvv", "123").
	////Set("month", "01").
	////Set("year", "30").
	//Set("notes", "Test Note-Remark for this transaction")

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
	//fmt.Println("payRsq.CreateOrderBody: ", payRsq.CreateOrderBody)
	if payRsq.CreateOrderBody.Error != "" {
		scanData.Code = payRsq.CreateOrderBody.Error
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
}

// WAPPAY First Pay WAPPAY 支付方式
// 入参 订单信息 渠道信息 项目信息
// 返回 支付的链接
func (p *PayImpl) WAPPAY(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (model.BodyMap, *appError.Error) {

	bm := make(model.BodyMap)
	client, err := NewClient(channelDepartInfo, requestId, constant.FynzonpayLogFileName)
	if err != nil {
		return bm, appError.CodeSupplierInitClientCode
	}

	productName := goutils.IfString(orderInfo.Body == "", orderInfo.Body, ProductName)

	//return scanData
	amount := goutils.Fen2Yuan(cast.ToInt64(orderInfo.TotalFee))
	bm.Set("api_token", client.SecretKey).
		Set("store_id", client.AppId).
		Set("cardsend", CardSend2).
		Set("client_ip", orderInfo.ClientIp).
		Set("action", Action).
		Set("source", "Host-Redirect-Card-Payment (Core PHP)").
		Set("source_url", "https://apiprod.sunnypay.net").
		Set("price", amount).
		Set("curr", Curr).
		Set("product_name", productName).
		Set("fullname", orderInfo.CustomerName).
		Set("email", orderInfo.CustomerEmail).
		Set("bill_street_1", "A97B North Block").
		Set("bill_street_2", "West Vinod Nagar").
		Set("bill_city", "New Delhi").
		Set("bill_state", "DL").
		Set("bill_country", "IND").
		Set("bill_zip", "110092").
		Set("bill_phone", orderInfo.CustomerPhone).
		Set("id_order", orderInfo.Sn).
		Set("notify_url", config.AppConfig.Urls.FynzonpayWappayNotifyUrl).
		Set("success_url", config.AppConfig.Urls.FynzonpayOrderReturnUrl).
		Set("error_url", config.AppConfig.Urls.FynzonpayWappayErrorUrl)

	// api_token=MTEyOTNfMTUxM18yMDIzMDExMjE5NTkzMQ&cardsend=CHECKOUT&client_ip=192.168.147.1&action=product&source=Encode-Checkout&source_url=https%3A%2F%2Fbrightsgaming.com&price=30.00&curr=INR&product_name=Testing+Product&fullname=DEV+PRAKASH+YADAV&email=ericluzhonghua%40gmail.com&bill_street_1=A97B+North+Block&bill_street_2=West+Vinod+Nagar&bill_city=New+Delhi&bill_state=DL&bill_country=IN&bill_zip=110092&bill_phone=9036830689&id_order=202202211811¬ify_url=https%3A%2F%2Fyourdomain.com%2Fnotify.php&success_url=https%3A%2F%2Fyourdomain.com%2Fsuccess.php&error_url=https%3A%2F%2Fyourdomain.com%2Ffailed.php&checkout_url=https%3A%2F%2Fyourdomain.com%2Fcheckout_url.php&ccno=5555555555554444&ccvv=123&month=01&year=30¬es=Remark+for+transaction"

	// ZE1XTEpCeW5lWWlHU09lUFB0ZnpUN2dDL1FhRnRKb2lIVHZia3Q0cDIxdG9wQnQ2SnFkMEl3eVpqMU05amx4dGdNcUhTaUYxNkZQckQzOTdYNFZxb0lzcTN1SVNEdkkxaHlMNmlybjArUUVMU0o5VjhOL0tRejFMa3lEcVFsTU1BWVoreXY0NERLUzBLM1MzR21oWTFmSXE0RW1QTW90ZDJLemtLTHRKSkR4blpUb1V5UXhpeFBjeXFZQ1JFcTlXR1FtMXRYSTFqQituK1daTXo2VjhuWUZEeVNDU01RR2p4VFl6aDlUcC9vNHdHYzdDS29zTVZTOVI4QW9rZFRPWWFETmxvVTNyb1hYNEZmZm9uN3ljTWJOdjR0NG5KWm84UXRTWVFWeDZ2RmZkdEtpK0cyZWw1N29oOWpFMWtWNGxGeW9samdjZ3ZkYkQxc1lUZWIvZ2p1NTVuQWNDZGV3bUIrdjM0R2JoQUhMWEZIYzVOZVR2SC9oM0Frblo5M2E1UkVTTkZWSHRaS0ZqZmtBOXJEMmpWV2g5eVIwVlVLT0xaazNFVDJia2JLaTZHR3NnRDE5eUNTTVJ3NVhlcm5IV1BJUXhyWW9FUGxXblJxRFhmQmppZ3BvRFV5ZDRhMGxoZlhnMFRURFZUaFBGZ1M5Wmk4ak1YSHJsMjQzWjhHK3M0ZXpNeWJwd2lVbUtRRXF6UEQ1Y3JZNWZobCtkQnJRbVkzdytLakxGY3BQZFFnWFJXd0VNZmY1ZlAyMjhUNkZIU0NldEJ4WWMvcHNnbHh4TU4vNm1nZG8rSklkT2xzZFpxNzZwQjQ5N0cyc1FDcFd0c0xKWHJHdW8yaUQ5cU14ZWhXYmhCSE1ReWdRdnRJNUhKYzZNLzBiUXM5N2RDejF6SVU1T0ZqbzJzWW94S3Q3SnJydjRSZE1EM1ZHQ00zR1ZGeHVDZEVSTWphQnJkaFpGNmp0TmVBVFRQU0Vncll1RmZUTGJNZzFVYTB3UDJock51MDIwM1o3R2Y1Q1Y5ZXlEajRNclRSaUpYU2xUZUd1aGF2UE9BZ3ZKNzVIMmtFcGdEcnlTWU5kTEl1VzREVmh0ZHhkZHZxTkNFL3p3YWNqQlB0amNtR0RHUy9FME9Ma1dCa0pKWmtzc3c2ZC9teFUzQU9vNEY4N3lVYVVkL3ZOb08xZXN1RXg3ZmVTTTBTTVFjQWNIYkdWQmpKYUVvdHJ6YklnNHBCUXBFTVBWREw2TVlvS0dHME9IZklQdW92WlcyZW1tUVA4UzZwM0tUM1hsbkl6ejRXckEvRi83bTV5V3Z2bVZJSDdjdGpkdjUxY29HTTV0UjFIZG5VUT0"
	param := bm.EncodeURLParams()
	url := OrderCreateUrl + "?" + param
	bm.Set("url", url)
	return bm, nil
}

// Payout 78Pay web 支付方式
// 入参 订单信息 渠道信息 项目信息
// 返回 支付的链接
func (p *PayImpl) Payout(requestId string, channelDepartInfo *model.AspChannelDepartConfig, payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutCreateData, *appError.Error) {
	var payoutCreateData interfaces.ThirdPayoutCreateData
	payoutCreateData.Code = DataSuccessCode
	client, err := NewClient(channelDepartInfo, requestId, constant.FynzonpayLogFileName)
	if err != nil {
		return &payoutCreateData, appError.CodeSupplierInternalChannelErrCode
	}

	//payoutCreateData.Msg = ""
	//payoutCreateData.TransactionID = "payRsq.CreatePayoutBody.TransactionID"
	//payoutCreateData.CashFee = cast.ToInt(payoutInfo.TotalFee)
	//payoutCreateData.CashFeeType = payoutInfo.CashFeeType
	//payoutCreateData.FinishTime = int64(payoutInfo.FinishTime)
	//payoutCreateData.TradeState = GetPayoutStatus("1")
	//payoutCreateData.CashFeeType = payoutInfo.CashFeeType
	//payoutCreateData.BankType = payoutInfo.BankType
	//return &payoutCreateData, nil

	amount := goutils.Fen2Yuan(cast.ToInt64(payoutInfo.TotalFee))
	bm := make(model.BodyMap)
	bm.Set("payout_token", client.ApplicationId).
		//Set("secret_key", payClient.PartnerId).
		Set("payout_secret_key", client.PayoutSignature).
		Set("checkout", "CURL").
		Set("client_ip", payoutInfo.ClientIp).
		Set("source", "Encode-Curl-API").
		//Set("source_url", "https://needsixgaming.com").
		Set("price", amount).
		Set("curr", Curr).
		Set("remarks", payoutInfo.Note).
		Set("request_id", payoutInfo.Sn).
		Set("product_name", payoutInfo.MchProjectId).
		Set("beneficiary_id", payoutInfo.BeneficiaryId).
		Set("notify_url", client.PayoutNotify)
	//Set("success_url", 1).
	//Set("error_url", payoutInfo.BankCard)
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

	if payRsq.CreatePayoutBody.Status != PayoutSuccessCode {
		payoutCreateData.Code = payRsq.CreatePayoutBody.Status
		payoutCreateData.Msg = payRsq.CreatePayoutBody.Reason
		payoutCreateData.TradeState = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED
		return &payoutCreateData, appError.CodeSupplierChannelErrCode
	}

	payoutCreateData.Code = payRsq.CreatePayoutBody.Status
	payoutCreateData.Msg = payRsq.StatusMsg
	payoutCreateData.TransactionID = payRsq.CreatePayoutBody.TransactionID
	payoutCreateData.CashFee = cast.ToInt(payoutInfo.TotalFee)
	payoutCreateData.CashFeeType = payoutInfo.CashFeeType
	payoutCreateData.FinishTime = int64(payoutInfo.FinishTime)
	payoutCreateData.TradeState = GetPayoutStatus(cast.ToString(payRsq.CreatePayoutBody.StatusNm))
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

	payClient, err := NewClient(channelDepartInfo, requestId, constant.FynzonpayLogFileName)
	if err != nil {
		return &thirdQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	// 如果没有传入 TransactionId 则返回最近的一条数据 会有错误的更新
	if orderInfo.TransactionId == "" {
		return &thirdQueryData, appError.CodeInvalidParamErrCode
	}

	bm := make(model.BodyMap)
	bm.Set("transaction_id", orderInfo.TransactionId).Set("api_token", payClient.SecretKey)
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

	//fmt.Println("payRsp.QueryOrderBody:", fmt.Sprintf("%+v", payRsp.QueryOrderBody))

	//有错误，设置错误消息
	if payRsp.QueryOrderBody.Error != "" {
		thirdQueryData.Code = Error
		// 错误的话，错误消息往上层传递
		thirdQueryData.Msg = payRsp.QueryOrderBody.Error
		return &thirdQueryData, nil
	}

	//根据上游支付状态，转换为本系统支付状态
	upstreamStatus := GetPaymentTradeState(payRsp.QueryOrderBody.StatusNm) // 字符串 success

	//fmt.Println("upstreamStatus: ", upstreamStatus)
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
		//params["finish_time"] = carbon.ParseByFormat(fmt.Sprintf("%s %s", payRsp.QueryOrderBody.ResponseDate, payRsp.QueryOrderBody.ResponseTime), "dmY H:i:s").Format("Y-m-d H:i:s")
		params["status"] = upstreamStatus
	}

	// 赋值给统一返回结构体
	thirdQueryData.TransactionID = payRsp.QueryOrderBody.TransactionID
	thirdQueryData.Msg = payRsp.QueryOrderBody.Descriptor
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

	payClient, err := NewClient(channelDepartInfo, requestId, constant.FynzonpayLogFileName)
	if err != nil {
		return &thirdPayoutQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	bm := make(model.BodyMap)
	bm.Set("payout_token", payClient.ApplicationId).
		//Set("secret_key", payClient.PartnerId).
		Set("payout_secret_key", payClient.PayoutSignature).
		//Set("checkout", "CURL").
		Set("client_ip", payoutInfo.ClientIp).
		Set("action", "paymentdetail").
		Set("source", "Encode-Curl-API").
		//Set("source_url", "https://apiprod.sunnypay.net").
		//Set("transaction_id", payoutInfo.TransactionId).
		Set("order_number", payoutInfo.Sn)
	//Set("transaction_id", "1130317845").
	//Set("order_number", "3ff87684-0a42-422d-a512-7156200fb675")
	//Set("order_number", "75bdb60c-ed9b-4854-8adf-3c3796aa4b90")
	// 组装参数 初始化请求的结构体
	//{"status":"success","msg":"apply successful","transaction_id":"L1208711582746222","balance":"97.00"}
	payRsp, qErr := payClient.QueryPayout(bm)
	// 获取响应结果 初始化响应的结构体
	if qErr != nil {
		return &thirdPayoutQueryData, appError.CodeSupplierHttpErrorCode
	}

	if payRsp.StatusCode != 200 {
		return &thirdPayoutQueryData, appError.CodeSupplierHttpCode
	}

	//有错误，设置错误消息
	if payRsp.QueryPayoutBody.Status != PayoutSuccessCode {
		thirdPayoutQueryData.Code = payRsp.QueryPayoutBody.Status
		// 错误的话，错误消息往上层传递
		thirdPayoutQueryData.Msg = payRsp.QueryPayoutBody.Remarks
		return &thirdPayoutQueryData, nil
	}
	queryPayoutData := payRsp.QueryPayoutBody
	transactionAmount := cast.ToFloat64(payRsp.QueryPayoutBody.TransactionAmount)
	if transactionAmount < 0 {
		transactionAmount = -transactionAmount
	}
	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if cast.ToUint(goutils.Yuan2Fen(transactionAmount)) != payoutInfo.TotalFee {
		logger.ApiWarn(payClient.LogFileName, payClient.RequestId, "payRsp.QueryPayoutBody.Amount != payoutInfo.TotalFee ", zap.Any("payoutInfo", payoutInfo))
		thirdPayoutQueryData.Code = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		// 错误的话，错误消息往上层传递
		thirdPayoutQueryData.Msg = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
		return &thirdPayoutQueryData, appError.CodeSupplierInternalChannelUpstreamErrCode // 道内部参数错误，失败
	}

	upstreamStatus := GetPayoutStatus(cast.ToString(queryPayoutData.TransactionStatus))
	params := map[string]string{}

	params["finish_time"] = cast.ToString(payoutInfo.FinishTime)
	params["transaction_id"] = payoutInfo.TransactionId
	params["status"] = payoutInfo.TradeState // 默认为数据库的状态

	//上游是非成功状态，并且数据库的状态不是成功状态，更新数据库
	//fmt.Println("upstreamStatus:", upstreamStatus)
	//fmt.Println("payoutInfo.TradeState:", payoutInfo.TradeState)
	if upstreamStatus != constant.PAYOUT_TRADE_STATE_SUCCESS && payoutInfo.TradeState != constant.PAYOUT_TRADE_STATE_SUCCESS {
		params["status"] = upstreamStatus // 设置最新的状态，成功状态不可逆转
	}

	//如果上游是成功状态，更新完成时间
	if upstreamStatus == constant.PAYOUT_TRADE_STATE_SUCCESS {
		params["finish_time"] = cast.ToString(goutils.GetDateTimeUnix())
		//params["transaction_id"] = queryPayoutData.TransactionId
		params["status"] = upstreamStatus // 设置最新的状态，成功状态不可逆转
	}

	// 赋值给统一返回结构体
	thirdPayoutQueryData.TransactionID = payoutInfo.TransactionId
	thirdPayoutQueryData.CashFee = cast.ToInt(payoutInfo.TotalFee)
	thirdPayoutQueryData.CashFeeType = payoutInfo.FeeType
	thirdPayoutQueryData.FinishTime = cast.ToInt64(params["finish_time"])
	thirdPayoutQueryData.TradeState = params["status"]
	thirdPayoutQueryData.Msg = queryPayoutData.Remarks

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

	var thirdAddBeneficiaryData interfaces.ThirdAddBeneficiary
	thirdAddBeneficiaryData.Code = "0000"
	//thirdAddBeneficiaryData.BenefiaryId = "113031010341"
	//return &thirdAddBeneficiaryData, nil

	payClient, err := NewClient(channelDepartInfo, requestId, constant.FynzonpayLogFileName)
	if err != nil {
		return &thirdAddBeneficiaryData, appError.CodeSupplierInternalChannelErrCode
	}

	//thirdAddBeneficiaryData.Code = "0000"
	//thirdAddBeneficiaryData.BenefiaryId = "113031010915"
	//return &thirdAddBeneficiaryData, nil
	bm := make(model.BodyMap)
	bm.Set("payout_token", payClient.ApplicationId).
		//Set("secret_key", payClient.PartnerId).
		Set("payout_secret_key", payClient.PayoutSignature).
		Set("client_ip", clientIp).
		Set("source", "Encode-Curl-API").
		Set("source_url", "https://needsixgaming.com").
		Set("beneficiary_nickname", reqBeneficiary.CustomerName).
		Set("beneficiary_name", reqBeneficiary.CustomerName).
		Set("beneficiaryEmailId", reqBeneficiary.CustomerEmail).
		Set("beneficiaryPhone", reqBeneficiary.CustomerPhone).
		Set("account_number", reqBeneficiary.BankCard).
		Set("beneficiary_ac_repeat", reqBeneficiary.BankCard).
		Set("beneficiary_bank_name", reqBeneficiary.BankCode).
		Set("bank_code1", reqBeneficiary.Ifsc).
		Set("notify_url", payClient.BeneficiaryNotify)
	// 组装参数 初始化请求的结构体
	//{\"status\":\"0000\",\"bene_id\":\"113031010341\",\"reason\":\"Success\",\"remark\":\"Beneficiary successfully added\",\"notify\":\"OK\"}
	payRsp, qErr := payClient.AddbeneFiciary(bm)
	// 获取响应结果 初始化响应的结构体
	if qErr != nil {
		return &thirdAddBeneficiaryData, appError.CodeSupplierHttpErrorCode
	}

	if payRsp.StatusCode != 200 {
		return &thirdAddBeneficiaryData, appError.CodeSupplierHttpCode
	}

	//有错误，设置错误消息
	if payRsp.CreateBeneficiaryBody.Status != AddBeneficiarySuccessCode {
		thirdAddBeneficiaryData.Code = payRsp.CreateBeneficiaryBody.Status
		// 错误的话，错误消息往上层传递
		thirdAddBeneficiaryData.Msg = payRsp.CreateBeneficiaryBody.Remark
		return &thirdAddBeneficiaryData, nil
	}
	thirdAddBeneficiaryData.BenefiaryId = payRsp.CreateBeneficiaryBody.BeneId
	return &thirdAddBeneficiaryData, nil
}

func (p *PayImpl) PayoutUpi(requestId string, channelDepartInfo *model.AspChannelDepartConfig, payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutCreateData, *appError.Error) {
	return nil, appError.NotImplementedErrCode
}

func (p *PayImpl) UpiValidate(requestId string, channelDepartInfo *model.AspChannelDepartConfig, upiValidateInfo *req.AspPayoutUpiValidate) (*interfaces.ThirdUpiValidate, *appError.Error) {
	return nil, appError.NotImplementedErrCode
}
