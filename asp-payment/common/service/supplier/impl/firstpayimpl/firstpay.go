package firstpayimpl

import (
	"asp-payment/api-server/req"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
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

func initClient(channelDepartInfo *model.AspChannelDepartConfig, request_id string) (*Client, error) {
	// fmt.Println("depart_id------------------------", depart_id)
	var client *Client
	// 2. 根据订单信息 查询出商户对应的 key 入参是 需要查询到对应的 asp_channel_depart 中对应 config 字段 并拿到上游asp_channel_depart参数
	// config 转 struct 字符串转 struct
	var channelConfigInfo model.AspChannelDepartConfigInfo
	goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo)

	// 初始化 client 需要参数 appid Signature
	// appid := "tpEarth123"
	// SerectKey := "d519a24f7f1c87e084763e77d6d1bb114bc693051e1774aa2f07a424b79db23f"
	// client = NewClient(appid, SerectKey)
	client = NewClient(channelConfigInfo.Appid, channelConfigInfo.Signature)

	client.RequestId = request_id
	client.LogFileName = constant.FirstPayLogFileName
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
	//return scanData
	amount := goutils.Fen2Yuan(cast.ToInt64(orderInfo.TotalFee))
	bm := make(model.BodyMap)
	bm.Set("app_order_id", orderInfo.Sn).
		Set("amount", amount).
		Set("phone", orderInfo.CustomerPhone).
		Set("user_name", orderInfo.CustomerName).
		Set("return_url", orderInfo.ReturnUrl)

	// 生成签名

	// 组装参数 初始化请求的结构体
	fpRsq, oErr := client.CreateOrder(bm)
	// 请求 curl

	// 获取响应结果 初始化响应的结构体
	if oErr != nil {
		return &scanData, appError.CodeSupplierHttpErrorCode
	}

	if fpRsq.Code != 200 {
		return &scanData, appError.CodeSupplierHttpCode
	}
	if fpRsq.Response.Code != DataSuccessCode {
		scanData.Code = fpRsq.Response.Code
		scanData.Msg = fpRsq.Response.Msg
		return &scanData, appError.CodeSupplierChannelErrCode
	}
	// 赋值给统一返回结构体
	fpRsq.Generate(&scanData)
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
		//payoutCreateData.Status = "01"
		//payoutCreateData.Msg = ClientInitErr
		return &payoutCreateData, appError.CodeSupplierInitClientCode
	}
	// return payoutCreateData
	amount := goutils.Fen2Yuan(cast.ToInt64(payoutInfo.TotalFee))
	bm := make(model.BodyMap)
	bm.Set("app_order_id", payoutInfo.Sn).
		Set("amount", amount).
		Set("phone", payoutInfo.CustomerPhone).
		Set("user_name", payoutInfo.CustomerName).
		Set("ifsc", payoutInfo.Ifsc).
		Set("bank_card", payoutInfo.BankCard).
		Set("bank_code", payoutInfo.BankCode).
		Set("vpa", payoutInfo.Vpa).
		Set("pay_type", payoutInfo.PayType)
	// 生成签名
	// fmt.Println("bm--------------", bm)
	// 组装参数 初始化请求的结构体
	fpRsq, pErr := client.CreatePayout(bm)
	// 请求 curl

	// 获取响应结果 初始化响应的结构体
	if pErr != nil {
		//payoutCreateData.Status = "02"
		//payoutCreateData.Msg = DoCurlErr
		return &payoutCreateData, appError.CodeSupplierHttpErrorCode
	}
	// fmt.Println("bm32132--------------", bm)
	if fpRsq.Code != 200 {
		//payoutCreateData.Status = "01"
		// 错误的话，错误消息往上层传递
		//payoutCreateData.Msg = fpRsq.Msg
		return &payoutCreateData, appError.CodeSupplierHttpCode
	}
	if fpRsq.Response.Code != DataSuccessCode {
		payoutCreateData.Code = fpRsq.Response.Code
		payoutCreateData.Msg = fpRsq.Response.Msg
		return &payoutCreateData, appError.CodeSupplierChannelErrCode
	}

	// 更新订单信息 入参结构体 返回error
	params := map[string]string{}
	if fpRsq.Response.CreatePayoutData.Status == StatusPayoutSuccess {
		params["status"] = constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS
		params["finish_time"] = goutils.Int642String(goutils.GetDateTimeUnix())

	} else if fpRsq.Response.CreatePayoutData.Status == StatusPayoutFailed {
		params["status"] = constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED
	} else {
		params["status"] = constant.PAYOUT_TRADE_STATE_PENDING
		params["finish_time"] = "0"
	}

	// 赋值给统一返回结构体
	fpRsq.Generate(&payoutCreateData)
	payoutCreateData.CashFeeType = payoutInfo.CashFeeType
	payoutCreateData.FinishTime = cast.ToInt64(params["finish_time"])
	payoutCreateData.TradeState = params["status"]
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
// 注意：上游提供的是异步查询的方法，发送回调的消息，成功则等待回调通知订单结果
func (p *PayImpl) PayQuery(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (*interfaces.ThirdQueryData, *appError.Error) {

	var thirdQueryData interfaces.ThirdQueryData
	thirdQueryData.Code = DataSuccessCode

	// fmt.Println("hello-------------------")
	client, err := initClient(channelDepartInfo, requestId)
	if err != nil {
		//thirdQueryData.Status = "01"
		//thirdQueryData.Msg = ClientInitErr
		return &thirdQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	// fmt.Println("thirdQueryData------------", thirdQueryData)
	// return thirdQueryData
	bm := make(model.BodyMap)
	// 组装参数 初始化请求的结构体
	fpRsq, qErr := client.QueryOrder(bm, orderInfo.TransactionId)
	// 请求 curl

	// 获取响应结果 初始化响应的结构体
	if qErr != nil {
		//thirdQueryData.Status = "01"
		//thirdQueryData.Msg = DoCurlErr
		return &thirdQueryData, appError.CodeSupplierHttpErrorCode
	}

	if fpRsq.Code != 200 {
		//thirdQueryData.Status = "01"
		// 错误的话，错误消息往上层传递
		//thirdQueryData.Msg = fpRsq.Msg
		return &thirdQueryData, appError.CodeSupplierHttpCode
	}
	if fpRsq.Response.Code != DataSuccessCode {
		thirdQueryData.Code = fpRsq.Response.Code
		thirdQueryData.Msg = fpRsq.Response.Msg
		return &thirdQueryData, nil
	}
	//上游不是同步返回，只能异步通知，原样返回给客户端
	// 赋值给统一返回结构体
	thirdQueryData.TransactionID = orderInfo.TransactionId
	thirdQueryData.CashFee = orderInfo.CashFee
	thirdQueryData.CashFeeType = orderInfo.CashFeeType
	thirdQueryData.FinishTime = int64(orderInfo.FinishTime)
	thirdQueryData.TradeState = orderInfo.TradeState

	return &thirdQueryData, nil
}

// PayoutQuery 查询上游订单状态 需要的client 中 channel_depart 中的 config 字段
// 上游提供的是异步查询的方法，发送回调的消息，成功则等待回调通知订单结果
func (p *PayImpl) PayoutQuery(requestId string, channelDepartInfo *model.AspChannelDepartConfig, payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutQueryData, *appError.Error) {

	var thirdPayoutQueryData interfaces.ThirdPayoutQueryData
	thirdPayoutQueryData.Code = DataSuccessCode
	// fmt.Println("payoutInfo---------------", payoutInfo)

	// fmt.Println("depart_id---------------", payoutInfo)
	// fmt.Println("channel_id---------------", channel_id)
	client, err := initClient(channelDepartInfo, requestId)
	if err != nil {
		//thirdPayoutQueryData.Status = "01"
		//thirdPayoutQueryData.Msg = ClientInitErr
		return &thirdPayoutQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	// return thirdPayoutQueryData
	// fmt.Println("client--------------", client)
	// 组装参数 初始化请求的结构体
	fpRsq, qErr := client.QueryPayout(payoutInfo.TransactionId, payoutInfo.OutTradeNo)
	// 请求 curl

	// 获取响应结果 初始化响应的结构体
	if qErr != nil {
		//thirdPayoutQueryData.Status = "02"
		//thirdPayoutQueryData.Msg = DoCurlErr
		return &thirdPayoutQueryData, appError.CodeSupplierHttpErrorCode
	}
	// fmt.Println("fpRsq----------------", fpRsq)
	if fpRsq.Code != 200 {
		//thirdPayoutQueryData.Status = "01"
		// 错误的话，错误消息往上层传递
		//thirdPayoutQueryData.Msg = fpRsq.Msg
		return &thirdPayoutQueryData, appError.CodeSupplierHttpCode
	}
	if fpRsq.Response.Code != DataSuccessCode {
		thirdPayoutQueryData.Code = fpRsq.Response.Code
		thirdPayoutQueryData.Msg = fpRsq.Response.Msg
		return &thirdPayoutQueryData, nil
	}
	// 提现查询更新逻辑 横列是 上游返回 竖列是提现当前的状态
	//                 已申请              已成功               已失败
	// 已申请          不更新              更新成功             更新失败
	// 已成功          不更新               不更新                待定
	// 已失败          不更新                待定                不更新
	//

	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if goutils.Yuan2Fen(cast.ToFloat64(fpRsq.Response.QueryPayoutData.Amount)) != cast.ToInt64(payoutInfo.TotalFee) {
		logger.ApiWarn(client.LogFileName, client.RequestId, "fpRsq.Response.QueryPayoutData.Amount != payoutInfo.TotalFee ", zap.Any("payoutInfo", payoutInfo))
		thirdPayoutQueryData.Code = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		// 错误的话，错误消息往上层传递
		thirdPayoutQueryData.Msg = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
		return &thirdPayoutQueryData, appError.CodeSupplierInternalChannelParamsFailedErrCode // 道内部参数错误，失败
	}

	// 如果当前的提现状态和上游返回的情况不一致的情况
	upstreamStatus := GetFirstPayPayoutStatus(fpRsq.Response.QueryPayoutData.Status)
	params := map[string]string{}
	params["finish_time"] = cast.ToString(payoutInfo.FinishTime)
	params["status"] = payoutInfo.TradeState // 默认值
	if payoutInfo.TradeState != upstreamStatus {
		if upstreamStatus != constant.PAYOUT_TRADE_STATE_PENDING {
			params["transaction_id"] = fpRsq.Response.QueryPayoutData.OrderId
			// 如果是成功 必须提现状态是 顺序的增长的，不能状态逆序
			// 如果上游返回成功 提现状态是 已申请 则修改
			if upstreamStatus == constant.PAYOUT_TRADE_STATE_SUCCESS && payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_PENDING {
				params["finish_time"] = cast.ToString(goutils.GetDateTimeUnix())
				params["status"] = upstreamStatus
				// 如果上游返回失败 提现状态是 已申请 则修改
			} else if upstreamStatus == constant.PAYOUT_TRADE_STATE_FAILED && payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_PENDING {
				params["status"] = upstreamStatus
			}
		} else {
			// 如果发生了这种情况 则 记录错误日志
			zap.L().Error("payoutInfo.TradeState != upstreamStatus error ", zap.String("request_id", client.RequestId), zap.Any("payoutInfo", payoutInfo))
		}
	}
	fpRsq.Generate(&thirdPayoutQueryData)
	// 赋值给统一返回结构体
	thirdPayoutQueryData.CashFeeType = payoutInfo.FeeType
	thirdPayoutQueryData.FinishTime = cast.ToInt64(params["finish_time"])
	thirdPayoutQueryData.TradeState = params["status"]
	// fmt.Println("thirdPayoutQueryData----------------", thirdPayoutQueryData)
	return &thirdPayoutQueryData, nil
}

// GetDepartAccountInfo 获取商户账户信息
func (p *PayImpl) GetDepartAccountInfo(requestId string, channelDepartInfo *model.AspChannelDepartConfig) (*interfaces.ThirdMerchantAccountQueryData, *appError.Error) {
	var thirdMerchantAccountQueryData interfaces.ThirdMerchantAccountQueryData
	thirdMerchantAccountQueryData.Code = DataSuccessCode

	client, err := initClient(channelDepartInfo, requestId)
	if err != nil {
		//thirdMerchantAccountQueryData.Status = "01"
		return &thirdMerchantAccountQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	// fmt.Println("initClient--------------------", client)
	// return thirdMerchantAccountQueryData
	// 组装参数 初始化请求的结构体
	fpRsq, qErr := client.QueryMerchantAccount()
	if qErr != nil {
		logger.ApiWarn(client.LogFileName, client.RequestId, "client.QueryMerchantAccount ", zap.String("request_id", client.RequestId), zap.Error(qErr))
		//thirdMerchantAccountQueryData.Status = "02"
		return &thirdMerchantAccountQueryData, appError.CodeSupplierHttpErrorCode
	}
	thirdMerchantAccountQueryData.Msg = fpRsq.Msg
	if fpRsq.Code != 200 {
		//thirdMerchantAccountQueryData.Status = "01"
		return &thirdMerchantAccountQueryData, appError.CodeSupplierHttpCode
	}

	if fpRsq.Response.Code != DataSuccessCode {
		thirdMerchantAccountQueryData.Code = fpRsq.Response.Code
		thirdMerchantAccountQueryData.Msg = fpRsq.Response.Msg
		return &thirdMerchantAccountQueryData, nil
	}

	fpRsq.Generate(&thirdMerchantAccountQueryData)
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
