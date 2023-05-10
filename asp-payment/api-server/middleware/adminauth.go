package middleware

import "C"
import (
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/goutils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AdminAuth() fiber.Handler {
	// Return new handler
	return func(c *fiber.Ctx) error {
		// 获取sign
		signature := c.Get("signature")
		params := make(map[string]interface{})
		// 获取post all
		_ = c.BodyParser(&params)
		params["sign"] = strings.TrimSpace(signature)
		params["Timestamp"] = c.Get("Timestamp")
		paySecret := config.AppConfig.Server.AdminSecret
		// 验证签名
		if !goutils.HmacSHA256Verify(params, paySecret) {
			return appError.UnauthenticatedErrCode // 签名错误
		}
		return c.Next()
	}
}
