package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"share/common/pkg/appError"
	"share/common/pkg/goutils"
	"share/common/service"
)

// SecretAuth 接口校验
func SecretAuth() fiber.Handler {
	// Return new handler
	return func(c *fiber.Ctx) error {

		//获取AppId
		appId := c.Get("AppId")
		if appId == "" {
			return c.JSON(fiber.Map{
				"code": appError.HeadParamsError.Code,
				"msg":  appError.HeadParamsError.Msg,
				"data": nil,
			})
		}

		//获取Signature
		signature := c.Get("Signature")
		if signature == "" {
			return c.JSON(fiber.Map{
				"code": appError.HeadParamsError.Code,
				"msg":  appError.HeadParamsError.Msg,
				"data": nil,
			})
		}

		//校验签名，查找从缓存获取应用信息是否存在
		appInfo, _ := service.App(c).GetCacheAppInfo(cast.ToInt(appId))
		if appInfo == nil {
			return c.JSON(fiber.Map{
				"code": appError.AppIdNoFound.Code,
				"msg":  appError.AppIdNoFound.Msg,
				"data": nil,
			})
		}

		//获取所有的POST参数
		params := make(map[string]any)
		_ = c.BodyParser(&params)
		params["sign"] = signature
		if !goutils.HmacSHA256Verify(params, appInfo.Secret) {
			return appError.Unauthenticated // 签名错误
		}

		//设置到上下文中
		c.Locals("AppSalt", appInfo.Salt)

		return c.Next()
	}
}
