package service

import (
	"context"
	"upload_log/service/system"
)

func UploadLogOss(ctx context.Context) {
	// 读取文件地址 写入到数据库 状态为待上传
	// set := system.ReadFileService{}
	// set.ReadFile("./")

	// 执行上传到 OSS 平台
	// 读取数据库 状态为待上传的数据 依次上传到 OSS
	// upload := system.UploadFileLogService{}

	// upload.UploadFileLog()

	// 删除已经上传好的文件列表
	// 读取数据库 状态为上传完成的数据 依次删除对应的文件
	clearFile := system.ClearFileService{}

	clearFile.ClearFile()
}
