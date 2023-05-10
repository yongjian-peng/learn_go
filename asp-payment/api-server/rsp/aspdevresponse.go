package rsp

import (
	"asp-payment/common/model"
	"github.com/spf13/cast"
)

// GenerateDevSuccessData 赋值给 创建收款接口 统一输出参数
func GenerateDevSuccessData(orderInfo *model.AspOrder) *ScanSuccessData {
	scanSuccessData := new(ScanSuccessData)
	//scanSuccessData.Id = orderInfo.Id
	scanSuccessData.Sn = orderInfo.Sn
	scanSuccessData.OutTradeNo = orderInfo.OutTradeNo
	scanSuccessData.AppId = cast.ToString(orderInfo.MchProjectId)
	scanSuccessData.FeeType = orderInfo.FeeType
	scanSuccessData.CashFeeType = orderInfo.CashFeeType
	scanSuccessData.TradeType = orderInfo.TradeType
	scanSuccessData.Provider = orderInfo.Provider
	scanSuccessData.TransactionID = orderInfo.TransactionId

	scanSuccessData.ClientIP = orderInfo.ClientIp // test
	scanSuccessData.OrderToken = orderInfo.OrderToken
	scanSuccessData.CreateTime = orderInfo.CreateTime
	scanSuccessData.FinishTime = orderInfo.FinishTime

	scanSuccessData.OrderAmount = orderInfo.TotalFee
	scanSuccessData.TotalFee = orderInfo.TotalFee

	scanSuccessData.PaymentLink = orderInfo.PaymentLink
	scanSuccessData.TradeState = orderInfo.TradeState
	scanSuccessData.NotifyURL = orderInfo.NotifyUrl
	scanSuccessData.ReturnURL = orderInfo.ReturnUrl

	scanSuccessData.CustomerName = orderInfo.CustomerName
	scanSuccessData.CustomerEmail = orderInfo.CustomerEmail
	scanSuccessData.CustomerPhone = orderInfo.CustomerPhone

	return scanSuccessData
}

// 赋值给 收款查询接口 统一输出参数
func GenerateDevOrderQuerySuccessData(orderInfo *model.AspOrder) *QuerySuccessData {
	// params := make(map[string]string)
	// params["orderNo"] = orderInfo.OutTradeNo
	// params["orderPrice"] = fmt.Sprintf("%d", orderInfo.TotalFee)
	// params["payKey"] = ""
	// params["payURL"] = payoutInfo.PayUrl

	// keys := goutils.SortMap(params)
	// sign := goutils.GetReleaseSign(params, keys, c.AspId.Key)
	querySuccessData := new(QuerySuccessData)

	querySuccessData.Id = orderInfo.Id
	querySuccessData.Sn = orderInfo.Sn
	querySuccessData.OutTradeNo = orderInfo.OutTradeNo
	querySuccessData.FeeType = orderInfo.FeeType
	querySuccessData.CashFeeType = orderInfo.CashFeeType
	querySuccessData.TradeType = orderInfo.TradeType
	querySuccessData.TradeState = orderInfo.TradeState
	querySuccessData.Provider = orderInfo.Provider
	querySuccessData.TransactionID = orderInfo.TransactionId
	querySuccessData.ClientIP = orderInfo.ClientIp
	querySuccessData.OrderToken = orderInfo.OrderToken
	querySuccessData.CreateTime = orderInfo.CreateTime
	querySuccessData.FinishTime = orderInfo.FinishTime
	querySuccessData.OrderAmount = orderInfo.TotalFee
	querySuccessData.TotalFee = orderInfo.TotalFee
	querySuccessData.PaymentLink = orderInfo.PaymentLink
	// querySuccessData.BankType = orderInfo.BankType
	querySuccessData.ReturnURL = orderInfo.ReturnUrl
	querySuccessData.NotifyURL = orderInfo.NotifyUrl
	querySuccessData.CustomerName = orderInfo.CustomerName
	querySuccessData.CustomerEmail = orderInfo.CustomerEmail
	querySuccessData.CustomerPhone = orderInfo.CustomerPhone

	return querySuccessData
}

// 赋值给 创建提现接口 统一输出参数
func GenerateDevPayoutSuccessData(payoutInfo *model.AspPayout) *PayoutSuccessData {

	payoutSuccessData := new(PayoutSuccessData)
	payoutSuccessData.Sn = payoutInfo.Sn
	payoutSuccessData.OutTradeNo = payoutInfo.OutTradeNo
	payoutSuccessData.AppId = cast.ToString(payoutInfo.MchProjectId)
	payoutSuccessData.FeeType = payoutInfo.FeeType
	payoutSuccessData.CashFeeType = payoutInfo.CashFeeType
	payoutSuccessData.TradeType = payoutInfo.TradeType
	payoutSuccessData.TradeState = payoutInfo.TradeState
	payoutSuccessData.Provider = payoutInfo.Provider
	payoutSuccessData.TransactionID = payoutInfo.TransactionId
	payoutSuccessData.ClientIP = payoutInfo.ClientIp
	payoutSuccessData.CreateTime = payoutInfo.CreateTime
	payoutSuccessData.FinishTime = payoutInfo.FinishTime
	payoutSuccessData.OrderAmount = payoutInfo.TotalFee
	payoutSuccessData.TotalFee = payoutInfo.TotalFee
	payoutSuccessData.BankType = payoutInfo.BankType
	payoutSuccessData.NotifyURL = payoutInfo.NotifyUrl
	payoutSuccessData.CustomerName = payoutInfo.CustomerName
	payoutSuccessData.CustomerEmail = payoutInfo.CustomerEmail
	payoutSuccessData.CustomerPhone = payoutInfo.CustomerPhone

	return payoutSuccessData
}

// 赋值给 审核提现接口 统一输出参数
func GenerateDevPayoutAuditSuccessData(payoutInfo *model.AspPayout) *PayoutSuccessData {
	// params := make(map[string]string)
	// params["orderNo"] = orderInfo.OutTradeNo
	// params["orderPrice"] = fmt.Sprintf("%d", orderInfo.TotalFee)
	// params["payKey"] = ""
	// params["payURL"] = payoutInfo.PayUrl

	// keys := goutils.SortMap(params)
	// sign := goutils.GetReleaseSign(params, keys, c.AspId.Key)
	payoutSuccessData := new(PayoutSuccessData)
	payoutSuccessData.Sn = payoutInfo.Sn
	payoutSuccessData.OutTradeNo = payoutInfo.OutTradeNo
	payoutSuccessData.AppId = cast.ToString(payoutInfo.MchProjectId)
	payoutSuccessData.FeeType = payoutInfo.FeeType
	payoutSuccessData.CashFeeType = payoutInfo.CashFeeType
	payoutSuccessData.TradeType = payoutInfo.TradeType
	payoutSuccessData.Provider = payoutInfo.Provider
	payoutSuccessData.TradeState = payoutInfo.TradeState
	payoutSuccessData.TransactionID = payoutInfo.TransactionId
	payoutSuccessData.ClientIP = payoutInfo.ClientIp
	payoutSuccessData.CreateTime = payoutInfo.CreateTime
	payoutSuccessData.FinishTime = payoutInfo.FinishTime
	payoutSuccessData.OrderAmount = payoutInfo.TotalFee
	payoutSuccessData.TotalFee = payoutInfo.TotalFee
	payoutSuccessData.BankType = payoutInfo.BankType
	payoutSuccessData.NotifyURL = payoutInfo.NotifyUrl
	payoutSuccessData.CustomerName = payoutInfo.CustomerName
	payoutSuccessData.CustomerEmail = payoutInfo.CustomerEmail
	payoutSuccessData.CustomerPhone = payoutInfo.CustomerPhone

	return payoutSuccessData
}

// 赋值给 创建提现接口 统一输出参数
func GenerateDevPayoutApplySuccessData(payoutInfo *model.AspPayout) *PayoutSuccessData {
	payoutSuccessData := new(PayoutSuccessData)
	payoutSuccessData.Sn = payoutInfo.Sn
	payoutSuccessData.OutTradeNo = payoutInfo.OutTradeNo
	payoutSuccessData.AppId = cast.ToString(payoutInfo.MchProjectId)
	payoutSuccessData.FeeType = payoutInfo.FeeType
	payoutSuccessData.CashFeeType = payoutInfo.CashFeeType
	payoutSuccessData.TradeType = payoutInfo.TradeType
	payoutSuccessData.Provider = payoutInfo.Provider
	payoutSuccessData.TradeState = payoutInfo.TradeState
	payoutSuccessData.TransactionID = payoutInfo.TransactionId
	payoutSuccessData.ClientIP = payoutInfo.ClientIp
	payoutSuccessData.CreateTime = payoutInfo.CreateTime
	payoutSuccessData.FinishTime = payoutInfo.FinishTime
	payoutSuccessData.OrderAmount = payoutInfo.TotalFee
	payoutSuccessData.TotalFee = payoutInfo.TotalFee
	payoutSuccessData.BankType = payoutInfo.BankType
	payoutSuccessData.NotifyURL = payoutInfo.NotifyUrl
	payoutSuccessData.CustomerName = payoutInfo.CustomerName
	payoutSuccessData.CustomerEmail = payoutInfo.CustomerEmail
	payoutSuccessData.CustomerPhone = payoutInfo.CustomerPhone

	return payoutSuccessData
}

// 赋值给 提现查询接口 统一输出参数
func GenerateDevPayoutQuerySuccessData(payoutInfo *model.AspPayout) *PayoutQuerySuccessData {
	payoutQuerySuccessData := new(PayoutQuerySuccessData)

	//payoutQuerySuccessData.Id = payoutInfo.Id
	payoutQuerySuccessData.Sn = payoutInfo.Sn
	payoutQuerySuccessData.OutTradeNo = payoutInfo.OutTradeNo
	payoutQuerySuccessData.FeeType = payoutInfo.FeeType
	payoutQuerySuccessData.CashFeeType = payoutInfo.CashFeeType
	payoutQuerySuccessData.TradeType = payoutInfo.TradeType
	payoutQuerySuccessData.Provider = payoutInfo.Provider
	payoutQuerySuccessData.TradeState = payoutInfo.TradeState
	payoutQuerySuccessData.TransactionID = payoutInfo.TransactionId
	payoutQuerySuccessData.ClientIP = payoutInfo.ClientIp
	payoutQuerySuccessData.CreateTime = payoutInfo.CreateTime
	payoutQuerySuccessData.FinishTime = payoutInfo.FinishTime
	payoutQuerySuccessData.OrderAmount = payoutInfo.TotalFee
	payoutQuerySuccessData.TotalFee = payoutInfo.TotalFee
	payoutQuerySuccessData.BankType = payoutInfo.BankType
	payoutQuerySuccessData.NotifyURL = payoutInfo.NotifyUrl
	payoutQuerySuccessData.CustomerName = payoutInfo.CustomerName
	payoutQuerySuccessData.CustomerEmail = payoutInfo.CustomerEmail
	payoutQuerySuccessData.CustomerPhone = payoutInfo.CustomerPhone

	return payoutQuerySuccessData
}
