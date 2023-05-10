package callback

import (
	"asp-payment/api-server/req/seveneightpayreq"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"asp-payment/common/service"
	"asp-payment/common/service/supplier/impl/amarquickpay"
	"asp-payment/common/service/supplier/impl/seveneight"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type AmarquickPayCallBackService struct {
	*service.Service
}

func NewAmarquickPayCallBackService(c *fiber.Ctx) *AmarquickPayCallBackService {
	return &AmarquickPayCallBackService{Service: service.NewService(c, constant.AmarquickPayCallBackLogFileName)}
}

func (s *AmarquickPayCallBackService) PayOrder() error {
	// 贯穿支付所需要的参数
	status := s.C.FormValue("STATUS")
	responseDateTime := s.C.FormValue("RESPONSE_DATE_TIME")
	responseCode := s.C.FormValue("RESPONSE_CODE")
	appId := s.C.FormValue("APP_ID")
	txnId := s.C.FormValue("TXN_ID")
	txnType := s.C.FormValue("TXNTYPE")
	returnUrl := s.C.FormValue("RETURN_URL")
	orderId := s.C.FormValue("ORDER_ID")
	txnKey := s.C.FormValue("TXN_KEY", "")
	acqId := s.C.FormValue("ACQ_ID", "")
	cardMask := s.C.FormValue("CARD_MASK", "")
	duplicateYn := s.C.FormValue("DUPLICATE_YN", "")
	mopType := s.C.FormValue("MOP_TYPE", "")
	paymentType := s.C.FormValue("PAYMENT_TYPE", "")
	responseMessage := s.C.FormValue("RESPONSE_MESSAGE", "")
	custPhone := s.C.FormValue("CUST_PHONE", "")
	custName := s.C.FormValue("CUST_NAME", "")
	custEmail := s.C.FormValue("CUST_EMAIL", "")
	currencyCode := s.C.FormValue("CURRENCY_CODE", "")
	amount := s.C.FormValue("AMOUNT", "")
	rrn := s.C.FormValue("RRN", "")
	origTxnId := s.C.FormValue("ORIG_TXN_ID", "")
	authCode := s.C.FormValue("AUTH_CODE", "")
	hash := s.C.FormValue("HASH", "")
	dateTm := strings.Replace(responseDateTime, "+", " ", -1)

	params := make(map[string]interface{})
	params["RESPONSE_DATE_TIME"] = dateTm
	params["RESPONSE_CODE"] = responseCode
	params["STATUS"] = status
	params["APP_ID"] = appId
	params["TXN_ID"] = txnId
	params["TXNTYPE"] = txnType
	params["RETURN_URL"] = returnUrl
	params["ORDER_ID"] = orderId
	if acqId != "" {
		params["ACQ_ID"] = acqId
	}
	if txnKey != "" {
		params["TXN_KEY"] = txnKey
	}
	if cardMask != "" {
		params["CARD_MASK"] = cardMask
	}
	if duplicateYn != "" {
		params["DUPLICATE_YN"] = duplicateYn
	}
	if mopType != "" {
		params["MOP_TYPE"] = mopType
	}
	if paymentType != "" {
		params["PAYMENT_TYPE"] = paymentType
	}
	if responseMessage != "" {
		params["RESPONSE_MESSAGE"] = responseMessage
	}
	if custPhone != "" {
		params["CUST_PHONE"] = custPhone
	}
	if custName != "" {
		params["CUST_NAME"] = custName
	}
	if custEmail != "" {
		params["CUST_EMAIL"] = custEmail
	}
	if currencyCode != "" {
		params["CURRENCY_CODE"] = currencyCode
	}
	if amount != "" {
		params["AMOUNT"] = amount
	}
	if rrn != "" {
		params["RRN"] = rrn
	}
	if origTxnId != "" {
		params["ORIG_TXN_ID"] = origTxnId
	}
	if authCode != "" {
		params["AUTH_CODE"] = authCode
	}

	logger.ApiInfo(s.LogFileName, s.RequestId, fmt.Sprintf("amarquickpay params:%+v", params))
	//根据订单号加锁
	lockName := fmt.Sprintf("amarquickpay:callback:order:lock:%s", orderId)
	if !s.Lock(lockName) {
		return s.C.SendString("Error")
	}
	defer s.UnLock(lockName)

	OrderServer := service.NewPayOrderServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["sn"] = orderId

	orderInfo, oErr := OrderServer.GetOrderInfo(orderInfoParam, nil)
	if oErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "OrderServer.GetOrderInfo  ", zap.Error(oErr))
		return s.C.SendString("Error")
	}

	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if cast.ToUint(amount) != orderInfo.TotalFee {
		logger.ApiWarn(s.LogFileName, s.RequestId, "amount != orderInfo.TotalFee ", zap.Any("orderInfo", orderInfo))
		failedParams := make(map[string]interface{})
		failedParams["supplier_return_code"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		failedParams["supplier_return_msg"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
		if _, oErr = OrderServer.SolveOrderPaySuccess(orderInfo.Id, failedParams, "firstPayCallbackUpdate", nil); oErr != nil {
			logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess firstPayCallbackUpdate", zap.Error(oErr))
			return appError.CodeUnknown // 服务器开小差了！
		}
		return s.C.SendString("Error")
	}

	// 如果已经是成功的状态了 则直接放回 不做操作
	if orderInfo.TradeState == constant.ORDER_TRADE_STATE_SUCCESS || orderInfo.TradeState == constant.ORDER_TRADE_STATE_FAILED {
		logger.ApiInfo(s.LogFileName, s.RequestId, fmt.Sprintf("orderInfo.TradeState == %s", orderInfo.TradeState))
		return s.C.SendString("OK")
	}

	merchantProjectServer := service.NewMerchantProjectServer(s.C)
	projectId := cast.ToString(orderInfo.MchProjectId)
	_, pErr := merchantProjectServer.GetMerchantProjectInfo(projectId)
	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetMerchantProjectInfo  ", zap.Error(pErr))
		return s.C.SendString("Error")
	}

	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(orderInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(orderInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, cErr := service.NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if cErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetChannelDepartInfo  ", zap.Error(cErr))
		return s.C.SendString("Error")
	}

	// config 转 struct 字符串转 struct
	var channelConfigInfo model.AspChannelDepartConfigInfo
	if err := goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "channelDepartInfo.Config to json ", zap.Error(err))
		return s.C.SendString("Error")
	}

	logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig channelConfigInfo ", zap.String("PartnerId", channelConfigInfo.PartnerId), zap.String("Signature", channelConfigInfo.Signature))

	signature := amarquickpay.GetSignature(params, channelConfigInfo.Signature)
	if signature != hash {
		logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("bodyReq.Sign", hash))
		return s.C.SendString("Error")
	}

	if status == "Captured" || status == "Approved" {
		//修改订单状态
		MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
		totalFee := cast.ToInt(orderInfo.TotalFee)
		addAmount := totalFee - orderInfo.ChargeFee - orderInfo.FixedAmount
		// 代收上游成功 待结算余额新增    收支流水记录新增
		orderSuccessUpdate := make(map[string]interface{})
		orderSuccessUpdate["trade_state"] = constant.ORDER_TRADE_STATE_SUCCESS
		orderSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
		orderSuccessUpdate["transaction_id"] = txnId

		errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
			if _, oErr = OrderServer.SolveOrderPaySuccess(orderInfo.Id, orderSuccessUpdate, "firstPayCallbackUpdate", tx); oErr != nil {
				logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess firstPayCallbackUpdate", zap.Error(oErr))
				return appError.CodeUnknown // 服务器开小差了！
			}
			// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
			errChange := MerchantProjectRepository.PayinOrderSuccess(orderInfo.MchProjectId, addAmount, orderInfo.Id, tx)
			if errChange != nil {
				return appError.NewError(constant.ChangeErrMsg).FormatMessage("seveneightPayOrderCallbackUpdate") // 更新代付订单错误 请重新提交
			}
			return nil
		})
		if errTrans != nil {
			return s.C.SendString("fail")
		}
		_ = service.NewSendQueueServer().SendNotifyQueue(orderInfo.Sn)

		return s.C.SendString("OK")
	}

	return s.C.SendString("fail")
}

func (s *AmarquickPayCallBackService) Payout() error {
	// 贯穿支付所需要的参数
	bodyReq := new(seveneightpayreq.CallBackPayoutReq)

	if err := s.C.BodyParser(bodyReq); err != nil {
		return s.C.SendString("Error")
	}
	bod, err := goutils.JsonEncode(bodyReq)
	payoutServer := service.NewPayoutServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单
	payoutInfoParam := make(map[string]interface{})
	payoutInfoParam["sn"] = bodyReq.OutTradeNo
	//根据订单号加锁
	lockName := fmt.Sprintf("seveneightpay:callback:payout:lock:%s", bodyReq.OutTradeNo)
	if !s.Lock(lockName) {
		return s.C.SendString("Error")
	}
	defer s.UnLock(lockName)
	fmt.Println("payoutInfoParam:", payoutInfoParam)
	payoutInfo, pErr := payoutServer.GetPayoutInfo(payoutInfoParam, nil)
	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "payoutServer.GetPayoutInfo  ", zap.Error(pErr))
		return s.C.SendString("Error")
	}

	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if cast.ToUint(bodyReq.Money) != payoutInfo.TotalFee {
		logger.ApiWarn(s.LogFileName, s.RequestId, "bodyReq.Money != orderInfo.TotalFee ", zap.Any("payoutInfo", payoutInfo))

		// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
		payoutFailedUpdate := make(map[string]interface{})
		payoutFailedUpdate["supplier_return_code"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		payoutFailedUpdate["supplier_return_msg"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE

		_, fErr := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "amarquickpayorderUpdate", nil)
		if fErr != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
		}
		return s.C.SendString("1")
	}

	// 如果已经是成功的状态了 则直接放回 不做操作
	if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_SUCCESS || payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_FAILED {
		logger.ApiInfo(s.LogFileName, s.RequestId, fmt.Sprintf("payoutInfo.TradeState == %s", payoutInfo.TradeState))
		return s.C.SendString("OK")
	}

	merchantProjectServer := service.NewMerchantProjectServer(s.C)
	appId := cast.ToString(payoutInfo.MchProjectId)
	// appError.MissMerchantProjectNotFoundErr.FormatMessage("xxxx")
	_, pEr := merchantProjectServer.GetMerchantProjectInfo(appId)
	if pEr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetMerchantProjectInfo  ", zap.Error(pEr))
		return s.C.SendString("Error")
	}

	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(payoutInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(payoutInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, cErr := service.NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if cErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetChannelDepartInfo  ", zap.Error(cErr))
		return s.C.SendString("Error")
	}

	// config 转 struct 字符串转 struct
	var channelConfigInfo model.AspChannelDepartConfigInfo
	if err = goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "channelDepartInfo.Config to json ", zap.Error(err))
		return s.C.SendString("Error")
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig channelConfigInfo ", zap.Any("channelConfigInfo", channelConfigInfo))

	var m map[string]interface{}
	_ = goutils.JsonDecode(bod, &m)
	logger.ApiInfo(s.LogFileName, s.RequestId, "json.Unmarshal ", zap.Any("PartnerId", m))
	signature := seveneight.GetSignature(m, channelConfigInfo.Signature)
	if signature != bodyReq.Sign {
		logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("reqHeader.Signature", bodyReq.Sign))
		return s.C.SendString("Error7")
	}

	// 如果当前的提现状态和上游返回的情况不一致的情况
	status := bodyReq.Status
	upstreamStatus := seveneight.GetPayoutCallBackStatus(status)
	if payoutInfo.TradeState != upstreamStatus {
		var isSendQueue bool
		isSendQueue = false

		totalFee := cast.ToInt(payoutInfo.TotalFee)
		changeAmount := totalFee + payoutInfo.ChargeFee + payoutInfo.FixedAmount
		// 如果是成功 必须提现状态是 顺序的增长的，不能状态逆序
		// 如果上游返回成功 提现状态是 已申请 则修改
		if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
			isSendQueue = true
			// 如果上游返回失败 提现状态是 已申请 则修改
		}

		errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
			// 当成功时候
			if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
				MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
				// 代付上游成功 解冻+代付成功  冻结金额减去  预扣金额记录新增  收支流水扣减
				payoutSuccessUpdate := make(map[string]interface{})
				payoutSuccessUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_SUCCESS
				payoutSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
				_, errSuccess := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutSuccessUpdate, "amarquickpayorderUpdate", tx)
				if errSuccess != nil {
					return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				errChange := MerchantProjectRepository.PayoutOrderChannelSuccess(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if errChange != nil {
					return appError.NewError("Error8") // 更新商户余额错误 请重新提交
				}
			}

			// 当失败的时候
			if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
				MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)

				// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
				payoutFailedUpdate := make(map[string]interface{})
				payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FAILED
				_, errFailed := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "amarquickpayorderUpdate", tx)
				if errFailed != nil {
					return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				errChange := MerchantProjectRepository.PayoutOrderChannelFailed(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if errChange != nil {
					return appError.NewError("Error9") // 更新商户余额错误 请重新提交
				}
			}
			return nil
		})
		if errTrans != nil {
			return s.C.SendString("Error10")
		}
		// 如果是 已完成状态 则需要添加到 redis 中 通知到下游商户 提现回调
		if isSendQueue == true {
			_ = service.NewSendQueueServer().SendNotifyQueue(payoutInfo.Sn)
		}
		if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
			_ = service.NewSendQueueServer().ManualSendNotifyQueue(payoutInfo.Sn)
		}
	}
	return s.C.SendString("OK")
}
