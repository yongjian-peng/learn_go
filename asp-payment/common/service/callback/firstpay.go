package callback

import (
	"asp-payment/api-server/req/firstpayreq"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"asp-payment/common/service"
	"asp-payment/common/service/supplier/impl/firstpayimpl"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FirstPayCallBackService struct {
	*service.Service
}

func NewFirstPayCallBackServer(c *fiber.Ctx) *FirstPayCallBackService {
	return &FirstPayCallBackService{Service: service.NewService(c, constant.FirstPayCallBackLogFileName)}
}

func (s *FirstPayCallBackService) PayOrder() error {
	// 贯穿支付所需要的参数
	bodyReq := new(firstpayreq.FirstPayCallBackReq)

	if err := s.C.BodyParser(bodyReq); err != nil {
		return err
	}
	bod, err := goutils.JsonEncode(bodyReq)
	OrderServer := service.NewPayOrderServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["sn"] = bodyReq.AppOrderId
	orderInfo, oErr := OrderServer.GetOrderInfo(orderInfoParam, nil)
	if oErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "logic.GetOrderInfo  ", zap.Error(oErr))
		return appError.OrderNotFoundErrCode // 订单不存在
	}

	// 如果已经是成功的状态了 则直接放回 不做操作
	if orderInfo.TradeState == constant.ORDER_TRADE_STATE_SUCCESS || orderInfo.TradeState == constant.ORDER_TRADE_STATE_FAILED {
		logger.ApiInfo(s.LogFileName, s.RequestId, fmt.Sprintf("orderInfo.TradeState == %s", orderInfo.TradeState))
		return s.Success("SUCCESS")
	}

	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if goutils.Yuan2Fen(cast.ToFloat64(bodyReq.Amount)) != cast.ToInt64(orderInfo.TotalFee) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "bodyReq.Amount != orderInfo.TotalFee ", zap.Any("orderInfo", orderInfo))
		failedParams := make(map[string]interface{})
		failedParams["supplier_return_code"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		failedParams["supplier_return_msg"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
		if _, oErr = OrderServer.SolveOrderPaySuccess(orderInfo.Id, failedParams, "firstPayCallbackUpdate", nil); oErr != nil {
			logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess firstPayCallbackUpdate", zap.Error(oErr))
			return appError.CodeUnknown // 服务器开小差了！
		}
		return appError.CodeSupplierInternalChannelParamsFailedErrCode // 道内部参数错误，失败
	}

	merchantProjectServer := service.NewMerchantProjectServer(s.C)
	appId := cast.ToString(orderInfo.MchProjectId)
	// appError.MissMerchantProjectNotFoundErr.FormatMessage("xxxx")
	_, pErr := merchantProjectServer.GetMerchantProjectInfo(appId)
	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "merchantProjectServer.GetMerchantProjectInfo  ", zap.Error(pErr))
		return err
	}

	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(orderInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(orderInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, cErr := service.NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if cErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetChannelDepartInfo  ", zap.Error(cErr))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartNotFoundErrMsg) // missing channel depart not found
	}

	if channelDepartInfo.Config == "" {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetChannelDepartInfo  config null")
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartNotFoundErrMsg) // missing channel depart not found
	}

	// config 转 struct 字符串转 struct
	var channelConfigInfo model.AspChannelDepartConfigInfo
	if err = goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig to json ", zap.Error(err))
		return err
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig channelConfigInfo ", zap.String("Appid", channelConfigInfo.Appid), zap.String("Signature", channelConfigInfo.Signature))

	signature := firstpayimpl.VerifySignature(channelConfigInfo.Signature, bod)

	if signature != s.Head.Signature {
		logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("reqHeader.Signature", s.Head.Signature))
		return appError.UnauthenticatedErrCode // 签名错误
	}

	// 4. 判断订单的状态，修改订单状态 以及相关信息
	// 如果当前的提现状态和上游返回的情况不一致的情况
	upstreamStatus := firstpayimpl.GetFirstPayPaymentStatus(bodyReq.Status)
	if orderInfo.TradeState != upstreamStatus {
		// 如果存在 代付中的情况则不做操作
		params := make(map[string]interface{})
		params["trade_state"] = orderInfo.TradeState // 默认值
		params["finish_time"] = cast.ToString(orderInfo.FinishTime)
		var isSendQueue bool
		isSendQueue = false
		// 如果成功 必须提现状态是 顺序的增长的，不能状态逆序
		// 如果上游返回成功 提现状态是 已申请 则修改
		if upstreamStatus == constant.ORDER_TRADE_STATE_SUCCESS {
			MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
			totalFee := cast.ToInt(orderInfo.TotalFee)
			amount := totalFee - orderInfo.ChargeFee - orderInfo.FixedAmount
			// 代收上游成功 待结算余额新增    收支流水记录新增
			orderSuccessUpdate := make(map[string]interface{})
			orderSuccessUpdate["trade_state"] = upstreamStatus
			orderSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
			errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
				if _, errSuccess := OrderServer.SolveOrderPaySuccess(orderInfo.Id, orderSuccessUpdate, "firstPayCallbackUpdate", tx); errSuccess != nil {
					logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess firstPayCallbackUpdate", zap.Error(errSuccess))
					return appError.CodeUnknown // 服务器开小差了！
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				mErr := MerchantProjectRepository.PayinOrderSuccess(orderInfo.MchProjectId, amount, orderInfo.Id, tx)
				if mErr != nil {
					return appError.NewError(constant.ChangeErrMsg).FormatMessage("firstPayOrderCallbackUpdate") // 更新代付订单错误 请重新提交
				}
				return nil
			})
			if errTrans != nil {
				return s.Error(appError.NewError(errTrans.Error()))
			}

			isSendQueue = true
			// 如果上游返回失败 提现状态是 已申请 则修改
		}
		if upstreamStatus == constant.ORDER_TRADE_STATE_FAILED {
			params["trade_state"] = upstreamStatus
			if _, oErr = OrderServer.SolveOrderPaySuccess(orderInfo.Id, params, "firstPayCallbackUpdate", nil); oErr != nil {
				logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess firstPayCallbackUpdate", zap.Error(oErr))
				return appError.CodeUnknown // 服务器开小差了！
			}
			_ = service.NewSendQueueServer().ManualSendNotifyQueue(orderInfo.Sn)
		}
		// TODO 如果是 已完成状态 则需要添加到 redis 中 通知到下游商户 提现回调
		if isSendQueue == true {
			_ = service.NewSendQueueServer().SendNotifyQueue(orderInfo.Sn)
		}
	}

	return s.Success("SUCCESS")
}

func (s *FirstPayCallBackService) Payout() error {

	// 贯穿支付所需要的参数
	bodyReq := new(firstpayreq.FirstPayPayoutCallBackReq)
	if err := s.C.BodyParser(bodyReq); err != nil {
		return err
	}
	bod, err := goutils.JsonEncode(bodyReq)
	payoutServer := service.NewPayoutServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单
	payoutInfoParam := make(map[string]interface{})
	payoutInfoParam["sn"] = bodyReq.AppOrderId

	payoutInfo, pErr := payoutServer.GetPayoutInfo(payoutInfoParam, nil)
	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "payoutServer.GetPayoutInfo  ", zap.Error(pErr))
		return appError.OrderNotFoundErrCode // 订单不存在
	}
	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if goutils.Yuan2Fen(cast.ToFloat64(bodyReq.Amount)) != cast.ToInt64(payoutInfo.TotalFee) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "bodyReq.Amount != payoutInfo.TotalFee ", zap.Any("payoutInfo", payoutInfo))
		payoutFailedUpdate := make(map[string]interface{})
		payoutFailedUpdate["supplier_return_code"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		payoutFailedUpdate["supplier_return_msg"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE

		_, fErr := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "firstpayOrderUpdate", nil)
		if fErr != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
		}
		return appError.CodeSupplierInternalChannelParamsFailedErrCode // 道内部参数错误，失败
	}

	// 如果已经是成功的状态了 则直接放回 不做操作
	if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_SUCCESS || payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_FAILED {
		logger.ApiInfo(s.LogFileName, s.RequestId, fmt.Sprintf("payoutInfo.TradeState == %s", payoutInfo.TradeState))
		return s.Success("SUCCESS")
	}

	merchantProjectServer := service.NewMerchantProjectServer(s.C)
	appId := cast.ToString(payoutInfo.MchProjectId)
	// appError.MissMerchantProjectNotFoundErr.FormatMessage("xxxx")
	_, err = merchantProjectServer.GetMerchantProjectInfo(appId)
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetMerchantProjectInfo  ", zap.Error(err))
		return err
	}

	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(payoutInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(payoutInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, err := service.NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetChannelDepartInfo  ", zap.Error(err))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartNotFoundErrMsg) // missing channel depart not found
	}

	if channelDepartInfo.Config == "" {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetChannelDepartInfo  config null")
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartNotFoundErrMsg) // missing channel depart not found
	}

	// config 转 struct 字符串转 struct
	var channelConfigInfo model.AspChannelDepartConfigInfo
	if err = goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig to json ", zap.Error(err))
		return err
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig channelConfigInfo ", zap.String("Appid", channelConfigInfo.Appid), zap.String("Signature", channelConfigInfo.Signature))

	signature := firstpayimpl.VerifySignature(channelConfigInfo.Signature, bod)

	if signature != s.Head.Signature {
		logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("reqHeader.Signature", s.Head.Signature))
		return appError.UnauthenticatedErrCode // 签名错误
	}

	// 4. 判断订单的状态，修改订单状态 以及相关信息
	// 如果当前的提现状态和上游返回的情况不一致的情况
	status := cast.ToInt(bodyReq.Status)
	upstreamStatus := firstpayimpl.GetFirstPayPaymentStatus(status)
	if payoutInfo.TradeState != upstreamStatus {
		params := make(map[string]interface{})
		params["trade_state"] = payoutInfo.TradeState // 默认值
		params["finish_time"] = payoutInfo.FinishTime
		var isSendQueue bool
		isSendQueue = false
		// 如果是成功 必须提现状态是 顺序的增长的，不能状态逆序
		// 如果上游返回成功 提现状态是 已申请 则修改
		if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
			params["trade_state"] = upstreamStatus
			isSendQueue = true
			// 如果上游返回失败 提现状态是 已申请 则修改
		} else if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
			params["trade_state"] = upstreamStatus
		}

		totalFee := cast.ToInt(payoutInfo.TotalFee)
		changeAmount := totalFee + payoutInfo.ChargeFee + payoutInfo.FixedAmount

		errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
			// 当成功时候
			if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
				MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)

				// 代付上游成功 解冻+代付成功  冻结金额减去  预扣金额记录新增  收支流水扣减
				payoutSuccessUpdate := make(map[string]interface{})
				payoutSuccessUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_SUCCESS
				payoutSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
				_, errSuccess := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutSuccessUpdate, "firstpayOrderUpdate", tx)
				if errSuccess != nil {
					return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
				}

				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				errChange := MerchantProjectRepository.PayoutOrderChannelSuccess(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if errChange != nil {
					return appError.NewError(constant.ChangeErrMsg).FormatMessage("firstPayPayoutCallbackUpdate") // 更新代付订单错误 请重新提交
				}
			}

			// 当失败的时候
			if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
				MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
				// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
				payoutFailedUpdate := make(map[string]interface{})
				payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FAILED
				_, errFailed := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "firstpayOrderUpdate", tx)
				if errFailed != nil {
					return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
				}

				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				errChange := MerchantProjectRepository.PayoutOrderChannelFailed(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if errChange != nil {
					return appError.NewError(constant.ChangeErrMsg).FormatMessage("firstPayPayoutCallbackUpdate") // 更新代付订单错误 请重新提交
				}
			}
			return nil
		})
		if errTrans != nil {
			return s.Error(appError.NewError(errTrans.Error()))
		}

		// TODO 如果是 已完成状态 则需要添加到 redis 中 通知到下游商户 提现回调
		if isSendQueue == true {
			_ = service.NewSendQueueServer().SendNotifyQueue(payoutInfo.Sn)
		}
		if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
			_ = service.NewSendQueueServer().ManualSendNotifyQueue(payoutInfo.Sn)
		}
	}

	return s.Success("SUCCESS")
}
