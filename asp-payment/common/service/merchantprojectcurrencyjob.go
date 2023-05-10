package service

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"github.com/golang-module/carbon/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type MerchantProjectCurrencyJobServer struct {
	*Service
}

func NewSimpleMerchantProjectCurrencyJobServer() *MerchantProjectCurrencyJobServer {
	return &MerchantProjectCurrencyJobServer{&Service{LogFileName: constant.MerchantProjectCurrencyLogFileName}}
}

// SyncMerchantProjectCurrencyJob 更新cp账户待结算余额 转到 可用余额中
func (s *MerchantProjectCurrencyJobServer) SyncMerchantProjectCurrencyJob() {
	// 查询到所有的cp项目余额 list
	merchantProjectServer := NewSimpleMerchantProjectServer()
	merchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, "")
	// 循环执行更新 应该具备事务功能
	merchantProjectCurrencyServerList, err := merchantProjectServer.GetMerchantProjectCurrencyList()
	if err != nil {
		logger.ApiError(s.LogFileName, s.RequestId, "merchantProjectServer.GetMerchantProjectCurrencyList err: ", zap.Error(err))
	}

	timeBegin := carbon.Parse(carbon.Now().AddDays(-1).Format("Y-m-d")).StartOfDay().Timestamp()
	//timeBegin := carbon.Parse(carbon.Now().Format("Y-m-d")).StartOfDay().Timestamp()
	timeEnd := timeBegin + 86399

	timeDay := timeBegin + 86400

	//var timeBegin int64
	//var timeEnd int64
	//
	//timeBegin = 1669019212
	//timeEnd = 1669107281

	changeAvailableTotalFee := 0

	for mchProjectId, merchantProjectCurrencyInfo := range merchantProjectCurrencyServerList {
		// 余额改怎么计算 定时脚本更新的订单的结果，一天cp项目总的待结算余额 实时去查询并循环中去更新
		merchantProjectCapitalFlowItem, errItem := merchantProjectServer.GetMerchantProjectCapitalFlowItem(constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYORDER, mchProjectId, timeBegin, timeEnd)
		if errItem != nil {
			logger.ApiError(s.LogFileName, s.RequestId, "merchantProjectServer.GetMerchantProjectCapitalFlowItem err: ")
			continue
		}
		merchantProjectTotalFee := merchantProjectCapitalFlowItem[0].TotalFee
		// 写入 资金待结算余额转可用余额日结明细表
		var aspMerchantProjectTransfersDayFlow model.AspMerchantProjectTransfersDayFlow
		aspMerchantProjectTransfersDayFlow.MchId = merchantProjectCurrencyInfo.MchId
		aspMerchantProjectTransfersDayFlow.MchProjectId = cast.ToInt(merchantProjectCurrencyInfo.MchProjectId)
		aspMerchantProjectTransfersDayFlow.MchProjectCurrencyId = cast.ToInt(merchantProjectCurrencyInfo.CurrencyId)
		aspMerchantProjectTransfersDayFlow.Currency = merchantProjectCurrencyInfo.Currency
		aspMerchantProjectTransfersDayFlow.TotalFee = merchantProjectTotalFee
		aspMerchantProjectTransfersDayFlow.Day = timeDay
		aspMerchantProjectTransfersDayFlow.Remark = constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_TO_AVAILABLE_TOTAL_FEE
		aspMerchantProjectTransfersDayFlow.CreateTime = carbon.Now().Timestamp()
		aspMerchantProjectTransfersDayFlow.UpdateTime = carbon.Now().Timestamp()
		errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
			errCreate := tx.Create(&aspMerchantProjectTransfersDayFlow).Error
			if errCreate != nil {
				logger.ApiWarn(s.LogFileName, s.RequestId, "AspMerchantProjectTransfersDayFlow Insert: ", zap.Error(errCreate))
				return errCreate
			}
			// 每天的总的代收成功的金额

			if merchantProjectTotalFee > 0 {
				// 查询到对应的cp项目余额 增加可用余额 减去待结算余额
				// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
				changeAvailableTotalFee = cast.ToInt(merchantProjectTotalFee)
				errChange := merchantProjectRepository.ChangeTotalFeeToAvailableTotalFee(mchProjectId, changeAvailableTotalFee, aspMerchantProjectTransfersDayFlow.Id, tx)
				if errChange != nil {
					logger.ApiError(s.LogFileName, s.RequestId, "MerchantProjectCurrent TO_AVAILABLE_TOTAL_FEE err")
					return errChange
				}
			}
			return nil
		})
		if errTrans != nil {
			logger.ApiInfo(s.LogFileName, s.RequestId, "merchantProjectTotalFee TO_AVAILABLE_TOTAL_FEE", zap.Int64("merchantProjectTotalFee", merchantProjectTotalFee), zap.Any("merchantProjectCurrencyBefore", merchantProjectCurrencyInfo))
			time.Sleep(2 * 100 * time.Millisecond)
			continue
		}

		afterCurrencyInfo, errProCur := merchantProjectServer.GetMerchantProjectCurrencyInfo(merchantProjectCurrencyInfo.Currency, mchProjectId)
		if errProCur != nil {
			logger.ApiError(s.LogFileName, s.RequestId, "merchantProjectServer.GetMerchantProjectInfoWithNoCache err: ", zap.Error(errProCur))
		}
		logger.ApiInfo(s.LogFileName, s.RequestId, "merchantProjectTotalFee TO_AVAILABLE_TOTAL_FEE", zap.Int64("merchantProjectTotalFee", merchantProjectTotalFee), zap.Any("merchantProjectCurrencyBefore", merchantProjectCurrencyInfo), zap.Any("merchantProjectCurrencyAfter", afterCurrencyInfo))
		// 延时 0.2秒
		time.Sleep(2 * 100 * time.Millisecond)
	}
	// 记录操作日志

}
