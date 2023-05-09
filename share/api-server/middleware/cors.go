package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Cors 跨域配置
func Cors() cors.Config {
	return cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With",
		AllowCredentials: true,
		ExposeHeaders:    "",
		MaxAge:           12 * 60 * 60,
	}
}
