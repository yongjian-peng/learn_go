package service

import (
	"asp-payment/api-server/req"
	"asp-payment/api-server/rsp"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goRedis"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"asp-payment/common/service/cashierdesk"
	"asp-payment/common/service/supplier"
	"asp-payment/common/service/supplier/interfaces"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"

	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type OrderServer struct {
	*Service
}

func NewPayOrderServer(c *fiber.Ctx) *OrderServer {
	return &OrderServer{Service: NewService(c, constant.OrderServerLogFileName)}
}
func NewSimplePayOrderServer() *OrderServer {
	return &OrderServer{&Service{LogFileName: constant.OrderServerLogFileName}}
}

// Pay 支付
func (s *OrderServer) Pay() error {
	// 贯穿支付所需要的参数
	reqAspPayment := new(req.AspPayment)
	if err := s.C.BodyParser(reqAspPayment); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "s.C.BodyParser(reqAspPayment) error ", zap.Error(err))
		return s.Error(appError.NewError(err.Error()))
	}
	// logger.ApiInfo(s.LogFileName, s.RequestId, "xxx")
	if err := s.DealUnifiedOrderParams(reqAspPayment, s.Head); err != nil {
		return err
	}
	lockName := goRedis.GetKey(fmt.Sprintf("pay:order:%s_%s", s.Head.AppId, reqAspPayment.UserID))
	flag := goRedis.Lock(lockName)
	if !flag {
		return appError.IsWaitErrCode
	}
	defer goRedis.UnLock(lockName)

	merchantProjectServer := NewMerchantProjectServer(s.C)
	MerchantProjectInfo, err := merchantProjectServer.GetMerchantProjectInfo(s.Head.AppId)
	if err != nil {
		return s.Error(err)
	}

	AspIdInfo, mErr := merchantProjectServer.GetIdInfo()
	if mErr != nil {
		return s.Error(mErr)
	}

	// 验签
	if err = s.CheckSign(reqAspPayment, AspIdInfo.Key); err != nil {
		return s.Error(err)
	}
	// 提现的验证 如果验证失败 则提示 验证收款订单号
	if err = s.VerifyPaymentBefore(MerchantProjectInfo, reqAspPayment); err != nil {
		return s.Error(err)
	}

	//验证黑名单列表
	if err = s.VerifyApiBlack(reqAspPayment.CustomerPhone, reqAspPayment.DeviceInfo, "", ""); err != nil {
		return s.Error(err)
	}

	// 查询到cp 项目的配置
	merchantProjectConfigInfo, pErr := merchantProjectServer.GetMerchantAccountConfigInfo()
	if pErr != nil {
		return s.Error(pErr)
	}
	merchantProjectCurrencyInfo, mpErr := merchantProjectServer.GetMerchantProjectCurrencyInfo(reqAspPayment.OrderCurrency, MerchantProjectInfo.Id)
	if mpErr != nil {
		return s.Error(mpErr)
	}

	// 查询到 商户的可用的支付渠道 过程中有赋值给 AspChannelDepartTradeType
	aspChannelDepartTradeType, cErr := NewDepartServer(s.C).ChooseChannelDepartTradeType(reqAspPayment)
	if cErr != nil {
		return s.Error(cErr)
	}

	paymentArr := strings.Split(aspChannelDepartTradeType.Payment, ".")
	if len(paymentArr) != 2 {
		return s.Error(appError.MissNotFoundErrCode.FormatMessage(constant.MissChannelDepartPaymentParamErrMsg))
	}
	reqAspPayment.PayParams.Adapter = paymentArr[1]

	err = merchantProjectServer.MerchantProjectUserInsertOrUpdate(reqAspPayment.UserID, MerchantProjectInfo)
	if err != nil {
		return s.Error(err) // 用户处理异常，请重试
	}

	var orderInfo *model.AspOrder
	// 生成订单记录 提前预插入订单数据
	orderInfo, err = s.Insert(reqAspPayment, MerchantProjectInfo, merchantProjectConfigInfo, merchantProjectCurrencyInfo, aspChannelDepartTradeType)
	if err != nil {
		return s.Error(err) // 写入数据失败，请重试
	}

	if config.IsTestEnv() && s.Head.AppId != constant.IGNORE_MERCHANT_PROJECT_ID {
		// 更新操作 修改订单状态即可
		scanSuccessData, err := s.DevPayAction(orderInfo)
		if err != nil {
			return s.Error(err)
		}
		return s.Success(scanSuccessData)
	}

	// 获取到上游数据 获取到对应的上游
	supplierCode := orderInfo.Adapter + "." + orderInfo.TradeType
	//fmt.Println("supplierCode: ", supplierCode)
	//return s.Error(appError.ChannelDepartTradeTypeNotFoundErrCode) // 渠道信息不存在
	// supplierCode := "firstpay.H5"
	paySupplier := supplier.GetPaySupplierByCode(supplierCode)
	if paySupplier == nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "thirdParty.GetPaySupplierByCode error ", zap.Error(err), zap.String("supplierCode", supplierCode), zap.Any("paySupplier", paySupplier))
		return s.Error(appError.ChannelDepartTradeTypeNotFoundErrCode) // 渠道信息不存在
	}

	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(orderInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(orderInfo.ChannelId)
	//查询账户在上游的配置
	var channelDepartInfo *model.AspChannelDepartConfig
	channelDepartInfo, err = NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if err != nil {
		return s.Error(err) // missing channel depart not found
	}
	var scanData *interfaces.ScanData

	if orderInfo.H5Type == constant.H5Type_H5 {
		//没有错误，就是成功状态 更新发版
		scanData, err = paySupplier.H5(s.RequestId, channelDepartInfo, orderInfo)
	} else if orderInfo.H5Type == constant.H5Type_WAPPAY {
		scanData = &interfaces.ScanData{}
		scanData.TransactionID = ""
		scanData.PaymentLink = fmt.Sprintf("%s/%s", config.AppConfig.Urls.WapPayCheckoutPreUrl, orderInfo.Sn)
		// 更新订单信息 入参结构体 返回error
		params := make(map[string]interface{})
		params["transaction_id"] = scanData.TransactionID
		params["cash_fee"] = orderInfo.CashFee
		params["payment_link"] = scanData.PaymentLink
		params["return_url"] = s.GetOrderReturnUrl(orderInfo) // 映射return_url
		// h5 返回的是支付的url 链接是需要用户根据url去完成支付的动作的。所以直接更新
		var newOrderInfo *model.AspOrder
		newOrderInfo, err = s.SolveOrderPaySuccess(orderInfo.Id, params, "orderPayUpdate", nil)
		if err != nil {
			return s.Error(appError.CodeUnknown) // 服务器开小差了！
		}
		// 统一返回参数 转换
		scanSuccessData := rsp.GenerateSuccessData(scanData, newOrderInfo)
		return s.Success(scanSuccessData)

	} else {
		return s.Error(appError.CodeInvalidParamErrCode)
	}

	if err != nil {
		//上游返回异常信息，直接修改订单状态为异常，并记录上游异常信息
		if err.Code == appError.CodeSupplierChannelErrCode.Code {
			//先记录日志
			logger.ApiWarn(s.LogFileName, s.RequestId, "request order up supplier rsp err", zap.Any("code", err.Code), zap.String("msg", err.Message))
			aspOrder := model.AspOrder{}
			params := make(map[string]interface{})
			params["trade_state"] = constant.ORDER_TRADE_STATE_PAYERROR
			params["supplier_return_code"] = scanData.Code
			params["supplier_return_msg"] = scanData.Msg
			//修改订单状态
			resOrderUpdate := database.DB.Model(&aspOrder).Where("id = ?", orderInfo.Id).Where("trade_state = ?", constant.ORDER_TRADE_STATE_PENDING).Updates(params)
			if resOrderUpdate.RowsAffected < 1 || resOrderUpdate.Error != nil {
				logger.ApiWarn(s.LogFileName, s.RequestId, "修改订单状态失败: ", zap.Int("orderId: ", orderInfo.Id), zap.Any("params", params))
			}
		}

		return s.Error(err)
	}
	//// 渠道映射系统的统一返回字符串 在 appError中有定义的 map
	//supplierStr := orderInfo.Adapter + "_" + scanData.Code
	//supplierError, ok := supplier.SupplierErrorMap[supplierStr]
	//if !ok {
	//	logger.ApiWarn(s.LogFileName, s.RequestId, "Response json new Status ", zap.String("newCode", supplierStr))
	//	return s.Error(appError.CodeSupplierInternalChannelErrCode)
	//}
	//
	//if supplierError.Code != appError.SUCCESS.Code {
	//	return s.Error(supplierError)
	//}

	// 更新订单信息 入参结构体 返回error
	params := make(map[string]interface{})
	params["transaction_id"] = scanData.TransactionID
	params["cash_fee"] = scanData.CashFee
	params["payment_link"] = scanData.PaymentLink
	params["payments_url"] = scanData.PaymentsURL
	// h5 返回的是支付的url 链接是需要用户根据url去完成支付的动作的。所以直接更新
	var newOrderInfo *model.AspOrder
	newOrderInfo, err = s.SolveOrderPaySuccess(orderInfo.Id, params, "orderPayUpdate", nil)
	if err != nil {
		return s.Error(appError.CodeUnknown) // 服务器开小差了！
	}
	// 统一返回参数 转换
	scanSuccessData := rsp.GenerateSuccessData(scanData, newOrderInfo)
	return s.Success(scanSuccessData)
}

func (s *OrderServer) DevPayAction(orderInfo *model.AspOrder) (*rsp.ScanSuccessData, *appError.Error) {
	orderSuccessUpdate := make(map[string]interface{})
	orderSuccessUpdate["transaction_id"] = goutils.RandomString(20)
	orderSuccessUpdate["cash_fee"] = orderInfo.CashFee
	orderSuccessUpdate["payment_link"] = orderInfo.PaymentLink
	orderSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
	orderSuccessUpdate["trade_state"] = constant.ORDER_TRADE_STATE_SUCCESS

	MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
	totalFee := cast.ToInt(orderInfo.TotalFee)
	amount := totalFee - orderInfo.ChargeFee - orderInfo.FixedAmount
	err := MerchantProjectRepository.ChangeMerchantProjectCurrentByTest(orderInfo.MchProjectId, 0, amount, 0, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYORDER, orderInfo.Id, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_SUCCESS, orderSuccessUpdate)
	if err != nil {
		return nil, appError.NewError(constant.ChangeErrMsg).FormatMessage("orderEnv") // 更新代付订单错误 请重新提交
	}
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["id"] = orderInfo.Id
	newOrderInfo, err := s.GetOrderInfo(orderInfoParam, nil)
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "s.GetOrderInfo  ", zap.Error(err))
		return nil, appError.OrderNotFoundErrCode // 订单不存在
	}
	_ = NewSendQueueServer().SendNotifyQueue(newOrderInfo.Sn)
	scanSuccessData := rsp.GenerateDevSuccessData(newOrderInfo)
	return scanSuccessData, nil
}

func (s *OrderServer) renderError(title, message string) error {
	return s.C.Render("order/error", fiber.Map{
		"title": title,
		"error": message,
	})
}

// CheckOut 收银台显示页面
func (s *OrderServer) CheckOut() error {
	sn := s.C.Params("sn", "")
	redirectType := s.C.Params("redirectType", "")
	payType := s.C.Params("payType", constant.CASHIERDESK_PAY_TYPE_UPI)
	title := "Recharge Amount"
	checkoutReq := &req.CheckOutReq{Sn: sn}
	err := checker.Struct(checkoutReq)
	if err != nil {
		return s.renderError(title, err.Message)
	}
	// 加锁判断 防止重复请求
	lockName := goRedis.GetKey(fmt.Sprintf("pay:orderCheckOut:%s", sn))
	if !goRedis.Lock(lockName) {
		return appError.IsWaitErrCode
	}
	defer goRedis.UnLock(lockName)

	//查询订单信息
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["sn"] = sn
	orderInfo, oErr := s.GetOrderInfo(orderInfoParam, nil)
	if oErr != nil {
		// 订单不存在
		return s.renderError(title, appError.OrderNotFoundErrCode.Message)
	}
	// 获取到上游数据 获取到对应的上游
	supplierCode := orderInfo.Adapter + "." + orderInfo.TradeType
	paySupplier := supplier.GetPaySupplierByCode(supplierCode)
	if paySupplier == nil {
		return s.renderError(title, appError.ChannelDepartTradeTypeNotFoundErrCode.Message)
	}

	cashierDeskInterface := cashierdesk.GetCashierDeskByCode(orderInfo.Adapter)
	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(orderInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(orderInfo.ChannelId)
	//查询账户在上游的配置
	var channelDepartInfo *model.AspChannelDepartConfig
	channelDepartInfo, err = NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if err != nil {
		return s.renderError(title, err.Message)
	}
	//fmt.Println("orderInfo: ", orderInfo)
	//fmt.Println("cashierDeskInterface: ", cashierDeskInterface)
	//return s.renderError(title, "err.Message")
	// 需要重定向到指定链接或者 深链地址
	if redirectType == "1" {
		payTypes := []string{constant.CASHIERDESK_PAY_TYPE_UPI, constant.CASHIERDESK_PAY_TYPE_GPAY, constant.CASHIERDESK_PAY_TYPE_PAYTM, constant.CASHIERDESK_PAY_TYPE_PHONEPE}
		if payType == "" || !goutils.InArray(payType, payTypes) {
			return s.renderError(title, "Request Error")
		}
		redirectUrl, _ := cashierDeskInterface.GetPaymentIntentUrl(orderInfo)
		if orderInfo.PaymentsUrl == "" {
			// 请求统一处理 wappay 统一请求处理上游
			bm, wErr := paySupplier.WAPPAY(s.RequestId, channelDepartInfo, orderInfo)

			fmt.Println("bm: ", bm)
			if wErr != nil {
				return s.renderError(title, wErr.Message)
			}
			// TODO 是否需要处理302跳转 待后期优化
			// 更新订单信息 入参结构体 返回error
			params := make(map[string]interface{})
			params["transaction_id"] = bm["transactionId"]
			params["payments_url"] = bm["url"]
			// wappay 返回的是支付的 upi的深链 需要记录到数据库
			orderInfo, err = s.SolveOrderPaySuccess(orderInfo.Id, params, "orderPayRedirectUpdate", nil)
			if err != nil {
				return s.Error(appError.CodeUnknown) // 服务器开小差了！
			}
			//redirectUrl := "<a url=\"https://baidu.com\"><a/>"
			redirectUrl, _ = cashierDeskInterface.GetPaymentIntentUrl(orderInfo)
		}

		upiPrefix := constant.CASHIERDESK_UPIPREFIX // upi://
		// 做多一层考虑，如果上游返回的upi地址前缀有变动，则抛出异常
		if upiPrefix != redirectUrl[:6] {
			return s.renderError(title, "Request Error")
		}
		switch payType {
		case constant.CASHIERDESK_PAY_TYPE_GPAY:
			upiPrefix = "gpay://upi/"
		case constant.CASHIERDESK_PAY_TYPE_PAYTM:
			upiPrefix = "paytmmp://"
		case constant.CASHIERDESK_PAY_TYPE_PHONEPE:
			upiPrefix = "phonepe://"
		}
		// 根据规则 替换前缀  upi:// 去掉 然后替换成对应的 应用 app 深链的标识
		redirectUrl = upiPrefix + redirectUrl[6:]

		//redirectUrl = "upi://pay?pa=604751@icici&pn=SURYA ENTERPRISES&tr=EZV2023032112052390684001&am=100.0&cu=INR&mc=5411"
		//redirectUrl = "https://www.baidu.com"
		// 返回 302 重定向
		return s.RedirectSuccess(redirectUrl)
	}

	if redirectType == "" {
		if orderInfo.H5Type == constant.H5Type_WAPPAY {

			// mypay 需要有更新 更新upi 深链地址
			// 包含事件 渲染页面 wappay 返回对接的加密等信息 再次进入到链接中 和 请求 wappay2 统一返回参数处理
			// 统一处理参数 组装参数 页面路径
			bm, wErr := cashierDeskInterface.Assembling(s.RequestId, channelDepartInfo, orderInfo)
			if wErr != nil {
				return s.renderError(title, wErr.Message)
			}
			// fmt.Println("bm: ", bm)

			redirectStr, fm, _ := cashierDeskInterface.Rendering(&bm)
			fmt.Println("fm: ", fm)
			return s.C.Render(redirectStr, fiber.Map{
				"sn":     sn,
				"amount": goutils.Fen2Yuan(cast.ToInt64(orderInfo.TotalFee)),
				"title":  title,
				"params": fm,
			})
		}
	}

	return s.renderError(title, appError.CodeInvalidParamErrCode.Message)
}

// QrCode 收银台二维码页面
func (s *OrderServer) QrCode() error {
	sn := s.C.Params("sn", "")
	title := "Recharge Amount"
	checkoutReq := &req.CheckOutReq{Sn: sn}
	err := checker.Struct(checkoutReq)
	if err != nil {
		return s.renderError(title, err.Message)
	}
	// 加锁判断 防止重复请求
	lockName := goRedis.GetKey(fmt.Sprintf("pay:orderQrCode:%s", sn))
	if !goRedis.Lock(lockName) {
		return appError.IsWaitErrCode
	}
	defer goRedis.UnLock(lockName)

	//查询订单信息
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["sn"] = sn
	orderInfo, oErr := s.GetOrderInfo(orderInfoParam, nil)
	if oErr != nil {
		// 订单不存在
		return s.renderError(title, appError.OrderNotFoundErrCode.Message)
	}

	// 获取到上游数据 获取到对应的上游
	supplierCode := orderInfo.Adapter + "." + orderInfo.TradeType
	paySupplier := supplier.GetPaySupplierByCode(supplierCode)
	if paySupplier == nil {
		return s.renderError(title, appError.ChannelDepartTradeTypeNotFoundErrCode.Message)
	}

	cashierDeskInterface := cashierdesk.GetCashierDeskByCode(orderInfo.Adapter)
	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(orderInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(orderInfo.ChannelId)
	//查询账户在上游的配置
	var channelDepartInfo *model.AspChannelDepartConfig
	channelDepartInfo, err = NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if err != nil {
		return s.renderError(title, err.Message)
	}

	redirectUrl, _ := cashierDeskInterface.GetPaymentIntentUrl(orderInfo)
	redirectQrUrl, _ := cashierDeskInterface.GetPaymentQrUrl(orderInfo)
	if orderInfo.PaymentsUrl == "" {
		// 请求统一处理 wappay 统一请求处理上游
		bm, wErr := paySupplier.WAPPAY(s.RequestId, channelDepartInfo, orderInfo)

		fmt.Println("bm: ", bm)
		if wErr != nil {
			return s.renderError(title, wErr.Message)
		}
		// TODO 是否需要处理302跳转 待后期优化
		// 更新订单信息 入参结构体 返回error
		params := make(map[string]interface{})
		params["transaction_id"] = bm["transactionId"]
		params["payments_url"] = bm["url"]
		// wappay 返回的是支付的 upi的深链 需要记录到数据库
		orderInfo, err = s.SolveOrderPaySuccess(orderInfo.Id, params, "orderPayRedirectUpdate", nil)
		if err != nil {
			return s.Error(appError.CodeUnknown) // 服务器开小差了！
		}
		//redirectUrl := "<a url=\"https://baidu.com\"><a/>"
		redirectUrl, _ = cashierDeskInterface.GetPaymentIntentUrl(orderInfo)
		redirectQrUrl, _ = cashierDeskInterface.GetPaymentQrUrl(orderInfo)
	}

	return s.C.Render("order/qrcode", fiber.Map{
		"sn":     sn,
		"amount": goutils.Fen2Yuan(cast.ToInt64(orderInfo.TotalFee)),
		"title":  title,
		"params": map[string]string{
			"redirectUrl":   redirectUrl,
			"redirectQrUrl": redirectQrUrl,
		},
	})
}
func (s *OrderServer) OrderQueryApi() error {
	// 贯穿支付所需要的参数
	reqPayQuery := new(req.AspOrderQuery)
	if err := s.C.QueryParser(reqPayQuery); err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	if err := s.DealUnifiedOrderQueryParams(reqPayQuery, s.Head); err != nil {
		return s.Error(err)
	}

	merchantProjectServer := NewMerchantProjectServer(s.C)
	_, err := merchantProjectServer.GetMerchantProjectInfo(s.Head.AppId)
	if err != nil {
		return s.Error(err)
	}
	AspIdInfo, mpErr := merchantProjectServer.GetIdInfo()
	if mpErr != nil {
		return s.Error(mpErr)
	}

	// 验签
	if err = s.CheckSignOrderQuery(AspIdInfo.Key, reqPayQuery, s.Head.Timestamp); err != nil {
		return s.Error(err)
	}

	// 1. 首先查询订单
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["sn"] = reqPayQuery.Sn
	orderInfo, oErr := s.GetOrderInfo(orderInfoParam, nil)

	if oErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "logic.GetOrderInfo  ", zap.Error(oErr))
		return s.Error(appError.OrderNotFoundErrCode) // 订单不存在
	}

	if orderInfo.MchProjectId != cast.ToInt(s.Head.AppId) {
		return s.Error(appError.CodeInvalidParamErrCode) // 请求参数错误
	}

	params := map[string]string{
		"is_call_upstream": reqPayQuery.IsCallUpstream,
	}
	// 2. 根据是否需要走上游
	isCallUpstream := s.IsNeedOrderQueryUpstream(orderInfo, params)

	// 3. 如果需要走上游 拼接订单渠道中
	//var querySuccessData *req.QuerySuccessData
	if isCallUpstream == false || (config.IsTestEnv() && s.Head.AppId != constant.IGNORE_MERCHANT_PROJECT_ID) {
		scanSuccessData := rsp.GenerateDevOrderQuerySuccessData(orderInfo)
		return s.Success(scanSuccessData)
	}

	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(orderInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(orderInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, ccErr := NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if ccErr != nil {
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartNotFoundErrMsg) // missing channel depart not found
	}

	newOrderInfo, oqErr := s.OrderQuery(orderInfo, channelDepartInfo)
	if oqErr != nil {
		return s.Error(oqErr)
	}
	// 统一返回参数 转换
	scanSuccessData := rsp.GenerateOrderQuerySuccessData(newOrderInfo)
	return s.Success(scanSuccessData)
}

func (s *OrderServer) ReturnSuccess() error {
	return s.C.Render("order/return_success", fiber.Map{
		"title": "payment success",
	})
	return nil
}

func (s *OrderServer) OrderQueryJob() error {

	return nil
}

func (s *OrderServer) OrderQuery(orderInfo *model.AspOrder, channelDepartInfo *model.AspChannelDepartConfig) (*model.AspOrder, *appError.Error) {
	scanQueryData, err := s.DoQueryUpstream(orderInfo, channelDepartInfo)
	if err != nil {
		return nil, err
	}

	upstreamStatus := scanQueryData.TradeState // 字符串 success
	isCapitalFlow := false
	newOrderInfo := &model.AspOrder{}
	orderInfoUpdateParam := make(map[string]interface{})
	// 中间有 代收中的状态 状态更改是不可逆序的 需要注意
	if orderInfo.TradeState != upstreamStatus {
		//upstreamCode := zpRsq.Response.QueryOrderData.Status        // 上游的状态 数字
		//orderCode := GetZPayPaymentTradeState(orderInfo.TradeState) // 订单的状态 转成数字
		orderInfoUpdateParam["transaction_id"] = scanQueryData.TransactionID
		orderInfoUpdateParam["cash_fee"] = scanQueryData.CashFee
		orderInfoUpdateParam["payment_link"] = orderInfo.PaymentLink
		orderInfoUpdateParam["finish_time"] = scanQueryData.FinishTime
		orderInfoUpdateParam["utr"] = orderInfo.Utr
		orderInfoUpdateParam["trade_state"] = upstreamStatus
		// 如果成功 必须提现状态是 顺序的增长的，不能状态逆序
		// 如果上游返回成功 提现状态是 已申请 则修改
		if upstreamStatus == constant.ORDER_TRADE_STATE_SUCCESS {
			errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
				MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
				totalFee := cast.ToInt(orderInfo.TotalFee)
				amount := totalFee - orderInfo.ChargeFee - orderInfo.FixedAmount
				// 代收上游成功 待结算余额新增    收支流水记录新增
				orderInfoUpdateParam["trade_state"] = upstreamStatus
				orderInfoUpdateParam["utr"] = scanQueryData.Utr
				orderInfoUpdateParam["finish_time"] = goutils.GetDateTimeUnix()
				if _, err = s.SolveOrderPaySuccess(orderInfo.Id, orderInfoUpdateParam, "orderQueryUpdate", tx); err != nil {
					logger.ApiError(s.LogFileName, s.RequestId, "s.SolveOrderPaySuccess orderQueryUpdate", zap.Error(err))
					return appError.CodeUnknown // 服务器开小差了！
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				errChange := MerchantProjectRepository.PayinOrderSuccess(orderInfo.MchProjectId, amount, orderInfo.Id, tx)
				if errChange != nil {
					return appError.NewError(constant.ChangeErrMsg).FormatMessage("OrderQueryUpdate ChangeMerchantProjectCurrent") // 更新代付订单错误 请重新提交
				}
				return nil
			})
			if errTrans != nil {
				return nil, appError.NewError(errTrans.Error())
			}

			isCapitalFlow = true
			// 查询到最新的数据
			orderInfoParam := make(map[string]interface{})
			orderInfoParam["id"] = orderInfo.Id
			newOrderInfo, _ = s.GetOrderInfo(orderInfoParam, nil)
		}

		if upstreamStatus == constant.ORDER_TRADE_STATE_FAILED {
			orderInfoUpdateParam["trade_state"] = upstreamStatus
			if newOrderInfo, err = s.SolveOrderPaySuccess(orderInfo.Id, orderInfoUpdateParam, "orderQueryUpdate", nil); err != nil {
				logger.ApiError(s.LogFileName, s.RequestId, "s.SolveOrderPaySuccess orderQueryUpdate", zap.Error(err))
				return nil, appError.CodeUnknown // 服务器开小差了！
			}
		}

		// TODO 如果是 已完成状态 则需要添加到 redis 中 通知到下游商户 提现回调
		if isCapitalFlow == true {
			// 写入数据到 asp_merchant_project_capital_flow 需要开启事务来处理
			// 修改 asp_merchant_project_currency 中的金额
			_ = NewSendQueueServer().SendNotifyQueue(orderInfo.Sn)

			//merchantProjectQueue := new(*req.MerchantProjectTotalFeeQueue)
			//merchantProjectQueue.ProductID = orderInfo.Id
			//merchantProjectQueue.ProductType = model.SEND_TYPE_ORDER_SUCCESS
		}
	} else if orderInfo.TransactionId == "" {
		orderInfoUpdateParam["transaction_id"] = scanQueryData.TransactionID
		orderInfoUpdateParam["cash_fee"] = orderInfo.CashFee
		orderInfoUpdateParam["payment_link"] = orderInfo.PaymentLink
		orderInfoUpdateParam["finish_time"] = orderInfo.FinishTime
		orderInfoUpdateParam["trade_state"] = orderInfo.TradeState
		newOrderInfo, err = s.SolveOrderPaySuccess(orderInfo.Id, orderInfoUpdateParam, "orderQueryUpdate", nil)
		if err != nil {
			return nil, appError.CodeUnknown
		}
	} else {
		newOrderInfo = orderInfo
	}
	return newOrderInfo, nil
}

func (s *OrderServer) DoQueryUpstream(orderInfo *model.AspOrder, channelDepartInfo *model.AspChannelDepartConfig) (*interfaces.ThirdQueryData, *appError.Error) {
	// 获取到上游数据 获取到对应的上游
	supplierCode := orderInfo.Adapter + "." + orderInfo.TradeType
	// supplierCode := "firstpay.H5"
	paySupplier := supplier.GetPaySupplierByCode(supplierCode)
	if paySupplier == nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "thirdParty.GetPaySupplierByCode error ", zap.String("err", "supplier == nil"))
		return nil, appError.ChannelDepartTradeTypeNotFoundErrCode // 渠道信息不存在
	}

	scanQueryData, scErr := paySupplier.PayQuery(s.RequestId, channelDepartInfo, orderInfo)
	//fmt.Println("scanQueryData: ", scanQueryData)
	if scErr != nil {
		//上游返回异常信息，直接修改订单状态为异常，并记录上游异常信息
		if scErr.Code == appError.CodeSupplierInternalChannelParamsFailedErrCode.Code {
			if orderInfo.TradeState != constant.ORDER_TRADE_STATE_PAYERROR {
				//先记录日志
				logger.ApiWarn(s.LogFileName, s.RequestId, "request order up supplier rsp err", zap.Any("code", scErr.Code), zap.String("msg", scErr.Message))
				aspOrder := model.AspOrder{}
				params := make(map[string]interface{})
				params["trade_state"] = constant.ORDER_TRADE_STATE_PAYERROR
				params["supplier_return_code"] = scanQueryData.Code
				params["supplier_return_msg"] = scanQueryData.Msg
				//修改订单状态
				resOrderUpdate := database.DB.Model(&aspOrder).Where("id = ?", orderInfo.Id).Where("trade_state = ?", constant.ORDER_TRADE_STATE_PENDING).Updates(params)
				if resOrderUpdate.RowsAffected < 1 || resOrderUpdate.Error != nil {
					logger.ApiWarn(s.LogFileName, s.RequestId, "修改订单状态失败: ", zap.Int("orderId: ", orderInfo.Id), zap.Any("params", params))
				}
			}
		}
		return nil, scErr
	}
	// 渠道映射系统的统一返回字符串 在 appError中有定义的 map
	supplierStr := orderInfo.Adapter + "_" + scanQueryData.Code

	//fmt.Println("suppliderStr:", supplierStr)

	supplierError, ok := appError.SupplierErrorMap[supplierStr]
	if !ok {
		logger.ApiWarn(s.LogFileName, s.RequestId, "Response json new Status ", zap.String("newCode", supplierStr))
		return nil, appError.CodeSupplierInternalChannelErrCode
	}
	if supplierError.Code != appError.SUCCESS.Code {
		return nil, supplierError
	}
	return scanQueryData, nil
}

// DealUnifiedOrderParams 处理请求参数 赋值一些基础值 例如：client_ip
func (s *OrderServer) DealUnifiedOrderParams(d *req.AspPayment, h *req.AspPaymentHeader) *appError.Error {

	if err := checker.Struct(d); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedOrderParams:", zap.Error(err))
		return err
	}

	if err := checker.Struct(h); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedOrderParams:", zap.Error(err))
		return err
	}

	d.PayParams.PaymentMethod = d.PaymentMethod
	d.PayParams.AppId = h.AppId
	d.PayParams.OrderID = d.OrderID
	d.PayParams.UserID = d.UserID
	d.PayParams.OrderCurrency = d.OrderCurrency
	d.PayParams.OrderAmount = d.OrderAmount
	d.PayParams.Timestamp = h.Timestamp
	d.PayParams.OrderName = d.OrderName
	d.PayParams.ReturnURL = d.ReturnURL
	d.PayParams.NotifyURL = d.NotifyURL
	d.PayParams.CustomerName = d.CustomerName
	d.PayParams.CustomerPhone = d.CustomerPhone
	d.PayParams.CustomerEmail = d.CustomerEmail
	d.PayParams.DeviceInfo = d.DeviceInfo
	d.PayParams.OrderNote = d.OrderNote
	d.PayParams.Attach = d.Attach
	d.PayParams.Sign = h.Signature

	// 接收参数 需要 e.orm
	d.PayParams.ClientIp = s.C.IP()

	paymentArr := strings.Split(d.PaymentMethod, ".")

	if len(paymentArr) != 2 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedOrderParams with invalid params payment_method:", zap.String("sign", constant.InvalidParamsErrMsg))
		return appError.CodeInvalidParamErrCode // 请求参数错误 请求被拒绝
	}

	d.PayParams.Provider = strings.ToLower(paymentArr[0])  // 统一处理转小写
	d.PayParams.TradeType = strings.ToUpper(paymentArr[1]) // 统一处理成转大写
	serialNumBer := goutils.GenerateSerialNumBer("", config.AppConfig.Server.Name, config.AppConfig.Server.Env)
	if serialNumBer == "" {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GenerateSerialNumBer Repeat")
		return appError.CodeUnknown // 服务器开小差了！
	}
	d.PayParams.Sn = serialNumBer
	return nil
}

// CheckSign 处理支付成功，加款等各项操作 更新订单上游返回的数据
func (s *OrderServer) CheckSign(d *req.AspPayment, paySecret string) *appError.Error {

	if s.Head.AppId == "" {
		return appError.CodeInvalidParamErrCode // 请求参数错误 请求被拒绝
	}

	// TODO 需要完善对应的字段 验证签名
	params := make(map[string]interface{})
	params["payment_method"] = strings.TrimSpace(d.PaymentMethod)
	params["order_id"] = strings.TrimSpace(d.OrderID)
	params["order_currency"] = strings.TrimSpace(d.OrderCurrency)
	params["order_amount"] = cast.ToString(d.OrderAmount)
	params["Timestamp"] = cast.ToString(d.Timestamp)
	params["order_name"] = strings.TrimSpace(d.OrderName)
	params["user_id"] = strings.TrimSpace(d.UserID)
	params["return_url"] = strings.TrimSpace(d.ReturnURL)
	params["notify_url"] = strings.TrimSpace(d.NotifyURL)
	params["customer_name"] = strings.TrimSpace(d.CustomerName)
	params["customer_phone"] = strings.TrimSpace(d.CustomerPhone)
	params["customer_email"] = strings.TrimSpace(d.CustomerEmail)
	params["device_info"] = strings.TrimSpace(d.DeviceInfo)
	params["order_note"] = strings.TrimSpace(d.OrderNote)
	params["sign"] = strings.TrimSpace(d.Sign)

	if !goutils.HmacSHA256Verify(params, paySecret) {
		return appError.UnauthenticatedErrCode // 签名错误
	}

	return nil
}

// VerifyPaymentBefore 初步判断是否能够代收
func (s *OrderServer) VerifyPaymentBefore(aspMerchantProject *model.AspMerchantProject, d *req.AspPayment) *appError.Error {

	// 判断商户自定义收款单是否存在
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["mch_id"] = aspMerchantProject.MchId
	orderInfoParam["out_trade_no"] = d.OrderID
	_, err := s.GetOrderInfo(orderInfoParam, nil)
	// 判断是否存在 已经存在则 err == nil 则需要处理
	if err == nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "logic.GetPayoutInfo: ", zap.String("", constant.ConflictExceptionErrMsg))
		return appError.CodeConflictException // 请求有冲突
	}
	return nil
}

func (s *OrderServer) GetOrderInfo(params map[string]interface{}, tx *gorm.DB) (*model.AspOrder, *appError.Error) {
	if tx == nil {
		tx = database.DB
	}
	var aspOrder *model.AspOrder

	for k, v := range params {
		switch k {
		case "id":
			tx = tx.Where("id = ?", v)
		case "sn":
			tx = tx.Where("sn = ?", v)
		case "transaction_id":
			tx = tx.Where("transaction_id = ?", v)
		case "mch_id":
			tx = tx.Where("mch_id = ?", v)
		case "out_trade_no":
			tx = tx.Where("out_trade_no = ?", v)
		}
	}

	if err := tx.First(&aspOrder).Error; err != nil {
		return nil, appError.OrderNotFoundErrCode
	}

	return aspOrder, nil

}

// Insert Get 获取SysApi对象with id
func (s *OrderServer) Insert(d *req.AspPayment, aspMerchantProject *model.AspMerchantProject, aspMerchantProjectConfig *model.AspMerchantProjectConfig, merchantProjectCurrencyInfo *model.AspMerchantProjectCurrency, deptTradeTypeInfo *req.DeptTradeTypeInfo) (*model.AspOrder, *appError.Error) {
	var data model.AspOrder
	// 赋值给 order 数据
	d.Generate(&data, aspMerchantProject, aspMerchantProjectConfig, merchantProjectCurrencyInfo, deptTradeTypeInfo)
	err := database.DB.Create(&data).Error
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspOrder Insert: ", zap.Error(err))
		return nil, appError.CodeInsertErr
	}
	return &data, nil
}

// SolveOrderPaySuccess 处理支付成功的加款等各项操作 更新订单上游返回的数据
// 所有相关订单的更新 都使用这个一个方法来处理
// OrderId 订单id 主键
// params 修改的数组
// updateStructType 修改订单类型
func (s *OrderServer) SolveOrderPaySuccess(OrderId int, params map[string]interface{}, updateStructType string, tx *gorm.DB) (*model.AspOrder, *appError.Error) {
	if tx == nil {
		tx = database.DB
	}
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["id"] = OrderId

	orderInfo, err := s.GetOrderInfo(orderInfoParam, tx)

	if err != nil {
		return nil, err
	}

	b, _ := goutils.JsonEncode(orderInfo)
	logger.ApiInfo(s.LogFileName, s.RequestId, "update before ", zap.String("updateStructType", updateStructType), zap.Any("params", params), zap.String("before", b))

	aspOrder := model.AspOrder{}
	resOrderUpdate := tx.Model(&aspOrder).Where("id = ?", OrderId).Where("trade_state = ? or trade_state = ?", constant.ORDER_TRADE_STATE_PENDING, constant.ORDER_TRADE_STATE_USERPAYING).Updates(params)

	if resOrderUpdate.RowsAffected < 1 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "Update Order error: ", zap.String("err: ", constant.UpdateOrderRowsErrMsg))
		return nil, appError.NewError(constant.ChangeErrMsg).FormatMessage("Order RowsAffected") // 更新代收订单行数 请重新提交
	}

	if resOrderUpdate.Error != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "order.update() ", zap.String("updateStructType", updateStructType), zap.Error(resOrderUpdate.Error))
		return nil, appError.NewError(resOrderUpdate.Error.Error())
	}

	afterOrderInfo, _ := s.GetOrderInfo(orderInfoParam, tx)
	b, _ = goutils.JsonEncode(afterOrderInfo)
	logger.ApiInfo(s.LogFileName, s.RequestId, "update after ", zap.String("updateStructType", updateStructType), zap.String("after", b))
	return afterOrderInfo, nil
}

// DealUnifiedOrderQueryParams 处理请求参数 赋值一些基础值 例如：client_ip
func (s *OrderServer) DealUnifiedOrderQueryParams(d *req.AspOrderQuery, h *req.AspPaymentHeader) *appError.Error {

	if err := checker.Struct(d); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedOrderQueryParams:", zap.Error(err))
		return err
	}

	if err := checker.Struct(h); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedOrderParams:", zap.Error(err))
		return err
	}
	return nil
}

func (s *OrderServer) IsNeedOrderQueryUpstream(orderInfo *model.AspOrder, params map[string]string) bool {
	// 这个参数是强制查询上游
	if _, ok := params["is_call_upstream"]; ok && params["is_call_upstream"] == "yes" {
		return true
	}

	// 判断结构体是否为空
	if orderInfo == nil || orderInfo.Id < 1 {
		return false
	}

	// 判断状态是否是
	status := []string{constant.ORDER_TRADE_STATE_SUCCESS, constant.ORDER_TRADE_STATE_FAILED}

	// strings.ToUpper 将字符串转成大写
	if goutils.InArray(strings.ToUpper(orderInfo.TradeState), status) {
		return false
	}

	return true
}

// GetCurrentTotalFeeOfToday 获取商户当天交易净额 （当天已完成支付成功净额 - 已提现成功的余额）
// 需要参数 mch_id
func (s *OrderServer) GetCurrentTotalFeeOfToday(mchAccountId string) (*int, error) {
	o := database.DB
	var orderAmount *model.AspOrder
	orderAmountList := make([]*req.OrderAmountList, 0)
	// 一天开始的时间戳
	beginTime := goutils.GetTodayBeginTimeStamp()
	// 统计的收款的状态
	tradeState := constant.ORDER_TRADE_STATE_SUCCESS

	// 统计当天收款的成功额度
	err := o.Model(orderAmount).Select("sum(cash_fee) as amount").
		Where("finish_time >= ?", beginTime).
		Where("mch_project_id = ?", mchAccountId).
		Where("trade_state = ?", tradeState).
		Find(&orderAmountList).Error
	if err != nil {
		return nil, err
	}

	// fmt.Println("orderAmountList------------------------", orderAmountList[0].Amount)

	var payoutAmount *model.AspPayout
	payoutAmountList := make([]*req.PayoutAmountList, 0)
	//payoutAmountList[1]= req.PayoutAmountList{}

	// 统计的提现的状态（成功 + 预提现的收款）
	tradeState = constant.PAYOUT_TRADE_STATE_SUCCESS
	// 统计当天提现的成功额度
	if err := database.DB.Model(payoutAmount).Select("sum(cash_fee) as amount").
		Where("finish_time >= ?", beginTime).
		Where("mch_project_id = ?", mchAccountId).
		Where("trade_state = ?", tradeState).
		Find(&payoutAmountList).Error; err != nil {
		return nil, err
	}
	// fmt.Println("payoutAmountList------------------------", payoutAmountList[0].Amount)
	result := orderAmountList[0].Amount - payoutAmountList[0].Amount
	return &result, nil
}

// CheckSignOrderQuery 订单查询 处理验签
func (s *OrderServer) CheckSignOrderQuery(paySecret string, reqPayQuery *req.AspOrderQuery, timestamp int) *appError.Error {
	Signature := s.Head.Signature
	// TODO 需要完善对应的字段 验证签名
	params := make(map[string]interface{})
	params["sn"] = strings.TrimSpace(reqPayQuery.Sn)
	params["Timestamp"] = timestamp
	params["sign"] = strings.TrimSpace(Signature)
	if !goutils.HmacSHA256Verify(params, paySecret) {
		return appError.UnauthenticatedErrCode // 签名错误
	}
	return nil
}

func (s *OrderServer) GetAvailableTotalFeeOfHistory() (int, error) {
	return 0, nil
}

func (s *OrderServer) GetAspOrderList(timeBegin, timeEnd int64) ([]*model.AspOrder, *appError.Error) {
	o := database.DB
	var orderList []*model.AspOrder
	var orderModel model.AspOrder

	err := o.Model(&orderModel).Where("create_time >= ?", timeBegin).
		Where("create_time <= ?", timeEnd).
		//Not(map[string]interface{}{"trade_state": []string{constant.ORDER_TRADE_STATE_SUCCESS, constant.ORDER_TRADE_STATE_FAILED}}).
		Where(map[string]interface{}{"trade_state": []string{constant.ORDER_TRADE_STATE_PENDING}}).
		Order("id desc").
		//Limit(10).
		//Offset(0).
		Find(&orderList).Error

	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "orderList", zap.Error(err))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		//err = (&MissNotFoundErrCode).FormatMessage(constant.MissOrderListNotFoundErrMsg)
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissOrderListNotFoundErrMsg)
	}
	if len(orderList) == 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(orderList) == 0 ", zap.String("error", constant.MissChannelDepartTradeTypeNotFoundErrMsg))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissOrderListNotFoundErrMsg)
	}
	return orderList, nil
}

func (s *OrderServer) GetOrderReturnUrl(orderInfo *model.AspOrder) string {
	returnUrl := orderInfo.ReturnUrl
	if orderInfo.TradeType != "" {
		switch orderInfo.TradeType {
		case constant.TradeTypeAmarquickPay:
			returnUrl = config.AppConfig.Urls.AmarquickpayWappayNotifyUrl
		}
	}
	return returnUrl
}
