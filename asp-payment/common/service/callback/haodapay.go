package callback

import (
	"asp-payment/api-server/req/haodapayreq"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"asp-payment/common/service"
	"asp-payment/common/service/supplier/impl/haodapay"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HaoDaPayCallBackService struct {
	*service.Service
}

func NewHaoDaPayCallBackService(c *fiber.Ctx) *HaoDaPayCallBackService {
	return &HaoDaPayCallBackService{Service: service.NewService(c, constant.HaodaPayCallBackLogFileName)}
}

func (s *HaoDaPayCallBackService) PayOrder() error {

	bodyReq := new(haodapayreq.CallBackReq)

	if errBody := s.C.BodyParser(bodyReq); errBody != nil {
		return s.Error(appError.NewError("Error Params"))
	}

	logger.ApiWarn(s.LogFileName, s.RequestId, "PayOrder ", zap.Any("bodyReq", bodyReq))
	//根据订单号加锁
	lockName := fmt.Sprintf("HaoDaPay:callback:order:lock:%s_%s_%s", constant.ProviderSunny, constant.TradeTypeHaoDaPay, bodyReq.Data.OrderID)
	if !s.Lock(lockName) {
		return s.Error(appError.NewError("Error Params"))
	}
	defer s.UnLock(lockName)
	// 按照文档要求成功返回的数据
	successMap := make(map[string]string)
	successMap["status"] = "success"
	successMap["message"] = "data received"

	//bod, _ := goutils.JsonEncode(bodyReq)
	OrderServer := service.NewPayOrderServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["sn"] = bodyReq.Data.OrderID

	orderInfo, oErr := OrderServer.GetOrderInfo(orderInfoParam, nil)
	if oErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "OrderServer.GetOrderInfo  ", zap.Error(oErr))
		return s.Error(appError.NewError("Error Params"))
	}

	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if goutils.Yuan2Fen(cast.ToFloat64(bodyReq.Data.Amount)) != cast.ToInt64(orderInfo.TotalFee) && cast.ToInt(bodyReq.Data.Amount) != 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "bodyReq.Money != orderInfo.TotalFee ", zap.Any("orderInfo", orderInfo))
		if orderInfo.TradeState != constant.ORDER_TRADE_STATE_PAYERROR {
			failedParams := make(map[string]interface{})
			//failedParams["trade_state"] = constant.ORDER_TRADE_STATE_PAYERROR
			failedParams["supplier_return_code"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
			failedParams["supplier_return_msg"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
			if _, oErr = OrderServer.SolveOrderPaySuccess(orderInfo.Id, failedParams, "haodaPayCallbackUpdate", nil); oErr != nil {
				logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess haodaPayCallbackUpdate", zap.Error(oErr))
				// the server has deserted!
				return s.Error(appError.NewError("the server has deserted!"))
			}
		}
		return s.Error(appError.NewError("Error 2"))
	}

	// 如果已经是成功的状态了 则直接放回 不做操作
	if orderInfo.TradeState == constant.ORDER_TRADE_STATE_SUCCESS || orderInfo.TradeState == constant.ORDER_TRADE_STATE_FAILED {
		logger.ApiInfo(s.LogFileName, s.RequestId, fmt.Sprintf("orderInfo.TradeState == %s", orderInfo.TradeState))
		return s.SuccessJson(successMap)
	}

	merchantProjectServer := service.NewMerchantProjectServer(s.C)
	projectId := cast.ToString(orderInfo.MchProjectId)
	_, pErr := merchantProjectServer.GetMerchantProjectInfo(projectId)
	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetMerchantProjectInfo  ", zap.Error(pErr))
		return s.Error(appError.NewError("Error 3"))
	}

	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(orderInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(orderInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, cErr := service.NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if cErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetChannelDepartInfo  ", zap.Error(cErr))
		return s.Error(appError.NewError("Error 4"))
	}

	// config 转 struct 字符串转 struct
	var channelConfigInfo haodapay.ChannelConfigInfo
	if err := goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "channelDepartInfo.Config to json ", zap.Error(err))
		return s.Error(appError.NewError("Error 5"))
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig channelConfigInfo ", zap.Any("channelConfigInfo", channelConfigInfo))
	var m = make(map[string]interface{})
	m["order_id"] = bodyReq.Data.OrderID
	m["reference"] = bodyReq.Data.Reference
	m["payer_UPIID"] = bodyReq.Data.PayerUPIID
	m["amount"] = bodyReq.Data.Amount
	m["UTR"] = bodyReq.Data.Utr
	fmt.Println("m: ", m)
	signature := haodapay.GetSignature(m, channelConfigInfo.PayinChecksumSecret)
	if signature != bodyReq.Data.Checksum {
		logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("bodyReq.Sign", bodyReq.Data.Checksum))
		return s.Error(appError.NewError("Error 6"))
	}

	upstreamStatus := haodapay.GetPaymentTradeState(bodyReq.Data.Status)
	// fmt.Println("upstreamStatus: ", upstreamStatus)

	if orderInfo.TradeState != upstreamStatus {
		transactionId := bodyReq.Data.Reference
		if upstreamStatus != constant.ORDER_TRADE_STATE_PENDING {
			params := make(map[string]interface{})
			params["transaction_id"] = transactionId
			params["trade_state"] = orderInfo.TradeState // 默认值
			params["finish_time"] = orderInfo.FinishTime
			var isSendQueue bool
			isSendQueue = false
			// 如果是成功 必须提现状态是 顺序的增长的，不能状态逆序
			// 如果上游返回成功 提现状态是 已申请 则修改
			if upstreamStatus == constant.ORDER_TRADE_STATE_SUCCESS {
				MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
				totalFee := cast.ToInt(orderInfo.TotalFee)
				amount := totalFee - orderInfo.ChargeFee - orderInfo.FixedAmount
				// 代收上游成功 待结算余额新增    收支流水记录新增
				orderSuccessUpdate := make(map[string]interface{})
				orderSuccessUpdate["transaction_id"] = transactionId
				orderSuccessUpdate["trade_state"] = upstreamStatus
				orderSuccessUpdate["utr"] = bodyReq.Data.Utr
				orderSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
				errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
					if _, errSuccess := OrderServer.SolveOrderPaySuccess(orderInfo.Id, orderSuccessUpdate, "haodaPayOrderCallbackUpdate", tx); errSuccess != nil {
						logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess haodaPayOrderCallbackUpdate", zap.Error(oErr))
						return appError.CodeUnknown // 服务器开小差了！
					}
					// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
					mErr := MerchantProjectRepository.PayinOrderSuccess(orderInfo.MchProjectId, amount, orderInfo.Id, tx)
					if mErr != nil {
						return appError.NewError(constant.ChangeErrMsg).FormatMessage("HaoDaPayOrderCallbackUpdate") // 更新代付订单错误 请重新提交
					}
					return nil
				})
				if errTrans != nil {
					return s.Error(appError.NewError("Error 7"))
				}
				isSendQueue = true
			}

			if upstreamStatus == constant.ORDER_TRADE_STATE_FAILED {
				params["transaction_id"] = transactionId
				params["trade_state"] = upstreamStatus
				if _, fErr := OrderServer.SolveOrderPaySuccess(orderInfo.Id, params, "HaoDaPayOrderCallbackUpdate", nil); fErr != nil {
					logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess HaoDaPayOrderCallbackUpdate", zap.Error(fErr))
					return s.Error(appError.NewError("Error 8"))
				}
				_ = service.NewSendQueueServer().ManualSendNotifyQueue(orderInfo.Sn)
			}

			// TODO 如果是 已完成状态 则需要添加到 redis 中 通知到下游商户 提现回调
			if isSendQueue == true {
				_ = service.NewSendQueueServer().SendNotifyQueue(orderInfo.Sn)
			}
		}
	}

	return s.SuccessJson(successMap)
}

func (s *HaoDaPayCallBackService) Payout() error {
	// 贯穿支付所需要的参数
	bodyReq := new(haodapayreq.CallBackPayoutReq)

	if err := s.C.BodyParser(bodyReq); err != nil {
		return s.Error(appError.NewError("Error Params"))
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "Payout ", zap.Any("bodyReq", bodyReq))
	//bod, err := goutils.JsonEncode(bodyReq)
	payoutServer := service.NewPayoutServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单 因为传到的是上游渠道id 所以需要加入 provider adapter 条件
	payoutInfoParam := make(map[string]interface{})
	payoutInfoParam["sn"] = bodyReq.Data.Reference
	payoutInfoParam["provider"] = constant.ProviderSunny
	payoutInfoParam["adapter"] = constant.TradeTypeHaoDaPay
	//根据订单号加锁
	lockName := fmt.Sprintf("HaoDaPay:callback:payout:lock:%s", bodyReq.Data.Reference)
	if !s.Lock(lockName) {
		return s.C.SendString("Error")
	}
	defer s.UnLock(lockName)

	// 按照文档要求成功返回的数据
	successMap := make(map[string]string)
	successMap["status"] = "success"
	successMap["message"] = "data received"

	payoutInfo, pErr := payoutServer.GetPayoutInfo(payoutInfoParam, nil)
	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "payoutServer.GetPayoutInfo  ", zap.Error(pErr))
		return s.Error(appError.NewError("Error Params"))
	}
	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if goutils.Yuan2Fen(cast.ToFloat64(bodyReq.Data.Amount)) != cast.ToInt64(payoutInfo.TotalFee) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "bodyReq.Money != orderInfo.TotalFee ", zap.Any("payoutInfo", payoutInfo))
		if payoutInfo.TradeState != constant.PAYOUT_TRADE_STATE_FAILED {
			// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
			payoutFailedUpdate := make(map[string]interface{})
			payoutFailedUpdate["supplier_return_code"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
			payoutFailedUpdate["supplier_return_msg"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
			_, fErr := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "haoDaPayPayoutCallbackUpdate", nil)
			if fErr != nil {
				return s.Error(appError.NewError("Error01")) // 更新商户余额错误 请重新提交
			}
		}
		return s.Error(appError.NewError("Error02"))
	}

	// 如果已经是成功的状态了 则直接放回 不做操作
	if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_SUCCESS || payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_FAILED {
		logger.ApiInfo(s.LogFileName, s.RequestId, fmt.Sprintf("payoutInfo.TradeState == %s", payoutInfo.TradeState))
		return s.SuccessJson(successMap)
	}

	merchantProjectServer := service.NewMerchantProjectServer(s.C)
	appId := cast.ToString(payoutInfo.MchProjectId)
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
		return s.Error(appError.NewError("Error03"))
	}

	// config 转 struct 字符串转 struct
	var channelConfigInfo haodapay.ChannelConfigInfo
	if err := goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "channelDepartInfo.Config to json ", zap.Error(err))
		return s.Error(appError.NewError("Error04"))
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig channelConfigInfo ", zap.Any("channelConfigInfo", channelConfigInfo))

	var m = make(map[string]interface{})
	signature := ""
	if payoutInfo.PayType == constant.PARAMS_PAY_TYPE_BANK {
		m["payout_id"] = bodyReq.Data.PayoutID
		m["reference"] = bodyReq.Data.Reference
		m["beneficiary_account_num"] = bodyReq.Data.BeneficiaryAccountNumber
		m["beneficiary_account_ifsc"] = bodyReq.Data.BeneficiaryAccountIfsc
		m["amount"] = bodyReq.Data.Amount
		m["UTR"] = bodyReq.Data.Utr
		fmt.Println("m: ", m)
		signature = haodapay.GetPayoutSignature(m, channelConfigInfo.PayoutChecksumSecret)
	} else if payoutInfo.PayType == constant.CASHIERDESK_PAY_TYPE_UPI {
		m["payout_id"] = bodyReq.Data.PayoutID
		m["reference"] = bodyReq.Data.Reference
		m["beneficiary_upi_handle"] = bodyReq.Data.BeneficiaryUpiHandle
		m["amount"] = bodyReq.Data.Amount
		m["UTR"] = bodyReq.Data.Utr
		fmt.Println("m: ", m)
		signature = haodapay.GetPayoutUpiSignature(m, channelConfigInfo.PayoutChecksumSecret)
	}

	if signature != bodyReq.Data.Checksum {
		logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("bodyReq.Sign", bodyReq.Data.Checksum))
		return s.Error(appError.NewError("Error05"))
	}
	// 如果当前的提现状态和上游返回的情况不一致的情况
	upstreamStatus := haodapay.GetPayoutCallBackStatus(bodyReq.Status)
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
				payoutSuccessUpdate["bank_utr"] = bodyReq.Data.Utr
				payoutSuccessUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_SUCCESS
				payoutSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
				_, errSuccess := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutSuccessUpdate, "haoDaPayPayoutOrderCallbackUpdate", tx)
				if errSuccess != nil {
					return appError.NewError("Error06") // 更新商户余额错误 请重新提交
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				err := MerchantProjectRepository.PayoutOrderChannelSuccess(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if err != nil {
					return appError.NewError("Error07") // 更新商户余额错误 请重新提交
				}
			}

			// 当失败的时候
			if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
				MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
				// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
				payoutFailedUpdate := make(map[string]interface{})
				payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FAILED
				_, errFailed := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "haoDaPayPayoutOrderCallbackUpdate", tx)
				if errFailed != nil {
					return appError.NewError("Error08") // 更新商户余额错误 请重新提交
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				err := MerchantProjectRepository.PayoutOrderChannelFailed(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if err != nil {
					return appError.NewError("Error09") // 更新商户余额错误 请重新提交
				}
			}
			return nil
		})
		if errTrans != nil {
			return s.Error(appError.NewError(errTrans.Error()))
		}
		// 如果是 已完成状态 则需要添加到 redis 中 通知到下游商户 提现回调
		if isSendQueue == true {
			_ = service.NewSendQueueServer().SendNotifyQueue(payoutInfo.Sn)
		}
		if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
			_ = service.NewSendQueueServer().ManualSendNotifyQueue(payoutInfo.Sn)
		}
	}
	return s.SuccessJson(successMap)
}
