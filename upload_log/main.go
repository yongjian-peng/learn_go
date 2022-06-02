package main

import (
	"context"
	"fmt"
	"upload_log/core"
	"upload_log/global"
	"upload_log/initialize"
	"upload_log/service"
)

func main() {

	fmt.Println("log")
	ctx := context.Background()
	// 初始化配置
	// 配置文件 yaml 链接数据库 封装工具 使用 vip
	// 初始化Viper这个是获取配置
	global.Viper = core.Viper()
	// 读取配置
	// 初始化Redis
	// initialize.Redis()

	// 初始化Gorm

	global.Gorm = initialize.Gorm()
	// using standard library "flag" package
	// global.Log2()
	// 初始化 数据库

	service.UploadLogOss(ctx)

	// 读取文件

	// 删除文件
}
