package callback

import (
	"asp-payment/api-server/req/abcpayreq"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"asp-payment/common/service"
	"asp-payment/common/service/supplier/impl/abcpay"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AbcPayCallBackService struct {
	*service.Service
}

func NewAbcPayCallBackService(c *fiber.Ctx) *AbcPayCallBackService {
	return &AbcPayCallBackService{Service: service.NewService(c, constant.AbcPayCallBackLogFileName)}
}

func (s *AbcPayCallBackService) PayOrder() error {

	// 贯穿支付所需要的参数
	bodyReq := new(abcpayreq.CallBackOrderReq)
	if err := s.C.BodyParser(bodyReq); err != nil {
		return s.C.SendString("Error 1")
	}

	//根据订单号加锁
	lockName := fmt.Sprintf("abcpay:callback:order:lock:%s", bodyReq.OutBizNo)
	if !s.Lock(lockName) {
		return s.C.SendString("Error 2")
	}
	defer s.UnLock(lockName)

	bod, _ := goutils.JsonEncode(bodyReq)
	OrderServer := service.NewPayOrderServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["sn"] = bodyReq.OutBizNo

	orderInfo, oErr := OrderServer.GetOrderInfo(orderInfoParam, nil)
	if oErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "OrderServer.GetOrderInfo  ", zap.Error(oErr))
		return s.C.SendString("Error 3")
	}

	// 如果已经是成功的状态了 则直接放回 不做操作
	if orderInfo.TradeState == constant.ORDER_TRADE_STATE_SUCCESS || orderInfo.TradeState == constant.ORDER_TRADE_STATE_FAILED {
		logger.ApiInfo(s.LogFileName, s.RequestId, fmt.Sprintf("orderInfo.TradeState == %s", orderInfo.TradeState))
		return s.C.SendString("OK")
	}

	status := bodyReq.Status
	upstreamStatus := abcpay.GetPaymentTradeState(cast.ToString(status))
	//成功的订单才修改
	if upstreamStatus != constant.ORDER_TRADE_STATE_SUCCESS {
		return s.C.SendString("Error 4")
	}

	merchantProjectServer := service.NewMerchantProjectServer(s.C)
	projectId := cast.ToString(orderInfo.MchProjectId)
	_, pErr := merchantProjectServer.GetMerchantProjectInfo(projectId)
	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetMerchantProjectInfo  ", zap.Error(pErr))
		return s.C.SendString("Error 5")
	}

	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(orderInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(orderInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, cErr := service.NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if cErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetChannelDepartInfo  ", zap.Error(cErr))
		return s.C.SendString("Error 6")
	}

	// config 转 struct 字符串转 struct
	var channelConfigInfo model.AspChannelDepartConfigInfo
	if err := goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "channelDepartInfo.Config to json ", zap.Error(err))
		return s.C.SendString("Error 7")
	}

	logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig channelConfigInfo ", zap.String("PartnerId", channelConfigInfo.PartnerId), zap.String("Signature", channelConfigInfo.Signature))
	var m map[string]interface{}

	_ = goutils.JsonDecode(bod, &m)
	signature := abcpay.GetFixSignature(m, []string{"appid", "status", "money", "order_id", "out_biz_no"}, channelConfigInfo.Signature)
	if signature != bodyReq.Sign {
		logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("bodyReq.Sign", bodyReq.Sign))
		return s.C.SendString("Error 8")
	}

	//修改订单状态
	MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
	totalFee := cast.ToInt(orderInfo.TotalFee)
	amount := totalFee - orderInfo.ChargeFee - orderInfo.FixedAmount
	// 代收上游成功 待结算余额新增    收支流水记录新增
	orderSuccessUpdate := make(map[string]interface{})
	orderSuccessUpdate["trade_state"] = constant.ORDER_TRADE_STATE_SUCCESS
	orderSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
	errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
		if _, oErr = OrderServer.SolveOrderPaySuccess(orderInfo.Id, orderSuccessUpdate, "firstPayCallbackUpdate", tx); oErr != nil {
			logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess abcPayCallbackUpdate", zap.Error(oErr))
			return appError.CodeUnknown // 服务器开小差了！
		}
		// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
		errChange := MerchantProjectRepository.PayinOrderSuccess(orderInfo.MchProjectId, amount, orderInfo.Id, tx)
		if errChange != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("abcPayOrderCallbackUpdate") // 更新代付订单错误 请重新提交
		}
		return nil
	})
	if errTrans != nil {
		return s.C.SendString("fail")
	}

	_ = service.NewSendQueueServer().SendNotifyQueue(orderInfo.Sn)

	return s.C.SendString("OK")
}

func (s *AbcPayCallBackService) Payout() error {
	// 贯穿支付所需要的参数
	bodyReq := new(abcpayreq.CallBackPayoutReq)

	if err := s.C.BodyParser(bodyReq); err != nil {
		return s.C.SendString("Error 1")
	}
	bod, err := goutils.JsonEncode(bodyReq)
	payoutServer := service.NewPayoutServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单
	payoutInfoParam := make(map[string]interface{})
	payoutInfoParam["sn"] = bodyReq.OutBizNo
	//根据订单号加锁
	lockName := fmt.Sprintf("abcpay:callback:payout:lock:%s", bodyReq.OutBizNo)
	if !s.Lock(lockName) {
		return s.C.SendString("Error")
	}
	defer s.UnLock(lockName)

	payoutInfo, pErr := payoutServer.GetPayoutInfo(payoutInfoParam, nil)
	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "payoutServer.GetPayoutInfo  ", zap.Error(pErr))
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
	if err = goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "channelDepartInfo.Config to json ", zap.Error(err))
		return s.C.SendString("Error")
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig channelConfigInfo ", zap.Any("channelConfigInfo", channelConfigInfo))

	var m map[string]interface{}
	_ = goutils.JsonDecode(bod, &m)
	logger.ApiInfo(s.LogFileName, s.RequestId, "json.Unmarshal ", zap.Any("PartnerId", m))
	signature := abcpay.GetFixSignature(m, []string{"appid", "orderstatus", "amount", "endtime"}, channelConfigInfo.PayoutSignature)
	if signature != bodyReq.Sign {
		logger.ApiWarn(s.LogFileName, s.RequestId, "signature != reqHeader.Signature ", zap.String("signature", signature), zap.String("reqHeader.Signature", bodyReq.Sign))
		return s.C.SendString("Error7")
	}

	// 如果当前的提现状态和上游返回的情况不一致的情况
	status := bodyReq.Orderstatus
	upstreamStatus := abcpay.GetPayoutCallBackStatus(cast.ToString(status))
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

				_, errSuccess := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutSuccessUpdate, "abcPayPayoutOrderCallbackUpdate", tx)
				if errSuccess != nil {
					return appError.NewError("Error06") // 更新商户余额错误 请重新提交
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				changeErr := MerchantProjectRepository.PayoutOrderChannelSuccess(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if changeErr != nil {
					return appError.NewError("Error8") // 更新商户余额错误 请重新提交
				}
			}

			// 当失败的时候
			if upstreamStatus == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
				MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
				// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
				payoutFailedUpdate := make(map[string]interface{})
				payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FAILED

				_, errFailed := payoutServer.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "abcPayPayoutOrderCallbackUpdate", tx)
				if errFailed != nil {
					return appError.NewError("Error09") // 更新商户余额错误 请重新提交
				}
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				err := MerchantProjectRepository.PayoutOrderChannelFailed(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
				if err != nil {
					return appError.NewError("Error10") // 更新商户余额错误 请重新提交
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
	return s.C.SendString("OK")
}
