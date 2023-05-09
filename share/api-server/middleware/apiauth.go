package middleware

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"share/common/pkg/appError"
	"share/common/pkg/goRedis"
)

// ApiAuth 接口校验
func ApiAuth() fiber.Handler {
	// Return new handler
	return func(c *fiber.Ctx) error {

		//获取access—token
		accessToken := c.Get("x-access-token")
		if accessToken == "" {
			return c.JSON(fiber.Map{
				"code": appError.TokenError.Code,
				"msg":  appError.TokenError.Msg,
				"data": nil,
			})
		}

		//校验token
		ctx := context.Background()
		tokenKey := goRedis.GetKey(fmt.Sprintf("api_access_token:%s", accessToken))
		uid := goRedis.Redis.Get(ctx, tokenKey).Val()
		if uid == "" {
			return c.JSON(fiber.Map{
				"code": appError.TokenError.Code,
				"msg":  appError.TokenError.Msg,
				"data": nil,
			})
		}

		//获取key的TTL，延迟有效期
		fmt.Println("token TTL:")
		fmt.Println(goRedis.Redis.TTL(ctx, tokenKey).Val().Seconds())

		//设置到上下文中
		c.Locals("UID", uid)

		return c.Next()
	}
}
