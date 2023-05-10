package appError

var (
	SUCCESS                                        = &Error{Code: 200, Message: "Operation successful"}                                                               // 操作成功
	RedirectSuccess                                = &Error{Code: 302, Message: "Operation successful"}                                                               // 操作成功
	CodeInvalidParamErrCode                        = &Error{Code: 400, Message: "Request parameter error Request denied"}                                             // 请求参数错误 请求被拒绝
	UnauthenticatedErrCode                         = &Error{Code: 401, Message: "Request rejected with signature error"}                                              // 签名错误请求被拒绝
	PermissiionDeniedErrCode                       = &Error{Code: 403, Message: "Permission error, please contact administrator"}                                     // 权限错误，请联系管理员
	CodeConflictException                          = &Error{Code: 409, Message: "Request has a conflict. For more information, see the accompanying error message"}   // 请求有冲突。有关详细信息，请参阅随附的错误消息
	CodeLimitExceededException                     = &Error{Code: 429, Message: "The request exceeded the rate limit. Retry after the specified time period"}         // 请求超出了速率限制。在指定的时间段后重试
	CodeSupplierLimitExceededException             = &Error{Code: 430, Message: "The request channel exceeded the rate limit. Retry after the specified time period"} // 请求渠道超出了速率限制。在指定的时间段后重试
	CodeUnknown                                    = &Error{Code: 500, Message: "The server has deserted!"}                                                           // 服务器开小差了！
	NotImplementedErrCode                          = &Error{Code: 501, Message: "The server does not implement this method"}                                          // 服务器没有实现该方法
	CodeUserNotExist                               = &Error{Code: 502, Message: "User handling exception, please try again"}                                          // 用户处理异常，请重试
	TokenError                                     = &Error{Code: 503, Message: "Login has expired, please login again"}                                              // 登录已过期,请重新登录
	DeadlineExceededErrCode                        = &Error{Code: 504, Message: "Request has expired"}                                                                // 请求已过期
	HeadParamsError                                = &Error{Code: 505, Message: "Head information parameter error"}                                                   // 头信息参数错误
	IsWaitErrCode                                  = &Error{Code: 506, Message: "Your operation has not yet completed, please wait"}                                  // 您的操作尚未完成，请稍等 --
	CodeInsertErr                                  = &Error{Code: 507, Message: "Write data failed, please try again"}                                                // 写入数据失败，请重试
	CodeSupplierInitClientCode                     = &Error{Code: 508, Message: "Channel initialization error Request rejected"}                                      // 渠道初始化错误 请求被拒绝
	CodeSupplierHttpCode                           = &Error{Code: 509, Message: "Channel network initialization error Request denied"}                                // 渠道网络初始化错误 请求被拒绝
	CodeSupplierHttpErrorCode                      = &Error{Code: 510, Message: "Channel network error Request denied"}                                               // 渠道网络错误 请求被拒绝
	CodeSupplierChannelErrCode                     = &Error{Code: 511, Message: "Channel assignment error Request denied"}                                            // 渠道分配错误 请求被拒绝
	CodeSupplierInternalChannelErrCode             = &Error{Code: 512, Message: "Channel internal error Request denied"}                                              // 渠道内部错误 请求被拒绝
	CodeSupplierInternalChannelAmountErrCode       = &Error{Code: 513, Message: "Internal channel error Merchant balance error Request denied"}                       // 渠道内部错误商户余额错误 请求被拒绝  Generic message
	CodeSupplierInternalChannelDepartErrCode       = &Error{Code: 514, Message: "Channel internal error merchant error request denied"}                               // 渠道内部错误商户错误 请求被拒绝 Generic message
	CodeSupplierInternalChannelUpstreamErrCode     = &Error{Code: 515, Message: "Channel Internal Error Upstream Error Request Rejected"}                             // 渠道内部错误上游错误 请求被拒绝
	CodeSupplierInternalChannelWaitStatusErrCode   = &Error{Code: 516, Message: "Awaiting callback result final status within the channel"}                           // 渠道内部等待回调结果最终状态
	CodeSupplierInternalChannelParamsFailedErrCode = &Error{Code: 517, Message: "Error in internal channel parameters, failed"}                                       // 渠道内部参数错误，失败
	UserIsBlack                                    = &Error{Code: 518, Message: "Request exception, please try again later"}

	MissNotFoundErrCode                           = &Error{Code: 4000, Message: "%s"}                                                                                                   // appid不存在 ++
	MissIdNotFoundErrCode                         = &Error{Code: 4002, Message: "missing AppId not found"}                                                                              // appid不存在 --
	MissIdNotAvailableErrCode                     = &Error{Code: 4003, Message: "AppId is not available"}                                                                               // appid不可用 --
	MissMerchantProjectNotFoundErrMsg             = &Error{Code: 4004, Message: "missing merchant project not found"}                                                                   // 商户项目不存在 --
	MissDepartListNotFoundErrMsg                  = &Error{Code: 4005, Message: "missing depart list not found"}                                                                        // 商户不存在 --
	ChannelDepartTradeTypeNotFoundErrCode         = &Error{Code: 4006, Message: "missing channel depart trade type not found"}                                                          // 渠道信息不存在 +++
	OrderNotFoundErrCode                          = &Error{Code: 4007, Message: "Order does not exist"}                                                                                 // 订单不存在 ++
	PayoutNotFoundErrCode                         = &Error{Code: 4008, Message: "Payout does not exist"}                                                                                // 订单不存在 ++
	MissIpFilterNotFoundErrMsg                    = &Error{Code: 4009, Message: "missing ip %s filter not found"}                                                                       // ip白名单不存在
	MerchantProjectCurrencyLimitExceededException = &Error{Code: 5000, Message: "The merchant amount has exceeded the limit, for more information, please refer to the merchant limit"} // 商户金额已超出限制，有关详细信息，请参阅商户限额
)

var SupplierErrorMap = map[string]*Error{
	"zpay_0000":         SUCCESS,
	"firstpay_200":      SUCCESS,
	"78pay_success":     SUCCESS,
	"abcpay_0":          SUCCESS,
	"fynzonpay_success": SUCCESS,
	"fynzonpay_0000":    SUCCESS,
	"haodapay_200":      SUCCESS,

	"zpay_initcode":     CodeSupplierInitClientCode,
	"firstpay_initcode": CodeSupplierInitClientCode,

	"zpay_httpcode":     CodeSupplierHttpCode,
	"firstpay_httpcode": CodeSupplierHttpCode,

	"zpay_httperrorcode":     CodeSupplierHttpErrorCode,
	"firstpay_httperrorcode": CodeSupplierHttpErrorCode,

	"zpay_0001":    CodeSupplierInternalChannelErrCode,       // 参数错误  request params[参数名] error.  --- 渠道内部错误 请求被拒绝
	"zpay_0002":    CodeSupplierInternalChannelDepartErrCode, // 商户错误 --- 渠道内部错误商户错误 请求被拒绝
	"firstpay_500": CodeSupplierInternalChannelWaitStatusErrCode,

	"firstpay_400": CodeSupplierInternalChannelParamsFailedErrCode,

	"zpay_0003": CodeSupplierInternalChannelErrCode, // 应用错误  Can't found application[applicationId=应用Id] record in db. --- 渠道内部错误 请求被拒绝

	"zpay_0004":     CodeSupplierChannelErrCode, // 通道错误 --- 渠道分配错误 请求被拒绝
	"firstpay_0104": NotImplementedErrCode,

	"zpay_0005": CodeSupplierInternalChannelErrCode, // 签名错误  Verify pay sign failed. --- 渠道内部错误 请求被拒绝

	"zpay_0006": CodeSupplierInternalChannelAmountErrCode, // 商户金额限制 partner amount limit . --- 渠道内部错误商户余额错误 请求被拒绝

	"zpay_0007": CodeSupplierInternalChannelUpstreamErrCode, // 上游错误  partner amount limit . --- 渠道内部错误上游错误 请求被拒绝

}
