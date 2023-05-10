package callback

import (
	"asp-payment/api-server/req/zpayreq"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"asp-payment/common/service"
	"asp-payment/common/service/supplier/impl/zpayimpl"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ZPayCallBackService struct {
	*service.Service
}

func NewZPayCallBackServer(c *fiber.Ctx) *ZPayCallBackService {
	return &ZPayCallBackService{Service: service.NewService(c, constant.ZPayCallBackLogFileName)}
}

func (s *ZPayCallBackService) PayOrder() error {

	// 贯穿支付所需要的参数
	bodyReq := new(zpayreq.ZPayCallBackReq)

	if err := s.C.QueryParser(bodyReq); err != nil {
		return err
	}

	bod, err := goutils.JsonEncode(bodyReq)
	OrderServer := service.NewPayOrderServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["sn"] = bodyReq.PartnerOrderNo

	orderInfo, oErr := OrderServer.GetOrderInfo(orderInfoParam, nil)
	if oErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "OrderServer.GetOrderInfo  ", zap.Error(oErr))
		return s.C.SendString("1")
	}
	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if cast.ToUint(bodyReq.Amount) != orderInfo.TotalFee {
		logger.ApiWarn(s.LogFileName, s.RequestId, "bodyReq.Amount != orderInfo.TotalFee ", zap.Any("orderInfo", orderInfo))
		failedParams := make(map[string]interface{})
		failedParams["supplier_return_code"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		failedParams["supplier_return_msg"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
		if _, oErr = OrderServer.SolveOrderPaySuccess(orderInfo.Id, failedParams, "firstPayCallbackUpdate", nil); oErr != nil {
			logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess firstPayCallbackUpdate", zap.Error(oErr))
			return appError.CodeUnknown // 服务器开小差了！
		}
		return s.C.SendString("1")
	}

	// 如果已经是成功的状态了 则直接放回 不做操作
	if orderInfo.TradeState == constant.ORDER_TRADE_STATE_SUCCESS || orderInfo.TradeState == constant.ORDER_TRADE_STATE_FAILED {
		logger.ApiInfo(s.LogFileName, s.RequestId, fmt.Sprintf("orderInfo.TradeState == %s", orderInfo.TradeState))
		return s.C.SendString("0")
	}

	merchantProjectServer := service.NewMerchantProjectServer(s.C)
	appId := cast.ToString(orderInfo.MchProjectId)
	// appError.MissMerchantProjectNotFoundErr.FormatMessage("xxxx")
	_, pErr := merchantProjectServer.GetMerchantProjectInfo(appId)
	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetMerchantProjectInfo  ", zap.Error(pErr))
		return s.C.SendString("1")
	}

	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(orderInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(orderInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, cErr := service.NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if cErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetChannelDepartInfo  ", zap.Error(cErr))
		return s.C.SendString("1")
	}

	// config 转 struct 字符串转 struct
	var channelConfigInfo model.AspChannelDepartConfigInfo
	if err = goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "channelDepartInfo.Config to json ", zap.Error(err))
		return s.C.SendString("1")
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig channelConfigInfo ", zap.String("PartnerId", channelConfigInfo.PartnerId), zap.String("Signature", channelConfigInfo.Signature))
	var m map[string]interface{}
	_ = goutils.JsonDecode(bod, &m)

	signMap := make(map[string]interface{})
	for key, val := range m {
		key = goutils.FirstLower(key)
		signMap[key] = val
	}

	signature := zpayimpl.VerifySignature(signMap, channelConfigInfo.Signature)
	if signature != bodyReq.Sign {
		logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("bodyReq.Sign", bodyReq.Sign))
		return s.C.SendString("1")
	}

	// 4. 判断订单的状态，修改订单状态 以及相关信息
	// 如果当前的提现状态和上游返回的情况不一致的情况
	upstreamStatus := zpayimpl.GetZPayPaymentCallBackStatus(bodyReq.Status)
	if orderInfo.TradeState != upstreamStatus {
		if upstreamStatus != constant.ORDER_TRADE_STATE_PENDING {
			params := make(map[string]interface{})
			params["transaction_id"] = orderInfo.TransactionId // 默认值
			params["trade_state"] = orderInfo.TradeState       // 默认值
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
				orderSuccessUpdate["transaction_id"] = bodyReq.OrderNo
				orderSuccessUpdate["trade_state"] = upstreamStatus
				orderSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
				errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
					if _, errSuccess := OrderServer.SolveOrderPaySuccess(orderInfo.Id, orderSuccessUpdate, "firstPayCallbackUpdate", tx); errSuccess != nil {
						logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess firstPayCallbackUpdate", zap.Error(oErr))
						return appError.CodeUnknown // 服务器开小差了！
					}
					// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
					mErr := MerchantProjectRepository.PayinOrderSuccess(orderInfo.MchProjectId, amount, orderInfo.Id, tx)
					if mErr != nil {
						return appError.NewError(constant.ChangeErrMsg).FormatMessage("zPayOrderCallbackUpdate") // 更新代付订单错误 请重新提交
					}
					return nil
				})
				if errTrans != nil {
					return s.Error(appError.NewError(errTrans.Error()))
				}
				isSendQueue = true
			}

			if upstreamStatus == constant.ORDER_TRADE_STATE_FAILED {
				params["transaction_id"] = bodyReq.OrderNo
				params["trade_state"] = upstreamStatus
				if _, fErr := OrderServer.SolveOrderPaySuccess(orderInfo.Id, params, "zPayCallbackUpdate", nil); fErr != nil {
					logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess zPayCallbackUpdate", zap.Error(err))
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

	return s.C.SendString("0")
}

func (s *ZPayCallBackService) Payout() error {
	// 贯穿支付所需要的参数
	bodyReq := new(zpayreq.ZPayPayoutCallBackReq)
	if err := s.C.QueryParser(bodyReq); err != nil {
		return err
	}
	bod, err := goutils.JsonEncode(bodyReq)
	payoutServer := service.NewPayoutServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单
	payoutInfoParam := make(map[string]interface{})
	payoutInfoParam["sn"] = bodyReq.PartnerWithdrawNo

	payoutInfo, pErr := payoutServer.GetPayoutInfo(payoutInfoParam, nil)
	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "payoutServer.GetPayoutInfo  ", zap.Error(pErr))
		return s.C.SendString("1")
	}
	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if cast.ToUint(bodyReq.Amount) != payoutInfo.TotalFee {
		logger.ApiWarn(s.LogFileName, s.RequestId, "bodyReq.Amount != orderInfo.TotalFee ", zap.Any("payoutInfo", payoutInfo))
		payoutFailedUpdate := make(map[string]interface{})
		payoutFailedUpdate["supplier_return_code"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
		payoutFailedUpdate["supplier_return_msg"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE

		_, fErr := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "zpayOrderUpdate", nil)
		if fErr != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
		}
		return s.C.SendString("1")
	}

	// 如果已经是成功的状态了 则直接放回 不做操作
	if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_SUCCESS || payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_FAILED {
		logger.ApiInfo(s.LogFileName, s.RequestId, fmt.Sprintf("payoutInfo.TradeState == %s", payoutInfo.TradeState))
		return s.C.SendString("0")
	}

	merchantProjectServer := service.NewMerchantProjectServer(s.C)
	appId := cast.ToString(payoutInfo.MchProjectId)
	// appError.MissMerchantProjectNotFoundErr.FormatMessage("xxxx")
	_, pEr := merchantProjectServer.GetMerchantProjectInfo(appId)
	if pEr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetMerchantProjectInfo  ", zap.Error(pEr))
		return s.C.SendString("2")
	}

	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(payoutInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(payoutInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, cErr := service.NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if cErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetChannelDepartInfo  ", zap.Error(cErr))
		return s.C.SendString("3")
	}

	// config 转 struct 字符串转 struct
	var channelConfigInfo model.AspChannelDepartConfigInfo
	if err = goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "channelDepartInfo.Config to json ", zap.Error(err))
		return s.C.SendString("4")
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig channelConfigInfo ", zap.Any("channelConfigInfo", channelConfigInfo))

	var m map[string]interface{}
	_ = goutils.JsonDecode(bod, &m)
	signMap := make(map[string]interface{})
	for key, val := range m {
		key = goutils.FirstLower(key)
		signMap[key] = val
	}
	logger.ApiInfo(s.LogFileName, s.RequestId, "json.Unmarshal ", zap.Any("PartnerId", signMap))
	signature := zpayimpl.VerifySignature(signMap, channelConfigInfo.PayoutSignature)

	if signature != bodyReq.Sign {
		logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("reqHeader.Signature", bodyReq.Sign))
		return s.C.SendString("5")
	}

	/// 提现查询更新逻辑 横列是 上游返回 竖列是提现当前的状态
	//                 已申请              已成功               已失败
	// 已申请          不更新              更新成功             更新失败
	// 已成功          不更新               不更新                待定
	// 已失败          不更新                待定                不更新
	//
	// 4. 判断订单的状态，修改订单状态 以及相关信息

	// 如果当前的提现状态和上游返回的情况不一致的情况
	status := cast.ToInt(bodyReq.Status)
	upstreamStatus := zpayimpl.GetZPayPayoutCallBackStatus(status)
	if payoutInfo.TradeState != upstreamStatus {
		//params := make(map[string]interface{})
		//params["trade_state"] = payoutInfo.TradeState // 默认值
		//params["finish_time"] = cast.ToString(payoutInfo.FinishTime)
		var isSendQueue bool
		isSendQueue = false

		totalFee := cast.ToInt(payoutInfo.TotalFee)
		changeAmount := totalFee + payoutInfo.ChargeFee + payoutInfo.FixedAmount
		// 如果是成功 必须提现状态是 顺序的增长的，不能状态逆序
		// 如果上游返回成功 提现状态是 已申请 则修改
		if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
			//params["finish_time"] = goutils.GetDateTimeUnix()
			//params["trade_state"] = upstreamStatus
			isSendQueue = true
			// 如果上游返回失败 提现状态是 已申请 则修改
		}
		errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
			// 当成功时候
			if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
				MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)

				// 代付上游成功 解冻+代付成功  冻结金额减去  预扣金额记录新增  收支流水扣减
				payoutSuccessUpdate := make(map[string]interface{})
				payoutSuccessUpdate["transaction_id"] = bodyReq.WithdrawNo
				payoutSuccessUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_SUCCESS
				payoutSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
				_, errSuccess := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutSuccessUpdate, "zpayOrderUpdate", tx)
				if errSuccess != nil {
					return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				sErr := MerchantProjectRepository.PayoutOrderChannelSuccess(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if sErr != nil {
					return appError.NewError("6") // 更新商户余额错误 请重新提交
				}
			}

			// 当失败的时候
			if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
				MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)

				// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
				payoutFailedUpdate := make(map[string]interface{})
				payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FAILED
				payoutFailedUpdate["transaction_id"] = bodyReq.WithdrawNo
				_, errFailed := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "zpayOrderUpdate", tx)
				if errFailed != nil {
					return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				fErr := MerchantProjectRepository.PayoutOrderChannelFailed(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if fErr != nil {
					return appError.NewError("7") // 更新商户余额错误 请重新提交
				}
			}
			return nil
		})
		if errTrans != nil {
			return s.C.SendString(errTrans.Error())
		}
		// TODO 如果是 已完成状态 则需要添加到 redis 中 通知到下游商户 提现回调
		if isSendQueue == true {
			_ = service.NewSendQueueServer().SendNotifyQueue(payoutInfo.Sn)
		}
		if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
			_ = service.NewSendQueueServer().ManualSendNotifyQueue(payoutInfo.Sn)
		}
	}
	return s.C.SendString("0")
}
