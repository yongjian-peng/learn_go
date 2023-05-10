package main

import (
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goRedis"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/service"
	"asp-payment/crond-server/pkg"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cast"
)

func listenSignal(cancel context.CancelFunc) {
	//chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//阻塞直至有信号传入
	for {
		select {
		case s := <-c:
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				logger.ApiInfo(constant.CrondServerLogFileName, "", "收到信号需要平滑关闭")
				cancel()
			default:
				cancel()
			}
		}
	}
}

func listenRedisSignal(preStartTime string, cancel context.CancelFunc) {

	for {
		//定时获取redis状态，如果启动时间不一致，停止原脚本任务
		startTimeKey := goRedis.GetKey("crond:crond:start:time")
		curStartTime := goRedis.Redis.Get(context.Background(), startTimeKey).Val()
		if curStartTime != preStartTime {
			cancel()
		}
		time.Sleep(time.Second * 1)
	}

}

// 设置启动时间，并获取启动时间
func getStartTime() string {
	startTimeKey := goRedis.GetKey("crond:crond:start:time")
	startTime := cast.ToString(goutils.GetCurTimeUnixSecond())
	goRedis.Redis.Set(context.Background(), startTimeKey, startTime, -1)
	return startTime
}

func main() {

	//初始化配置 更新关联 common 0417-003
	config.InitConfig()

	//定义一个WaitGroup，阻塞主线程执行
	var wg sync.WaitGroup
	defer wg.Wait()
	// 父context(利用根context得到)
	ctx, cancel := context.WithCancel(context.Background())
	//监听信号
	go listenSignal(cancel)
	//监听redis
	go listenRedisSignal(getStartTime(), cancel)

	//运行的job列表 添加测试 同步 common
	jobs := make([]*pkg.CrondJob, 0)
	// 父context的子协程
	orderCrondServer := service.NewSimpleOrderCrondJobServer()
	payoutCrondServer := service.NewSimplePayoutJobServer()
	payoutAuditCrondServer := service.NewSimplePayoutAuditJobServer()
	jobs = append(jobs, &pkg.CrondJob{Name: "syncOrderJob", Job: orderCrondServer.SyncOrderCrondJob, Ctx: ctx, Wg: &wg})
	jobs = append(jobs, &pkg.CrondJob{Name: "SyncPayoutCrondJob", Job: payoutCrondServer.SyncPayoutCrondJob, Ctx: ctx, Wg: &wg})
	jobs = append(jobs, &pkg.CrondJob{Name: "SyncPayoutAuditCrondJob", Job: payoutAuditCrondServer.SyncPayoutAuditCrondJob, Ctx: ctx, Wg: &wg})

	for _, job := range jobs {
		job.Run()
	}

}
