package service

import (
	"asp-payment/api-server/req"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChannelConfigServer struct {
	*Service
}

func NewChannelConfigServer(c *fiber.Ctx) *ChannelConfigServer {
	return &ChannelConfigServer{Service: NewService(c, constant.ChannelConfigLogFileName)}
}

func NewSimpleChannelConfigServer() *ChannelConfigServer {
	return &ChannelConfigServer{&Service{LogFileName: constant.ChannelConfigLogFileName}}
}

// GetChannelConfigInfo 获取渠道信息
func (s *ChannelConfigServer) GetChannelConfigInfo(channelId int) (*model.AspChannelConfig, *appError.Error) {
	var channelConfigInfo model.AspChannelConfig
	// 获取渠道的列表数据 格式 [channel_id]channel_config
	channelConfigList, err := s.GetChannelConfigListList()
	if err != nil {
		return nil, err
	}
	if channelConfig, ok := channelConfigList[channelId]; ok {
		// 转换成 model
		channelConfig.Generate(&channelConfigInfo)
		return &channelConfigInfo, nil
	}

	o := database.DB
	dbErr := o.Where("status = ?", 1).Where("id = ?", channelId).First(&channelConfigInfo).Error
	// TODO 需要 统一记录日志 待处理
	if dbErr != nil && errors.Is(dbErr, gorm.ErrRecordNotFound) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelConfig ", zap.Error(err))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissChannelConfigNotFoundErrMsg) // 渠道信息不存
	}
	return &channelConfigInfo, nil
}

// GetChannelConfigInfo 获取渠道信息 不带缓存
func (s *ChannelConfigServer) GetChannelConfigInfo2(params map[string]interface{}) (*model.AspChannelConfig, *appError.Error) {
	var channelConfigInfo model.AspChannelConfig
	// 获取渠道的列表数据 格式 [channel_id]channel_config
	o := database.DB
	for k, v := range params {
		switch k {
		case "id":
			o = o.Where("id = ?", v)
		case "name":
			o = o.Where("name = ?", v)
		}
	}

	dbErr := o.Where("status = ?", 1).First(&channelConfigInfo).Error
	// TODO 需要 统一记录日志 待处理
	if dbErr != nil && errors.Is(dbErr, gorm.ErrRecordNotFound) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelConfig ", zap.Error(dbErr))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissChannelConfigNotFoundErrMsg) // 渠道信息不存
	}
	return &channelConfigInfo, nil
}

func (s *ChannelConfigServer) GetChannelConfigIds() (channelConfigIds []int, err *appError.Error) {
	// 获取渠道的列表数据 格式 [channel_id]channel_config
	channelConfigList, err := s.GetChannelConfigListList()
	if err != nil {
		return channelConfigIds, err
	}
	if len(channelConfigList) > 0 {
		// map 说两个 Go 循环里的坑 参考 https://mp.weixin.qq.com/s/QtFkh5d7Y-n2i4JI6tUaNA
		// 基本类型不使用指针，引用类型才使用指针，这是常规做法
		for channelId, _ := range channelConfigList {
			channelConfigIds = append(channelConfigIds, channelId)
		}
	} else {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(channelConfigList) > 0  ", zap.String("err", constant.MissChannelConfigNotFoundErrMsg))

		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissChannelConfigNotFoundErrMsg)
	}
	return
}

// GetChannelConfigListList 获取 渠道的列表
func (s *ChannelConfigServer) GetChannelConfigListList() (map[int]*req.ChannelConfigList, *appError.Error) {

	// deptIds := make([]int, 0)

	channelConfigPluckList := make(map[int]*req.ChannelConfigList, 0)

	// redisKey := fmt.Sprintf(redis.KEY_CHANNEL_LIST_STRING)
	redisKey := constant.GetRedisKey(constant.KEY_CHANNEL_LIST_STRING)

	channelConfigValue, err := s.Redis().Get(context.Background(), redisKey).Result()
	if err != constant.Nil {
		if len(channelConfigValue) > 0 {
			// JSON转结构体
			err = goutils.JsonDecode(channelConfigValue, &channelConfigPluckList)
			if err != nil {
				logger.ApiWarn(s.LogFileName, s.RequestId, "channelConfigValue.json ", zap.Error(err))
				return nil, appError.NewError(err.Error())
			}
			return channelConfigPluckList, nil
		}
	}

	channelConfigList := make([]req.ChannelConfigList, 0)
	o := database.DB
	var data model.AspChannelConfig
	err = o.Model(&data).
		Where("status = ?", 1).
		Find(&channelConfigList).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelConfig ", zap.Error(err))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissChannelConfigNotFoundErrMsg)
	}
	if len(channelConfigList) == 0 {
		logger.ApiWarn(s.LogFileName, s.RequestId, "len(channelConfigList) == 0 ", zap.String("err", constant.MissChannelConfigNotFoundErrMsg))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		return nil, (&MissNotFoundErrCode).FormatMessage(constant.MissChannelConfigNotFoundErrMsg)
	}

	for i := 0; i < len(channelConfigList); i++ {
		channelConfigPluckList[channelConfigList[i].Id] = &channelConfigList[i]
	}

	s.SetChannelConfigListToCache(channelConfigPluckList, redisKey)

	return channelConfigPluckList, nil
}

// SetChannelConfigListToCache 设置 channel_config 数据的 redis 的缓存
func (s *ChannelConfigServer) SetChannelConfigListToCache(Data map[int]*req.ChannelConfigList, redisKey string) bool {

	// 转成JSON
	channelConfigJsonByte, err := json.Marshal(Data)
	if err != nil {
		return true
	}
	err = s.Redis().Set(context.Background(), redisKey, string(channelConfigJsonByte), constant.EXPIRED_TIME).Err()
	return err == nil
}

// GetChannelDepartInfo 获取商户的渠道 中的配置
func (s *ChannelConfigServer) GetChannelDepartInfo(params map[string]string) (*model.AspChannelDepartConfig, *appError.Error) {
	o := database.DB
	var data model.AspChannelDepartConfig

	redisKey := fmt.Sprintf(constant.KEY_CHANNEL_DEPART_INFO_STRING, params["channel_id"], params["depart_id"])
	redisKey = constant.GetRedisKey(redisKey)

	channelDepartValue, err := s.Redis().Get(context.Background(), redisKey).Result()
	if err != constant.Nil {
		if len(channelDepartValue) > 0 {
			// JSON转结构体
			err = goutils.JsonDecode(channelDepartValue, &data)
			if err != nil {
				logger.ApiWarn(s.LogFileName, s.RequestId, "channelDepartValue ", zap.Error(err))
				return nil, appError.MissNotFoundErrCode.FormatMessage(constant.MissChannelDepartNotFoundErrMsg)
			}
			return &data, nil
		}
	}
	err = o.Where("depart_id = ?", params["depart_id"]).
		Where("channel_id = ?", params["channel_id"]).
		Where("status = ?", 1).
		First(&data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ApiWarn(s.LogFileName, s.RequestId, "AspChannelDepartConfig ", zap.Error(err))
		return nil, appError.NewError(err.Error())
	}
	s.SetChannelDepartInfoToCache(&data, redisKey)
	return &data, nil
}

// SetChannelDepartInfoToCache 私有
// 设置Redis key信息
func (s *ChannelConfigServer) SetChannelDepartInfoToCache(Data *model.AspChannelDepartConfig, redisKey string) bool {
	// 转成JSON
	idJsonByte, err := json.Marshal(Data)
	if err != nil {
		return true
	}
	err = s.Redis().Set(context.Background(), redisKey, string(idJsonByte), constant.EXPIRED_TIME).Err()
	return err == nil
}
