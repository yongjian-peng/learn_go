package service

import (
	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	*Service
}

func User(c *fiber.Ctx) *UserService {
	return &UserService{Service: NewService(c, "user")}
}

// UpdateUserConfig 修改配置
func (s *UserService) UpdateUserConfig() error {

	//11 22
	return s.Success(nil)
}
