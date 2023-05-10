package middleware

import (
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/goRedis"
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Auth() fiber.Handler {
	// Return new handler
	return func(c *fiber.Ctx) error {

		//获取jwt token
		token := c.Get("Authorization")
		//验证前端传过来的token格式，不为空，开头为Bearer
		if len(token) >= 7 && strings.HasPrefix(token, "Bearer ") {
			token = token[7:] //截取字符
		}

		if token == "" {
			return c.JSON(fiber.Map{
				"code":    appError.TokenError.Code,
				"message": appError.TokenError.Message,
				"data":    nil,
			})
		}

		//校验token
		uid := goRedis.Redis.Get(context.Background(), goRedis.GetKey(fmt.Sprintf("auth:%s", token))).Val()
		if uid == "" {
			return c.JSON(fiber.Map{
				"code":    appError.TokenError.Code,
				"message": appError.TokenError.Message,
				"data":    nil,
			})
		}

		//设置到上下文中
		c.Locals("UID", uid)

		return c.Next()
	}
}
