package service

import (
	"asp-payment/api-server/req"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/go-redis/redis/v7"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DepartServer struct {
	*Service
}

func NewDepartServer(c *fiber.Ctx) *DepartServer {
	return &DepartServer{Service: NewService(c, constant.ChannelDepartTradeTypeLogFileName)}
}

func NewSimpleDepartServer() *DepartServer {
	return &DepartServer{&Service{LogFileName: constant.ChannelDepartTradeTypeLogFileName}}
}

// ChooseChannelDepartTradeType 收款 查询商户对应的渠道配置
func (s *DepartServer) ChooseChannelDepartTradeType(reqAspPayment *req.AspPayment) (*req.DeptTradeTypeInfo, *appError.Error) {
	// 贯穿支付所需要的参数
	appId := s.Head.AppId
	payment := reqAspPayment.PaymentMethod
	currency := reqAspPayment.OrderCurrency

	paymentArr := strings.Split(payment, ".")
	//provider := paymentArr[0]
	tradeType := strings.ToUpper(paymentArr[1])

	deptTradeType, err := s.ChooseChannelDepartTradeTypeInfoOfTwo(appId, tradeType, currency)

	if err != nil {
		return nil, err
	}

	return deptTradeType, nil

}

// ChoosePayoutChannelDepartTradeType  收款 查询商户对应的渠道配置
func (s *DepartServer) ChoosePayoutChannelDepartTradeType(reqAspPayout *req.AspPayout) (*req.DeptTradeTypeInfo, *appError.Error) {
	// 贯穿支付所需要的参数
	appId := s.Head.AppId
	//payment := reqAspPayout.Provider + "." + reqAspPayout.TradeType
	currency := reqAspPayout.OrderCurrency

	deptTradeType, err := s.ChooseChannelDepartTradeTypeInfoOfTwo(appId, reqAspPayout.TradeType, currency)

	if err != nil {
		return nil, err
	}

	return deptTradeType, nil

}

// ChooseMerchantProjectChannelDepartTradeType  查询cp对应的渠道配置
func (s *DepartServer) ChooseMerchantProjectChannelDepartTradeType(reqBody *req.AspMerchantProjectChannelReq) (*req.DeptTradeTypeInfo, *appError.Error) {
	// 贯穿支付所需要的参数
	appId := reqBody.MchProId
	//payment := reqAspPayout.Provider + "." + reqAspPayout.TradeType
	currency := reqBody.Currency

	deptTradeType, err := s.ChooseChannelDepartTradeTypeInfoOfTwo(appId, reqBody.TradeType, currency)

	if err != nil {
		return nil, err
	}

	return deptTradeType, nil

}

// ChooseChannelDepartTradeTypeInfo 选择的可以使用的渠道
func (s *DepartServer) ChooseChannelDepartTradeTypeInfo(appId, payment, currency string) (*req.DeptTradeTypeInfo, *appError.Error) {
	// var deptTradeTypeList req.DeptTradeTypeInfo
	// 查询cp项目关联  AspMerchantProjectChannelDepartTradeTypeLink 中的 departs
	departMerchantProjectLinkPluckList, err := s.GetDepartMerchantProjectLinkList(appId)
	if err != nil {
		return nil, err
	}

	deptIds, dmpErr := s.GetDepartMerchantProjectLinkIds(departMerchantProjectLinkPluckList)
	if dmpErr != nil {
		return nil, dmpErr
	}
	// channelConfigIds, err := NewChannelConfigServer(s.C).GetChannelConfigIds()
	// if err != nil {
	// 	return nil, err
	// }

	// 获取 cp 项目 可以使用的所有的内部商户列表 结构[[depart_id] => depart]
	departList, dpErr := s.GetDepartsList()
	if dpErr != nil {
		return nil, dpErr
	}
	newChannelConfigServer := NewChannelConfigServer(s.C)
	channelConfigList, ncErr := newChannelConfigServer.GetChannelConfigListList()
	if ncErr != nil {
		return nil, ncErr
	}

	channelConfigIds, nccErr := newChannelConfigServer.GetChannelConfigIds()
	if nccErr != nil {
		return nil, nccErr
	}

	// fmt.Println("channelConfigIds------------------", channelConfigIds)

	// 然后再次 根据 depart_ids 查询到 asp_channel_depart_tradetype 表中的数据 条件 depart_ids payment
	channelDepartTradeTypeList, dmtErr := s.GetDepartMerchantTradeTypeList(channelConfigIds, deptIds, payment)
	if dmtErr != nil {
		return nil, dmtErr
	}

	if len(*channelDepartTradeTypeList) == 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(channelDepartTradeTypeList) == 0: ", zap.String("err", constant.MissChannelConfigNotFoundErrMsg))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissChannelConfigNotFoundErrMsg)
	}

	// fmt.Println("channelDepartTradeTypeList-------------------", channelDepartTradeTypeList)

	// 然后查询出 所有的 channel_config 结构[[channel_id] => channel]

	// 做聚合排序 得到 唯一的一个渠道 payment 即返回对应的值

	channelDepartTradeTypeInfo, fttErr := s.FilterTradeTypeList(channelDepartTradeTypeList, channelConfigList, departList, departMerchantProjectLinkPluckList)

	if fttErr != nil {
		return nil, fttErr
	}

	// fmt.Println("channelDepartTradeTypeInfo-------------------", channelDepartTradeTypeInfo)
	if len(channelDepartTradeTypeInfo) == 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(channelDepartTradeTypeInfo) == 0: ", zap.String("err", constant.MissChannelConfigNotFoundErrMsg))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissChannelConfigNotFoundErrMsg)
	}
	return channelDepartTradeTypeInfo[0], nil
}

// ChooseDepartInfo 添加受益人 选择可以使用的商户
func (s *DepartServer) ChooseDepartInfo(appId string) (*req.AspMerchantProjectDepartList, *appError.Error) {
	// var deptTradeTypeList req.DeptTradeTypeList
	// 查询cp项目关联  AspDepartMerchantProjectLink 中的 departs
	departMerchantProjectLinkPluckList, err := s.GetDepartMerchantProjectLinkList(appId)
	if err != nil {
		return nil, err
	}
	// 获取 cp 项目 可以使用的所有的内部商户列表 结构[[depart_id] => depart]
	departList, dpErr := s.GetDepartsList()
	if dpErr != nil {
		return nil, dpErr
	}

	chooseDepartList, fttErr := s.FilterDepartList(departList, departMerchantProjectLinkPluckList)

	if fttErr != nil {
		return nil, fttErr
	}

	// fmt.Println("channelDepartTradeTypeInfo-------------------", channelDepartTradeTypeInfo)
	if len(chooseDepartList) == 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(channelDepartTradeTypeInfo) == 0: ", zap.String("err", constant.MissChannelConfigNotFoundErrMsg))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissChannelConfigNotFoundErrMsg)
	}
	return chooseDepartList[0], nil
}

// ChooseChannelDepartTradeTypeInfoOfTwo 选择cp支付渠道内部商户支付方式
// 修改了执行逻辑 asp_depart_merchant_project_link cp 关联的内部商户表 新增了渠道 支付方式字段
// 查询方式 做了调整 直接查询到对应的渠道 内部商户 支付方式 然后查询 asp_channel_depart_tradetype 支付方式 后 即可使用
func (s *DepartServer) ChooseChannelDepartTradeTypeInfoOfTwo(appId, tradeType, currency string) (*req.DeptTradeTypeInfo, *appError.Error) {
	// 根据支付方式 cp项目id 查询 cp可以使用的 渠道内部商户支付方式信息
	merchantProjectChannelDepartTradeTypeLinkInfo, err := s.GetMerchantProjectChannelDepartTradeTypeLinkInfo(appId, tradeType)
	if err != nil {
		return nil, appError.NewError(constant.MissMerchantProjectChannelDepartLinkNotFoundErrMsg)
	}
	if merchantProjectChannelDepartTradeTypeLinkInfo == nil {
		return nil, appError.MissNotFoundErrCode.FormatMessage(constant.MissMerchantProjectChannelDepartLinkNotFoundErrMsg)
	}
	// 根据渠道内部商户支付方式信息 channel_id depart_id trade_type 查询 渠道内部商户支付方式信息 渠道内部商户支付方式信息
	channelDepartTradeTypeInfo, cdtErr := s.GetChannelDepartTradeTypeInfo(merchantProjectChannelDepartTradeTypeLinkInfo.ChannelId, merchantProjectChannelDepartTradeTypeLinkInfo.DepartId, merchantProjectChannelDepartTradeTypeLinkInfo.TradeType)
	if cdtErr != nil {
		return nil, appError.MissNotFoundErrCode.FormatMessage(constant.MissMerchantProjectChannelDepartLinkNotFoundErrMsg)
	}
	if channelDepartTradeTypeInfo == nil {
		return nil, appError.MissNotFoundErrCode.FormatMessage(constant.MissMerchantProjectChannelDepartLinkNotFoundErrMsg)
	}

	deptTradeTypeInfo := new(req.DeptTradeTypeInfo)
	deptTradeTypeInfo.ID = channelDepartTradeTypeInfo.Id
	deptTradeTypeInfo.DepartID = cast.ToInt(channelDepartTradeTypeInfo.DepartId)
	deptTradeTypeInfo.ChannelID = cast.ToUint(channelDepartTradeTypeInfo.ChannelId)
	deptTradeTypeInfo.Provider = channelDepartTradeTypeInfo.Provider
	deptTradeTypeInfo.Payment = channelDepartTradeTypeInfo.Payment
	deptTradeTypeInfo.TradeType = channelDepartTradeTypeInfo.TradeType
	deptTradeTypeInfo.H5Type = channelDepartTradeTypeInfo.H5Type
	deptTradeTypeInfo.InFeeRate = channelDepartTradeTypeInfo.InFeeRate
	deptTradeTypeInfo.OutFeeRate = channelDepartTradeTypeInfo.OutFeeRate
	deptTradeTypeInfo.DayUpperLimit = channelDepartTradeTypeInfo.DayUpperLimit
	deptTradeTypeInfo.UpperLimit = channelDepartTradeTypeInfo.UpperLimit
	deptTradeTypeInfo.LowerLimit = channelDepartTradeTypeInfo.LowerLimit
	deptTradeTypeInfo.FixedAmount = channelDepartTradeTypeInfo.FixedAmount
	deptTradeTypeInfo.FixedCurrency = channelDepartTradeTypeInfo.FixedCurrency
	deptTradeTypeInfo.InFeeRateUpdating = channelDepartTradeTypeInfo.InFeeRateUpdating
	deptTradeTypeInfo.OutFeeRateUpdating = channelDepartTradeTypeInfo.OutFeeRateUpdating
	deptTradeTypeInfo.Sort = 0
	deptTradeTypeInfo.DepartSort = 0
	deptTradeTypeInfo.DepartMerchantProjectLinkSort = 0
	// 返回可以使用的渠道内部商户支付方式
	return deptTradeTypeInfo, nil
}

// GetDepartMerchantProjectLinkList 获取cp 项目 关联的depart_ids => 关联id信息（排序值，状态等信息）
func (s *DepartServer) GetDepartMerchantProjectLinkList(merchantProjectId string) (map[int]*req.DepartMerchantProjectLinkList, *appError.Error) {
	departMerchantProjectLinkPluckList := make(map[int]*req.DepartMerchantProjectLinkList, 0)

	redisKey := fmt.Sprintf(constant.KEY_DEPART_MERCHANT_PROJECT_LINK_LIST_STRING, merchantProjectId)
	redisKey = constant.GetRedisKey(redisKey)
	// fmt.Println("redisKey------------------", redisKey)
	merchantProjectLinkListValue := s.Redis().Get(context.Background(), redisKey).Val()
	//fmt.Println("merchantProjectLinkListValue: ", merchantProjectLinkListValue)
	//fmt.Println("err != redis.Nil: ", err != redis.Nil)
	if merchantProjectLinkListValue != "" {
		err := goutils.JsonDecode(merchantProjectLinkListValue, &departMerchantProjectLinkPluckList)
		// fmt.Println("deptIds------------------", deptIds)
		if err != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "merchantProjectLinkListValue: ", zap.Error(err))
			MissNotFoundErrCode := *appError.MissNotFoundErrCode
			return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissMerchantProjectChannelDepartLinkNotFoundErrMsg)
		}
		return departMerchantProjectLinkPluckList, nil
	}
	o := database.DB
	departList := make([]*req.DepartMerchantProjectLinkList, 0)
	var data model.AspMerchantProjectChannelDepartTradeTypeLink
	err := o.Model(&data).Where("status = ?", 1).Where("mch_project_id = ? ", merchantProjectId).Find(&departList).Error
	//fmt.Println("departList: ", departList)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspMerchantProjectChannelDepartTradeTypeLink: ", zap.Error(err))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissMerchantProjectChannelDepartLinkNotFoundErrMsg)
	}
	if len(departList) == 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(departList) == 0: ", zap.String("err", constant.MissMerchantProjectChannelDepartLinkNotFoundErrMsg))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissMerchantProjectChannelDepartLinkNotFoundErrMsg)
	}
	for i := 0; i < len(departList); i++ {
		departMerchantProjectLinkPluckList[departList[i].DepartId] = departList[i]
	}
	s.SetMerchantProjectLinkListToCache(departMerchantProjectLinkPluckList, redisKey)
	return departMerchantProjectLinkPluckList, nil
}

// GetDepartMerchantProjectLinkIds 获取cp 项目 关联的depart_ids 数组 转换
// departMerchantProjectLinkPluckList 查询出来的内部商户depart_id =》 商户项目内部商户信息 关联数据结构
func (s *DepartServer) GetDepartMerchantProjectLinkIds(departMerchantProjectLinkPluckList map[int]*req.DepartMerchantProjectLinkList) ([]int, *appError.Error) {

	departIds := make([]int, 0)

	for departId, _ := range departMerchantProjectLinkPluckList {
		departIds = append(departIds, departId)
	}

	return departIds, nil
}

// GetDepartMerchantTradeTypeList 获取商户渠道列表
func (s *DepartServer) GetDepartMerchantTradeTypeList(channelConfigIds []int, merchantProjectDepartIds []int, payment string) (*[]model.AspChannelDepartTradeType, *appError.Error) {
	deptTradeTypeList := make([]model.AspChannelDepartTradeType, 0)

	paymentArr := strings.Split(payment, ".")
	provider := paymentArr[0]
	tradeType := strings.ToUpper(paymentArr[1])

	// 新增对渠道的ids 的绑定
	//redisKey := fmt.Sprintf(constant.KEY_CHANNEL_DEPART_PROVIDER_TRADETYPE_LIST, appId, currency, provider, trade_type)
	//redisKey = constant.GetRedisKey(redisKey)
	//
	//departMerchantTradeTypeListValue, err := s.Redis().Get(context.Background(), redisKey).Result()
	//if err != redis.Nil {
	//	if len(departMerchantTradeTypeListValue) > 0 {
	//		// JSON转结构体
	//		err = goutils.JsonDecode(departMerchantTradeTypeListValue, &deptTradeTypeList)
	//		if err != nil {
	//			logger.ApiWarn(s.LogFileName, s.RequestId, "departMerchantTradeTypeListValue: ", zap.Error(err))
	//			return nil, err
	//		}
	//		return deptTradeTypeList, nil
	//	}
	//}

	o := database.DB

	var data model.AspChannelDepartTradeType
	err := o.
		Model(&data).
		Where("channel_id in ? ", channelConfigIds).
		Where("depart_id in ? ", merchantProjectDepartIds).
		Where("provider = ? ", provider).
		Where("trade_type = ? ", tradeType).
		//Where("disabled = ?", 0).         // 0 上级控制未禁止
		//Where("selected = ?", 1).         // 1 自己选择已启用
		//Where("channel_disabled = ?", 0). // 0 渠道未禁止
		//Where("is_config = ?", 1).        // 1 已配置
		Where("status = ?", 1). // 1 成功
		Find(&deptTradeTypeList).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartTradeType: ", zap.Error(err))
		return nil, appError.NewError(err.Error())
	}
	// fmt.Println(deptTradeTypeList, "deptTradeTypeList-------------------")
	// for i := 0; i < len(deptList); i++ {
	// 	deptIds = append(deptIds, deptList[i].DeptId)
	// }
	if len(deptTradeTypeList) == 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(deptTradeTypeList) == 0 ", zap.String("err", constant.MissChannelDepartTradeTypeNotFoundErrMsg))
		return nil, appError.ChannelDepartTradeTypeNotFoundErrCode
	}
	//s.SetDepartMerchantTradeTypeListToCache(deptTradeTypeList, redisKey)
	return &deptTradeTypeList, nil
}

// SetDepartMerchantTradeTypeListToCache 设置 内部商户&cp项目关联表 setMerchantProjectLinkList
func (s *DepartServer) SetDepartMerchantTradeTypeListToCache(Data []*model.AspChannelDepartTradeType, redisKey string) bool {

	// 转成JSON
	idJsonByte, err := json.Marshal(Data)
	if err != nil {
		return true
	}
	err = s.Redis().Set(context.Background(), redisKey, string(idJsonByte), constant.EXPIRED_TIME).Err()
	return err == nil
}

// SetMerchantProjectLinkListToCache 私有
// 设置Redis
func (s *DepartServer) SetMerchantProjectLinkListToCache(Data map[int]*req.DepartMerchantProjectLinkList, redisKey string) bool {
	// 转成JSON
	idJsonByte, err := json.Marshal(Data)
	if err != nil {
		return true
	}
	err = s.Redis().Set(context.Background(), redisKey, string(idJsonByte), constant.EXPIRED_TIME).Err()
	return err == nil
}

// GetDepartsList 获取所有的 cp 项目 可以使用的所有内部商户列表
// depart_ids 获取的
func (s *DepartServer) GetDepartsList() (map[int]*req.AspDepartList, *appError.Error) {
	departPluckList := make(map[int]*req.AspDepartList, 0)

	// redisKey := fmt.Sprintf(redis.KEY_DEPART_LIST_STRING)
	redisKey := constant.GetRedisKey(constant.KEY_DEPART_LIST_STRING)

	channelConfigValue := s.Redis().Get(context.Background(), redisKey).Val()
	if channelConfigValue != "" {
		err := goutils.JsonDecode(channelConfigValue, &departPluckList)
		if err != nil {
			logger.ApiWarn(s.LogFileName, s.RequestId, "channelConfigValue: ", zap.Error(err))
			return nil, appError.NewError(err.Error())
		}
		return departPluckList, nil
	}

	departList := make([]*req.AspDepartList, 0)
	o := database.DB
	var data model.AspDeparts
	err := o.Model(&data).
		Where("status = ?", 1).
		Where("depart_type = ?", 2).
		Find(&departList).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspDeparts: ", zap.Error(err))
		return nil, appError.MissDepartListNotFoundErrMsg
	}
	if len(departList) == 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(departList) == 0: ", zap.String("err", constant.MissDepartListNotFoundErrMsg))
		return nil, appError.MissDepartListNotFoundErrMsg
	}

	for i := 0; i < len(departList); i++ {
		departPluckList[departList[i].Id] = departList[i]
	}

	s.SetDepartListToCache(departPluckList, redisKey)

	return departPluckList, nil
}

// GetDepartInfo 获取 departs 信息
func (s *DepartServer) GetDepartInfo(departId int) (*model.AspDeparts, error) {
	var departInfo model.AspDeparts

	redisKey := fmt.Sprintf(constant.KEY_DEPART_INFO_STRING, departId)
	redisKey = constant.GetRedisKey(redisKey)

	idValue, err := s.Redis().Get(context.Background(), redisKey).Result()
	if err != redis.Nil {
		if len(idValue) > 0 {
			// JSON转结构体
			err = goutils.JsonDecode(idValue, &departInfo)
			if err != nil {
				logger.ApiInfo(s.LogFileName, s.RequestId, "GetMerchantAccountInfo ", zap.Error(err))
				return nil, err
			}
			return &departInfo, nil
		}
	}

	o := database.DB
	err = o.Where("status = ?", 1).
		Where("id = ?", departId).
		First(&departInfo).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ApiInfo(s.LogFileName, s.RequestId, "AspDeparts ", zap.Error(err))
		return &departInfo, err
	}
	s.SetDepartInfoToCache(&departInfo, redisKey)
	return &departInfo, nil
}

// SetDepartListToCache 设置 channel_config 数据的 redis 的缓存
func (s *DepartServer) SetDepartListToCache(Data map[int]*req.AspDepartList, redisKey string) bool {
	// 转成JSON
	channelConfigJsonByte, err := json.Marshal(Data)
	if err != nil {
		return true
	}
	err = s.Redis().Set(context.Background(), redisKey, string(channelConfigJsonByte), constant.EXPIRED_TIME).Err()
	return err == nil
}

// SetDepartInfoToCache 私有
// 设置Redis 商户详情 信息
func (s *DepartServer) SetDepartInfoToCache(Data *model.AspDeparts, redisKey string) bool {
	// 转成JSON
	departJsonByte, err := json.Marshal(Data)
	if err != nil {
		return true
	}
	err = s.Redis().Set(context.Background(), redisKey, string(departJsonByte), constant.EXPIRED_TIME).Err()
	return err == nil
}

// FilterTradeTypeList 获取cp 项目 关联的depart_ids DeptTradeTypeInfo
// deptTradeTypeList 待排序中的cp 项目 能使用的所有的渠道
// channelConfigList 可以使用的cp 项目 渠道 使用到的渠道的sort值
// departsList 可以使用的cp 项目 渠道 使用到的代理的sort值
func (s *DepartServer) FilterTradeTypeList(deptTradeTypeList *[]model.AspChannelDepartTradeType, channelConfigList map[int]*req.ChannelConfigList, departList map[int]*req.AspDepartList, departMerchantProjectLinkPluckList map[int]*req.DepartMerchantProjectLinkList) ([]*req.DeptTradeTypeInfo, *appError.Error) {
	// 这里需要和渠道的排序 目前已有的 渠道中的排序  商户中的渠道的 tradetype (tradetype 中包含有channel的id)

	reqDeptTradeTypeList := make([]*req.DeptTradeTypeInfo, 0, len(*deptTradeTypeList))

	//fmt.Println("deptTradeTypeList[0]:", deptTradeTypeList[0])
	//fmt.Println("deptTradeTypeList[1]:", deptTradeTypeList[1])
	//fmt.Println("deptTradeTypeList[2]:", deptTradeTypeList[2])
	//fmt.Println("deptTradeTypeList[3]:", deptTradeTypeList[3])

	var ok bool

	for i := 0; i < len(*deptTradeTypeList); i++ {
		departId := cast.ToInt((*deptTradeTypeList)[i].DepartId)
		// 当商户和内部商户关联表 asp_depart_merchant_project_link 有值
		// depart 给关闭了 这种情况 则跳过
		if _, ok = departList[departId]; !ok {
			continue
		}

		departTradeTypeInfo := &req.DeptTradeTypeInfo{}
		departTradeTypeInfo.ID = (*deptTradeTypeList)[i].Id
		departTradeTypeInfo.ChannelID = cast.ToUint((*deptTradeTypeList)[i].ChannelId)
		departTradeTypeInfo.DepartID = cast.ToInt((*deptTradeTypeList)[i].DepartId)
		departTradeTypeInfo.Provider = (*deptTradeTypeList)[i].Provider
		departTradeTypeInfo.Payment = (*deptTradeTypeList)[i].Payment
		departTradeTypeInfo.TradeType = (*deptTradeTypeList)[i].TradeType
		departTradeTypeInfo.InFeeRate = (*deptTradeTypeList)[i].InFeeRate
		departTradeTypeInfo.OutFeeRate = (*deptTradeTypeList)[i].OutFeeRate
		departTradeTypeInfo.DayUpperLimit = (*deptTradeTypeList)[i].DayUpperLimit
		departTradeTypeInfo.UpperLimit = (*deptTradeTypeList)[i].UpperLimit
		departTradeTypeInfo.LowerLimit = (*deptTradeTypeList)[i].LowerLimit
		departTradeTypeInfo.FixedAmount = (*deptTradeTypeList)[i].FixedAmount
		departTradeTypeInfo.FixedCurrency = (*deptTradeTypeList)[i].FixedCurrency
		departTradeTypeInfo.InFeeRateUpdating = (*deptTradeTypeList)[i].InFeeRateUpdating
		departTradeTypeInfo.OutFeeRateUpdating = (*deptTradeTypeList)[i].OutFeeRateUpdating
		departTradeTypeInfo.Sort = channelConfigList[(*deptTradeTypeList)[i].ChannelId].Sort // 为了新增到了排序的值 来渠道排序

		departTradeTypeInfo.DepartSort = departList[departId].Sort // 为了新增到了排序的值 来内部商户排序

		if _, ok = departMerchantProjectLinkPluckList[departId]; ok {
			departTradeTypeInfo.DepartMerchantProjectLinkSort = departMerchantProjectLinkPluckList[departId].Sort
		} else {
			departTradeTypeInfo.DepartMerchantProjectLinkSort = 1
		}

		reqDeptTradeTypeList = append(reqDeptTradeTypeList, departTradeTypeInfo)
	}

	// fmt.Println("reqDeptTradeTypeList-----------------", reqDeptTradeTypeList)

	sort.Slice(reqDeptTradeTypeList, func(i, j int) bool {
		if reqDeptTradeTypeList[i].Sort != reqDeptTradeTypeList[j].Sort {
			return reqDeptTradeTypeList[i].Sort > reqDeptTradeTypeList[j].Sort
		}
		return reqDeptTradeTypeList[i].ChannelID > reqDeptTradeTypeList[j].ChannelID
	})
	sort.Slice(reqDeptTradeTypeList, func(i, j int) bool {
		if reqDeptTradeTypeList[i].DepartSort != reqDeptTradeTypeList[j].DepartSort {
			return reqDeptTradeTypeList[i].DepartSort > reqDeptTradeTypeList[j].DepartSort
		}
		return reqDeptTradeTypeList[i].DepartID > reqDeptTradeTypeList[j].DepartID
	})

	sort.Slice(reqDeptTradeTypeList, func(i, j int) bool {
		return reqDeptTradeTypeList[i].DepartMerchantProjectLinkSort > reqDeptTradeTypeList[j].DepartMerchantProjectLinkSort
	})

	// 渠道配置 需要转换成 id object 的结构
	return reqDeptTradeTypeList, nil
}

// FilterDepartList 获取cp 项目 关联的depart
// departsList 可以使用的cp 项目 渠道 使用到的代理的sort值
// departMerchantProjectLinkPluckList 可以使用的cp 项目 渠道 使用到的代理的sort值
func (s *DepartServer) FilterDepartList(departList map[int]*req.AspDepartList, departMerchantProjectLinkPluckList map[int]*req.DepartMerchantProjectLinkList) ([]*req.AspMerchantProjectDepartList, *appError.Error) {
	// 这里需要和渠道的排序 目前已有的 渠道中的排序  商户中的渠道的 tradetype (tradetype 中包含有channel的id)

	reqAspMerchantProjectDepartList := make([]*req.AspMerchantProjectDepartList, 0)

	var ok bool

	for _, depart := range departList {
		departId := depart.Id
		// 当商户和内部商户关联表 asp_depart_merchant_project_link 有值
		// depart 给关闭了 这种情况 则跳过
		if _, ok = departMerchantProjectLinkPluckList[departId]; !ok {
			continue
		}

		departInfo := new(req.AspMerchantProjectDepartList)
		departInfo.Id = depart.Id
		departInfo.ParentId = depart.ParentId
		departInfo.DepartType = depart.DepartType
		departInfo.Title = depart.Title
		departInfo.CurrencyId = depart.CurrencyId
		departInfo.Sort = depart.Sort
		departInfo.MchProjectSort = departMerchantProjectLinkPluckList[departId].Sort
		departInfo.ContactsName = depart.ContactsName
		departInfo.ContactsPhone = depart.ContactsPhone
		departInfo.ContactsEmail = depart.ContactsEmail
		departInfo.BankName = depart.BankName
		departInfo.BankOfDeposit = depart.BankOfDeposit
		departInfo.BankAccount = depart.BankAccount
		departInfo.BankAccountName = depart.BankAccountName
		departInfo.CompanyName = depart.CompanyName
		departInfo.CompanyShort = depart.CompanyShort
		departInfo.CompanyEmail = depart.CompanyEmail
		departInfo.Status = depart.Status
		departInfo.CreateTime = depart.CreateTime
		departInfo.UpdateTime = depart.UpdateTime
		departInfo.Remark = depart.Remark
		departInfo.DefaultChannelId = depart.DefaultChannelId
		departInfo.SettlementType = depart.SettlementType
		reqAspMerchantProjectDepartList = append(reqAspMerchantProjectDepartList, departInfo)
	}
	//for i := 0; i < len(departList); i++ {
	//
	//}

	sort.Slice(reqAspMerchantProjectDepartList, func(i, j int) bool {
		return reqAspMerchantProjectDepartList[i].MchProjectSort > reqAspMerchantProjectDepartList[j].MchProjectSort
	})

	// 渠道配置 需要转换成 id object 的结构
	return reqAspMerchantProjectDepartList, nil
}

// GetMerchantProjectChannelDepartTradeTypeLinkInfo 获取cp渠道内部商户支付方式信息
func (s *DepartServer) GetMerchantProjectChannelDepartTradeTypeLinkInfo(appId, tradeType string) (*model.AspMerchantProjectChannelDepartTradeTypeLink, error) {
	// 先去查询缓存
	key := constant.GetRedisKey(fmt.Sprintf(constant.KEY_MERCHANT_PROJECT_CHANNEL_DEPART_TRADE_TYPE_INFO_STRING, appId, tradeType))
	return repository.AspMerchantProjectChannelDepartTradeTypeLinkRepository.GetCacheInfo(key, func() (*model.AspMerchantProjectChannelDepartTradeTypeLink, error) {
		return repository.AspMerchantProjectChannelDepartTradeTypeLinkRepository.FindOne(database.NewSqlCondition().Where("mch_project_id = ?", appId).Where("channel_status = ?", 1).Where("trade_status = ?", 1).Where("status = ?", 1).Where("trade_type = ?", tradeType).Desc("sort"))
	})
}

// GetChannelDepartTradeTypeInfo 获取渠道内部商户支付方式
func (s *DepartServer) GetChannelDepartTradeTypeInfo(ChannelId, DepartId int, tradeType string) (*model.AspChannelDepartTradeType, error) {
	// 先去查询缓存
	key := constant.GetRedisKey(fmt.Sprintf(constant.KEY_CHANNEL_DEPART_TRADE_TYPE_INFO_STRING, ChannelId, DepartId, tradeType))
	return repository.AspChannelDepartTradeTypeRepository.GetCacheInfo(key, func() (*model.AspChannelDepartTradeType, error) {
		return repository.AspChannelDepartTradeTypeRepository.FindOne(database.NewSqlCondition().Where("channel_id = ?", ChannelId).Where("depart_id = ?", DepartId).Where("trade_type = ?", tradeType).Where("status = ?", 1))
	})
}
