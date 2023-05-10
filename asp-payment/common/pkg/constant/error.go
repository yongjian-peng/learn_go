package constant

const (
	MissChannelDepartPaymentParamErrMsg                = "missing channel depart payment parameter"
	MissChannelDepartNotFoundErrMsg                    = "missing channel depart not found"
	MissChannelDepartPaymentCallBackConfigErrMsg       = "missing channel depart payment callback config"
	MissMerchantProjectConfigNotFoundErrMsg            = "missing merchant project config not found"
	MissMerchantProjectCurrencyListNotFoundErrMsg      = "missing merchant project currency list not found"
	MissMerchantProjectCurrencyNotFoundErrMsg          = "missing merchant project currency not found"
	MissChannelConfigNotFoundErrMsg                    = "missing channel config not found"
	MissDepartListNotFoundErrMsg                       = "missing depart list not found"
	MissChannelDepartTradeTypeNotFoundErrMsg           = "missing channel depart trade type not found"
	MissMerchantProjectChannelDepartLinkNotFoundErrMsg = "missing merchant project channel depart link not found"
	MissOrderListNotFoundErrMsg                        = "missing order list not found"
	MissPayoutListNotFoundErrMsg                       = "missing payout list not found"
	MissIpFilterNotFoundErrMsg                         = "missing ip filter not found"
	MissSystemPayoutStatusNotFoundErrMsg               = "missing system payout status not found"
	MissBankCategoryErrMsg                             = "missing bank_category not found or status was failed"
	MissPayoutUpiValidateErrMsg                        = "missing upi_validate not found or status was failed"
	AccountIsNotActiveErrMsg                           = "Your account is not active at the moment. Please contact ASP Team." // 你的账户目前没有激活。请联系 ASP 团队。
	ConflictExceptionErrMsg                            = "Requests are conflicting"                                           // 请求有冲突
	CodeChannelConfigNotFoundErrMsg                    = "missing channel supplier not found"                                 // 渠道信息不存在
	InvalidParamsErrMsg                                = "Request parameter error Request denied"                             //请求参数错误 请求被拒绝
	ChangeErrMsg                                       = "Update %s error Please resubmit"                                    // 更新%s错误 请重新提交
	InsertErrMsg                                       = "Create %s error Please resubmit"                                    // 创建%s错误 请重新提交
	UpdateMerchantProjectCurrencyRowsErrMsg            = "update merchant project currency rows error"                        // 修改cp账户金额 行数错误
	UpdateOrderRowsErrMsg                              = "update order rows error"                                            // 修改代收订单 行数错误
	UpdatePayoutRowsErrMsg                             = "update payout rows error"                                           // 修改代收订单 行数错误
	PayoutOffErrMsg                                    = "payout is off, please wait"                                         // 代付已关闭，请稍后请求
)
