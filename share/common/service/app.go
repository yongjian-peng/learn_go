package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"share/common/model"
	"share/common/pkg/goRedis"
	"share/common/repository"
	"sync"
)

var appService *AppService
var AppServiceOnce sync.Once

type AppService struct {
	*Service
}

func App(c *fiber.Ctx) *AppService {
	return &AppService{Service: NewService(c, "appService")}
}

func (s *AppService) GetCacheAppInfo(appId int) (*model.App, error) {
	key := goRedis.GetKey(fmt.Sprintf("appinfo:%d", appId))
	return repository.AppRepository.GetCacheInfo(key, func() (*model.App, error) {
		return repository.AppRepository.Take("id = ?", appId)
	})
}
