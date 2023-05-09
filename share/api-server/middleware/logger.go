package middleware

import (
	"fmt"
	"share/common/pkg/config"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 配置
func Logger() logger.Config {
	return logger.Config{
		TimeFormat: "2006-01-02 15:04:05.000000",
		Format: `{"time":"${time}","serverPort":"${locals:serverPort}","ip":"${ip}","requestid":"${locals:requestid}","url":"${url}","method":"${method}","body":${locals:body},"requestHeader":${locals:requestHeader},"responseBody":${locals:responseBody},"status":${status},"latency":"${latency}"}
`,
		Output: &lumberjack.Logger{
			Filename:   fmt.Sprintf("%s%s.log", "./logs/", config.AppFileName), // 日志文件路径
			MaxSize:    10,                                                     // 最大M
			MaxBackups: 5,                                                      // 最多保留多少个备份
			MaxAge:     24,                                                     // days
			Compress:   false,                                                  // 是否压缩 disabled by default
		},
	}
}
