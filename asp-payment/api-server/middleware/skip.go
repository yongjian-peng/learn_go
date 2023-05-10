package middleware

import (
	"asp-payment/common/pkg/appError"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/skip"
)

func Skip() fiber.Handler {
	return skip.New(func(ctx *fiber.Ctx) error {
		//如果是ajax option请求，则跳过
		return ctx.JSON(fiber.Map{
			"code":    appError.SUCCESS.Code,
			"message": appError.SUCCESS.Message,
			"data":    nil,
		})
	}, func(c *fiber.Ctx) bool {
		//非option请求，进行next处理
		return c.Method() != fiber.MethodOptions
	})
}
