package service

import (
	"asp-payment/api-server/req"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/logger"
	"context"
	"encoding/json"

	"go.uber.org/zap"
)

type SendQueueServer struct {
	*Service
}

func NewSendQueueServer() *SendQueueServer {
	return &SendQueueServer{&Service{LogFileName: constant.OrderServerLogFileName}}
}

func (s *SendQueueServer) SendNotifyQueue(sn string) error {

	err := s.Redis().LPush(context.Background(), constant.KEY_ORDER_NOTIFY_QUEUE, sn).Err()
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "redis.KEY_ORDER_NOTIFY_QUEUE err ", zap.String("sn", sn), zap.Error(err))
	}
	return err
}

// ManualSendNotifyQueue 手动加入回调队列，没有状态的判断
func (s *SendQueueServer) ManualSendNotifyQueue(sn string) error {

	err := s.Redis().LPush(context.Background(), constant.KEY_ORDER_NOTIFY_MANUAL_QUEUE, sn).Err()
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "redis.KEY_ORDER_NOTIFY_MANUAL_QUEUE err ", zap.String("sn", sn), zap.Error(err))
	}
	return err
}

func (s *SendQueueServer) SendMerchantProjectAmountQueue(requestId string, merchantProjectQueue *req.MerchantProjectTotalFeeQueue) error {
	// 需要的字段 参数

	merchantProjectQueueJson, err := json.Marshal(merchantProjectQueue)
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "redis.SendMerchantProjectAmountQueue err ", zap.String("request_id", requestId), zap.Error(err))
		return err
	}

	queueString := string(merchantProjectQueueJson)

	err = s.Redis().LPush(context.Background(), constant.KEY_MERCHANT_PROJECT_AMOUNT_QUEUE, queueString).Err()
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "redis.KEY_MERCHANT_PROJECT_TOTAL_FEE_QUEUE err ", zap.String("queueString", queueString), zap.Error(err))
	}
	return err
}
