package appError

var (
	SUCCESS            = &Error{Code: 200, Msg: "Operation successful"} // 操作成功
	Unauthenticated    = &Error{Code: 401, Msg: "Request rejected with signature error"}
	ServerError        = &Error{Code: 500, Msg: "The server has deserted!"}                          // 服务器开小差了！
	ParameterError     = &Error{Code: 501, Msg: "Request parameter validate error, %s"}              //参数校验错误
	ParameterTypeError = &Error{Code: 502, Msg: "Request parameter type error"}                      //参数类型错误
	HeadParamsError    = &Error{Code: 503, Msg: "Head information parameter error"}                  // 头信息参数错误
	AppIdNoFound       = &Error{Code: 504, Msg: "AppId not found"}                                   // 登录已过期,请重新登录
	IsWait             = &Error{Code: 505, Msg: "Your operation has not yet completed, please wait"} // 您的操作尚未完成，请稍等
	TokenError         = &Error{Code: 506, Msg: "Login has expired, please login again"}             // 登录已过期,请重新登录
	PermissionsDenied  = &Error{Code: 507, Msg: "No Permission Access Denied"}                       // 无权限
	UserNotRegister    = &Error{Code: 508, Msg: "User Not Register"}                                 // 未注册
	AccountIsDisable   = &Error{Code: 509, Msg: "Account Is Disable"}                                // 用户账号已禁用
)
