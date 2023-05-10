package repository

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MerchantProjectRepository struct {
	RequestId   string
	LogFileName string
}

func NewMerchantProjectRepository(LogFileName, RequestId string) *MerchantProjectRepository {
	return &MerchantProjectRepository{RequestId: RequestId, LogFileName: LogFileName}
}

// ChangeMerchantProjectCurrent 更新商户余额
func (r *MerchantProjectRepository) ChangeMerchantProjectCurrent(merchantProjectId, availableTotalFee, totalFee, freezeFee, businessType, businessTypeSourceId int, remark string, tx *gorm.DB) *appError.Error {
	// 记录日志
	logger.ApiWarn(r.LogFileName, r.RequestId, "ChangeMerchantProjectCurrent Params", zap.Int("merchantProjectId", merchantProjectId), zap.Int("availableTotalFee", availableTotalFee), zap.Int("totalFee", totalFee), zap.Int("freezeFee", freezeFee), zap.Int("businessType", businessType), zap.Int("businessTypeSourceId", businessTypeSourceId), zap.String("remark", remark))

	if tx == nil {
		tx = database.DB
	}
	merchantProjectCurrencyBefore := model.AspMerchantProjectCurrency{}
	err := tx.Where("mch_project_id = ?", merchantProjectId).First(&merchantProjectCurrencyBefore).Error
	if err != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "FindBefore merchantAccount error: ", zap.Error(err))
		return appError.NewError(err.Error())
	}
	merchantProjectCurrencyUpdate := model.AspMerchantProjectCurrency{}
	merchantProjectCurrencyQuery := tx.Model(&merchantProjectCurrencyUpdate).Where("mch_project_id = ?", merchantProjectId)
	updateData := map[string]interface{}{}
	// 可用余额 冻结
	// 可用余额 解冻
	if availableTotalFee > 0 {
		updateData["available_total_fee"] = gorm.Expr("available_total_fee + ?", availableTotalFee)
	}
	if availableTotalFee < 0 {
		operateAvailableTotalFee := availableTotalFee
		// 则传入的是 负值 （负 - 负） 变成了增加了 应该是减去操作的余额
		operateAvailableTotalFee = -(availableTotalFee)
		updateData["available_total_fee"] = gorm.Expr("available_total_fee - ?", operateAvailableTotalFee)
		merchantProjectCurrencyQuery.Where("available_total_fee >= ?", operateAvailableTotalFee)
	}
	// 冻结金额 冻结
	// 冻结金额 解冻
	if freezeFee > 0 {
		updateData["freeze_fee"] = gorm.Expr("freeze_fee + ?", freezeFee)
	}
	if freezeFee < 0 {
		operateFreezeFee := freezeFee
		// 则传入的是 负值 （负 - 负） 变成了增加了 应该是减去操作的余额
		operateFreezeFee = -(freezeFee)
		updateData["freeze_fee"] = gorm.Expr("freeze_fee - ?", operateFreezeFee)
		merchantProjectCurrencyQuery.Where("freeze_fee >= ?", operateFreezeFee)
	}
	if totalFee > 0 {
		updateData["total_fee"] = gorm.Expr("total_fee + ?", totalFee)
	}

	if totalFee < 0 {
		operateTotalFee := totalFee
		// 则传入的是 负值 （负 - 负） 变成了增加了 应该是减去操作的余额
		operateTotalFee = -(totalFee)
		updateData["total_fee"] = gorm.Expr("total_fee - ?", operateTotalFee)
		merchantProjectCurrencyQuery.Where("total_fee >= ?", operateTotalFee)
	}
	//账户操作
	resultMerchantProjectCurrencyUpdate := merchantProjectCurrencyQuery.Updates(updateData)
	// 判断影响的行数
	if resultMerchantProjectCurrencyUpdate.RowsAffected < 1 {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.String("err: ", constant.UpdateMerchantProjectCurrencyRowsErrMsg))
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent RowsAffected") // 更新商户余额错误 请重新提交
	}
	if resultMerchantProjectCurrencyUpdate.Error != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.Error(err))
		return appError.NewError(err.Error())
	}
	// 记录金额流水
	errRecord := r.MerchantProjectCurrencyRecord(merchantProjectId, availableTotalFee, totalFee, freezeFee, businessType, businessTypeSourceId, remark, tx)
	if errRecord != nil {
		return errRecord
	}
	return nil
}

// PayinOrderSuccess 代收成功
func (r *MerchantProjectRepository) PayinOrderSuccess(merchantProjectId, totalFee, sourceId int, tx *gorm.DB) *appError.Error {
	// 记录日志
	logger.ApiWarn(r.LogFileName, r.RequestId, "ChangeMerchantProjectCurrent Params", zap.Int("merchantProjectId", merchantProjectId), zap.Int("totalFee", totalFee), zap.Int("sourceId", sourceId))

	// 保证接收到的参数是正整数
	if totalFee <= 0 {
		totalFee = -(totalFee)
	}

	if tx == nil {
		tx = database.DB
	}
	merchantProjectCurrencyBefore := model.AspMerchantProjectCurrency{}
	err := tx.Where("mch_project_id = ?", merchantProjectId).First(&merchantProjectCurrencyBefore).Error
	if err != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "FindBefore merchantAccount error: ", zap.Error(err))
		return appError.NewError(err.Error())
	}
	merchantProjectCurrencyUpdate := model.AspMerchantProjectCurrency{}
	merchantProjectCurrencyQuery := tx.Model(&merchantProjectCurrencyUpdate).Where("mch_project_id = ?", merchantProjectId)
	updateData := make(map[string]interface{})

	updateData["total_fee"] = gorm.Expr("total_fee + ?", totalFee)
	//账户操作
	resultMerchantProjectCurrencyUpdate := merchantProjectCurrencyQuery.Updates(updateData)
	// 判断影响的行数
	if resultMerchantProjectCurrencyUpdate.RowsAffected < 1 {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.String("err: ", constant.UpdateMerchantProjectCurrencyRowsErrMsg))
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent RowsAffected") // 更新商户余额错误 请重新提交
	}
	if resultMerchantProjectCurrencyUpdate.Error != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.Error(err))
		return appError.NewError(err.Error())
	}
	// 记录金额流水
	errRecord := r.CapitalFlowInsert(merchantProjectId, 0, totalFee, 0, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYORDER, sourceId, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_SUCCESS, tx)
	if errRecord != nil {
		return errRecord
	}
	return nil
}

// PayoutOrderFreezeFee 代付申请冻结余额
func (r *MerchantProjectRepository) PayoutOrderFreezeFee(merchantProjectId, changeAmount, sourceId int, tx *gorm.DB) *appError.Error {
	// 记录日志
	logger.ApiWarn(r.LogFileName, r.RequestId, "ChangeMerchantProjectCurrent Params", zap.Int("merchantProjectId", merchantProjectId), zap.Int("changeAmount", changeAmount), zap.Int("sourceId", sourceId))
	// 保证接收到的参数是正整数
	if changeAmount <= 0 {
		changeAmount = -(changeAmount)
	}

	if tx == nil {
		tx = database.DB
	}
	merchantProjectCurrencyBefore := model.AspMerchantProjectCurrency{}
	err := tx.Where("mch_project_id = ?", merchantProjectId).First(&merchantProjectCurrencyBefore).Error
	if err != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "FindBefore merchantAccount error: ", zap.Error(err))
		return appError.NewError(err.Error())
	}
	merchantProjectCurrencyUpdate := model.AspMerchantProjectCurrency{}
	merchantProjectCurrencyQuery := tx.Model(&merchantProjectCurrencyUpdate).Where("mch_project_id = ?", merchantProjectId)
	updateData := make(map[string]interface{})
	// 可用余额 冻结
	updateData["available_total_fee"] = gorm.Expr("available_total_fee - ?", changeAmount)
	merchantProjectCurrencyQuery.Where("available_total_fee >= ?", changeAmount)
	// 冻结金额 冻结
	updateData["freeze_fee"] = gorm.Expr("freeze_fee + ?", changeAmount)

	//账户操作
	resultMerchantProjectCurrencyUpdate := merchantProjectCurrencyQuery.Updates(updateData)
	// 判断影响的行数
	if resultMerchantProjectCurrencyUpdate.RowsAffected < 1 {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.String("err: ", constant.UpdateMerchantProjectCurrencyRowsErrMsg))
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent RowsAffected") // 更新商户余额错误 请重新提交
	}
	if resultMerchantProjectCurrencyUpdate.Error != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.Error(err))
		return appError.NewError(err.Error())
	}
	changeAvailableTotalFee := -changeAmount
	changeTotalFee := 0
	changeFreezeFee := changeAmount
	// 记录金额流水
	errPreCapitalFlow := r.PreCapitalFlowInsert(merchantProjectId, changeAvailableTotalFee, changeTotalFee, changeFreezeFee, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_FREEZE, sourceId, merchantProjectCurrencyBefore.MchId, merchantProjectCurrencyBefore.Currency, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_FREEZE_SUCCESS, tx)
	if errPreCapitalFlow != nil {
		return errPreCapitalFlow
	}
	errCapitalFlow := r.CapitalFlowInsert(merchantProjectId, changeAvailableTotalFee, changeTotalFee, changeFreezeFee, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_FREEZE, sourceId, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_FREEZE_SUCCESS, tx)
	if errCapitalFlow != nil {
		return errCapitalFlow
	}
	return nil
}

// PayoutOrderChannelSuccess 代付请求上游返回成功
func (r *MerchantProjectRepository) PayoutOrderChannelSuccess(merchantProjectId, changeAmount, sourceId int, tx *gorm.DB) *appError.Error {
	// 记录日志
	logger.ApiWarn(r.LogFileName, r.RequestId, "ChangeMerchantProjectCurrent Params", zap.Int("merchantProjectId", merchantProjectId), zap.Int("changeAmount", changeAmount), zap.Int("sourceId", sourceId))
	// 保证接收到的参数是正整数
	if changeAmount <= 0 {
		changeAmount = -(changeAmount)
	}

	if tx == nil {
		tx = database.DB
	}
	merchantProjectCurrencyBefore := model.AspMerchantProjectCurrency{}
	err := tx.Where("mch_project_id = ?", merchantProjectId).First(&merchantProjectCurrencyBefore).Error
	if err != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "FindBefore merchantAccount error: ", zap.Error(err))
		return appError.NewError(err.Error())
	}
	merchantProjectCurrencyUpdate := model.AspMerchantProjectCurrency{}
	merchantProjectCurrencyQuery := tx.Model(&merchantProjectCurrencyUpdate).Where("mch_project_id = ?", merchantProjectId)
	updateData := make(map[string]interface{})
	updateData["freeze_fee"] = gorm.Expr("freeze_fee - ?", changeAmount)
	merchantProjectCurrencyQuery.Where("freeze_fee >= ?", changeAmount)
	//账户操作
	resultMerchantProjectCurrencyUpdate := merchantProjectCurrencyQuery.Updates(updateData)
	// 判断影响的行数
	if resultMerchantProjectCurrencyUpdate.RowsAffected < 1 {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.String("err: ", constant.UpdateMerchantProjectCurrencyRowsErrMsg))
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent RowsAffected") // 更新商户余额错误 请重新提交
	}
	if resultMerchantProjectCurrencyUpdate.Error != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.Error(err))
		return appError.NewError(err.Error())
	}
	changeAvailableTotalFee := 0
	changeTotalFee := 0
	changeFreezeFee := -changeAmount
	// 记录金额流水
	errPreCapitalFlow := r.PreCapitalFlowInsert(merchantProjectId, changeAvailableTotalFee, changeTotalFee, changeFreezeFee, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_UNFREEZE, sourceId, merchantProjectCurrencyBefore.MchId, merchantProjectCurrencyBefore.Currency, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_UNFREEZE_SUCCESS, tx)
	if errPreCapitalFlow != nil {
		return errPreCapitalFlow
	}
	errCapitalFlow := r.CapitalFlowInsert(merchantProjectId, changeAvailableTotalFee, changeTotalFee, changeFreezeFee, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_UNFREEZE, sourceId, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_UNFREEZE_SUCCESS, tx)
	if errCapitalFlow != nil {
		return errCapitalFlow
	}
	return nil
}

// PayoutOrderChannelFailed 代付请求上游返回失败
func (r *MerchantProjectRepository) PayoutOrderChannelFailed(merchantProjectId, changeAmount, sourceId int, tx *gorm.DB) *appError.Error {
	// 记录日志
	logger.ApiWarn(r.LogFileName, r.RequestId, "ChangeMerchantProjectCurrent Params", zap.Int("merchantProjectId", merchantProjectId), zap.Int("changeAmount", changeAmount), zap.Int("sourceId", sourceId))
	// 保证接收到的参数是正整数
	if changeAmount <= 0 {
		changeAmount = -(changeAmount)
	}

	if tx == nil {
		tx = database.DB
	}
	merchantProjectCurrencyBefore := model.AspMerchantProjectCurrency{}
	err := tx.Where("mch_project_id = ?", merchantProjectId).First(&merchantProjectCurrencyBefore).Error
	if err != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "FindBefore merchantAccount error: ", zap.Error(err))
		return appError.NewError(err.Error())
	}
	merchantProjectCurrencyUpdate := model.AspMerchantProjectCurrency{}
	merchantProjectCurrencyQuery := tx.Model(&merchantProjectCurrencyUpdate).Where("mch_project_id = ?", merchantProjectId)
	updateData := map[string]interface{}{}
	// 可用余额 冻结
	// 可用余额 解冻
	updateData["available_total_fee"] = gorm.Expr("available_total_fee + ?", changeAmount)
	updateData["freeze_fee"] = gorm.Expr("freeze_fee - ?", changeAmount)
	merchantProjectCurrencyQuery.Where("freeze_fee >= ?", changeAmount)
	//账户操作
	resultMerchantProjectCurrencyUpdate := merchantProjectCurrencyQuery.Updates(updateData)
	// 判断影响的行数
	if resultMerchantProjectCurrencyUpdate.RowsAffected < 1 {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.String("err: ", constant.UpdateMerchantProjectCurrencyRowsErrMsg))
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent RowsAffected") // 更新商户余额错误 请重新提交
	}
	if resultMerchantProjectCurrencyUpdate.Error != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.Error(err))
		return appError.NewError(err.Error())
	}
	changeAvailableTotalFee := changeAmount
	changeTotalFee := 0
	changeFreezeFee := -changeAmount
	// 记录金额流水
	errCapitalFlow := r.CapitalFlowInsert(merchantProjectId, changeAvailableTotalFee, changeTotalFee, changeFreezeFee, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_UNFREEZE, sourceId, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_UNFREEZE_FAILED, tx)
	if errCapitalFlow != nil {
		return errCapitalFlow
	}
	errPreCapitalFlow := r.PreCapitalFlowInsert(merchantProjectId, changeAvailableTotalFee, changeTotalFee, changeFreezeFee, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_UNFREEZE, sourceId, merchantProjectCurrencyBefore.MchId, merchantProjectCurrencyBefore.Currency, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_UNFREEZE_SUCCESS, tx)
	if errPreCapitalFlow != nil {
		return errPreCapitalFlow
	}
	return nil
}

// PayoutOrderAuditReturn 代付管理员撤销
func (r *MerchantProjectRepository) PayoutOrderAuditReturn(merchantProjectId, changeAmount, sourceId int, tx *gorm.DB) *appError.Error {
	// 记录日志
	logger.ApiWarn(r.LogFileName, r.RequestId, "ChangeMerchantProjectCurrent Params", zap.Int("merchantProjectId", merchantProjectId), zap.Int("changeAmount", changeAmount), zap.Int("sourceId", sourceId))
	// 保证接收到的参数是正整数
	if changeAmount <= 0 {
		changeAmount = -(changeAmount)
	}
	if tx == nil {
		tx = database.DB
	}
	merchantProjectCurrencyBefore := model.AspMerchantProjectCurrency{}
	err := tx.Where("mch_project_id = ?", merchantProjectId).First(&merchantProjectCurrencyBefore).Error
	if err != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "FindBefore merchantAccount error: ", zap.Error(err))
		return appError.NewError(err.Error())
	}
	merchantProjectCurrencyUpdate := model.AspMerchantProjectCurrency{}
	merchantProjectCurrencyQuery := tx.Model(&merchantProjectCurrencyUpdate).Where("mch_project_id = ?", merchantProjectId)
	updateData := map[string]interface{}{}
	// 可用余额 冻结
	// 可用余额 解冻
	updateData["available_total_fee"] = gorm.Expr("available_total_fee + ?", changeAmount)
	updateData["freeze_fee"] = gorm.Expr("freeze_fee - ?", changeAmount)
	merchantProjectCurrencyQuery.Where("freeze_fee >= ?", changeAmount)
	//账户操作
	resultMerchantProjectCurrencyUpdate := merchantProjectCurrencyQuery.Updates(updateData)
	// 判断影响的行数
	if resultMerchantProjectCurrencyUpdate.RowsAffected < 1 {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.String("err: ", constant.UpdateMerchantProjectCurrencyRowsErrMsg))
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent RowsAffected") // 更新商户余额错误 请重新提交
	}
	if resultMerchantProjectCurrencyUpdate.Error != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.Error(err))
		return appError.NewError(err.Error())
	}
	changeAvailableTotalFee := changeAmount
	changeTotalFee := 0
	changeFreezeFee := -changeAmount
	// 记录金额流水
	errCapitalFlow := r.CapitalFlowInsert(merchantProjectId, changeAvailableTotalFee, changeTotalFee, changeFreezeFee, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_UNFREEZE, sourceId, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_UNFREEZE_REVOKE, tx)
	if errCapitalFlow != nil {
		return errCapitalFlow
	}
	errPreCapitalFlow := r.PreCapitalFlowInsert(merchantProjectId, changeAvailableTotalFee, changeTotalFee, changeFreezeFee, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_UNFREEZE, sourceId, merchantProjectCurrencyBefore.MchId, merchantProjectCurrencyBefore.Currency, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_UNFREEZE_SUCCESS, tx)
	if errPreCapitalFlow != nil {
		return errPreCapitalFlow
	}
	return nil
}

// ChangeTotalFeeToAvailableTotalFee 每天待结算余额转可用余额 待结算余额减少 可用余额增加 记录流水
func (r *MerchantProjectRepository) ChangeTotalFeeToAvailableTotalFee(merchantProjectId, changeAmount, sourceId int, tx *gorm.DB) *appError.Error {
	// 记录日志
	logger.ApiWarn(r.LogFileName, r.RequestId, "ChangeMerchantProjectCurrent Params", zap.Int("merchantProjectId", merchantProjectId), zap.Int("changeAmount", changeAmount), zap.Int("sourceId", sourceId))
	// 保证接收到的参数是正整数
	if changeAmount <= 0 {
		changeAmount = -(changeAmount)
	}

	if tx == nil {
		tx = database.DB
	}
	merchantProjectCurrencyBefore := model.AspMerchantProjectCurrency{}
	err := tx.Where("mch_project_id = ?", merchantProjectId).First(&merchantProjectCurrencyBefore).Error
	if err != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "FindBefore merchantAccount error: ", zap.Error(err))
		return appError.NewError(err.Error())
	}
	merchantProjectCurrencyUpdate := model.AspMerchantProjectCurrency{}
	merchantProjectCurrencyQuery := tx.Model(&merchantProjectCurrencyUpdate).Where("mch_project_id = ?", merchantProjectId)
	updateData := map[string]interface{}{}

	updateData["available_total_fee"] = gorm.Expr("available_total_fee + ?", changeAmount)

	// 则传入的是 负值 （负 - 负） 变成了增加了 应该是减去操作的余额
	updateData["total_fee"] = gorm.Expr("total_fee - ?", changeAmount)
	merchantProjectCurrencyQuery.Where("total_fee >= ?", changeAmount)
	//账户操作
	resultMerchantProjectCurrencyUpdate := merchantProjectCurrencyQuery.Updates(updateData)
	// 判断影响的行数
	if resultMerchantProjectCurrencyUpdate.RowsAffected < 1 {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.String("err: ", constant.UpdateMerchantProjectCurrencyRowsErrMsg))
		return appError.NewError(constant.ChangeErrMsg).FormatMessage("MerchantProjectCurrent RowsAffected") // 更新商户余额错误 请重新提交
	}
	if resultMerchantProjectCurrencyUpdate.Error != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.Error(err))
		return appError.NewError(err.Error())
	}
	// 记录金额流水
	changeAvailableTotalFee := changeAmount
	changeTotalFee := -changeAmount
	changeFreezeFee := 0
	// 记录金额流水
	errCapitalFlow := r.CapitalFlowInsert(merchantProjectId, changeAvailableTotalFee, changeTotalFee, changeFreezeFee, constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_TO_AVAILABLE_TOTAL_FEE, sourceId, constant.MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_TO_AVAILABLE_TOTAL_FEE, tx)
	if errCapitalFlow != nil {
		return errCapitalFlow
	}
	return nil
}

// ChangeMerchantProjectCurrent004
/*
* 账户金额操作
// 添加 账户余额变动队列
// 预扣金额 的情况 审核代付 提现中的订单
// 金额变动操作
//
//	 操作类别          可用余额        冻结余额     待结算余额       结算中余额      预扣金额记录      收支流水记录
//	备注             （脚本每天会统计）
//
// 代收成功              无             无          增加           无              无              新增

// 代付审核成功后则
// 代付申请冻结余额       减去           增加        无             无               新增             无
// 代付审核失败           无            无          无             无               无             无    （审核失败后 订单则关闭 ，不可能再申请）
// 代付上游成功           无            减去          无            无             新增             扣减
// 管理员撤销             增加           减去        无            无               新增              无 （如果）
// 代付上游失败           增加           减去        无            无               新增              无 （如果）

// 提现申请              减去           增加        无             新增             新增             无
// 提现成功              无             减去        无             减去             新增             扣减
// 提现失败              增加           减去        无             减去             新增            无
// 财务操作
// 每天执行的脚本 统计 可用余额 和 待结算余额 计算截止前一天的时间 统计总的待结算余额 然后 可用余额 增加 待结算余额 减少 保证在一个事务中执行
@param merchantProjectId cp项目id
@param availableTotalFee 操作的可用余额代付和提现使用 增加 正数 减少 负数
@param totalFee 待结算余额 增加 正数 减少 负数
@param freezeFee 冻结金额 增加 正数 减少 负数
@param businessType 业务类型 1:代收 2:代付 3:充值 4:结算
@param businessTypeSourceId 具体业务类型来源表关联id (类型关联表 1: asp_order 2: asp_payout 3: asp_merchant_project_withdrawal)
@param remark 备注
@param params 代收或者代付 修改对应的参数
@ 注意： 当金额传的是负数的情况 已经处理了 （负 - 负）得正 则需要处理
*/

func (r *MerchantProjectRepository) MerchantProjectCurrencyRecord(merchantProjectId int, availableTotalFee int, totalFee int, freezeFee int, businessType int, businessTypeSourceId int, remark string, tx *gorm.DB) *appError.Error {
	merchantProjectCurrencyAfter := model.AspMerchantProjectCurrency{}
	err := tx.Where("mch_project_id = ?", merchantProjectId).First(&merchantProjectCurrencyAfter).Error
	if err != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "FindAfter merchantAccount error: ", zap.Error(err))
		return appError.NewError(err.Error())
	}

	// 当代付的情况 记录 冻结记录
	if businessType == constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_FREEZE || businessType == constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_UNFREEZE {
		preFlowStatus := constant.MERCHANT_PROJECT_PRE_FLOW_BUSINESS_STATUS_SUCCESS // 冻结记录锁定的状态(1.锁定中
		// 冻结记录状态为 解锁
		if businessType == constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_UNFREEZE {
			preFlowStatus = constant.MERCHANT_PROJECT_PRE_FLOW_BUSINESS_STATUS_FAILD
		}
		var aspMerchantProjectPreFlowModel model.AspMerchantProjectPreFlow
		aspMerchantProjectPreFlowModel.MchId = merchantProjectCurrencyAfter.MchId
		aspMerchantProjectPreFlowModel.MchProjectId = merchantProjectId
		aspMerchantProjectPreFlowModel.MchProjectCurrencyId = merchantProjectId
		aspMerchantProjectPreFlowModel.Currency = merchantProjectCurrencyAfter.Currency
		aspMerchantProjectPreFlowModel.PreTotalFee = cast.ToInt64(freezeFee)
		aspMerchantProjectPreFlowModel.BusinessType = businessType
		aspMerchantProjectPreFlowModel.BusinessSourceId = cast.ToInt64(businessTypeSourceId)
		aspMerchantProjectPreFlowModel.Status = preFlowStatus
		aspMerchantProjectPreFlowModel.Remark = remark
		aspMerchantProjectPreFlowModel.CreateTime = goutils.GetDateTimeUnix()
		aspMerchantProjectPreFlowModel.UpdateTime = goutils.GetDateTimeUnix()
		// 赋值给 order 数据
		//d.Generate(&data, aspMerchantProject, aspMerchantProjectConfig, deptTradeTypeInfo)
		err = tx.Create(&aspMerchantProjectPreFlowModel).Error
		if err != nil {
			logger.ApiWarn(r.LogFileName, r.RequestId, "AspMerchantProjectPreFlow Insert: ", zap.Error(err))
			return appError.NewError(err.Error())
		}
	}

	errCapitalFlow := r.CapitalFlowInsert(merchantProjectId, availableTotalFee, totalFee, freezeFee, businessType, businessTypeSourceId, remark, tx)
	if errCapitalFlow != nil {
		return errCapitalFlow
	}
	return nil
}

func (r *MerchantProjectRepository) PreCapitalFlowInsert(merchantProjectId int, availableTotalFee int, totalFee int, freezeFee int, businessType int, businessTypeSourceId int, mchId int, currency string, remark string, tx *gorm.DB) *appError.Error {
	preFlowStatus := constant.MERCHANT_PROJECT_PRE_FLOW_BUSINESS_STATUS_SUCCESS // 冻结记录锁定的状态(1.锁定中
	// 冻结记录状态为 解锁
	if businessType == constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_UNFREEZE {
		preFlowStatus = constant.MERCHANT_PROJECT_PRE_FLOW_BUSINESS_STATUS_FAILD
	}
	var aspMerchantProjectPreFlowModel model.AspMerchantProjectPreFlow
	aspMerchantProjectPreFlowModel.MchId = mchId
	aspMerchantProjectPreFlowModel.MchProjectId = merchantProjectId
	aspMerchantProjectPreFlowModel.MchProjectCurrencyId = merchantProjectId
	aspMerchantProjectPreFlowModel.Currency = currency
	aspMerchantProjectPreFlowModel.PreTotalFee = cast.ToInt64(freezeFee)
	aspMerchantProjectPreFlowModel.BusinessType = businessType
	aspMerchantProjectPreFlowModel.BusinessSourceId = cast.ToInt64(businessTypeSourceId)
	aspMerchantProjectPreFlowModel.Status = preFlowStatus
	aspMerchantProjectPreFlowModel.Remark = remark
	aspMerchantProjectPreFlowModel.CreateTime = goutils.GetDateTimeUnix()
	aspMerchantProjectPreFlowModel.UpdateTime = goutils.GetDateTimeUnix()
	// 赋值给 order 数据
	//d.Generate(&data, aspMerchantProject, aspMerchantProjectConfig, deptTradeTypeInfo)
	err := tx.Create(&aspMerchantProjectPreFlowModel).Error
	if err != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "AspMerchantProjectPreFlow Insert: ", zap.Error(err))
		return appError.NewError(err.Error())
	}
	return nil
}

func (r *MerchantProjectRepository) CapitalFlowInsert(merchantProjectId int, availableTotalFee int, totalFee int, freezeFee int, businessType int, businessTypeSourceId int, remark string, tx *gorm.DB) *appError.Error {
	merchantProjectCurrencyAfter := model.AspMerchantProjectCurrency{}
	err := tx.Where("mch_project_id = ?", merchantProjectId).First(&merchantProjectCurrencyAfter).Error
	if err != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "FindAfter merchantAccount error: ", zap.Error(err))
		return appError.NewError(err.Error())
	}

	var aspMerchantProjectCapitalFlowModel model.AspMerchantProjectCapitalFlow
	aspMerchantProjectCapitalFlowModel.MchId = merchantProjectCurrencyAfter.MchId
	aspMerchantProjectCapitalFlowModel.MchProjectId = merchantProjectId
	aspMerchantProjectCapitalFlowModel.AccountCurrencyId = cast.ToInt(merchantProjectCurrencyAfter.CurrencyId)
	aspMerchantProjectCapitalFlowModel.Currency = merchantProjectCurrencyAfter.Currency
	aspMerchantProjectCapitalFlowModel.TotalFee = cast.ToInt64(totalFee) // 待结算余额
	aspMerchantProjectCapitalFlowModel.TotalFeeSurplus = merchantProjectCurrencyAfter.TotalFee
	aspMerchantProjectCapitalFlowModel.TotalFeeBefore = aspMerchantProjectCapitalFlowModel.TotalFeeSurplus - aspMerchantProjectCapitalFlowModel.TotalFee
	aspMerchantProjectCapitalFlowModel.FreezeFee = cast.ToInt64(freezeFee) // 冻结余额
	aspMerchantProjectCapitalFlowModel.FreezeFeeSurplus = merchantProjectCurrencyAfter.FreezeFee
	aspMerchantProjectCapitalFlowModel.FreezeFeeBefore = aspMerchantProjectCapitalFlowModel.FreezeFeeSurplus - aspMerchantProjectCapitalFlowModel.FreezeFee
	aspMerchantProjectCapitalFlowModel.AvailableFee = cast.ToInt64(availableTotalFee) // 可用余额
	aspMerchantProjectCapitalFlowModel.AvailableFeeSurplus = merchantProjectCurrencyAfter.AvailableTotalFee
	aspMerchantProjectCapitalFlowModel.AvailableFeeBefore = aspMerchantProjectCapitalFlowModel.AvailableFeeSurplus - aspMerchantProjectCapitalFlowModel.AvailableFee
	aspMerchantProjectCapitalFlowModel.BusinessType = businessType
	aspMerchantProjectCapitalFlowModel.BusinessSourceId = cast.ToInt64(businessTypeSourceId)
	aspMerchantProjectCapitalFlowModel.Remark = remark
	aspMerchantProjectCapitalFlowModel.CreateTime = goutils.GetDateTimeUnix()
	aspMerchantProjectCapitalFlowModel.UpdateTime = goutils.GetDateTimeUnix()
	err = tx.Create(&aspMerchantProjectCapitalFlowModel).Error
	if err != nil {
		logger.ApiWarn(r.LogFileName, r.RequestId, "AspMerchantProjectCapitalFlow Insert: ", zap.Error(err))
		return appError.NewError(err.Error())
	}
	return nil
}

func (r *MerchantProjectRepository) ChangeMerchantProjectCurrentByTest(merchantProjectId, availableTotalFee, totalFee, freezeFee, businessType, businessTypeSourceId int, remark string, params map[string]interface{}) *appError.Error {
	// 记录日志
	logger.ApiWarn(r.LogFileName, r.RequestId, "ChangeMerchantProjectCurrent Params", zap.Int("merchantProjectId", merchantProjectId), zap.Int("availableTotalFee", availableTotalFee), zap.Int("totalFee", totalFee), zap.Int("freezeFee", freezeFee), zap.Int("businessType", businessType), zap.Int("businessTypeSourceId", businessTypeSourceId), zap.String("remark", remark), zap.Any("params", params))

	// 包含的操作有 更新账户可用余额 待结算余额 冻结余额 记录冻结日志 余额变动日志 更新代收订单状态 更新代付订单状态
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 更新代收
		if businessType == constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYORDER {
			aspOrder := model.AspOrder{}
			err := tx.Model(&aspOrder).Where("id = ?", businessTypeSourceId).Where("trade_state = ? or trade_state = ?", constant.ORDER_TRADE_STATE_PENDING, constant.ORDER_TRADE_STATE_USERPAYING).Updates(params).Error
			if err != nil {
				logger.ApiWarn(r.LogFileName, r.RequestId, "Update Order error: ", zap.Error(err))
				return err
			}
		}

		if businessType == constant.MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_FREEZE {
			// 冻结 解冻代付成功 解冻代付失败 解冻代付取消
			aspPayout := model.AspPayout{}

			payoutUpdate := tx.Model(&aspPayout).Where("id = ?", businessTypeSourceId)

			// 代付成功
			if params["trade_state"] == constant.PAYOUT_TRADE_STATE_SUCCESS {
				payoutUpdate.Where("trade_state = ?", constant.PAYOUT_TRADE_STATE_APPLY)
			}
			//代付失败
			if params["trade_state"] == constant.PAYOUT_TRADE_STATE_FAILED {
				payoutUpdate.Where("trade_state = ?", constant.PAYOUT_TRADE_STATE_APPLY)
			}
			err := payoutUpdate.Updates(params).Error
			if err != nil {
				logger.ApiWarn(r.LogFileName, r.RequestId, "Update aspPayout error: ", zap.Error(err))
				return err
			}
		}

		merchantProjectCurrencyBefore := model.AspMerchantProjectCurrency{}
		err := tx.Where("mch_project_id = ?", merchantProjectId).First(&merchantProjectCurrencyBefore).Error
		if err != nil {
			logger.ApiWarn(r.LogFileName, r.RequestId, "FindBefore merchantAccount error: ", zap.Error(err))
			return err
		}
		b, err := goutils.JsonEncode(merchantProjectCurrencyBefore)
		logger.ApiInfo(r.LogFileName, r.RequestId, "update before ", zap.String("before", b))
		//goutils.Dump(merchantProjectCurrencyBefore)
		merchantProjectCurrencyUpdate := model.AspMerchantProjectCurrency{}
		merchantProjectCurrencyQuery := tx.Model(&merchantProjectCurrencyUpdate).Where("mch_project_id = ?", merchantProjectId)
		updateData := map[string]interface{}{}
		// 可用余额 冻结
		// 可用余额 解冻
		if availableTotalFee > 0 {
			updateData["available_total_fee"] = gorm.Expr("available_total_fee + ?", availableTotalFee)
		}
		if availableTotalFee < 0 {
			operateAvailableTotalFee := availableTotalFee
			// 则传入的是 负值 （负 - 负） 变成了增加了 应该是减去操作的余额
			operateAvailableTotalFee = -(availableTotalFee)
			updateData["available_total_fee"] = gorm.Expr("available_total_fee - ?", operateAvailableTotalFee)
			merchantProjectCurrencyQuery.Where("available_total_fee >= ?", operateAvailableTotalFee)
		}
		// 冻结金额 冻结
		// 冻结金额 解冻
		if freezeFee > 0 {
			updateData["freeze_fee"] = gorm.Expr("freeze_fee + ?", freezeFee)
		}
		if freezeFee < 0 {
			operateFreezeFee := freezeFee
			// 则传入的是 负值 （负 - 负） 变成了增加了 应该是减去操作的余额
			operateFreezeFee = -(freezeFee)
			updateData["freeze_fee"] = gorm.Expr("freeze_fee - ?", operateFreezeFee)
			merchantProjectCurrencyQuery.Where("freeze_fee >= ?", operateFreezeFee)
		}
		if totalFee > 0 {
			updateData["total_fee"] = gorm.Expr("total_fee + ?", totalFee)
		}

		if totalFee < 0 {
			operateTotalFee := totalFee
			// 则传入的是 负值 （负 - 负） 变成了增加了 应该是减去操作的余额
			operateTotalFee = -(totalFee)
			updateData["total_fee"] = gorm.Expr("total_fee - ?", operateTotalFee)
			merchantProjectCurrencyQuery.Where("total_fee >= ?", operateTotalFee)
		}
		//账户操作
		err = merchantProjectCurrencyQuery.Updates(updateData).Error
		if err != nil {
			logger.ApiWarn(r.LogFileName, r.RequestId, "Update merchantAccount Error", zap.Error(err))
			return err
		}
		merchantProjectCurrencyAfter := model.AspMerchantProjectCurrency{}
		err = tx.Where("mch_project_id = ?", merchantProjectId).First(&merchantProjectCurrencyAfter).Error
		if err != nil {
			logger.ApiWarn(r.LogFileName, r.RequestId, "FindAfter merchantAccount error: ", zap.Error(err))
			return err
		}
		b, err = goutils.JsonEncode(merchantProjectCurrencyAfter)
		logger.ApiInfo(r.LogFileName, r.RequestId, "update after ", zap.String("after", b))

		var aspMerchantProjectCapitalFlowModel model.AspMerchantProjectCapitalFlow
		aspMerchantProjectCapitalFlowModel.MchId = merchantProjectCurrencyBefore.MchId
		aspMerchantProjectCapitalFlowModel.MchProjectId = merchantProjectId
		aspMerchantProjectCapitalFlowModel.AccountCurrencyId = cast.ToInt(merchantProjectCurrencyBefore.CurrencyId)
		aspMerchantProjectCapitalFlowModel.Currency = merchantProjectCurrencyBefore.Currency
		aspMerchantProjectCapitalFlowModel.TotalFee = cast.ToInt64(totalFee) // 待结算余额
		aspMerchantProjectCapitalFlowModel.TotalFeeSurplus = merchantProjectCurrencyAfter.TotalFee
		aspMerchantProjectCapitalFlowModel.TotalFeeBefore = aspMerchantProjectCapitalFlowModel.TotalFeeSurplus - aspMerchantProjectCapitalFlowModel.TotalFee
		aspMerchantProjectCapitalFlowModel.FreezeFee = cast.ToInt64(freezeFee) // 冻结余额
		aspMerchantProjectCapitalFlowModel.FreezeFeeSurplus = merchantProjectCurrencyAfter.FreezeFee
		aspMerchantProjectCapitalFlowModel.FreezeFeeBefore = aspMerchantProjectCapitalFlowModel.FreezeFeeSurplus - aspMerchantProjectCapitalFlowModel.FreezeFee
		aspMerchantProjectCapitalFlowModel.AvailableFee = cast.ToInt64(availableTotalFee) // 可用余额
		aspMerchantProjectCapitalFlowModel.AvailableFeeSurplus = merchantProjectCurrencyAfter.AvailableTotalFee
		aspMerchantProjectCapitalFlowModel.AvailableFeeBefore = aspMerchantProjectCapitalFlowModel.AvailableFeeSurplus - aspMerchantProjectCapitalFlowModel.AvailableFee
		aspMerchantProjectCapitalFlowModel.BusinessType = businessType
		aspMerchantProjectCapitalFlowModel.BusinessSourceId = cast.ToInt64(businessTypeSourceId)
		aspMerchantProjectCapitalFlowModel.Remark = remark
		aspMerchantProjectCapitalFlowModel.CreateTime = goutils.GetDateTimeUnix()
		aspMerchantProjectCapitalFlowModel.UpdateTime = goutils.GetDateTimeUnix()
		err = tx.Create(&aspMerchantProjectCapitalFlowModel).Error

		if err != nil {
			logger.ApiWarn(r.LogFileName, r.RequestId, "AspMerchantProjectCapitalFlow Insert: ", zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		return appError.NewError(err.Error())
	}

	return nil
}
