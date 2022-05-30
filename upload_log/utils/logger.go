package utils

// 简单的封装日志相关类

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	file     *os.File
	fileName string
}

func (l *Logger) Write(b []byte) (n int, err error) {
	if l.file == nil {
		// 新建一个file连接
		l.createNewFile()
	}
	if l.fileName != nowFileName() {
		l.createNewFile()
	}
	n, err = l.file.Write(b)
	return n, err
}

// 获取最新的文件名
func nowFileName() string {
	// ioc, _ := time.LoadLocation("Asia/Shanghai")
	// timeDir := fmt.Sprint(time.Now().In(ioc).Format("2006/01/02/15")) // 2006-01-02 15:04:05
	timeDir := fmt.Sprint(time.Now().Format("2006/01/02/15")) // 2006-01-02 15:04:05
	path := "./log/" + timeDir + ".log"
	return path
}

// 新建文件
func (l *Logger) createNewFile() error {
	fileName := nowFileName()

	dir, _ := filepath.Split(fileName) // 识别目录与文件
	os.MkdirAll(dir, os.ModePerm)      // 创建多层目录，如果存在不会报错

	// 打开该文件，如果不存在则创建
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// 打开失败，尝试创建
		fmt.Println("打开日志文件失败")
		return err
	}
	// 关闭原来的文件
	if l.file != nil {
		l.closeFile()
	}
	l.file = file
	l.fileName = fileName
	return nil
}

// 关闭文件
func (l *Logger) closeFile() error {
	l.file.Sync()
	return l.file.Close()
}
