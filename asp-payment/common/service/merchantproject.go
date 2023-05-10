package service

import (
	"asp-payment/api-server/req"
	"asp-payment/api-server/rsp"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goRedis"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	reqCommon "asp-payment/common/req"
	"asp-payment/common/service/supplier"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MerchantProjectServer struct {
	*Service
}

func NewMerchantProjectServer(c *fiber.Ctx) *MerchantProjectServer {
	return &MerchantProjectServer{Service: NewService(c, constant.MerchantServerLogFileName)}
}

func NewSimpleMerchantProjectServer() *MerchantProjectServer {
	return &MerchantProjectServer{&Service{LogFileName: constant.MerchantServerLogFileName}}
}

// MerchantAccountQuery  商户账户信息
func (s *MerchantProjectServer) MerchantAccountQuery() error {
	lockName := goRedis.GetKey(fmt.Sprintf("pay:mchQuery:%s", s.Head.AppId))
	flag := goRedis.Lock(lockName)
	if !flag {
		return appError.IsWaitErrCode
	}
	defer goRedis.UnLock(lockName)
	// 贯穿支付所需要的参数
	reqAspMerchantAccountQuery := &req.AspMerchantAccountQuery{}

	if err := s.C.BodyParser(reqAspMerchantAccountQuery); err != nil {
		return err
	}

	if err := s.DealUnifiedMerchantAccountQueryParams(s.Head); err != nil {
		return err
	}

	// TODO 获取用户信息 后期需要验证 用户的账号和密码

	merchantProjectServer := NewMerchantProjectServer(s.C)
	// appError.MissMerchantProjectNotFoundErr.FormatMessage("xxxx")
	merchantProjectInfo, err := merchantProjectServer.GetMerchantProjectInfo(s.Head.AppId)
	if err != nil {
		return err
	}
	AspIdInfo, err := merchantProjectServer.GetIdInfo()
	if err != nil {
		return err
	}

	// 验证签名 可以使用 默认的 util 来实现 验证签名
	err = s.CheckSignMerchantAccountQuery(reqAspMerchantAccountQuery, AspIdInfo.Key, s.Head.Timestamp)
	if err != nil {
		return err
	}
	// fmt.Println("logic.CheckSignMerchantAccountQuery--------------------", req)
	// 获取 外部商户项目 金额信息
	// 需要的参数 币种id 外部商户项目id 返回 信息
	merchantProjectCurrencyInfo, err := s.GetMerchantProjectCurrencyInfo(reqAspMerchantAccountQuery.Currency, merchantProjectInfo.Id)
	if err != nil {
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return (&MissNotFoundErrCode).FormatMessage(constant.MissMerchantProjectCurrencyNotFoundErrMsg) // 商户项目金额异常
	}
	// fmt.Println("s.GetMerchantProjectCurrencyInfo--------------------", req.AspMerchantProjectCurrency)
	scanSuccessData := rsp.GenerateMerchantAccountQuerySuccessData(merchantProjectCurrencyInfo)
	return s.Success(scanSuccessData)
}

var MerchantAccountTradeType = "MERCHANT_ACCOUNT"

// 需要的入参 查询商户信息 需要有商户的 账号和密码 验证通过后，才可以执行查询上游 传 depart_id 绑定的depart_id 用户 channel_id 渠道id
// 怎么获取对应的渠道 ?? 根据 channel_id 和 depart_id 去查询

// 查找商户是否存在

// 查找渠道是否存在

// MerchantAccountChannelQuery 寻找  channel_depart 表中的 config
// 查询上游数据 需要有 账户的信息
//
//	商户账户信息
func (s *MerchantProjectServer) MerchantAccountChannelQuery() error {
	lockName := goRedis.GetKey("pay:mchAccQuery")
	if !goRedis.Lock(lockName) {
		return appError.IsWaitErrCode
	}
	defer goRedis.UnLock(lockName)
	// 贯穿支付所需要的参数
	reqMerchantAccount := &req.AspMerchantAccount{}
	if err := s.C.BodyParser(reqMerchantAccount); err != nil {
		return appError.NewError(err.Error())
	}

	if err := s.DealUnifiedMerchantAccountParams(reqMerchantAccount); err != nil {
		return s.Error(err)
	}
	// TODO 获取用户信息 后期需要验证 用户的账号和密码
	// fmt.Println("logic.CheckSignMerchantAccount--------------------", req)
	// 获取渠道的信息 过程中有赋值给 AspChannelConfig
	channelConfigInfo, cErr := NewChannelConfigServer(s.C).GetChannelConfigInfo(reqMerchantAccount.ChannelId)
	if cErr != nil {
		return cErr
	}
	// fmt.Println("s.GetChannelConfigInfo--------------------", req)
	// 获取到上游数据 获取到对应的上游
	// supplierCode := req.AspChannelDepartTradeType.TradeType

	supplierCode := channelConfigInfo.Name + "." + MerchantAccountTradeType

	// supplierCode := "firstpay.MERCHANT_ACCOUNT"
	// supplierCode = "firstpay.MERCHANT_ACCOUNT"
	paySupplier := supplier.GetPaySupplierByCode(supplierCode)
	if paySupplier == nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "thirdParty.GetPaySupplierByCode error ", zap.String("err", "supplier == nil"))
		return appError.ChannelDepartTradeTypeNotFoundErrCode // 渠道信息不存在
	}

	channelDepartInfoParam := map[string]string{}
	channelDepartInfoParam["depart_id"] = cast.ToString(reqMerchantAccount.DepartId)
	channelDepartInfoParam["channel_id"] = cast.ToString(reqMerchantAccount.ChannelId)
	//查询账户在上游的配置
	channelDepartInfo, err := NewChannelConfigServer(s.C).GetChannelDepartInfo(channelDepartInfoParam)
	if err != nil {
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartNotFoundErrMsg) // missing channel depart not found
	}

	thirdMerchantAccountQueryData, dErr := paySupplier.GetDepartAccountInfo(s.RequestId, channelDepartInfo)
	if dErr != nil {
		return dErr
	}
	// 渠道映射系统的统一返回字符串 在 appError中有定义的 map
	supplierStr := channelConfigInfo.Name + "_" + thirdMerchantAccountQueryData.Code
	supplierError, ok := supplier.SupplierErrorMap[supplierStr]
	if !ok {
		logger.ApiWarn(s.LogFileName, s.RequestId, "Response json new Status ", zap.String("newCode", supplierStr))
		return appError.CodeSupplierInternalChannelErrCode
	}
	if supplierError.Code != appError.SUCCESS.Code {
		return supplierError
	}
	// 统一返回参数 转换
	scanSuccessData := rsp.GenerateMerchantAccountSuccessData(thirdMerchantAccountQueryData)
	return s.Success(scanSuccessData)
}

// GetMerchantProjectCurrentChannel  获取cp 使用的当前渠道
func (s *MerchantProjectServer) GetMerchantProjectCurrentChannel() error {
	lockName := goRedis.GetKey("pay:mchprochannel")
	if !goRedis.Lock(lockName) {
		return appError.IsWaitErrCode
	}
	defer goRedis.UnLock(lockName)
	// 贯穿查询需要的参数
	reqBody := &req.AspMerchantProjectChannelReq{}
	if bodyErr := s.C.BodyParser(reqBody); bodyErr != nil {
		return s.Error(appError.NewError(bodyErr.Error()))
	}

	if err := s.DealUnifiedMerchantProjectChannelParams(reqBody); err != nil {
		return s.Error(err)
	}
	// 验证是否存在 cp
	_, mchErr := s.GetMerchantProjectInfo(reqBody.MchProId)
	if mchErr != nil {
		return s.Error(mchErr)
	}

	newDepartServer := NewDepartServer(s.C)
	// 查询到 商户的可用的支付渠道 过程中有赋值给 AspChannelDepartTradeType
	aspChannelDepartTradeType, cErr := newDepartServer.ChooseMerchantProjectChannelDepartTradeType(reqBody)
	if cErr != nil {
		return s.Error(cErr)
	}

	// payment = sunny.paytm
	paymentArr := strings.Split(aspChannelDepartTradeType.Payment, ".")
	if len(paymentArr) != 2 {
		return s.Error(appError.MissNotFoundErrCode.FormatMessage(constant.MissChannelDepartPaymentParamErrMsg))
	}

	// 获取到上游数据 获取到对应的上游
	supplierCode := paymentArr[1] + "." + aspChannelDepartTradeType.TradeType
	//fmt.Println("supplierCode: ", supplierCode)
	//return s.Error(appError.ChannelDepartTradeTypeNotFoundErrCode) // 渠道信息不存在
	// supplierCode := "firstpay.H5"
	paySupplier := supplier.GetPaySupplierByCode(supplierCode)
	if paySupplier == nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "thirdParty.GetPaySupplierByCode error ", zap.String("supplierCode", supplierCode), zap.Any("paySupplier", paySupplier))
		return s.Error(appError.ChannelDepartTradeTypeNotFoundErrCode) // 渠道信息不存在
	}

	channelInfo := rsp.GenerateMerchantProjectSuccessData(aspChannelDepartTradeType)
	departId := cast.ToInt(aspChannelDepartTradeType.DepartID)
	departInfo, deErr := newDepartServer.GetDepartInfo(departId)
	if deErr != nil {
		return s.Error(appError.NewError(deErr.Error()))
	}
	channelInfo.DepartName = departInfo.Title
	channelInfo.Payment = paymentArr[1]
	return s.Success(channelInfo)
}

// GetMerchantProjectCurrencyList 获取cp项目余额账户列表
func (s *MerchantProjectServer) GetMerchantProjectCurrencyList() (map[int]*model.AspMerchantProjectCurrency, error) {
	merchantProjectCurrencyPluckList := make(map[int]*model.AspMerchantProjectCurrency, 0)
	merchantProjectCurrencyList := make([]*model.AspMerchantProjectCurrency, 0)
	var merchantProjectCurrencyInfo model.AspMerchantProjectCurrency
	err := database.DB.Model(&merchantProjectCurrencyInfo).Where("status = ?", 1).Find(&merchantProjectCurrencyList).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "merchantProjectCurrencyList err: ", zap.Error(err))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissMerchantProjectCurrencyListNotFoundErrMsg)
	}

	if len(merchantProjectCurrencyList) == 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(merchantProjectCurrencyList) == 0 ", zap.String("err", constant.MissMerchantProjectCurrencyListNotFoundErrMsg))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissMerchantProjectCurrencyListNotFoundErrMsg)
	}

	for i := 0; i < len(merchantProjectCurrencyList); i++ {
		mchProjectId := cast.ToInt(merchantProjectCurrencyList[i].MchProjectId)
		merchantProjectCurrencyPluckList[mchProjectId] = merchantProjectCurrencyList[i]
	}
	return merchantProjectCurrencyPluckList, nil
}

func (s *MerchantProjectServer) GetMerchantProjectInfo(appId string) (*model.AspMerchantProject, *appError.Error) {
	var data model.AspMerchantProject

	redisKey := fmt.Sprintf(constant.KEY_MERCHANT_PROJECT_INFO_STRING, appId)
	redisKey = constant.GetRedisKey(redisKey)

	merchantProjectValue := s.Redis().Get(context.Background(), redisKey).Val()
	if merchantProjectValue != "" {
		// JSON转结构体
		err := goutils.JsonDecode(merchantProjectValue, &data)
		if err != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "merchantProjectValue", zap.Error(err))
			return nil, appError.MissMerchantProjectNotFoundErrMsg
		}
		return &data, nil
	}

	o := database.DB

	err := o.Where("status = ? and id = ? ", 1, appId).First(&data).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appError.MissMerchantProjectNotFoundErrMsg
	}

	s.setMerchantProjectInfoToCache(&data, redisKey)
	return &data, nil
}

// 私有
// 设置Redis
func (s *MerchantProjectServer) setMerchantProjectInfoToCache(Data *model.AspMerchantProject, redisKey string) bool {
	// 转成JSON
	idJsonByte, err := json.Marshal(Data)
	if err != nil {
		return true
	}
	err = s.Redis().Set(context.Background(), redisKey, string(idJsonByte), constant.EXPIRED_TIME).Err()
	return err == nil
}

func (s *MerchantProjectServer) GetIdInfo() (*model.AspId, *appError.Error) {
	appId := s.Head.AppId
	var data *model.AspId

	redisKey := fmt.Sprintf(constant.KEY_APPID_STRING, appId)
	redisKey = constant.GetRedisKey(redisKey)

	idValue := s.Redis().Get(context.Background(), redisKey).Val()
	if idValue != "" {
		// JSON转结构体
		err := goutils.JsonDecode(idValue, &data)
		if err != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "idValue", zap.Error(err))
			return nil, appError.MissIdNotFoundErrCode
		}
		return data, nil
	}

	o := database.DB

	err := o.Where("status = ? and id = ? ", 1, appId).First(&data).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appError.MissIdNotFoundErrCode
	}

	s.SetIdInfoToCache(data, redisKey)
	return data, nil
}

// SetIdInfoToCache 私有
// 设置Redis
func (s *MerchantProjectServer) SetIdInfoToCache(Data *model.AspId, redisKey string) bool {
	// 转成JSON
	idJsonByte, err := json.Marshal(Data)
	if err != nil {
		return true
	}
	err = s.Redis().Set(context.Background(), redisKey, string(idJsonByte), constant.EXPIRED_TIME).Err()
	return err == nil
}

// MerchantProjectUserInsertOrUpdate 处理玩家 数据 判断是否有用户 如果有则更新 如果没有 则插入新的数据
func (s *MerchantProjectServer) MerchantProjectUserInsertOrUpdate(uid string, merchantProjectInfo *model.AspMerchantProject) *appError.Error {

	// 应该首先查询 是否有
	var data model.AspMerchantProjectUser
	err := database.DB.Where("uid = ?", uid).
		Where("mch_id = ?", merchantProjectInfo.MchId).
		Where("mch_project_id = ?", merchantProjectInfo.Id).
		First(&data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		s.Generate(&data, uid, merchantProjectInfo)
		err = database.DB.Create(&data).Error
		if err != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "AspMerchantProjectUser.Create(&data) ", zap.Error(err))
			return appError.CodeUserNotExist
		}
	}

	return nil
}

func (s *MerchantProjectServer) Generate(data *model.AspMerchantProjectUser, uid string, m *model.AspMerchantProject) {
	data.MchId = m.MchId
	data.MchProjectId = m.Id
	data.Uid = uid
	data.Phone = ""
	data.VipLevel = 0
	data.Status = constant.MERCHANT_PROJECT_USER_STATUS_NORMAL
	data.Email = ""
	data.Remark = ""
	data.CreateTime = cast.ToUint64(goutils.GetDateTimeUnix())
	data.UpdateTime = cast.ToUint64(goutils.GetDateTimeUnix())
}

// DealUnifiedMerchantAccountParams 处理请求参数 参数的验证
func (s *MerchantProjectServer) DealUnifiedMerchantAccountParams(d *req.AspMerchantAccount) *appError.Error {

	if err := checker.Struct(d); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedMerchantAccountParams ", zap.Error(err))
		return err
	}

	return nil
}

// DealUnifiedMerchantAccountQueryParams 处理请求参数 参数的验证
func (s *MerchantProjectServer) DealUnifiedMerchantAccountQueryParams(h *req.AspPaymentHeader) *appError.Error {

	if err := checker.Struct(h); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedMerchantAccountQueryParams ", zap.Error(err))
		return err
	}

	return nil
}

// DealUnifiedMerchantProjectChannelParams 处理请求参数 参数的验证
func (s *MerchantProjectServer) DealUnifiedMerchantProjectChannelParams(d *req.AspMerchantProjectChannelReq) *appError.Error {

	if err := checker.Struct(d); err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "DealUnifiedMerchantProjectChannelParams ", zap.Error(err))
		return err
	}

	return nil
}

// CheckSignMerchantAccountQuery 查询商户余额 验签
func (s *MerchantProjectServer) CheckSignMerchantAccountQuery(d *req.AspMerchantAccountQuery, paySecret string, timestamp int) *appError.Error {

	// TODO 需要完善对应的字段 验证签名 如果没有传参数 这里需要处理
	params := make(map[string]interface{})
	params["Timestamp"] = timestamp
	params["currency"] = d.Currency
	params["sign"] = strings.TrimSpace(s.Head.Signature)

	if !goutils.HmacSHA256Verify(params, paySecret) {
		return appError.UnauthenticatedErrCode // 签名错误
	}
	return nil
}

// GetMerchantAccountConfigInfo 获取 departs 信息
func (s *MerchantProjectServer) GetMerchantAccountConfigInfo() (*model.AspMerchantProjectConfig, *appError.Error) {
	var data model.AspMerchantProjectConfig

	appId := s.Head.AppId

	redisKey := fmt.Sprintf(constant.KEY_MERCHANT_PROJECT_CONFIG_INFO_STRING, appId)
	redisKey = constant.GetRedisKey(redisKey)

	merchantProjectConfigValue, err := s.Redis().Get(context.Background(), redisKey).Result()
	if err != constant.Nil {
		if len(merchantProjectConfigValue) > 0 {
			// JSON转结构体
			err = goutils.JsonDecode(merchantProjectConfigValue, &data)
			if err != nil {
				logger.ApiWarn(s.LogFileName, s.RequestId, "GetMerchantAccountConfigInfo: ", zap.Error(err))
				MissNotFoundErrCode := *appError.MissNotFoundErrCode
				return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissMerchantProjectConfigNotFoundErrMsg)
			}
			return &data, nil
		}
	}
	o := database.DB

	err = o.Where("mch_project_id = ?", appId).
		First(&data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissMerchantProjectConfigNotFoundErrMsg)
	}
	s.SetMerchantAccountConfigInfoToCache(&data, redisKey)
	return &data, nil
}

// SetMerchantAccountConfigInfoToCache 私有
// 设置Redis cp项目配置 信息
func (s *MerchantProjectServer) SetMerchantAccountConfigInfoToCache(Data *model.AspMerchantProjectConfig, redisKey string) bool {
	// 转成JSON
	merchantAccountJsonByte, err := json.Marshal(Data)
	if err != nil {
		return true
	}
	err = s.Redis().Set(context.Background(), redisKey, string(merchantAccountJsonByte), constant.EXPIRED_TIME).Err()
	return err == nil
}

// GetMerchantProjectCurrencyInfo 获取外部商户项目 的金额
func (s *MerchantProjectServer) GetMerchantProjectCurrencyInfo(currency string, MerchantProjectId int) (*model.AspMerchantProjectCurrency, *appError.Error) {
	var merchantProjectCurrency model.AspMerchantProjectCurrency

	//redisKey := fmt.Sprintf(constant.KEY_MERCHANT_PROJECT_CURRENCY_INFO_STRING, MerchantProjectId)
	//redisKey = constant.GetRedisKey(redisKey)
	//
	//merchantProjectCurrencyValue, err := s.Redis().Get(context.Background(), redisKey).Result()
	//if err != constant.Nil {
	//	if len(merchantProjectCurrencyValue) > 0 {
	//		// JSON转结构体
	//		err = goutils.JsonDecode(merchantProjectCurrencyValue, &merchantProjectCurrency)
	//		if err != nil {
	//			logger.ApiWarn(s.LogFileName, s.RequestId, "GetMerchantProjectCurrencyInfo: ", zap.Error(err))
	//			return nil, appError.NewError(constant.MissMerchantProjectCurrencyNotFoundErrMsg)
	//		}
	//		return &merchantProjectCurrency, nil
	//	}
	//}

	o := database.DB

	err := o.Where("status = ?", 1).
		Where("currency = ?", currency).
		Where("mch_project_id = ?", MerchantProjectId).
		First(&merchantProjectCurrency).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "GetMerchantProjectCurrencyInfo", zap.Error(err))
		return nil, appError.NewError(constant.MissMerchantProjectCurrencyNotFoundErrMsg)
	}
	//s.SetMerchantProjectCurrencyInfoToCache(&merchantProjectCurrency, redisKey)
	return &merchantProjectCurrency, nil
}

// SetMerchantProjectCurrencyInfoToCache 设置Redis cp项目配置 信息
func (s *MerchantProjectServer) SetMerchantProjectCurrencyInfoToCache(Data *model.AspMerchantProjectCurrency, redisKey string) bool {
	// 转成JSON
	merchantAccountJsonByte, err := json.Marshal(Data)
	if err != nil {
		return true
	}
	err = s.Redis().Set(context.Background(), redisKey, string(merchantAccountJsonByte), constant.EXPIRED_TIME).Err()
	return err == nil
}

// GetMerchantProjectCapitalFlowItem 查询cp项目每天收入、支付出的金额
// @cashType 资金方向类型，
// @mchProjectId cp项目id
// @beginTime 开始时间
// @endTime 结束时间
func (s *MerchantProjectServer) GetMerchantProjectCapitalFlowItem(businessType, mchProjectId int, beginTime, endTime int64) ([]*reqCommon.MerchantProjectCapitalFlowItem, error) {
	var merchantProjectCapitalFlow model.AspMerchantProjectCapitalFlow
	merchantProjectCapitalFlowItem := make([]*reqCommon.MerchantProjectCapitalFlowItem, 0)

	err := database.DB.Model(&merchantProjectCapitalFlow).Select("sum(total_fee) as total_fee").Where("mch_project_id = ?", mchProjectId).Where("business_type = ?", businessType).Where("create_time >= ?", beginTime).Where("create_time <= ?", endTime).Find(&merchantProjectCapitalFlowItem).Error
	if err != nil {
		return nil, err
	}
	return merchantProjectCapitalFlowItem, nil
}
