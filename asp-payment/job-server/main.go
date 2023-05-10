package main

import (
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goRedis"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/service"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cast"
)

func listenSignal(c *cron.Cron, cancel context.CancelFunc) {
	//chan
	cs := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(cs, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//阻塞直至有信号传入
	for {
		select {
		case s := <-cs:
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				logger.ApiInfo(constant.JobsServerLogFileName, "", "收到信号需要平滑关闭")
				//关闭着计划任务, 但是不能关闭已经在执行中的任务.
				ctx := c.Stop()
				//等待任务结束
				<-ctx.Done()
				cancel()
			default:
				//关闭着计划任务, 但是不能关闭已经在执行中的任务.
				ctx := c.Stop()
				//等待任务结束
				<-ctx.Done()
				cancel()
			}
		}
	}
}

func listenRedisSignal(preStartTime string, c *cron.Cron, cancel context.CancelFunc) {

	for {
		//定时获取redis状态，如果启动时间不一致，停止原脚本任务
		startTimeKey := goRedis.GetKey("crond:job:start:time")
		curStartTime := goRedis.Redis.Get(context.Background(), startTimeKey).Val()
		if curStartTime != preStartTime {
			ctx := c.Stop()
			//等待任务结束
			<-ctx.Done()
			cancel()
		}
		time.Sleep(time.Second * 1)
	}

}

// 设置启动时间，并获取启动时间
func getStartTime() string {
	startTimeKey := goRedis.GetKey("crond:job:start:time")
	startTime := cast.ToString(goutils.GetCurTimeUnixSecond())
	goRedis.Redis.Set(context.Background(), startTimeKey, startTime, -1)
	return startTime
}

func main() {

	//初始化配置 更新关联 common 017-001
	config.InitConfig()

	//创建任务调度器实例
	c := cron.New(cron.WithSeconds())
	// 初始化更新cp项目待结算余额到可用余额中
	merchantProjectCurrencyJobServer := service.NewSimpleMerchantProjectCurrencyJobServer()
	//注册任务到调度器，注册的任务都是异步执行的。
	//注册任务到调度器，注册的任务都是异步执行的。
	// 秒(Seconds) 分(Minutes) 时(Hours) 日(Day of month) 月(Month) 星期(Day of week)
	//每隔5秒执行一次：*/5 * * * * ?
	//每隔1分钟执行一次：0 */1 * * * ?
	//每天23点执行一次：0 0 23 * * ?
	//每天凌晨1点执行一次：0 0 1 * * ?
	//每月1号凌晨1点执行一次：0 0 1 1 * ?
	//在26分、29分、33分执行一次：0 26,29,33 * * * ?
	//每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
	// 每天6点执行脚本
	c.AddFunc("0 0 6 * * *", func() {
		merchantProjectCurrencyJobServer.SyncMerchantProjectCurrencyJob()
	})

	statisticServer := service.NewSimpleStatisticServer()
	//每天4天中执行脚本
	c.AddFunc("0 0 4 * * *", func() {
		//统计15天内的数据
		for i := 1; i <= 15; i++ {
			statisticServer.OrderDay(-1 * i)
		}
		for i := 1; i <= 15; i++ {
			statisticServer.OrderDayMerchant(-1 * i)
		}
		//统计上个月的数据
		statisticServer.OrderMonth(-1)
		statisticServer.OrderMonthMerchant(-1)
	})

	//每0,30分钟时执行一次，每天当前日报表
	c.AddFunc("0 0,30 * * * ?", func() {
		statisticServer.OrderDay(0)
		statisticServer.OrderDayMerchant(0)
	})

	//启动计划任务
	c.Start()

	// 父context(利用根context得到)
	ctx, cancel := context.WithCancel(context.Background())

	//监听信号，平滑关闭
	go listenSignal(c, cancel)
	//监听redis平滑关闭
	go listenRedisSignal(getStartTime(), c, cancel)

	//取消事件监听
	<-ctx.Done()
}
