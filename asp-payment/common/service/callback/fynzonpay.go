package callback

import (
	"asp-payment/api-server/req/fynzonpayreq"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"asp-payment/common/service"
	"asp-payment/common/service/supplier/impl/fynzonpay"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FynzonPayCallBackService struct {
	*service.Service
}

func NewFynzonPayCallBackService(c *fiber.Ctx) *FynzonPayCallBackService {
	return &FynzonPayCallBackService{Service: service.NewService(c, constant.FynzonPayCallBackLogFileName)}
}

func (s *FynzonPayCallBackService) PayOrder() error {
	// 返回参数 两种情况
	// 1 transaction_id=8457250803&status_nm=1&status=Approved&price=100.00&curr=INR&id_order=102023022011351800000001&cardtype=&reason=Transaction+Successful+-+Success&fullname=ericluzhong&email=ericluzhong%40gmail.com&address=A97B+North+Block%2CWest+Vinod+Nagar&city=New+Delhi&state=DL&country=IN&phone=9036830689&product_name=Product&amt=100.00&memail=sales%40needsixgaming.com&company=Need+Six+Gaming&bussinessurl=https%3A%2F%2Fneedsixgaming.com%2F&contact_us_url=&customer_service_no=9776361401&tdate=2023-02-20+11%3A37%3A31&descriptor=&zip=110092&authurl=https%3A%2F%2Fbo.fynzonpay.com%2Fauthurl.do%3Ftransaction_id%3D8457250803&callbacks=notify&notify_via=curl_base_notify
	// 2 /notify/order/return_success?status_nm=1&authurl=&status=Approved&amount=100.00&transaction_id=8457250803&descriptor=&tdate=2023-02-20+11%3A37%3A31&curr=INR&reason=Transaction+Successful+-+Success&id_order=102023022011351800000001
	//
	// status_nm=1&authurl=&status=Approved&amount=100.00&transaction_id=8457250803&descriptor=&tdate=2023-02-20+11%3A37%3A31&curr=INR&reason=Transaction+Successful+-+Success&id_order=102023022011351800000001,"requestHeader
	// 贯穿支付所需要的参数
	bodyReq := new(fynzonpayreq.CallBackReq)

	if err := s.C.BodyParser(bodyReq); err != nil {
		return s.C.SendString("Error")
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "PayOrder ", zap.Any("bodyReq", bodyReq))
	//根据订单号加锁
	lockName := fmt.Sprintf("Fynzonpay:callback:order:lock:%s", bodyReq.IDOrder)
	if !s.Lock(lockName) {
		return s.C.SendString("Error")
	}
	defer s.UnLock(lockName)

	//bod, _ := goutils.JsonEncode(bodyReq)
	OrderServer := service.NewPayOrderServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["sn"] = bodyReq.IDOrder

	orderInfo, oErr := OrderServer.GetOrderInfo(orderInfoParam, nil)
	if oErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "OrderServer.GetOrderInfo  ", zap.Error(oErr))
		return s.C.SendString("Error")
	}

	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if goutils.Yuan2Fen(cast.ToFloat64(bodyReq.Amount)) != cast.ToInt64(orderInfo.TotalFee) && cast.ToInt(bodyReq.Amount) != 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "bodyReq.Money != orderInfo.TotalFee ", zap.Any("orderInfo", orderInfo))
		if orderInfo.TradeState != constant.ORDER_TRADE_STATE_PAYERROR {
			failedParams := make(map[string]interface{})
			//failedParams["trade_state"] = constant.ORDER_TRADE_STATE_PAYERROR
			failedParams["supplier_return_code"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
			failedParams["supplier_return_msg"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
			if _, oErr = OrderServer.SolveOrderPaySuccess(orderInfo.Id, failedParams, "fynzonPayCallbackUpdate", nil); oErr != nil {
				logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess fynzonPayCallbackUpdate", zap.Error(oErr))
				return appError.CodeUnknown // 服务器开小差了！
			}
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
	//var m map[string]interface{}
	//
	//_ = goutils.JsonDecode(bod, &m)
	//signature := fynzon.GetSignature(m, channelConfigInfo.Signature)
	//if signature != bodyReq.Sign {
	//	logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("bodyReq.Sign", bodyReq.Sign))
	//	return s.C.SendString("Error")
	//}

	upstreamStatus := fynzonpay.GetPaymentTradeState(bodyReq.StatusNm)
	// fmt.Println("upstreamStatus: ", upstreamStatus)

	if orderInfo.TradeState != upstreamStatus {
		if upstreamStatus != constant.ORDER_TRADE_STATE_PENDING {
			orderInfo.TransactionId = bodyReq.TransactionID
			// 执行查询逻辑 判断金额和状态是否一致
			scanQueryData, scErr := service.NewPayOrderServer(s.C).DoQueryUpstream(orderInfo, channelDepartInfo)
			if scErr != nil {
				s.C.SendString("Error")
			}
			if scanQueryData.TradeState != upstreamStatus || cast.ToUint(scanQueryData.CashFee) != orderInfo.TotalFee {
				logger.ApiError(s.LogFileName, s.RequestId, "Second diff OrderInfo Err", zap.Any("scanQueryData", scanQueryData))
				return s.C.SendString("Error")
			}
			params := make(map[string]interface{})
			params["transaction_id"] = bodyReq.TransactionID // 默认值
			params["trade_state"] = orderInfo.TradeState     // 默认值
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
				orderSuccessUpdate["transaction_id"] = bodyReq.TransactionID
				orderSuccessUpdate["trade_state"] = upstreamStatus
				orderSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
				errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
					if _, errSuccess := OrderServer.SolveOrderPaySuccess(orderInfo.Id, orderSuccessUpdate, "fynzonPayOrderCallbackUpdate", tx); errSuccess != nil {
						logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess fynzonPayOrderCallbackUpdate", zap.Error(oErr))
						return appError.CodeUnknown // 服务器开小差了！
					}
					// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
					mErr := MerchantProjectRepository.PayinOrderSuccess(orderInfo.MchProjectId, amount, orderInfo.Id, tx)
					if mErr != nil {
						return appError.NewError(constant.ChangeErrMsg).FormatMessage("FynzonPayOrderCallbackUpdate") // 更新代付订单错误 请重新提交
					}
					return nil
				})
				if errTrans != nil {
					return s.C.SendString(errTrans.Error())
				}
				isSendQueue = true
			}

			if upstreamStatus == constant.ORDER_TRADE_STATE_FAILED {
				params["transaction_id"] = bodyReq.TransactionID
				params["trade_state"] = upstreamStatus
				if _, fErr := OrderServer.SolveOrderPaySuccess(orderInfo.Id, params, "FynzonPayOrderCallbackUpdate", nil); fErr != nil {
					logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess FynzonPayOrderCallbackUpdate", zap.Error(fErr))
					return s.C.SendString("1")
				}
				_ = service.NewSendQueueServer().ManualSendNotifyQueue(orderInfo.Sn)
			}

			// TODO 如果是 已完成状态 则需要添加到 redis 中 通知到下游商户 提现回调
			if isSendQueue == true {
				_ = service.NewSendQueueServer().SendNotifyQueue(orderInfo.Sn)
			}
		}
	}

	return s.C.SendString("OK")
}

// Beneficiary 添加受益人回调
func (s *FynzonPayCallBackService) Beneficiary() error {

	// 贯穿支付所需要的参数
	bodyReq := new(fynzonpayreq.CallBackBeneficiaryReq)
	if err := s.C.BodyParser(bodyReq); err != nil {
		return s.C.SendString("Error")
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "Beneficiary ", zap.Any("bodyReq", bodyReq))
	//根据订单号加锁
	lockName := fmt.Sprintf("Fynzonpay:callback:order:lock:%s_%s_%s", constant.ProviderSunny, constant.TradeTypeFynzonPay, bodyReq.BeneId)
	if !s.Lock(lockName) {
		return s.C.SendString("Error")
	}
	defer s.UnLock(lockName)

	if bodyReq.Status != "0000" || bodyReq.BeneId == "" {
		return s.C.SendString("OK")
	}
	payoutServer := service.NewPayoutServer(s.C)
	// 1. 查询到受益人
	beneficiaryInfoParam := make(map[string]interface{})
	beneficiaryInfoParam["provider"] = constant.ProviderSunny
	beneficiaryInfoParam["adapter"] = constant.TradeTypeFynzonPay
	beneficiaryInfoParam["trade_type"] = constant.TradeType_PAYOUT
	beneficiaryInfoParam["benefiary_id"] = bodyReq.BeneId
	BenefiaryInfo, beErr := payoutServer.BeneficiaryInfo(beneficiaryInfoParam)
	// 2. 更新受益人状态
	if beErr == nil && BenefiaryInfo.TradeState != constant.BENEFICIARY_TRADE_STATE_SUCCESS {
		beneficiaryUpdateParam := make(map[string]interface{})
		beneficiaryUpdateParam["trade_state"] = constant.BENEFICIARY_TRADE_STATE_SUCCESS
		_, _ = payoutServer.SolveBeneficiarySuccess(BenefiaryInfo.Id, beneficiaryUpdateParam, "fynzonpayCallBackUpdate", nil)
	}
	return s.C.SendString("OK")
}
func (s *FynzonPayCallBackService) Payout() error {
	// 贯穿支付所需要的参数
	bodyReq := new(fynzonpayreq.CallBackPayoutReq)

	if err := s.C.BodyParser(bodyReq); err != nil {
		return s.C.SendString("Error")
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "Payout ", zap.Any("bodyReq", bodyReq))
	//bod, err := goutils.JsonEncode(bodyReq)
	payoutServer := service.NewPayoutServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单 因为传到的是上游渠道id 所以需要加入 provider adapter 条件
	payoutInfoParam := make(map[string]interface{})
	payoutInfoParam["transaction_id"] = bodyReq.TransactionID
	payoutInfoParam["provider"] = constant.ProviderSunny
	payoutInfoParam["adapter"] = constant.TradeTypeFynzonPay
	//根据订单号加锁
	lockName := fmt.Sprintf("Fynzonpay:callback:payout:lock:%s", bodyReq.TransactionID)
	if !s.Lock(lockName) {
		return s.C.SendString("Error")
	}
	defer s.UnLock(lockName)

	payoutInfo, pErr := payoutServer.GetPayoutInfo(payoutInfoParam, nil)
	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "payoutServer.GetPayoutInfo  ", zap.Error(pErr))
		return s.C.SendString("Error")
	}
	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if goutils.Yuan2Fen(cast.ToFloat64(bodyReq.PayoutAmount)) != cast.ToInt64(payoutInfo.TotalFee) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "bodyReq.Money != orderInfo.TotalFee ", zap.Any("payoutInfo", payoutInfo))
		if payoutInfo.TradeState != constant.PAYOUT_TRADE_STATE_FAILED {
			// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
			payoutFailedUpdate := make(map[string]interface{})
			payoutFailedUpdate["supplier_return_code"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
			payoutFailedUpdate["supplier_return_msg"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
			_, fErr := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "fynzonPayPayoutCallbackUpdate", nil)
			if fErr != nil {
				return s.C.SendString("fynzonPayPayoutCallbackUpdate Error") // 更新商户余额错误 请重新提交
			}
		}
		return s.C.SendString("Error")
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
	if err := goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "channelDepartInfo.Config to json ", zap.Error(err))
		return s.C.SendString("Error")
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig channelConfigInfo ", zap.Any("channelConfigInfo", channelConfigInfo))

	//var m map[string]interface{}
	//_ = goutils.JsonDecode(bod, &m)
	//logger.ApiInfo(s.LogFileName, s.RequestId, "json.Unmarshal ", zap.Any("PartnerId", m))
	//signature := Fynzon.GetSignature(m, channelConfigInfo.Signature)
	//if signature != bodyReq.Sign {
	//	logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("reqHeader.Signature", bodyReq.Sign))
	//	return s.C.SendString("Error7")
	//}
	// 如果当前的提现状态和上游返回的情况不一致的情况
	upstreamStatus := fynzonpay.GetPayoutCallBackStatus(bodyReq.TransactionStatus)
	if payoutInfo.TradeState != upstreamStatus {

		scanQueryData, qErr := payoutServer.RequestUpstreamQueryPayout(payoutInfo, channelDepartInfo)
		if qErr != nil {
			return s.C.SendString(qErr.Error())
		}
		if scanQueryData.TradeState != upstreamStatus {
			return s.C.SendString("Error10")
		}
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
				_, errSuccess := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutSuccessUpdate, "fynzonPayPayoutOrderCallbackUpdate", tx)
				if errSuccess != nil {
					return appError.NewError("Error8") // 更新商户余额错误 请重新提交
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				err := MerchantProjectRepository.PayoutOrderChannelSuccess(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if err != nil {
					return appError.NewError("Error9") // 更新商户余额错误 请重新提交
				}
			}

			// 当失败的时候
			if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
				MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
				// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
				payoutFailedUpdate := make(map[string]interface{})
				payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FAILED
				_, errFailed := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "fynzonPayPayoutOrderCallbackUpdate", tx)
				if errFailed != nil {
					return appError.NewError("Error10") // 更新商户余额错误 请重新提交
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				err := MerchantProjectRepository.PayoutOrderChannelFailed(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if err != nil {
					return appError.NewError("Error11") // 更新商户余额错误 请重新提交
				}
			}
			return nil
		})
		if errTrans != nil {
			return s.C.SendString(errTrans.Error())
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
