package mylogger

import (
	"io"
	"os"
	"time"
)

const (
	// LOGPATH LOGPATH/time.Now().FORMAT(FORMAT)/*.log
	LOGPATH = "log/"
	// FORMAT .
	FORMAT = "20060102"
	// LineFeed 换行
	LineFeed = "\r\n"
)

var path = LOGPATH + time.Now().Format(FORMAT) + "/"

// WriteLog return error

func WriteLog(filename, msg string) error {
	if !IsExist(path) {
		return CreateDir(path)
	}
	var (
		err error
		f   *os.File
	)
	f, err = os.OpenFile(path+filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	_, err = io.WriteString(f, LineFeed+"[log:]"+time.Now().Format("2006-01-02 15:04:05")+" "+msg)
	// logger = log.New(logFile, "[qSkiptool]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
	defer f.Close()
	return err
}

// CreateDir 文件夹创建
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

// IsExist 判断文件夹、文件是否存在 存在返回 true
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
