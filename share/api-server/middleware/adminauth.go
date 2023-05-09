package middleware

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"share/common/pkg/appError"
	"share/common/pkg/constant"
	"share/common/pkg/goRedis"
	"share/common/service/admin"
)

// AdminAuth 管理端登录校验
func AdminAuth() fiber.Handler {
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
		uid := goRedis.Redis.Get(context.Background(), goRedis.GetKey(fmt.Sprintf("%s:%s", constant.AdminAccessToken, accessToken))).Val()
		if uid == "" {
			return c.JSON(fiber.Map{
				"code": appError.TokenError.Code,
				"msg":  appError.TokenError.Msg,
				"data": nil,
			})
		}

		//设置到上下文中
		c.Locals("UID", uid)
		c.Locals("AccessToken", accessToken)

		return c.Next()
	}
}

// PermissionsAuth  权限校验
func PermissionsAuth() fiber.Handler {
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
		uid := goRedis.Redis.Get(context.Background(), goRedis.GetKey(fmt.Sprintf("%s:%s", constant.AdminAccessToken, accessToken))).Val()
		if uid == "" {
			return c.JSON(fiber.Map{
				"code": appError.TokenError.Code,
				"msg":  appError.TokenError.Msg,
				"data": nil,
			})
		}

		//设置到上下文中
		c.Locals("UID", uid)
		c.Locals("AccessToken", accessToken)

		//超级管理员不校验权限
		if cast.ToInt(uid) == 1 {
			return c.Next()
		}

		//校验Api权限
		urlPath := c.Path()
		//fmt.Println("uid:", uid)
		routes := admin.AdminUser(c).GetUserPermissionsRoutes(cast.ToInt(uid))
		if !lo.Contains[string](routes, urlPath) {
			return c.JSON(fiber.Map{
				"code": appError.PermissionsDenied.Code,
				"msg":  appError.PermissionsDenied.Msg,
				"data": nil,
			})
		}

		return c.Next()
	}
}
