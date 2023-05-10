package service

import (
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type PayoutJobServer struct {
	*Service
}

var ticketPayout = time.Tick(5 * time.Minute)

func NewSimplePayoutJobServer() *PayoutJobServer {
	return &PayoutJobServer{&Service{LogFileName: constant.PayoutServerJobLogFileName}}
}

func (s *PayoutJobServer) SyncPayoutCrondJob() {
	<-ticketPayout
	newSamplePayoutServer := NewSimplePayoutServer()
	newSampleChannelConfigServer := NewSimpleChannelConfigServer()
	newSendQueueServer := NewSendQueueServer()
	// 查询到前 n天的数据
	queryBeforeDay := config.AppConfig.JobConfig.QueryBeforeDay
	timeEnd := goutils.GetDateTimeUnix()
	timeBegin := timeEnd - (cast.ToInt64(queryBeforeDay) * 864000)

	payoutList, err := newSamplePayoutServer.GetAspPayoutList(timeBegin, timeEnd)
	//logger.ApiWarn(s.LogFileName, s.RequestId, "orderList", zap.Any("orderList", payoutList))
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "OrderCrondJob.GetAspPayoutList", zap.Error(err))
		//time.Sleep(time.Duration(300) * time.Second)
		return
	}

	payoutListCount := len(payoutList)
	logger.ApiInfo(s.LogFileName, s.RequestId, "GetAspPayoutList count: ", zap.Any("count", payoutListCount))

	//var channelDepartInfo *model.AspChannelDepartConfig
	//var newPayoutInfo *model.AspPayout
	currentCount := 0
	var isSendQueue bool
	db := database.DB
	for _, payoutInfo := range payoutList {
		currentCount++
		// 一次最多300条
		if currentCount > constant.OrderCrondMaxUpdateCount {
			break
		}
		//fmt.Println("key", key, "order: ", orderInfo)
		channelDepartInfoParam := map[string]string{}
		channelDepartInfoParam["depart_id"] = cast.ToString(payoutInfo.DepartId)
		channelDepartInfoParam["channel_id"] = cast.ToString(payoutInfo.ChannelId)
		channelDepartInfo, errDepart := newSampleChannelConfigServer.GetChannelDepartInfo(channelDepartInfoParam)
		if errDepart != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "OrderCrondJob.GetAspPayoutList", zap.Any("payoutInfo", payoutInfo), zap.Error(errDepart))
			continue
		}
		isSendQueue = false
		scanQueryData, errQuery := newSamplePayoutServer.RequestUpstreamQueryPayout(payoutInfo, channelDepartInfo)
		if errQuery != nil {
			total := cast.ToInt(payoutInfo.TotalFee)
			changeAmount := total + payoutInfo.ChargeFee + payoutInfo.FixedAmount
			merchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
			// 是否上游返回异常 + 更新订单 + 释放可用金额 + 事务
			errUpdateFailed := newSamplePayoutServer.PayoutQueryUpdateFailed(merchantProjectRepository, scanQueryData, payoutInfo, errQuery, changeAmount, db)
			if errUpdateFailed != nil {
				logger.ApiWarn(s.LogFileName, s.RequestId, "PayoutQueryUpdateFailed", zap.Any("orderInfo", payoutInfo), zap.Error(errUpdateFailed))
				continue
			}
			if errQuery.Code == appError.CodeSupplierInternalChannelParamsFailedErrCode.Code {
				_ = newSendQueueServer.ManualSendNotifyQueue(payoutInfo.Sn)
			}
			logger.ApiWarn(s.LogFileName, s.RequestId, "RequestUpstreamQueryPayout", zap.Any("orderInfo", payoutInfo), zap.Error(errQuery))
			continue
		}
		if scanQueryData.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
			isSendQueue = true
		}
		errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
			totalFee := cast.ToInt(payoutInfo.TotalFee)
			changeAmount := totalFee + payoutInfo.ChargeFee + payoutInfo.FixedAmount
			_, errPayoutUpdate := newSamplePayoutServer.PayoutQueryUpdate(payoutInfo, scanQueryData, changeAmount, tx)
			if errPayoutUpdate != nil {
				return errPayoutUpdate
			}
			return nil
		})
		if errTrans != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "OrderCrondJob.RequestUpstreamQueryPayout", zap.Any("orderInfo", payoutInfo), zap.Error(errQuery))
			continue
		}
		// 发送回调到redis
		if isSendQueue == true {
			_ = newSendQueueServer.SendNotifyQueue(payoutInfo.Sn)
		}
		if scanQueryData.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
			_ = newSendQueueServer.ManualSendNotifyQueue(payoutInfo.Sn)
		}
		time.Sleep(time.Duration(20) * time.Millisecond)
		//fmt.Println("time payout : ", time.Now())
	}

	//fmt.Println("channelDepartInfo: ", channelDepartInfo)
	//_ = newPayoutInfo

}

func (s *PayoutJobServer) SyncOrderPushJob() {
	fmt.Println("SyncPayoutCrondJob is start")
	time.Sleep(3 * time.Second)
	fmt.Println("SyncPayoutCrondJob is end")
}
