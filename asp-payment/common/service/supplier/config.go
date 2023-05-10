package supplier

import (
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/service/supplier/impl/abcpay"
	"asp-payment/common/service/supplier/impl/amarquickpay"
	"asp-payment/common/service/supplier/impl/firstpayimpl"
	"asp-payment/common/service/supplier/impl/fynzonpay"
	"asp-payment/common/service/supplier/impl/haodapay"
	"asp-payment/common/service/supplier/impl/mypay"
	"asp-payment/common/service/supplier/impl/seveneight"
	"asp-payment/common/service/supplier/impl/zpayimpl"
	"asp-payment/common/service/supplier/interfaces"
)

var firstPayServer = firstpayimpl.NewPayImpl()
var zPayServer = zpayimpl.NewZPayImpl()
var seveneightServer = seveneight.NewPayImpl()
var amarquickPayServer = amarquickpay.NewPayImpl()
var abcPayServer = abcpay.NewPayImpl()
var fynzonPayServer = fynzonpay.NewPayImpl()
var myPayServer = mypay.NewPayImpl()
var haodaPayServer = haodapay.NewPayImpl()

// 注册各种上游的支付接口
var firstPayAdapterTradeTypes = []string{"firstpay.PAYOUT", "firstpay.H5", "firstpay.WAPPAY", "firstpay.MERCHANT_ACCOUNT", "firstpay.BENEFICIARY", "firstpay.PAYOUTUPI", "firstpay.UPIVALIDATE"}
var zPayAdapterTradeTypes = []string{"zpay.PAYOUT", "zpay.H5", "zpay.WAPPAY", "zpay.MERCHANT_ACCOUNT", "zpay.BENEFICIARY", "zpay.PAYOUTUPI", "zpay.UPIVALIDATE"}
var sevenEightPayAdapterTradeTypes = []string{"78pay.PAYOUT", "78pay.H5", "78pay.WAPPAY", "78pay.MERCHANT_ACCOUNT", "78pay.BENEFICIARY", "78pay.PAYOUTUPI", "78pay.UPIVALIDATE"}
var amarquickPayAdapterTradeTypes = []string{"amarquickpay.PAYOUT", "amarquickpay.H5", "amarquickpay.WAPPAY", "amarquickpay.MERCHANT_ACCOUNT", "amarquickpay.BENEFICIARY", "amarquickpay.PAYOUTUPI", "amarquickpay.UPIVALIDATE"}
var fynzonPayAdapterTradeTypes = []string{"fynzonpay.PAYOUT", "fynzonpay.H5", "fynzonpay.WAPPAY", "fynzonpay.MERCHANT_ACCOUNT", "fynzonpay.BENEFICIARY", "fynzonpay.PAYOUTUPI", "fynzonpay.UPIVALIDATE"}
var myPayAdapterTradeTypes = []string{"mypay.PAYOUT", "mypay.H5", "mypay.WAPPAY", "mypay.MERCHANT_ACCOUNT", "mypay.PAYOUTUPI", "mypay.UPIVALIDATE"}

var haodaPayAdapterTradeTypes = []string{"haodapay.PAYOUT", "haodapay.H5", "haodapay.WAPPAY", "haodapay.MERCHANT_ACCOUNT", "haodapay.BENEFICIARY", "haodapay.PAYOUTUPI", "haodapay.UPIVALIDATE"}
var abcPayAdapterTradeTypes = []string{"abcpay.PAYOUT", "abcpay.H5", "abcpay.WAPPAY", "abcpay.MERCHANT_ACCOUNT", "abcpay.PAYOUTUPI", "abcpay.UPIVALIDATE"}

func GetPaySupplierByCode(payCode string) interfaces.PayInterface {
	if goutils.InSlice(payCode, firstPayAdapterTradeTypes) {
		return firstPayServer
	}
	if goutils.InSlice(payCode, zPayAdapterTradeTypes) {
		return zPayServer
	}
	if goutils.InSlice(payCode, sevenEightPayAdapterTradeTypes) {
		return seveneightServer
	}
	if goutils.InSlice(payCode, amarquickPayAdapterTradeTypes) {
		return amarquickPayServer
	}
	if goutils.InSlice(payCode, fynzonPayAdapterTradeTypes) {
		return fynzonPayServer
	}
	if goutils.InSlice(payCode, myPayAdapterTradeTypes) {
		return myPayServer
	}
	if goutils.InSlice(payCode, haodaPayAdapterTradeTypes) {
		return haodaPayServer
	}
	if goutils.InSlice(payCode, abcPayAdapterTradeTypes) {
		return abcPayServer
	}
	return nil
}

var SupplierErrorMap = map[string]*appError.Error{
	"zpay_0000":             appError.SUCCESS,
	"firstpay_200":          appError.SUCCESS,
	"78pay_success":         appError.SUCCESS, //状态转换
	"amarquickpay_Captured": appError.SUCCESS, //状态转换
	"mypay_Success":         appError.SUCCESS, //状态转换
	"haodapay_200":          appError.SUCCESS, //状态转换

	"zpay_initcode":     appError.CodeSupplierInitClientCode,
	"firstpay_initcode": appError.CodeSupplierInitClientCode,

	"zpay_httpcode":     appError.CodeSupplierHttpCode,
	"firstpay_httpcode": appError.CodeSupplierHttpCode,

	"zpay_httperrorcode":     appError.CodeSupplierHttpErrorCode,
	"firstpay_httperrorcode": appError.CodeSupplierHttpErrorCode,

	"zpay_0001":    appError.CodeSupplierInternalChannelErrCode,       // 参数错误  request params[参数名] error.  --- 渠道内部错误 请求被拒绝
	"zpay_0002":    appError.CodeSupplierInternalChannelDepartErrCode, // 商户错误 --- 渠道内部错误商户错误 请求被拒绝
	"firstpay_500": appError.CodeSupplierInternalChannelWaitStatusErrCode,

	"firstpay_400": appError.CodeSupplierInternalChannelParamsFailedErrCode,

	"zpay_0003": appError.CodeSupplierInternalChannelErrCode, // 应用错误  Can't found application[applicationId=应用Id] record in db. --- 渠道内部错误 请求被拒绝

	"zpay_0004":     appError.CodeSupplierChannelErrCode, // 通道错误 --- 渠道分配错误 请求被拒绝
	"firstpay_0104": appError.NotImplementedErrCode,

	"zpay_0005": appError.CodeSupplierInternalChannelErrCode, // 签名错误  Verify pay sign failed. --- 渠道内部错误 请求被拒绝

	"zpay_0006": appError.CodeSupplierInternalChannelAmountErrCode, // 商户金额限制 partner amount limit . --- 渠道内部错误商户余额错误 请求被拒绝

	"zpay_0007": appError.CodeSupplierInternalChannelUpstreamErrCode, // 上游错误  partner amount limit . --- 渠道内部错误上游错误 请求被拒绝

}
