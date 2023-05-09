package router

import (
	"fmt"
	"share/api-server/middleware"
	"share/common/pkg/config"
	"share/common/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Api(app *fiber.App) {

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code": 0,
			"msg":  "Success",
			"data": "from ip:" + c.IP(),
		})
	})

	app.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code": 0,
			"msg":  "Success",
			"data": fmt.Sprintf("time:%s port:%d version : %s", time.Now().Format("2006-01-02 15:04:05.000000"), c.Locals("serverPort").(int), config.Version),
		})
	})

	api := app.Group("/api")
	{
		// 添加路由
		user := api.Group("user")
		{
			user.Post("/register", middleware.SecretAuth(), func(ctx *fiber.Ctx) error {
				return service.User(ctx).Register()
			})

			user.Post("/login", middleware.SecretAuth(), func(ctx *fiber.Ctx) error {
				return service.User(ctx).Login()
			})

			user.Post("/info", middleware.ApiAuth(), func(ctx *fiber.Ctx) error {
				return service.User(ctx).Info()
			})
		}
	}

}
