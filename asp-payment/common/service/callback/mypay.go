package callback

import (
	"asp-payment/api-server/req/mypayreq"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"asp-payment/common/service"
	"asp-payment/common/service/supplier/impl/haodapay"
	"asp-payment/common/service/supplier/impl/mypay"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MyPayCallBackService struct {
	*service.Service
}

func NewMyPayCallBackService(c *fiber.Ctx) *MyPayCallBackService {
	return &MyPayCallBackService{Service: service.NewService(c, constant.MyPayCallBackLogFileName)}
}

func (s *MyPayCallBackService) PayOrder() error {

	bodyReq := new(mypayreq.CallBackReq)

	if errBody := s.C.BodyParser(bodyReq); errBody != nil {
		return s.Error(appError.NewError("Error Params"))
	}

	logger.ApiWarn(s.LogFileName, s.RequestId, "PayOrder ", zap.Any("bodyReq", bodyReq))
	//根据订单号加锁
	lockName := fmt.Sprintf("mypay:callback:order:lock:%s_%s_%s", constant.ProviderSunny, constant.TradeTypeMyPay, bodyReq.UserTxnID)
	if !s.Lock(lockName) {
		return s.Error(appError.NewError("Error Params"))
	}
	defer s.UnLock(lockName)

	//bod, _ := goutils.JsonEncode(bodyReq)
	OrderServer := service.NewPayOrderServer(s.C)
	// 1. 根据订单的唯一标识 去查询订单
	orderInfoParam := make(map[string]interface{})
	orderInfoParam["transaction_id"] = bodyReq.UserTxnID

	orderInfo, oErr := OrderServer.GetOrderInfo(orderInfoParam, nil)
	if oErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "OrderServer.GetOrderInfo  ", zap.Error(oErr))
		return s.Error(appError.NewError("Error Params"))
	}

	// 判断渠道返回的金额和支付金额是否一致 因为出现了支付5万 成功订单金额是1 的情况
	if goutils.Yuan2Fen(cast.ToFloat64(bodyReq.Amount)) != cast.ToInt64(orderInfo.TotalFee) && cast.ToInt(bodyReq.Amount) != 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "bodyReq.Money != orderInfo.TotalFee ", zap.Any("orderInfo", orderInfo))
		if orderInfo.TradeState != constant.ORDER_TRADE_STATE_PAYERROR {
			failedParams := make(map[string]interface{})
			//failedParams["trade_state"] = constant.ORDER_TRADE_STATE_PAYERROR
			failedParams["supplier_return_code"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE
			failedParams["supplier_return_msg"] = constant.SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE
			if _, oErr = OrderServer.SolveOrderPaySuccess(orderInfo.Id, failedParams, "myPayCallbackUpdate", nil); oErr != nil {
				logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess myPayCallbackUpdate", zap.Error(oErr))
				// the server has deserted!
				return s.Error(appError.NewError("the server has deserted!"))
			}
		}
		return s.Error(appError.NewError("Error 2"))
	}

	// 如果已经是成功的状态了 则直接放回 不做操作
	if orderInfo.TradeState == constant.ORDER_TRADE_STATE_SUCCESS || orderInfo.TradeState == constant.ORDER_TRADE_STATE_FAILED {
		logger.ApiInfo(s.LogFileName, s.RequestId, fmt.Sprintf("orderInfo.TradeState == %s", orderInfo.TradeState))
		return s.Success("")
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

	upstreamStatus := mypay.GetPaymentTradeState(bodyReq.Status)
	// fmt.Println("upstreamStatus: ", upstreamStatus)

	if orderInfo.TradeState != upstreamStatus {
		transactionId := bodyReq.MyPayTransactionID
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
				orderSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
				errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
					if _, errSuccess := OrderServer.SolveOrderPaySuccess(orderInfo.Id, orderSuccessUpdate, "myPayOrderCallbackUpdate", tx); errSuccess != nil {
						logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess myPayOrderCallbackUpdate", zap.Error(oErr))
						return appError.CodeUnknown // 服务器开小差了！
					}
					// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
					mErr := MerchantProjectRepository.PayinOrderSuccess(orderInfo.MchProjectId, amount, orderInfo.Id, tx)
					if mErr != nil {
						return appError.NewError(constant.ChangeErrMsg).FormatMessage("myPayOrderCallbackUpdate") // 更新代付订单错误 请重新提交
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
					logger.ApiError(s.LogFileName, s.RequestId, "OrderServer.SolveOrderPaySuccess myPayOrderCallbackUpdate", zap.Error(fErr))
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

	return s.Success("")
}

func (s *MyPayCallBackService) Payout() error {
	return nil
}
