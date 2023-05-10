package service

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goRedis"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"github.com/golang-module/carbon/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type StatisticServer struct {
	*Service
}

func NewSimpleStatisticServer() *StatisticServer {
	return &StatisticServer{&Service{LogFileName: constant.StatisticServerLogFileName}}
}

type IncomeOrderResult struct {
	TotalFee      int64 `json:"total_fee,omitempty"`
	TotalNum      int   `json:"total_num,omitempty"`
	ProfitFee     int64 `json:"profit_fee,omitempty"`
	TotalOrderNum int   `json:"total_order_num,omitempty"`
}

type PayoutOrderResult struct {
	TotalFee        int64 `json:"total_fee,omitempty"`
	TotalNum        int   `json:"total_num,omitempty"`
	ProfitFee       int64 `json:"profit_fee,omitempty"`
	SuccessOrderNum int   `json:"success_order_num,omitempty"`
	TotalOrderNum   int   `json:"total_order_num,omitempty"`
}

// OrderBase 订单统计
func (s *StatisticServer) OrderBase(dateType uint, date int) {

	var curDate int64 = 0
	var dayStartTime int64 = 0
	var endStartTime int64 = 0
	//日报表
	if dateType == 1 {
		curDate = carbon.Parse(carbon.Now().AddDays(date).Format("Y-m-d")).StartOfDay().Timestamp()
		dayStartTime = curDate
		endStartTime = carbon.Parse(carbon.Now().AddDays(date).Format("Y-m-d")).EndOfDay().Timestamp()
	} else {
		curDate = carbon.Parse(carbon.Now().AddMonths(date).Format("Y-m-d H:i:s")).StartOfMonth().Timestamp()
		dayStartTime = curDate
		endStartTime = carbon.Parse(carbon.Now().AddMonths(date).Format("Y-m-d H:i:s")).EndOfMonth().Timestamp()
	}
	logger.ApiInfo(s.LogFileName, s.RequestId, "OrderBase", zap.Any("dateType", dateType), zap.Any("date", date), zap.Any("dayStartTime", dayStartTime), zap.Any("endStartTime", endStartTime))
	//获取渠道列表
	aspChannelList, _ := repository.NewRepository[*model.AspChannelConfig](database.DB, goRedis.Redis).Find(database.NewSqlCondition())
	//获取内部商户列表
	aspDepartList, _ := repository.NewRepository[*model.AspDeparts](database.DB, goRedis.Redis).Find(database.NewSqlCondition())
	//CP项目id
	aspMerchantProjectList, _ := repository.NewRepository[*model.AspMerchantProject](database.DB, goRedis.Redis).Find(database.NewSqlCondition())
	//订单统计
	aspOrderStatisticRepository := repository.NewRepository[*model.AspOrderStatistic](database.DB, goRedis.Redis)
	//项目
	aspMerchantProjectCurrencyRepository := repository.NewRepository[*model.AspMerchantProjectCurrency](database.DB, goRedis.Redis)
	//渠道配置
	aspChannelConfigRepository := repository.NewRepository[*model.AspChannelConfig](database.DB, goRedis.Redis)
	aspOrder := model.AspOrder{}
	aspPayout := model.AspPayout{}
	//trade_type 支付方式
	TradeTypeList := []string{constant.TradeType_H5, constant.TradeType_PAYOUT}
	for _, channel := range aspChannelList {
		for _, depart := range aspDepartList {
			for _, project := range aspMerchantProjectList {
				for _, tradeType := range TradeTypeList {
					channelId := channel.Id
					departId := depart.Id
					projectId := project.Id
					incomeSuccessRate := 0 //代收成功率（0-100 之间的值）
					payoutSuccessRate := 0 //代付成功率（0-100 之间的值）
					incomeOrderResult := IncomeOrderResult{}
					payoutOrderResult := PayoutOrderResult{}

					//代收
					database.DB.Table(aspOrder.TableName()).Select("sum(total_fee) as total_fee,count(id) as total_num,sum(total_charge_fee) as profit_fee").
						Where("channel_id=? AND depart_id=? AND mch_project_id=? AND trade_type=? AND trade_state IN ? and create_time>=? and create_time<=?",
							channelId, departId, projectId, tradeType, []string{constant.ORDER_TRADE_STATE_SUCCESS}, dayStartTime, endStartTime).Find(&incomeOrderResult)
					//查询
					database.DB.Table(aspOrder.TableName()).Select("count(id) as total_order_num").
						Where("channel_id=? AND depart_id=? AND mch_project_id=? AND trade_type=? AND trade_state IN ? and create_time>=? and create_time<=?",
							channelId, departId, projectId, tradeType, []string{constant.ORDER_TRADE_STATE_SUCCESS, constant.ORDER_TRADE_STATE_USERPAYING, constant.ORDER_TRADE_STATE_PAYERROR, constant.ORDER_TRADE_STATE_FAILED},
							dayStartTime, endStartTime).Find(&incomeOrderResult)

					//代收成功率
					if incomeOrderResult.TotalOrderNum > 0 {
						incomeSuccessRate = cast.ToInt((incomeOrderResult.TotalNum / incomeOrderResult.TotalOrderNum) * 100)
					}
					//代付订单
					database.DB.Table(aspPayout.TableName()).Select("sum(total_fee) as total_fee,count(id) as total_num,sum(total_charge_fee) as profit_fee").Where("channel_id=? AND depart_id=? AND mch_project_id=? AND trade_type=? AND trade_state IN ? and create_time>=? and create_time<=?", channelId, departId, projectId, tradeType, []string{constant.PAYOUT_TRADE_STATE_SUCCESS}, dayStartTime, endStartTime).Find(&payoutOrderResult)

					//查询成功计订单数
					database.DB.Table(aspPayout.TableName()).Select("count(id) as success_order_num").Where("channel_id=? AND depart_id=? AND mch_project_id=? AND trade_type=? AND trade_state IN ? and create_time>=? and create_time<=?", channelId, departId, projectId, tradeType, []string{constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS, constant.PAYOUT_TRADE_STATE_SUCCESS},
						dayStartTime, endStartTime).Find(&payoutOrderResult)

					//查询总计订单数
					database.DB.Table(aspPayout.TableName()).Select("count(id) as total_order_num").Where("channel_id=? AND depart_id=? AND mch_project_id=? AND trade_type=? AND trade_state IN ? and create_time>=? and create_time<=?", channelId, departId, projectId, tradeType, []string{constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING, constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED, constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS, constant.PAYOUT_TRADE_STATE_SUCCESS, constant.PAYOUT_TRADE_STATE_FAILED}, dayStartTime, endStartTime).Find(&payoutOrderResult)

					//代付成功率
					if payoutOrderResult.TotalOrderNum > 0 {
						payoutSuccessRate = cast.ToInt((payoutOrderResult.SuccessOrderNum / payoutOrderResult.TotalOrderNum) * 100)
					}

					netProfitFee := incomeOrderResult.ProfitFee + payoutOrderResult.ProfitFee
					netTotalFee := incomeOrderResult.TotalFee - incomeOrderResult.ProfitFee - payoutOrderResult.TotalFee - payoutOrderResult.ProfitFee

					//查询统计是否已存在
					statisticQueryWhere := "period = ? and date = ? and data_name = ? and channel_id = ? and depart_id = ? and mch_project_id = ? and trade_type = ?"
					orderStatistic, _ := aspOrderStatisticRepository.FindOne(database.NewSqlCondition().Where(statisticQueryWhere, dateType, curDate, constant.ORDER_STATISTIC_DATA_NAME_ALL, channelId, departId, projectId, tradeType))

					//根据项目id查询外部商户id
					projectCurrencyInfo, _ := aspMerchantProjectCurrencyRepository.FindOne(database.NewSqlCondition().Where("mch_project_id = ?", projectId))

					//根据渠道id获取渠道配置
					channelConfigInfo, _ := aspChannelConfigRepository.FindOne(database.NewSqlCondition().Where("id = ?", channelId))

					//如果历史没有统计，就插入
					if orderStatistic == nil {
						orderStatistic = &model.AspOrderStatistic{}
						orderStatistic.Period = dateType
						orderStatistic.Date = cast.ToUint64(curDate)
						orderStatistic.DataName = constant.ORDER_STATISTIC_DATA_NAME_ALL
						orderStatistic.ChannelId = channelId
						orderStatistic.DepartId = cast.ToUint(departId)
						orderStatistic.MchId = projectCurrencyInfo.MchId
						orderStatistic.CurrencyId = cast.ToInt(projectCurrencyInfo.CurrencyId)
						orderStatistic.MchProjectId = projectId
						orderStatistic.Adapter = channelConfigInfo.Name
						orderStatistic.TradeType = tradeType
						orderStatistic.TotalFee = incomeOrderResult.TotalFee
						orderStatistic.TotalNum = incomeOrderResult.TotalNum
						orderStatistic.ProfitFee = incomeOrderResult.ProfitFee
						orderStatistic.SuccessRate = incomeSuccessRate
						orderStatistic.PayoutTotalFee = payoutOrderResult.TotalFee
						orderStatistic.PayoutNum = payoutOrderResult.TotalNum
						orderStatistic.PayoutProfitFee = payoutOrderResult.ProfitFee
						orderStatistic.PayoutSuccessRate = payoutSuccessRate
						orderStatistic.NetProfitFee = netProfitFee
						orderStatistic.NetTotalFee = netTotalFee
						orderStatistic.CreateTime = cast.ToUint64(carbon.Now().Timestamp())
						orderStatistic.UpdateTime = cast.ToUint64(carbon.Now().Timestamp())
						aspOrderStatisticRepository.Create(orderStatistic)
						//fmt.Println("============")
						//fmt.Println(fmt.Sprintf("channelId:%d,departId:%d,projectId:%d,tradeType:%s,totalFee:%d,totalNum:%d,profitFee:%d,incomeSuccessRate:%d", channelId, departId, projectId, tradeType, incomeOrderResult.TotalFee, incomeOrderResult.TotalNum, incomeOrderResult.ProfitFee, incomeSuccessRate))
						//fmt.Println(fmt.Sprintf("income:%+v", incomeOrderResult))
						//fmt.Println(fmt.Sprintf("payout:%+v", payoutOrderResult))
						//fmt.Println(fmt.Sprintf("netProfitFee:%+v", netProfitFee))
						//fmt.Println(fmt.Sprintf("netTotalFee:%+v", netTotalFee))
						//fmt.Println(fmt.Sprintf("orderStatistic:%+v", orderStatistic))
						//fmt.Println("============")
						continue
					}

					updateData := make(map[string]any)
					updateData["total_fee"] = incomeOrderResult.TotalFee
					updateData["total_num"] = incomeOrderResult.TotalNum
					updateData["profit_fee"] = incomeOrderResult.ProfitFee
					updateData["success_rate"] = incomeSuccessRate
					updateData["payout_total_fee"] = payoutOrderResult.TotalFee
					updateData["payout_num"] = payoutOrderResult.TotalNum
					updateData["payout_profit_fee"] = payoutOrderResult.ProfitFee
					updateData["payout_success_rate"] = payoutSuccessRate
					updateData["net_profit_fee"] = cast.ToUint64(netProfitFee)
					updateData["net_total_fee"] = netTotalFee
					updateData["update_time"] = carbon.Now().Timestamp()
					aspOrderStatisticRepository.Updates(updateData, "id = ?", orderStatistic.Id)

					//fmt.Println("============")
					//fmt.Println(fmt.Sprintf("channelId:%d,departId:%d,projectId:%d,tradeType:%s,totalFee:%d,totalNum:%d,profitFee:%d,incomeSuccessRate:%d", channelId, departId, projectId, tradeType, incomeOrderResult.TotalFee, incomeOrderResult.TotalNum, incomeOrderResult.ProfitFee, incomeSuccessRate))
					//fmt.Println(fmt.Sprintf("income:%+v", incomeOrderResult))
					//fmt.Println(fmt.Sprintf("payout:%+v", payoutOrderResult))
					//fmt.Println(fmt.Sprintf("netProfitFee:%+v", netProfitFee))
					//fmt.Println(fmt.Sprintf("netTotalFee:%+v", netTotalFee))
					//fmt.Println(fmt.Sprintf("orderStatistic:%+v", orderStatistic))
					//fmt.Println("============")

				}
			}
		}
	}
}

// OrderDay 日订单统计
func (s *StatisticServer) OrderDay(day int) {
	s.OrderBase(1, day)
}

// OrderMonth 日订单统计
func (s *StatisticServer) OrderMonth(month int) {
	s.OrderBase(2, month)
}

type MerchantOrderStatisticResult struct {
	TotalFee              int64   `json:"total_fee,omitempty"`
	TotalNum              int     `json:"total_num,omitempty"`
	TotalOrderNum         int     `json:"total_order_num,omitempty"`
	ProfitFee             int64   `json:"profit_fee,omitempty"`
	SuccessRate           float64 `json:"success_rate,omitempty"`
	PayoutTotalFee        int64   `json:"payout_total_fee,omitempty"`
	PayoutNum             int     `json:"payout_num,omitempty"`
	PayoutSuccessOrderNum int     `json:"payout_success_order_num,omitempty"`
	PayoutTotalOrderNum   int     `json:"payout_total_order_num,omitempty"`
	PayoutProfitFee       int64   `json:"payout_profit_fee,omitempty"`
	PayoutSuccessRate     float64 `json:"payout_success_rate,omitempty"`
	NetProfitFee          int64   `json:"net_profit_fee,omitempty"`
	NetTotalFee           int64   `json:"net_total_fee,omitempty"`
}

func (s *StatisticServer) OrderBaseMerchant(dateType uint, date int) {
	var curDate int64 = 0
	var dayStartTime int64 = 0
	var endStartTime int64 = 0
	//日报表
	if dateType == 1 {
		curDate = carbon.Parse(carbon.Now().AddDays(date).Format("Y-m-d")).StartOfDay().Timestamp()
		dayStartTime = curDate
		endStartTime = carbon.Parse(carbon.Now().AddDays(date).Format("Y-m-d")).EndOfDay().Timestamp()
	} else {
		curDate = carbon.Parse(carbon.Now().AddMonths(date).Format("Y-m-d H:i:s")).StartOfMonth().Timestamp()
		dayStartTime = curDate
		endStartTime = carbon.Parse(carbon.Now().AddMonths(date).Format("Y-m-d H:i:s")).EndOfMonth().Timestamp()
	}
	logger.ApiInfo(s.LogFileName, s.RequestId, "OrderBaseMerchant", zap.Any("dateType", dateType), zap.Any("date", date), zap.Any("dayStartTime", dayStartTime), zap.Any("endStartTime", endStartTime))
	//获取内部商户列表
	aspMerchantList, _ := repository.NewRepository[*model.AspMerchant](database.DB, goRedis.Redis).Find(database.NewSqlCondition())
	aspOrderStatisticRepository := repository.NewRepository[*model.AspOrderStatistic](database.DB, goRedis.Redis)
	//获取币种列表
	aspCurrencyList, _ := repository.NewRepository[*model.AspCurrency](database.DB, goRedis.Redis).Find(database.NewSqlCondition())

	//订单统计
	aspOrderStatistic := model.AspOrderStatistic{}
	aspOrder := model.AspOrder{}
	aspPayout := model.AspPayout{}

	for _, merchant := range aspMerchantList {

		for _, currency := range aspCurrencyList {

			currencyId := currency.Id
			incomeSuccessRate := 0 //代收成功率（0-100 之间的值）
			payoutSuccessRate := 0 //代付成功率（0-100 之间的值）
			mchId := merchant.Id

			merchantOrderStatisticResult := MerchantOrderStatisticResult{}
			//代付订单
			database.DB.Table(aspOrderStatistic.TableName()).Select("sum(total_fee) as total_fee,sum(total_num) as total_num,sum(profit_fee) as profit_fee,sum(payout_total_fee) as payout_total_fee,sum(payout_num) as payout_num,sum(payout_profit_fee) as payout_profit_fee,sum(net_profit_fee) as net_profit_fee,sum(net_total_fee) as net_total_fee").Where("date = ? and mch_id = ? and period = ? and data_name = ? and currency_id = ?", curDate, mchId, dateType, constant.ORDER_STATISTIC_DATA_NAME_MCH, currencyId).Find(&merchantOrderStatisticResult)

			//查询代收总计订单数
			database.DB.Table(aspOrder.TableName()).Select("count(id) as total_order_num").
				Where("mch_id=? AND currency_id = ? AND trade_state IN ? and create_time>=? and create_time<=?",
					mchId, currencyId, []string{constant.ORDER_TRADE_STATE_SUCCESS, constant.ORDER_TRADE_STATE_USERPAYING, constant.ORDER_TRADE_STATE_PAYERROR, constant.ORDER_TRADE_STATE_FAILED},
					dayStartTime, endStartTime).Find(&merchantOrderStatisticResult)
			//代收成功率
			if merchantOrderStatisticResult.TotalOrderNum > 0 {
				incomeSuccessRate = cast.ToInt((merchantOrderStatisticResult.TotalNum / merchantOrderStatisticResult.TotalOrderNum) * 100)
			}
			//查询代付成功计订单数
			database.DB.Table(aspPayout.TableName()).Select("count(id) as payout_success_order_num").Where("mch_id=? AND currency_id = ?  AND trade_state IN ? and create_time>=? and create_time<=?", mchId, currencyId, []string{constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS, constant.PAYOUT_TRADE_STATE_SUCCESS},
				dayStartTime, endStartTime).Find(&merchantOrderStatisticResult)

			//查询代付总计订单数
			database.DB.Table(aspPayout.TableName()).Select("count(id) as payout_total_order_num").Where("mch_id=? AND currency_id = ? AND trade_state IN ? and create_time>=? and create_time<=?", mchId, currencyId, []string{constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING, constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED, constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS, constant.PAYOUT_TRADE_STATE_SUCCESS, constant.PAYOUT_TRADE_STATE_FAILED}, dayStartTime, endStartTime).Find(&merchantOrderStatisticResult)
			//代付成功率
			if merchantOrderStatisticResult.PayoutTotalOrderNum > 0 {
				payoutSuccessRate = cast.ToInt((merchantOrderStatisticResult.PayoutSuccessOrderNum / merchantOrderStatisticResult.PayoutTotalOrderNum) * 100)
			}

			orderStatistic, _ := aspOrderStatisticRepository.FindOne(database.NewSqlCondition().Where("period = ? and date = ? and data_name = ? and mch_id = ? AND currency_id = ?", 1, curDate, "MCH", mchId, currencyId))

			netProfitFee := merchantOrderStatisticResult.ProfitFee + merchantOrderStatisticResult.PayoutProfitFee
			netTotalFee := merchantOrderStatisticResult.TotalFee - merchantOrderStatisticResult.ProfitFee - merchantOrderStatisticResult.PayoutTotalFee - merchantOrderStatisticResult.PayoutProfitFee

			if orderStatistic == nil {
				orderStatistic = &model.AspOrderStatistic{}
				orderStatistic.Period = dateType
				orderStatistic.Date = cast.ToUint64(curDate)
				orderStatistic.DataName = constant.ORDER_STATISTIC_DATA_NAME_MCH
				orderStatistic.ChannelId = 0
				orderStatistic.DepartId = 0
				orderStatistic.MchId = mchId
				orderStatistic.CurrencyId = currencyId
				orderStatistic.TotalFee = merchantOrderStatisticResult.TotalFee
				orderStatistic.TotalNum = merchantOrderStatisticResult.TotalNum
				orderStatistic.ProfitFee = merchantOrderStatisticResult.ProfitFee
				orderStatistic.SuccessRate = incomeSuccessRate
				orderStatistic.PayoutTotalFee = merchantOrderStatisticResult.PayoutTotalFee
				orderStatistic.PayoutNum = merchantOrderStatisticResult.PayoutNum
				orderStatistic.PayoutProfitFee = merchantOrderStatisticResult.PayoutProfitFee
				orderStatistic.PayoutSuccessRate = payoutSuccessRate
				orderStatistic.NetProfitFee = netProfitFee
				orderStatistic.NetTotalFee = netTotalFee
				orderStatistic.PayoutSuccessRate = payoutSuccessRate
				orderStatistic.SuccessRate = incomeSuccessRate
				orderStatistic.CreateTime = uint64(carbon.Now().Timestamp())
				orderStatistic.UpdateTime = uint64(carbon.Now().Timestamp())
				aspOrderStatisticRepository.Create(orderStatistic)
				continue
			}

			updateData := make(map[string]any)
			updateData["total_fee"] = merchantOrderStatisticResult.TotalFee
			updateData["total_num"] = merchantOrderStatisticResult.TotalNum
			updateData["profit_fee"] = merchantOrderStatisticResult.ProfitFee
			updateData["success_rate"] = incomeSuccessRate
			updateData["payout_total_fee"] = merchantOrderStatisticResult.PayoutTotalFee
			updateData["payout_num"] = merchantOrderStatisticResult.PayoutNum
			updateData["payout_profit_fee"] = merchantOrderStatisticResult.PayoutProfitFee
			updateData["payout_success_rate"] = payoutSuccessRate
			updateData["net_profit_fee"] = cast.ToUint64(netProfitFee)
			updateData["net_total_fee"] = netTotalFee
			updateData["update_time"] = carbon.Now().Timestamp()
			aspOrderStatisticRepository.Updates(updateData, "id = ?", orderStatistic.Id)

		}

	}
}

// OrderDayMerchant 外部商户每日统计
func (s *StatisticServer) OrderDayMerchant(day int) {
	s.OrderBaseMerchant(1, day)
}

// OrderMonthMerchant 外部商户每月统计
func (s *StatisticServer) OrderMonthMerchant(month int) {
	s.OrderBaseMerchant(2, month)
}
