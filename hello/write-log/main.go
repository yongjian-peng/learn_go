package main

import (
	"log"
	"os"
	"time"
)

var logger *log.Logger

func init() {
	file := "./" + time.Now().Format("20220419") + ".txt"

	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	logger = log.New(logFile, "[qSkiptool]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出

	return
}

func main() {
	logger.Println("Hello:GO")
	return
}
