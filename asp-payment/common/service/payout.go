package service

import (
	"asp-payment/api-server/req"
	"asp-payment/api-server/rsp"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goRedis"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	reqCommon "asp-payment/common/req"
	thirdParty "asp-payment/common/service/supplier"
	"asp-payment/common/service/supplier/interfaces"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type PayoutServer struct {
	*Service
}

func NewPayoutServer(c *fiber.Ctx) *PayoutServer {
	return &PayoutServer{Service: NewService(c, constant.PayoutServerLogFileName)}
}

func NewSimplePayoutServer() *PayoutServer {
	return &PayoutServer{&Service{LogFileName: constant.PayoutServerLogFileName}}
}

var (
	payoutPayment        = "sunny.payout"
	tradeTypeBeneficiary = "fynzonpay.BENEFICIARY"
	tradeTypeUpiValidate = "haodapay.UPIVALIDATE"
	tradeTypePayoutUpi   = "haodapay.PAYOUTUPI"
)

func (s *PayoutServer) Payout() error {

	// 贯穿支付所需要的参数
	reqAspPayout := new(req.AspPayout)
	if errBody := s.C.BodyParser(reqAspPayout); errBody != nil {
		return s.Error(appError.NewError(errBody.Error()))
	}
	reqAspPayout.PayoutCode = "00"
	// logger.ApiInfo(s.LogFileName, s.RequestId, "xxx")
	if err := s.DealUnifiedPayoutParams(reqAspPayout, s.Head); err != nil {
		return s.Error(appError.NewError(err.Error()))
	}
	lockName := goRedis.GetKey(fmt.Sprintf("pay:payout:%s_%s", s.Head.AppId, reqAspPayout.UserID))
	flag := goRedis.Lock(lockName)
	if !flag {
		return appError.IsWaitErrCode
	}
	defer goRedis.UnLock(lockName)

	// 验证银行卡编码 如果支付类型是 银行卡
	if reqAspPayout.PayType == constant.PARAMS_PAY_TYPE_BANK {
		BankCategoryInfo := s.getBankCategoryInfo(reqAspPayout.BankCode)
		if BankCategoryInfo == nil {
			return appError.MissNotFoundErrCode.FormatMessage(constant.MissBankCategoryErrMsg)
		}
	}
	merchantProjectServer := NewMerchantProjectServer(s.C)
	// appError.MissMerchantProjectNotFoundErr.FormatMessage("xxxx")
	merchantProjectInfo, pErr := merchantProjectServer.GetMerchantProjectInfo(s.Head.AppId)
	if pErr != nil {
		return s.Error(pErr)
	}
	AspIdInfo, mErr := merchantProjectServer.GetIdInfo()
	if mErr != nil {
		return s.Error(mErr)
	}
	// 验证签名 可以使用 默认的 util 来实现 验证签名
	if err := s.CheckSignPayout(reqAspPayout, AspIdInfo.Key); err != nil {
		return err
	}
	clientIp := s.C.IP()
	// 提现的验证 如果验证失败 则提示 验证收款订单号
	if err := s.VerifyPayoutBefore(merchantProjectInfo, reqAspPayout, clientIp); err != nil {
		return s.Error(err)
	}

	//验证黑名单列表 bankcard是客户端的，后台数据库是bankcode
	if err := s.VerifyApiBlack(reqAspPayout.CustomerPhone, reqAspPayout.DeviceInfo, reqAspPayout.BankCard, ""); err != nil {
		return s.Error(err)
	}

	// 查询到 商户的可用的支付渠道 过程中有赋值给 AspChannelDepartTradeType
	aspChannelDepartTradeType, err := NewDepartServer(s.C).ChoosePayoutChannelDepartTradeType(reqAspPayout)
	if err != nil {
		return s.Error(err)
	}

	paymentArr := strings.Split(aspChannelDepartTradeType.Payment, ".")
	if len(paymentArr) != 2 {
		return s.Error(appError.MissNotFoundErrCode.FormatMessage(constant.MissChannelDepartPaymentParamErrMsg))
	}
	reqAspPayout.PayoutParams.Adapter = paymentArr[1]

	// 代付 fynzonpay 渠道需要验证 受益人信息
	if paymentArr[1] == constant.TradeTypeFynzonPay {
		// 需要参数 接收参数 过程中有对 d.PayoutParams.BenefiaryId 更改
		if beErr := s.GetPayoutBenefiary(reqAspPayout); beErr != nil {
			return s.Error(beErr)
		}
	}

	err = merchantProjectServer.MerchantProjectUserInsertOrUpdate(reqAspPayout.UserID, merchantProjectInfo)
	if err != nil {
		return s.Error(appError.CodeUserNotExist) // 用户处理异常，请重试
	}

	// 查询到cp 项目的配置
	merchantProjectConfigInfo, aErr := merchantProjectServer.GetMerchantAccountConfigInfo()
	if aErr != nil {
		return s.Error(aErr)
	}

	merchantProjectCurrencyInfo, cErr := merchantProjectServer.GetMerchantProjectCurrencyInfo(reqAspPayout.OrderCurrency, merchantProjectInfo.Id)
	if cErr != nil {
		return s.Error(cErr)
	}
	availableTotalFee := cast.ToInt(merchantProjectCurrencyInfo.AvailableTotalFee)
	// 提现的是否可以提交到上游去 如果验证失败 则提示
	if err = s.VerifyPayoutAfter(reqAspPayout, availableTotalFee, merchantProjectConfigInfo); err != nil {
		return s.Error(err)
	}

	// 生成订单记录 提前预插入订单数据
	payoutInfo, dErr := s.Insert(reqAspPayout, merchantProjectInfo, merchantProjectConfigInfo, merchantProjectCurrencyInfo, aspChannelDepartTradeType)
	if dErr != nil {
		return s.Error(appError.CodeInsertErr) // 写入数据失败，请重试
	}
	// 初始化操作金额
	changeAvailableTotalFee := 0
	changeTotalFee := 0
	changeFreezeFee := 0
	changeAmount := reqAspPayout.OrderAmount + payoutInfo.ChargeFee + payoutInfo.FixedAmount
	//fmt.Println("reqAspPayout.PayoutCode: ", reqAspPayout.PayoutCode)
	// 商户余额不足 或者是 金额超过了设置的范围了
	if reqAspPayout.PayoutCode != "00" {
		scanSuccessData := rsp.GenerateDevPayoutSuccessData(payoutInfo)
		return s.Success(scanSuccessData)
	}

	// 测试环境 只有这一个 id 走上游 测试环境是不能请求到上游的
	if config.IsTestEnv() && s.Head.AppId != constant.IGNORE_MERCHANT_PROJECT_ID {
		// 更新操作 修改订单状态即可
		payoutSuccessUpdate := make(map[string]interface{})
		payoutSuccessUpdate["transaction_id"] = goutils.RandomString(20)
		payoutSuccessUpdate["cash_fee"] = payoutInfo.CashFee
		payoutSuccessUpdate["cash_fee_type"] = payoutInfo.CashFeeType
		payoutSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
		payoutSuccessUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_SUCCESS

		MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
		changeAvailableTotalFee = -changeAmount
		changeTotalFee = 0
		changeFreezeFee = 0

		// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
		errChange := MerchantProjectRepository.ChangeMerchantProjectCurrentByTest(payoutInfo.MchProjectId, changeAvailableTotalFee, changeTotalFee, changeFreezeFee, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_FREEZE, payoutInfo.Id, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_SUCCESS, payoutSuccessUpdate)
		if errChange != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("PayoutDev") // 更新代付订单错误 请重新提交
		}
		payoutInfoParam := make(map[string]interface{})
		payoutInfoParam["id"] = cast.ToString(payoutInfo.Id)
		newPayoutInfo, errInfo := s.GetPayoutInfo(payoutInfoParam, nil)
		if errInfo != nil {
			return appError.CodeUnknown // 服务器开小差了！
		}
		_ = NewSendQueueServer().SendNotifyQueue(newPayoutInfo.Sn)
		scanSuccessData := rsp.GenerateDevPayoutSuccessData(newPayoutInfo)
		return s.Success(scanSuccessData)
	}
	// 冻结可用余额  可用余额减去  冻结金额增加  预扣金额记录新增 更新订单状态为冻结成功

	MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
	db := database.DB
	payoutInfo, err = s.FreezeAndUpdate(MerchantProjectRepository, payoutInfo, changeAmount, db)
	if err != nil {
		return s.Error(err)
	}
	// 查询渠道内部商户配置信息 提供添加受益人 验证 upi 方法使用
	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(payoutInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(payoutInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, errChCo := NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if errChCo != nil {
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return s.Error((&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartNotFoundErrMsg))
	}
	currencyId := cast.ToInt(merchantProjectCurrencyInfo.CurrencyId)

	// 代付 fynzonpay 渠道需要验证 受益人信息
	if paymentArr[1] == constant.TradeTypeFynzonPay {
		// 当受益人id 为空则需要请求到上游
		if reqAspPayout.BenefiaryId == "" {
			beneficiaryUp, errUpbe := s.RequestUpstreamAddBeneficiary(reqAspPayout, channelDepartInfo)
			if errUpbe != nil {
				return s.Error(errUpbe)
			}
			payoutInfo.BeneficiaryId = beneficiaryUp.BenefiaryId
			reqAspPayout.BenefiaryId = beneficiaryUp.BenefiaryId
			// 插入到受益人表
			_, errInsBe := s.InsertBeneficiar(payoutInfo, merchantProjectInfo, channelDepartInfo, currencyId)
			if errInsBe != nil {
				return s.Error(errInsBe)
			}
			// 更新订单受益人id
			beParams := make(map[string]interface{})
			beParams["beneficiary_id"] = beneficiaryUp.BenefiaryId
			resultPayoutUpdate := database.DB.Model(payoutInfo).Where("id = ?", payoutInfo.Id).Updates(beParams)
			if resultPayoutUpdate.Error != nil {
				logger.ApiWarn(s.LogFileName, s.RequestId, "payout.update() ", zap.Error(resultPayoutUpdate.Error))
				return s.Error(appError.NewError(resultPayoutUpdate.Error.Error()))
			}
			payoutInfo.BeneficiaryId = beneficiaryUp.BenefiaryId
		}
	}

	// 需要验证 upi 的 合法性
	if payoutInfo.PayType == constant.CASHIERDESK_PAY_TYPE_UPI {
		if upiVaErr := s.PayoutUpiValidate(reqAspPayout, payoutInfo, merchantProjectInfo, channelDepartInfo, currencyId); upiVaErr != nil {
			return s.Error(upiVaErr)
		}
	}

	// 执行请求上游 + 对应更新
	scanSuccessData, sErr := s.RequestUpstream(payoutInfo)

	if sErr != nil {
		// 是否上游返回异常 + 更新订单 + 释放可用金额
		errUpdateFailed := s.PayoutCreateUpdateFailed(MerchantProjectRepository, scanSuccessData, payoutInfo, changeAmount, sErr, db)
		if errUpdateFailed != nil {
			return errUpdateFailed
		}
		if sErr.Code == appError.CodeSupplierChannelErrCode.Code {
			_ = NewSendQueueServer().ManualSendNotifyQueue(payoutInfo.Sn)
		}
		return s.Error(sErr)
	}

	var isSendQueue bool
	isSendQueue = false
	if scanSuccessData.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
		isSendQueue = true
	}

	errTrans := database.DB.Transaction(func(tx *gorm.DB) error {

		errUpdate := s.PayoutCreateUpdate(MerchantProjectRepository, payoutInfo, changeAmount, scanSuccessData, tx)
		if errUpdate != nil {
			return errUpdate
		}
		return nil
	})

	if errTrans != nil {
		return s.Error(appError.NewError(errTrans.Error()))
	}
	// 发送回调到redis
	if isSendQueue == true {
		_ = NewSendQueueServer().SendNotifyQueue(payoutInfo.Sn)
	}
	// 事务前查询的到订单状态 不可见 所以 commit 后 再查询一次
	newPayoutInfoParam := make(map[string]interface{})
	newPayoutInfoParam["id"] = payoutInfo.Id
	newPayoutInfo, errNewPayout := s.GetPayoutInfo(newPayoutInfoParam, nil)
	if errNewPayout != nil {
		return s.Error(appError.CodeUnknown) // 服务器开小差了！
	}

	payoutSuccessData := rsp.GeneratePayoutApplySuccessData(newPayoutInfo)
	return s.Success(payoutSuccessData)
}

func (s *PayoutServer) PayoutAuditApi() error {
	lockName := goRedis.GetKey(fmt.Sprintf("pay:audit:%s", s.Head.AppId))
	if !goRedis.Lock(lockName) {
		return appError.IsWaitErrCode
	}
	defer goRedis.UnLock(lockName)
	// 贯穿支付所需要的参数
	reqAspPayoutAudit := new(req.AspPayoutAudit)
	if err := s.C.BodyParser(reqAspPayoutAudit); err != nil {
		return err
	}
	if err := s.DealUnifiedPayoutAuditParams(reqAspPayoutAudit); err != nil {
		return err
	}

	payoutSuccessData, errDo := s.PayoutAudit(reqAspPayoutAudit)
	// 更新已完成状态 修改审核状态
	params := make(map[string]interface{})
	params["is_checkout"] = constant.PAYOUT_IS_CHECKOUT_SUCCESS // 是否审核完成状态
	payoutId := reqAspPayoutAudit.PayoutID
	_, errAuditUpdate := s.SolvePayoutSuccess(payoutId, params, "payoutAuditUpdateCheckout", nil)
	if errAuditUpdate != nil {
		return errAuditUpdate
	}

	if errDo != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "PayoutAudit ", zap.Any("errDo", errDo))
		return errDo
	}
	return s.Success(payoutSuccessData)
}

// PayoutAuditJob 审核脚本执行
func (s *PayoutServer) PayoutAuditJob(payoutAuditInfo *reqCommon.AspPayoutAuditJob) error {
	// 验证参数
	if err := s.DealUnifiedPayoutAuditJobParams(payoutAuditInfo); err != nil {
		return err
	}
	lockName := goRedis.GetKey(fmt.Sprintf("pay:auditjob:%s", payoutAuditInfo.Id))
	flag := goRedis.Lock(lockName)
	if !flag {
		return appError.IsWaitErrCode
	}
	defer goRedis.UnLock(lockName)
	// 贯穿支付所需要的参数
	reqAspPayoutAudit := new(req.AspPayoutAudit)
	reqAspPayoutAudit.OperationID = payoutAuditInfo.OperationID
	reqAspPayoutAudit.PayoutID = payoutAuditInfo.Id
	reqAspPayoutAudit.Action = payoutAuditInfo.Status
	_, err := s.PayoutAudit(reqAspPayoutAudit)

	// 更新已完成状态 修改审核状态
	params := make(map[string]interface{})
	params["is_checkout"] = constant.PAYOUT_IS_CHECKOUT_SUCCESS // 是否审核完成状态
	payoutId := cast.ToInt(payoutAuditInfo.Id)
	_, errAuditUpdate := s.SolvePayoutSuccess(payoutId, params, "payoutAuditUpdateCheckout", nil)
	if errAuditUpdate != nil {
		return errAuditUpdate
	}

	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "PayoutAudit ", zap.Any("err", err))
		return err
	}
	return nil
}

// PayoutAudit 审核
func (s *PayoutServer) PayoutAudit(reqAspPayoutAudit *req.AspPayoutAudit) (*rsp.PayoutSuccessData, *appError.Error) {
	payoutInfoParam := make(map[string]interface{})
	payoutInfoParam["id"] = reqAspPayoutAudit.PayoutID
	payoutInfo, pErr := s.GetPayoutInfo(payoutInfoParam, nil)

	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "s.GetPayoutInfo  ", zap.Error(pErr))
		return nil, appError.OrderNotFoundErrCode // 订单不存在
	}

	// 如果是申请状态 则提交到上游 并且更新状态 其他的状态 则是异常的状态 记录为异常
	if !s.IsNeedPayoutAuditOperate(payoutInfo, reqAspPayoutAudit) {
		// 输出结果
		//payoutAuditSuccessData := rsp.GeneratePayoutApplySuccessData(payoutInfo)
		// 记录异常日志 来处理
		logger.ApiWarn(s.LogFileName, s.RequestId, "Payout TradeState error", zap.Any("payout", payoutInfo))
		return nil, appError.CodeInvalidParamErrCode
	}
	total := cast.ToInt(payoutInfo.TotalFee)
	changeAmount := total + payoutInfo.ChargeFee + payoutInfo.FixedAmount
	merchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
	var payoutCreateData *interfaces.ThirdPayoutCreateData
	var payoutQueryData *interfaces.ThirdPayoutQueryData
	db := database.DB
	// 需要请求上游的情况 通过的三种情况 执行请求到上游
	// 请求到上游后，如果有返回异常 需要更新订单 返回 为了事务的外面处理 commit 则方法没有封装一个具体的方法
	if reqAspPayoutAudit.Action == "pass" && payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_APPLY {
		// 请求上游前 不加事务 代付中有乐观锁的
		// 冻结可用余额  可用余额减去  冻结金额增加  预扣金额记录新增 更新订单状态为冻结成功
		payoutFreezeInfo, errFree := s.FreezeAndUpdate(merchantProjectRepository, payoutInfo, changeAmount, db)
		if errFree != nil {
			return nil, errFree
		}
		payoutInfo = payoutFreezeInfo
	}
	if reqAspPayoutAudit.Action == "pass" {
		if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_APPLY || payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS {
			// 执行请求上游
			payoutCreatePassData, errUpstreamCreate := s.RequestUpstream(payoutInfo)
			if errUpstreamCreate != nil {
				// 是否上游返回异常 + 更新订单 + 释放可用金额 + 事务
				errUpdateFailed := s.PayoutCreateUpdateFailed(merchantProjectRepository, payoutCreatePassData, payoutInfo, changeAmount, errUpstreamCreate, db)
				if errUpdateFailed != nil {
					return nil, errUpdateFailed
				}
				// 手动添加到队列，发送到下游
				if errUpstreamCreate.Code == appError.CodeSupplierChannelErrCode.Code {
					_ = NewSendQueueServer().ManualSendNotifyQueue(payoutInfo.Sn)
				}
				return nil, errUpstreamCreate
			}
			payoutCreateData = payoutCreatePassData
		}
		if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING {
			payoutQueryPassData, errUpstreamQuery := s.PayoutAuditUpstreamQuery(payoutInfo)
			if errUpstreamQuery != nil {
				// 是否上游返回异常 + 更新订单 + 释放可用金额 + 事务
				errUpdateFailed := s.PayoutQueryUpdateFailed(merchantProjectRepository, payoutQueryPassData, payoutInfo, errUpstreamQuery, changeAmount, db)
				if errUpdateFailed != nil {
					return nil, errUpdateFailed
				}
				if errUpstreamQuery.Code == appError.CodeSupplierInternalChannelParamsFailedErrCode.Code {
					_ = NewSendQueueServer().ManualSendNotifyQueue(payoutInfo.Sn)
				}
				return nil, errUpstreamQuery
			}
			payoutQueryData = payoutQueryPassData
		}
		logger.ApiWarn(s.LogFileName, s.RequestId, "payoutAuditUpdateCheckout ", zap.Any("payoutCreateData", payoutCreateData), zap.Any("payoutQueryData", payoutQueryData))
	}
	// 是否发送回调
	var isSendQueue bool
	isSendQueue = false
	if payoutCreateData != nil && payoutCreateData.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
		isSendQueue = true
	}
	if payoutQueryData != nil && payoutQueryData.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
		isSendQueue = true
	}

	var payoutSuccessData *rsp.PayoutSuccessData
	// 不需要请求上游的情况
	// 审核拒绝的逻辑
	errTrans := database.DB.Transaction(func(tx *gorm.DB) error {
		if reqAspPayoutAudit.Action == "return" {
			errReturn := s.ReturnAudit(payoutInfo, changeAmount, merchantProjectRepository, tx)
			if errReturn != nil {
				return errReturn
			}
		}
		// 审核通过的逻辑
		if reqAspPayoutAudit.Action == "pass" {
			errPass := s.PassAudit(payoutInfo, changeAmount, merchantProjectRepository, payoutCreateData, payoutQueryData, tx)
			if errPass != nil {
				return errPass
			}
		}
		// 需要记录审核的流水的记录
		errIst := s.InsertAuditRecord(reqAspPayoutAudit, payoutInfo, tx)
		if errIst != nil {
			return appError.CodeUnknown // 服务器开小差了！
		}

		return nil
	})

	if errTrans != nil {
		return nil, appError.NewError(errTrans.Error())
	}
	// 发送回调到redis
	if isSendQueue == true {
		_ = NewSendQueueServer().SendNotifyQueue(payoutInfo.Sn)
	}
	// 事务前查询的到订单状态 不可见 所以 commit 后 再查询一次
	newPayoutInfoParam := make(map[string]interface{})
	newPayoutInfoParam["id"] = reqAspPayoutAudit.PayoutID
	newPayoutInfo, errNewPayout := s.GetPayoutInfo(newPayoutInfoParam, nil)
	if errNewPayout != nil {
		return nil, appError.CodeUnknown // 服务器开小差了！
	}
	payoutSuccessData = rsp.GeneratePayoutApplySuccessData(newPayoutInfo)
	return payoutSuccessData, nil
}

// PayoutAuditUpstreamQuery 审核代收 需要走上游查询的情况
func (s *PayoutServer) PayoutAuditUpstreamQuery(payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutQueryData, *appError.Error) {
	var payoutQueryData *interfaces.ThirdPayoutQueryData
	if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING {
		channelDepartInfoParam := map[string]string{}
		channelDepartInfoParam["depart_id"] = cast.ToString(payoutInfo.DepartId)
		channelDepartInfoParam["channel_id"] = cast.ToString(payoutInfo.ChannelId)
		//查询账户在上游的配置
		channelDepartInfo, err5 := NewSimpleChannelConfigServer().GetChannelDepartInfo(channelDepartInfoParam)
		if err5 != nil {
			MissNotFoundErrCode := *appError.MissNotFoundErrCode
			return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartNotFoundErrMsg) // missing channel depart not found
		}
		// 查询代付状态 并更新
		scanSuccessQueryData, err4 := s.RequestUpstreamQueryPayout(payoutInfo, channelDepartInfo)
		if err4 != nil {
			return scanSuccessQueryData, err4
		}
		payoutQueryData = scanSuccessQueryData
	}
	return payoutQueryData, nil
}

func (s *PayoutServer) PassAudit(payoutInfo *model.AspPayout, changeAmount int, merchantProjectRepository *repository.MerchantProjectRepository, payoutCreateData *interfaces.ThirdPayoutCreateData, payoutQueryData *interfaces.ThirdPayoutQueryData, tx *gorm.DB) *appError.Error {
	if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_APPLY || payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS {
		errUpdate := s.PayoutCreateUpdate(merchantProjectRepository, payoutInfo, changeAmount, payoutCreateData, tx)
		if errUpdate != nil {
			return errUpdate
		}
	}
	if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING {
		_, errPayoutUpdate := s.PayoutQueryUpdate(payoutInfo, payoutQueryData, changeAmount, tx)
		if errPayoutUpdate != nil {
			return errPayoutUpdate
		}
	}
	if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
		errChannelFailed := s.AuditPassChannelFailed(payoutInfo, changeAmount, merchantProjectRepository, tx)
		if errChannelFailed != nil {
			return errChannelFailed
		}
	}
	if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
		errChannelSuccess := s.AuditPassChannelSuccess(payoutInfo, changeAmount, merchantProjectRepository, tx)
		if errChannelSuccess != nil {
			return errChannelSuccess
		}
	}
	return nil
}

func (s *PayoutServer) AuditPassChannelSuccess(payoutInfo *model.AspPayout, changeAmount int, merchantProjectRepository *repository.MerchantProjectRepository, tx *gorm.DB) *appError.Error {
	// 解冻 + 更新代付成功
	// 代付上游成功 解冻+代付成功  冻结金额减去  预扣金额记录新增  收支流水扣减
	payoutSuccessUpdate := make(map[string]interface{})
	payoutSuccessUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_SUCCESS
	payoutSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()

	_, errSuccess := s.SolvePayoutSuccess(payoutInfo.Id, payoutSuccessUpdate, "AuditPassChannelSuccess", tx)
	if errSuccess != nil {
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("AuditPassChannelSuccess") // 更新商户余额错误 请重新提交
	}
	// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
	err := merchantProjectRepository.PayoutOrderChannelSuccess(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
	if err != nil {
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_SUCCESS") // 更新商户余额错误 请重新提交
	}
	return nil
}

func (s *PayoutServer) AuditPassChannelFailed(payoutInfo *model.AspPayout, changeAmount int, merchantProjectRepository *repository.MerchantProjectRepository, tx *gorm.DB) *appError.Error {
	// 解冻 + 更新代付失败
	// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
	payoutFailedUpdate := make(map[string]interface{})
	payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FAILED

	_, errFailed := s.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "AuditPassChannelFailed", tx)
	if errFailed != nil {
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("AuditPassChannelFailed") // 更新商户余额错误 请重新提交
	}
	// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
	err := merchantProjectRepository.PayoutOrderChannelFailed(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
	if err != nil {
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
	}
	return nil
}

func (s *PayoutServer) ReturnAudit(payoutInfo *model.AspPayout, changeAmount int, merchantProjectRepository *repository.MerchantProjectRepository, tx *gorm.DB) *appError.Error {
	var err *appError.Error
	// 更新状态为 审核中
	params := make(map[string]interface{})
	if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_APPLY {
		// 更新为拒绝
		// 更新状态为 审核中
		params["trade_state"] = constant.PAYOUT_TRADE_STATE_RETURN // 审核失败
		params["finish_time"] = goutils.GetDateTimeUnix()          // 修改时间
		_, err = s.SolvePayoutSuccess(payoutInfo.Id, params, "ReturnAudit", tx)
		if err != nil {
			return appError.CodeUnknown // 服务器开小差了！
		}
		return nil
	}
	if payoutInfo.TradeState == constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS {
		// 解冻 + 更新代付取消
		// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
		payoutFailedUpdate := make(map[string]interface{})
		payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_REVOKE

		_, errFailed := s.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "ReturnAudit", tx)
		if errFailed != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("ReturnAudit") // 更新商户余额错误 请重新提交
		}
		// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
		err = merchantProjectRepository.PayoutOrderAuditReturn(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
		if err != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
		}
	}
	return nil
}

func (s *PayoutServer) PayoutQuery() error {
	// 贯穿支付所需要的参数
	reqPayoutQuery := new(req.AspPayoutQuery)

	if errBody := s.C.QueryParser(reqPayoutQuery); errBody != nil {
		return s.Error(appError.NewError(errBody.Error()))
	}

	if err := s.DealUnifiedPayoutQueryParams(reqPayoutQuery, s.Head); err != nil {
		return s.Error(err)
	}

	merchantProjectServer := NewMerchantProjectServer(s.C)
	// appError.MissMerchantProjectNotFoundErr.FormatMessage("xxxx")
	_, err := merchantProjectServer.GetMerchantProjectInfo(s.Head.AppId)
	if err != nil {
		return err
	}
	AspIdInfo, err := merchantProjectServer.GetIdInfo()
	if err != nil {
		return s.Error(err)
	}

	// 验签
	if err = s.CheckSignPayoutQuery(AspIdInfo.Key, reqPayoutQuery, s.Head.Timestamp); err != nil {
		return err
	}
	// 1. 首先查询订单
	payoutInfoParam := make(map[string]interface{})
	payoutInfoParam["sn"] = reqPayoutQuery.Sn
	payoutInfo, pErr := s.GetPayoutInfo(payoutInfoParam, nil)

	if pErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "logic.GetOrderInfo  ", zap.Error(pErr))
		return s.Error(appError.PayoutNotFoundErrCode) // 代付订单不存在
	}

	params := map[string]string{
		"is_call_upstream": reqPayoutQuery.IsCallUpstream,
	}
	// 2. 根据是否需要走上游
	isCallUpstream := s.IsNeedPayoutQueryUpstream(payoutInfo, params)

	// 3. 如果需要走上游 拼接订单渠道中
	//var querySuccessData *req.QuerySuccessData
	if isCallUpstream == false || (config.IsTestEnv() && s.Head.AppId != constant.IGNORE_MERCHANT_PROJECT_ID) {
		scanSuccessData := rsp.GenerateDevPayoutQuerySuccessData(payoutInfo)
		return s.Success(scanSuccessData)
	}
	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(payoutInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(payoutInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, cErr := NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if cErr != nil {
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartNotFoundErrMsg) // missing channel depart not found
	}
	scanQueryData, qErr := s.RequestUpstreamQueryPayout(payoutInfo, channelDepartInfo)
	newPayoutInfo := &model.AspPayout{}
	if qErr != nil {
		total := cast.ToInt(payoutInfo.TotalFee)
		changeAmount := total + payoutInfo.ChargeFee + payoutInfo.FixedAmount
		merchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)
		// 是否上游返回异常 + 更新订单 + 释放可用金额 + 事务
		errUpdateFailed := s.PayoutQueryUpdateFailed(merchantProjectRepository, scanQueryData, payoutInfo, qErr, changeAmount, database.DB)
		if errUpdateFailed != nil {
			return errUpdateFailed
		}
		if qErr.Code == appError.CodeSupplierInternalChannelParamsFailedErrCode.Code {
			_ = NewSendQueueServer().ManualSendNotifyQueue(payoutInfo.Sn)
		}
		return qErr
	}
	var isSendQueue bool
	isSendQueue = false
	if scanQueryData.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
		isSendQueue = true
	}
	errTrans := database.DB.Transaction(func(tx *gorm.DB) error {

		totalFee := cast.ToInt(payoutInfo.TotalFee)
		changeAmount := totalFee + payoutInfo.ChargeFee + payoutInfo.FixedAmount
		payoutUpdateInfo, errPayoutUpdate := s.PayoutQueryUpdate(payoutInfo, scanQueryData, changeAmount, tx)
		if errPayoutUpdate != nil {
			return errPayoutUpdate
		}
		newPayoutInfo = payoutUpdateInfo
		return nil
	})
	if errTrans != nil {
		return s.Error(appError.NewError(errTrans.Error()))
	}
	// 发送回调到redis
	if isSendQueue == true {
		_ = NewSendQueueServer().SendNotifyQueue(payoutInfo.Sn)
	}
	// 统一返回参数 转换
	scanSuccessData := rsp.GeneratePayoutQuerySuccessData(newPayoutInfo)
	return s.Success(scanSuccessData)
}

// RequestUpstreamAddBeneficiary 请求上游执行添加受益人 + 写入到受益人表
func (s *PayoutServer) RequestUpstreamAddBeneficiary(reqPayout *req.AspPayout, channelDepartInfo *model.AspChannelDepartConfig) (*interfaces.ThirdAddBeneficiary, *appError.Error) {
	// 贯穿支付所需要的参数
	reqBeneficiary := new(req.AspBeneficiary)
	reqBeneficiary.CustomerName = reqPayout.CustomerName
	reqBeneficiary.CustomerPhone = reqPayout.CustomerPhone
	reqBeneficiary.CustomerEmail = reqPayout.CustomerEmail
	reqBeneficiary.Ifsc = reqPayout.Ifsc
	reqBeneficiary.BankCard = reqPayout.BankCard
	reqBeneficiary.BankCode = reqPayout.BankCode

	supplier := thirdParty.GetPaySupplierByCode(tradeTypeBeneficiary)
	if supplier == nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "thirdParty.GetPaySupplierByCode error ", zap.String("err", "supplier == nil"))
		return nil, appError.ChannelDepartTradeTypeNotFoundErrCode // 渠道信息不存在
	}

	scanData, err := supplier.AddBeneficiary(s.RequestId, s.C.IP(), channelDepartInfo, reqBeneficiary)
	if err != nil {
		return nil, err
	}
	// 渠道映射系统的统一返回字符串 在 appError中有定义的 map
	supplierStr := constant.TradeTypeFynzonPay + "_" + scanData.Code
	supplierError, ok := appError.SupplierErrorMap[supplierStr]
	if !ok {
		logger.ApiWarn(s.LogFileName, s.RequestId, "Response json new Status ", zap.String("newCode", supplierStr))
		return nil, appError.CodeSupplierInternalChannelErrCode
	}
	if supplierError.Code != appError.SUCCESS.Code {
		return nil, supplierError
	}
	return scanData, nil
}

// DealUnifiedPayoutParams 处理请求参数 赋值一些基础值 例如：client_ip
func (s *PayoutServer) DealUnifiedPayoutParams(d *req.AspPayout, h *req.AspPaymentHeader) *appError.Error {

	if err := checker.Struct(d); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedPayoutParams:", zap.Error(err))
		return err
	}

	if err := checker.Struct(h); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedPayoutParams:", zap.Error(err))
		return err
	}

	d.PayoutParams.AppId = h.AppId
	d.PayoutParams.OrderID = d.OrderID
	d.PayoutParams.UserID = d.UserID
	d.PayoutParams.OrderCurrency = d.OrderCurrency
	d.PayoutParams.OrderAmount = d.OrderAmount
	d.PayoutParams.Timestamp = h.Timestamp
	d.PayoutParams.OrderName = d.OrderName
	d.PayoutParams.NotifyURL = d.NotifyURL
	d.PayoutParams.CustomerName = d.CustomerName
	d.PayoutParams.CustomerPhone = d.CustomerPhone
	d.PayoutParams.CustomerEmail = d.CustomerEmail
	d.PayoutParams.DeviceInfo = d.DeviceInfo
	d.PayoutParams.OrderNote = d.OrderNote
	d.PayoutParams.Ifsc = d.Ifsc
	d.PayoutParams.BankCard = d.BankCard
	d.PayoutParams.BankCode = d.BankCode
	d.PayoutParams.Vpa = d.Vpa
	d.PayoutParams.PayType = d.PayType
	d.PayoutParams.Address = d.Address
	d.PayoutParams.City = d.City
	d.PayoutParams.Sign = h.Signature

	// 接收参数 需要 e.orm
	d.PayoutParams.ClientIp = s.C.IP()

	paymentArr := strings.Split(payoutPayment, ".")

	if len(paymentArr) != 2 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedPayoutParams with invalid params payment_method:", zap.String("sign", constant.InvalidParamsErrMsg))
		return appError.CodeInvalidParamErrCode // 请求参数错误 请求被拒绝
	}

	serialNumBer := goutils.GenerateSerialNumBer("Payout", config.AppConfig.Server.Name, config.AppConfig.Server.Env)
	if serialNumBer == "" {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GenerateSerialNumBer Repeat")
		return appError.CodeUnknown // 服务器开小差了！
	}

	d.PayoutParams.Provider = strings.ToLower(paymentArr[0])  // 统一处理转小写aps
	d.PayoutParams.TradeType = strings.ToUpper(paymentArr[1]) // Payout 统一处理成转大写
	d.PayoutParams.Sn = serialNumBer                          // 生成订单号的 开头编号
	// fmt.Println(d.PayParams, "d.PayParams")
	return nil
}

// VerifyPayoutAfter 初步判断是否能够收款
func (s *PayoutServer) VerifyPayoutAfter(d *req.AspPayout, availableTotalFee int, merchantProjectConfigInfo *model.AspMerchantProjectConfig) *appError.Error {
	orderFee := d.OrderAmount + goutils.ChargeFee(d.OrderAmount, merchantProjectConfigInfo.InFeeRate) + merchantProjectConfigInfo.FixedOutAmount
	// fmt.Println("availableTotalFee+currentTotalFee ----------------", todayTotalFee)
	AuditUpperLimit := merchantProjectConfigInfo.OutAuditUpperLimit
	if d.OrderAmount > AuditUpperLimit {
		return appError.MerchantProjectCurrencyLimitExceededException // 商户金额已超出限制，有关详细信息，请参阅商户限额
	}
	// 当提现金额 大于可用金额 则进入到审核的流程
	if availableTotalFee < orderFee {
		// _ = e.AddError(apis.ConflictExceptionErr) // 请求有冲突
		d.PayoutCode = "01" // 提现金额不足
		logger.ApiWarn(s.LogFileName, s.RequestId, "availableTotalFee < orderFee  ")
		return nil
	}
	UpperLimit := merchantProjectConfigInfo.OutUpperLimit
	LowerLimit := merchantProjectConfigInfo.OutLowerLimit

	// fmt.Println("UpperLimit LowerLimit amount ---------------", UpperLimit, LowerLimit, amount)
	if d.OrderAmount > UpperLimit || d.OrderAmount < LowerLimit {
		d.PayoutCode = "02" // 单笔金额限制不通过
		// 已超出限制。有关详细信息，请参阅随附的错误消息
		logger.ApiWarn(s.LogFileName, s.RequestId, "d.OrderAmount > UpperLimit || d.OrderAmount < LowerLimit")
		return nil
	}

	outDayUpperLimit := merchantProjectConfigInfo.OutDayUpperLimit
	outDayUpperNum := merchantProjectConfigInfo.OutDayUpperNum

	// 查询商户今日的代付总额 入参 项目id 日期
	currentOutDayUpperLimit, err := s.GetMerchantProjectPayoutCurrentDayAllTotalFee(merchantProjectConfigInfo.MchProjectId)
	if err != nil {
		return err
	}
	// 如果当天的交易额度 大于或者等于配置 代付上限额度
	if currentOutDayUpperLimit >= outDayUpperLimit {
		d.PayoutCode = "03" // 商户当天总金额限制不通过
		logger.ApiWarn(s.LogFileName, s.RequestId, "currentOutDayUpperLimit >= outDayUpperLimit")
		return nil
	}
	currentOutDayUpperNum := 0
	// 查询商户今日的代付总笔数 入参 项目id 日期
	currentOutDayUpperNum, err = s.GetMerchantProjectPayoutCurrentDayAllNum(merchantProjectConfigInfo.MchProjectId)
	if err != nil {
		return err
	}
	// 如果当天的交易额度 大于或者等于配置 代付上限额度
	if currentOutDayUpperNum >= outDayUpperNum {
		d.PayoutCode = "04" // 商户当天交易数量限制不通过
		logger.ApiWarn(s.LogFileName, s.RequestId, "currentOutDayUpperNum >= outDayUpperNum")
	}
	return nil
}

func (s *PayoutServer) GetIpFilters(ipType string) ([]string, *appError.Error) {
	appId := s.Head.AppId
	tx := database.DB
	ipFiltersAddrList := make([]*req.IpFiltersAddr, 0)
	var res []string

	table := tx.Select("asp_ip_filter.ip_addr").Table("asp_ip_filter")
	table = table.Where("asp_ip_filter.type_id = ? and asp_ip_filter.type= ?", appId, ipType)
	err := table.Find(&ipFiltersAddrList).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "asp_ip_filter: ", zap.String("", constant.MissIpFilterNotFoundErrMsg))
		return nil, appError.MissIpFilterNotFoundErrMsg
	}

	if len(ipFiltersAddrList) == 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(ipFiltersAddrList) == 0: ", zap.String("", constant.MissIpFilterNotFoundErrMsg))
		return nil, appError.MissIpFilterNotFoundErrMsg
	}

	for i := 0; i < len(ipFiltersAddrList); i++ {
		res = append(res, ipFiltersAddrList[i].IpAddr)
	}

	return res, nil
}

// VerifyPayoutBefore 初步判断是否能够收款
func (s *PayoutServer) VerifyPayoutBefore(aspMerchantProject *model.AspMerchantProject, d *req.AspPayout, clientIp string) *appError.Error {

	// 判断系统是否开启代付
	systemPayoutStatusInfo, err := s.getSystemPayoutStatus(constant.SYSTEM_PAYOUT_STATUS)
	if err != nil {
		return err
	}
	if systemPayoutStatusInfo.Data == constant.SYSTEM_PAYOUT_STATUS_OFF {
		return appError.NewError(constant.PayoutOffErrMsg)
	}
	// 判断商户自定义收款单是否存在
	payoutInfoParam := make(map[string]interface{})
	payoutInfoParam["mch_id"] = aspMerchantProject.MchId
	payoutInfoParam["out_trade_no"] = d.OrderID
	_, err = s.GetPayoutInfo(payoutInfoParam, nil)

	// 判断是否存在 已经存在则不会有 err 则需要处理
	if err == nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "logic.GetPayoutInfo: ", zap.String("", constant.ConflictExceptionErrMsg))
		return appError.CodeConflictException // 请求有冲突
	}
	// 判断ip 白名单设置
	if ipErr := s.VerifyIpFilter("1", s.Head.AppId, clientIp); ipErr != nil {
		return ipErr
	}
	return nil
}

// VerifyIpFilter 验证IP白名单
func (s *PayoutServer) VerifyIpFilter(typeCode string, typeId string, ip string) *appError.Error {

	if ip != "" {
		exits := goRedis.Redis.SIsMember(context.Background(), s.getIpFilterKey(typeCode, typeId), ip).Val()
		if !exits {
			logger.ApiWarn(s.LogFileName, s.RequestId, "utils.InArray.clientIP: ", zap.String("ip", constant.MissIpFilterNotFoundErrMsg))
			missIpFilterNotFoundErrMsg := *appError.MissIpFilterNotFoundErrMsg
			return (&missIpFilterNotFoundErrMsg).FormatMessage(ip) //
		}
	}

	return nil
}

// GetPayoutBenefiary 判断受益人信息是否合法
func (s *PayoutServer) GetPayoutBenefiary(d *req.AspPayout) *appError.Error {

	params := make(map[string]interface{})
	params["customer_name"] = d.CustomerName
	params["customer_phone"] = d.CustomerPhone
	params["ifsc"] = d.Ifsc
	params["bank_card"] = d.BankCard
	// 查询受益人信息
	BenefiaryInfo, beErr := s.BeneficiaryInfo(params)
	if beErr == nil {
		// 存在 则返回
		if BenefiaryInfo.Id > 0 {
			d.PayoutParams.BenefiaryId = BenefiaryInfo.BenefiaryId
		}
	}

	return nil
}

// Insert Get 获取SysApi对象with id
func (s *PayoutServer) Insert(d *req.AspPayout, aspMerchantProject *model.AspMerchantProject, aspMerchantProjectConfig *model.AspMerchantProjectConfig, merchantProjectCurrencyInfo *model.AspMerchantProjectCurrency, aspChannelDepartTradeType *req.DeptTradeTypeInfo) (*model.AspPayout, *appError.Error) {
	var data model.AspPayout
	// 赋值给 order 数据
	d.Generate(&data, aspMerchantProject, aspMerchantProjectConfig, merchantProjectCurrencyInfo, aspChannelDepartTradeType)
	err := database.DB.Create(&data).Error
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspPayout Insert: ", zap.Error(err))
		return nil, appError.NewError(err.Error())
	}
	return &data, nil
}

// InsertBeneficiar 写入到受益人表
func (s *PayoutServer) InsertBeneficiar(payoutInfo *model.AspPayout, merchantProjectInfo *model.AspMerchantProject, channelDepartInfo *model.AspChannelDepartConfig, currencyId int) (*model.AspBeneficiary, *appError.Error) {
	var data model.AspBeneficiary
	// 赋值给 order 数据
	serialNumBer := goutils.GenerateSerialNumBer("Benefiary", config.AppConfig.Server.Name, config.AppConfig.Server.Env)
	if serialNumBer == "" {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GenerateSerialNumBer Repeat")
		return nil, appError.CodeUnknown // 服务器开小差了！
	}
	data.MchProjectBenefiary = serialNumBer // 生成订单号的 开头编号
	data.MchId = merchantProjectInfo.MchId
	data.DepartId = cast.ToInt(channelDepartInfo.DepartId)
	data.ChannelId = cast.ToUint(channelDepartInfo.ChannelId)
	data.CurrencyId = currencyId
	data.MchProjectId = merchantProjectInfo.Id
	data.SpbillCreateIp = s.C.IP()
	data.NotifyUrl = ""
	data.CustomerId = payoutInfo.CustomerName
	data.CustomerName = payoutInfo.CustomerName
	data.CustomerEmail = payoutInfo.CustomerEmail
	data.CustomerPhone = payoutInfo.CustomerPhone
	data.Ifsc = payoutInfo.Ifsc
	data.BankCard = payoutInfo.BankCard
	data.BankCode = payoutInfo.BankCode
	data.TradeState = constant.BENEFICIARY_TRADE_STATE_SUCCESS
	data.Provider = payoutInfo.Provider
	data.Adapter = payoutInfo.Adapter
	data.TradeType = payoutInfo.TradeType
	data.BenefiaryId = payoutInfo.BeneficiaryId
	data.CreateTime = cast.ToUint64(goutils.GetDateTimeUnix())
	data.FinishTime = data.CreateTime
	data.SupplierReturnCode = ""
	data.SupplierReturnMsg = ""
	err := database.DB.Create(&data).Error
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspBeneficiary Insert: ", zap.Error(err))
		return nil, appError.NewError(err.Error())
	}
	return &data, nil
}

// DealUnifiedPayoutAuditParams 处理请求参数 验证
func (s *PayoutServer) DealUnifiedPayoutAuditParams(d *req.AspPayoutAudit) error {

	if err := checker.Struct(d); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedPayoutParams ", zap.Error(err))
		return err
	}
	return nil
}

// DealUnifiedPayoutAuditJobParams 处理审核脚本参数 验证
func (s *PayoutServer) DealUnifiedPayoutAuditJobParams(d *reqCommon.AspPayoutAuditJob) error {

	if err := checker.Struct(d); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedPayoutAuditJobParams ", zap.Error(err))
		return err
	}
	return nil
}

// DealUnifiedPayoutQueryParams 处理请求参数 参数的验证
func (s *PayoutServer) DealUnifiedPayoutQueryParams(d *req.AspPayoutQuery, h *req.AspPaymentHeader) *appError.Error {

	if err := checker.Struct(d); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedPayoutQueryParams ", zap.Error(err))
		return err
	}

	if err := checker.Struct(h); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedPayoutQueryParams ", zap.Error(err))
		return err
	}

	return nil
}

// CheckSignPayout 处理支付成功，加款等各项操作 更新订单上游返回的数据
func (s *PayoutServer) CheckSignPayout(d *req.AspPayout, paySecret string) *appError.Error {

	// TODO 需要完善对应的字段 验证签名 如果没有传参数 这里需要处理
	params := make(map[string]interface{})
	params["order_id"] = strings.TrimSpace(d.OrderID)
	params["order_currency"] = strings.TrimSpace(d.OrderCurrency)
	params["order_amount"] = cast.ToString(d.OrderAmount)
	params["Timestamp"] = cast.ToString(d.Timestamp)
	params["order_name"] = strings.TrimSpace(d.OrderName)
	params["user_id"] = strings.TrimSpace(d.UserID)
	params["notify_url"] = strings.TrimSpace(d.NotifyURL)
	params["customer_name"] = strings.TrimSpace(d.CustomerName)
	params["customer_phone"] = strings.TrimSpace(d.CustomerPhone)
	params["customer_email"] = strings.TrimSpace(d.CustomerEmail)
	params["device_info"] = strings.TrimSpace(d.DeviceInfo)
	params["order_note"] = strings.TrimSpace(d.OrderNote)
	params["ifsc"] = strings.TrimSpace(d.Ifsc)
	params["bank_card"] = strings.TrimSpace(d.BankCard)
	params["bank_code"] = strings.TrimSpace(d.BankCode)
	params["vpa"] = strings.TrimSpace(d.Vpa)
	params["pay_type"] = strings.TrimSpace(d.PayType)
	params["address"] = strings.TrimSpace(d.Address)
	params["city"] = strings.TrimSpace(d.City)
	params["sign"] = strings.TrimSpace(d.Sign)

	if !goutils.HmacSHA256Verify(params, paySecret) {
		return appError.UnauthenticatedErrCode // 签名错误
	}
	return nil
}

func (s *PayoutServer) IsNeedPayoutQueryUpstream(payoutInfo *model.AspPayout, params map[string]string) bool {
	// 这个参数是强制查询上游
	if _, ok := params["is_call_upstream"]; ok && params["is_call_upstream"] == "yes" {
		return true
	}

	// 判断结构体是否为空
	if payoutInfo == nil || payoutInfo.Id < 1 {
		return false
	}

	// 判断状态是否是
	status := []string{constant.PAYOUT_TRADE_STATE_SUCCESS, constant.PAYOUT_TRADE_STATE_FAILED, constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS, constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED}

	// strings.ToUpper 将字符串转成大写
	if goutils.InArray(strings.ToUpper(payoutInfo.TradeState), status) {
		return false
	}

	return true
}

// IsNeedPayoutAuditOperate 判断审核操作 是否合法
func (s *PayoutServer) IsNeedPayoutAuditOperate(payoutInfo *model.AspPayout, payoutAudit *req.AspPayoutAudit) bool {
	// 判断结构体是否为空
	if payoutInfo == (&model.AspPayout{}) {
		return false
	}

	// 判断状态是否是
	var status []string
	status = []string{constant.PAYOUT_TRADE_STATE_APPLY, constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS, constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING, constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED, constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS}

	// strings.ToUpper 将字符串转成大写
	if !goutils.InArray(strings.ToUpper(payoutInfo.TradeState), status) {
		return false
	}

	if payoutAudit.Action == "pass" {
		var passStatus []string
		passStatus = []string{constant.PAYOUT_TRADE_STATE_APPLY, constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS, constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING, constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED, constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS}

		// strings.ToUpper 将字符串转成大写
		if !goutils.InArray(strings.ToUpper(payoutInfo.TradeState), passStatus) {
			return false
		}
	}

	if payoutAudit.Action == "return" {
		var returnStatus []string
		returnStatus = []string{constant.PAYOUT_TRADE_STATE_APPLY, constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS}
		// strings.ToUpper 将字符串转成大写
		if !goutils.InArray(strings.ToUpper(payoutInfo.TradeState), returnStatus) {
			return false
		}
	}
	return true
}

// CheckSignPayoutQuery 提现查询 处理验签
func (s *PayoutServer) CheckSignPayoutQuery(paySecret string, reqPayoutQuery *req.AspPayoutQuery, timestamp int) *appError.Error {
	// TODO 需要完善对应的字段 验证签名
	params := make(map[string]interface{})
	params["sn"] = strings.TrimSpace(reqPayoutQuery.Sn)
	params["is_call_upstream"] = strings.TrimSpace(reqPayoutQuery.IsCallUpstream)
	params["Timestamp"] = timestamp
	params["sign"] = strings.TrimSpace(s.Head.Signature)
	// fmt.Println("req.AspPayoutQuery-----------", d)
	// fmt.Println("req.AspPayoutQueryHeaderReq-----------", h)
	if !goutils.HmacSHA256Verify(params, paySecret) {
		return appError.NewError("签名错误")
	}
	return nil
}

func (s *PayoutServer) GetPayoutInfo(params map[string]interface{}, tx *gorm.DB) (*model.AspPayout, *appError.Error) {
	if tx == nil {
		tx = database.DB
	}
	var payoutInfo *model.AspPayout

	for k, v := range params {
		switch k {
		case "id":
			tx = tx.Where("id = ?", v)
		case "sn":
			tx = tx.Where("sn = ?", v)
		case "mch_id":
			tx = tx.Where("mch_id = ?", v)
		case "out_trade_no":
			tx = tx.Where("out_trade_no = ?", v)
		case "transaction_id":
			tx = tx.Where("transaction_id = ?", v)
		case "provider":
			tx = tx.Where("provider = ?", v)
		case "adapter":
			tx = tx.Where("adapter = ?", v)
		}
	}

	if err := tx.First(&payoutInfo).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return payoutInfo, appError.PayoutNotFoundErrCode
	}

	return payoutInfo, nil

}

// SolvePayoutSuccess 处理提现成功的加款等各项操作 更新提现上游返回的数据
// 所有相关提现的更新 都使用这个一个方法来处理
// OrderId 提现id 主键
// params 修改的数组
// updateStructType 修改提现类型
// host 数据库的链接 框架需要
// request_id 系统中产生的唯一请求的id
// 新增审核状态 加入乐观锁 2023-02-16
func (s *PayoutServer) SolvePayoutSuccess(payoutId int, params map[string]interface{}, updateStructType string, tx *gorm.DB) (*model.AspPayout, *appError.Error) {
	o := database.DB
	if tx != nil {
		o = tx
	}

	payoutInfoParam := make(map[string]interface{})
	payoutInfoParam["id"] = payoutId

	payoutInfo, err := s.GetPayoutInfo(payoutInfoParam, tx)

	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspPayout: ", zap.Error(err))
		return nil, err
	}

	b, _ := json.Marshal(payoutInfo)
	logger.ApiInfo(s.LogFileName, s.RequestId, "update payout before: ", zap.String("updateStructType", updateStructType), zap.Any("params", params), zap.String("before", string(b)))

	aspPayout := model.AspPayout{}

	resultPayoutUpdate := o.Model(&aspPayout).Where("id = ?", payoutId)

	// 参数 状态 审核状态
	_, tradeStateOk := params["trade_state"]
	_, isCheckOutOk := params["is_checkout"]

	if tradeStateOk {
		switch params["trade_state"] {
		case constant.PAYOUT_TRADE_STATE_RETURN:
			resultPayoutUpdate.Where("trade_state = ?", constant.PAYOUT_TRADE_STATE_APPLY)
			break
		case constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED:
			resultPayoutUpdate.Where("trade_state = ? or trade_state = ?", constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING, constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS)
			break
		// 冻结
		case constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS:
			resultPayoutUpdate.Where("trade_state = ?", constant.PAYOUT_TRADE_STATE_APPLY)
			break
		case constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING:
			resultPayoutUpdate.Where("trade_state = ?", constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS)
			break
		case constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS:
			resultPayoutUpdate.Where("trade_state = ? or trade_state = ?", constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING, constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS)
			break
		//解冻代付取消
		case constant.PAYOUT_TRADE_STATE_REVOKE:
			resultPayoutUpdate.Where("trade_state = ?", constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS)
			break
		// 解冻代付成功
		case constant.PAYOUT_TRADE_STATE_SUCCESS:
			resultPayoutUpdate.Where("trade_state = ? or trade_state = ?", constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING, constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS)
			break
		//解冻代付失败
		case constant.PAYOUT_TRADE_STATE_FAILED:
			resultPayoutUpdate.Where("trade_state = ? or trade_state = ? or trade_state = ?", constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING, constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS, constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED)
			break
		default:
			resultPayoutUpdate.Where("trade_state = ?", constant.PAYOUT_TRADE_STATE_PENDING)
			break
		}
	}
	// 当存在脚本批量审核，修改审核状态时候。
	if isCheckOutOk {
		resultPayoutUpdate.Where("is_checkout = ? or is_checkout = ?", constant.PAYOUT_IS_CHECKOUT_APPLY, constant.PAYOUT_IS_CHECKOUT_PENDING)
	}

	resultPayoutUpdate.Updates(params)

	if resultPayoutUpdate.RowsAffected < 1 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "Update Order error: ", zap.String("err: ", constant.UpdatePayoutRowsErrMsg))
		return nil, appError.NewError(constant.ChangeErrMsg).FormatMessage("Payout RowsAffected") // 更新代收订单行数 请重新提交
	}

	if resultPayoutUpdate.Error != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "payout.update() ", zap.Error(resultPayoutUpdate.Error))
		return nil, appError.NewError(resultPayoutUpdate.Error.Error())
	}

	afterPayoutInfo, _ := s.GetPayoutInfo(payoutInfoParam, tx)
	b, _ = json.Marshal(afterPayoutInfo)
	logger.ApiInfo(s.LogFileName, s.RequestId, "update payout after: ", zap.String("updateStructType", updateStructType), zap.String("after", string(b)))
	return afterPayoutInfo, nil
}

func (s *PayoutServer) InsertAuditRecord(d *req.AspPayoutAudit, payoutInfo *model.AspPayout, tx *gorm.DB) *appError.Error {
	o := database.DB
	if tx != nil {
		o = tx
	}
	// 初始化 审核记录
	var data model.AspAuditRecords

	s.GenerateAuditRecord(&data, d, payoutInfo)
	err := o.Create(&data).Error
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspMerchantProjectUser.Create(&data) ", zap.Any("AuditRecord", data), zap.Error(err))
		return appError.NewError(err.Error())
	}
	return nil
}

func (s *PayoutServer) GenerateAuditRecord(data *model.AspAuditRecords, d *req.AspPayoutAudit, payoutInfo *model.AspPayout) {
	data.AuditType = constant.AUDIT_RECORDS_AUDIT_TYPE_PAYOUT
	data.ProjectId = payoutInfo.MchProjectId
	data.Sn = payoutInfo.Sn
	data.DepartId = payoutInfo.DepartId
	data.ChannelId = payoutInfo.ChannelId
	data.MchId = payoutInfo.MchId
	data.Status = s.GetAuditRecordStatus(d.Action)
	data.OperateId = d.OperationID
	data.OperateTime = goutils.GetDateTimeUnix()
}

// GetAuditRecordStatus 审核输入状态 转换
func (s *PayoutServer) GetAuditRecordStatus(actionType string) int {
	var auditRecordStatus = make(map[string]int)
	auditRecordStatus["pass"] = constant.AUDIT_RECORDS_AUDIT_TYPE_STATUS_PASS
	auditRecordStatus["return"] = constant.AUDIT_RECORDS_AUDIT_TYPE_STATUS_RETURN

	status := 0
	if firstpayStatus, ok := auditRecordStatus[actionType]; ok {
		status = firstpayStatus
	} else {
		status = constant.AUDIT_RECORDS_AUDIT_TYPE_STATUS_DOING
	}

	return status
}

// FreezeAndUpdate 冻结成功 更改可用余额 修改代付状态
func (s *PayoutServer) FreezeAndUpdate(merchantProjectRepository *repository.MerchantProjectRepository, payoutInfo *model.AspPayout, changeAmount int, tx *gorm.DB) (*model.AspPayout, *appError.Error) {
	// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
	err := merchantProjectRepository.PayoutOrderFreezeFee(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
	if err != nil {
		return nil, appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent FREEZE_SUCCESS") // 更新商户余额错误 请重新提交
	}
	payoutFreezeUpdate := make(map[string]interface{})
	payoutFreezeUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS
	_, errFreeze := s.SolvePayoutSuccess(payoutInfo.Id, payoutFreezeUpdate, "FreezeAndUpdate", tx)
	if errFreeze != nil {
		return nil, appError.NewError(constant.ChangeErrMsg).FormatMessage("FreezeAndUpdate") // 更新商户余额错误 请重新提交
	}
	payoutInfoParam := make(map[string]interface{})
	payoutInfoParam["id"] = payoutInfo.Id
	newPayoutInfo, _ := s.GetPayoutInfo(payoutInfoParam, nil)
	return newPayoutInfo, nil
}

// RequestUpstream 请求上游
func (s *PayoutServer) RequestUpstream(payoutInfo *model.AspPayout) (*interfaces.ThirdPayoutCreateData, *appError.Error) {
	// 获取到上游数据 获取到对应的上游
	supplierCode := payoutInfo.Adapter + "." + payoutInfo.TradeType
	//fmt.Println("supplierCode", supplierCode)
	// supplierCode := "firstpay.H5"
	supplierServer := thirdParty.GetPaySupplierByCode(supplierCode)
	if supplierServer == nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "thirdParty.GetPaySupplierByCode error ", zap.String("err", "supplierServer == nil"), zap.Any("supplierCode", supplierCode), zap.Any("supplierServer", supplierServer))
		return nil, appError.ChannelDepartTradeTypeNotFoundErrCode // 渠道信息不存在
	}
	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(payoutInfo.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(payoutInfo.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, err := NewSimpleChannelConfigServer().GetChannelDepartInfo(channelDepartInfoParam)
	if err != nil {
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartNotFoundErrMsg) // missing channel depart not found
	}
	scanData := new(interfaces.ThirdPayoutCreateData)
	pErr := new(appError.Error)
	// haodapay 新增upi代付支付方式
	if payoutInfo.PayType == constant.PARAMS_PAY_TYPE_BANK {
		scanData, pErr = supplierServer.Payout(s.RequestId, channelDepartInfo, payoutInfo)
	} else if payoutInfo.PayType == constant.PARAMS_PAY_TYPE_UPI {
		scanData, pErr = supplierServer.PayoutUpi(s.RequestId, channelDepartInfo, payoutInfo)
	}

	if pErr != nil {
		return scanData, pErr
	}
	return scanData, nil
}

// PayoutCreateUpdateFailed 创建代收订单上游返回的错误 更新订单+释放可用金额
func (s *PayoutServer) PayoutCreateUpdateFailed(merchantProjectRepository *repository.MerchantProjectRepository, scanData *interfaces.ThirdPayoutCreateData, payoutInfo *model.AspPayout, changeAmount int, pErr *appError.Error, tx *gorm.DB) *appError.Error {
	//上游返回异常信息,并记录上游异常信息
	if pErr.Code == appError.CodeSupplierChannelErrCode.Code {
		//记录日志
		logger.ApiWarn(s.LogFileName, s.RequestId, "request payout up supplier rsp err", zap.Any("code", pErr.Code), zap.String("msg", pErr.Message))
		// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
		payoutFailedUpdate := make(map[string]interface{})
		payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FAILED
		payoutFailedUpdate["supplier_return_code"] = scanData.Code
		payoutFailedUpdate["supplier_return_msg"] = scanData.Msg

		_, errFreeze := s.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "PayoutCreateUpdateFailed", tx)
		if errFreeze != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("PayoutCreateUpdateFailed") // 更新商户余额错误 请重新提交
		}
		// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
		err := merchantProjectRepository.PayoutOrderChannelFailed(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
		if err != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
		}
	}
	return nil
}

func (s *PayoutServer) PayoutCreateUpdate(merchantProjectRepository *repository.MerchantProjectRepository, payoutInfo *model.AspPayout, changeAmount int, scanData *interfaces.ThirdPayoutCreateData, tx *gorm.DB) *appError.Error {
	params := make(map[string]interface{})
	params["transaction_id"] = goutils.IfString(scanData.TransactionID != "", scanData.TransactionID, payoutInfo.TransactionId)
	params["cash_fee"] = scanData.CashFee
	params["cash_fee_type"] = payoutInfo.FeeType
	params["finish_time"] = "0"
	params["trade_state"] = scanData.TradeState // 赋值代付订单状态
	params["bank_utr"] = scanData.BankUtr       // Bank UTR No

	// 请求上游后操作
	upErr := s.RequestUpstreamAndUpdate(merchantProjectRepository, params, payoutInfo, changeAmount, tx)
	if upErr != nil {
		return upErr
	}
	return nil
}

// RequestUpstreamAndUpdate 请求上游 执行操作
func (s *PayoutServer) RequestUpstreamAndUpdate(merchantProjectRepository *repository.MerchantProjectRepository, params map[string]interface{}, payoutInfo *model.AspPayout, changeAmount int, tx *gorm.DB) *appError.Error {
	//var isSendQueue bool
	//isSendQueue = false
	//if params["trade_state"] == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
	//	isSendQueue = true
	//}
	// 修改代付订单状态 包含（代付订单 代付中情况）
	_, err := s.SolvePayoutSuccess(payoutInfo.Id, params, "RequestUpstreamAndUpdate", tx)
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "RequestUpstreamAndUpdate", zap.Any("params", params), zap.Int("Id", payoutInfo.Id))
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("PayoutOrder") // 更新代付订单错误 请重新提交
	}

	// 当成功时候
	if params["trade_state"] == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
		// 代付上游成功 解冻+代付成功  冻结金额减去  预扣金额记录新增  收支流水扣减
		payoutSuccessUpdate := make(map[string]interface{})
		payoutSuccessUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_SUCCESS
		payoutSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()

		_, errSuccess := s.SolvePayoutSuccess(payoutInfo.Id, payoutSuccessUpdate, "RequestUpstreamAndUpdate", tx)
		if errSuccess != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("RequestUpstreamAndUpdate") // 更新商户余额错误 请重新提交
		}
		// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
		err = merchantProjectRepository.PayoutOrderChannelSuccess(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
		if err != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_SUCCESS") // 更新商户余额错误 请重新提交
		}
	}

	// 当失败的时候
	if params["trade_state"] == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
		// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
		payoutFailedUpdate := make(map[string]interface{})
		payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FAILED

		_, errFailed := s.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "RequestUpstreamAndUpdate", tx)
		if errFailed != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("RequestUpstreamAndUpdate") // 更新商户余额错误 请重新提交
		}
		// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
		err = merchantProjectRepository.PayoutOrderChannelFailed(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
		if err != nil {
			return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
		}
	}

	return nil
}

// RequestUpstreamQueryPayout 请求上游执行订单查询 + 更新响应结果
func (s *PayoutServer) RequestUpstreamQueryPayout(payoutInfo *model.AspPayout, channelDepartInfo *model.AspChannelDepartConfig) (*interfaces.ThirdPayoutQueryData, *appError.Error) {
	// 获取到上游数据 获取到对应的上游
	supplierCode := payoutInfo.Adapter + "." + payoutInfo.TradeType
	// supplierCode := "firstpay.H5"
	supplier := thirdParty.GetPaySupplierByCode(supplierCode)
	if supplier == nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "thirdParty.GetPaySupplierByCode error ", zap.String("err", "supplier == nil"), zap.Any("supplierCode", supplierCode), zap.Any("supplier", supplier))
		return nil, appError.ChannelDepartTradeTypeNotFoundErrCode // 渠道信息不存在
	}

	scanQueryData, err := supplier.PayoutQuery(s.RequestId, channelDepartInfo, payoutInfo)
	if err != nil {
		return scanQueryData, err
	}
	// 渠道映射系统的统一返回字符串 在 appError中有定义的 map
	supplierStr := payoutInfo.Adapter + "_" + scanQueryData.Code
	supplierError, ok := appError.SupplierErrorMap[supplierStr]
	if !ok {
		logger.ApiWarn(s.LogFileName, s.RequestId, "Response json new Status ", zap.String("newCode", supplierStr))
		return scanQueryData, appError.CodeSupplierInternalChannelErrCode
	}
	if supplierError.Code != appError.SUCCESS.Code {
		return scanQueryData, supplierError
	}

	fmt.Println("scanQueryData:", fmt.Sprintf("%+v", scanQueryData))

	return scanQueryData, nil
}

func (s *PayoutServer) PayoutQueryUpdateFailed(merchantProjectRepository *repository.MerchantProjectRepository, scanQueryData *interfaces.ThirdPayoutQueryData, payoutInfo *model.AspPayout, errUpQuery *appError.Error, changeAmount int, tx *gorm.DB) *appError.Error {
	//上游返回异常信息,并记录上游异常信息
	if errUpQuery.Code == appError.CodeSupplierInternalChannelParamsFailedErrCode.Code {
		if payoutInfo.TradeState != constant.PAYOUT_TRADE_STATE_FAILED {
			//记录日志
			logger.ApiWarn(s.LogFileName, s.RequestId, "request payout up supplier rsp err", zap.Any("code", errUpQuery.Code), zap.String("msg", errUpQuery.Message))

			// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
			payoutFailedUpdate := make(map[string]interface{})
			payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FAILED
			payoutFailedUpdate["supplier_return_code"] = scanQueryData.Code
			payoutFailedUpdate["supplier_return_msg"] = scanQueryData.Msg

			_, errFailed := s.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "PayoutQueryUpdateFailed", tx)
			if errFailed != nil {
				return appError.NewError(constant.ChangeErrMsg).FormatMessage("PayoutQueryUpdateFailed") // 更新商户余额错误 请重新提交
			}
			// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
			cErr := merchantProjectRepository.PayoutOrderChannelFailed(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
			if cErr != nil {
				return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
			}
		}
	}
	return nil
}

func (s *PayoutServer) PayoutQueryUpdate(payoutInfo *model.AspPayout, scanQueryData *interfaces.ThirdPayoutQueryData, changeAmount int, tx *gorm.DB) (*model.AspPayout, *appError.Error) {
	upstreamStatus := scanQueryData.TradeState
	//isCapitalFlow := false
	newPayoutInfo := &model.AspPayout{}
	var err *appError.Error
	payoutUpdateParams := make(map[string]interface{})

	// 当系统想查询代付状态的 请求到上游后
	// 判断状态是否是 最终成功或者最终失败
	IsNeedStatus := []string{constant.PAYOUT_TRADE_STATE_SUCCESS, constant.PAYOUT_TRADE_STATE_FAILED}
	if payoutInfo.TradeState != upstreamStatus && !goutils.InArray(strings.ToUpper(payoutInfo.TradeState), IsNeedStatus) {

		payoutUpdateParams["transaction_id"] = scanQueryData.TransactionID
		payoutUpdateParams["cash_fee"] = scanQueryData.CashFee
		payoutUpdateParams["cash_fee_type"] = payoutInfo.FeeType
		payoutUpdateParams["bank_utr"] = payoutInfo.BankUtr // 默认值
		payoutUpdateParams["trade_state"] = upstreamStatus  // 默认值
		// 如果成功 必须提现状态是 顺序的增长的，不能状态逆序
		// 如果上游返回成功 提现状态是 已申请 则修改
		newPayoutInfo, err = s.SolvePayoutSuccess(payoutInfo.Id, payoutUpdateParams, "PayPayoutQueryUpdate", tx)
		if err != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "PayoutQueryUpdate", zap.Any("params", payoutUpdateParams), zap.Int("Id", payoutInfo.Id))
			return nil, appError.NewError(constant.ChangeErrMsg).FormatMessage("PayPayoutQueryUpdate") // 更新代付订单错误 请重新提交
		}

		// 当成功时候
		if scanQueryData.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS {
			MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)

			// 代付上游成功 解冻+代付成功  冻结金额减去  预扣金额记录新增  收支流水扣减
			payoutSuccessUpdate := make(map[string]interface{})
			payoutSuccessUpdate["bank_utr"] = scanQueryData.BankUtr
			payoutSuccessUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_SUCCESS
			payoutSuccessUpdate["finish_time"] = goutils.GetDateTimeUnix()
			_, errSuccess := s.SolvePayoutSuccess(payoutInfo.Id, payoutSuccessUpdate, "PayoutQueryUpdate", tx)
			if errSuccess != nil {
				return nil, appError.NewError(constant.ChangeErrMsg).FormatMessage("PayoutQueryUpdate") // 更新商户余额错误 请重新提交
			}
			// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
			err = MerchantProjectRepository.PayoutOrderChannelSuccess(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
			if err != nil {
				return nil, appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_SUCCESS") // 更新商户余额错误 请重新提交
			}
			payoutInfoParam := make(map[string]interface{})
			payoutInfoParam["id"] = payoutInfo.Id
			newPayoutInfo, err = s.GetPayoutInfo(payoutInfoParam, tx)
		}

		// 当失败的时候
		if scanQueryData.TradeState == constant.PAYOUT_TRADE_STATE_CHANNEL_FAILED {
			MerchantProjectRepository := repository.NewMerchantProjectRepository(s.LogFileName, s.RequestId)

			// 代付上游失败解冻  可用余额增加   冻结余额释放   预扣金额记录新增
			payoutFailedUpdate := make(map[string]interface{})
			payoutFailedUpdate["trade_state"] = constant.PAYOUT_TRADE_STATE_FAILED
			_, errFailed := s.SolvePayoutSuccess(payoutInfo.Id, payoutFailedUpdate, "PayoutQueryUpdateFailed", tx)
			if errFailed != nil {
				return nil, appError.NewError(constant.ChangeErrMsg).FormatMessage("PayoutQueryUpdateFailed") // 更新商户余额错误 请重新提交
			}
			// cp项目id  操作的可用余额代付  待结算余额  冻结金额  业务类型 1:代收 2:代付  具体业务类型来源表关联id  备注  修改对应的参数
			err = MerchantProjectRepository.PayoutOrderChannelFailed(payoutInfo.MchProjectId, changeAmount, payoutInfo.Id, tx)
			if err != nil {
				return nil, appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent UNFREEZE_FAILED") // 更新商户余额错误 请重新提交
			}
			payoutInfoParam := make(map[string]interface{})
			payoutInfoParam["id"] = payoutInfo.Id
			newPayoutInfo, err = s.GetPayoutInfo(payoutInfoParam, tx)
		}
		//if isCapitalFlow == true {
		//	_ = NewSendQueueServer().SendNotifyQueue(payoutInfo.Sn)
		//}
	} else {
		payoutInfo.TradeState = upstreamStatus
		newPayoutInfo = payoutInfo
	}
	return newPayoutInfo, nil
}

func (s *PayoutServer) GetAspPayoutList(timeBegin, timeEnd int64) ([]*model.AspPayout, *appError.Error) {
	o := database.DB
	var payoutList []*model.AspPayout
	var payoutModel model.AspPayout

	err := o.Model(&payoutModel).Where("create_time >= ?", timeBegin).
		Where("create_time <= ?", timeEnd).
		Where(map[string]interface{}{"trade_state": []string{constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING}}).
		Order("id desc").
		//Limit(10).
		//Offset(0).
		Find(&payoutList).Error

	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "payoutList", zap.Error(err))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissPayoutListNotFoundErrMsg)
	}
	if len(payoutList) == 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(payoutList) == 0 ", zap.String("error", constant.MissChannelDepartTradeTypeNotFoundErrMsg))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissPayoutListNotFoundErrMsg)
	}
	return payoutList, nil
}

// getBankCategoryInfo 根据银行编码查询系统配置中的对应的编码
func (s *PayoutServer) getBankCategoryInfo(code string) *model.AspBankCategory {
	bankCategory, _ := repository.NewRepository[*model.AspBankCategory](database.DB, goRedis.Redis).FindOne(database.NewSqlCondition().Where("code = ?", code).Where("status = ?", 1))
	return bankCategory
}

// getSystemPayoutStatus 获取系统是否开启代付配置
func (s *PayoutServer) getSystemPayoutStatus(code string) (*model.AspConfig, *appError.Error) {
	var systemPayoutStatus model.AspConfig

	redisKey := constant.GetNoExpiredRedisKey(constant.KEY_SYSTEM_CONFIG_PAYOUT_STATUS)
	systemPayoutStatusValue, err := s.Redis().Get(context.Background(), redisKey).Result()
	if len(systemPayoutStatusValue) > 0 {
		err = goutils.JsonDecode(systemPayoutStatusValue, &systemPayoutStatus)
		if err != nil {
			logger.ApiError(s.LogFileName, s.RequestId, "getSystemPayoutStatus err: ", zap.Error(err))
			return nil, appError.NewError(err.Error())
		}
		return &systemPayoutStatus, nil
	}
	systemPayoutStatusInfo, _ := repository.NewRepository[*model.AspConfig](database.DB, goRedis.Redis).FindOne(database.NewSqlCondition().Where("name = ?", code))

	if systemPayoutStatusInfo == nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "systemPayoutStatusInfo == nil")
		return nil, appError.NewError(constant.MissSystemPayoutStatusNotFoundErrMsg)
	}

	s.SetSystemPayoutStatusToCache(systemPayoutStatusInfo, redisKey)

	return systemPayoutStatusInfo, nil
}

func (s *PayoutServer) SetSystemPayoutStatusToCache(systemPayoutStatusInfo *model.AspConfig, redisKey string) bool {
	// 转成JSON
	departJsonByte, err := json.Marshal(systemPayoutStatusInfo)
	if err != nil {
		return true
	}
	err = s.Redis().Set(context.Background(), redisKey, string(departJsonByte), constant.NO_EXPIRED_TIME).Err()
	return err == nil
}

// GetMerchantProjectPayoutCurrentDayAllTotalFee 获取商户代付当天总支付金额
func (s *PayoutServer) GetMerchantProjectPayoutCurrentDayAllTotalFee(mchProjectId string) (int, *appError.Error) {
	// 一天开始的时间戳
	timeBegin := carbon.Parse(carbon.Now().Format("Y-m-d")).StartOfDay().Timestamp()
	// 计算当天支付金额状态
	tradeState := []string{constant.PAYOUT_TRADE_STATE_APPLY, constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS, constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING, constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS, constant.PAYOUT_TRADE_STATE_SUCCESS}

	var payoutAmount *model.AspPayout
	payoutAmountList := make([]*req.PayoutAmountList, 0)
	// 统计当天提现的成功额度
	err := database.DB.Model(payoutAmount).Select("sum(total_fee) as amount").
		Where("create_time >= ?", timeBegin).
		Where("mch_project_id = ?", mchProjectId).
		Where(map[string]interface{}{"trade_state": tradeState}).
		Find(&payoutAmountList).Error
	if err != nil {
		return 0, appError.NewError(err.Error())
	}
	return payoutAmountList[0].Amount, nil
}

// GetMerchantProjectPayoutCurrentDayAllNum 获取商户当天的总支付笔数
func (s *PayoutServer) GetMerchantProjectPayoutCurrentDayAllNum(mchProjectId string) (int, *appError.Error) {
	// 一天开始的时间戳
	timeBegin := carbon.Parse(carbon.Now().Format("Y-m-d")).StartOfDay().Timestamp()
	// 计算当天支付金额状态
	tradeState := []string{constant.PAYOUT_TRADE_STATE_APPLY, constant.PAYOUT_TRADE_STATE_FREEZE_SUCCESS, constant.PAYOUT_TRADE_STATE_CHANNEL_PENDING, constant.PAYOUT_TRADE_STATE_CHANNEL_SUCCESS, constant.PAYOUT_TRADE_STATE_SUCCESS}

	var payoutCount *model.AspPayout
	payoutCountList := make([]*req.PayoutCountList, 0)
	// 统计当天提现的成功额度
	if err := database.DB.Model(payoutCount).Select("count(*) as count").
		Where("create_time >= ?", timeBegin).
		Where("mch_project_id = ?", mchProjectId).
		Where(map[string]interface{}{"trade_state": tradeState}).
		Find(&payoutCountList).Error; err != nil {
		return 0, appError.NewError(err.Error())
	}
	return payoutCountList[0].Count, nil
}

func (s *PayoutServer) BeneficiaryInfo(params map[string]interface{}) (*model.AspBeneficiary, *appError.Error) {
	o := database.DB
	var beneficiaryInfo model.AspBeneficiary

	for k, v := range params {
		switch k {
		case "id":
			o = o.Where("id = ?", v)
		case "mch_project_benefiary":
			o = o.Where("mch_project_benefiary = ?", v)
		case "customer_name":
			o = o.Where("customer_name = ?", v)
		case "customer_phone":
			o = o.Where("customer_phone = ?", v)
		case "bank_card":
			o = o.Where("bank_card = ?", v)
		case "ifsc":
			o = o.Where("ifsc = ?", v)
		case "provider":
			o = o.Where("provider = ?", v)
		case "adapter":
			o = o.Where("adapter = ?", v)
		case "trade_type":
			o = o.Where("trade_type = ?", v)
		case "benefiary_id":
			o = o.Where("benefiary_id = ?", v)
		}
	}

	if err := o.First(&beneficiaryInfo).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appError.CodeInvalidParamErrCode
	}
	return &beneficiaryInfo, nil
}

// SolveBeneficiarySuccess 处理更新受益人
func (s *PayoutServer) SolveBeneficiarySuccess(beneficiaryId int, params map[string]interface{}, updateStructType string, tx *gorm.DB) (*model.AspBeneficiary, *appError.Error) {
	o := database.DB
	if tx != nil {
		o = tx
	}

	beneficiaryInfoParam := make(map[string]interface{})
	beneficiaryInfoParam["id"] = beneficiaryId

	beneficiaryInfo, err := s.BeneficiaryInfo(beneficiaryInfoParam)

	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "SolveBeneficiarySuccess: ", zap.Error(err))
		return nil, err
	}

	b, _ := json.Marshal(beneficiaryInfo)
	logger.ApiInfo(s.LogFileName, s.RequestId, "update Beneficiary before: ", zap.String("updateStructType", updateStructType), zap.Any("params", params), zap.String("before", string(b)))

	aspBeneficiary := model.AspBeneficiary{}

	resultBeneficiaryUpdate := o.Model(&aspBeneficiary).Where("id = ?", beneficiaryId)

	// 参数 状态 审核状态
	_, tradeStateOk := params["trade_state"]

	if tradeStateOk {
		switch params["trade_state"] {
		case constant.BENEFICIARY_TRADE_STATE_SUCCESS:
			resultBeneficiaryUpdate.Where("trade_state = ?", constant.BENEFICIARY_TRADE_STATE_PENDING)
			break
		}
	}

	resultBeneficiaryUpdate.Updates(params)

	if resultBeneficiaryUpdate.RowsAffected < 1 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "Update Order error: ", zap.String("err: ", "update Beneficiary rows error"))
		return nil, appError.NewError(constant.ChangeErrMsg).FormatMessage("Beneficiary RowsAffected") // 更新代收订单行数 请重新提交
	}

	if resultBeneficiaryUpdate.Error != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "Beneficiary.update() ", zap.Error(resultBeneficiaryUpdate.Error))
		return nil, appError.NewError(resultBeneficiaryUpdate.Error.Error())
	}

	afterBeneficiaryInfo, _ := s.BeneficiaryInfo(beneficiaryInfoParam)
	b, _ = json.Marshal(afterBeneficiaryInfo)
	logger.ApiInfo(s.LogFileName, s.RequestId, "update Beneficiary after: ", zap.String("updateStructType", updateStructType), zap.String("after", string(b)))
	return afterBeneficiaryInfo, nil
}

// PayoutUpiValidate 验证upi是否合法
func (s *PayoutServer) PayoutUpiValidate(d *req.AspPayout, payoutInfo *model.AspPayout, merchantProjectInfo *model.AspMerchantProject, channelDepartInfo *model.AspChannelDepartConfig, currencyId int) *appError.Error {
	// 首先查询是否有验证过的数据
	if upiInfoErr := s.GetPayoutUpiValidate(d); upiInfoErr != nil {
		// 请求到上游，验证upi合法性
		_, upiErr := s.RequestUpstreamPayoutUpiValidate(d, channelDepartInfo)
		if upiErr != nil {
			return upiErr
		}
		// 记录缓存，写入到数据库
		_, insUpiErr := s.InsertUpiValidate(payoutInfo, merchantProjectInfo, channelDepartInfo, currencyId)
		if insUpiErr != nil {
			return insUpiErr
		}
	}
	return nil
}

// GetPayoutBenefiary 判断受益人信息是否合法
func (s *PayoutServer) GetPayoutUpiValidate(d *req.AspPayout) *appError.Error {
	// 需要的参数 vpa 客户名称 客户手机号 支付方式
	params := make(map[string]interface{})
	params["customer_phone"] = d.CustomerPhone
	params["vpa"] = d.Vpa
	// 查询受益人信息
	payoutUpiValidateInfo, beErr := s.PayoutUpiValidateInfo(params)
	if beErr == nil {
		// 存在 则返回
		if payoutUpiValidateInfo.Id > 0 {
			return nil
		}
	}
	return appError.MissNotFoundErrCode.FormatMessage(constant.MissPayoutUpiValidateErrMsg)
}

func (s *PayoutServer) PayoutUpiValidateInfo(params map[string]interface{}) (*model.AspUpiValidate, *appError.Error) {
	o := database.DB
	var payoutUpiValidateInfo model.AspUpiValidate

	for k, v := range params {
		switch k {
		case "id":
			o = o.Where("id = ?", v)
		case "vpa":
			o = o.Where("vpa = ?", v)
		case "customer_name":
			o = o.Where("customer_name = ?", v)
		case "customer_phone":
			o = o.Where("customer_phone = ?", v)
		case "provider":
			o = o.Where("provider = ?", v)
		case "adapter":
			o = o.Where("adapter = ?", v)
		case "trade_type":
			o = o.Where("trade_type = ?", v)
		}
	}

	if err := o.First(&payoutUpiValidateInfo).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appError.CodeInvalidParamErrCode
	}
	return &payoutUpiValidateInfo, nil
}

// RequestUpstreamPayoutUpiValidate 请求上游执行验证upi
func (s *PayoutServer) RequestUpstreamPayoutUpiValidate(reqPayout *req.AspPayout, channelDepartInfo *model.AspChannelDepartConfig) (*interfaces.ThirdUpiValidate, *appError.Error) {
	// 请求到上游需要的参数结构体
	reqPayoutUpiValidate := new(req.AspPayoutUpiValidate)
	reqPayoutUpiValidate.Vpa = reqPayout.Vpa

	supplier := thirdParty.GetPaySupplierByCode(tradeTypeUpiValidate)
	if supplier == nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "thirdParty.GetPaySupplierByCode error ", zap.String("err", "supplier == nil"))
		return nil, appError.ChannelDepartTradeTypeNotFoundErrCode // 渠道信息不存在
	}

	scanData, err := supplier.UpiValidate(s.RequestId, channelDepartInfo, reqPayoutUpiValidate)
	if err != nil {
		return nil, err
	}
	// 渠道映射系统的统一返回字符串 在 appError中有定义的 map
	supplierStr := constant.TradeTypeFynzonPay + "_" + scanData.Code
	supplierError, ok := appError.SupplierErrorMap[supplierStr]
	if !ok {
		logger.ApiWarn(s.LogFileName, s.RequestId, "Response json new Status ", zap.String("newCode", supplierStr))
		return nil, appError.CodeSupplierInternalChannelErrCode
	}
	if supplierError.Code != appError.SUCCESS.Code {
		return nil, supplierError
	}
	return scanData, nil
}

// InsertUpiValidate 写入到upi验证表
func (s *PayoutServer) InsertUpiValidate(payoutInfo *model.AspPayout, merchantProjectInfo *model.AspMerchantProject, channelDepartInfo *model.AspChannelDepartConfig, currencyId int) (*model.AspUpiValidate, *appError.Error) {
	var data model.AspUpiValidate

	data.MchId = merchantProjectInfo.MchId
	data.DepartId = cast.ToInt(channelDepartInfo.DepartId)
	data.ChannelId = cast.ToUint(channelDepartInfo.ChannelId)
	data.CurrencyId = currencyId
	data.MchProjectId = merchantProjectInfo.Id
	data.SpbillCreateIp = s.C.IP()
	data.CustomerId = payoutInfo.CustomerName
	data.CustomerName = payoutInfo.CustomerName
	data.CustomerEmail = payoutInfo.CustomerEmail
	data.CustomerPhone = payoutInfo.CustomerPhone
	data.Vpa = payoutInfo.Vpa
	data.TradeState = constant.BENEFICIARY_TRADE_STATE_SUCCESS
	data.Provider = payoutInfo.Provider
	data.Adapter = payoutInfo.Adapter
	data.TradeType = payoutInfo.TradeType
	data.CreateTime = cast.ToUint64(goutils.GetDateTimeUnix())
	data.FinishTime = data.CreateTime
	data.UnionpayAppend = "{}"
	data.SupplierReturnCode = ""
	data.SupplierReturnMsg = ""
	err := database.DB.Create(&data).Error
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspUpiValidate Insert: ", zap.Error(err))
		return nil, appError.NewError(err.Error())
	}
	return &data, nil
}
