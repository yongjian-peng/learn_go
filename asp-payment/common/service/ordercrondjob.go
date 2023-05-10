package service

import (
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

type OrderCrondJobServer struct {
	*Service
}

func NewSimpleOrderCrondJobServer() *OrderCrondJobServer {
	return &OrderCrondJobServer{&Service{LogFileName: constant.OrderServerJobLogFileName}}
}

var ticket = time.Tick(5 * time.Minute)

func (s *OrderCrondJobServer) SyncOrderCrondJob() {
	<-ticket
	// 查询最近的订单 状态为
	//fmt.Println("SyncOrderCrondJob is start")
	//time.Sleep(3 * time.Second)
	//fmt.Println("SyncOrderCrondJob is end")

	newSampleOrderServer := NewSimplePayOrderServer()
	newSampleChannelConfigServer := NewSimpleChannelConfigServer()
	// 查询到前 n天的数据
	queryBeforeDay := config.AppConfig.JobConfig.QueryBeforeDay
	timeEnd := goutils.GetDateTimeUnix()
	timeBegin := timeEnd - (cast.ToInt64(queryBeforeDay) * 864000)

	orderList, err := newSampleOrderServer.GetAspOrderList(timeBegin, timeEnd)
	//logger.ApiWarn(s.LogFileName, s.RequestId, "orderList", zap.Any("orderList", orderList))
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "OrderCrondJob.GetAspOrderList", zap.Error(err))
		//time.Sleep(time.Duration(300) * time.Second)
		return
	}

	orderListCount := len(orderList)
	logger.ApiInfo(s.LogFileName, s.RequestId, "GetAspOrderList count: ", zap.Any("count", orderListCount))

	//var channelDepartInfo *model.AspChannelDepartConfig
	//var newOrderInfo *model.AspOrder

	currentCount := 0

	for _, orderInfo := range orderList {
		currentCount++
		// 一次最多300条
		if currentCount%constant.OrderCrondMaxUpdateCount == 0 {
			break
		}
		//fmt.Println("order: ", orderInfo)
		//fmt.Println("&order: ", &orderInfo)
		channelDepartInfoParam := map[string]string{}
		channelDepartInfoParam["depart_id"] = cast.ToString(orderInfo.DepartId)
		channelDepartInfoParam["channel_id"] = cast.ToString(orderInfo.ChannelId)
		channelDepartInfo, errDepart := newSampleChannelConfigServer.GetChannelDepartInfo(channelDepartInfoParam)
		if errDepart != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "OrderCrondJob.GetChannelDepartInfo", zap.Any("orderInfo", orderInfo), zap.Error(errDepart))
			continue
		}
		_, errQuery := newSampleOrderServer.OrderQuery(orderInfo, channelDepartInfo)
		//logger.ApiInfo(s.LogFileName, s.RequestId, "newOrderInfo: ", zap.Any("newOrderInfo", newOrderInfo))
		if errQuery != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "OrderCrondJob.OrderQuery", zap.Any("orderInfo", orderInfo), zap.Error(errQuery))
			continue
		}
		time.Sleep(time.Duration(20) * time.Millisecond)
		//fmt.Println("time: ", time.Now())
	}
}
