package service

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-module/carbon/v2"
	"share/common/pkg/appError"
	"share/common/pkg/config"
	"share/common/pkg/goRedis"
	"share/common/pkg/goutils"
	"share/common/pkg/validator/check"
	"share/common/req"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/cast"
)

var checker = check.NewCheck()
var ctx = context.Background()

type Service struct {
	C           *fiber.Ctx
	Head        *req.Header
	RequestId   string
	LogFileName string
	Uid         int64
}

func NewService(c *fiber.Ctx, logFileName string) *Service {
	reqHeader := &req.Header{}
	reqHeader.AppId = c.Get("AppId")
	reqHeader.Version = c.Get("Version")
	reqHeader.Signature = c.Get("Signature")

	headerJson, _ := goutils.JsonEncode(*reqHeader)
	c.Locals("requestHeader", headerJson)
	uid := c.Locals("UID")
	requestId := cast.ToString(c.Locals(requestid.ConfigDefault.ContextKey))

	return &Service{C: c, RequestId: requestId, Head: reqHeader, LogFileName: logFileName, Uid: cast.ToInt64(uid)}
}

func (s *Service) GetRequestId() string {
	return cast.ToString(s.C.Locals(requestid.ConfigDefault.ContextKey))
}

func (s *Service) Redis() *redis.Client {
	return goRedis.Redis
}

func (s *Service) Success(data interface{}) error {
	resMap := fiber.Map{
		"code": appError.SUCCESS.Code,
		"msg":  appError.SUCCESS.Msg,
		"data": data,
	}
	//null转{}输出
	if data == nil {
		resMap["data"] = fiber.Map{}
	}
	responseBody, _ := goutils.JsonEncode(resMap)
	s.C.Locals("responseBody", responseBody)
	return s.C.JSON(resMap)
}

func (s *Service) Error(appErr *appError.Error) error {
	resMap := fiber.Map{
		"code": appErr.Code,
		"msg":  appErr.Msg,
		"data": fiber.Map{},
	}
	responseBody, _ := goutils.JsonEncode(resMap)
	s.C.Locals("responseBody", responseBody)
	return s.C.JSON(resMap)
}

func (s *Service) SuccessToCode(str string) error {
	s.C.Locals("responseBody", str)
	return s.C.SendString(str)
}

func (s *Service) QueryArgsAll() (map[string]string, error) {
	data := make(map[string]string)
	var err error
	s.C.Context().QueryArgs().VisitAll(func(key, val []byte) {
		if err != nil {
			return
		}
		k := utils.UnsafeString(key)
		v := utils.UnsafeString(val)

		if strings.Contains(k, "[") {
			k, err = goutils.ParseParamSquareBrackets(k)
		}
		data[k] = v
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) PostArgsAll() (map[string]string, error) {
	data := make(map[string]string)
	var err error
	s.C.Context().PostArgs().VisitAll(func(key, val []byte) {
		k := utils.UnsafeString(key)
		v := utils.UnsafeString(val)

		if strings.Contains(k, "[") {
			k, err = goutils.ParseParamSquareBrackets(k)
		}
		data[k] = v
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) ReqArgsAll() (map[string]string, error) {
	data := make(map[string]string)
	query, _ := s.QueryArgsAll()
	post, _ := s.PostArgsAll()
	for i, v := range query {
		for j, w := range post {
			if i == j {
				data[i] = w
			} else {
				if _, ok := data[i]; !ok {
					data[i] = v
				}
				if _, ok := data[j]; !ok {
					data[j] = w
				}
			}
		}
	}
	return data, nil
}

// ValidatePost 绑定并校验数据
func (s *Service) ValidatePost(out interface{}) *appError.Error {
	//绑定参数
	if err := s.C.BodyParser(out); err != nil {
		return appError.ParameterTypeError
	}
	//校验参数
	if err := checker.Struct(out); err != nil {
		return err
	}
	return nil
}

// ValidateGet 绑定并校验数据
func (s *Service) ValidateGet(out interface{}) *appError.Error {
	//绑定参数
	if err := s.C.QueryParser(out); err != nil {
		fmt.Println("parser err:", err.Error())
		return appError.ParameterTypeError
	}
	//校验参数
	if err := checker.Struct(out); err != nil {
		return err
	}
	return nil
}

func (s *Service) Lock(lockName string) bool {
	return goRedis.Lock(lockName)
}

func (s *Service) UnLock(lockName string) {
	goRedis.UnLock(lockName)
}

func (s *Service) getSequenceKey(currentTime int64) string {
	return fmt.Sprintf("%s:%s:snow:seq:order:%d", config.AppConfig.Server.Name, config.AppConfig.Server.Env, currentTime)
}

func (s *Service) sequence(key string) int64 {

	autoIncrement := goRedis.Redis.Incr(context.Background(), key).Val()
	// 设置过期时间
	if autoIncrement <= 1 {
		goRedis.Redis.Expire(context.Background(), key, 3*time.Second)
	}

	return autoIncrement
}

// GetSequenceId 自增id 19位
func (s *Service) GetSequenceId() string {
	maxOrderSequenceLength := 13 //每微秒可生成 8191 个
	//通过redis获取单位微秒内的自增id
	microsecond := carbon.Now().TimestampMicro()
	key := s.getSequenceKey(microsecond)
	sequence := s.sequence(key)
	for sequence > (-1 ^ (-1 << maxOrderSequenceLength)) {
		//暂停1微秒数
		time.Sleep(1 * time.Microsecond)
		microsecond = carbon.Now().TimestampMicro()
		key = s.getSequenceKey(microsecond)
		sequence = s.sequence(key)
	}
	strMicrosecond := cast.ToString(microsecond)
	return fmt.Sprintf("%s%s%s", carbon.CreateFromTimestampMicro(microsecond).Format("ymdHis"), strMicrosecond[10:13], fmt.Sprintf("%04d", sequence))
}

func (s *Service) GetAppService() *AppService {
	AppServiceOnce.Do(func() {
		appService = App(s.C)
	})
	return appService
}
