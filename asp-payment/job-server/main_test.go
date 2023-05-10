package main

import (
	"asp-payment/common/pkg/config"
	"asp-payment/common/service"
	"testing"
)

func initConfig() {
	//初始化配置
	config.InitConfig()
}

func TestOrder(t *testing.T) {
	//merchantProjectCurrencyJobServer := service.NewSimpleMerchantProjectCurrencyJobServer()
	initConfig()
	//merchantProjectCurrencyJobServer.SyncMerchantProjectCurrencyJob()
	statisticServer := service.NewSimpleStatisticServer()
	//
	//for i := 0; i < 15; i++ {
	//	statisticServer.OrderDay(-1 * i)
	//}
	//for i := 0; i < 15; i++ {
	//	statisticServer.OrderDayMerchant(-1 * i)
	//}
	////统计上个月的数据
	//statisticServer.OrderMonth(-1)
	//statisticServer.OrderMonthMerchant(-1)
	//
	statisticServer.OrderDay(0)
	statisticServer.OrderDayMerchant(0)
}
