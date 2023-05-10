package rsp

import (
	"asp-payment/api-server/req"
	"asp-payment/common/model"
	"asp-payment/common/pkg/constant"
	supplier "asp-payment/common/service/supplier/interfaces"
	"github.com/spf13/cast"
)

// GenerateSuccessData 赋值给 创建收款接口 统一输出参数
func GenerateSuccessData(scanData *supplier.ScanData, orderInfo *model.AspOrder) *ScanSuccessData {
	// params := make(map[string]string)
	// params["orderNo"] = orderInfo.OutTradeNo
	// params["orderPrice"] = fmt.Sprintf("%d", orderInfo.TotalFee)
	// params["payKey"] = ""
	// params["payURL"] = scanData.PayUrl

	// keys := goutils.SortMap(params)
	// sign := goutils.GetReleaseSign(params, keys, c.AspId.Key)
	scanSuccessData := new(ScanSuccessData)

	//scanSuccessData.Id = orderInfo.Id
	scanSuccessData.Sn = orderInfo.Sn
	scanSuccessData.OutTradeNo = orderInfo.OutTradeNo
	scanSuccessData.AppId = cast.ToString(orderInfo.MchProjectId)
	scanSuccessData.FeeType = orderInfo.FeeType
	scanSuccessData.CashFeeType = orderInfo.CashFeeType
	scanSuccessData.TradeType = orderInfo.TradeType
	scanSuccessData.Provider = orderInfo.Provider
	scanSuccessData.TransactionID = scanData.TransactionID

	scanSuccessData.ClientIP = orderInfo.ClientIp
	scanSuccessData.OrderToken = orderInfo.OrderToken
	scanSuccessData.CreateTime = orderInfo.CreateTime
	scanSuccessData.FinishTime = orderInfo.FinishTime

	scanSuccessData.OrderAmount = orderInfo.TotalFee
	scanSuccessData.TotalFee = orderInfo.TotalFee
	// scanSuccessData.OrderPrice = fmt.Sprintf("%d", orderInfo.TotalFee)

	scanSuccessData.PaymentLink = scanData.PaymentLink
	scanSuccessData.TradeState = orderInfo.TradeState
	scanSuccessData.NotifyURL = orderInfo.NotifyUrl
	scanSuccessData.ReturnURL = orderInfo.ReturnUrl

	scanSuccessData.CustomerName = orderInfo.CustomerName
	scanSuccessData.CustomerEmail = orderInfo.CustomerEmail
	scanSuccessData.CustomerPhone = orderInfo.CustomerPhone

	return scanSuccessData
}

// 赋值给 收款查询接口 统一输出参数
func GenerateOrderQuerySuccessData(orderInfo *model.AspOrder) *QuerySuccessData {
	// params := make(map[string]string)
	// params["orderNo"] = orderInfo.OutTradeNo
	// params["orderPrice"] = fmt.Sprintf("%d", orderInfo.TotalFee)
	// params["payKey"] = ""
	// params["payURL"] = scanData.PayUrl

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
func GeneratePayoutSuccessData(scanData *supplier.ThirdPayoutCreateData, payoutInfo *model.AspPayout) *PayoutSuccessData {
	// params := make(map[string]string)
	// params["orderNo"] = orderInfo.OutTradeNo
	// params["orderPrice"] = fmt.Sprintf("%d", orderInfo.TotalFee)
	// params["payKey"] = ""
	// params["payURL"] = scanData.PayUrl

	// keys := goutils.SortMap(params)
	// sign := goutils.GetReleaseSign(params, keys, c.AspId.Key)
	payoutSuccessData := new(PayoutSuccessData)
	payoutSuccessData.Sn = payoutInfo.Sn
	payoutSuccessData.OutTradeNo = payoutInfo.OutTradeNo
	payoutSuccessData.AppId = cast.ToString(payoutInfo.MchProjectId)
	payoutSuccessData.FeeType = payoutInfo.FeeType
	payoutSuccessData.CashFeeType = payoutInfo.CashFeeType
	payoutSuccessData.TradeType = payoutInfo.TradeType
	payoutSuccessData.TradeState = GetPayoutTradeState(payoutInfo.TradeState)
	payoutSuccessData.Provider = payoutInfo.Provider
	payoutSuccessData.TransactionID = scanData.TransactionID
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
func GeneratePayoutSuccessData2(payoutInfo *model.AspPayout) *PayoutSuccessData {
	payoutSuccessData := new(PayoutSuccessData)
	payoutSuccessData.Sn = payoutInfo.Sn
	payoutSuccessData.OutTradeNo = payoutInfo.OutTradeNo
	payoutSuccessData.AppId = cast.ToString(payoutInfo.MchProjectId)
	payoutSuccessData.FeeType = payoutInfo.FeeType
	payoutSuccessData.CashFeeType = payoutInfo.CashFeeType
	payoutSuccessData.TradeType = payoutInfo.TradeType
	payoutSuccessData.Provider = payoutInfo.Provider
	payoutSuccessData.TradeState = GetPayoutTradeState(payoutInfo.TradeState)
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
func GeneratePayoutAuditSuccessData(scanData *supplier.ThirdPayoutCreateData, payoutInfo *model.AspPayout) *PayoutSuccessData {
	// params := make(map[string]string)
	// params["orderNo"] = orderInfo.OutTradeNo
	// params["orderPrice"] = fmt.Sprintf("%d", orderInfo.TotalFee)
	// params["payKey"] = ""
	// params["payURL"] = scanData.PayUrl

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
	payoutSuccessData.TradeState = GetPayoutTradeState(payoutInfo.TradeState)
	payoutSuccessData.TransactionID = scanData.TransactionID
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
func GeneratePayoutApplySuccessData(payoutInfo *model.AspPayout) *PayoutSuccessData {
	payoutSuccessData := new(PayoutSuccessData)
	payoutSuccessData.Sn = payoutInfo.Sn
	payoutSuccessData.OutTradeNo = payoutInfo.OutTradeNo
	payoutSuccessData.AppId = cast.ToString(payoutInfo.MchProjectId)
	payoutSuccessData.FeeType = payoutInfo.FeeType
	payoutSuccessData.CashFeeType = payoutInfo.CashFeeType
	payoutSuccessData.TradeType = payoutInfo.TradeType
	payoutSuccessData.Provider = payoutInfo.Provider
	payoutSuccessData.TradeState = GetPayoutTradeState(payoutInfo.TradeState)
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

// GeneratePayoutQuerySuccessData 赋值给 提现查询接口 统一输出参数
func GeneratePayoutQuerySuccessData(payoutInfo *model.AspPayout) *PayoutQuerySuccessData {
	payoutQuerySuccessData := new(PayoutQuerySuccessData)

	//payoutQuerySuccessData.Id = payoutInfo.Id
	payoutQuerySuccessData.Sn = payoutInfo.Sn
	payoutQuerySuccessData.OutTradeNo = payoutInfo.OutTradeNo
	payoutQuerySuccessData.FeeType = payoutInfo.FeeType
	payoutQuerySuccessData.CashFeeType = payoutInfo.CashFeeType
	payoutQuerySuccessData.TradeType = payoutInfo.TradeType
	payoutQuerySuccessData.Provider = payoutInfo.Provider
	payoutQuerySuccessData.TradeState = GetPayoutTradeState(payoutInfo.TradeState)
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

func GenerateMerchantAccountSuccessData(thirdMerchantAccountQueryData *supplier.ThirdMerchantAccountQueryData) *MerchantAccountChannelQuerySuccessData {
	merchantAccountChannelQuerySuccessData := new(MerchantAccountChannelQuerySuccessData)
	merchantAccountChannelQuerySuccessData.Balance = thirdMerchantAccountQueryData.Balance
	merchantAccountChannelQuerySuccessData.BalanceUnsettled = thirdMerchantAccountQueryData.UnsettledBalance // 未结算余额

	return merchantAccountChannelQuerySuccessData
}

func GenerateMerchantAccountQuerySuccessData(merchantProjectCurrency *model.AspMerchantProjectCurrency) *MerchantAccountQuerySuccessData {
	merchantAccountQuerySuccessData := new(MerchantAccountQuerySuccessData)
	merchantAccountQuerySuccessData.Balance = cast.ToInt64(merchantProjectCurrency.AvailableTotalFee)       // 可用余额
	merchantAccountQuerySuccessData.BalanceIng = cast.ToInt64(merchantProjectCurrency.SettlementInProgress) // 结算中余额
	merchantAccountQuerySuccessData.BalanceFreeze = cast.ToInt64(merchantProjectCurrency.FreezeFee)         // 冻结余额
	merchantAccountQuerySuccessData.BalanceUnsettled = cast.ToInt64(merchantProjectCurrency.TotalFee)       // 待结算余额
	merchantAccountQuerySuccessData.AppId = cast.ToInt(merchantProjectCurrency.MchProjectId)

	return merchantAccountQuerySuccessData
}

func GenerateMerchantProjectSuccessData(deptTradeTypeInfo *req.DeptTradeTypeInfo) *MerchantProjectChannelSuccessData {
	successData := new(MerchantProjectChannelSuccessData)
	successData.DepartID = deptTradeTypeInfo.DepartID
	successData.ChannelID = deptTradeTypeInfo.ChannelID
	successData.Provider = deptTradeTypeInfo.Provider
	successData.Payment = deptTradeTypeInfo.Payment
	successData.TradeType = deptTradeTypeInfo.TradeType
	successData.H5Type = deptTradeTypeInfo.H5Type

	return successData
}

// GetPayoutTradeState 状态的映射关系
func GetPayoutTradeState(tradeState string) string {
	tradeStates := make(map[string]string)
	tradeStates[constant.PAYOUT_TRADE_STATE_APPLY] = constant.PAYOUT_TRADE_STATE_PENDING
	tradeStates[constant.PAYOUT_TRADE_STATE_RETURN] = constant.PAYOUT_TRADE_STATE_RETURN
	tradeStates[constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS] = constant.PAYOUT_TRADE_STATE_PENDING
	tradeStates[constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING] = constant.PAYOUT_TRADE_STATE_PENDING
	tradeStates[constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED] = constant.PAYOUT_TRADE_STATE_FAILED
	tradeStates[constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS] = constant.PAYOUT_TRADE_STATE_SUCCESS
	tradeStates[constant.PAYOUT_TRADE_STATE_SUCCESS] = constant.PAYOUT_TRADE_STATE_SUCCESS
	tradeStates[constant.PAYOUT_TRADE_STATE_FAILED] = constant.PAYOUT_TRADE_STATE_FAILED
	tradeStates[constant.PAYOUT_TRADE_STATE_REVOKE] = constant.PAYOUT_TRADE_STATE_RETURN

	status := ""
	if payoutStatus, ok := tradeStates[tradeState]; ok {
		status = payoutStatus
	} else {
		status = constant.PAYOUT_TRADE_STATE_FAILED
	}
	return status
}
