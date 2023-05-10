package service

import (
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/req"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

// 代付审核队列脚本

type PayoutAuditJobServer struct {
	*Service
}

var ticketPayoutAudit = time.Tick(1 * time.Minute)

var newSimplePayoutServer = NewSimplePayoutServer()

func NewSimplePayoutAuditJobServer() *PayoutAuditJobServer {
	return &PayoutAuditJobServer{&Service{LogFileName: constant.PayoutAuditServerJobLogFileName}}
}

func (s *PayoutAuditJobServer) SyncPayoutAuditCrondJob() {
	<-ticketPayoutAudit
	// 首先查询redis中的key
	redis := s.Redis()
	// 循环操作

	currentCount := 0
	// redis list key
	redisKey := constant.GetRedisKey(constant.KEY_PAYOUT_AUDIT_LIST)

	for {
		currentCount++
		if currentCount%constant.PayoutAuditCrondMaxCount == 0 {
			break
		}
		lLen := redis.LLen(context.Background(), redisKey)
		if lLen.Err() != nil {
			continue
		}
		if lLen.Val() <= 0 {
			continue
		}

		auditInfo := redis.RPop(context.Background(), redisKey)
		//fmt.Println("auditInfo.Err(): ", auditInfo.Err())
		//fmt.Println("auditInfo.Val(): ", auditInfo.Val())

		if auditInfo.Err() != nil {
			continue
		}
		logger.ApiWarn(s.LogFileName, s.RequestId, "auditInfo", zap.Any("auditInfo", auditInfo))
		var payoutAuditJobInfo req.AspPayoutAuditJob
		if err := goutils.JsonDecode(auditInfo.Val(), &payoutAuditJobInfo); err != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "payoutAuditJobInfo to json ", zap.Error(err))
			continue
		}
		// 执行审核操作
		errAudit := newSimplePayoutServer.PayoutAuditJob(&payoutAuditJobInfo)
		if errAudit != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "payoutAuditJobInfo to json ", zap.Error(errAudit))
			continue
		}

		// 更新代付审核状态

		//fmt.Println("payoutAuditJobInfo: ", payoutAuditJobInfo)

		time.Sleep(time.Duration(20) * time.Millisecond)
	}

	//fmt.Println("currentCount is OK: ", currentCount)

	//
}

func (s *PayoutAuditJobServer) SyncPayoutauditAddListCrondJob() {
	// 首先查询redis中的key
	redis := s.Redis()
	fmt.Println("redis: ", redis)
	// 循环操作
	// APPLY pass  CHANNEL_PENDING pass success FREEZE_SUCCESS return
	// CHANNEL_PENDING pass failed

	auditInfo := new(req.AspPayoutAuditJob)
	auditInfo.Id = 1173
	auditInfo.Status = "pass"
	auditInfo.OperationID = 1

	s.InsertRedis(auditInfo, redis)
	//auditInfo.Id = "869"
	//auditInfo.Status = "pass"
	//auditInfo.OperationID = "1"
	//s.InsertRedis(auditInfo, redis)
	//
	//fmt.Println("size12: ", size)
	//fmt.Println("err: ", err)

	//
}

func (s *PayoutAuditJobServer) InsertRedis(auditInfo *req.AspPayoutAuditJob, redis *redis.Client) {
	auditString, errJson := goutils.JsonEncode(auditInfo)

	fmt.Println("auditString: ", auditString)
	fmt.Println("errJson: ", errJson)

	// redis list key
	redisKey := constant.GetRedisKey(constant.KEY_PAYOUT_AUDIT_LIST)

	fmt.Println("redisKey: ", redisKey)

	redis.LPush(context.Background(), redisKey, auditString).Result()
}
