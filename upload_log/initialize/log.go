package initialize

import (
	"log"
	"upload_log/utils"
)

func InitLog() {
	// 初始化日志
	log.SetOutput(&utils.Logger{})
	log.SetFlags(log.LstdFlags | log.Llongfile | log.Lmicroseconds) // log.Lshortfile  | log.LUTC
	log.SetPrefix("[BOP] ")
}
