package service

import (
	"asp-payment/api-server/req"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/goRedis"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/pkg/validator/check"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"go.uber.org/zap"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/cast"
)

var checker = check.NewCheck()

type Service struct {
	C           *fiber.Ctx
	Head        *req.AspPaymentHeader
	RequestId   string
	LogFileName string
}

func NewService(c *fiber.Ctx, logFileName string) *Service {
	reqHeader := &req.AspPaymentHeader{}
	reqHeader.AppId = c.Get("AppId")
	timestamp := cast.ToInt(c.Get("Timestamp"))
	reqHeader.Timestamp = timestamp
	reqHeader.Version = c.Get("Version")
	reqHeader.Signature = c.Get("Signature")

	headerJson, _ := goutils.JsonEncode(*reqHeader)
	c.Locals("requestHeader", headerJson)

	requestId := cast.ToString(c.Locals(requestid.ConfigDefault.ContextKey))

	return &Service{C: c, RequestId: requestId, Head: reqHeader, LogFileName: logFileName}
}

func (s *Service) GetRequestId() string {
	return cast.ToString(s.C.Locals(requestid.ConfigDefault.ContextKey))
}

func (s *Service) Redis() *redis.Client {
	return goRedis.Redis
}

func (s *Service) Success(data interface{}) error {
	resMap := fiber.Map{
		"code":    appError.SUCCESS.Code,
		"message": appError.SUCCESS.Message,
		"data":    data,
	}
	responseBody, _ := goutils.JsonEncode(resMap)
	s.C.Locals("responseBody", responseBody)
	return s.C.JSON(resMap)
}

func (s *Service) SuccessJson(data interface{}) error {
	responseBody, _ := goutils.JsonEncode(data)
	s.C.Locals("responseBody", responseBody)
	return s.C.JSON(data)
}

//func (s *Service) RedirectSuccess(url string, body string) error {
//	s.C.Set("Location", url)
//	s.C.Locals("responseBody", body)
//	return s.C.Status(302).SendString(body)
//}

func (s *Service) RedirectSuccess(deepLink string) error {
	s.C.Locals("responseBody", deepLink)
	s.C.Set("Access-Control-Allow-Origin", "*")
	return s.C.Redirect(deepLink, 302)
}

func (s *Service) Error(appErr *appError.Error) error {
	resMap := fiber.Map{
		"code":    appErr.Code,
		"message": appErr.Message,
		"data":    nil,
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

func (s *Service) getKey(typeCode string) string {
	return fmt.Sprintf("api:black:set:%s", typeCode)
}

func (s *Service) getIpFilterKey(typeCode string, typeId string) string {
	return fmt.Sprintf("api:ipfilter:set:%s_%s", typeCode, typeId)
}

// VerifyApiBlack 验证黑名单
func (s *Service) VerifyApiBlack(phone string, device string, bankcode string, ip string) *appError.Error {

	//检查phone
	if phone != "" {
		exits := goRedis.Redis.SIsMember(context.Background(), s.getKey("phone"), phone).Val()
		if exits {

			logger.ApiWarn(s.LogFileName, s.RequestId, "黑名单用户拒绝请求,black_info：", zap.String("phone", phone))

			return appError.UserIsBlack
		}
	}

	if device != "" {
		exits := goRedis.Redis.SIsMember(context.Background(), s.getKey("device"), device).Val()
		if exits {
			logger.ApiWarn(s.LogFileName, s.RequestId, "黑名单用户拒绝请求,black_info：", zap.String("device", device))
			return appError.UserIsBlack
		}
	}

	if bankcode != "" {
		exits := goRedis.Redis.SIsMember(context.Background(), s.getKey("bankcode"), bankcode).Val()
		if exits {
			logger.ApiWarn(s.LogFileName, s.RequestId, "黑名单用户拒绝请求,bankcode：", zap.String("bankcode", bankcode))
			return appError.UserIsBlack
		}
	}

	if ip != "" {
		exits := goRedis.Redis.SIsMember(context.Background(), s.getKey("ip"), ip).Val()
		if exits {
			logger.ApiWarn(s.LogFileName, s.RequestId, "黑名单用户拒绝请求,ip：", zap.String("ip", ip))
			return appError.UserIsBlack
		}
	}

	return nil
}

func (s *Service) Lock(lockName string) bool {
	return goRedis.Lock(lockName)
}

func (s *Service) UnLock(lockName string) {
	goRedis.UnLock(lockName)
}
