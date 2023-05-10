package middleware

import (
	"asp-payment/common/pkg/config"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ServerPort New creates a new middleware handler
func ServerPort() fiber.Handler {
	// Return new handler
	return func(c *fiber.Ctx) error {
		//记录日志的server port
		c.Locals("serverPort", config.AppConfig.Server.Port)
		//获取请求的body，屏蔽掉空格
		body := strings.NewReplacer(" ", "", "　", "", "\t", "", "\n", "", "\r", "").Replace(string(c.Body()))
		if body != "" {
			c.Locals("body", body)
		} else {
			c.Locals("body", "{}")
		}
		return c.Next()
	}
}
